package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	"github.com/service/filter_activity_type"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type filterActivityTypeRepository struct {
	Conn *sql.DB
}


// NewexperienceRepository will create an object that represent the article.repository interface
func NewFilterActivityTypeRepository(Conn *sql.DB) filter_activity_type.Repository {
	return &filterActivityTypeRepository{Conn}
}
func (m *filterActivityTypeRepository) fetchJoinExpType(ctx context.Context, query string, args ...interface{}) ([]*models.FilterActivityTypeJoin, error) {
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

	result := make([]*models.FilterActivityTypeJoin, 0)
	for rows.Next() {
		t := new(models.FilterActivityTypeJoin)
		err = rows.Scan(
			&t.Id      ,
			&t.CreatedBy   ,
			&t.CreatedDate  ,
			&t.ModifiedBy   ,
			&t.ModifiedDate ,
			&t.DeletedBy    ,
			&t.DeletedDate  ,
			&t.IsDeleted    ,
			&t.IsActive     ,
			&t.ExpTypeId    ,
			&t.ExpId        ,
			&t.ExpTypeName  ,
			&t.ExpTypeIcon  ,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m filterActivityTypeRepository) GetJoinExpType(ctx context.Context, expId string) ([]*models.FilterActivityTypeJoin,error) {
	query := `SELECT fat.*,et.exp_type_name,et.exp_type_icon 
				FROM filter_activity_types fat 
				JOIN experience_types et ON fat.exp_type_id = et.id
				WHERE fat.exp_id = ?`

	result, err := m.fetchJoinExpType(ctx, query,expId)
	if err != nil {
			logrus.Error(err)
			return nil, err
	}
	return result,nil
}
func (m filterActivityTypeRepository) DeleteByExpId(ctx context.Context, expId string) error {
	query := "DELETE FROM filter_activity_types WHERE exp_id = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
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
func (m filterActivityTypeRepository) Insert(ctx context.Context, a *models.FilterActivityType) error {
	query := `INSERT filter_activity_types SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? ,
				deleted_by=? , deleted_date=? , is_deleted=? , is_active=? , exp_type_id=? , exp_id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.ExpTypeId, a.ExpId)
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
