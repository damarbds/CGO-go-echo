package models

import "time"

type Schedule struct {
	Id                 string     `json:"id" validate:"required"`
	CreatedBy          string     `json:"created_by":"required"`
	CreatedDate        time.Time  `json:"created_date" validate:"required"`
	ModifiedBy         *string    `json:"modified_by"`
	ModifiedDate       *time.Time `json:"modified_date"`
	DeletedBy          *string    `json:"deleted_by"`
	DeletedDate        *time.Time `json:"deleted_date"`
	IsDeleted          int        `json:"is_deleted" validate:"required"`
	IsActive           int        `json:"is_active" validate:"required"`
	TransId				string	`json:"trans_id"`
	DepartureTime		string	`json:"departure_time"`
	ArrivalTime			string	`json:"arrival_time"`
	Day 				string		`json:"day"`
	Month				string	`json:"month"`
	Year				int 	`json:"year"`
	Price 				string	`json:"price"`
	//Currency 			int 	`json:"currency"`
	DepartureTimeoptionId	*int `json:"departure_timeoption_id"`
	ArrivalTimeoptionId		*int `json:"arrival_timeoption_id"`
}

type PriceObj struct {
	AdultPrice	float64	`json:"adult_price"`
	ChildrenPrice float64	`json:"children_price"`
	Currency 	 int		`json:"currency"`
}