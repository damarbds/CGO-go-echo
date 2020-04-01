package repository

import (
	"database/sql"
	"github.com/models"
	"github.com/product/experience_add_ons"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type experienceAddOnsRepository struct {
	Conn *sql.DB
}


// NewexperienceRepository will create an object that represent the article.Repository interface
func NewexperienceRepository(Conn *sql.DB) experience_add_ons.Repository {
	return &experienceAddOnsRepository{Conn}
}
func (m *experienceAddOnsRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.ExperienceAddOn, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]*models.ExperienceAddOn, 0)
	for rows.Next() {
		t := new(models.ExperienceAddOn)
		err = rows.Scan(
			&t.Id,
			&t.	CreatedBy,
			&t.	CreatedDate ,
			&t.	ModifiedBy  ,
			&t.	ModifiedDate  ,
			&t.	DeletedBy ,
			&t.	DeletedDate   ,
			&t.	IsDeleted   ,
			&t.	IsActive   ,
			&t.	Name	,
			&t.	Desc ,
			&t.	Currency 	,
			&t.	Amount 			,
			&t.	ExpId			,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (e experienceAddOnsRepository) GetByExpId(ctx context.Context, exp_id string) ([]*models.ExperienceAddOn, error) {
	query := `select * from experience_add_ons where exp_id =? AND is_deleted = 0 AND is_active = 1`

	res, err := e.fetch(ctx, query,exp_id)
	if err != nil {
		return nil, err
	}
	return res, err
}
