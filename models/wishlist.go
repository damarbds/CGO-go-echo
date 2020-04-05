package models

import (
	"time"
)

type Wishlist struct {
	Id           string         `json:"id" validate:"required"`
	CreatedBy    string         `json:"created_by" validate:"required"`
	CreatedDate  time.Time      `json:"created_date" validate:"required"`
	ModifiedBy   *string        `json:"modified_by"`
	ModifiedDate *time.Time     `json:"modified_date"`
	DeletedBy    *string        `json:"deleted_by"`
	DeletedDate  *time.Time     `json:"deleted_date"`
	IsDeleted    int            `json:"is_deleted" validate:"required"`
	IsActive     int            `json:"is_active" validate:"required"`
	TransId      string `json:"trans_id"`
	ExpId        string `json:"exp_id"`
	UserId       string         `json:"user_id"`
}

type WishlistIn struct {
	TransID string `json:"trans_id,omitempty"`
	ExpID   string `json:"exp_id,omitempty"`
}
