package usecase

import (
	"context"
	"math"
	"time"

	"github.com/auth/merchant"
	"github.com/transactions/balance_history"

	"github.com/models"
)

type balanceHistoryUsecase struct {
	merchantUsecase    merchant.Usecase
	balanceHistoryRepo balance_history.Repository
	contextTimeout     time.Duration
}

func NewBalanceHistoryUsecase(bh balance_history.Repository, mu merchant.Usecase, timeout time.Duration) balance_history.Usecase {
	return &balanceHistoryUsecase{
		balanceHistoryRepo: bh,
		merchantUsecase:    mu,
		contextTimeout:     timeout,
	}
}
func (b balanceHistoryUsecase) List(c context.Context, merchantId, status string, page int, limit, offset *int, month, year string) (*models.BalanceHistoryDtoWithPagination, error) {
	ctx, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	getBhistory, err := b.balanceHistoryRepo.GetAll(ctx, merchantId, status, limit, offset, month, year)
	if err != nil {
		return nil, err
	}
	var result models.BalanceHistoryDtoWithPagination
	for _, element := range getBhistory {
		dto := models.BalanceHistoryDto{
			Id:            element.Id,
			MerchantId:    element.MerchantId,
			Status:        element.Status,
			Amount:        element.Amount,
			DateOfRequest: element.DateOfRequest,
			DateOfPayment: element.DateOfPayment,
			Remarks:       element.Remarks,
		}
		result.Data = append(result.Data, &dto)
	}

	totalRecords, _ := b.balanceHistoryRepo.Count(ctx, merchantId, status)
	var totalPage int
	var prev int
	var next int
	if limit != nil {
		totalPage = int(math.Ceil(float64(totalRecords) / float64(*limit)))
		prev = page
		next = page
		if page != 1 {
			prev = page - 1
		}

		if page != totalPage {
			next = page + 1
		}
	}

	meta := &models.MetaPagination{
		Page:          page,
		Total:         totalPage,
		TotalRecords:  totalRecords,
		Prev:          prev,
		Next:          next,
		RecordPerPage: len(getBhistory),
	}
	result.Meta = meta

	return &result, nil
}
func (b balanceHistoryUsecase) Create(c context.Context, bHistory models.NewBalanceHistoryCommand, token string) (*models.NewBalanceHistoryCommand, error) {
	ctx, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	currentMerchant, err := b.merchantUsecase.ValidateTokenMerchant(ctx, token)
	if err != nil {
		return nil, err
	}
	layoutFormat := "2006-01-02 15:04:05"

	dateOfPayment, err := time.Parse(layoutFormat, bHistory.DateOfPayment)
	balance := models.BalanceHistory{
		Id:            "",
		CreatedBy:     currentMerchant.MerchantEmail,
		CreatedDate:   time.Now(),
		ModifiedBy:    nil,
		ModifiedDate:  nil,
		DeletedBy:     nil,
		DeletedDate:   nil,
		IsDeleted:     0,
		IsActive:      0,
		Status:        bHistory.Status,
		MerchantId:    currentMerchant.Id,
		Amount:        bHistory.Amount,
		DateOfRequest: time.Now(),
		DateOfPayment: dateOfPayment,
		Remarks:       bHistory.Remarks,
	}
	create, err := b.balanceHistoryRepo.Insert(ctx, balance)
	if err != nil {
		return nil, err
	}
	bHistory.Id = *create
	return &bHistory, nil
}
