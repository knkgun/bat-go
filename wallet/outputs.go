package wallet

import (
	"fmt"

	"github.com/brave-intl/bat-go/utils/altcurrency"
	walletutils "github.com/brave-intl/bat-go/utils/wallet"
	uuid "github.com/satori/go.uuid"
)

const (
	// InvalidCurrency - wallet currency is invalid
	InvalidCurrency = "invalid"
	// BATCurrency - wallet currency is BAT
	BATCurrency = "BAT"
	// BTCCurrency - wallet currency is BTC
	BTCCurrency = "BTC"
	// ETHCurrency - wallet currency is ETH
	ETHCurrency = "ETH"
	// LTCCurrency - wallet currency is LTC
	LTCCurrency = "LTC"
)

// BraveProviderDetailsV3 - details about the provider
type BraveProviderDetailsV3 struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UpholdProviderDetailsV3 - details about the provider
type UpholdProviderDetailsV3 struct {
	Name             string `json:"name"`
	LinkingID        string `json:"linkingId"`
	AnonymousAddress string `json:"anonymousAddress"`
}

// ResponseV3 - wallet creation response
type ResponseV3 struct {
	PaymentID              string      `json:"paymentId"`
	DepositAccountProvider interface{} `json:"depositAccountProvider,omitempty"`
	WalletProvider         interface{} `json:"walletProvider,omitempty"`
	AltCurrency            string      `json:"altcurrency"`
	PublicKey              string      `json:"publicKey"`
}

func convertAltCurrency(a *altcurrency.AltCurrency) string {
	if a == nil {
		return BATCurrency
	}
	switch *a {
	case altcurrency.BAT:
		return BATCurrency
	case altcurrency.BTC:
		return BTCCurrency
	case altcurrency.ETH:
		return ETHCurrency
	case altcurrency.LTC:
		return LTCCurrency
	default:
		return InvalidCurrency
	}
}

func responseV3ToInfo(resp ResponseV3) *walletutils.Info {
	alt, _ := altcurrency.FromString(resp.AltCurrency)
	info := walletutils.Info{
		ID:          resp.PaymentID,
		AltCurrency: &alt,
		PublicKey:   resp.PublicKey,
	}
	if up, ok := resp.DepositAccountProvider.(UpholdProviderDetailsV3); ok {
		info.ProviderID = resp.PaymentID
		info.Provider = up.Name
		if up.LinkingID != "" {
			info.ProviderID = string(up.LinkingID)
		}
		if up.AnonymousAddress != "" {
			anonymousAddress := uuid.Must(uuid.FromString(up.AnonymousAddress))
			info.AnonymousAddress = &anonymousAddress
		}
	}
	// bp, ok := resp.WalletProvider.(BraveProviderDetailsV3)
	// fmt.Println(bp)
	if resp.WalletProvider["id"] != "" {
		info.Provider = bp.Name
		if bp.Name == "uphold" {
			info.ProviderID = string(bp.ID)
		} else if bp.Name == "brave" {
			info.ProviderID = string(bp.ID)
		}
		fmt.Println(bp)
	}
	return &info
}

// func responseV3ToInfo(response ResponseV3) walletutils.Info {
// 	alt, _ := altcurrency.FromString(response.AltCurrency)
// 	info := walletutils.Info{
// 		ID:          response.PaymentID,
// 		AltCurrency: &alt,
// 	}
// 	if response.WalletProvider != nil {
// 		switch wp := response.WalletProvider.(type) {
// 		case UpholdProviderDetailsV3:
// 			anonymousAddress := uuid.Must(uuid.FromString(wp.AnonymousAddress))
// 			info.AnonymousAddress = &anonymousAddress
// 			if wp.LinkingID != "" {
// 				providerLinkingID := uuid.Must(uuid.FromString(wp.LinkingID))
// 				info.ProviderLinkingID = &providerLinkingID
// 			}
// 			info.ProviderID = "not"
// 		case BraveProviderDetailsV3:
// 			info.Provider = wp.Name
// 			if wp.Name == "uphold" {
// 				info.ProviderID = wp.ID
// 			}
// 		}
// 	}
// 	return info
// }

func infoToResponseV3(info *walletutils.Info) ResponseV3 {
	var (
		linkingID        string
		anonymousAddress string
		altCurrency      string = convertAltCurrency(info.AltCurrency)
	)
	if info == nil {
		return ResponseV3{}
	}

	if info.ProviderLinkingID == nil {
		linkingID = ""
	} else {
		linkingID = info.ProviderLinkingID.String()
	}

	if info.AnonymousAddress == nil {
		anonymousAddress = ""
	} else {
		anonymousAddress = info.AnonymousAddress.String()
	}

	resp := ResponseV3{
		PaymentID:   info.ID,
		AltCurrency: altCurrency,
		PublicKey:   info.PublicKey,
	}

	// if this is linked to uphold, add the default account provider
	if info.Provider == "uphold" {
		if anonymousAddress == "" {
			// no linked anon card
			resp.WalletProvider = BraveProviderDetailsV3{
				Name: "brave",
				ID:   info.ID,
			}
		} else {
			resp.WalletProvider = BraveProviderDetailsV3{
				Name: "uphold",
				ID:   info.ProviderID,
			}
		}
		if linkingID != "" {
			resp.DepositAccountProvider = UpholdProviderDetailsV3{
				Name:             info.Provider,
				LinkingID:        linkingID,
				AnonymousAddress: anonymousAddress,
			}
		}
	} else if info.Provider == "brave" {
		// no linked anon card
		resp.WalletProvider = BraveProviderDetailsV3{
			Name: "brave",
			ID:   info.ID,
		}
	} else {
		resp.WalletProvider = BraveProviderDetailsV3{
			Name: info.Provider,
			ID:   info.ProviderID,
		}
	}
	return resp
}
