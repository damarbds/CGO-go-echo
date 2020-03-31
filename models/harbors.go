package models

import "time"

type Harbors struct {
	Id                   string    `json:"id" validate:"required"`
	CreatedBy            string    `json:"created_by":"required"`
	CreatedDate          time.Time `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int       `json:"is_deleted" validate:"required"`
	IsActive             int       `json:"is_active" validate:"required"`
	HarborsName			string		`json:"harbors_name"`
	HarborsLongitude	float64		`json:"harbors_longitude"`
	HarborsLatitude		float64		`json:"harbors_latitude"`
	HarborsImage		string		`json:"harbors_image"`
	CityId				int			`json:"city_id"`
}
