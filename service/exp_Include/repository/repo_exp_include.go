package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	"github.com/service/exp_Include"
	"github.com/sirupsen/logrus"
)

type expIncludeRepository struct {
	Conn *sql.DB
}



// NewExpPaymentRepository will create an object that represent the exp_payment.repository interface
func NewExpIncludeRepository(Conn *sql.DB) exp_Include.Repository {
	return &expIncludeRepository{Conn}
}

func (m *expIncludeRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.ExperienceIncludeJoin, error) {
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

	result := make([]*models.ExperienceIncludeJoin, 0)
	for rows.Next() {
		t := new(models.ExperienceIncludeJoin)
		err = rows.Scan(
			&t.Id 	,
			&t.ExpId ,
			&t.IncludeId,
			&t.IncludeName,
			&t.IncludeIcon	,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m *expIncludeRepository) GetByExpIdJoin(ctx context.Context, expId string) ([]*models.ExperienceIncludeJoin, error) {
	query := `SELECT ei.*,i.include_name,i.include_icon
				FROM experience_includes ei 
				JOIN includes i ON ei.include_id = i.id
				WHERE ei.exp_id = ? `

	res, err := m.fetch(ctx, query,expId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return res, nil
}

func (m *expIncludeRepository) Insert(ctx context.Context, a *models.ExperienceInclude) error {
	query := `INSERT experience_includes SET exp_id=?,include_id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx,a.ExpId,a.IncludeId)
	if err != nil {
		return err
	}

	return nil
}

func (s *expIncludeRepository) Delete(ctx context.Context, expId string) error {
	query := "DELETE FROM experience_includes WHERE exp_id = ?"

	stmt, err := s.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, expId)
	if err != nil {

		return err
	}

	//rowsAfected, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}
	//
	//if rowsAfected != 1 {
	//	err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", rowsAfected)
	//	return err
	//}

	return nil
}