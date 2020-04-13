package models

import "time"

type ExperiencePayment struct {
	Id               string     `json:"id" validate:"required"`
	CreatedBy        string     `json:"created_by":"required"`
	CreatedDate      time.Time  `json:"created_date" validate:"required"`
	ModifiedBy       *string    `json:"modified_by"`
	ModifiedDate     *time.Time `json:"modified_date"`
	DeletedBy        *string    `json:"deleted_by"`
	DeletedDate      *time.Time `json:"deleted_date"`
	IsDeleted        int        `json:"is_deleted" validate:"required"`
	IsActive         int        `json:"is_active" validate:"required"`
	ExpPaymentTypeId string     `json:"exp_payment_type_id"`
	ExpId            string     `json:"exp_id"`
	PriceItemType    int        `json:"price_item_type"`
	Currency         int        `json:"currency"`
	Price            float64    `json:"price"`
	CustomPrice      *float64   `json:"custom_price"`
}

type ExperiencePaymentJoinType struct {
	Id                 string     `json:"id" validate:"required"`
	CreatedBy          string     `json:"created_by":"required"`
	CreatedDate        time.Time  `json:"created_date" validate:"required"`
	ModifiedBy         *string    `json:"modified_by"`
	ModifiedDate       *time.Time `json:"modified_date"`
	DeletedBy          *string    `json:"deleted_by"`
	DeletedDate        *time.Time `json:"deleted_date"`
	IsDeleted          int        `json:"is_deleted" validate:"required"`
	IsActive           int        `json:"is_active" validate:"required"`
	ExpPaymentTypeId   string     `json:"exp_payment_type_id"`
	ExpId              string     `json:"exp_id"`
	PriceItemType      int        `json:"price_item_type"`
	Currency           int        `json:"currency"`
	Price              float64    `json:"price"`
	CustomPrice        *float64   `json:"custom_price"`
	ExpPaymentTypeName string     `json:"exp_payment_type_name"`
	ExpPaymentTypeDesc string     `json:"exp_payment_type_desc"`
}
