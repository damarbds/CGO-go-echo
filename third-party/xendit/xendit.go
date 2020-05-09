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
	defXenditSecretKey           = "xnd_development_Dms6iAkgd6b4p5f9jpLdP41uaCVBdCLPNqJ00XDiFQL0oIpsTZYVLlERGFnxi"
	defXenditVerifyTokenCallback = "41a617f2a5fed878b48570432f3d5e1d68f49666e2fa22270a9c694532fba50a"
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
	ExternalID string
	Amount     float64
}

func (cc *CreditCard) CreateCharge(ctx context.Context) (*xendit.CardCharge, error) {
	data := &card.CreateChargeParams{
		TokenID:          "",
		ExternalID:       "",
		Amount:           0,
		AuthenticationID: "",
		CardCVN:          "",
		Capture:          nil,
		Currency:         "",
		IsRecurring:      nil,
	}

	resCc, err := cc.CreateChargeWithContext(ctx, data)
	if err != nil {
		return nil, err
	}

	return resCc, nil
}
