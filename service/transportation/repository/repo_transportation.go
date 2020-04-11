package repository

import (
	"database/sql"
	"github.com/models"
	"github.com/service/transportation"
	"golang.org/x/net/context"
	"time"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type transportationRepository struct {
	Conn *sql.DB
}



// NewpromoRepository will create an object that represent the article.Repository interface
func NewTransportationRepository(Conn *sql.DB) transportation.Repository {
	return &transportationRepository{Conn}
}

func (t transportationRepository) Insert(ctx context.Context, a models.Transportation) (*string, error) {
	query := `INSERT transportations SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , 
				deleted_date=? , is_deleted=? , is_active=? , trans_name=?,harbors_source_id=?,harbors_dest_id=?,merchant_id=?,
				trans_capacity=?,trans_title=?,trans_status=?,trans_images=?,return_trans_id=?,boat_details=?,transcoverphoto=?,
				class=?`
	stmt, err := t.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil,err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.TransName,a.HarborsSourceId,
		a.HarborsDestId,a.MerchantId,a.TransCapacity,a.TransTitle,a.TransStatus,a.TransImages,a.ReturnTransId,
		a.BoatDetails,a.Transcoverphoto,a.Class)
	if err != nil {
		return nil,err
	}

	//lastID, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}

	//a.Id = lastID
	return &a.Id,nil
}

func (t transportationRepository) Update(ctx context.Context, a models.Transportation) (*string, error) {
	query := `UPDATE transportations SET modified_by=?, modified_date=? , deleted_by=? , 
				deleted_date=? , is_deleted=? , is_active=? , trans_name=?,harbors_source_id=?,harbors_dest_id=?,merchant_id=?,
				trans_capacity=?,trans_title=?,trans_status=?,trans_images=?,return_trans_id=?,boat_details=?,transcoverphoto=?,
				class=?`
	stmt, err := t.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil,err
	}
	_, err = stmt.ExecContext(ctx, a.ModifiedBy, time.Now(), nil, nil, 0, 1, a.TransName,a.HarborsSourceId,
		a.HarborsDestId,a.MerchantId,a.TransCapacity,a.TransTitle,a.TransStatus,a.TransImages,a.ReturnTransId,
		a.BoatDetails,a.Transcoverphoto,a.Class)
	if err != nil {
		return nil,err
	}

	//lastID, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}

	//a.Id = lastID
	return &a.Id,nil
}