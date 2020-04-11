package models

import "time"

type Transportation struct {
	Id                   string     `json:"id" validate:"required"`
	CreatedBy            string     `json:"created_by" validate:"required"`
	CreatedDate          time.Time  `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int        `json:"is_deleted" validate:"required"`
	IsActive             int        `json:"is_active" validate:"required"`
	TransName 			string		`json:"trans_name"`
	HarborsSourceId		string		`json:"harbors_source_id"`
	HarborsDestId		string		`json:"harbors_dest_id"`
	MerchantId			string		`json:"merchant_id"`
	TransCapacity		int 		`json:"trans_capacity"`
	TransTitle			string 		`json:"trans_title"`
	TransStatus			int			`json:"trans_status"`
	TransImages 		string		`json:"trans_images"`
	ReturnTransId		*string		`json:"return_trans_id"`
	BoatDetails 		string		`json:"boat_details"`
	Transcoverphoto string			`json:"transcoverphoto"`
	Class			string			`json:"class"`
}

type NewCommandTransportation struct {
	Id                   string     `json:"id"`
	TransName 			string		`json:"trans_name"`
	TransCapacity		int 		`json:"trans_capacity"`
	TransTitle			string 		`json:"trans_title"`
	Status			int			`json:"status"`
	BoatDetails 		BoatDetailsObj		`json:"boat_details"`
	Transcoverphoto string			`json:"transcoverphoto"`
	Class			string			`json:"class"`
	Facilities		[]string	`json:"facilities"`
	TransImages 	[]CoverPhotosObj	`json:"trans_photos"`
	DepartureRoute 	RouteObj		`json:"departure_route"`
	ReturnRoute		*RouteObj		`json:"return_route"`
}
type RouteObj struct {
	HarborsIdFrom 		string		`json:"harbors_id_from"`
	HarborsIdTo			string		`json:"harbors_id_to"`
	Time 				[]TimeObj	`json:"time"`
	Schedule 		[]YearObj	`json:"schedule"`

}
type YearObj struct {
	Year 	int	`json:"year"`
	Month 	[]MonthObj	`json:"month"`
} 
type MonthObj struct {
	Month 		string	`json:"month"`
	DayPrice 	[]DayPriceObj	`json:"day_price"`
} 
type DayPriceObj struct {
	Day 		string	`json:"day"`
	AdultPrice	float64	`json:"adult_price"`
	ChildrenPrice float64	`json:"children_price"`
	Currency 	 string		`json:"currency"`
}
type TimeObj struct {
	DepartureTime		string		`json:"departure_time"`
	ArrivalTime			string		`json:"arrival_time"`
}
type BoatDetailsObj struct {
	Length 			float64		`json:"length"`
	Width 			float64		`json:"width"`
	Machine 		string		`json:"machine"`
	Cabin 			int 		`json:"cabin"`
}