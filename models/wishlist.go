package models

import (
	"database/sql"
	"time"
)

type Wishlist struct {
	Id           string     `json:"id" validate:"required"`
	CreatedBy    string     `json:"created_by" validate:"required"`
	CreatedDate  time.Time  `json:"created_date" validate:"required"`
	ModifiedBy   *string    `json:"modified_by"`
	ModifiedDate *time.Time `json:"modified_date"`
	DeletedBy    *string    `json:"deleted_by"`
	DeletedDate  *time.Time `json:"deleted_date"`
	IsDeleted    int        `json:"is_deleted" validate:"required"`
	IsActive     int        `json:"is_active" validate:"required"`
	TransId      string     `json:"trans_id"`
	ExpId        string     `json:"exp_id"`
	UserId       string     `json:"user_id"`
}

type WishlistObj struct {
	Id           string         `json:"id" validate:"required"`
	CreatedBy    string         `json:"created_by" validate:"required"`
	CreatedDate  time.Time      `json:"created_date" validate:"required"`
	ModifiedBy   *string        `json:"modified_by"`
	ModifiedDate *time.Time     `json:"modified_date"`
	DeletedBy    *string        `json:"deleted_by"`
	DeletedDate  *time.Time     `json:"deleted_date"`
	IsDeleted    int            `json:"is_deleted" validate:"required"`
	IsActive     int            `json:"is_active" validate:"required"`
	TransId      sql.NullString `json:"trans_id"`
	ExpId        sql.NullString `json:"exp_id"`
	UserId       string         `json:"user_id"`
}

type WishlistIn struct {
	TransID string `json:"trans_id,omitempty"`
	ExpID   string `json:"exp_id,omitempty"`
}

type WishlistOut struct {
	WishlistID  string   `json:"wishlist_id"`
	Type        string   `json:"type"`
	ExpID       string   `json:"exp_id"`
	ExpTitle    string   `json:"exp_title"`
	ExpType     []string `json:"exp_type"`
	Rating      float64  `json:"rating"`
	CountRating int      `json:"count_rating"`
	Currency    string   `json:"currency"`
	Price       float64  `json:"price"`
	PaymentType string   `json:"payment_type"`
	CoverPhoto  string   `json:"cover_photo"`
	ListPhoto 	[]ExpPhotosObj	`json:"list_photo"`
}
type WishlistOutWithPagination struct {
	Data []*WishlistOut `json:"data"`
	Meta *MetaPagination    `json:"meta"`
}
