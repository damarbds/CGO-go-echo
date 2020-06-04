package repository

import (
	"database/sql"
	"time"

	guuid "github.com/google/uuid"
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


// NewexperienceRepository will create an object that represent the article.repository interface
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
			&t.CreatedBy,
			&t.CreatedDate,
			&t.ModifiedBy,
			&t.ModifiedDate,
			&t.DeletedBy,
			&t.DeletedDate,
			&t.IsDeleted,
			&t.IsActive,
			&t.Name,
			&t.Desc,
			&t.Currency,
			&t.Amount,
			&t.ExpId,
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

	res, err := e.fetch(ctx, query, exp_id)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (m *experienceAddOnsRepository) Insert(ctx context.Context, a models.ExperienceAddOn) (string, error) {
	id := guuid.New()
	a.Id = id.String()
	query := `INSERT experience_add_ons SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , 
				deleted_by=? , deleted_date=? , is_deleted=? , is_active=? ,experience_add_ons.name=? , 
				experience_add_ons.desc = ? , currency=? , amount=?,exp_id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return "", err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.Name, a.Desc, a.Currency,
		a.Amount, a.ExpId)
	if err != nil {
		return "", err
	}

	//lastID, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}

	//a.Id = lastID
	return a.Id, nil
}

func (m *experienceAddOnsRepository) Update(ctx context.Context, a models.ExperienceAddOn) error {
	query := `UPDATE experience_add_ons SET modified_by=?, modified_date=? , 
				deleted_by=? , deleted_date=? , is_deleted=? , is_active=? ,experience_add_ons.name=? , 
				experience_add_ons.desc = ?,currency=?,amount=?,exp_id=?
				WHERE id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.ModifiedBy, time.Now(), nil, nil, 0, 1, a.Name, a.Desc, a.Currency,
		a.Amount, a.ExpId, a.Id)
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
func (m *experienceAddOnsRepository) Deletes(ctx context.Context, ids []string, expId string, deletedBy string) error {
	query := `UPDATE  experience_add_ons SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE exp_id=?`
	for index, id := range ids {
		if index == 0 && index != (len(ids)-1) {
			query = query + ` AND (id !=` + id
		} else if index == 0 && index == (len(ids)-1) {
			query = query + ` AND (id !=` + id + ` ) `
		} else if index == (len(ids) - 1) {
			query = query + ` OR id !=` + id + ` ) `
		} else {
			query = query + ` OR id !=` + id
		}
	}
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, deletedBy, time.Now(), 1, 0, expId)
	if err != nil {
		return err
	}

	//lastID, err := res.RowsAffected()
	if err != nil {
		return err
	}

	//a.Id = lastID
	return nil
}


func (m *experienceAddOnsRepository) DeleteByExpId(ctx context.Context, expId string, deletedBy string) error {
	query := `UPDATE experience_add_ons SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE exp_id=?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, deletedBy, time.Now(), 1, 0, expId)
	if err != nil {
		return err
	}

	//lastID, err := res.RowsAffected()
	if err != nil {
		return err
	}

	//a.Id = lastID
	return nil
}
