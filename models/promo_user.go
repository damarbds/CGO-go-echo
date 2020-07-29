package models

type PromoUser struct {
	Id 		int 		`json:"id"`
	PromoId string		`json:"promo_id"`
	UserId 	string		`json:"user_id"`
}
