package models

import "time"

type Notification struct {
	Id           string     `json:"id" validate:"required"`
	CreatedBy    string     `json:"created_by" validate:"required"`
	CreatedDate  time.Time  `json:"created_date" validate:"required"`
	ModifiedBy   *string    `json:"modified_by"`
	ModifiedDate *time.Time `json:"modified_date"`
	DeletedBy    *string    `json:"deleted_by"`
	DeletedDate  *time.Time `json:"deleted_date"`
	IsDeleted    int        `json:"is_deleted" validate:"required"`
	IsActive     int        `json:"is_active" validate:"required"`
	MerchantId   string     `json:"merchant_id"`
	Type         int        `json:"type"`
	Title        string     `json:"title"`
	Desc         string     `json:"desc"`
	ExpId 		*string	`json:"exp_id"`
	ScheduleId  *string `json:"schedule_id"`
	BookingExpId *string	`json:"booking_exp_id"`
	IsRead 		*int `json:"is_read"`
}

type NotifDto struct {
	Id    string    `json:"id"`
	Type  string    `json:"type"`
	Title string    `json:"title"`
	Desc  string    `json:"desc"`
	Date  time.Time `json:"date"`
	OrderId *string `json:"order_id"`
	ExpId *string	`json:"exp_id"`
	ExpTitle *string `json:"exp_title"`
	TransId *string `json:"trans_id"`
	TransName *string `json:"trans_name"`
	DepartureTime    *string    `json:"departure_time"`
	ArrivalTime      *string    `json:"arrival_time"`
	TripDuration     *string    `json:"trip_duration"`
	HarborSourceName *string    `json:"harbor_source_name"`
	HarborDestName   *string    `json:"harbor_dest_name"`
}

type NotifWithPagination struct {
	Data []*NotifDto `json:"data"`
	Meta *MetaPagination   `json:"meta"`
}

type FCMPushNotif struct {
	To 	 string `json:"to"`
	Data DataFCMPushNotif `json:"data"`
}
type DataFCMPushNotif struct {
	Title string `json:"title"`
	Message string `json:"message"`
}

type NotificationRead struct {
	NotificationId string `json:"notification_id"`
	IsRead int `json:"is_read"`
}