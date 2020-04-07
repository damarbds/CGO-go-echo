package models

import "time"

type Admin struct {
	Id                   string    `json:"id" validate:"required"`
	CreatedBy            string    `json:"created_by":"required"`
	CreatedDate          time.Time `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int       `json:"is_deleted" validate:"required"`
	IsActive             int       `json:"is_active" validate:"required"`
	Name 				string	`json:"name"`
	Email 				string	`json:"email"`
}

type NewCommandAdmin struct {
	Id                   string    `json:"id"`
	Name 				string	`json:"name"`
	Email 				string	`json:"email"`
	Password 			string `json:"password"`
}
type AdminDto struct {
	Id                   string    `json:"id"`
	Name 				string	`json:"name"`
	Email 				string	`json:"email"`
} 
