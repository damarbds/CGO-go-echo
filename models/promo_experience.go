package models

type PromoExperience struct {
	Id 		int 		`json:"id"`
	PromoId string		`json:"promo_id"`
	ExperienceId string	`json:"experience_id"`
}
