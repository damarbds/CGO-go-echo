package repository

import (
	"database/sql"
	guuid "github.com/google/uuid"
	"github.com/models"
	"github.com/product/experience_add_ons"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
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

func (m *experienceAddOnsRepository) Insert(ctx context.Context, a models.ExperienceAddOn) error {
	id := guuid.New()
	a.Id = id.String()
	query := `INSERT experience_add_ons SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , 
				deleted_by=? , deleted_date=? , is_deleted=? , is_active=? ,name=?,desc=?,currency=?,amount=?,exp_id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.Name,a.Desc,a.Currency,
		a.Amount,a.ExpId)
	if err != nil {
		return err
	}

	//lastID, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}

	//a.Id = lastID
	return nil
}