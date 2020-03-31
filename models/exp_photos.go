package models

import "time"

type ExpPhotos struct {
	Id                   string    `json:"id" validate:"required"`
	CreatedBy            string    `json:"created_by":"required"`
	CreatedDate          time.Time `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int       `json:"is_deleted" validate:"required"`
	IsActive             int       `json:"is_active" validate:"required"`
	ExpPhotoFolder		string		`json:"exp_photo_folder"`
	ExpPhotoImage 		string		`json:"exp_photo_image"`
	ExpId				string 		`json:"exp_id"`
}

type ExpPhotosDto struct {
	Id                   string    `json:"id" validate:"required"`
	ExpPhotoFolder		string		`json:"exp_photo_folder"`
	ExpPhotoImage 		[]ExpPhotoImageObject		`json:"exp_photo_image"`
	ExpId				string 		`json:"exp_id"`
}
type ExpPhotoImageObject struct {
	Original		string 			`json:"original"`
	Thumbnail		string			`json:"thumbnail"`
}