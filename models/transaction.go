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
	BookingExpId        string     `json:"booking_exp_id"`
	PromoId             *string    `json:"promo_id"`
	PaymentMethodId     string     `json:"payment_method_id"`
	ExperiencePaymentId string     `json:"experience_payment_id"`
	Status              int        `json:"status"`
	TotalPrice          float64    `json:"total_price"`
	Currency            string     `json:"currency"`
}

type PaymentTransaction struct {
	Status        int    `json:"status"`
	Message       string `json:"message"`
	TransactionID string `json:"transaction_id"`
}

type TransactionIn struct {
	BookingType         int     `json:"booking_type,omitempty"`
	BookingExpId        string  `json:"booking_exp_id"`
	PromoId             string  `json:"promo_id"`
	PaymentMethodId     string  `json:"payment_method_id"`
	ExperiencePaymentId string  `json:"experience_payment_id"`
	Status              int     `json:"status,omitempty"`
	TotalPrice          float64 `json:"total_price,omitempty"`
	Currency            string  `json:"currency"`
}

type ConfirmPaymentIn struct {
	TransactionID     string `json:"transaction_id"`
	TransactionStatus int    `json:"transaction_status,omitempty"`
	BookingStatus     int    `json:"booking_status,omitempty"`
}

type TransactionOut struct {
	TransactionId     string    `json:"transaction_id"`
	ExpId             string    `json:"exp_id"`
	ExpType           string    `json:"exp_type"`
	ExpTitle          string    `json:"exp_title"`
	BookingExpId      string    `json:"booking_exp_id"`
	BookingCode       string    `json:"booking_code"`
	BookingDate       time.Time `json:"booking_date"`
	CheckInDate       time.Time `json:"check_in_date"`
	BookedBy          string    `json:"booked_by"`
	GuestDesc         string    `json:"guest_desc"`
	Email             string    `json:"email"`
	TransactionStatus int       `json:"transaction_status"`
	BookingStatus     int       `json:"booking_status"`
}

type TransactionDto struct {
	TransactionId string        `json:"transaction_id"`
	ExpId         string        `json:"exp_id"`
	ExpTitle      string        `json:"exp_title"`
	ExpType       []string      `json:"exp_type"`
	BookingExpId  string        `json:"booking_exp_id"`
	BookingCode   string        `json:"booking_code"`
	BookingDate   time.Time     `json:"booking_date"`
	CheckInDate   time.Time     `json:"check_in_date"`
	BookedBy      []BookedByObj `json:"booked_by"`
	Guest         int           `json:"guest"`
	Email         string        `json:"email"`
	Status        string        `json:"status"`
}

type TransactionWithPagination struct {
	Data []*TransactionDto `json:"data"`
	Meta *MetaPagination   `json:"meta"`
}

type TotalTransaction struct {
	TransactionCount      int     `json:"transaction_count"`
	TransactionValueTotal float64 `json:"transaction_value_total"`
}
