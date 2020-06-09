package models

import "time"

type Rule struct {
	Id           int            `json:"id" validate:"required"`
	CreatedBy    string         `json:"created_by" validate:"required"`
	CreatedDate  time.Time      `json:"created_date" validate:"required"`
	ModifiedBy   *string        `json:"modified_by"`
	ModifiedDate *time.Time     `json:"modified_date"`
	DeletedBy    *string        `json:"deleted_by"`
	DeletedDate  *time.Time     `json:"deleted_date"`
	IsDeleted    int            `json:"is_deleted" validate:"required"`
	IsActive     int            `json:"is_active" validate:"required"`
	RuleName	 string			`json:"rule_name"`
	RuleIcon	 string			`json:"rule_icon"`
}

type RuleDto struct {
	Id           int    `json:"id"`
	RuleName	 string	`json:"rule_name"`
	RuleIcon  string	`json:"rule_icon"`
}
type NewCommandRule struct {
	Id           int    `json:"id"`
	RuleName string `json:"rule_name"`
	RuleIcon string `json:"rule_icon"`
}
type RuleDtoWithPagination struct {
	Data []*RuleDto 	`json:"data"`
	Meta *MetaPagination 	`json:"meta"`
}