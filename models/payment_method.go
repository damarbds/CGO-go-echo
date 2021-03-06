package models

import (
	"database/sql"
	"time"
)

type PaymentMethod struct {
	Id           string         `json:"id" validate:"required"`
	CreatedBy    string         `json:"created_by" validate:"required"`
	CreatedDate  time.Time      `json:"created_date" validate:"required"`
	ModifiedBy   *string        `json:"modified_by"`
	ModifiedDate *time.Time     `json:"modified_date"`
	DeletedBy    *string        `json:"deleted_by"`
	DeletedDate  *time.Time     `json:"deleted_date"`
	IsDeleted    int            `json:"is_deleted" validate:"required"`
	IsActive     int            `json:"is_active" validate:"required"`
	Name         string         `json:"name"`
	Type         int            `json:"type"`
	Desc         sql.NullString `json:"desc"`
	Icon         string         `json:"icon"`
	MidtransPaymentCode *string	`json:"midtrans_payment_code"`
}

type PaymentMethodObject struct {
	Id   string `json:"id" validate:"required"`
	Name string `json:"name"`
	Type int    `json:"type"`
	Desc string `json:"desc"`
	Icon string `json:"icon"`
	MidtransPaymentCode *string	`json:"midtrans_payment_code"`
}
