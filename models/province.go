package models

import "time"

type Province struct {
	Id           int        `json:"id" validate:"required"`
	CreatedBy    string     `json:"created_by":"required"`
	CreatedDate  time.Time  `json:"created_date" validate:"required"`
	ModifiedBy   *string    `json:"modified_by"`
	ModifiedDate *time.Time `json:"modified_date"`
	DeletedBy    *string    `json:"deleted_by"`
	DeletedDate  *time.Time `json:"deleted_date"`
	IsDeleted    int        `json:"is_deleted" validate:"required"`
	IsActive     int        `json:"is_active" validate:"required"`
	ProvinceName string     `json:"province_name"`
	CountryId    int        `json:"country_id"`
	ProvinceNameTransportation *string `json:"province_name_transportation"`
}
type ProvinceDto struct {
	Id           int        `json:"id"`
	ProvinceName string     `json:"province_name"`
	CountryId    int        `json:"country_id"`
}
type NewCommandProvince struct {
	Id           int        `json:"id"`
	ProvinceName string     `json:"province_name"`
	CountryId    int        `json:"country_id"`
	ProvinceNameTransportation *string `json:"province_name_transportation"`
}

type ProvinceDtoWithPagination struct {
	Data []*ProvinceDto `json:"data"`
	Meta *MetaPagination    `json:"meta"`
}