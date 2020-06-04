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
}

type NotifDto struct {
	Id    string    `json:"id"`
	Type  string    `json:"type"`
	Title string    `json:"title"`
	Desc  string    `json:"desc"`
	Date  time.Time `json:"date"`
}
