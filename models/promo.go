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
	StartTripPeriod	*string	`json:"start_trip_period"`
	EndTripPeriod	*string	`json:"end_trip_period"`
	IsAnyTripPeriod int	`json:"is_any_trip_period"`
	//Currency 	*int 		`json:"currency"`
	HowToGet	string		`json:"how_to_get"`
	HowToUse	string		`json:"how_to_use"`
	TermCondition	string	`json:"term_condition"`
	Disclaimer	string		`json:"disclaimer"`
	MaxDiscount	float32		`json:"max_discount"`
	MaxUsage   	*int 		`json:"max_usage"`
	ProductionCapacity	*int `json:"production_capacity"`
	CurrencyId *int 	`json:"currency_id"`
	PromoProductType *int `json:"promo_product_type"`
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
	StartTripPeriod	*string	`json:"start_trip_period"`
	EndTripPeriod	*string	`json:"end_trip_period"`
	IsAnyTripPeriod	int	`json:"is_any_trip_period"`
	Currency 	*int 		`json:"currency"`
	MaxDiscount	float32		`json:"max_value"`
	MaxUsage   	*int 		`json:"max_usage"`
	HowToGet	string		`json:"how_to_get"`
	HowToUse	string		`json:"how_to_use"`
	TermCondition	string	`json:"term_condition"`
	Disclaimer	string		`json:"disclaimer"`
	ProductionCapacity	*int `json:"production_capacity"`
	MerchantId	[]string	`json:"merchant_id"`
	PromoProductType *int `json:"promo_product_type"`
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
	StartTripPeriod	string	`json:"start_trip_period"`
	EndTripPeriod	string	`json:"end_trip_period"`
	MaxDiscount	float32		`json:"max_discount"`
	HowToGet	string		`json:"how_to_get"`
	HowToUse	string		`json:"how_to_use"`
	IsAnyTripPeriod	int	`json:"is_any_trip_period"`
	TermCondition	string	`json:"term_condition"`
	Disclaimer	string		`json:"disclaimer"`
	Currency 	int 		`json:"currency"`
	MaxUsage   	int 		`json:"max_usage"`
	ProductionCapacity int 	`json:"production_capacity"`
	MerchantId 	[]string	`json:"merchant_id"`
	PromoProductType *int `json:"promo_product_type"`
}
type PromoWithPagination struct {
	Data []*PromoDto `json:"data"`
	Meta *MetaPagination    `json:"meta"`
}