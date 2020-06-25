package models

import "time"

type Transaction struct {
	Id                  string     `json:"id" validate:"required"`
	CreatedBy           string     `json:"created_by" validate:"required"`
	CreatedDate         time.Time  `json:"created_date" validate:"required"`
	ModifiedBy          *string    `json:"modified_by"`
	ModifiedDate        *time.Time `json:"modified_date"`
	DeletedBy           *string    `json:"deleted_by"`
	DeletedDate         *time.Time `json:"deleted_date"`
	IsDeleted           int        `json:"is_deleted" validate:"required"`
	IsActive            int        `json:"is_active" validate:"required"`
	BookingType         int        `json:"booking_type"`
	BookingExpId        *string    `json:"booking_exp_id"`
	PromoId             *string    `json:"promo_id"`
	PaymentMethodId     string     `json:"payment_method_id"`
	ExperiencePaymentId *string    `json:"experience_payment_id"`
	Status              int        `json:"status"`
	TotalPrice          float64    `json:"total_price"`
	Currency            string     `json:"currency"`
	OrderId             *string    `json:"order_id"`
	VaNumber            *string    `json:"va_number"`
	ExChangeRates 		*float64	`json:"ex_change_rates"`
	ExChangeCurrency 	*string		`json:"ex_change_currency"`
}
type TransactionWithBooking struct {
	ExpTitle 		string 		`json:"exp_title"`
	BookedBy 		string 		`json:"booked_by"`
	BookedByEmail 	string		`json:"booked_by_email"`
	BookingDate   	time.Time	`json:"booking_date"`
	TotalPrice 		float64		`json:"total_price"`
	Price 			float64		`json:"price"`
	ExpDuration 	int 		`json:"exp_duration"`
	OrderId			*string		`json:"order_id"`
	MerchantName 	string		`json:"merchant_name"`
	MerchantPhone   string 		`json:"merchant_phone"`
	ExpPaymentDeadlineAmount *int 		`json:"exp_payment_deadline_amount"`
	ExpPaymentDeadlineType 	*string		`json:"exp_payment_deadline_type"`
}
type TransactionWMerchant struct {
	Id                  string     `json:"id" validate:"required"`
	CreatedBy           string     `json:"created_by" validate:"required"`
	CreatedDate         time.Time  `json:"created_date" validate:"required"`
	ModifiedBy          *string    `json:"modified_by"`
	ModifiedDate        *time.Time `json:"modified_date"`
	DeletedBy           *string    `json:"deleted_by"`
	DeletedDate         *time.Time `json:"deleted_date"`
	IsDeleted           int        `json:"is_deleted" validate:"required"`
	IsActive            int        `json:"is_active" validate:"required"`
	BookingType         int        `json:"booking_type"`
	BookingExpId        *string    `json:"booking_exp_id"`
	PromoId             *string    `json:"promo_id"`
	PaymentMethodId     string     `json:"payment_method_id"`
	ExperiencePaymentId *string    `json:"experience_payment_id"`
	Status              int        `json:"status"`
	TotalPrice          float64    `json:"total_price"`
	Currency            string     `json:"currency"`
	OrderId             *string    `json:"order_id"`
	VaNumber            *string    `json:"va_number"`
	ExChangeRates 		*float64	`json:"ex_change_rates"`
	ExChangeCurrency 	*string		`json:"ex_change_currency"`
	MerchantId          string     `json:"merchant_id"`
	OrderIdBook         string     `json:"order_id_book"`
	BookedBy            string     `json:"booked_by"`
	ExpTitle            string     `json:"exp_title"`
	BookingDate         time.Time  `json:"booking_date"`
}

type PaymentTransaction struct {
	Status        int    `json:"status"`
	Message       string `json:"message"`
	TransactionID string `json:"transaction_id"`
}
type PaymentTransactionBookingDate struct {
	Status        int    `json:"status"`
	Message       string `json:"message"`
	BookingDate string `json:"transaction_id"`
}
type TransactionIn struct {
	BookingType         int     `json:"booking_type,omitempty"`
	BookingId           string  `json:"booking_id"`
	OrderId             string  `json:"order_id"`
	PaypalOrderId       string  `json:"paypal_order_id"`
	CcTokenId           string  `json:"cc_token_id"`
	CcAuthId            string  `json:"cc_auth_id"`
	PromoId             string  `json:"promo_id"`
	PaymentMethodId     string  `json:"payment_method_id"`
	ExperiencePaymentId string  `json:"experience_payment_id"`
	Status              int     `json:"status,omitempty"`
	TotalPrice          float64 `json:"total_price,omitempty"`
	Currency            string  `json:"currency"`
	Points              float64 `json:"points"`
	ExChangeRates 		float64	`json:"ex_change_rates"`
	ExChangeCurrency 	string		`json:"ex_change_currency"`
}

type ConfirmPaymentIn struct {
	TransactionID     string `json:"transaction_id"`
	TransactionStatus int    `json:"transaction_status,omitempty"`
	BookingStatus     int    `json:"booking_status,omitempty"`
}
type ConfirmTransactionPayment struct {
	ExpId 	string	`json:"exp_id"`
	TransId string	`json:"trans_id"`
	TransactionStatus int 	`json:"transaction_status"`
	BookingStatus int `json:"booking_status"`
	BookingDate string `json:"booking_date"`
}

type TransactionOut struct {
	TransactionId       string    `json:"transaction_id"`
	ExpId               string    `json:"exp_id"`
	ExpType             string    `json:"exp_type"`
	ExpTitle            string    `json:"exp_title"`
	BookingExpId        string    `json:"booking_exp_id"`
	BookingCode         string    `json:"booking_code"`
	BookingDate         time.Time `json:"booking_date"`
	CheckInDate         time.Time `json:"check_in_date"`
	BookedBy            string    `json:"booked_by"`
	GuestDesc           string    `json:"guest_desc"`
	Email               string    `json:"email"`
	TransactionStatus   int       `json:"transaction_status"`
	BookingStatus       int       `json:"booking_status"`
	TotalPrice          float64   `json:"total_price"`
	ExperiencePaymentId *string   `json:"experience_payment_id"`
	MerchantName        string    `json:"merchant_name"`
	OrderId             *string   `json:"order_id"`
	ExpDuration 		*int 		`json:"exp_duration"`
	ProvinceName 		*string		`json:"province_name"`
	CountryName 		*string 		`json:"country_name"`
}

type TransactionDto struct {
	TransactionId         string                    `json:"transaction_id"`
	ExpId                 string                    `json:"exp_id"`
	ExpTitle              string                    `json:"exp_title"`
	ExpType               []string                  `json:"exp_type"`
	BookingExpId          string                    `json:"booking_exp_id"`
	BookingCode           string                    `json:"booking_code"`
	BookingDate           time.Time                 `json:"booking_date"`
	CheckInDate           time.Time                 `json:"check_in_date"`
	BookedBy              []BookedByObj             `json:"booked_by"`
	Guest                 int                       `json:"guest"`
	Email                 string                    `json:"email"`
	Status                string                    `json:"status"`
	TotalPrice            float64                   `json:"total_price"`
	ExperiencePaymentType *ExperiencePaymentTypeDto `json:"experience_payment_type"`
	Merchant              string                    `json:"merchant"`
	OrderId               *string                   `json:"order_id"`
	GuestCount 			   TotalGuestTransportation 			`json:"guest_count"`
	ExpDuration 			int 		`json:"exp_duration"`
	ProvinceName 			string		`json:"province_name"`
	CountryName 		string			`json:"country_name"`
}

type TransactionWithPagination struct {
	Data []*TransactionDto `json:"data"`
	Meta *MetaPagination   `json:"meta"`
}

type TotalTransaction struct {
	TransactionCount      int     `json:"transaction_count"`
	TransactionValueTotal float64 `json:"transaction_value_total"`
}
