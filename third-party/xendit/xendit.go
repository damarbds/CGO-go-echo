package xendit

import (
	"context"
	"time"

	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/card"
	"github.com/xendit/xendit-go/client"
	"github.com/xendit/xendit-go/virtualaccount"
)

var XenClient *client.API

const (
	defXenditSecretKey = "xnd_development_Dms6iAkgd6b4p5f9jpLdP41uaCVBdCLPNqJ00XDiFQL0oIpsTZYVLlERGFnxi"

	// Production
	defXenditSecretKeyProd = "xnd_production_wUqt0xBrasJpktiTTgOgOIojpewhY455AGFik0AxizdVAL1pIUYBic8EGeStyDs"
)

func XenditSetup() {
	XenClient = client.New(defXenditSecretKey)
}

type VirtualAccount struct {
	*virtualaccount.Client
	ExternalID string
	BankCode   string
	Name       string
	ExpireDate *time.Time
}

type VACallbackRequest struct {
	ID                       string `json:"id"`
	MerchantCode             string `json:"merchant_code"`
	TransactionTimestamp     string `json:"transaction_timestamp"`
	Amount                   string `json:"amount"`
	BankCode                 string `json:"bank_code"`
	AccountNumber            string `json:"account_number"`
	ExternalID               string `json:"external_id"`
	OwnerID                  string `json:"owner_id"`
	CallbackVirtualAccountID string `json:"callback_virtual_account_id"`
	PaymentID                string `json:"payment_id"`
	Created                  string `json:"created"`
	Updated                  string `json:"updated"`
}

func (va *VirtualAccount) CreateFixedVA(ctx context.Context) (*xendit.VirtualAccount, error) {
	data := &virtualaccount.CreateFixedVAParams{
		ExternalID:     va.ExternalID,
		BankCode:       va.BankCode,
		Name:           va.Name,
		ExpirationDate: va.ExpireDate,
	}

	resVa, err := va.CreateFixedVAWithContext(ctx, data)
	if err != nil {
		return nil, err
	}

	return resVa, nil
}

type CreditCard struct {
	*card.Client
	TokenID    string
	AuthID     string
	ExternalID string
	Amount     float64
	IsCapture  bool
}

func (cc *CreditCard) CreateCharge(ctx context.Context) (*xendit.CardCharge, error) {
	data := &card.CreateChargeParams{
		TokenID:          cc.TokenID,
		AuthenticationID: cc.AuthID,
		ExternalID:       cc.ExternalID,
		Amount:           cc.Amount,
		Capture:          &cc.IsCapture,
	}

	resCc, err := cc.CreateChargeWithContext(ctx, data)
	if err != nil {
		return nil, err
	}

	return resCc, nil
}
