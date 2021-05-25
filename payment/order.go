package payment

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/brave-intl/bat-go/utils/datastore"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	stripe "github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/checkout/session"
	"github.com/stripe/stripe-go/v71/customer"
	macaroon "gopkg.in/macaroon.v2"
)

//StripePaymentMethod - the label for stripe payment method
const StripePaymentMethod = "stripe"

// Order includes information about a particular order
type Order struct {
	ID         uuid.UUID            `json:"id" db:"id"`
	CreatedAt  time.Time            `json:"createdAt" db:"created_at"`
	Currency   string               `json:"currency" db:"currency"`
	UpdatedAt  time.Time            `json:"updatedAt" db:"updated_at"`
	TotalPrice decimal.Decimal      `json:"totalPrice" db:"total_price"`
	MerchantID string               `json:"merchantId" db:"merchant_id"`
	Location   datastore.NullString `json:"location" db:"location"`
	Status     string               `json:"status" db:"status"`
	Items      []OrderItem          `json:"items"`
	Metadata   Metadata             `json:"metadata" db:"metadata"`
}

// OrderItem includes information about a particular order item
type OrderItem struct {
	ID             uuid.UUID            `json:"id" db:"id"`
	OrderID        uuid.UUID            `json:"orderId" db:"order_id"`
	SKU            string               `json:"sku" db:"sku"`
	CreatedAt      *time.Time           `json:"createdAt" db:"created_at"`
	UpdatedAt      *time.Time           `json:"updatedAt" db:"updated_at"`
	Currency       string               `json:"currency" db:"currency"`
	Quantity       int                  `json:"quantity" db:"quantity"`
	Price          decimal.Decimal      `json:"price" db:"price"`
	Subtotal       decimal.Decimal      `json:"subtotal" db:"subtotal"`
	Location       datastore.NullString `json:"location" db:"location"`
	Description    datastore.NullString `json:"description" db:"description"`
	CredentialType string               `json:"credentialType" db:"credential_type"`
	PaymentMethods Methods              `json:"paymentMethods" db:"payment_methods"`
	Metadata       Metadata             `json:"metadata" db:"metadata"`
}

// Methods type is a string slice holding payments
type Methods []string

// Scan the src sql type into the passed JSONStringArray
func (pm *Methods) Scan(src interface{}) error {
	var x []sql.NullString
	var v = pq.Array(&x)

	if err := v.Scan(src); err != nil {
		return err
	}
	for i := 0; i < len(x); i++ {
		if x[i].Valid {
			*pm = append(*pm, x[i].String)
		}
	}

	return nil
}

// Value the driver.Value representation
func (pm *Methods) Value() (driver.Value, error) {
	return pq.Array(pm), nil
}

// CreateOrderItemFromMacaroon creates an order item from a macaroon
func CreateOrderItemFromMacaroon(sku string, quantity int) (*OrderItem, error) {
	macBytes, err := macaroon.Base64Decode([]byte(sku))
	if err != nil {
		return nil, err
	}
	mac := &macaroon.Macaroon{}
	err = mac.UnmarshalBinary(macBytes)
	if err != nil {
		return nil, err
	}

	// TODO Figure out how to verify macaroon using library
	// I think we have to call .Bind()
	// https://github.com/go-macaroon/macaroon#func-macaroon-bind
	// I think we simply want to verify the signature and not the caveats?
	// SO maybe VerifySignature
	// https://github.com/go-macaroon/macaroon#func-macaroon-verifysignature

	caveats := mac.Caveats()
	orderItem := OrderItem{}
	orderItem.Quantity = quantity
	orderItem.Location.String = mac.Location()
	orderItem.Location.Valid = true

	for i := 0; i < len(caveats); i++ {
		caveat := mac.Caveats()[i]
		values := strings.Split(string(caveat.Id), "=")
		key := strings.TrimSpace(values[0])
		value := strings.TrimSpace(values[1])

		switch key {
		case "sku":
			orderItem.SKU = value
		case "price", "amount":
			orderItem.Price, err = decimal.NewFromString(value)
			if err != nil {
				return nil, err
			}
		case "description":
			orderItem.Description.String = value
			orderItem.Description.Valid = true
			if err != nil {
				return nil, err
			}
		case "currency":
			orderItem.Currency = value
		case "credential_type":
			orderItem.CredentialType = value
		case "payment_methods":
			orderItem.PaymentMethods = strings.Split(value, ",")
		case "metadata":
			err := json.Unmarshal([]byte(value), &orderItem.Metadata)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal macaroon metadata: %w", err)
			}
		}
	}
	newQuantity, err := decimal.NewFromString(strconv.Itoa(orderItem.Quantity))
	if err != nil {
		return nil, err
	}

	orderItem.Subtotal = orderItem.Price.Mul(newQuantity)

	return &orderItem, nil
}

// IsStripePayable returns true if every item is payable by Stripe
func (order Order) IsStripePayable() bool {
	for _, item := range order.Items {

		// check stripe in payment
		if !strings.Contains(strings.Join(item.PaymentMethods, ","), StripePaymentMethod) {
			return false
		}
		// TODO: make sure we have a stripe_product_id caveat
		// TODO: if not we need to look into subscription trials:
		/// -> https://stripe.com/docs/billing/subscriptions/trials
	}
	return true
}

// CreateCheckoutSessionResponse - the structure of a checkout session response
type CreateCheckoutSessionResponse struct {
	SessionID string `json:"checkoutSessionId"`
}

// CreateStripeCheckoutSession - Create a Stripe Checkout Session for an Order
func (order Order) CreateStripeCheckoutSession(email, successURI, cancelURI string) (CreateCheckoutSessionResponse, error) {
	// Create customer if not already created
	i := customer.List(&stripe.CustomerListParams{
		Email: stripe.String(email),
	})

	matchingCustomers := 0
	for i.Next() {
		matchingCustomers++
	}

	var customerID string
	if matchingCustomers > 0 {
		customerID = i.Customer().ID
	} else {
		customer, err := customer.New(&stripe.CustomerParams{
			Email: stripe.String(email),
		})
		if err != nil {
			return CreateCheckoutSessionResponse{}, fmt.Errorf("failed to create stripe customer: %w", err)
		}
		customerID = customer.ID
	}

	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(customerID),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode:              stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL:        stripe.String(successURI),
		CancelURL:         stripe.String(cancelURI),
		ClientReferenceID: stripe.String(order.ID.String()),
		SubscriptionData:  &stripe.CheckoutSessionSubscriptionDataParams{},
		LineItems:         order.CreateStripeLineItems(),
	}

	params.SubscriptionData.AddMetadata("orderID", order.ID.String())

	session, err := session.New(params)
	if err != nil {
		return CreateCheckoutSessionResponse{}, fmt.Errorf("failed to create stripe session: %w", err)
	}

	data := CreateCheckoutSessionResponse{
		SessionID: session.ID,
	}
	return data, nil
}

// CreateStripeLineItems - create line items for a checkout session with stripe
func (order Order) CreateStripeLineItems() []*stripe.CheckoutSessionLineItemParams {
	lineItems := make([]*stripe.CheckoutSessionLineItemParams, len(order.Items))
	for index, item := range order.Items {
		// since we are creating stripe line item, we can assume
		// that the stripe product is embedded in macaroon as metadata
		lineItems[index] = &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(item.Metadata["stripe_product_id"]),
			Quantity: stripe.Int64(int64(item.Quantity)),
		}
	}
	return lineItems
}

// IsPaid returns true if the order is paid
func (order Order) IsPaid() bool {
	return order.Status == "paid"
}
