package models

type ExperienceFacilities struct {
	Id 			int 	`json:"id"`
	ExpId 		*string 	`json:"exp_id"`
	TransId 	*string 	`json:"trans_id"`
	FacilitiesId int	`json:"facilities_id"`
	Amount 		int 	`json:"amount"`
}

type ExperienceFacilitiesJoin struct {
	Id 			int 	`json:"id"`
	ExpId 		*string 	`json:"exp_id"`
	TransId 	*string 	`json:"trans_id"`
	FacilitiesId int	`json:"facilities_id"`
	Amount 		int 	`json:"amount"`
	FacilityName string         `json:"facility_name"`
	IsNumerable  int            `json:"is_numerable"`
	FacilityIcon *string`json:"facility_icon"`
}
