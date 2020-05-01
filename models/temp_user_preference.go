package models

type TempUserPreference struct {
	Id 			int 		`json:"id"`
	ProvinceId  int 		`json:"province_id"`
	CityId 		int 		`json:"city_id"`
	HarborsId 	string		`json:"harbors_id"`
}
