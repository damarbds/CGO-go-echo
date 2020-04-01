package usecase

import (
	"github.com/models"
	"github.com/product/experience_add_ons"
	"golang.org/x/net/context"
	"time"
)

type experience_add_ons_Usecase struct {
	experience_add_onsRepo experience_add_ons.Repository
	contextTimeout         time.Duration
}


// NewharborsUsecase will create new an harborsUsecase object representation of harbors.Usecase interface
func NewharborsUsecase(a experience_add_ons.Repository, timeout time.Duration) experience_add_ons.Usecase {
	return &experience_add_ons_Usecase{
		experience_add_onsRepo:    a,
		contextTimeout: timeout,
	}
}


func (e experience_add_ons_Usecase) GetByExpId(ctx context.Context, exp_id string) ([]*models.ExperienceAddOnDto, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	res, err := e.experience_add_onsRepo.GetByExpId(ctx,exp_id)
	if err != nil {
		return nil, err
	}
	var ExpExperienceAddOns []*models.ExperienceAddOnDto
	for _, element := range res {
		var currency string
		if element.Currency == 1 {
			currency = "USD"
		} else {
			currency = "IDR"
		}
		expExperienceAddOn := models.ExperienceAddOnDto{
			Id:      element.Id,
			Name:     element.Name,
			Desc:     element.Desc,
			Currency: currency,
			Amount:   element.Amount,
			ExpId:    element.ExpId,
		}
		ExpExperienceAddOns = append(ExpExperienceAddOns,&expExperienceAddOn)
	}

	return ExpExperienceAddOns, nil
}
