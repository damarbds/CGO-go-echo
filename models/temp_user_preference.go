package models

import "time"

type TempUserPreference struct {
	Id 			int 		`json:"id"`
	ProvinceId  *int 		`json:"province_id"`
	CityId 		*int 		`json:"city_id"`
	HarborsId 	*string		`json:"harbors_id"`
}

type TempUserPreferenceWithJoin struct {
	CityId                  int        `json:"city_id"`
	CityName                string     `json:"city_name"`
	CityDesc                string     `json:"city_desc"`
	CityPhotos              *string    `json:"city_photos"`
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
