package usecase

import (
	"context"
	"github.com/models"
	"github.com/service/transportation"
	"time"
)

type transportationUsecase struct {
	transRepo transportation.Repository
	contextTimeout  time.Duration
}

func NewTransportationUsecase(t transportation.Repository, timeout time.Duration) transportation.Usecase {
	return &transportationUsecase{
		transRepo:      t,
		contextTimeout: timeout,
	}
}

func (t transportationUsecase) List(ctx context.Context) ([]*models.TimeOptionDto, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	list, err := t.transRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	timeOptions := make([]*models.TimeOptionDto, len(list))
	for i, item := range list {
		timeOptions[i] = &models.TimeOptionDto{
			Id:        item.Id,
			StartTime: item.StartTime,
			EndTime:   item.EndTime,
		}
	}

	return timeOptions, nil
}