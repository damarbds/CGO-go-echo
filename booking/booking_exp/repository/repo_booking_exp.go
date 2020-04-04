package repository

import (
	"database/sql"
	"github.com/booking/booking_exp"
	guuid "github.com/google/uuid"
	"github.com/models"
	"github.com/sirupsen/logrus"
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
				deleted_date=?,is_deleted=?,is_active=?,exp_id=?,order_id=?,guest_desc=?,booked_by=?,booked_by_email=?,
				booking_date=?,user_id=?,status=?,ticket_code=?,ticket_qr_code=?,experience_add_on_id=?`

	stmt, err := b.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil,err
	}

	_, error := stmt.ExecContext(ctx,a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1,a.ExpId,a.OrderId,a.GuestDesc,a.BookedBy,
		a.BookedByEmail,a.BookingDate,a.UserId,a.Status,a.TicketCode,a.TicketQRCode,a.ExperienceAddOnId)
	if error != nil {
		return nil,err
	}

	return a, nil
}

func (b bookingExpRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.BookingExpJoin, error) {
	rows, err := b.Conn.QueryContext(ctx, query, args...)
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

	result := make([]*models.BookingExpJoin, 0)
	for rows.Next() {
		t := new(models.BookingExpJoin)
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
			&t.ExpId				,
			&t.OrderId		,
			&t.GuestDesc	,
			&t.BookedBy	,
			&t.BookedByEmail	,
			&t.BookingDate 	,
			&t.UserId			,
			&t.Status 			,
			&t.TicketCode		,
			&t.TicketQRCode	,
			&t.ExperienceAddOnId ,
			&t.ExpTitle		,
			&t.ExpType 		,
			&t.ExpPickupPlace 	,
			&t.ExpPickupTime 	,
			&t.TotalPrice 		,
			&t.PaymentType 	,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (b bookingExpRepository) GetEmailByID(ctx context.Context, bookingId string) (string, error) {
	var email string
	query := `SELECT booked_by_email as email FROM booking_exps WHERE id = ?`

	err := b.Conn.QueryRowContext(ctx, query, bookingId).Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", models.ErrNotFound
		}
		return "", err
	}

	return email, err
}

func (b bookingExpRepository) GetDetailBookingID(ctx context.Context, bookingId string) (*models.BookingExpJoin, error) {
	var booking *models.BookingExpJoin
	query := `select  a.*, b.exp_title,b.exp_type,b.exp_pickup_place,b.exp_pickup_time,t.total_price ,pm.name as payment_type from booking_exps a
			join experiences b on a.exp_id = b.id
			join experience_payments c on b.id = c.exp_id
            join transactions t on t.booking_exp_id = a.id            
            join payment_methods pm on pm.id = t.payment_method_id
            where a.id = ? AND a.is_active = 1 AND a.is_deleted = 0`

	list, err := b.fetch(ctx, query, bookingId)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		booking = list[0]
	} else {
		return nil, models.ErrNotFound
	}



	return booking, err
}

