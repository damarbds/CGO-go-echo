package models

type PromoExperienceTransport struct {
	Id 		int 		`json:"id"`
	PromoId string		`json:"promo_id"`
	ExperienceId *string	`json:"experience_id"`
	TransportationId *string `json:"transportation_id"`
}