package models

import (
	"time"
)

type ExpInspiration struct {
	Id            string     `json:"id" validate:"required"`
	CreatedBy     string     `json:"created_by" validate:"required"`
	CreatedDate   time.Time  `json:"created_date" validate:"required"`
	ModifiedBy    *string    `json:"modified_by"`
	ModifiedDate  *time.Time `json:"modified_date"`
	DeletedBy     *string    `json:"deleted_by"`
	DeletedDate   *time.Time `json:"deleted_date"`
	IsDeleted     int        `json:"is_deleted" validate:"required"`
	IsActive      int        `json:"is_active" validate:"required"`
	ExpId         string     `json:"exp_id"`
	ExpTitle      string     `json:"exp_title"`
	ExpDesc       string     `json:"exp_desc"`
	ExpCoverPhoto string     `json:"exp_cover_photo"`
}

type ExpInspirationObject struct {
	ExpInspirationID string  `json:"exp_inspiration_id"`
	ExpId            string  `json:"exp_id"`
	ExpTitle         string  `json:"exp_title"`
	ExpDesc          string  `json:"exp_desc"`
	ExpCoverPhoto    string  `json:"exp_cover_photo"`
	ExpType          string  `json:"exp_type"`
	Rating           float64 `json:"rating"`
}
type ExpInspirationDto struct {
	ExpInspirationID string   `json:"exp_inspiration_id"`
	ExpId            string   `json:"exp_id"`
	ExpTitle         string   `json:"exp_title"`
	ExpDesc          string   `json:"exp_desc"`
	ExpCoverPhoto    string   `json:"exp_cover_photo"`
	ExpType          []string `json:"exp_type"`
	Rating           float64  `json:"rating"`
	CountRating      int      `json:"count_rating"`
}
