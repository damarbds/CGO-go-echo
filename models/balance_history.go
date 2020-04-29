package models

import "time"

type BalanceHistory struct {
	Id            string     `json:"id" validate:"required"`
	CreatedBy     string     `json:"created_by":"required"`
	CreatedDate   time.Time  `json:"created_date" validate:"required"`
	ModifiedBy    *string    `json:"modified_by"`
	ModifiedDate  *time.Time `json:"modified_date"`
	DeletedBy     *string    `json:"deleted_by"`
	DeletedDate   *time.Time `json:"deleted_date"`
	IsDeleted     int        `json:"is_deleted" validate:"required"`
	IsActive      int        `json:"is_active" validate:"required"`
	Status        int        `json:"status"`
	MerchantId    string     `json:"merchant_id"`
	Amount        float64    `json:"amount"`
	DateOfRequest time.Time  `json:"date_of_request"`
	DateOfPayment time.Time  `json:"date_of_payment"`
	Remarks       string     `json:"remarks"`
}
type NewBalanceHistoryCommand struct {
	Id 			string 		`json:"id"`
	Status 		int		`json:"status"`
	Amount 		float64		`json:"amount"`
	DateOfPayment string	`json:"date_of_payment"`
	Remarks 	string	`json:"remarks"`
}
type NewBalanceHistoryAmountCommand struct {
	Id 			string 		`json:"id"`
	Amount 		float64		`json:"amount"`
}
type NewBalanceHistoryConfirmCommand struct {
	Id 			string 		`json:"id"`
	MerchantId 	string		`json:"merchant_id"`
	Amount 		float64		`json:"amount"`
	Action 		string		`json:"action"`
}
type BalanceHistoryDto struct {
	Id 			string 		`json:"id"`
	MerchantId 	string		`json:"merchant_id"`
	Status 		int		`json:"status"`
	Amount 		float64		`json:"amount"`
	DateOfRequest time.Time	`json:"date_of_request"`
	DateOfPayment time.Time	`json:"date_of_payment"`
	Remarks 	string	`json:"remarks"`
}
type BalanceHistoryDtoWithPagination struct {
	Data []*BalanceHistoryDto `json:"data"`
	Meta *MetaPagination   `json:"meta"`

}
