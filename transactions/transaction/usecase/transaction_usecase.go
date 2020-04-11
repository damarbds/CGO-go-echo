package usecase

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"github.com/models"
	"github.com/transactions/transaction"
)

type transactionUsecase struct {
	transactionRepo transaction.Repository
	contextTimeout  time.Duration
}

func NewTransactionUsecase(t transaction.Repository, timeout time.Duration) transaction.Usecase {
	return &transactionUsecase{
		transactionRepo: t,
		contextTimeout:  timeout,
	}
}

func (t transactionUsecase) List(ctx context.Context, status string, page, limit, offset int) (*models.TransactionWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	list, err := t.transactionRepo.List(ctx, status, limit, offset)
	if err != nil {
		return nil, err
	}

	transactions := make([]*models.TransactionDto, len(list))
	for i, item := range list {
		var expType []string
		if item.ExpType != "" {
			if errUnmarshal := json.Unmarshal([]byte(item.ExpType), &expType); errUnmarshal != nil {
				return nil, errUnmarshal
			}
		}

		var guestDesc []models.GuestDescObj
		if item.GuestDesc != "" {
			if errUnmarshal := json.Unmarshal([]byte(item.GuestDesc), &guestDesc); errUnmarshal != nil {
				return nil, errUnmarshal
			}
		}

		var bookedBy []models.BookedByObj
		if item.BookedBy != "" {
			if errUnmarshal := json.Unmarshal([]byte(item.BookedBy), &bookedBy); errUnmarshal != nil {
				return nil, errUnmarshal
			}
		}

		transactions[i] = &models.TransactionDto{
			TransactionId: item.TransactionId,
			ExpId:         item.ExpId,
			ExpType:       expType,
			BookingExpId:  item.BookingExpId,
			BookingCode:   item.BookingCode,
			BookingDate:   item.BookingDate,
			CheckInDate:   item.CheckInDate,
			BookedBy:      bookedBy,
			Guest:         len(guestDesc),
			Email:         item.Email,
		}
	}
	totalRecords, _ := t.transactionRepo.Count(ctx, status)
	totalPage := int(math.Ceil(float64(totalRecords) / float64(limit)))
	prev := page
	next := page
	if page != 1 {
		prev = page - 1
	}

	if page != totalPage {
		next = page + 1
	}
	meta := &models.MetaPagination{
		Page:          page,
		Total:         totalPage,
		TotalRecords:  totalRecords,
		Prev:          prev,
		Next:          next,
		RecordPerPage: len(list),
	}

	response := &models.TransactionWithPagination{
		Data: transactions,
		Meta: meta,
	}

	return response, nil
}

func (t transactionUsecase) CountSuccess(ctx context.Context) (*models.Count, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	count, err := t.transactionRepo.CountSuccess(ctx)
	if err != nil {
		return nil, err
	}

	return &models.Count{Count: count}, nil
}
