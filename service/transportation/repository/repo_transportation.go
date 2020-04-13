package repository

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/models"
	"github.com/service/transportation"
	"golang.org/x/net/context"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type transportationRepository struct {
	Conn *sql.DB
}

func NewTransportationRepository(Conn *sql.DB) transportation.Repository {
	return &transportationRepository{Conn}
}

func (t transportationRepository) GetTransCount(ctx context.Context, merchantId string) (int, error) {
	query := `SELECT count(*) as count FROM transportations WHERE merchant_id = ?`

	rows, err := t.Conn.QueryContext(ctx, query, merchantId)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	count, err := checkCount(rows)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	return count, nil
}

func (t transportationRepository) CountFilterSearch(ctx context.Context, query string) (int, error) {
	rows, err := t.Conn.QueryContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	count, err := checkCount(rows)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	return count, nil
}

func checkCount(rows *sql.Rows) (count int, err error) {
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}

func (t transportationRepository) FilterSearch(ctx context.Context, query string, limit, offset int) ([]*models.TransSearch, error) {
	query = query + ` LIMIT ? OFFSET ?`
	res, err := t.fetchSearchTrans(ctx, query, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return res, nil
}

func (t *transportationRepository) fetchSearchTrans(ctx context.Context, query string, args ...interface{}) ([]*models.TransSearch, error) {
	rows, err := t.Conn.QueryContext(ctx, query, args...)
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

	result := make([]*models.TransSearch, 0)
	for rows.Next() {
		t := new(models.TransSearch)
		err = rows.Scan(
			&t.ScheduleId,
			&t.DepartureDate,
			&t.DepartureTime,
			&t.ArrivalTime,
			&t.Price,
			&t.TransId,
			&t.TransName,
			&t.TransImages,
			&t.HarborSourceId,
			&t.HarborSourceName,
			&t.HarborDestId,
			&t.HarborDestName,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (t transportationRepository) Insert(ctx context.Context, a models.Transportation) (*string, error) {
	query := `INSERT transportations SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , 
				deleted_date=? , is_deleted=? , is_active=? , trans_name=?,harbors_source_id=?,harbors_dest_id=?,merchant_id=?,
				trans_capacity=?,trans_title=?,trans_status=?,trans_images=?,return_trans_id=?,boat_details=?,transcoverphoto=?,
				class=?`
	stmt, err := t.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.TransName, a.HarborsSourceId,
		a.HarborsDestId, a.MerchantId, a.TransCapacity, a.TransTitle, a.TransStatus, a.TransImages, a.ReturnTransId,
		a.BoatDetails, a.Transcoverphoto, a.Class)
	if err != nil {
		return nil, err
	}

	//lastID, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}

	//a.Id = lastID
	return &a.Id, nil
}

func (t transportationRepository) Update(ctx context.Context, a models.Transportation) (*string, error) {
	query := `UPDATE transportations SET modified_by=?, modified_date=? , deleted_by=? , 
				deleted_date=? , is_deleted=? , is_active=? , trans_name=?,harbors_source_id=?,harbors_dest_id=?,merchant_id=?,
				trans_capacity=?,trans_title=?,trans_status=?,trans_images=?,return_trans_id=?,boat_details=?,transcoverphoto=?,
				class=? WHERE id=?`
	stmt, err := t.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.ExecContext(ctx, a.ModifiedBy, time.Now(), nil, nil, 0, 1, a.TransName, a.HarborsSourceId,
		a.HarborsDestId, a.MerchantId, a.TransCapacity, a.TransTitle, a.TransStatus, a.TransImages, a.ReturnTransId,
		a.BoatDetails, a.Transcoverphoto, a.Class, a.Id)
	if err != nil {
		return nil, err
	}

	//lastID, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}

	//a.Id = lastID
	return &a.Id, nil
}
