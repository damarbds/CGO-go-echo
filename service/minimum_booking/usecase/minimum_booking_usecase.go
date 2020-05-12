package usecase

import (
	"context"
	"github.com/service/minimum_booking"
	"math"
	"time"

	"github.com/models"
)

type minimumBookingUsecase struct {
	minimumBookingRepo minimum_booking.Repository
	contextTimeout time.Duration
}


// NewharborsUsecase will create new an harborsUsecase object representation of harbors.Usecase interface
func NewminimumBookingUsecase(minimumBookingRepo minimum_booking.Repository, timeout time.Duration) minimum_booking.Usecase {
	return &minimumBookingUsecase{
		minimumBookingRepo:minimumBookingRepo,
		contextTimeout: timeout,
	}
}

func (m minimumBookingUsecase) GetAll(c context.Context, page, limit, size int) (*models.MinimumBookingDtoWithPagination, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	getAll ,err := m.minimumBookingRepo.GetAll(ctx,limit,size)
	if err != nil {
		return nil,err
	}

	minimumBookingDtos := make([]*models.MinimumBookingDto,0)
	for _,element := range getAll{

		dto := models.MinimumBookingDto{
			Id:                   element.Id,
			MinimumBookingDesc:   element.MinimumBookingDesc,
			MinimumBookingAmount: element.MinimumBookingAmount,
		}

		minimumBookingDtos = append(minimumBookingDtos,&dto)
	}
	totalRecords, _ := m.minimumBookingRepo.GetCount(ctx)

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
		RecordPerPage: len(minimumBookingDtos),
	}

	response := &models.MinimumBookingDtoWithPagination{
		Data: minimumBookingDtos,
		Meta: meta,
	}
	return response, nil
}
