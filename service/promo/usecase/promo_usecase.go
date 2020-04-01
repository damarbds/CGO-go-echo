package usecase

import (
	"github.com/models"
	"github.com/service/promo"
	"golang.org/x/net/context"
	"time"
)

type promoUsecase struct {
	promoRepo     promo.Repository
	contextTimeout time.Duration
}

// NewArticleUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewArticleUsecase(p promo.Repository,timeout time.Duration) promo.Usecase {
	return &promoUsecase{
		promoRepo:p,
		contextTimeout: timeout,
	}
}

func (p promoUsecase) Fetch(ctx context.Context, page *int, size *int) ([]*models.PromoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	promoList, err := p.promoRepo.Fetch(ctx,page,size)
	if err != nil {
		return nil, err
	}
	var promoDto []*models.PromoDto
	for _, element := range promoList {
		promo:=models.PromoDto{
			Id:         element.Id,
			PromoCode:  element.PromoCode,
			PromoName:  element.PromoName,
			PromoDesc:  element.PromoDesc,
			PromoValue: element.PromoValue,
			PromoType:  element.PromoType,
			PromoImage: element.PromoImage,
		}
		promoDto = append(promoDto,&promo)
	}
	return promoDto,nil
}
