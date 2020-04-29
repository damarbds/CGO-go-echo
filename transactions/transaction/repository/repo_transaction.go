package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

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

func (t transactionRepository) UpdateStatus(ctx context.Context, status int, transactionId, bookingId string) error {
	query := `UPDATE transactions SET status = ? WHERE (id = ? OR booking_exp_id = ? OR order_id = ?)`

	stmt, err := t.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, status, transactionId, bookingId, bookingId)
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

func (t transactionRepository) List(ctx context.Context, startDate, endDate, search, status string, limit, offset int) ([]*models.TransactionOut, error) {
	var transactionStatus int
	var bookingStatus int

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
		ep.id as experience_payment_id
	FROM
		transactions t
		JOIN booking_exps b ON t.booking_exp_id = b.id
		JOIN experiences e ON b.exp_id = e.id
		JOIN experience_payments ep ON e.id = ep.exp_id
	WHERE 
		t.is_deleted = 0
		AND t.is_active = 1
	`

	queryT := `
	SELECT
		t.id AS transaction_id,
		trans_id,
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
		tr.class as trans_class
	FROM
		transactions t
		JOIN booking_exps b ON t.booking_exp_id = b.id
		JOIN transportations tr ON b.trans_id = tr.id
	WHERE
		t.is_deleted = 0
		AND t.is_active = 1`

	if search != "" {
		keyword := `'%` + search + `%'`
		query = query + ` AND (LOWER(b.booked_by) LIKE LOWER(` + keyword + `) OR LOWER(b.order_id) LIKE LOWER(` + keyword + `))`
		queryT = queryT + ` AND (LOWER(b.booked_by) LIKE LOWER(` + keyword + `) OR LOWER(b.order_id) LIKE LOWER(` + keyword + `))`
	}
	if startDate != "" && endDate != "" {
		query = query + ` AND DATE(b.created_date) BETWEEN '` + startDate + `' AND '` + endDate + `'`
		queryT = queryT + ` AND DATE(b.created_date) BETWEEN '` + startDate + `' AND '` + endDate + `'`
	}
	unionQuery := query + ` UNION ` + queryT
	queryWithoutStatus := unionQuery + ` ORDER BY booking_date DESC LIMIT ? OFFSET ?`

	list, err := t.fetchWithJoin(ctx, queryWithoutStatus, limit, offset)
	if status != "" {
		if status == "pending" {
			transactionStatus = 0
		} else if status == "waitingApproval" {
			transactionStatus = 1
		} else if status == "confirm" {
			transactionStatus = 2
		}
		querySt := query + ` AND t.status = ` + strconv.Itoa(transactionStatus)
		queryTSt := queryT + ` AND t.status = ` + strconv.Itoa(transactionStatus)
		unionQuery = querySt + ` UNION ` + queryTSt
		queryWithStatus := unionQuery + ` ORDER BY booking_date DESC LIMIT ? OFFSET ?`
		list, err = t.fetchWithJoin(ctx, queryWithStatus, limit, offset)

		if status == "failed" {
			transactionStatus = 3
			cancelledStatus := 4
			querySt = query + ` AND t.status IN (` + strconv.Itoa(transactionStatus) + `,` + strconv.Itoa(cancelledStatus) + `)`
			queryTSt = queryT + ` AND t.status IN (` + strconv.Itoa(transactionStatus) + `,` + strconv.Itoa(cancelledStatus) + `)`
			unionQuery = querySt + ` UNION ` + queryTSt
			fmt.Println(unionQuery)
			queryWithStatus = unionQuery + ` ORDER BY booking_date DESC LIMIT ? OFFSET ?`
			list, err = t.fetchWithJoin(ctx, queryWithStatus, limit, offset)
		}

		if status == "boarded" {
			transactionStatus = 1
			bookingStatus = 3
			querySt = query + ` AND t.status = ` + strconv.Itoa(transactionStatus) + ` AND b.status = ` + strconv.Itoa(bookingStatus)
			queryTSt = queryT + ` AND t.status = ` + strconv.Itoa(transactionStatus) + ` AND b.status = ` + strconv.Itoa(bookingStatus)
			unionQuery = querySt + ` UNION ` + queryTSt
			queryWithStatus = unionQuery + ` ORDER BY booking_date DESC LIMIT ? OFFSET ?`
			list, err = t.fetchWithJoin(ctx, queryWithStatus, limit, offset)
		}
	}
	if err != nil {
		fmt.Println("err1", err.Error())
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
			&t.MerchantId,
			&t.OrderIdBook,
			&t.BookedBy,
			&t.ExpTitle,
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

func (t transactionRepository) Count(ctx context.Context, startDate, endDate, search, status string) (int, error) {
	query := `
	SELECT
		count(*) as count
	FROM 
		transactions t
		JOIN booking_exps b ON t.booking_exp_id = b.id
		JOIN experiences e ON b.exp_id = e.id
		JOIN experience_payments ep ON e.id = ep.exp_id
	WHERE
		t.is_deleted = 0
		AND t.is_active = 1`

	queryT := `
	SELECT
		count(*) as count
	FROM
		transactions t
		JOIN booking_exps b ON t.booking_exp_id = b.id
		JOIN transportations tr ON b.trans_id = tr.id
	WHERE
		t.is_deleted = 0
		AND t.is_active = 1`

	if search != "" {
		keyword := `'%` + search + `%'`
		query = query + ` AND (LOWER(b.booked_by) LIKE LOWER(` + keyword + `) OR LOWER(b.order_id) LIKE LOWER(` + keyword + `))`
		queryT = queryT + ` AND (LOWER(b.booked_by) LIKE LOWER(` + keyword + `) OR LOWER(b.order_id) LIKE LOWER(` + keyword + `))`
	}
	if startDate != "" && endDate != "" {
		query = query + ` AND DATE(b.created_date) BETWEEN '` + startDate + `' AND '` + endDate + `'`
		queryT = queryT + ` AND DATE(b.created_date) BETWEEN '` + startDate + `' AND '` + endDate + `'`
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
			transactionStatus = 1
			bookingStatus := 3
			querySt = query + ` AND t.status = ` + strconv.Itoa(transactionStatus) + ` AND b.status = ` + strconv.Itoa(bookingStatus)
			queryTSt = queryT + ` AND t.status = ` + strconv.Itoa(transactionStatus) + ` AND b.status = ` + strconv.Itoa(bookingStatus)
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

func (m transactionRepository) GetById(ctx context.Context, id string) (*models.TransactionWMerchant, error) {
	query := `SELECT t.*,e.merchant_id,b.order_id as order_id_book,b.booked_by,e.exp_title FROM transactions t
				join booking_exps b on t.booking_exp_id = b.id
				join experiences e on b.exp_id = e.id WHERE t.id = ?`

	list, err := m.fetch(ctx, query, id)
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
