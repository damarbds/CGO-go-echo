package models

import "time"

type ExperienceAddOn struct {
	Id                   string    `json:"id" validate:"required"`
	CreatedBy            string    `json:"created_by":"required"`
	CreatedDate          time.Time `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int       `json:"is_deleted" validate:"required"`
	IsActive             int       `json:"is_active" validate:"required"`
	Name				string		`json:"name"`
	Desc 				string		`json:"desc"`
	Currency 			int			`json:"currency"`
	Amount 				float64		`json:"amount"`
	ExpId				string		`json:"exp_id"`
}


type ExperienceAddOnDto struct {
	Id                   string    `json:"id" validate:"required"`
	Name				string		`json:"name"`
	Desc 				string		`json:"desc"`
	Currency 			string			`json:"currency"`
	Amount 				float64		`json:"amount"`
	ExpId				string		`json:"exp_id"`
}
