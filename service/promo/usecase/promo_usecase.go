package usecase

import (
	"github.com/models"
	"github.com/service/promo"
	"golang.org/x/net/context"
	"time"
)

type promoUsecase struct {
	promoRepo      promo.Repository
	contextTimeout time.Duration
}

// NewPromoUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewPromoUsecase(p promo.Repository, timeout time.Duration) promo.Usecase {
	return &promoUsecase{
		promoRepo:      p,
		contextTimeout: timeout,
	}
}

func (p promoUsecase) Fetch(ctx context.Context, page *int, size *int) ([]*models.PromoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	promoList, err := p.promoRepo.Fetch(ctx, page, size)
	if err != nil {
		return nil, err
	}
	var promoDto []*models.PromoDto
	for _, element := range promoList {
		resPromo := models.PromoDto{
			Id:         element.Id,
			PromoCode:  element.PromoCode,
			PromoName:  element.PromoName,
			PromoDesc:  element.PromoDesc,
			PromoValue: element.PromoValue,
			PromoType:  element.PromoType,
			PromoImage: element.PromoImage,
		}
		promoDto = append(promoDto, &resPromo)
	}

	return promoDto, nil
}

func (p promoUsecase) GetByCode(ctx context.Context, code string) (*models.PromoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	promos, err := p.promoRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	promoDto := make([]*models.PromoDto, len(promos))
	for i, p := range promos {
		promoDto[i] = &models.PromoDto{
			Id:         p.Id,
			PromoCode:  p.PromoCode,
			PromoName:  p.PromoName,
			PromoDesc:  p.PromoDesc,
			PromoValue: p.PromoValue,
			PromoType:  p.PromoType,
			PromoImage: p.PromoImage,
		}
	}

	return promoDto[0], nil
}
