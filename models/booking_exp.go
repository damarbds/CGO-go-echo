package models

import "time"

type BookingExp struct {
	Id                   string    `json:"id" validate:"required"`
	CreatedBy            string    `json:"created_by":"required"`
	CreatedDate          time.Time `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int       `json:"is_deleted" validate:"required"`
	IsActive             int       `json:"is_active" validate:"required"`
	ExpId				string		`json:"exp_id"`
	OrderId				string		`json:"order_id"`
	GuestDesc			string		`json:"guest_desc"`
	BookedBy			string		`json:"booked_by"`
	BookedByEmail		string		`json:"booked_by_email"`
	BookingDate 		time.Time	`json:"booking_date"`
	UserId				*string		`json:"user_id"`
	Status 				int			`json:"status"`
	TicketCode			string		`json:"ticket_code"`
	TicketQRCode		string		`json:"ticket_qr_code"`
	ExperienceAddOnId 	*string		`json:"experience_add_on_id"`
}

type NewBookingExpCommand struct {
	Id                  string    `json:"id"`
	ExpId				string		`json:"exp_id"`
	GuestDesc			string		`json:"guest_desc"`
	BookedBy			string		`json:"booked_by"`
	BookedByEmail		string		`json:"booked_by_email"`
	BookingDate 		string	`json:"booking_date"`
	UserId				*string		`json:"user_id"`
	Status 				string			`json:"status"`
	OrderId				string		`json:"order_id"`
	TicketCode			string		`json:"ticket_code"`
	TicketQRCode		string		`json:"ticket_qr_code"`
	ExperienceAddOnId 	*string		`json:"experience_add_on_id"`
}