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
type BookingExpJoin struct {
	Id                   string    `json:"id" validate:"required"`
	CreatedBy            string    `json:"created_by" validate:"required"`
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
	ExpTitle			string		`json:"exp_title"`
	ExpType 			string		`json:"exp_type"`
	ExpPickupPlace 		string		`json:"exp_pickup_place"`
	ExpPickupTime 		string		`json:"exp_pickup_time"`
	ExpDuration int `json:"exp_duration"`
	TotalPrice 			float64		`json:"total_price"`
	PaymentType 		string		`json:"payment_type"`
	City string `json:"city"`
	Province string `json:"province"`
	Country string `json:"country"`
}
type BookingExpDetailDto struct {
	Id                   string    `json:"id" validate:"required"`
	GuestDesc			[]GuestDescObj		`json:"guest_desc"`
	BookedBy			[]BookedByObj		`json:"booked_by"`
	BookedByEmail		string		`json:"booked_by_email"`
	BookingDate 		time.Time	`json:"booking_date"`
	UserId				*string		`json:"user_id"`
	Status 				int			`json:"status"`
	//TicketCode			string		`json:"ticket_code"`
	OrderId				string		`json:"order_id"`
	TicketQRCode		string		`json:"ticket_qr_code"`
	ExperienceAddOnId 	*string		`json:"experience_add_on_id"`
	ExpId				string		`json:"exp_id"`
	ExpTitle			string		`json:"exp_title"`
	ExpType 			[]string	`json:"exp_type"`
	ExpPickupPlace 		string		`json:"exp_pickup_place"`
	ExpPickupTime 		string		`json:"exp_pickup_time"`
	TotalPrice 			float64		`json:"total_price"`
	PaymentType 		string		`json:"payment_type"`
}
type BookedByObj struct {
	Title 			string		`json:"title"`
	FullName 		string		`json:"fullname"`
	Email 			string		`json:"email"`
	PhoneNumber 	string		`json:"phonenumber"`
	IdType			string		`json:"idtype"`
	IdNumber 		string		`json:"idnumber"`
}
type GuestDescObj struct {
	FullName 		string	`json:"fullname"`
	IdType 			string	`json:"idtype"`
	IdNumber 		string	`json:"idnumber"`
	Type 			string	`json:"type"`
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
type MyBooking struct {
	ExpId string `json:"exp_id"`
	ExpTitle string `json:"exp_title"`
	BookingDate time.Time `json:"booking_date"`
	ExpDuration int `json:"exp_duration"`
	TotalGuest int `json:"total_guest"`
	City string `json:"city"`
	Province string `json:"province"`
	Country string `json:"country"`
}