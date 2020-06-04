package models

import "time"

type Currency struct {
	Id 		int 		`json:"id"`
	CreatedBy    string         `json:"created_by" validate:"required"`
	CreatedDate  time.Time      `json:"created_date" validate:"required"`
	ModifiedBy   *string        `json:"modified_by"`
	ModifiedDate *time.Time     `json:"modified_date"`
	DeletedBy    *string        `json:"deleted_by"`
	DeletedDate  *time.Time     `json:"deleted_date"`
	IsDeleted    int            `json:"is_deleted" validate:"required"`
	IsActive     int            `json:"is_active" validate:"required"`
	Code 	string		`json:"code"`
	Name 	string		`json:"name"`
	Symbol string		`json:"symbol"`
}
type CurrencyDto struct {
	Id 		int 		`json:"id"`
	Code 	string		`json:"code"`
	Name 	string		`json:"name"`
	Symbol string		`json:"symbol"`
}
type NewCommandCurrency struct {
	Id 		int 		`json:"id"`
	Code 	string		`json:"code"`
	Name 	string		`json:"name"`
	Symbol string		`json:"symbol"`
}
type CurrencyDtoWithPagination struct {
	Data []*CurrencyDto `json:"data"`
	Meta *MetaPagination    `json:"meta"`
}

type CurrencyExChangeRate struct {
	Rates 		Rates `json:"rates"`
	Base 		string 	`json:"base"`
	Date		string 	`json:"date"`
}

type Rates struct {
	IDR 		float64	`json:"IDR"`
	USD 		float64 `json:"USD"`
}