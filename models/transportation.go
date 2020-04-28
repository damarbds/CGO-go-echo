package models

import "time"

type Transportation struct {
	Id              string     `json:"id" validate:"required"`
	CreatedBy       string     `json:"created_by" validate:"required"`
	CreatedDate     time.Time  `json:"created_date" validate:"required"`
	ModifiedBy      *string    `json:"modified_by"`
	ModifiedDate    *time.Time `json:"modified_date"`
	DeletedBy       *string    `json:"deleted_by"`
	DeletedDate     *time.Time `json:"deleted_date"`
	IsDeleted       int        `json:"is_deleted" validate:"required"`
	IsActive        int        `json:"is_active" validate:"required"`
	TransName       string     `json:"trans_name"`
	HarborsSourceId *string     `json:"harbors_source_id"`
	HarborsDestId   *string     `json:"harbors_dest_id"`
	MerchantId      string     `json:"merchant_id"`
	TransCapacity   int        `json:"trans_capacity"`
	TransTitle      string     `json:"trans_title"`
	TransStatus     int        `json:"trans_status"`
	TransImages     string     `json:"trans_images"`
	ReturnTransId   *string    `json:"return_trans_id"`
	BoatDetails     string     `json:"boat_details"`
	Transcoverphoto string     `json:"transcoverphoto"`
	Class           string     `json:"class"`
	TransFacilities *string  `json:"trans_facilities"`
}

type NewCommandTransportation struct {
	Id              string           `json:"id"`
	TransName       string           `json:"trans_name"`
	TransCapacity   int              `json:"trans_capacity"`
	TransTitle      string           `json:"trans_title"`
	Status          int              `json:"status"`
	BoatDetails     BoatDetailsObj   `json:"boat_details"`
	Transcoverphoto string           `json:"transcoverphoto"`
	Class           string           `json:"class"`
	Facilities      []ExpFacilitiesObject         `json:"facilities"`
	TransImages     []CoverPhotosObj `json:"trans_photos"`
	DepartureRoute  RouteObj         `json:"departure_route"`
	ReturnRoute     *RouteObj        `json:"return_route"`
}

type TransportationDto struct {
	Id              string           `json:"id"`
	TransName       string           `json:"trans_name"`
	TransCapacity   int              `json:"trans_capacity"`
	TransTitle      string           `json:"trans_title"`
	Status          int              `json:"status"`
	BoatDetails     BoatDetailsObj   `json:"boat_details"`
	Transcoverphoto string           `json:"transcoverphoto"`
	Class           string           `json:"class"`
	Facilities      []ExpFacilitiesObject         `json:"facilities"`
	TransImages     []CoverPhotosObj `json:"trans_photos"`
	DepartureRoute  RouteObj         `json:"departure_route"`
	ReturnRoute     *RouteObj        `json:"return_route"`
}
type RouteObj struct {
	Id            string    `json:"id"`
	HarborsIdFrom *string    `json:"harbors_id_from"`
	HarborsIdTo   *string    `json:"harbors_id_to"`
	Time          []TimeObj `json:"time"`
	Schedule      []YearObj `json:"schedule"`
}
type YearObj struct {
	Year  int        `json:"year"`
	Month []MonthObj `json:"month"`
}
type MonthObj struct {
	Month    string        `json:"month"`
	DayPrice []DayPriceObj `json:"day_price"`
}
type DayPriceObj struct {
	DepartureDate string  `json:"departure_date"`
	Day           string  `json:"day"`
	AdultPrice    float64 `json:"adult_price"`
	ChildrenPrice float64 `json:"children_price"`
	Currency      string  `json:"currency"`
}
type TimeObj struct {
	DepartureTime string `json:"departure_time"`
	ArrivalTime   string `json:"arrival_time"`
}
type BoatDetailsObj struct {
	Length  float64 `json:"length"`
	Width   float64 `json:"width"`
	Machine string  `json:"machine"`
	Cabin   int     `json:"cabin"`
}
type TransSearch struct {
	ScheduleId       *string `json:"schedule_id"`
	DepartureDate    *string `json:"departure_date"`
	DepartureTime    *string `json:"departure_time"`
	ArrivalTime      *string `json:"arrival_time"`
	Price            *string `json:"price"`
	TransId          string `json:"trans_id"`
	TransName        string `json:"trans_name"`
	TransImages      string `json:"trans_images"`
	TransStatus int `json:"trans_status"`
	HarborSourceId   *string `json:"harbor_source_id"`
	HarborSourceName *string `json:"harbor_source_name"`
	HarborDestId     *string `json:"harbor_dest_id"`
	HarborDestName   *string `json:"harbor_dest_name"`
	MerchantName 		  string		`json:"merchant_name"`
	MerchantPicture 	 *string		`json:"merchant_picture"`
	Class				  *string 		`json:"class"`
	TransFacilities		*string		`json:"trans_facilities"`
	TransCapacity 		int			`json:"trans_capacity"`
}
type TransPriceObj struct {
	AdultPrice    float64 `json:"adult_price"`
	ChildrenPrice float64 `json:"children_price"`
	Currency      int     `json:"currency"`
	CurrencyLabel string  `json:"currency_label"`
	PriceType     string  `json:"price_type"`
}
type TransportationSearchObj struct {
	ScheduleId            *string           `json:"schedule_id"`
	DepartureDate         *string           `json:"departure_date"`
	DepartureTime         *string           `json:"departure_time"`
	ArrivalTime           *string           `json:"arrival_time"`
	TripDuration          *string           `json:"trip_duration"`
	TransportationId      string           `json:"transportation_id"`
	TransportationName    string           `json:"transportation_name"`
	TransportationImages  []CoverPhotosObj `json:"transportation_images"`
	TransportationStatus string `json:"transportation_status"`
	HarborSourceId        *string           `json:"harbor_source_id"`
	HarborSourceName      *string           `json:"harbor_source_name"`
	HarborDestinationId   *string           `json:"harbor_destination_id"`
	HarborDestinationName *string           `json:"harbor_destination_name"`
	Price                 TransPriceObj    `json:"price"`
	MerchantName 		  string		`json:"merchant_name"`
	MerchantPicture 	 *string		`json:"merchant_picture"`
	Class				  *string 		`json:"class"`
	TransFacilities 	[]ExpFacilitiesObject `json:"trans_facilities"`
}
type FilterSearchTransWithPagination struct {
	Data []*TransportationSearchObj `json:"data"`
	Meta *MetaPagination            `json:"meta"`
}