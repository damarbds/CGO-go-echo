package models

import "time"

type Promo struct {
	Id           string     `json:"id" validate:"required"`
	CreatedBy    string     `json:"created_by" validate:"required"`
	CreatedDate  time.Time  `json:"created_date" validate:"required"`
	ModifiedBy   *string    `json:"modified_by"`
	ModifiedDate *time.Time `json:"modified_date"`
	DeletedBy    *string    `json:"deleted_by"`
	DeletedDate  *time.Time `json:"deleted_date"`
	IsDeleted    int        `json:"is_deleted" validate:"required"`
	IsActive     int        `json:"is_active" validate:"required"`
	PromoCode    string     `json:"promo_code"`
	PromoName    string     `json:"promo_name"`
	PromoDesc    string     `json:"promo_desc"`
	PromoValue   float64    `json:"promo_value"`
	PromoType    int        `json:"promo_type"`
	PromoImage   string     `json:"promo_image"`
	StartDate 	 *string		`json:"start_date"`
	EndDate 	*string		`json:"end_date"`
	//Currency 	*int 		`json:"currency"`
	MaxUsage   	*int 		`json:"max_usage"`
	ProductionCapacity	*int `json:"production_capacity"`
	CurrencyId *int 	`json:"currency_id"`
	//VoucherValueOptionType *string 	`json:"voucher_value_option_type"`
}
type PromoDto struct {
	Id         string  `json:"id" validate:"required"`
	PromoCode  string  `json:"promo_code"`
	PromoName  string  `json:"promo_name"`
	PromoDesc  string  `json:"promo_desc"`
	PromoValue float64 `json:"promo_value"`
	PromoType  int     `json:"promo_type"`
	PromoImage string  `json:"promo_image"`
	StartDate 	 *string		`json:"start_date"`
	EndDate 	*string		`json:"end_date"`
	Currency 	*int 		`json:"currency"`
	MaxUsage   	*int 		`json:"max_usage"`
	VoucherValueOptionType *string 	`json:"voucher_value_option_type"`
}
type NewCommandPromo struct {
	Id         string  `json:"id"`
	PromoCode  string  `json:"promo_code"`
	PromoName  string  `json:"promo_name"`
	PromoDesc  string  `json:"promo_desc"`
	PromoValue float64 `json:"promo_value"`
	PromoType  int     `json:"promo_type"`
	PromoImage string  `json:"promo_image"`
	StartDate 	 string		`json:"start_date"`
	EndDate 	string		`json:"end_date"`
	Currency 	int 		`json:"currency"`
	MaxUsage   	int 		`json:"max_usage"`
	VoucherValueOptionType string 	`json:"voucher_value_option_type"`
}
type PromoWithPagination struct {
	Data []*PromoDto `json:"data"`
	Meta *MetaPagination    `json:"meta"`
}