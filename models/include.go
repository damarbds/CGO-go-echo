package models

import (
	"time"
)

type Include struct {
	Id           int            `json:"id" validate:"required"`
	CreatedBy    string         `json:"created_by" validate:"required"`
	CreatedDate  time.Time      `json:"created_date" validate:"required"`
	ModifiedBy   *string        `json:"modified_by"`
	ModifiedDate *time.Time     `json:"modified_date"`
	DeletedBy    *string        `json:"deleted_by"`
	DeletedDate  *time.Time     `json:"deleted_date"`
	IsDeleted    int            `json:"is_deleted" validate:"required"`
	IsActive     int            `json:"is_active" validate:"required"`
	IncludeName	 string			`json:"include_name"`
	IncludeIcon	 string			`json:"include_icon"`
}

type IncludeDto struct {
	Id           int    `json:"id"`
	IncludeName	 string	`json:"include_name"`
	IncludeIcon  string	`json:"include_icon"`
}
type NewCommandInclude struct {
	Id           int    `json:"id"`
	IncludeName string `json:"include_name"`
	IncludeIcon string `json:"include_icon"`
}
type IncludeDtoWithPagination struct {
	Data []*IncludeDto 	`json:"data"`
	Meta *MetaPagination 	`json:"meta"`
}