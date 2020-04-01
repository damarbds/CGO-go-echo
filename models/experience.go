package models

import "time"

type Experience struct {
	Id                      string     `json:"id" validate:"required"`
	CreatedBy               string     `json:"created_by":"required"`
	CreatedDate             time.Time  `json:"created_date" validate:"required"`
	ModifiedBy              *string    `json:"modified_by"`
	ModifiedDate            *time.Time `json:"modified_date"`
	DeletedBy               *string    `json:"deleted_by"`
	DeletedDate             *time.Time `json:"deleted_date"`
	IsDeleted               int        `json:"is_deleted" validate:"required"`
	IsActive                int        `json:"is_active" validate:"required"`
	ExpTitle                string     `json:"exp_title"`
	ExpType                 string     `json:"exp_type"`
	ExpTripType             string     `json:"exp_trip_type"`
	ExpBookingType          string     `json:"exp_booking_type"`
	ExpDesc                 string     `json:"exp_desc"`
	ExpMaxGuest             int        `json:"exp_max_guest"`
	ExpPickupPlace          string     `json:"exp_pickup_place"`
	ExpPickupTime           string     `json:"exp_pickup_time"`
	ExpPickupPlaceLongitude float64    `json:"exp_pickup_place_longitude"`
	ExpPickupPlaceLatitude  float64    `json:"exp_pickup_place_latitude"`
	ExpPickupPlaceMapsName  string     `json:"exp_pickup_place_maps_name"`
	ExpInternary            string     `json:"exp_internary"`
	ExpFacilities           string     `json:"exp_facilities"`
	ExpInclusion            string     `json:"exp_inclusion"`
	ExpRules                string     `json:"exp_rules"`
	Status                  int        `json:"status"`
	Rating                  float64    `json:"rating"`
	ExpLocationLatitude     float64    `json:"exp_location_latitude"`
	ExpLocationLongitude    float64    `json:"exp_location_longitude"`
	ExpLocationName         string     `json:"exp_location_name"`
	ExpCoverPhoto           *string    `json:"exp_cover_photo"`
	ExpDuration             int        `json:"exp_duration"`
	MinimumBookingId        string     `json:"minimum_booking_id"`
	MerchantId              string     `json:"merchant_id"`
	HarborsId               string     `json:"harbors_id"`
}

type ExperienceDto struct {
	Id                      string                `json:"id" validate:"required"`
	ExpTitle                string                `json:"exp_title"`
	ExpType                 []string              `json:"exp_type"`
	ExpTripType             string                `json:"exp_trip_type"`
	ExpBookingType          string                `json:"exp_booking_type"`
	ExpDesc                 string                `json:"exp_desc"`
	ExpMaxGuest             int                   `json:"exp_max_guest"`
	ExpPickupPlace          string                `json:"exp_pickup_place"`
	ExpPickupTime           string                `json:"exp_pickup_time"`
	ExpPickupPlaceLongitude float64               `json:"exp_pickup_place_longitude"`
	ExpPickupPlaceLatitude  float64               `json:"exp_pickup_place_latitude"`
	ExpPickupPlaceMapsName  string                `json:"exp_pickup_place_maps_name"`
	ExpInternary            ExpItineraryObject    `json:"exp_internary"`
	ExpFacilities           []ExpFacilitiesObject `json:"exp_facilities"`
	ExpInclusion            []ExpInclusionObject  `json:"exp_inclusion"`
	ExpRules                []ExpRulesObject      `json:"exp_rules"`
	Status                  int                   `json:"status"`
	Rating                  float64               `json:"rating"`
	ExpLocationLatitude     float64               `json:"exp_location_latitude"`
	ExpLocationLongitude    float64               `json:"exp_location_longitude"`
	ExpLocationName         string                `json:"exp_location_name"`
	ExpCoverPhoto           *string               `json:"exp_cover_photo"`
	ExpDuration             int                   `json:"exp_duration"`
	MinimumBookingId        string                `json:"minimum_booking_id"`
	MerchantId              string                `json:"merchant_id"`
	HarborsName             string                `json:"harbors_name"`
	City                    string                `json:"city"`
	Province                string                `json:"province"`
}
type ExpItineraryObject struct {
	Item []ItemObject `json:"item"`
}
type ItemObject struct {
	Day       int              `json:"day"`
	Itinerary []ItineraryObjet `json:"itinerary"`
}
type ItineraryObjet struct {
	Time     string `json:"time"`
	Activity string `json:"activity"`
}
type ExpFacilitiesObject struct {
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Amount int    `json:"amount"`
}

type ExpInclusionObject struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}
type ExpRulesObject struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type ExpSearch struct {
	Id          string  `json:"id" validate:"required"`
	ExpTitle    string  `json:"exp_title"`
	ExpType     string  `json:"exp_type"`
	Rating      float64 `json:"rating"`
}

type ExpSearchObject struct {
	Id          string   `json:"id" validate:"required"`
	ExpTitle    string   `json:"exp_title"`
	ExpType     []string `json:"exp_type"`
	Rating      float64  `json:"rating"`
	CountRating int      `json:"count_rating"`
	Currency    string   `json:"currency"`
	Price       float64  `json:"price"`
	PaymentType string   `json:"payment_type"`
}
type ExperienceUserDiscoverPreferenceDto struct {
	Id                      string                `json:"id" validate:"required"`
	ExpTitle                string                `json:"exp_title"`
	ExpType                 []string              `json:"exp_type"`
	Rating                  float64               `json:"rating"`
	CountRating				int 					`json:"count_rating"`
	Currency 				string					`json:"currency"`
	Price 					float64 					`json:"price"`
	Payment_type			string				  `json:"payment_type"`
}
type ExpUserDiscoverPreference struct {
	CityId  				int		`json:"city_id"`
	CityName				string	`json:"city_name"`
	CityDesc				string	`json:"city_desc"`
	Id                      string     `json:"id" validate:"required"`
	CreatedBy               string     `json:"created_by":"required"`
	CreatedDate             time.Time  `json:"created_date" validate:"required"`
	ModifiedBy              *string    `json:"modified_by"`
	ModifiedDate            *time.Time `json:"modified_date"`
	DeletedBy               *string    `json:"deleted_by"`
	DeletedDate             *time.Time `json:"deleted_date"`
	IsDeleted               int        `json:"is_deleted" validate:"required"`
	IsActive                int        `json:"is_active" validate:"required"`
	ExpTitle                string     `json:"exp_title"`
	ExpType                 string     `json:"exp_type"`
	ExpTripType             string     `json:"exp_trip_type"`
	ExpBookingType          string     `json:"exp_booking_type"`
	ExpDesc                 string     `json:"exp_desc"`
	ExpMaxGuest             int        `json:"exp_max_guest"`
	ExpPickupPlace          string     `json:"exp_pickup_place"`
	ExpPickupTime           string     `json:"exp_pickup_time"`
	ExpPickupPlaceLongitude float64    `json:"exp_pickup_place_longitude"`
	ExpPickupPlaceLatitude  float64    `json:"exp_pickup_place_latitude"`
	ExpPickupPlaceMapsName  string     `json:"exp_pickup_place_maps_name"`
	ExpInternary            string     `json:"exp_internary"`
	ExpFacilities           string     `json:"exp_facilities"`
	ExpInclusion            string     `json:"exp_inclusion"`
	ExpRules                string     `json:"exp_rules"`
	Status                  int        `json:"status"`
	Rating                  float64    `json:"rating"`
	ExpLocationLatitude     float64    `json:"exp_location_latitude"`
	ExpLocationLongitude    float64    `json:"exp_location_longitude"`
	ExpLocationName         string     `json:"exp_location_name"`
	ExpCoverPhoto           *string    `json:"exp_cover_photo"`
	ExpDuration             int        `json:"exp_duration"`
	MinimumBookingId        string     `json:"minimum_booking_id"`
	MerchantId              string     `json:"merchant_id"`
	HarborsId               string     `json:"harbors_id"`

}

type ExpUserDiscoverPreferenceDto struct {
	CityId	int 	`json:"city_id"`
	City	string 	`json:"city"`
	CityDesc	string	`json:"city_desc"`
	Item 	[]ExperienceUserDiscoverPreferenceDto `json:"item"`
}

