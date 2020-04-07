package usecase

import (
	"context"
	"time"

	"github.com/models"
	"github.com/service/facilities"
)

type facilityUsecase struct {
	facilityRepo   facilities.Repository
	contextTimeout time.Duration
}

func NewFacilityUsecase(f facilities.Repository, timeout time.Duration) facilities.Usecase {
	return &facilityUsecase{
		facilityRepo:   f,
		contextTimeout: timeout,
	}
}

func (f facilityUsecase) List(ctx context.Context) ([]*models.FacilityDto, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	res, err := f.facilityRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]*models.FacilityDto, len(res))
	for i, n := range res {
		results[i] = &models.FacilityDto{
			Id:           n.Id,
			FacilityName: n.FacilityName,
			FacilityIcon: n.FacilityIcon.String,
			IsNumerable:  n.IsNumerable,
		}
	}

	return results, nil
}
