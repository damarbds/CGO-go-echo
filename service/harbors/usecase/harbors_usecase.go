package usecase

import (
	"context"
	"time"

	"github.com/service/harbors"

	"github.com/models"
)

type harborsUsecase struct {
	harborsRepo    harbors.Repository
	contextTimeout time.Duration
}

// NewharborsUsecase will create new an harborsUsecase object representation of harbors.Usecase interface
func NewharborsUsecase(a harbors.Repository, timeout time.Duration) harbors.Usecase {
	return &harborsUsecase{
		harborsRepo:    a,
		contextTimeout: timeout,
	}
}
func (m harborsUsecase) GetAllWithJoinCPC(c context.Context, page *int,size *int) ([]*models.HarborsWCPCDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	res, err := m.harborsRepo.GetAllWithJoinCPC(ctx, page,size)
	if err != nil {
		return nil, err
	}
	var harborss []*models.HarborsWCPCDto
	for _, element := range res {
		harbors := models.HarborsWCPCDto{
			Id:               element.Id,
			HarborsName:      element.HarborsName,
			HarborsLongitude: element.HarborsLongitude,
			HarborsLatitude:  element.HarborsLatitude,
			HarborsImage:     element.HarborsImage,
			City:             element.CityName,
			Province:         element.ProvinceName,
			Country:          element.CountryName,
		}
		harborss = append(harborss, &harbors)
	}

	return harborss, nil
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
