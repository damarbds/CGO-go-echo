package usecase

import (
	"context"
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

func (t transactionUsecase) CountSuccess(ctx context.Context) (*models.Count, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	count, err := t.transactionRepo.CountSuccess(ctx)
	if err != nil {
		return nil, err
	}

	return &models.Count{Count: count}, nil
}
