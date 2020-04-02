package models

import "time"

type Payment struct {
	Id                   string    `json:"id" validate:"required"`
	CreatedBy            string    `json:"created_by":"required"`
	CreatedDate          time.Time `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int       `json:"is_deleted" validate:"required"`
	IsActive             int       `json:"is_active" validate:"required"`
	BookingExpId		string		`json:"booking_exp_id"`
	PromoId				string		`json:"promo_id"`
	PaymentMethodId 	string 		`json:"payment_method_id"`
	ExperiencePaymentId	string		`json:"experience_payment_id"`
	Status 				int 		`json:"status"`
}
