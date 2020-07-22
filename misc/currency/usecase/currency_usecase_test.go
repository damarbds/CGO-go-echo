package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/misc/currency/mocks"
	"github.com/misc/currency/usecase"
	"github.com/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	timeoutContext = time.Second * 30
)

func TestExchangeRatesApi(t *testing.T) {
	mockExChangeRateRepo := new(mocks.Repository)
	mockExChangeRate := models.ExChangeRate{
			Id:    1,
			Date:  "2020-09-17",
			From:  "USD",
			To:    "IDR",
			Rates: 1212112,
		}

	t.Run("success-from-usd", func(t *testing.T) {
		mockExChangeRateRepo.On("GetByDate", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(mockExChangeRate, nil).Once()
		u := usecase.NewCurrencyUsecase(mockExChangeRateRepo, timeoutContext)

		a, err := u.ExchangeRatesApi(context.TODO(), mockExChangeRate.From,mockExChangeRate.To)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockExChangeRateRepo.AssertExpectations(t)
	})
	t.Run("success-from-idr", func(t *testing.T) {
		mockExChangeRateRepo.On("GetByDate", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(mockExChangeRate, nil).Once()
		u := usecase.NewCurrencyUsecase(mockExChangeRateRepo, timeoutContext)

		a, err := u.ExchangeRatesApi(context.TODO(), mockExChangeRate.To,mockExChangeRate.From)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockExChangeRateRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockExChangeRateRepo.On("GetByDate", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(nil, errors.New("UnExpected")).Once()
		u := usecase.NewCurrencyUsecase(mockExChangeRateRepo, timeoutContext)

		a, err := u.ExchangeRatesApi(context.TODO(), mockExChangeRate.From,mockExChangeRate.To)

		assert.Error(t, err)
		assert.Nil(t, a)

		mockExChangeRateRepo.AssertExpectations(t)
	})
}
