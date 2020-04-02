package repository

import (
	"context"
	"database/sql"
	"github.com/service/filter_activity_type"

	"time"


	"github.com/models"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type filterActivityTypeRepository struct {
	Conn *sql.DB
}

func (f filterActivityTypeRepository) GetByExpId(context context.Context, expId string) ([]*models.FilterActivityType, error) {
	panic("implement me")
}

// NewexperienceRepository will create an object that represent the article.Repository interface
func NewFilterActivityTypeRepository(Conn *sql.DB) filter_activity_type.Repository {
	return &filterActivityTypeRepository{Conn}
}
