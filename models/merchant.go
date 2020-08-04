package models

import (
	"time"
)
type NewCommandMerchantRegistrationEmail struct {
	Type string `json:"type"`
	RegistrantName string	`json:"registrant_name"`
	Email 		string	`json:"email"`
	PhoneNumber string	`json:"phone_number"`
	ServiceType string `json:"service_type"`
	TypeOfShip string	`json:"type_of_ship"`
	CompanyName string	`json:"company_name"`
	CompanyAddress string	`json:"company_address"`
	CompanyPhoneNumber string	`json:"company_phone_number"`
}
type  NewCommandContactUs struct {
	FullName 	string	`json:"full_name"`
	Email string	`json:"email"`
	Message string	`json:"message"`
}
type Merchant struct {
	Id            string     `json:"id" validate:"required"`
	CreatedBy     string     `json:"created_by":"required"`
	CreatedDate   time.Time  `json:"created_date" validate:"required"`
	ModifiedBy    *string    `json:"modified_by"`
	ModifiedDate  *time.Time `json:"modified_date"`
	DeletedBy     *string    `json:"deleted_by"`
	DeletedDate   *time.Time `json:"deleted_date"`
	IsDeleted     int        `json:"is_deleted" validate:"required"`
	IsActive      int        `json:"is_active" validate:"required"`
	MerchantName  string     `json:"merchant_name" validate:"required"`
	MerchantDesc  string     `json:"merchant_desc"`
	MerchantEmail string     `json:"merchant_email" validate:"required"`
	Balance       float64    `json:"balance"`
	PhoneNumber   *string    `json:"phone_number"`
	MerchantPicture *string 	`json:"merchant_picture"`
}

type NewCommandMerchant struct {
	Id               string  `json:"id"`
	MerchantName     string  `json:"merchant_name" validate:"required"`
	MerchantDesc     string  `json:"merchant_desc"`
	MerchantEmail    string  `json:"merchant_email" validate:"required"`
	MerchantPassword string  `json:"merchant_password"`
	Balance          float64 `json:"balance"`
	PhoneNumber   string    `json:"phone_number"`
	MerchantPicture *string 	`json:"merchant_picture"`
}

type MerchantInfoDto struct {
	Id            string     `json:"id"`
	UserMerchantId string	`json:"user_merchant_id"`
	CreatedDate   time.Time  `json:"created_date"`
	UpdatedDate   *time.Time `json:"updated_date"`
	IsActive      int        `json:"is_active"`
	MerchantName  string     `json:"merchant_name" validate:"required"`
	MerchantDesc  string     `json:"merchant_desc"`
	MerchantEmail string     `json:"merchant_email" validate:"required"`
	Balance       float64    `json:"balance"`
	PhoneNumber   *string    `json:"phone_number"`
	MerchantPicture *string 	`json:"merchant_picture"`
}

type MerchantDto struct {
	Id            string     `json:"id"`
	CreatedDate   time.Time  `json:"created_date"`
	UpdatedDate   *time.Time `json:"updated_date"`
	IsActive      int        `json:"is_active"`
	MerchantName  string     `json:"merchant_name" validate:"required"`
	MerchantDesc  string     `json:"merchant_desc"`
	MerchantEmail string     `json:"merchant_email" validate:"required"`
	Password 		string 		`json:"password"`
	Balance       float64    `json:"balance"`
	PhoneNumber   *string    `json:"phone_number"`
	MerchantPicture *string 	`json:"merchant_picture"`
}
type MerchantExperience struct {
	Id 				string	`json:"id"`
	MerchantName	string	`json:"merchant_name"`
}
type MerchantTransport struct {
	Id 				string	`json:"id"`
	MerchantName	string	`json:"merchant_name"`
}
type LoginMerchant struct {
	MerchantEmail string `json:"merchant_email"`
	Password      string `json:"password"`
	Type          string `json:"type"`
}

type MerchantWithPagination struct {
	Data []*MerchantInfoDto `json:"data"`
	Meta *MetaPagination    `json:"meta"`
}
