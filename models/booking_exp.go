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
	ScheduleId         *string    `json:"schedule_id"`
	ExpTitle           string     `json:"exp_title"`
	ExpType            *string    `json:"exp_type"`
	ExpDuration        int        `json:"exp_duration"`
	CityName           string     `json:"city_name"`
	ProvinceName       string     `json:"province_name"`
	CountryName        string     `json:"country_name"`
	StatusTransaction  int        `json:"status_transaction"`
	TransName          *string    `json:"trans_name"`
	TransTitle         *string    `json:"trans_title"`
	TransStatus        *string    `json:"trans_status"`
	TransClass         *string    `json:"trans_class"`
	DepartureDate      *time.Time `json:"departure_date"`
	DepartureTime      *string    `json:"departure_time"`
	ArrivalTime        *string    `json:"arrival_time"`
	HarborSourceName   *string    `json:"harbor_source_name"`
	HarborDestName     *string    `json:"harbor_dest_name"`
}
type BookingHistoryDto struct {
	Category string            `json:"category"`
	Items    []ItemsHistoryDto `json:"items"`
}
type BookingHistoryDtoWithPagination struct {
	Data []*BookingHistoryDto `json:"data"`
	Meta *MetaPagination      `json:"meta"`
}
type ItemsHistoryDto struct {
	OrderId            string                   `json:"order_id"`
	ExpId              string                   `json:"exp_id"`
	ExpTitle           string                   `json:"exp_title"`
	ExpType            []string                 `json:"exp_type"`
	TransId            string                   `json:"trans_id"`
	TransName          string                   `json:"trans_name"`
	TransFrom          string                   `json:"trans_from"`
	TransTo            string                   `json:"trans_to"`
	TransDepartureTime *string                  `json:"trans_departure_time"`
	TransArrivalTime   *string                  `json:"trans_arrival_time"`
	TripDuration       string                   `json:"trip_duration"`
	TransClass         string                   `json:"trans_class"`
	TransGuest         TotalGuestTransportation `json:"trans_guest"`
	ExpBookingDate     time.Time                `json:"exp_booking_date"`
	ExpDuration        int                      `json:"exp_duration"`
	TotalGuest         int                      `json:"total_guest"`
	ExpGuest           TotalGuestTransportation `json:"exp_guest"`
	City               string                   `json:"city"`
	Province           string                   `json:"province"`
	Country            string                   `json:"country"`
	Status             string                   `json:"status"`
	IsReview           bool                     `json:"is_review"`
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
	ScheduleId         *string    `json:"schedule_id"`
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
	ExpId                  *string        `json:"exp_id"`
	OrderId                string         `json:"order_id"`
	GuestDesc              string         `json:"guest_desc"`
	BookedBy               string         `json:"booked_by"`
	BookedByEmail          string         `json:"booked_by_email"`
	BookingDate            time.Time      `json:"booking_date"`
	ExpiredDatePayment     *time.Time     `json:"expired_date_payment"`
	UserId                 *string        `json:"user_id"`
	Status                 int            `json:"status"`
	TransactionStatus      *int           `json:"transaction_status"`
	VaNumber               *string        `json:"va_number"`
	TicketCode             string         `json:"ticket_code"`
	TicketQRCode           string         `json:"ticket_qr_code"`
	ExperienceAddOnId      *string        `json:"experience_add_on_id"`
	TransId                *string        `json:"trans_id"`
	PaymentUrl             *string        `json:"payment_url"`
	ScheduleId             *string        `json:"schedule_id"`
	ExpTitle               *string        `json:"exp_title"`
	ExpType                *string        `json:"exp_type"`
	ExpPickupPlace         *string        `json:"exp_pickup_place"`
	ExpPickupTime          *string        `json:"exp_pickup_time"`
	ExpDuration            *int           `json:"exp_duration"`
	TotalPrice             *float64       `json:"total_price"`
	PaymentType            *string        `json:"payment_type"`
	HarborsName            *string        `json:"harbors_name"`
	City                   string         `json:"city"`
	Province               string         `json:"province"`
	Country                string         `json:"country"`
	ExperiencePaymentId    string         `json:"experience_payment_id"`
	Currency               int            `json:"currency"`
	AccountBank            *string        `json:"account_bank"`
	Icon                   *string        `json:"icon"`
	CreatedDateTransaction *time.Time     `json:"created_date_transaction"`
	MerchantName           sql.NullString `json:"merchant_name"`
	MerchantPhone          sql.NullString `json:"merchant_phone"`
	MerchantPicture        sql.NullString `json:"merchant_picture"`
	TransName              *string        `json:"trans_name"`
	TransTitle             *string        `json:"trans_title"`
	TransStatus            *string        `json:"trans_status"`
	TransClass             *string        `json:"trans_class"`
	DepartureDate          *time.Time     `json:"departure_date"`
	DepartureTime          *string        `json:"departure_time"`
	ArrivalTime            *string        `json:"arrival_time"`
	HarborSourceName       *string        `json:"harbor_source_name"`
	HarborDestName         *string        `json:"harbor_dest_name"`
	ExChangeRates 			*float64	`json:"ex_change_rates"`
	ExChangeCurrency 		*string		`json:"ex_change_currency"`
	ExpPaymentDeadlineAmount *int 		`json:"exp_payment_deadline_amount"`
	ExpPaymentDeadlineType 	*string		`json:"exp_payment_deadline_type"`
	ReturnTransId			*string `json:"return_trans_id"`
	OriginalPrice 		*float64	`json:"original_price"`
}

type BookingTransportationDetail struct {
	TransID          string    `json:"trans_id"`
	TransName        string    `json:"trans_name"`
	TransTitle       string    `json:"trans_title"`
	TransStatus      string    `json:"trans_status"`
	TransClass       string    `json:"trans_class"`
	DepartureDate    time.Time `json:"departure_date"`
	DepartureTime    string    `json:"departure_time"`
	ArrivalTime      string    `json:"arrival_time"`
	TripDuration     string    `json:"trip_duration"`
	HarborSourceName string    `json:"harbor_source_name"`
	HarborDestName   string    `json:"harbor_dest_name"`
	MerchantName     string    `json:"merchant_name"`
	MerchantPhone    string    `json:"merchant_phone"`
	MerchantPicture  string    `json:"merchant_picture"`
	TotalGuest       int       `json:"total_guest"`
	ReturnTransId 	*string		`json:"return_trans_id"`
}

type BookingExpDetail struct {
	ExpId           string               `json:"exp_id"`
	ExpTitle        string               `json:"exp_title"`
	ExpType         []string             `json:"exp_type"`
	ExpPickupPlace  string               `json:"exp_pickup_place"`
	ExpPickupTime   string               `json:"exp_pickup_time"`
	MerchantName    string               `json:"merchant_name"`
	MerchantPhone   string               `json:"merchant_phone"`
	MerchantPicture string               `json:"merchant_picture"`
	TotalGuest      int                  `json:"total_guest"`
	City            string               `json:"city"`
	ProvinceName    string               `json:"province_name"`
	HarborsName     string               `json:"harbors_name"`
	ExperienceAddOn []ExperienceAddOnObj `json:"experience_add_on"`
	ExpDuration     int                  `json:"exp_duration"`
	CountryName 	string 				 `json:"country_name"`
	ExpPaymentDeadlineAmount *int 		`json:"exp_payment_deadline_amount"`
	ExpPaymentDeadlineType 	*string		`json:"exp_payment_deadline_type"`
}

type BookingExpDetailDto struct {
	Id                     string                        `json:"id" validate:"required"`
	GuestDesc              []GuestDescObj                `json:"guest_desc"`
	BookedBy               []BookedByObj                 `json:"booked_by"`
	BookedByEmail          string                        `json:"booked_by_email"`
	BookingDate            time.Time                     `json:"booking_date"`
	ExpiredDatePayment     *time.Time                    `json:"expired_date_payment"`
	CreatedDateTransaction *time.Time                    `json:"created_date_transaction"`
	UserId                 *string                       `json:"user_id"`
	Status                 int                           `json:"status"`
	TransactionStatus      *int                          `json:"transaction_status"`
	OrderId                string                        `json:"order_id"`
	TicketQRCode           string                        `json:"ticket_qr_code"`
	ExperienceAddOnId      *string                       `json:"experience_add_on_id"`
	TotalPrice             *float64                      `json:"total_price"`
	Currency               string                        `json:"currency"`
	PaymentType            *string                       `json:"payment_type"`
	AccountNumber          string                        `json:"account_number"`
	AccountHolder          string                        `json:"account_holder"`
	BankIcon               *string                       `json:"bank_icon"`
	ExperiencePaymentId    string                        `json:"experience_payment_id"`
	Experience             []BookingExpDetail            `json:"experience,omitempty"`
	Transportation         []BookingTransportationDetail `json:"transportation,omitempty"`
	ExperiencePaymentType  *ExperiencePaymentTypeDto     `json:"experience_payment_type"`
	IsReview               bool                          `json:"is_review"`
	ReviewDesc             *string                       `json:"review_desc"`
	GuideReview            *float64                      `json:"guide_review"`
	ActivitiesReview       *float64                      `json:"activities_review"`
	ServiceReview          *float64                      `json:"service_review"`
	CleanlinessReview      *float64                      `json:"cleanliness_review"`
	ValueReview            *float64                      `json:"value_review"`
	MidtransUrl            *string                       `json:"midtrans_url"`
	ExChangeRates 			*float64	`json:"ex_change_rates"`
	ExChangeCurrency 		*string		`json:"ex_change_currency"`
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
type GuestDescObjForHTML struct {
	No 		int 	`json:"no"`
	FullName string `json:"fullname"`
	Type     string `json:"type"`
	IdType   string `json:"idtype"`
	IdNumber string `json:"idnumber"`
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
	TransId           *string `json:"trans_id"`
	ScheduleId        *string `json:"schedule_id"`
}
type MyBookingWithPagination struct {
	Data []*MyBooking    `json:"data"`
	Meta *MetaPagination `json:"meta"`
}
type MyBooking struct {
	OrderId            string                   `json:"order_id"`
	ExpType            []string                 `json:"exp_type"`
	ExpId              string                   `json:"exp_id"`
	ExpTitle           string                   `json:"exp_title"`
	TransId            string                   `json:"trans_id"`
	TransName          string                   `json:"trans_name"`
	TransFrom          string                   `json:"trans_from"`
	TransTo            string                   `json:"trans_to"`
	TransDepartureTime *string                  `json:"trans_departure_time"`
	TransArrivalTime   *string                  `json:"trans_arrival_time"`
	TripDuration       string                   `json:"trip_duration"`
	TransClass         string                   `json:"trans_class"`
	TransGuest         TotalGuestTransportation `json:"trans_guest"`
	BookingDate        time.Time                `json:"booking_date"`
	ExpDuration        int                      `json:"exp_duration"`
	ExpGuest           TotalGuestTransportation `json:"exp_guest"`
	TotalGuest         int                      `json:"total_guest"`
	City               string                   `json:"city"`
	Province           string                   `json:"province"`
	Country            string                   `json:"country"`
}

type TotalGuestTransportation struct {
	Adult    int `json:"adult"`
	Children int `json:"children"`
	Infant 	 int 	`json:"infant"`
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
	ScheduleId         *string    `json:"schedule_id"`
	TotalPrice         float64    `json:"total_price"`
}
