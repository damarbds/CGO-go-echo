package models

import "time"

type Promo struct {
	Id                   string    `json:"id" validate:"required"`
	CreatedBy            string    `json:"created_by":"required"`
	CreatedDate          time.Time `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int       `json:"is_deleted" validate:"required"`
	IsActive             int       `json:"is_active" validate:"required"`
	PromoCode			string		`json:"promo_code"`
	PromoName 			string		`json:"promo_name"`
	PromoDesc 			string		`json:"promo_desc"`
	PromoValue 			float64		`json:"promo_value"`
	PromoType			int			`json:"promo_type"`
	PromoImage 			string		`json:"promo_image"`
}
