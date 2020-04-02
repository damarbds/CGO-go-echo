package repository

import (
	"database/sql"
	"github.com/booking/booking_exp"
	guuid "github.com/google/uuid"
	"github.com/models"
	"golang.org/x/net/context"
	"time"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type bookingExpRepository struct {
	Conn *sql.DB
}


// NewMysqlArticleRepository will create an object that represent the article.Repository interface
func NewbookingExpRepository(Conn *sql.DB) booking_exp.Repository{
	return &bookingExpRepository{Conn}
}

func (b bookingExpRepository) Insert(ctx context.Context, a *models.BookingExp) (*models.BookingExp,error) {
	id:= guuid.New()
	a.Id = id.String()
	query := `INSERT  booking_exps SET id=?,created_by=?,created_date=?,modified_by=?,modified_date=?,deleted_by=?,
				deleted_date=?,is_deleted=?,is_active=?,exp_id=?,guest_desc=?,booked_by=?,booked_by_email=?,
				booking_date=?,user_id=?,status=?,ticket_code=?,ticket_qr_code=?`

	stmt, err := b.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil,err
	}

	_, error := stmt.ExecContext(ctx,a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1,a.ExpId,a.GuestDesc,a.BookedBy,
		a.BookedByEmail,a.BookingDate,a.UserId,a.Status,a.TicketCode,a.TicketQRCode)
	if error != nil {
		return nil,err
	}

	//lastID, err := res.LastInsertId()
	//if err != nil {
	//	return err
	//}
	//
	//a.Id = lastID
	return a,nil
}
