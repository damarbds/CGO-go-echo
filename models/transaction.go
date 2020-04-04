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
	PromoId             string     `json:"promo_id"`
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
	BookingType         int        `json:"booking_type,omitempty"`
	BookingExpId        string     `json:"booking_exp_id"`
	PromoId             string     `json:"promo_id"`
	PaymentMethodId     string     `json:"payment_method_id"`
	ExperiencePaymentId string     `json:"experience_payment_id"`
	Status              int        `json:"status,omitempty"`
	TotalPrice          float64    `json:"total_price,omitempty"`
	Currency            string     `json:"currency"`
}