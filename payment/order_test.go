// +build integration

package payment

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/asaskevich/govalidator"
	macarooncmd "github.com/brave-intl/bat-go/cmd/macaroon"
	appctx "github.com/brave-intl/bat-go/utils/context"
	"github.com/brave-intl/bat-go/utils/cryptography"
	"github.com/stretchr/testify/suite"
)

type OrderTestSuite struct {
	service *Service
	suite.Suite
}

func TestOrderTestSuite(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}

func (suite *OrderTestSuite) SetupSuite() {
	govalidator.SetFieldsRequiredByDefault(true)
	pg, err := NewPostgres("", false, "")
	suite.Require().NoError(err, "Failed to get postgres conn")

	m, err := pg.NewMigrate()
	suite.Require().NoError(err, "Failed to create migrate instance")

	ver, dirty, _ := m.Version()
	if dirty {
		suite.Require().NoError(m.Force(int(ver)))
	}
	if ver > 0 {
		suite.Require().NoError(m.Down(), "Failed to migrate down cleanly")
	}

	EncryptionKey = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0"
	InitEncryptionKeys()

	suite.Require().NoError(pg.Migrate(), "Failed to fully migrate")
	suite.service = &Service{
		Datastore: pg,
	}
}

func (suite *OrderTestSuite) TearDownTest() {
	suite.CleanDB()
}

func (suite *OrderTestSuite) CleanDB() {
	tables := []string{"api_keys"}

	pg, err := NewPostgres("", false, "")
	suite.Require().NoError(err, "Failed to get postgres conn")

	for _, table := range tables {
		_, err = pg.RawDB().Exec("delete from " + table)
		suite.Require().NoError(err, "Failed to get clean table")
	}
}

func (suite *OrderTestSuite) TestCreateOrderItemFromMacaroon() {
	// encrypt merchant key
	cipher, nonce, err := cryptography.EncryptMessage(byteEncryptionKey, []byte("testing123"))
	suite.Require().NoError(err)

	// create key in db for our brave.com location
	_, err = suite.service.Datastore.CreateKey("brave.com", "brave.com", hex.EncodeToString(cipher), hex.EncodeToString(nonce[:]))
	suite.Require().NoError(err)

	c := macarooncmd.Caveats{
		"sku":                     "sku",
		"price":                   "5.01",
		"description":             "coffee",
		"currency":                "usd",
		"credential_type":         "time_bound",
		"allowed_payment_methods": "stripe",
		"metadata": `
				{
					"stripe_product_id":"stripe_product_id",
					"stripe_success_url":"stripe_success_url",
					"stripe_cancel_url":"stripe_cancel_url"
				}
			`,
	}

	// create sku using key
	t := macarooncmd.Token{
		ID: "id", Version: 2, Location: "brave.com",
		FirstPartyCaveats: []macarooncmd.Caveats{c},
	}

	sku, err := t.Generate("testing123")
	suite.Require().NoError(err)

	// hacky add to skuMap
	skuMap["development"][sku] = true

	ctx := context.WithValue(context.Background(), appctx.EnvironmentCTXKey, "development")

	orderItem, apm, err := suite.service.CreateOrderItemFromMacaroon(ctx, sku, 1)
	suite.Require().NoError(err)

	suite.Assert().Equal("stripe", strings.Join(*apm, ","))
	suite.Assert().Equal("usd", orderItem.Currency)
	suite.Assert().Equal("sku", orderItem.SKU)
	suite.Assert().Equal("5.01", orderItem.Price.String())
	suite.Assert().Equal("coffee", orderItem.Description.String)
	suite.Assert().Equal("brave.com", orderItem.Location.String)

	badsku, err := t.Generate("321testing")
	suite.Require().NoError(err)

	ctx = context.WithValue(context.Background(), appctx.EnvironmentCTXKey, "development")
	_, _, err = suite.service.CreateOrderItemFromMacaroon(ctx, badsku, 1)
	suite.Require().Equal(err.Error(), "Invalid SKU Token provided in request")
}
