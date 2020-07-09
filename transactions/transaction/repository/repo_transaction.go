package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"github.com/models"

	"github.com/sirupsen/logrus"
	"github.com/transactions/transaction"
)

type transactionRepository struct {
	Conn *sql.DB
}

func NewTransactionRepository(Conn *sql.DB) transaction.Repository {
	return &transactionRepository{Conn: Conn}
}

func (t transactionRepository) GetTransactionByExpIdORTransId(ctx context.Context, date string, expId string, transId string, merchantId string) ([]*models.TransactionOut, error) {
	var query string
	queryExp := `
	SELECT
		t.id as transaction_id,
		e.id as exp_id,
		e.exp_type,
		e.exp_title,
		booking_exp_id,
		b.order_id as booking_code,
		b.created_date as booking_date,
		b.booking_date as check_in_date,
		b.booked_by,
		guest_desc,
		b.booked_by_email as email,
		t.status as transaction_status,
		b.status as booking_status,
		t.total_price,
		t.experience_payment_id as experience_payment_id,
		merchant_name,
		b.order_id,
		e.exp_duration,
		p.province_name,
		co.country_name,
		t.promo_id,
		t.points,
		t.original_price,
		null as arrival_time,
		null as departure_time
	FROM
		transactions t
		JOIN experience_payments ep ON t.experience_payment_id = ep.id
		JOIN booking_exps b ON t.booking_exp_id = b.id
		JOIN experiences e ON b.exp_id = e.id
		JOIN merchants m ON e.merchant_id = m.id
		JOIN harbors  h ON e.harbors_id = h.id
		JOIN cities  c ON h.city_id = c.id
		JOIN provinces p on c.province_id = p.id
		JOIN countries co on p.country_id = co.id
		WHERE 
		t.is_deleted = 0 AND
		t.is_active = 1 `

	queryTrans := `
	SELECT
		t.id AS transaction_id,
		b.trans_id,
		trans_name,
		trans_title,
		booking_exp_id,
		b.order_id AS booking_code,
		b.created_date AS booking_date,
		b.booking_date AS check_in_date,
		b.booked_by,
		guest_desc,
		b.booked_by_email AS email,
		t.status AS transaction_status,
		b.status AS booking_status,
		t.total_price,
		tr.class as trans_class,
		merchant_name,
		b.order_id,
		tr.trans_capacity as exp_duration,
		hs.harbors_name AS province_name,
		h.harbors_name AS country_name,
		t.promo_id,
		t.points,
		t.original_price,
		s.arrival_time,
		s.departure_time
	FROM
		transactions t
		JOIN booking_exps b ON t.booking_exp_id = b.id OR t.order_id = b.order_id
		JOIN transportations tr ON b.trans_id = tr.id
		JOIN merchants m ON tr.merchant_id = m.id
		JOIN harbors h ON tr.harbors_dest_id = h.id
		JOIN harbors hs ON tr.harbors_source_id = hs.id
		JOIN schedules s ON b.schedule_id = s.id
		WHERE 
		t.is_deleted = 0 AND
		t.is_active = 1 `

	if merchantId != "" {
		queryExp = queryExp + ` AND e.merchant_id = '` + merchantId + `' `
		queryTrans = queryTrans + ` AND tr.merchant_id = '` + merchantId + `' `
	}
	if date != "" {
		queryExp = queryExp + ` AND DATE(b.booking_date) = '` + date + `' `
		queryTrans = queryTrans + ` AND DATE(b.booking_date) = '` + date + `' `
	}
	if expId != "" {
		query = queryExp + ` AND b.exp_id = '` + expId + `' `
	} else if transId != "" {
		query = queryTrans + ` AND b.trans_id = '` + transId + `' `
	}

	list, err := t.fetchWithJoin(ctx, query)
	if err != nil {
		return nil, err
	}

	return list, nil
}
func (t transactionRepository) GetTransactionByDate(ctx context.Context, date string, isExperience bool, isTransportation bool, merchantId string) ([]*models.TransactionByDate, error) {
	queryExp := `SELECT DISTINCT
				b.exp_id,
				e.exp_title,
				null as trans_id,
				null as departure_time,
				null as arrival_time,
				null as harbors_dest,
				null as harbors_source,
				e.exp_max_guest as capacity

				FROM transactions t 
				JOIN booking_exps b ON t.booking_exp_id = b.id
				JOIN experiences e ON b.exp_id = e.id
				WHERE t.is_deleted = 0 
					AND t.is_active = 1
					AND t.status IN (0,1,2,5)
					AND DATE(b.booking_date) = '` + date + `' `

	queryTrans := `SELECT DISTINCT
				null as exp_id,
				null as exp_title,
				b.trans_id,
				s.departure_time,
				s.arrival_time,
				h.harbors_name as harbors_dest,
				hs.harbors_name as harbors_source,
				trans.trans_capacity as capacity
				
				FROM transactions t 
				JOIN booking_exps b ON t.booking_exp_id = b.id OR t.order_id = b.order_id
				JOIN transportations trans ON b.trans_id = trans.id
				JOIN schedules s ON b.schedule_id = s.id
				JOIN harbors h ON trans.harbors_dest_id = h.id
				JOIN harbors hs ON trans.harbors_source_id = hs.id

				WHERE t.is_deleted = 0 
					AND t.is_active = 1
					AND t.status IN (0,1,2,5)
					AND DATE(b.booking_date) = '` + date + `' `
	var query string
	if merchantId != "" {
		queryExp = queryExp + ` AND e.merchant_id = '` + merchantId + `' `
		queryTrans = queryTrans + ` AND trans.merchant_id = '` + merchantId + `' `
	}
	if isTransportation == true && isExperience == false {
		query = queryTrans
	} else if isExperience == true && isTransportation == false {
		query = queryExp
	} else {
		query = queryExp + ` UNION ` + queryTrans
	}
	rows, err := t.Conn.QueryContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	result := make([]*models.TransactionByDate, 0)
	for rows.Next() {
		t := new(models.TransactionByDate)
		err = rows.Scan(
			&t.ExpId,
			&t.ExpTitle,
			&t.TransId,
			&t.DepartureTime,
			&t.ArrivalTime,
			&t.HarborsDest,
			&t.HarborsSource,
			&t.Capacity,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (t transactionRepository) GetCountTransactionByPromoId(ctx context.Context, promoId string, userId string) (int, error) {
	query := `select count(*) from transactions a
											join booking_exps b on a.booking_exp_id = b.id 
											where a.status in (0,1,2,5)
												and a.promo_id = ? `

	if userId != "" {
		query = query + ` and b.user_id = '` + userId + `'`
	} else {
		query = query + ` and b.user_id != '' `
	}

	rows, err := t.Conn.QueryContext(ctx, query, promoId)
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

func (t transactionRepository) GetIdTransactionExpired(ctx context.Context) ([]*string, error) {
	query := `SELECT t.id
					FROM transactions t 
					JOIN booking_exps b ON t.booking_exp_id = b.id OR t.order_id = b.order_id
					WHERE t.status = 0 
					AND b.expired_date_payment <= ?
					ORDER BY b.expired_date_payment DESC`

	rows, err := t.Conn.QueryContext(ctx, query, time.Now())
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	var transactionIds []*string
	for rows.Next() {
		var transactionId string
		err = rows.Scan(&transactionId)
		if err != nil {
			return nil, err
		}

		transactionIds = append(transactionIds, &transactionId)
	}

	return transactionIds, nil
}

func (t transactionRepository) GetTransactionDownPaymentByDate(ctx context.Context) ([]*models.TransactionWithBooking, error) {
	threeDay := time.Now().AddDate(0, 0, 3).Format("2006-01-02")
	twoDay := time.Now().AddDate(0, 0, 2).Format("2006-01-02")
	tenDay := time.Now().AddDate(0, 0, 10).Format("2006-01-02")
	threetyDay := time.Now().AddDate(0, 0, 30).Format("2006-01-02")

	query := `
	SELECT 
		e.exp_title,
		be.booked_by,
		be.booked_by_email,
		be.booking_date,
		t.total_price,
		ep.price ,
		e.exp_duration,
		be.order_id,
		m.merchant_name,
		m.phone_number as merchant_phone,
		e.exp_payment_deadline_amount,
		e.exp_payment_deadline_type
	FROM transactions t
	JOIN experience_payments ep on ep.id = t.experience_payment_id
	JOIN booking_exps be on be.id = t.booking_exp_id
	JOIN experiences e on e.id = be.exp_id
	JOIN merchants m on m.id = e.merchant_id	
	WHERE 
		ep.exp_payment_type_id = '86e71b8d-acc3-4ade-80c0-de67b9100633' AND 
		t.total_price != ep.price AND 
		(DATE(be.booking_date) = ? OR DATE(be.booking_date) = ? OR DATE(be.booking_date) = ? OR DATE(be.booking_date) = ?)`

	rows, err := t.Conn.QueryContext(ctx, query, threeDay, twoDay, tenDay, threetyDay)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	result := make([]*models.TransactionWithBooking, 0)
	for rows.Next() {
		t := new(models.TransactionWithBooking)
		err = rows.Scan(
			&t.ExpTitle,
			&t.BookedBy,
			&t.BookedByEmail,
			&t.BookingDate,
			&t.TotalPrice,
			&t.Price,
			&t.ExpDuration,
			&t.OrderId,
			&t.MerchantName,
			&t.MerchantPhone,
			&t.ExpPaymentDeadlineAmount,
			&t.ExpPaymentDeadlineType,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (t transactionRepository) GetCountByExpId(ctx context.Context, date string, expId string,isTransaction bool) ([]*string, error) {
	query := `
	select b.guest_desc from transactions a
										join booking_exps b on a.order_id = b.order_id 
										join experiences c on c.id = b.exp_id 
										where a.status < 3 and (b.status = 1 or b.status = 3)
										and date(b.booking_date) = ? and exp_id = ?`
	if isTransaction == true {
		query = `select b.guest_desc from transactions a
										join booking_exps b on a.order_id = b.order_id 
										join experiences c on c.id = b.exp_id 
										where a.status in (0,1,2,5)
										and date(b.booking_date) = ? 
										and exp_id = ?
										and a.is_deleted = 0 
										and a.is_active = 1`
	}
	rows, err := t.Conn.QueryContext(ctx, query, date, expId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	var bookingGuestResult []*string
	for rows.Next() {
		var bookingDesc string
		err = rows.Scan(&bookingDesc)
		if err != nil {
			return nil, err
		}
		bookingGuestResult = append(bookingGuestResult,&bookingDesc)
	}

	return bookingGuestResult, nil
}
func (t transactionRepository) GetCountByTransId(ctx context.Context, transId string,isTransaction bool,date string) ([]*string, error) {
	query := `select b.guest_desc from transactions a
											join booking_exps b on a.order_id = b.order_id 
											join transportations c on c.id = b.trans_id 
											join schedules d on b.schedule_id = d.id
											where a.status < 3 and b.trans_id = ?`
	if isTransaction == true {
		query = `select b.guest_desc from transactions a
											join booking_exps b on a.order_id = b.order_id 
											join transportations c on c.id = b.trans_id 
											join schedules d on b.schedule_id = d.id
											where a.status in (0,1,2,5) 
											and b.trans_id = ?
											and a.is_deleted = 0 
											and a.is_active = 1`
	}
	if date != ""{
		query = query + ` and date(b.booking_date) = '` + date + `' `
	}
	rows, err := t.Conn.QueryContext(ctx, query, transId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	var bookingGuestResult []*string
	for rows.Next() {
		var bookingDesc string
		err = rows.Scan(&bookingDesc)
		if err != nil {
			return nil, err
		}
		bookingGuestResult = append(bookingGuestResult,&bookingDesc)
	}

	return bookingGuestResult, nil
}
func (t transactionRepository) UpdateAfterPayment(ctx context.Context, status int, vaNumber string, transactionId, bookingId string) error {
	query := `UPDATE transactions SET status = ?, va_number = ? WHERE (id = ? OR booking_exp_id = ? OR order_id = ?)`

	stmt, err := t.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, status, vaNumber, transactionId, bookingId, bookingId)
	if err != nil {
		return err
	}

	return nil
}

func (t transactionRepository) CountThisMonth(ctx context.Context) (*models.TotalTransaction, error) {
	query := `
	SELECT
		count(CAST(created_date AS DATE)) as transaction_count,
		SUM(total_price) as transaction_value_total
	FROM
		transactions
	WHERE
		is_deleted = 0
		AND is_active = 1
		AND status = 2
		AND created_date BETWEEN date_add(CURRENT_DATE, interval - DAY(CURRENT_DATE) + 1 DAY)
		AND CURRENT_DATE`

	rows, err := t.Conn.QueryContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	total := new(models.TotalTransaction)
	for rows.Next() {
		err = rows.Scan(&total.TransactionCount, &total.TransactionValueTotal)
		if err != nil {
			return nil, err
		}
	}

	return total, nil
}

func (t transactionRepository) List(ctx context.Context, startDate, endDate, search, status string, limit, offset *int, merchantId string, isTransportation bool, isExperience bool, isSchedule bool, tripType, paymentType, activityType string, confirmType string,class string,departureTimeStart string,departureTimeEnd string,arrivalTimeStart string,arrivalTimeEnd string) ([]*models.TransactionOut, error) {
	var transactionStatus int
	//var bookingStatus int

	query := `
	SELECT
		t.id as transaction_id,
		e.id as exp_id,
		e.exp_type,
		e.exp_title,
		booking_exp_id,
		b.order_id as booking_code,
		b.created_date as booking_date,
		b.booking_date as check_in_date,
		b.booked_by,
		guest_desc,
		b.booked_by_email as email,
		t.status as transaction_status,
		b.status as booking_status,
		t.total_price,
		t.experience_payment_id as experience_payment_id,
		merchant_name,
		b.order_id,
		e.exp_duration,
		p.province_name,
		co.country_name,
		t.promo_id,
		t.points,
		t.original_price,
		null as arrival_time,
		null as departure_time
	FROM
		transactions t
		JOIN experience_payments ep ON t.experience_payment_id = ep.id
		JOIN booking_exps b ON t.booking_exp_id = b.id
		JOIN experiences e ON b.exp_id = e.id
		JOIN merchants m ON e.merchant_id = m.id
		JOIN harbors  h ON e.harbors_id = h.id
		JOIN cities  c ON h.city_id = c.id
		JOIN provinces p on c.province_id = p.id
		JOIN countries co on p.country_id = co.id
		WHERE 
		t.is_deleted = 0 AND
		t.is_active = 1 `

	queryT := `
	SELECT
		t.id AS transaction_id,
		b.trans_id,
		trans_name,
		trans_title,
		booking_exp_id,
		b.order_id AS booking_code,
		b.created_date AS booking_date,
		b.booking_date AS check_in_date,
		b.booked_by,
		guest_desc,
		b.booked_by_email AS email,
		t.status AS transaction_status,
		b.status AS booking_status,
		t.total_price,
		tr.class as trans_class,
		merchant_name,
		b.order_id,
		tr.trans_capacity as exp_duration,
		hs.harbors_name AS province_name,
		h.harbors_name AS country_name,
		t.promo_id,
		t.points,
		t.original_price,
		s.arrival_time,
		s.departure_time
	FROM
		transactions t
		JOIN booking_exps b ON t.booking_exp_id = b.id OR t.order_id = b.order_id
		JOIN transportations tr ON b.trans_id = tr.id
		JOIN merchants m ON tr.merchant_id = m.id
		JOIN harbors h ON tr.harbors_dest_id = h.id
		JOIN harbors hs ON tr.harbors_source_id = hs.id
		JOIN schedules s ON b.schedule_id = s.id
		WHERE 
		t.is_deleted = 0 AND
		t.is_active = 1 `

	if merchantId != "" {
		query = query + ` AND e.merchant_id = '` + merchantId + `' `
		queryT = queryT + ` AND tr.merchant_id = '` + merchantId + `' `
	}



	if tripType == "0" {
		query = query + ` AND e.exp_trip_type = 'Private Trip' `
	} else if tripType == "1" {
		query = query + ` AND e.exp_trip_type = 'Share Trip' `
	} else if tripType == "2" {
		query = query + ` AND (e.exp_trip_type = 'Share Trip' OR e.exp_trip_type = 'Private Trip') `
	}

	if paymentType == "0" {
		query = query + ` AND ep.exp_payment_type_id = '8a5e3eef-a6db-4584-a280-af5ab18a979b' `
	} else if paymentType == "1" {
		query = query + ` AND ep.exp_payment_type_id = '86e71b8d-acc3-4ade-80c0-de67b9100633' `
	} else if paymentType == "2" {
		query = query + ` AND (ep.exp_payment_type_id = '8a5e3eef-a6db-4584-a280-af5ab18a979b' 
								OR ep.exp_payment_type_id = '86e71b8d-acc3-4ade-80c0-de67b9100633') `
	}

	if activityType != "[]" && activityType != "" {
		var activityTypes []int
		if activityType != "" {
			if errUnmarshal := json.Unmarshal([]byte(activityType), &activityTypes); errUnmarshal != nil {
				return nil, errUnmarshal
			}
		}

		if len(activityTypes) != 0 {
			for index, id := range activityTypes {
				if index == 0 && index != (len(activityTypes)-1) {
					query = query + ` AND (e.id = (SELECT distinct exp_id FROM filter_activity_types where exp_id = e.id and exp_type_id = ` + strconv.Itoa(id) + ` )`
				} else if index == 0 && index == (len(activityTypes)-1) {
					query = query + ` AND (e.id = (SELECT distinct exp_id FROM filter_activity_types where exp_id = e.id and exp_type_id = ` + strconv.Itoa(id) + ` )` + ` ) `
				} else if index == (len(activityTypes) - 1) {
					query = query + ` OR e.id = (SELECT distinct exp_id FROM filter_activity_types where exp_id = e.id and exp_type_id = ` + strconv.Itoa(id) + ` )` + ` )`
				} else {
					query = query + ` OR e.id = (SELECT distinct exp_id FROM filter_activity_types where exp_id = e.id and exp_type_id = ` + strconv.Itoa(id) + ` )`
				}
			}

		}
	}

	if confirmType == "0" {
		query = query + ` AND t.status not in (1) `
	} else if confirmType == "1" {
		query = query + ` AND t.status not in (5)' `
	} else if confirmType == "2" {

	}

	if search != "" {
		keyword := `'%` + search + `%'`
		query = query + ` AND (LOWER(b.booked_by) LIKE LOWER(` + keyword + `) OR LOWER(b.order_id) LIKE LOWER(` + keyword + `))`
		queryT = queryT + ` AND (LOWER(b.booked_by) LIKE LOWER(` + keyword + `) OR LOWER(b.order_id) LIKE LOWER(` + keyword + `))`
	}
	if startDate != "" && endDate != "" {
		if isSchedule == true {
			query = query + ` AND DATE(b.booking_date) BETWEEN '` + startDate + `' AND '` + endDate + `'`
			queryT = queryT + ` AND DATE(b.booking_date) BETWEEN '` + startDate + `' AND '` + endDate + `'`
		} else {
			query = query + ` AND DATE(b.created_date) BETWEEN '` + startDate + `' AND '` + endDate + `'`
			queryT = queryT + ` AND DATE(b.created_date) BETWEEN '` + startDate + `' AND '` + endDate + `'`
		}

	}
	if isTransportation == true && isExperience == false {
		query = query + ` AND b.trans_id != '' `
		queryT = queryT + ` AND b.trans_id != '' `

		if class == "0"{
			queryT = queryT + ` AND tr.class = 'Economy' `
		}else if class == "1"{
			queryT = queryT + ` AND tr.class = 'Executive' `
		}

		if departureTimeStart != "" && departureTimeEnd != ""{
			queryT = queryT + ` AND s.departure_time between '` +  departureTimeStart + `' AND '` + departureTimeEnd + `' `
		}
		if arrivalTimeStart != "" && arrivalTimeEnd != ""{
			queryT = queryT + ` AND s.arrival_time between '` +  arrivalTimeStart + `' AND '` + arrivalTimeEnd + `' `
		}

	} else if isExperience == true && isTransportation == false {
		query = query + ` AND b.exp_id != '' `
		queryT = queryT + ` AND b.exp_id != '' `
	}

	unionQuery := query + ` UNION ` + queryT
	if tripType != "" && activityType == "" {
		unionQuery = query
	} else if activityType != "" && tripType == "" {
		unionQuery = query
	} else if tripType != "" && activityType != "" {
		unionQuery = query
	}

	if limit != nil && offset != nil {
		unionQuery = unionQuery +
			` ORDER BY booking_date DESC LIMIT ` + strconv.Itoa(*limit) +
			` OFFSET ` + strconv.Itoa(*offset) + ` `

	}

	list, err := t.fetchWithJoin(ctx, unionQuery)
	if status != "" {
		if status == "pending" {
			transactionStatus = 0
		} else if status == "waitingApproval" {
			transactionStatus = 1
		} else if status == "confirm" {
			transactionStatus = 2
			query = query + ` AND b.booking_date > CURRENT_DATE `
			queryT = queryT + ` AND b.booking_date > CURRENT_DATE `
		} else if status == "upcoming" {
			transactionStatus = 2
			query = query + ` AND date_add(b.booking_date, interval + 14 DAY) >= DATE(CURRENT_DATE) `
			queryT = queryT + ` AND date_add(b.booking_date, interval + 14 DAY) >= DATE(CURRENT_DATE) `
		} else if status == "finished" {
			transactionStatus = 7
			//query = query + ` AND b.booking_date < CURRENT_DATE `
			//queryT = queryT + ` AND b.booking_date < CURRENT_DATE `
		}
		var querySt string
		var queryTSt string
		if status == "waitingApproval" {
			querySt = query + ` AND (t.status = ` + strconv.Itoa(transactionStatus) + ` OR t.status = 5 ) `
			queryTSt = queryT + ` AND (t.status = ` + strconv.Itoa(transactionStatus) + ` OR t.status = 5 ) `
		} else {
			querySt = query + ` AND t.status = ` + strconv.Itoa(transactionStatus)
			queryTSt = queryT + ` AND t.status = ` + strconv.Itoa(transactionStatus)
		}
		unionQuery = querySt + ` UNION ` + queryTSt
		if tripType != "" && activityType == "" {
			unionQuery = querySt
		} else if activityType != "" && tripType == "" {
			unionQuery = querySt
		} else if tripType != "" && activityType != "" {
			unionQuery = querySt
		}

		if limit != nil && offset != nil {
			unionQuery = unionQuery +
				` ORDER BY booking_date DESC LIMIT ` + strconv.Itoa(*limit) +
				` OFFSET ` + strconv.Itoa(*offset) + ` `

		}
		list, err = t.fetchWithJoin(ctx, unionQuery)

		if status == "failed" {
			transactionStatus = 3
			cancelledStatus := 4
			querySt = query + ` AND t.status IN (` + strconv.Itoa(transactionStatus) + `,` + strconv.Itoa(cancelledStatus) + `)`
			queryTSt = queryT + ` AND t.status IN (` + strconv.Itoa(transactionStatus) + `,` + strconv.Itoa(cancelledStatus) + `)`
			unionQuery = querySt + ` UNION ` + queryTSt
			if tripType != "" && activityType == "" {
				unionQuery = querySt
			} else if activityType != "" && tripType == "" {
				unionQuery = querySt
			} else if tripType != "" && activityType != "" {
				unionQuery = querySt
			}

			if limit != nil && offset != nil {
				unionQuery = unionQuery +
					` ORDER BY booking_date DESC LIMIT ` + strconv.Itoa(*limit) +
					` OFFSET ` + strconv.Itoa(*offset) + ` `

			}
			list, err = t.fetchWithJoin(ctx, unionQuery)
		}

		if status == "boarded" {
			//transactionStatus = 2
			//bookingStatus = 3
			//querySt = query + ` AND t.status = ` + strconv.Itoa(transactionStatus) + ` AND b.status = ` + strconv.Itoa(bookingStatus)
			//queryTSt = queryT + ` AND t.status = ` + strconv.Itoa(transactionStatus) + ` AND b.status = ` + strconv.Itoa(bookingStatus)
			transactionStatus = 6
			//bookingStatus = 3
			querySt = query + ` AND t.status = ` + strconv.Itoa(transactionStatus)
			queryTSt = queryT + ` AND t.status = ` + strconv.Itoa(transactionStatus)
			unionQuery = querySt + ` UNION ` + queryTSt
			if tripType != "" && activityType == "" {
				unionQuery = querySt
			} else if activityType != "" && tripType == "" {
				unionQuery = querySt
			} else if tripType != "" && activityType != "" {
				unionQuery = querySt
			}

			if limit != nil && offset != nil {
				unionQuery = unionQuery +
					` ORDER BY booking_date DESC LIMIT ` + strconv.Itoa(*limit) +
					` OFFSET ` + strconv.Itoa(*offset) + ` `

			}
			list, err = t.fetchWithJoin(ctx, unionQuery)
		}
	}
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (t transactionRepository) fetchWithJoin(ctx context.Context, query string, args ...interface{}) ([]*models.TransactionOut, error) {
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

	result := make([]*models.TransactionOut, 0)
	for rows.Next() {
		t := new(models.TransactionOut)
		err = rows.Scan(
			&t.TransactionId,
			&t.ExpId,
			&t.ExpType,
			&t.ExpTitle,
			&t.BookingExpId,
			&t.BookingCode,
			&t.BookingDate,
			&t.CheckInDate,
			&t.BookedBy,
			&t.GuestDesc,
			&t.Email,
			&t.TransactionStatus,
			&t.BookingStatus,
			&t.TotalPrice,
			&t.ExperiencePaymentId,
			&t.MerchantName,
			&t.OrderId,
			&t.ExpDuration,
			&t.ProvinceName,
			&t.CountryName,
			&t.PromoId,
			&t.Points,
			&t.OriginalPrice,
			&t.ArrivalTime,
			&t.DepartureTime,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (t transactionRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.TransactionWMerchant, error) {
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

	result := make([]*models.TransactionWMerchant, 0)
	for rows.Next() {
		t := new(models.TransactionWMerchant)
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
			&t.BookingType,
			&t.BookingExpId,
			&t.PromoId,
			&t.PaymentMethodId,
			&t.ExperiencePaymentId,
			&t.Status,
			&t.TotalPrice,
			&t.Currency,
			&t.OrderId,
			&t.VaNumber,
			&t.ExChangeRates,
			&t.ExChangeCurrency,
			&t.Points,
			&t.OriginalPrice,
			&t.Remarks,
			&t.MerchantId,
			&t.OrderIdBook,
			&t.BookedBy,
			&t.ExpTitle,
			&t.BookingDate,
			&t.ExpId,
			&t.TransId,
			&t.BookingStatus,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (t transactionRepository) CountSuccess(ctx context.Context) (int, error) {
	query := `SELECT count(*) as count FROM transactions WHERE is_deleted = 0 AND is_active = 1 AND status = 2`

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

func (t transactionRepository) Count(ctx context.Context, startDate, endDate, search, status string, merchantId string, isTransportation bool, isExperience bool,isSchedule bool,tripType,paymentType,activityType string,confirmType string,class string,departureTimeStart string,departureTimeEnd string,arrivalTimeStart string,arrivalTimeEnd string) (int, error) {
	query := `
	SELECT
		count(*) as count
	FROM 
		transactions t
		JOIN experience_payments ep ON t.experience_payment_id = ep.id
		JOIN booking_exps b ON t.booking_exp_id = b.id
		JOIN experiences e ON b.exp_id = e.id
		JOIN merchants m ON e.merchant_id = m.id
		JOIN harbors  h ON e.harbors_id = h.id
		JOIN cities  c ON h.city_id = c.id
		JOIN provinces p on c.province_id = p.id
		JOIN countries co on p.country_id = co.id
	WHERE
		t.is_deleted = 0
		AND t.is_active = 1`

	queryT := `
	SELECT
		count(*) as count
	FROM
		transactions t
		JOIN booking_exps b ON t.booking_exp_id = b.id OR t.order_id = b.order_id
		JOIN transportations tr ON b.trans_id = tr.id
		JOIN merchants m ON tr.merchant_id = m.id
		JOIN harbors h ON tr.harbors_dest_id = h.id
		JOIN harbors hs ON tr.harbors_source_id = hs.id
		JOIN schedules s ON b.schedule_id = s.id
	WHERE
		t.is_deleted = 0
		AND t.is_active = 1`

	if merchantId != "" {
		query = query + ` AND e.merchant_id = '` + merchantId + `' `
		queryT = queryT + ` AND tr.merchant_id = '` + merchantId + `' `
	}
	if tripType == "0" {
		query = query + ` AND e.exp_trip_type = 'Private Trip' `
	} else if tripType == "1" {
		query = query + ` AND e.exp_trip_type = 'Share Trip' `
	} else if tripType == "2" {
		query = query + ` AND (e.exp_trip_type = 'Share Trip' OR e.exp_trip_type = 'Private Trip') `
	}

	if paymentType == "0" {
		query = query + ` AND ep.exp_payment_type_id = '8a5e3eef-a6db-4584-a280-af5ab18a979b' `
	} else if paymentType == "1" {
		query = query + ` AND ep.exp_payment_type_id = '86e71b8d-acc3-4ade-80c0-de67b9100633' `
	} else if paymentType == "2" {
		query = query + ` AND (ep.exp_payment_type_id = '8a5e3eef-a6db-4584-a280-af5ab18a979b' 
								OR ep.exp_payment_type_id = '86e71b8d-acc3-4ade-80c0-de67b9100633') `
	}

	if activityType != "[]" && activityType != "" {
		var activityTypes []int
		if activityType != "" {
			if errUnmarshal := json.Unmarshal([]byte(activityType), &activityTypes); errUnmarshal != nil {
				return 0, errUnmarshal
			}
		}

		if len(activityTypes) != 0 {
			for index, id := range activityTypes {
				if index == 0 && index != (len(activityTypes)-1) {
					query = query + ` AND (e.id = (SELECT distinct exp_id FROM filter_activity_types where exp_id = e.id and exp_type_id = ` + strconv.Itoa(id) + ` )`
				} else if index == 0 && index == (len(activityTypes)-1) {
					query = query + ` AND (e.id = (SELECT distinct exp_id FROM filter_activity_types where exp_id = e.id and exp_type_id = ` + strconv.Itoa(id) + ` )` + ` ) `
				} else if index == (len(activityTypes) - 1) {
					query = query + ` OR e.id = (SELECT distinct exp_id FROM filter_activity_types where exp_id = e.id and exp_type_id = ` + strconv.Itoa(id) + ` )` + ` )`
				} else {
					query = query + ` OR e.id = (SELECT distinct exp_id FROM filter_activity_types where exp_id = e.id and exp_type_id = ` + strconv.Itoa(id) + ` )`
				}
			}

		}
	}

	if confirmType == "0" {
		query = query + ` AND t.status not in (1) `
	} else if confirmType == "1" {
		query = query + ` AND t.status not in (5)' `
	} else if confirmType == "2" {

	}

	if isTransportation == true && isExperience == false {
		query = query + ` AND b.trans_id != '' `
		queryT = queryT + ` AND b.trans_id != '' `

		if class == "0"{
			queryT = queryT + ` AND tr.class = 'Economy' `
		}else if class == "1"{
			queryT = queryT + ` AND tr.class = 'Executive' `
		}

		if departureTimeStart != "" && departureTimeEnd != ""{
			queryT = queryT + ` AND s.departure_time between '` +  departureTimeStart + `' AND '` + departureTimeEnd + `' `
		}
		if arrivalTimeStart != "" && arrivalTimeEnd != ""{
			queryT = queryT + ` AND s.arrival_time between '` +  arrivalTimeStart + `' AND '` + arrivalTimeEnd + `' `
		}

	} else if isExperience == true && isTransportation == false {
		query = query + ` AND b.exp_id != '' `
		queryT = queryT + ` AND b.exp_id != '' `
	}

	if search != "" {
		keyword := `'%` + search + `%'`
		query = query + ` AND (LOWER(b.booked_by) LIKE LOWER(` + keyword + `) OR LOWER(b.order_id) LIKE LOWER(` + keyword + `))`
		queryT = queryT + ` AND (LOWER(b.booked_by) LIKE LOWER(` + keyword + `) OR LOWER(b.order_id) LIKE LOWER(` + keyword + `))`
	}
	if startDate != "" && endDate != "" {
		query = query + ` AND DATE(b.created_date) BETWEEN '` + startDate + `' AND '` + endDate + `'`
		queryT = queryT + ` AND DATE(b.created_date) BETWEEN '` + startDate + `' AND '` + endDate + `'`
	}
	if isTransportation == true && isExperience == false {
		query = query + ` AND b.trans_id != '' `
		queryT = queryT + ` AND b.trans_id != '' `
	} else if isExperience == true && isTransportation == false {
		query = query + ` AND b.exp_id != '' `
		queryT = queryT + ` AND b.exp_id != '' `
	}
	unionQuery := query + ` UNION` + queryT
	rows, err := t.Conn.QueryContext(ctx, unionQuery)
	var transactionStatus int
	if status != "" {
		if status == "pending" {
			transactionStatus = 0
		} else if status == "waitingApproval" {
			transactionStatus = 1
		} else if status == "confirm" {
			transactionStatus = 2
			query = query + ` AND b.booking_date > CURRENT_DATE `
			queryT = queryT + ` AND b.booking_date > CURRENT_DATE `
		} else if status == "upcoming" {
			transactionStatus = 2
			query = query + ` AND date_add(b.booking_date, interval + 14 DAY) >= DATE(CURRENT_DATE) `
			queryT = queryT + ` AND date_add(b.booking_date, interval + 14 DAY) >= DATE(CURRENT_DATE) `
		} else if status == "finished" {
			transactionStatus = 7
			//query = query + ` AND b.booking_date < CURRENT_DATE `
			//queryT = queryT + ` AND b.booking_date < CURRENT_DATE `
		}

		querySt := query + ` AND t.status = ` + strconv.Itoa(transactionStatus)
		queryTSt := queryT + ` AND t.status = ` + strconv.Itoa(transactionStatus)
		unionQuery = querySt + ` UNION ` + queryTSt
		rows, err = t.Conn.QueryContext(ctx, unionQuery)

		if status == "failed" {
			transactionStatus = 3
			cancelledStatus := 4
			querySt = query + ` AND t.status IN (` + strconv.Itoa(transactionStatus) + `,` + strconv.Itoa(cancelledStatus) + `)`
			queryTSt = queryT + ` AND t.status IN (` + strconv.Itoa(transactionStatus) + `,` + strconv.Itoa(cancelledStatus) + `)`
			unionQuery = querySt + ` UNION ` + queryTSt
			rows, err = t.Conn.QueryContext(ctx, unionQuery)
		}

		if status == "boarded" {
			//transactionStatus = 1
			//bookingStatus := 3
			//querySt = query + ` AND t.status = ` + strconv.Itoa(transactionStatus) + ` AND b.status = ` + strconv.Itoa(bookingStatus)
			//queryTSt = queryT + ` AND t.status = ` + strconv.Itoa(transactionStatus) + ` AND b.status = ` + strconv.Itoa(bookingStatus)
			transactionStatus = 6
			//bookingStatus = 3
			querySt = query + ` AND t.status = ` + strconv.Itoa(transactionStatus)
			queryTSt = queryT + ` AND t.status = ` + strconv.Itoa(transactionStatus)
			unionQuery = querySt + ` UNION ` + queryTSt
			rows, err = t.Conn.QueryContext(ctx, unionQuery)
		}
		if err != nil {
			logrus.Error(err)
			return 0, err
		}
	}

	count, err := checkCountUnion(rows)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	return count, nil
}
func (m transactionRepository) GetByBookingDate(ctx context.Context, bookingDate string, transId string, expId string) ([]*models.TransactionWMerchant, error) {
	query := `SELECT t.*,e.merchant_id,b.order_id as order_id_book,b.booked_by,e.exp_title,b.booking_date ,b.exp_id,b.trans_id,b.status as booking_status
				FROM transactions t
				join booking_exps b on t.booking_exp_id = b.id
				join experiences e on b.exp_id = e.id 
				WHERE DATE(b.booking_date) = ?`
	if expId != "" {
		query = query + ` AND b.exp_id = '` + expId + `' `
	}
	var list []*models.TransactionWMerchant
	if transId != "" {
		query := `SELECT t.*,e.merchant_id,b.order_id as order_id_book,b.booked_by,e.trans_title,b.booking_date,b.exp_id,b.trans_id,b.status as booking_status
				FROM transactions t
				join booking_exps b on t.booking_exp_id = b.id OR t.order_id = b.order_id
				join transportations e on b.trans_id = e.id 
				WHERE DATE(b.booking_date) = ?`
		if transId != "" {
			query = query + ` AND b.trans_id = '` + transId + `' `
		}
		list, _ = m.fetch(ctx, query, bookingDate)
	} else {
		list, _ = m.fetch(ctx, query, bookingDate)
	}
	//if err != nil {
	//	return nil, err
	//}

	if len(list) > 0 {
		return list, nil
	} else {
		return make([]*models.TransactionWMerchant, 0), models.ErrNotFound
	}
	return nil, nil
}

func (m transactionRepository) GetById(ctx context.Context, id string) (*models.TransactionWMerchant, error) {
	query := `SELECT t.*,e.merchant_id,b.order_id as order_id_book,b.booked_by,e.exp_title,b.booking_date,b.exp_id,b.trans_id ,b.status as booking_status
				FROM transactions t
				join booking_exps b on t.booking_exp_id = b.id
				join experiences e on b.exp_id = e.id WHERE t.id = ?`

	list, err := m.fetch(ctx, query, id)
	if len(list) == 0 {
		query := `SELECT t.*,e.merchant_id,b.order_id as order_id_book,b.booked_by,e.trans_title,b.booking_date,b.exp_id,b.trans_id ,b.status as booking_status
				FROM transactions t
				join booking_exps b on t.booking_exp_id = b.id OR t.order_id = b.order_id
				join transportations e on b.trans_id = e.id WHERE t.id = ?`

		list, _ = m.fetch(ctx, query, id)
	}
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		res := list[0]
		return res, nil
	} else {
		return nil, models.ErrNotFound
	}
	return nil, nil
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

func checkCountUnion(rows *sql.Rows) (result int, err error) {
	var count int
	results := make([]int, 2)
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
		results = append(results, count)
	}

	if len(results) > 0 {
		for _, r := range results {
			result += r
		}
	}
	return result, nil
}
