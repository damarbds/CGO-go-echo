package models

type ExperienceInclude struct {
	Id 			int 	`json:"id"`
	ExpId 		string 	`json:"exp_id"`
	IncludeId	int `json:"include_id"`
}

type ExperienceIncludeJoin struct {
	Id 			int 	`json:"id"`
	ExpId 		string 	`json:"exp_id"`
	IncludeId	int `json:"include_id"`
	IncludeName	 string	`json:"include_name"`
	IncludeIcon  string	`json:"include_icon"`
}
