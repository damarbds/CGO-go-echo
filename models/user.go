package models

import (
	"time"
)

type User struct {
	Id                   string    `json:"id" validate:"required"`
	CreatedBy            string    `json:"created_by" validate:"required"`
	CreatedDate          time.Time `json:"created_date" validate:"required"`
	ModifiedBy           *string    `json:"modified_by"`
	ModifiedDate         *time.Time `json:"modified_date"`
	DeletedBy            *string    `json:"deleted_by"`
	DeletedDate          *time.Time `json:"deleted_date"`
	IsDeleted            int       `json:"is_deleted" validate:"required"`
	IsActive             int       `json:"is_active" validate:"required"`
	UserEmail            string    `json:"user_email" validate:"required"`
	FullName             string    `json:"full_name"`
	PhoneNumber          int       `json:"phone_number" validate:"required"`
	VerificationSendDate time.Time `json:"verification_send_date"`
	VerificationCode     int       `json:"verification_code"`
	ProfilePictUrl       string    `json:"profile_pict_url"`
	Address              string    `json:"address" validate:"required"`
	Dob                  time.Time `json:"dob" validate:"required"`
	Gender               int       `json:"gender" validate:"required"`
	IdType               int       `json:"id_type"`
	IdNumber             string    `json:"id_number"`
	ReferralCode         int       `json:"referral_code"`
	Points               int       `json:"points"`
}
type NewCommandUser struct {
	Id                   string    `json:"id"`
	UserEmail            string    `json:"user_email" validate:"required"`
	Password			 string 	`json:"password"`
	FullName             string    `json:"full_name"`
	PhoneNumber          int       `json:"phone_number" validate:"required"`
	VerificationSendDate string `json:"verification_send_date"`
	VerificationCode     int       `json:"verification_code"`
	ProfilePictUrl       string    `json:"profile_pict_url"`
	Address              string    `json:"address" validate:"required"`
	Dob                  string `json:"dob" validate:"required"`
	Gender               int       `json:"gender" validate:"required"`
	IdType               int       `json:"id_type"`
	IdNumber             string    `json:"id_number"`
	ReferralCode         int       `json:"referral_code"`
	Points               int       `json:"points"`
}
type UserInfoDto struct {
	Id                   string    `json:"id"`
	UserEmail            string    `json:"user_email" validate:"required"`
	FullName             string    `json:"full_name"`
	PhoneNumber          int       `json:"phone_number" validate:"required"`
	ProfilePictUrl       string    `json:"profile_pict_url"`
}