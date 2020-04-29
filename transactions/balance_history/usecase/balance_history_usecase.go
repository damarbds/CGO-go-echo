package usecase

import (
	"context"
	"github.com/auth/admin"
	"math"
	"time"

	"github.com/auth/merchant"
	"github.com/transactions/balance_history"

	"github.com/models"
)

type balanceHistoryUsecase struct {
	merchantRepo 		merchant.Repository
	adminUsecase 		admin.Usecase
	merchantUsecase    merchant.Usecase
	balanceHistoryRepo balance_history.Repository
	contextTimeout     time.Duration
}


func NewBalanceHistoryUsecase(mr merchant.Repository,au admin.Usecase,bh balance_history.Repository, mu merchant.Usecase, timeout time.Duration) balance_history.Usecase {
	return &balanceHistoryUsecase{
		merchantRepo:mr,
		adminUsecase:au,
		balanceHistoryRepo: bh,
		merchantUsecase:    mu,
		contextTimeout:     timeout,
	}
}
func (b balanceHistoryUsecase) UpdateAmount(c context.Context, amount models.NewBalanceHistoryAmountCommand, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	currentUser , err := b.adminUsecase.ValidateTokenAdmin(ctx,token)
	if err != nil {
		return nil,models.ErrUnAuthorize
	}

	getBalance , err := b.balanceHistoryRepo.GetById(ctx,amount.Id)

	balanceHistory := models.BalanceHistory{
		Id:            getBalance.Id,
		CreatedBy:     "",
		CreatedDate:   time.Time{},
		ModifiedBy:    &currentUser.Name,
		ModifiedDate:  nil,
		DeletedBy:     nil,
		DeletedDate:   nil,
		IsDeleted:     0,
		IsActive:      0,
		Status:        getBalance.Status,
		MerchantId:    getBalance.MerchantId,
		Amount:        amount.Amount,
		DateOfRequest: getBalance.DateOfRequest,
		DateOfPayment: getBalance.DateOfPayment,
		Remarks:       getBalance.Remarks,
	}

	_, err = b.balanceHistoryRepo.Update(ctx,balanceHistory)
	if err != nil {
		return nil,err
	}
	result := models.ResponseDelete{
		Id:      balanceHistory.Id,
		Message: "Success Update Amount",
	}
	return &result,nil
}
func (b balanceHistoryUsecase) ConfirmWithdraw(c context.Context,command models.NewBalanceHistoryConfirmCommand,token string)(*models.ResponseDelete,error) {
	ctx, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	currentUser , err := b.adminUsecase.ValidateTokenAdmin(ctx,token)
	if err != nil {
		return nil,models.ErrUnAuthorize
	}

	getBalance , err := b.balanceHistoryRepo.GetById(ctx,command.Id)

	var status int
	if command.Action == "accept"{
		status = 2
		getMerchant,err := b.merchantRepo.GetByID(ctx,command.MerchantId)
		if err != nil {
			return nil,err
		}
		getMerchant.ModifiedBy = &currentUser.Name
		getMerchant.Balance = getMerchant.Balance - command.Amount
		err = b.merchantRepo.Update(ctx,getMerchant)
		if err != nil {
			return nil,err
		}
	}else if command.Action == "decline"{
		status = 0
		command.Amount = getBalance.Amount
	}
	balanceHistory := models.BalanceHistory{
		Id:            getBalance.Id,
		CreatedBy:     "",
		CreatedDate:   time.Time{},
		ModifiedBy:    &currentUser.Name,
		ModifiedDate:  nil,
		DeletedBy:     nil,
		DeletedDate:   nil,
		IsDeleted:     0,
		IsActive:      0,
		Status:        status,
		MerchantId:    command.MerchantId,
		Amount:        command.Amount,
		DateOfRequest: getBalance.DateOfRequest,
		DateOfPayment: time.Now(),
		Remarks:       getBalance.Remarks,
	}

	_, err = b.balanceHistoryRepo.Update(ctx,balanceHistory)
	if err != nil {
		return nil,err
	}
	result := models.ResponseDelete{
		Id:      balanceHistory.Id,
		Message: "Success Confirm With Admin",
	}
	return &result,nil
}

func (b balanceHistoryUsecase) List(c context.Context, merchantId, status string, page int, limit, offset *int, month, year string,token string,isAdmin bool) (*models.BalanceHistoryDtoWithPagination, error) {
	ctx, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()
	getBhistory := make([]*models.BalanceHistory,0)
	if isAdmin == true {

		_,err := b.adminUsecase.ValidateTokenAdmin(ctx ,token)
		if err != nil {
			return nil,models.ErrUnAuthorize
		}
		query, err := b.balanceHistoryRepo.GetAll(ctx, merchantId, status, limit, offset, month, year)
		if err != nil {
			return nil, err
		}
		getBhistory = query
	}else {

		_,err := b.merchantUsecase.ValidateTokenMerchant(ctx ,token)
		if err != nil {
			return nil,models.ErrUnAuthorize
		}
		query, err := b.balanceHistoryRepo.GetAll(ctx, merchantId, status, limit, offset, month, year)
		if err != nil {
			return nil, err
		}
		getBhistory = query
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
