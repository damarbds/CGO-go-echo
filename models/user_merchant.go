package models

import "time"

type UserMerchant struct {
	Id                   string     `json:"id" validate:"required"`
	CreatedBy            string     `json:"created_by" validate:"required"`
	CreatedDate          time.Time  `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int        `json:"is_deleted" validate:"required"`
	IsActive             int        `json:"is_active" validate:"required"`
	FullName 			string  	`json:"full_name"`
	Email 				string 		`json:"email"`
	PhoneNumber 		string		`json:"phone_number"`
	MerchantId 			string		`json:"merchant_id"`
}

type NewCommandUserMerchant struct {
	Id                   string     `json:"id"`
	FullName 			string		`json:"full_name"`
	Email 				string 		`json:"email"`
	Password 			string		`json:"password"`
	PhoneNumber 		string		`json:"phone_number"`
	MerchantId 			string		`json:"merchant_id"`
}
type UserMerchantDto struct {
	Id                   string     `json:"id" validate:"required"`
	FullName 			string		`json:"full_name"`
	Email 				string 		`json:"email"`
	Password 			string		`json:"password"`
	PhoneNumber 		string		`json:"phone_number"`
	MerchantId 			string		`json:"merchant_id"`
}

type UserMerchantInfoDto struct {
	Id                   string     `json:"id" validate:"required"`
	CreatedDate 		time.Time	`json:"created_date"`
	UpdatedDate 		*time.Time 	`json:"update_date"`
	IsActive 			int 		`json:"is_active"`
	FullName 			string		`json:"full_name"`
	Email 				string 		`json:"email"`
	PhoneNumber 		string		`json:"phone_number"`
	MerchantId 			string		`json:"merchant_id"`
}

type UserMerchantWithPagination struct {
	Data []*UserMerchantInfoDto `json:"data"`
	Meta *MetaPagination    `json:"meta"`
}