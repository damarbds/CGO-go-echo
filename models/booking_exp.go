package models

import (
	"database/sql"
	"time"
)

type BookingExpHistory struct {
	Id                 string     `json:"id" validate:"required"`
	CreatedBy          string     `json:"created_by":"required"`
	CreatedDate        time.Time  `json:"created_date" validate:"required"`
	ModifiedBy         *string    `json:"modified_by"`
	ModifiedDate       *time.Time `json:"modified_date"`
	DeletedBy          *string    `json:"deleted_by"`
	DeletedDate        *time.Time `json:"deleted_date"`
	IsDeleted          int        `json:"is_deleted" validate:"required"`
	IsActive           int        `json:"is_active" validate:"required"`
	ExpId              string     `json:"exp_id"`
	OrderId            string     `json:"order_id"`
	GuestDesc          string     `json:"guest_desc"`
	BookedBy           string     `json:"booked_by"`
	BookedByEmail      string     `json:"booked_by_email"`
	BookingDate        time.Time  `json:"booking_date"`
	ExpiredDatePayment *time.Time `json:"expired_date_payment"`
	UserId             *string    `json:"user_id"`
	Status             int        `json:"status"`
	TicketCode         string     `json:"ticket_code"`
	TicketQRCode       string     `json:"ticket_qr_code"`
	ExperienceAddOnId  *string    `json:"experience_add_on_id"`
	TransId            *string    `json:"trans_id"`
	PaymentUrl         *string    `json:"payment_url"`
	ScheduleId 	   *string	`json:"schedule_id"`
	ExpTitle           string     `json:"exp_title"`
	ExpType            *string    `json:"exp_type"`
	ExpDuration        int        `json:"exp_duration"`
	CityName           string     `json:"city_name"`
	ProvinceName       string     `json:"province_name"`
	CountryName        string     `json:"country_name"`
	StatusTransaction  int        `json:"status_transaction"`
}
type BookingHistoryDto struct {
	Category string            `json:"category"`
	Items    []ItemsHistoryDto `json:"items"`
}
type ItemsHistoryDto struct {
	ExpId          string    `json:"exp_id"`
	ExpTitle       string    `json:"exp_title"`
	ExpType        []string  `json:"exp_type"`
	ExpBookingDate time.Time `json:"exp_booking_date"`
	ExpDuration    int       `json:"exp_duration"`
	TotalGuest     int       `json:"total_guest"`
	City           string    `json:"city"`
	Province       string    `json:"province"`
	Country        string    `json:"country"`
	Status         int       `json:"status"`
}
type BookingExp struct {
	Id                 string     `json:"id" validate:"required"`
	CreatedBy          string     `json:"created_by":"required"`
	CreatedDate        time.Time  `json:"created_date" validate:"required"`
	ModifiedBy         *string    `json:"modified_by"`
	ModifiedDate       *time.Time `json:"modified_date"`
	DeletedBy          *string    `json:"deleted_by"`
	DeletedDate        *time.Time `json:"deleted_date"`
	IsDeleted          int        `json:"is_deleted" validate:"required"`
	IsActive           int        `json:"is_active" validate:"required"`
	ExpId              *string    `json:"exp_id"`
	OrderId            string     `json:"order_id"`
	GuestDesc          string     `json:"guest_desc"`
	BookedBy           string     `json:"booked_by"`
	BookedByEmail      string     `json:"booked_by_email"`
	BookingDate        time.Time  `json:"booking_date"`
	ExpiredDatePayment *time.Time `json:"expired_date_payment"`
	UserId             *string    `json:"user_id"`
	Status             int        `json:"status"`
	TicketCode         string     `json:"ticket_code"`
	TicketQRCode       string     `json:"ticket_qr_code"`
	ExperienceAddOnId  *string    `json:"experience_add_on_id"`
	TransId            *string    `json:"trans_id"`
	PaymentUrl         *string    `json:"payment_url"`
	ScheduleId 	   *string	`json:"schedule_id"`
}
type BookingExpJoin struct {
	Id                     string         `json:"id" validate:"required"`
	CreatedBy              string         `json:"created_by" validate:"required"`
	CreatedDate            time.Time      `json:"created_date" validate:"required"`
	ModifiedBy             *string        `json:"modified_by"`
	ModifiedDate           *time.Time     `json:"modified_date"`
	DeletedBy              *string        `json:"deleted_by"`
	DeletedDate            *time.Time     `json:"deleted_date"`
	IsDeleted              int            `json:"is_deleted" validate:"required"`
	IsActive               int            `json:"is_active" validate:"required"`
	ExpId                  string         `json:"exp_id"`
	OrderId                string         `json:"order_id"`
	GuestDesc              string         `json:"guest_desc"`
	BookedBy               string         `json:"booked_by"`
	BookedByEmail          string         `json:"booked_by_email"`
	BookingDate            time.Time      `json:"booking_date"`
	ExpiredDatePayment     *time.Time     `json:"expired_date_payment"`
	UserId                 *string        `json:"user_id"`
	Status                 int            `json:"status"`
	TransactionStatus      int            `json:"transaction_status"`
	TicketCode             string         `json:"ticket_code"`
	TicketQRCode           string         `json:"ticket_qr_code"`
	ExperienceAddOnId      *string        `json:"experience_add_on_id"`
	TransId                *string        `json:"trans_id"`
	PaymentUrl             *string        `json:"payment_url"`
	ScheduleId 	   *string	`json:"schedule_id"`
	ExpTitle               string         `json:"exp_title"`
	ExpType                string         `json:"exp_type"`
	ExpPickupPlace         string         `json:"exp_pickup_place"`
	ExpPickupTime          string         `json:"exp_pickup_time"`
	ExpDuration            int            `json:"exp_duration"`
	TotalPrice             float64        `json:"total_price"`
	PaymentType            string         `json:"payment_type"`
	City                   string         `json:"city"`
	Province               string         `json:"province"`
	Country                string         `json:"country"`
	ExperiencePaymentId    string         `json:"experience_payment_id"`
	Currency               int            `json:"currency"`
	AccountBank            string         `json:"account_bank"`
	Icon                   string         `json:"icon"`
	CreatedDateTransaction time.Time      `json:"created_date_transaction"`
	MerchantName           sql.NullString `json:"merchant_name"`
	MerchantPhone          sql.NullString `json:"merchant_phone"`
	MerchantPicture        sql.NullString `json:"merchant_picture"`
}
type BookingExpDetailDto struct {
	Id                     string         `json:"id" validate:"required"`
	GuestDesc              []GuestDescObj `json:"guest_desc"`
	BookedBy               []BookedByObj  `json:"booked_by"`
	BookedByEmail          string         `json:"booked_by_email"`
	BookingDate            time.Time      `json:"booking_date"`
	ExpiredDatePayment     *time.Time     `json:"expired_date_payment"`
	CreatedDateTransaction time.Time      `json:"created_date_transaction"`
	UserId                 *string        `json:"user_id"`
	Status                 int            `json:"status"`
	TransactionStatus      int            `json:"transaction_status"`
	//TicketCode			string		`json:"ticket_code"`
	OrderId             string   `json:"order_id"`
	TicketQRCode        string   `json:"ticket_qr_code"`
	ExperienceAddOnId   *string  `json:"experience_add_on_id"`
	ExpId               string   `json:"exp_id"`
	ExpTitle            string   `json:"exp_title"`
	ExpType             []string `json:"exp_type"`
	ExpPickupPlace      string   `json:"exp_pickup_place"`
	ExpPickupTime       string   `json:"exp_pickup_time"`
	TotalPrice          float64  `json:"total_price"`
	Currency            string   `json:"currency"`
	PaymentType         string   `json:"payment_type"`
	AccountNumber       string   `json:"account_number"`
	AccountHolder       string   `json:"account_holder"`
	BankIcon            string   `json:"bank_icon"`
	ExperiencePaymentId string   `json:"experience_payment_id"`
	MerchantName        string   `json:"merchant_name"`
	MerchantPhone       string   `json:"merchant_phone"`
	MerchantPicture     string   `json:"merchant_picture"`
}
type AccountDesc struct {
	AccNumber string `json:"acc_number"`
	AccHolder string `json:"acc_holder"`
}
type BookedByObj struct {
	Title       string      `json:"title"`
	FullName    string      `json:"fullname"`
	Email       string      `json:"email"`
	PhoneNumber interface{} `json:"phonenumber"`
	IdType      string      `json:"idtype"`
	IdNumber    string      `json:"idnumber"`
}
type GuestDescObj struct {
	FullName string `json:"fullname"`
	IdType   string `json:"idtype"`
	IdNumber string `json:"idnumber"`
	Type     string `json:"type"`
}
type NewBookingExpCommand struct {
	Id                string  `json:"id"`
	ExpId             string  `json:"exp_id"`
	GuestDesc         string  `json:"guest_desc"`
	BookedBy          string  `json:"booked_by"`
	BookedByEmail     string  `json:"booked_by_email"`
	BookingDate       string  `json:"booking_date"`
	UserId            *string `json:"user_id"`
	Status            string  `json:"status"`
	OrderId           string  `json:"order_id"`
	TicketCode        string  `json:"ticket_code"`
	TicketQRCode      string  `json:"ticket_qr_code"`
	ExperienceAddOnId *string `json:"experience_add_on_id"`
	TransId           *string  `json:"trans_id"`
	PaymentUrl        string  `json:"payment_url"`
}
type MyBooking struct {
	ExpId       string    `json:"exp_id"`
	ExpTitle    string    `json:"exp_title"`
	BookingDate time.Time `json:"booking_date"`
	ExpDuration int       `json:"exp_duration"`
	TotalGuest  int       `json:"total_guest"`
	City        string    `json:"city"`
	Province    string    `json:"province"`
	Country     string    `json:"country"`
}
type BookingGrowth struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}
type BookingGrowthDto struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}
type BookingTransactionExp struct {
	Id                 string     `json:"id" validate:"required"`
	CreatedBy          string     `json:"created_by":"required"`
	CreatedDate        time.Time  `json:"created_date" validate:"required"`
	ModifiedBy         *string    `json:"modified_by"`
	ModifiedDate       *time.Time `json:"modified_date"`
	DeletedBy          *string    `json:"deleted_by"`
	DeletedDate        *time.Time `json:"deleted_date"`
	IsDeleted          int        `json:"is_deleted" validate:"required"`
	IsActive           int        `json:"is_active" validate:"required"`
	ExpId              *string    `json:"exp_id"`
	OrderId            string     `json:"order_id"`
	GuestDesc          string     `json:"guest_desc"`
	BookedBy           string     `json:"booked_by"`
	BookedByEmail      string     `json:"booked_by_email"`
	BookingDate        time.Time  `json:"booking_date"`
	ExpiredDatePayment *time.Time `json:"expired_date_payment"`
	UserId             *string    `json:"user_id"`
	Status             int        `json:"status"`
	TicketCode         string     `json:"ticket_code"`
	TicketQRCode       string     `json:"ticket_qr_code"`
	ExperienceAddOnId  *string    `json:"experience_add_on_id"`
	TransId            *string    `json:"trans_id"`
	PaymentUrl         *string    `json:"payment_url"`
	ScheduleId 	   *string	`json:"schedule_id"`
	TotalPrice         float64    `json:"total_price"`
}
