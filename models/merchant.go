package models

import(
	"time"
)

type Merchant struct {
	Id                   string    `json:"id" validate:"required"`
	CreatedBy            string    `json:"created_by":"required"`
	CreatedDate          time.Time `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int       `json:"is_deleted" validate:"required"`
	IsActive             int       `json:"is_active" validate:"required"`
	MerchantName		 string		`json:"merchant_name" validate:"required"`
	MerchantDesc		 string		`json:"merchant_desc"`
	MerchantEmail		 string		`json:"merchant_email" validate:"required"`
	Balance				 float64	`json:"balance"`
}

type NewCommandMerchant struct {
	Id                   string    `json:"id"`
	MerchantName		 string		`json:"merchant_name" validate:"required"`
	MerchantDesc		 string		`json:"merchant_desc"`
	MerchantEmail		 string		`json:"merchant_email" validate:"required"`
	MerchantPassword	 string		`json:"merchant_password"`
	Balance				 float64	`json:"balance"`
}

type MerchantInfoDto struct {
	Id                   string    `json:"id"`
	MerchantName		 string		`json:"merchant_name" validate:"required"`
	MerchantDesc		 string		`json:"merchant_desc"`
	MerchantEmail		 string		`json:"merchant_email" validate:"required"`
	Balance				 float64	`json:"balance"`
}
type LoginMerchant struct {
	MerchantEmail		string		`json:"merchant_email"`
	Password			string		`json:"password"`
	Type 				string		`json:"type"`
}
