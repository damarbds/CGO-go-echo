package models

import "time"

type MinimumBooking struct {
	Id                   string     `json:"id" validate:"required"`
	CreatedBy            string     `json:"created_by":"required"`
	CreatedDate          time.Time  `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int        `json:"is_deleted" validate:"required"`
	IsActive             int        `json:"is_active" validate:"required"`
	MinimumBookingDesc   string     `json:"minimum_booking_desc"`
	MinimumBookingAmount int        `json:"minimum_booking_amount"`
}

type MinimumBookingDto struct {
	Id                   string     `json:"id" validate:"required"`
	MinimumBookingDesc   string     `json:"minimum_booking_desc"`
	MinimumBookingAmount int        `json:"minimum_booking_amount"`
}

type MinimumBookingDtoWithPagination struct {
	Data []*MinimumBookingDto `json:"data"`
	Meta *MetaPagination    `json:"meta"`
}
