package models

import "time"

type Transportation struct {
	Id                   string     `json:"id" validate:"required"`
	CreatedBy            string     `json:"created_by" validate:"required"`
	CreatedDate          time.Time  `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int        `json:"is_deleted" validate:"required"`
	IsActive             int        `json:"is_active" validate:"required"`
	TransName 			string		`json:"trans_name"`
	HarborsSourceId		string		`json:"harbors_source_id"`
	HarborsDestId		string		`json:"harbors_dest_id"`
	MerchantId			string		`json:"merchant_id"`
	TransCapacity		int 		`json:"trans_capacity"`
	TransTitle			string 		`json:"trans_title"`
	TransStatus			int			`json:"trans_status"`
	TransImages 		string		`json:"trans_images"`
	ReturnTransId		string		`json:"return_trans_id"`
}
