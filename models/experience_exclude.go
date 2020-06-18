package models

type ExperienceExclude struct {
	Id 		int `json:"id"`
	ExpId 	string `json:"exp_id"`
	ExcludeId int `json:"exclude_id"`
}

type ExperienceExcludeJoin struct {
	Id 		int `json:"id"`
	ExpId 	string `json:"exp_id"`
	ExcludeId int `json:"exclude_id"`
	ExcludeName	string	`json:"exclude_name"`
	ExcludeIcon	string	`json:"exclude_icon"`
}


