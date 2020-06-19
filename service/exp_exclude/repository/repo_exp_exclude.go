package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	"github.com/service/exp_exclude"
	"github.com/sirupsen/logrus"
)

type expExcludeRepository struct {
	Conn *sql.DB
}



// NewExpPaymentRepository will create an object that represent the exp_payment.repository interface
func NewExpExcludeRepository(Conn *sql.DB) exp_exclude.Repository {
	return &expExcludeRepository{Conn}
}

func (m *expExcludeRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.ExperienceExcludeJoin, error) {
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

	result := make([]*models.ExperienceExcludeJoin, 0)
	for rows.Next() {
		t := new(models.ExperienceExcludeJoin)
		err = rows.Scan(
			&t.Id 	,
			&t.ExpId ,
			&t.ExcludeId,
			&t.ExcludeName,
			&t.ExcludeIcon	,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m *expExcludeRepository) GetByExpIdJoin(ctx context.Context, expId string) ([]*models.ExperienceExcludeJoin, error) {
	query := `SELECT ee.*,e.exclude_name, e.exclude_icon
				FROM experience_excludes ee 
				JOIN excludes e ON ee.exclude_id = e.id
				WHERE ee.exp_id = ? `

	res, err := m.fetch(ctx, query,expId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return res, nil
}

func (m *expExcludeRepository) Insert(ctx context.Context, a *models.ExperienceExclude) error {
	query := `INSERT experience_excludes SET exp_id=?,exclude_id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx,a.ExpId,a.ExcludeId)
	if err != nil {
		return err
	}

	return nil
}

func (s *expExcludeRepository) Delete(ctx context.Context, expId string) error {
	query := "DELETE FROM experience_excludes WHERE exp_id = ?"

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