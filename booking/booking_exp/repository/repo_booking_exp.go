package repository

import (
	"database/sql"
	"time"

	"github.com/booking/booking_exp"
	guuid "github.com/google/uuid"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type bookingExpRepository struct {
	Conn *sql.DB
}

// NewMysqlArticleRepository will create an object that represent the article.Repository interface
func NewbookingExpRepository(Conn *sql.DB) booking_exp.Repository {
	return &bookingExpRepository{Conn}
}

func (b bookingExpRepository) fetchGrowth(ctx context.Context, query string, args ...interface{}) ([]*models.BookingGrowth, error) {
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

	result := make([]*models.BookingGrowth, 0)
	for rows.Next() {
		t := new(models.BookingGrowth)
		err = rows.Scan(
			&t.Date,
			&t.Count,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (b bookingExpRepository) GetGrowthByMerchantID(ctx context.Context, merchantId string) ([]*models.BookingGrowth, error) {
	query := `
	SELECT
		cast(c.booking_date AS date) as date,
		count(*) as count
	FROM
		merchants a
		JOIN experiences b ON a.id = b.merchant_id
		JOIN booking_exps c ON b.id = c.exp_id
		JOIN transactions d ON c.id = d.booking_exp_id
	WHERE
		a.id = ?
		AND c.status = 1
		AND d.status = 2
	GROUP BY
		cast(c.booking_date AS date)`

	res, err := b.fetchGrowth(ctx, query, merchantId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return res, nil
}

func (b bookingExpRepository) GetByUserID(ctx context.Context, transactionStatus, bookingStatus int, userId string) ([]*models.BookingExpJoin, error) {
	query := `
	SELECT
		a.*,
		b.exp_title,
		b.exp_type,
		b.exp_duration,
		b.exp_pickup_place,
		b.exp_pickup_time,
		t.total_price,
		pm.name AS payment_type,
		city_name AS city,
		province_name AS province,
		country_name AS country,
		c.id as experience_payment_id
	FROM
		booking_exps a
		JOIN experiences b ON a.exp_id = b.id
		JOIN experience_payments c ON b.id = c.exp_id
		JOIN transactions t ON t.booking_exp_id = a.id
		JOIN payment_methods pm ON pm.id = t.payment_method_id
		JOIN harbors h ON b.harbors_id = h.id
		JOIN cities ci ON h.city_id = ci.id
		JOIN provinces p ON ci.province_id = p.id
		JOIN countries co ON p.country_id = co.id
	WHERE
		a.status = ?
		AND a.is_active = 1
		AND a.is_deleted = 0
		AND (t.status = ? OR t.status IS NULL)
		AND a.user_id = ?
	`

	list, err := b.fetch(ctx, query, bookingStatus, transactionStatus, userId)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (b bookingExpRepository) UpdateStatus(ctx context.Context, bookingId string, expiredDatePayment time.Time) error {
	query := `UPDATE booking_exps SET status = 1, expired_date_payment = ? WHERE id = ?`
	stmt, err := b.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, expiredDatePayment, bookingId)
	if err != nil {
		return err
	}

	return nil
}

func (b bookingExpRepository) Insert(ctx context.Context, a *models.BookingExp) (*models.BookingExp, error) {
	id := guuid.New()
	a.Id = id.String()
	query := `INSERT  booking_exps SET id=?,created_by=?,created_date=?,modified_by=?,modified_date=?,deleted_by=?,
				deleted_date=?,is_deleted=?,is_active=?,exp_id=?,order_id=?,guest_desc=?,booked_by=?,booked_by_email=?,
				booking_date=?,user_id=?,status=?,ticket_code=?,ticket_qr_code=?,experience_add_on_id=?`

	stmt, err := b.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	_, error := stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.ExpId, a.OrderId, a.GuestDesc, a.BookedBy,
		a.BookedByEmail, a.BookingDate, a.UserId, a.Status, a.TicketCode, a.TicketQRCode, a.ExperienceAddOnId)
	if error != nil {
		return nil, err
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
			&t.ExpId,
			&t.OrderId,
			&t.GuestDesc,
			&t.BookedBy,
			&t.BookedByEmail,
			&t.BookingDate,
			&t.ExpiredDatePayment,
			&t.UserId,
			&t.Status,
			&t.TicketCode,
			&t.TicketQRCode,
			&t.ExperienceAddOnId,
			&t.ExpTitle,
			&t.ExpType,
			&t.ExpDuration,
			&t.ExpPickupPlace,
			&t.ExpPickupTime,
			&t.TotalPrice,
			&t.PaymentType,
			&t.City,
			&t.Province,
			&t.Country,
			&t.ExperiencePaymentId,
			&t.Currency,
			&t.AccountBank,
			&t.Icon,
			&t.CreatedDateTransaction,
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
	query := `select  a.*, b.exp_title,b.exp_type,b.exp_duration,b.exp_pickup_place,b.exp_pickup_time,t.total_price ,pm.name as payment_type,
			city_name as city, province_name as province, country_name as country, c.id as experience_payment_id,c.currency,pm.desc as account_bank,pm.icon,t.created_date as created_date_transaction from booking_exps a
			join experiences b on a.exp_id = b.id
			join experience_payments c on b.id = c.exp_id
            join transactions t on t.booking_exp_id = a.id            
            join payment_methods pm on pm.id = t.payment_method_id
			JOIN harbors h ON b.harbors_id = h.id
			JOIN cities ci ON h.city_id = ci.id
			JOIN provinces p ON ci.province_id = p.id
			JOIN countries co ON p.country_id = co.id
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
func (b bookingExpRepository) fetchQueryHistory(ctx context.Context, query string, args ...interface{}) ([]*models.BookingExpHistory, error) {
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

	result := make([]*models.BookingExpHistory, 0)
	for rows.Next() {
		t := new(models.BookingExpHistory)
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
			&t.ExpId,
			&t.OrderId,
			&t.GuestDesc,
			&t.BookedBy,
			&t.BookedByEmail,
			&t.BookingDate,
			&t.ExpiredDatePayment,
			&t.UserId,
			&t.Status,
			&t.TicketCode,
			&t.TicketQRCode,
			&t.ExperienceAddOnId,
			&t.ExpTitle,
			&t.ExpType,
			&t.ExpDuration,
			&t.CityName,
			&t.ProvinceName,
			&t.CountryName,
			&t.StatusTransaction,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (b bookingExpRepository) QueryHistoryPer30DaysByUserId(ctx context.Context, userId string) ([]*models.BookingExpHistory, error) {
	query := `select a.*, b.exp_title,b.exp_type,b.exp_duration,d.city_name,e.province_name,f.country_name,g.status as status_transaction from booking_exps a
					join experiences b on a.exp_id = b.id
					join harbors c on b.harbors_id = c.id
					join cities d on c.city_id = d.id
					join provinces e on d.province_id = e.id
					join countries f on e.country_id = f.id
					join transactions g on g.booking_exp_id = a.id
					where a.user_id = ?
					and (g.created_date >= (NOW() - INTERVAL 1 MONTH) or g.modified_date >= (NOW() - INTERVAL 1 MONTH))`

	list, err := b.fetchQueryHistory(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	result := list
	return result, err
}

func (b bookingExpRepository) QueryHistoryPerMonthByUserId(ctx context.Context, userId string, yearMonth string) ([]*models.BookingExpHistory, error) {
	//dtstr1 := "2010-01-23 11:44:20"
	date := yearMonth + "-" + "01" + " 00:00:00"
	//dt,_ := time.Parse(date, dtstr1)
	query := `select a.*, b.exp_title,b.exp_type,b.exp_duration,d.city_name,e.province_name,f.country_name,g.status as status_transaction from booking_exps a
					join experiences b on a.exp_id = b.id 
					join harbors c on b.harbors_id = c.id
					join cities d on c.city_id = d.id 
					join provinces e on d.province_id = e.id 
					join countries f on e.country_id = f.id
					join transactions g on g.booking_exp_id = a.id
					where a.user_id = ?
					and (g.created_date >= ? or g.modified_date >= ?)
`

	list, err := b.fetchQueryHistory(ctx, query, userId, date, date)
	if err != nil {
		return nil, err
	}

	result := list
	return result, err
}
