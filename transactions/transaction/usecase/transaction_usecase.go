package usecase

import (
	"context"
	"encoding/json"
	"math"
	"strings"
	"time"

	"github.com/models"
	"github.com/service/exp_payment"
	"github.com/transactions/transaction"
)

type transactionUsecase struct {
	experiencePaymentTypeRepo exp_payment.Repository
	transactionRepo           transaction.Repository
	contextTimeout            time.Duration
}

func NewTransactionUsecase(ep exp_payment.Repository, t transaction.Repository, timeout time.Duration) transaction.Usecase {
	return &transactionUsecase{
		experiencePaymentTypeRepo: ep,
		transactionRepo:           t,
		contextTimeout:            timeout,
	}
}

func (t transactionUsecase) CountThisMonth(ctx context.Context) (*models.TotalTransaction, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	total, err := t.transactionRepo.CountThisMonth(ctx)
	if err != nil {
		return nil, err
	}

	return total, nil
}

func (t transactionUsecase) List(ctx context.Context, startDate, endDate, search, status string, page, limit, offset int) (*models.TransactionWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	list, err := t.transactionRepo.List(ctx, startDate, endDate, search, status, limit, offset)
	if err != nil {
		return nil, err
	}

	transactions := make([]*models.TransactionDto, len(list))
	for i, item := range list {
		var expType []string
		if item.ExpType != "" {
			if !strings.Contains(item.ExpType, "]") {
				// Default type for Transportation
				expType = []string{"Transportation"}
			} else {
				if errUnmarshal := json.Unmarshal([]byte(item.ExpType), &expType); errUnmarshal != nil {
					return nil, errUnmarshal
				}
			}
		}
		var experiencePaymentType *models.ExperiencePaymentTypeDto
		if item.ExperiencePaymentId != "" {
			if item.ExperiencePaymentId == "Economy" || item.ExperiencePaymentId == "Executive" {
				// Default Payment Type for Transportation
				experiencePaymentType = &models.ExperiencePaymentTypeDto{
					Id:   "8a5e3eef-a6db-4584-a280-af5ab18a979b",
					Name: "Full Payment",
					Desc: "Full Payment",
				}
			} else {
				query, err := t.experiencePaymentTypeRepo.GetByExpID(ctx, item.ExpId)
				if err != nil {
					return nil, err
				}
				for _, element := range query {
					if element.Id == item.ExperiencePaymentId {
						paymentType := models.ExperiencePaymentTypeDto{
							Id:   element.ExpPaymentTypeId,
							Name: element.ExpPaymentTypeName,
							Desc: element.ExpPaymentTypeDesc,
						}
						experiencePaymentType = &paymentType
					}
				}
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

		var status string
		if item.TransactionStatus == 0 {
			status = "Pending"
		} else if item.TransactionStatus == 1 {
			status = "Waiting approval"
		} else if item.TransactionStatus == 2 {
			status = "Confirm"
		} else if item.TransactionStatus == 3 || item.TransactionStatus == 4 {
			status = "Failed"
		} else if item.TransactionStatus == 1 && item.BookingStatus == 3 {
			status = "Boarded"
		}

		transactions[i] = &models.TransactionDto{
			TransactionId:         item.TransactionId,
			ExpId:                 item.ExpId,
			ExpTitle:              item.ExpTitle,
			ExpType:               expType,
			BookingExpId:          item.BookingExpId,
			BookingCode:           item.BookingCode,
			BookingDate:           item.BookingDate,
			CheckInDate:           item.CheckInDate,
			BookedBy:              bookedBy,
			Guest:                 len(guestDesc),
			Email:                 item.Email,
			Status:                status,
			TotalPrice:            item.TotalPrice,
			ExperiencePaymentType: experiencePaymentType,
		}
	}
	totalRecords, _ := t.transactionRepo.Count(ctx, startDate, endDate, search, status)
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
