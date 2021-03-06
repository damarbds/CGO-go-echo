package midtrans

import "github.com/veritrans/go-midtrans"

var Midclient midtrans.Client
var CoreGateway midtrans.CoreGateway
var SnapGateway midtrans.SnapGateway

const (
	//dev
	MidtransServerKey   = "SB-Mid-server-GHKQBfmGEt_NfZLQSi_2nFzV"
	MidtransClientKey   = "SB-Mid-client-g6CpjvcnJVhosUrt"
	TransactionEndpoint = "https://app.sandbox.midtrans.com/snap/v1/transactions"
	MidtransAPIEnvType  = midtrans.Sandbox

	// Production
	// MidtransServerKey   = "Mid-server-_ctQOnsKzxl7gY4FyWjboPgx"
	// MidtransClientKey   = "Mid-client-hcjkcpg5p77ob8Mc"
	// TransactionEndpoint = "https://app.midtrans.com/snap/v1/transactions"
	// MidtransAPIEnvType  = midtrans.Production
)

func SetupMidtrans() {
	Midclient = midtrans.NewClient()
	Midclient.ServerKey = MidtransServerKey
	Midclient.ClientKey = MidtransClientKey

	Midclient.APIEnvType = MidtransAPIEnvType

	CoreGateway = midtrans.CoreGateway{
		Client: Midclient,
	}

	SnapGateway = midtrans.SnapGateway{
		Client: Midclient,
	}
}

type MidtransPaymentCallback struct {
	OrderId           string `form:"order_id"`
	StatusCode        string `form:"status_code"`
	TransactionStatus string `form:"transaction_status"`
}

type MidtransCallback struct {
	VaNumber               []VaNumber `json:"va_numbers,omitempty"`
	TransactionTime        string     `json:"transaction_time,omitempty"`
	TransactionStatus      string     `json:"transaction_status,omitempty"`
	TransactionId          string     `json:"transaction_id,omitempty"`
	StatusMessage          string     `json:"status_message,omitempty"`
	StatusCode             string     `json:"status_code,omitempty"`
	SignatureKey           string     `json:"signature_key,omitempty"`
	BillKey                string     `json:"bill_key,omitempty"`
	BillerCode             string     `json:"biller_code,omitempty"`
	PermataVaNumber        string     `json:"permata_va_number,omitempty"`
	PaymentType            string     `json:"payment_type,omitempty"`
	OrderId                string     `json:"order_id,omitempty"`
	MaskedCard             string     `json:"masked_card,omitempty"`
	GrossAmount            string     `json:"gross_amount,omitempty"`
	FraudStatus            string     `json:"fraud_status,omitempty"`
	Currency               string     `json:"currency,omitempty"`
	ChannelResponseMessage string     `json:"channel_response_message,omitempty"`
	ChannelResponseCode    string     `json:"channel_response_code,omitempty"`
	CardType               string     `json:"card_type,omitempty"`
	Bank                   string     `json:"bank,omitempty"`
	ApprovalCode           string     `json:"approval_code,omitempty"`
}
type VaNumber struct {
	Number string `json:"va_number,omitempty"`
	Bank   string `json:"bank,omitempty"`
}
type TransactionDetails struct {
	OrderID     string  `json:"order_id"`
	GrossAmount float64 `json:"gross_amount"`
}

type CustomerDetail struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type OptionColorTheme struct {
	Primary     string `json:"primary"`
	PrimaryDark string `json:"primary_dark"`
	Secondary   string `json:"secondary"`
}

type MidtransCharge struct {
	TransactionDetails TransactionDetails `json:"transaction_details"`
	OptionColorTheme   OptionColorTheme   `json:"option_color_theme"`
	EnablePayment      []string           `json:"enabled_payments"`
	CustomerDetail     CustomerDetail     `json:"customer_details"`
}
