package models

import "time"

type Experience struct {
	Id                   string    `json:"id" validate:"required"`
	CreatedBy            string    `json:"created_by":"required"`
	CreatedDate          time.Time `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int       `json:"is_deleted" validate:"required"`
	IsActive             int       `json:"is_active" validate:"required"`
	ExpTile				 string 	`json:"exp_tile"`
	ExpType				 string		`json:"exp_type"`
	ExpTripType			 string		`json:"exp_trip_type"`
	ExpBookingType		 string		`json:"exp_booking_type"`
	ExpDesc				 string		`json:"exp_desc"`
	ExpMaxGuest			 int		 `json:"exp_max_guest"`
	ExpPickupPlace		string		`json:"exp_pickup_place"`
	ExpPickupPlaceLongitude	float64		`json:"exp_pickup_place_longitude"`
	ExpPickupPlaceLatitude	float64		`json:"exp_pickup_place_latitude"`
	ExpPickupPlaceMapsName string	`json:"exp_pickup_place_maps_name"`
	ExpInternary		string		`json:"exp_internary"`
	ExpFacilities		string		`json:"exp_facilities"`
	ExpInclusion		string		`json:"exp_inclusion"`
	ExpRules			string		`json:"exp_rules"`
	Status 				int			`json:"status"`
	Rating				float64		`json:"rating"`
	ExpLocationLatitude		float64		`json:"exp_location_latitude"`
	ExpLocationLongitude	float64		`json:"exp_location_longitude"`
	ExpLocationName		string		`json:"exp_location_name"`
	ExpDuration			int `json:"exp_duration"`
	MinimumBookingId 	string 	`json:"minimum_booking_id"`
	MerchantId			string	`json:"merchant_id"`

}






