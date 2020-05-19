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

// NewMysqlArticleRepository will create an object that represent the article.repository interface
func NewbookingExpRepository(Conn *sql.DB) booking_exp.Repository {
	return &bookingExpRepository{Conn}
}

func (b bookingExpRepository) GetDetailTransportBookingID(ctx context.Context, bookingId, bookingCode string) ([]*models.BookingExpJoin, error) {
	q := `
	SELECT
		a.*,
		b.trans_name,
		b.trans_title,
		b.trans_status,
		b.class AS trans_class,
		s.departure_date,
		s.departure_time,
		s.arrival_time,
		t.total_price,
		pm.name AS payment_type,
		pm.desc AS account_bank,
		pm.icon,
		t.status AS transaction_status,
		t.va_number,
		t.created_date AS created_date_transaction,
		m.merchant_name,
		m.phone_number AS merchant_phone,
		m.merchant_picture,
		hs.harbors_name AS harbor_source_name,
		h.harbors_name AS harbor_dest_name,
		t.ex_change_rates,
		t.ex_change_currency
	FROM
		booking_exps a
		LEFT JOIN transactions t ON t.booking_exp_id = a.id
			OR t.order_id = a.order_id
		LEFT JOIN payment_methods pm ON t.payment_method_id = pm.id
		JOIN transportations b ON a.trans_id = b.id
		LEFT JOIN schedules s ON b.id = s.trans_id
			AND a.schedule_id = s.id
		JOIN harbors h ON b.harbors_dest_id = h.id
		JOIN harbors hs ON b.harbors_source_id = hs.id
		JOIN merchants m ON b.merchant_id = m.id
	WHERE
		a.is_active = 1
		AND a.is_deleted = 0
		AND(a.id = ?
			OR a.order_id = ?)`

	list, err := b.fetchDetailTransport(ctx, q, bookingId, bookingCode)
	if err != nil {
		return nil, err
	}

	return list, err
}

func (b bookingExpRepository) CheckBookingCode(ctx context.Context, bookingCode string) bool {
	var code string
	query := `SELECT order_id as code FROM booking_exps WHERE order_id = ?`

	_ = b.Conn.QueryRowContext(ctx, query, bookingCode).Scan(&code)

	if bookingCode == code {
		return true
	}

	return false
}

func (b bookingExpRepository) GetByID(ctx context.Context, bookingId string) (*models.BookingTransactionExp, error) {
	query := `SELECT a.*, t.total_price from booking_exps a JOIN transactions t ON t.booking_exp_id = a.id where (a.id = ? OR a.order_id = ?)`

	list, err := b.fetchBooking(ctx, query, bookingId, bookingId)
	if err != nil {
		return nil, err
	}

	if len(list) < 1 {
		query = `SELECT a.*, t.total_price from booking_exps a LEFT JOIN transactions t ON t.order_id = a.order_id where a.order_id = ?`

		list, err = b.fetchBooking(ctx, query, bookingId)
		if err != nil {
			return nil, err
		}
	}

	return list[0], nil
}

func (b bookingExpRepository) fetchBooking(ctx context.Context, query string, args ...interface{}) ([]*models.BookingTransactionExp, error) {
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

	result := make([]*models.BookingTransactionExp, 0)
	for rows.Next() {
		t := new(models.BookingTransactionExp)
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
			&t.TransId,
			&t.PaymentUrl,
			&t.ScheduleId,
			&t.TotalPrice,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (b bookingExpRepository) UpdatePaymentUrl(ctx context.Context, bookingId, paymentUrl string) error {
	query := `UPDATE booking_exps SET payment_url = ? WHERE (id = ? OR order_id = ?)`

	stmt, err := b.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, paymentUrl, bookingId, bookingId)
	if err != nil {
		return err
	}

	return nil
}

func (b bookingExpRepository) CountThisMonth(ctx context.Context) (int, error) {
	query := `
	SELECT
		count(CAST(created_date AS DATE)) as count
	FROM
		booking_exps
	WHERE
		is_deleted = 0
		AND is_active = 1
		AND status = 1
		AND created_date BETWEEN date_add(CURRENT_DATE, interval - DAY(CURRENT_DATE) + 1 DAY)
		AND CURRENT_DATE`

	rows, err := b.Conn.QueryContext(ctx, query)
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

func (b bookingExpRepository) GetCountByBookingDateExp(ctx context.Context, bookingDate string, expId string) (int, error) {
	query := `
	SELECT
		count(*) as count
	FROM
		booking_exps
	WHERE
		is_deleted = 0
		AND is_active = 1
		AND (status = 1 OR status = 3)
		AND date(booking_date) = ?
		AND exp_id =?`

	rows, err := b.Conn.QueryContext(ctx, query, bookingDate, expId)
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

func (b bookingExpRepository) GetBookingTransByUserID(ctx context.Context, bookingIds []*string) ([]*models.BookingExpJoin, error) {
	query := `
	
	SELECT
		a.*,
		b.trans_name,
		b.trans_title,
		b.trans_status,
		b.class AS trans_class,
		s.departure_date,
		s.departure_time,
		s.arrival_time,
		t.total_price,
		pm.name AS payment_type,
		pm.desc AS account_bank,
		pm.icon,
		t.status AS transaction_status,
		t.created_date AS created_date_transaction,
		m.merchant_name,
		m.phone_number AS merchant_phone,
		m.merchant_picture,
		hs.harbors_name AS harbor_source_name,
		h.harbors_name AS harbor_dest_name,
		t.ex_change_rates,
		t.ex_change_currency
	FROM
		booking_exps a
		LEFT JOIN transactions t ON t.booking_exp_id = a.id
			OR t.order_id = a.order_id
		LEFT JOIN payment_methods pm ON t.payment_method_id = pm.id
		JOIN transportations b ON a.trans_id = b.id
		LEFT JOIN schedules s ON b.id = s.trans_id
			AND a.schedule_id = s.id
		JOIN harbors h ON b.harbors_dest_id = h.id
		JOIN harbors hs ON b.harbors_source_id = hs.id
		JOIN merchants m ON b.merchant_id = m.id`

	result := make([]*models.BookingExpJoin, 0)
	if len(bookingIds) != 0 {
		for index, id := range bookingIds {
			if index == 0 && index != (len(bookingIds)-1) {
				query = query + ` AND (a.id = '` + *id + `' `
			} else if index == 0 && index == (len(bookingIds)-1) {
				query = query + ` AND (a.id = '` + *id + `' ) `
			} else if index == (len(bookingIds) - 1) {
				query = query + ` OR  a.id = '` + *id + `' ) `
			} else {
				query = query + ` OR  a.id = '` + *id + `' `
			}
		}

		list, err := b.fetchDetailTransport(ctx, query)
		if err != nil {
			return nil, err
		}
		result = list
	}
	return result, nil
}

func (b bookingExpRepository) GetBookingExpByUserID(ctx context.Context, bookingIds []*string) ([]*models.BookingExpJoin, error) {

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
		h.harbors_name,
		city_name AS city,
		province_name AS province,
		country_name AS country,
		a.id AS experience_payment_id,
		a.status as currency,
		pm.desc AS account_bank,
		pm.icon,
		t.status AS transaction_status,
		t.va_number,
		t.created_date AS created_date_transaction,
		m.merchant_name,
		m.phone_number as merchant_phone,
		m.merchant_picture,
		t.ex_change_rates,
		t.ex_change_currency
	FROM
		booking_exps a
		JOIN experiences b ON a.exp_id = b.id
		JOIN transactions t ON t.booking_exp_id = a.id
		JOIN payment_methods pm ON pm.id = t.payment_method_id
		JOIN harbors h ON b.harbors_id = h.id
		JOIN cities ci ON h.city_id = ci.id
		JOIN provinces p ON ci.province_id = p.id
		JOIN countries co ON p.country_id = co.id
		JOIN merchants m ON b.merchant_id = m.id`

	result := make([]*models.BookingExpJoin, len(bookingIds))
	if len(bookingIds) != 0 {
		for index, id := range bookingIds {
			if index == 0 && index != (len(bookingIds)-1) {
				query = query + ` AND (a.id = '` + *id + `' `
			} else if index == 0 && index == (len(bookingIds)-1) {
				query = query + ` AND (a.id = '` + *id + `' ) `
			} else if index == (len(bookingIds) - 1) {
				query = query + ` OR  a.id = '` + *id + `' ) `
			} else {
				query = query + ` OR  a.id = '` + *id + `' `
			}
		}

		list, err := b.fetch(ctx, query)
		if err != nil {
			return nil, err
		}
		result = list
	}

	return result, nil
}

func (b bookingExpRepository) GetBookingCountByUserID(ctx context.Context, status string, userId string) (int, error) {

	query := `
	SELECT COUNT(*) as count 
	FROM
		booking_exps a
		JOIN experiences b ON a.exp_id = b.id
		JOIN transactions t ON t.booking_exp_id = a.id
	WHERE
		a.is_active = 1
		AND a.is_deleted = 0
		AND a.user_id = ?`

	if status == "confirm" {
		//bookingStatus = 1
		query = query + ` 	AND t.status = 2 
							AND DATE(a.booking_date) >= CURRENT_DATE `
	} else if status == "waiting" {
		//transactionStatus = 1
		//bookingStatus = 1
		query = query + ` 	AND t.status IN (1,5) 
							AND DATE(a.booking_date) >= CURRENT_DATE `
	} else if status == "pending" {
		//transactionStatus = 0
		//bookingStatus = 1
		query = query + ` 	AND t.status = 0 
							AND DATE(a.booking_date) >= CURRENT_DATE 
							AND a.expired_date_payment < CURRENT_TIMESTAMP `
	}

	rows, err := b.Conn.QueryContext(ctx, query, userId)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	resultCount, err := checkCount(rows)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	return resultCount, nil
}

func (b bookingExpRepository) GetBookingIdByUserID(ctx context.Context, status string, userId string, limit, offset int) ([]*string, error) {

	var bookingIds []*string
	query := `SELECT DISTINCT a.id
					FROM 
						booking_exps a
					JOIN experiences b ON a.exp_id = b.id
					JOIN transactions t ON t.booking_exp_id = a.id
					WHERE
						a.is_active = 1
					AND a.is_deleted = 0
					AND a.user_id = ?`

	if status == "confirm" {
		//bookingStatus = 1
		query = query + ` 	AND t.status = 2 
							AND DATE(a.booking_date) >= CURRENT_DATE `
	} else if status == "waiting" {
		//transactionStatus = 1
		//bookingStatus = 1
		query = query + ` 	AND t.status IN (1,5) 
							AND DATE(a.booking_date) >= CURRENT_DATE `
	} else if status == "pending" {
		//transactionStatus = 0
		//bookingStatus = 1
		query = query + ` 	AND t.status = 0 
							AND DATE(a.booking_date) >= CURRENT_DATE
							AND DATE_SUB(a.expired_date_payment, INTERVAL 7 HOUR) > DATE_ADD(NOW(), INTERVAL 7 HOUR)`
	}

	if limit != 0 {
		query = query + ` LIMIT ? OFFSET ?`

		rows, err := b.Conn.QueryContext(ctx, query, userId, limit, offset)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		for rows.Next() {
			id := new(string)
			err = rows.Scan(&id)
			if err != nil {
				return nil, err
			}
			bookingIds = append(bookingIds, id)
		}
	} else {
		rows, err := b.Conn.QueryContext(ctx, query, userId)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		for rows.Next() {
			id := new(string)
			err = rows.Scan(&id)
			if err != nil {
				return nil, err
			}
			bookingIds = append(bookingIds, id)
		}
	}
	return bookingIds, nil
}
func (b bookingExpRepository) UpdateStatus(ctx context.Context, bookingId string, expiredDatePayment time.Time) error {
	query := `UPDATE booking_exps SET status = 1, expired_date_payment = ? WHERE (id = ? OR order_id = ?)`

	stmt, err := b.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, expiredDatePayment, bookingId, bookingId)
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
				booking_date=?,user_id=?,status=?,ticket_code=?,ticket_qr_code=?,experience_add_on_id=?,payment_url=?,trans_id=?,schedule_id=?`

	stmt, err := b.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.ExpId, a.OrderId, a.GuestDesc, a.BookedBy,
		a.BookedByEmail, a.BookingDate, a.UserId, a.Status, a.TicketCode, a.TicketQRCode, a.ExperienceAddOnId, a.PaymentUrl, a.TransId, a.ScheduleId)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (b bookingExpRepository) fetchDetailTransport(ctx context.Context, query string, args ...interface{}) ([]*models.BookingExpJoin, error) {
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
			&t.TransId,
			&t.PaymentUrl,
			&t.ScheduleId,
			&t.TransName,
			&t.TransTitle,
			&t.TransStatus,
			&t.TransClass,
			&t.DepartureDate,
			&t.DepartureTime,
			&t.ArrivalTime,
			&t.TotalPrice,
			&t.PaymentType,
			&t.AccountBank,
			&t.Icon,
			&t.TransactionStatus,
			&t.VaNumber,
			&t.CreatedDateTransaction,
			&t.MerchantName,
			&t.MerchantPhone,
			&t.MerchantPicture,
			&t.HarborSourceName,
			&t.HarborDestName,
			&t.ExChangeRates,
			&t.ExChangeCurrency,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
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
			&t.TransId,
			&t.PaymentUrl,
			&t.ScheduleId,
			&t.ExpTitle,
			&t.ExpType,
			&t.ExpDuration,
			&t.ExpPickupPlace,
			&t.ExpPickupTime,
			&t.TotalPrice,
			&t.PaymentType,
			&t.HarborsName,
			&t.City,
			&t.Province,
			&t.Country,
			&t.ExperiencePaymentId,
			&t.Currency,
			&t.AccountBank,
			&t.Icon,
			&t.TransactionStatus,
			&t.VaNumber,
			&t.CreatedDateTransaction,
			&t.MerchantName,
			&t.MerchantPhone,
			&t.MerchantPicture,
			&t.ExChangeRates,
			&t.ExChangeCurrency,
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
	query := `SELECT booked_by_email as email FROM booking_exps WHERE (id = ? OR order_id = ?) LIMIT 1`

	err := b.Conn.QueryRowContext(ctx, query, bookingId, bookingId).Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", models.ErrNotFound
		}
		return "", err
	}

	return email, err
}

func (b bookingExpRepository) GetDetailBookingID(ctx context.Context, bookingId, bookingCode string) (*models.BookingExpJoin, error) {
	var booking *models.BookingExpJoin
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
		h.harbors_name,
		city_name AS city,
		province_name AS province,
		country_name AS country,
		c.id AS experience_payment_id,
		c.currency,
		pm.desc AS account_bank,
		pm.icon,
		t.status AS transaction_status,
		t.va_number,
		t.created_date AS created_date_transaction,
		m.merchant_name,
		m.phone_number as merchant_phone,
		m.merchant_picture,
		t.ex_change_rates,
		t.ex_change_currency
	FROM
		booking_exps a
		JOIN experiences b ON a.exp_id = b.id
		JOIN transactions t ON t.booking_exp_id = a.id		
		JOIN experience_payments c ON t.experience_payment_id = c.id
		JOIN payment_methods pm ON pm.id = t.payment_method_id
		JOIN harbors h ON b.harbors_id = h.id
		JOIN cities ci ON h.city_id = ci.id
		JOIN provinces p ON ci.province_id = p.id
		JOIN countries co ON p.country_id = co.id
		JOIN merchants m ON b.merchant_id = m.id
	WHERE
		a.is_active = 1
		AND a.is_deleted = 0
		AND(a.id = ?
			OR a.order_id = ?)
		AND c.is_active = 1
        AND c.is_deleted = 0`

	list, err := b.fetch(ctx, query, bookingId, bookingCode)
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
func (b bookingExpRepository) fetchQueryExpHistory(ctx context.Context, query string, args ...interface{}) ([]*models.BookingExpHistory, error) {
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
			&t.TransId,
			&t.PaymentUrl,
			&t.ScheduleId,
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
func (b bookingExpRepository) QueryHistoryPer30DaysExpByUserId(ctx context.Context, bookingIds []*string) ([]*models.BookingExpHistory, error) {
	query := `
		select a.*, 
				b.exp_title,
				b.exp_type,
				b.exp_duration,
				d.city_name,
				e.province_name,
				f.country_name,
				g.status as status_transaction 
				from 
					booking_exps a
					join experiences b on a.exp_id = b.id
					join harbors c on b.harbors_id = c.id
					join cities d on c.city_id = d.id
					join provinces e on d.province_id = e.id
					join countries f on e.country_id = f.id
					join transactions g on g.booking_exp_id = a.id`

	result := make([]*models.BookingExpHistory, 0)

	if len(bookingIds) != 0 {
		for index, id := range bookingIds {
			if index == 0 && index != (len(bookingIds)-1) {
				query = query + ` AND (a.id = '` + *id + `' `
			} else if index == 0 && index == (len(bookingIds)-1) {
				query = query + ` AND (a.id = '` + *id + `' ) `
			} else if index == (len(bookingIds) - 1) {
				query = query + ` OR  a.id = '` + *id + `' ) `
			} else {
				query = query + ` OR  a.id = '` + *id + `' `
			}
		}

		list, err := b.fetchQueryExpHistory(ctx, query)
		if err != nil {
			return nil, err
		}
		result = list
	}
	return result, nil
}

func (b bookingExpRepository) QueryHistoryPerMonthExpByUserId(ctx context.Context, bookingIds []*string) ([]*models.BookingExpHistory, error) {

	query := `select a.*, 
				b.exp_title,
				b.exp_type,
				b.exp_duration,
				d.city_name,
				e.province_name,
				f.country_name,
				g.status as status_transaction 
			from 
					booking_exps a
					join experiences b on a.exp_id = b.id 
					join harbors c on b.harbors_id = c.id
					join cities d on c.city_id = d.id 
					join provinces e on d.province_id = e.id 
					join countries f on e.country_id = f.id
					join transactions g on g.booking_exp_id = a.id
`
	result := make([]*models.BookingExpHistory, 0)
	if len(bookingIds) != 0 {
		for index, id := range bookingIds {
			if index == 0 && index != (len(bookingIds)-1) {
				query = query + ` AND (a.id = '` + *id + `' `
			} else if index == 0 && index == (len(bookingIds)-1) {
				query = query + ` AND (a.id = '` + *id + `' ) `
			} else if index == (len(bookingIds) - 1) {
				query = query + ` OR  a.id = '` + *id + `' ) `
			} else {
				query = query + ` OR  a.id = '` + *id + `' `
			}
		}

		list, err := b.fetchQueryExpHistory(ctx, query)
		if err != nil {
			return nil, err
		}
		result = list
	}
	return result, nil
}

func (b bookingExpRepository) QueryHistoryPer30DaysTransByUserId(ctx context.Context, bookingIds []*string) ([]*models.BookingExpJoin, error) {
	query := `
		
	SELECT
		a.*,
		b.trans_name,
		b.trans_title,
		b.trans_status,
		b.class AS trans_class,
		s.departure_date,
		s.departure_time,
		s.arrival_time,
		t.total_price,
		(SELECT pm.name FROM payment_methods pm WHERE id = t.payment_method_id LIMIT 0,1) AS payment_type,
		(SELECT pm.desc FROM payment_methods pm WHERE id = t.payment_method_id LIMIT 0,1) AS account_bank,
		(SELECT pm.icon FROM payment_methods pm WHERE id = t.payment_method_id LIMIT 0,1),
		t.status AS transaction_status,
		t.va_number ,
		t.created_date AS created_date_transaction,
		m.merchant_name,
		m.phone_number AS merchant_phone,
		m.merchant_picture,
		hs.harbors_name AS harbor_source_name,
		h.harbors_name AS harbor_dest_name,
		t.ex_change_rates,
		t.ex_change_currency
	FROM
		booking_exps a
		LEFT JOIN transactions t ON t.booking_exp_id = a.id
			OR t.order_id = a.order_id
		JOIN transportations b ON a.trans_id = b.id
		LEFT JOIN schedules s ON b.id = s.trans_id
			AND a.schedule_id = s.id
		JOIN harbors h ON b.harbors_dest_id = h.id
		JOIN harbors hs ON b.harbors_source_id = hs.id
		JOIN merchants m ON b.merchant_id = m.id`

	result := make([]*models.BookingExpJoin, 0)
	if len(bookingIds) != 0 {
		for index, id := range bookingIds {
			if index == 0 && index != (len(bookingIds)-1) {
				query = query + ` AND (a.id = '` + *id + `' `
			} else if index == 0 && index == (len(bookingIds)-1) {
				query = query + ` AND (a.id = '` + *id + `' ) `
			} else if index == (len(bookingIds) - 1) {
				query = query + ` OR  a.id = '` + *id + `' ) `
			} else {
				query = query + ` OR  a.id = '` + *id + `' `
			}
		}

		list, err := b.fetchDetailTransport(ctx, query)
		if err != nil {
			return nil, err
		}
		result = list
	}
	return result, nil
}

func (b bookingExpRepository) QueryHistoryPerMonthTransByUserId(ctx context.Context, bookingIds []*string) ([]*models.BookingExpJoin, error) {

	query := `SELECT
		a.*,
		b.trans_name,
		b.trans_title,
		b.trans_status,
		b.class AS trans_class,
		s.departure_date,
		s.departure_time,
		s.arrival_time,
		t.total_price,
		pm.name AS payment_type,
		pm.desc AS account_bank,
		pm.icon,
		t.status AS transaction_status,
		t.va_number,
		t.created_date AS created_date_transaction,
		m.merchant_name,
		m.phone_number AS merchant_phone,
		m.merchant_picture,
		hs.harbors_name AS harbor_source_name,
		h.harbors_name AS harbor_dest_name,
		t.ex_change_rates,
		t.ex_change_currency
	FROM
		booking_exps a
		LEFT JOIN transactions t ON t.booking_exp_id = a.id
			OR t.order_id = a.order_id
		LEFT JOIN payment_methods pm ON t.payment_method_id = pm.id
		JOIN transportations b ON a.trans_id = b.id
		LEFT JOIN schedules s ON b.id = s.trans_id
			AND a.schedule_id = s.id
		JOIN harbors h ON b.harbors_dest_id = h.id
		JOIN harbors hs ON b.harbors_source_id = hs.id
		JOIN merchants m ON b.merchant_id = m.id`

	result := make([]*models.BookingExpJoin, 0)
	if len(bookingIds) != 0 {
		for index, id := range bookingIds {
			if index == 0 && index != (len(bookingIds)-1) {
				query = query + ` AND (a.id = '` + *id + `' `
			} else if index == 0 && index == (len(bookingIds)-1) {
				query = query + ` AND (a.id = '` + *id + `' ) `
			} else if index == (len(bookingIds) - 1) {
				query = query + ` OR  a.id = '` + *id + `' ) `
			} else {
				query = query + ` OR  a.id = '` + *id + `' `
			}
		}

		list, err := b.fetchDetailTransport(ctx, query)
		if err != nil {
			return nil, err
		}
		result = list
	}
	return result, nil
}
func (b bookingExpRepository) QueryCountHistoryByUserId(ctx context.Context, userId string, yearMonth string) (int, error) {
	var count int
	if yearMonth != "" {
		date := yearMonth + "-" + "01" + " 00:00:00"
		query := `SELECT COUNT(*) as count 
					FROM 
						booking_exps a 
					LEFT JOIN 
						transactions t ON t.booking_exp_id = a.id OR t.order_id = a.order_id
					where a.user_id = ?
					and (t.created_date >= ? or t.modified_date >= ?)`
		rows, err := b.Conn.QueryContext(ctx, query, userId, date, date)
		if err != nil {
			logrus.Error(err)
			return 0, err
		}
		if err != nil {
			logrus.Error(err)
			return 0, err
		}

		resultCount, err := checkCount(rows)
		if err != nil {
			logrus.Error(err)
			return 0, err
		}
		count = resultCount
	} else {
		query := `SELECT COUNT(*) as count 
					FROM 
						booking_exps a 
					LEFT JOIN 
						transactions t ON t.booking_exp_id = a.id OR t.order_id = a.order_id	
					WHERE 
					a.user_id = ?
					and (t.created_date >= (NOW() - INTERVAL 1 MONTH) or t.modified_date >= (NOW() - INTERVAL 1 MONTH))`

		rows, err := b.Conn.QueryContext(ctx, query, userId)
		if err != nil {
			logrus.Error(err)
			return 0, err
		}
		if err != nil {
			logrus.Error(err)
			return 0, err
		}

		resultCount, err := checkCount(rows)
		if err != nil {
			logrus.Error(err)
			return 0, err
		}
		count = resultCount
	}

	return count, nil
}

func (b bookingExpRepository) QuerySelectIdHistoryByUserId(ctx context.Context, userId string, yearMonth string, limit, offset int) ([]*string, error) {
	var bookingIds []*string
	if yearMonth != "" {
		date := yearMonth + "-" + "01" + " 00:00:00"
		query := `SELECT DISTINCT a.id
					FROM 
						booking_exps a 
					LEFT JOIN 
						transactions t ON t.booking_exp_id = a.id OR t.order_id = a.order_id
					where a.user_id = ?
					and (t.created_date >= ? or t.modified_date >= ?)
					and t.status IN (1,2,5) 
					and DATE(a.booking_date) < CURRENT_DATE
				UNION 
					SELECT DISTINCT a.id
					FROM 
						booking_exps a 
					LEFT JOIN 
						transactions t ON t.booking_exp_id = a.id OR t.order_id = a.order_id
					where a.user_id = ?
					and (t.created_date >= ? or t.modified_date >= ?)
					and t.status IN (3,4)
				UNION 
					SELECT DISTINCT a.id
					FROM 
						booking_exps a 
					LEFT JOIN 
						transactions t ON t.booking_exp_id = a.id OR t.order_id = a.order_id
					where a.user_id = ?
					and (t.created_date >= ? or t.modified_date >= ?)
					and t.status = 0
					AND a.expired_date_payment < DATE_ADD(NOW(), INTERVAL 7 HOUR) `
		if limit != 0 {
			query = query + ` LIMIT ? OFFSET ?`
			rows, err := b.Conn.QueryContext(ctx, query, userId, date, date, userId, date, date, userId, date, date, limit, offset)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}

			for rows.Next() {
				var bookingId string
				err = rows.Scan(&bookingId)
				if err != nil {
					return nil, err
				}
				bookingIds = append(bookingIds, &bookingId)
			}
		} else {
			rows, err := b.Conn.QueryContext(ctx, query, userId, date, date, userId, date, date, userId, date, date)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}

			for rows.Next() {
				var bookingId string
				err = rows.Scan(&bookingId)
				if err != nil {
					return nil, err
				}
				bookingIds = append(bookingIds, &bookingId)
			}
		}

	} else {
		query := `SELECT DISTINCT a.id
					FROM 
						cgo_indonesia.booking_exps a 
					LEFT JOIN 
						cgo_indonesia.transactions t ON t.booking_exp_id = a.id OR t.order_id = a.order_id	
					WHERE 
					a.user_id = ?
					and (t.created_date >= (NOW() - INTERVAL 1 MONTH) or t.modified_date >= (NOW() - INTERVAL 1 MONTH))
					and t.status IN (1,2,5)
                    AND DATE(a.booking_date) < CURRENT_DATE 
				UNION
					SELECT DISTINCT a.id
					FROM 
						cgo_indonesia.booking_exps a 
					LEFT JOIN 
						cgo_indonesia.transactions t ON t.booking_exp_id = a.id OR t.order_id = a.order_id	
					WHERE 
					a.user_id = ?
					and (t.created_date >= (NOW() - INTERVAL 1 MONTH) or t.modified_date >= (NOW() - INTERVAL 1 MONTH))
					and t.status IN (3,4)
 				UNION 
 					SELECT DISTINCT a.id
					FROM 
						cgo_indonesia.booking_exps a 
					LEFT JOIN 
						cgo_indonesia.transactions t ON t.booking_exp_id = a.id OR t.order_id = a.order_id	
					WHERE 
					a.user_id = ?
					and (t.created_date >= (NOW() - INTERVAL 1 MONTH) or t.modified_date >= (NOW() - INTERVAL 1 MONTH))
					and t.status = 0 
					AND a.expired_date_payment < DATE_ADD(NOW(), INTERVAL 7 HOUR)`
		if limit != 0 {
			query = query + ` LIMIT ? OFFSET ?`
			rows, err := b.Conn.QueryContext(ctx, query, userId, userId, userId, limit, offset)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}

			for rows.Next() {
				var bookingId string
				err = rows.Scan(&bookingId)
				if err != nil {
					return nil, err
				}
				bookingIds = append(bookingIds, &bookingId)
			}
		} else {
			rows, err := b.Conn.QueryContext(ctx, query, userId, userId, userId)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}

			for rows.Next() {
				var bookingId string
				err = rows.Scan(&bookingId)
				if err != nil {
					return nil, err
				}
				bookingIds = append(bookingIds, &bookingId)
			}
		}
	}

	return bookingIds, nil
}
