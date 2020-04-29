package models

type PromoMerchant struct {
	Id 		int 		`json:"id"`
	PromoId string		`json:"promo_id"`
	MerchantId 	string	`json:"merchant_id"`
}
