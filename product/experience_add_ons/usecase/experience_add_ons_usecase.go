package usecase

import (
	"github.com/misc/currency"
	"github.com/models"
	"github.com/product/experience_add_ons"
	"golang.org/x/net/context"
	"time"
)

type experience_add_ons_Usecase struct {
	experience_add_onsRepo experience_add_ons.Repository
	contextTimeout         time.Duration
	currencyUsecase	currency.Usecase
}

// NewharborsUsecase will create new an harborsUsecase object representation of harbors.Usecase interface
func NewharborsUsecase(currencyUsecase	currency.Usecase,a experience_add_ons.Repository, timeout time.Duration) experience_add_ons.Usecase {
	return &experience_add_ons_Usecase{
		currencyUsecase:currencyUsecase,
		experience_add_onsRepo: a,
		contextTimeout:         timeout,
	}
}

func (e experience_add_ons_Usecase) GetByExpId(ctx context.Context, exp_id string,currencyPrice string) ([]*models.ExperienceAddOnDto, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	res, err := e.experience_add_onsRepo.GetByExpId(ctx, exp_id)
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
		if currencyPrice == "USD"{
			if currency == "IDR"{
				convertCurrency ,_ := e.currencyUsecase.ExchangeRatesApi(ctx,"IDR","USD")
				calculatePrice := convertCurrency.Rates.USD * element.Amount
				element.Amount = calculatePrice
				currency = "USD"
			}
		}else if currencyPrice =="IDR"{
			if currency == "USD"{
				convertCurrency ,_ := e.currencyUsecase.ExchangeRatesApi(ctx,"USD","IDR")
				calculatePrice := convertCurrency.Rates.IDR * element.Amount
				element.Amount = calculatePrice
				currency = "IDR"
			}
		}
		expExperienceAddOn := models.ExperienceAddOnDto{
			Id:       element.Id,
			Name:     element.Name,
			Desc:     element.Desc,
			Currency: currency,
			Amount:   element.Amount,
			ExpId:    element.ExpId,
		}
		ExpExperienceAddOns = append(ExpExperienceAddOns, &expExperienceAddOn)
	}

	return ExpExperienceAddOns, nil
}
