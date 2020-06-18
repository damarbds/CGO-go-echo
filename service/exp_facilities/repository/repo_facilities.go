package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	"github.com/service/exp_facilities"
	"github.com/sirupsen/logrus"
)

type expFacilitiesRepository struct {
	Conn *sql.DB
}



// NewExpPaymentRepository will create an object that represent the exp_payment.repository interface
func NewExpFacilitiesRepository(Conn *sql.DB) exp_facilities.Repository {
	return &expFacilitiesRepository{Conn}
}

func (m *expFacilitiesRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.ExperienceFacilitiesJoin, error) {
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

	result := make([]*models.ExperienceFacilitiesJoin, 0)
	for rows.Next() {
		t := new(models.ExperienceFacilitiesJoin)
		err = rows.Scan(
			&t.Id 	,
			&t.ExpId 	,
			&t.TransId ,
			&t.FacilitiesId ,
			&t.Amount,
			&t.FacilityName,
			&t.IsNumerable  ,
			&t.FacilityIcon ,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m *expFacilitiesRepository) GetJoin(ctx context.Context, expId string,transId string) ([]*models.ExperienceFacilitiesJoin, error) {
	var res []*models.ExperienceFacilitiesJoin
	query := `SELECT ef.* , f.facility_name,f.is_numerable,f.facility_icon
				FROM experience_facilities ef 
				JOIN facilities f ON ef.facilities_id = f.id`
	if expId != ""{
		query = query + " WHERE ef.exp_id = ?"
		result, err := m.fetch(ctx, query,expId)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		res = result
	}else if transId != ""{
		query = query + " WHERE ef.trans_id = ?"
		result, err := m.fetch(ctx, query,transId)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		res = result
	}

	return res, nil
}

func (m *expFacilitiesRepository) Insert(ctx context.Context, a *models.ExperienceFacilities) error {
	query := `INSERT experience_facilities SET exp_id=?,trans_id=?,facilities_id=?,amount=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx,a.ExpId,a.TransId,a.FacilitiesId,a.Amount)
	if err != nil {
		return err
	}

	return nil
}

func (s *expFacilitiesRepository) Delete(ctx context.Context, expId string,transId string) error {
	query := "DELETE FROM experience_facilities WHERE "
	if expId != ""{
		query = query + "exp_id = ?"
		stmt, err := s.Conn.PrepareContext(ctx, query)
		if err != nil {
			return err
		}

		_, err = stmt.ExecContext(ctx, expId)
		if err != nil {

			return err
		}

	}else if transId != ""{
		query = query + "trans_id = ?"
		stmt, err := s.Conn.PrepareContext(ctx, query)
		if err != nil {
			return err
		}

		_, err = stmt.ExecContext(ctx, transId)
		if err != nil {

			return err
		}

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