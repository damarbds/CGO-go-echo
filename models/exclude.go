package models

import "time"

type Exclude struct {
	Id           int            `json:"id" validate:"required"`
	CreatedBy    string         `json:"created_by" validate:"required"`
	CreatedDate  time.Time      `json:"created_date" validate:"required"`
	ModifiedBy   *string        `json:"modified_by"`
	ModifiedDate *time.Time     `json:"modified_date"`
	DeletedBy    *string        `json:"deleted_by"`
	DeletedDate  *time.Time     `json:"deleted_date"`
	IsDeleted    int            `json:"is_deleted" validate:"required"`
	IsActive     int            `json:"is_active" validate:"required"`
	ExcludeName	 string			`json:"include_name"`
	ExcludeIcon	 string			`json:"include_icon"`
}

type ExcludeDto struct {
	Id 			int		`json:"id"`
	ExcludeName	string	`json:"exclude_name"`
	ExcludeIcon	string	`json:"exclude_icon"`
}

type NewCommandExclude struct {
	Id 			int		`json:"id"`
	ExcludeName	string	`json:"exclude_name"`
	ExcludeIcon	string	`json:"exclude_icon"`
}

type ExcludeDtoWithPagination struct {
	Data []*ExcludeDto 		`json:"data"`
	Meta *MetaPagination	`json:"meta"`
}