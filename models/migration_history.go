package models

import "time"

type MigrationHistory struct {
	Id            int       `json:"id" validate:"required"`
	DescMigration string    `json:"desc_migration"`
	Date          time.Time `json:"date"`
}
