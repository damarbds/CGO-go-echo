package models

import "time"

type Guide struct {
	Id               string     `json:"id" validate:"required"`
	CreatedBy        string     `json:"created_by":"required"`
	CreatedDate      time.Time  `json:"created_date" validate:"required"`
	ModifiedBy       *string    `json:"modified_by"`
	ModifiedDate     *time.Time `json:"modified_date"`
	DeletedBy        *string    `json:"deleted_by"`
	DeletedDate      *time.Time `json:"deleted_date"`
	IsDeleted        int        `json:"is_deleted" validate:"required"`
	IsActive         int        `json:"is_active" validate:"required"`
	GuideName 		string `json:"guide_name"`
	GuidePhoto 		string `json:"guide_photo"`
	GuideGender 	string `json:"guide_gender"`
	IsCertified 	int `json:"is_certified"`
	GuideDesc 		string `json:"guide_desc"`
	LicenceNumber string `json:"licence_number"`
	GuideReview 	float64 `json:"guide_review"`
}
