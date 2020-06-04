package usecase

import (
	"github.com/service/exclusion_service"
	"math"
	"time"

	"github.com/auth/admin"
	"github.com/models"
	"golang.org/x/net/context"
)

type exclusionServicesUsecase struct {
	exclusionServicesRepo exclusion_service.Repository
	adminUsecase   admin.Usecase
	contextTimeout time.Duration
}

func (m exclusionServicesUsecase) List(ctx context.Context, page, limit, offset int, token string) (*models.ExclusionServiceWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	//_, err := m.adminUsecase.ValidateTokenAdmin(ctx, token)
	//if err != nil {
	//	return nil, models.ErrUnAuthorize
	//}

	list, err := m.exclusionServicesRepo.Fetch(ctx, limit,offset)
	if err != nil {
		return nil, err
	}

	exclusionService := make([]*models.ExclusionServiceDto, len(list))
	for i, item := range list {
		exclusionService[i] = &models.ExclusionServiceDto{
			Id:                   item.Id,
			ExclusionServiceName: item.ExclusionServiceName,
			ExclusionServiceType: item.ExclusionServiceType,
		}
	}
	totalRecords, _ := m.exclusionServicesRepo.GetCount(ctx)
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

	response := &models.ExclusionServiceWithPagination{
		Data: exclusionService,
		Meta: meta,
	}

	return response, nil
}

// NewPromoUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewExclusionServicesUsecase(exclusionServicesRepo exclusion_service.Repository,au admin.Usecase, timeout time.Duration) exclusion_service.Usecase {
	return &exclusionServicesUsecase{
		exclusionServicesRepo : exclusionServicesRepo,
		adminUsecase:   au,
		contextTimeout: timeout,
	}
}