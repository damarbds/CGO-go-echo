package usecase_test

import (
	"context"
	"errors"
	_currencyUsecaseMock "github.com/misc/currency/mocks"
	"github.com/models"
	"github.com/product/experience_add_ons/mocks"
	"github.com/product/experience_add_ons/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

var (
	timeoutContext = time.Second * 30
	mockExperienceAddOn = []*models.ExperienceAddOn{
		&models.ExperienceAddOn{
			Id:           "asdasdasd",
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			Name:         "qeqwe",
			Desc:         "adasd",
			Currency:     1,
			Amount:       11212,
			ExpId:        "sfdsf",
		},
	}
)

func TestGetByExpId(t *testing.T) {
	mockExperienceAddOnDtoRepo := new(mocks.Repository)
	mockCurrencyUsecase := new(_currencyUsecaseMock.Usecase)
	mockcurrencyication := models.CurrencyExChangeRate{
		Date:  "2020-09-17",
		Base:  "USD",
		Rates: models.Rates{
			IDR: 12312312,
			USD: 3242,
		},
	}
	t.Run("success-with-currency-idr", func(t *testing.T) {
		mockExperienceAddOnDtoRepo.On("GetByExpId", mock.Anything, mock.AnythingOfType("string")).Return(mockExperienceAddOn, nil).Once()
		mockCurrencyUsecase.On("ExchangeRatesApi", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(mockcurrencyication, nil).Once()
		u := usecase.NewharborsUsecase(mockCurrencyUsecase,mockExperienceAddOnDtoRepo, timeoutContext)

		a, err := u.GetByExpId(context.TODO(), mockExperienceAddOn[0].ExpId,"IDR")

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockExperienceAddOnDtoRepo.AssertExpectations(t)
	})
	t.Run("success-with-currency-usd", func(t *testing.T) {
		mockExperienceAddOnDtoRepo.On("GetByExpId", mock.Anything, mock.AnythingOfType("string")).Return(mockExperienceAddOn, nil).Once()
		mockCurrencyUsecase.On("ExchangeRatesApi", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(mockcurrencyication, nil).Once()
		u := usecase.NewharborsUsecase(mockCurrencyUsecase,mockExperienceAddOnDtoRepo, timeoutContext)

		a, err := u.GetByExpId(context.TODO(), mockExperienceAddOn[0].ExpId,"USD")

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockExperienceAddOnDtoRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockExperienceAddOnDtoRepo.On("GetByExpId", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("UnExpected")).Once()
		mockCurrencyUsecase.On("ExchangeRatesApi", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(mockcurrencyication, nil).Once()
		u := usecase.NewharborsUsecase(mockCurrencyUsecase,mockExperienceAddOnDtoRepo, timeoutContext)

		_, err := u.GetByExpId(context.TODO(), mockExperienceAddOn[0].ExpId,"USD")

		assert.Error(t, err)
		//assert.NotNil(t, a)

		mockExperienceAddOnDtoRepo.AssertExpectations(t)
	})

}
