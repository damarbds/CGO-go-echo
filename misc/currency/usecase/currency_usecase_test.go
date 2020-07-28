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
		mockExChangeRateRepo.On("GetByDate", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(&mockExChangeRate, nil).Once()
		u := usecase.NewCurrencyUsecase(mockExChangeRateRepo, timeoutContext)

		a, err := u.ExchangeRatesApi(context.TODO(), mockExChangeRate.From,mockExChangeRate.To)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockExChangeRateRepo.AssertExpectations(t)
	})
	t.Run("success-from-idr", func(t *testing.T) {
		mockExChangeRateRepo.On("GetByDate", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(&mockExChangeRate, nil).Once()
		u := usecase.NewCurrencyUsecase(mockExChangeRateRepo, timeoutContext)

		a, err := u.ExchangeRatesApi(context.TODO(), mockExChangeRate.To,mockExChangeRate.From)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockExChangeRateRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockExChangeRateRepo.On("GetByDate", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(nil, errors.New("UnExpected")).Once()
		u := usecase.NewCurrencyUsecase(mockExChangeRateRepo, timeoutContext)

		_, err := u.ExchangeRatesApi(context.TODO(), mockExChangeRate.From,mockExChangeRate.To)

		assert.Error(t, err)
		//assert.Nil(t, a)

		mockExChangeRateRepo.AssertExpectations(t)
	})
}

func TestExchangeRatesWithApi(t *testing.T) {
	mockExChangeRateRepo := new(mocks.Repository)
	mockExChangeRate := models.ExChangeRate{
		Id:    1,
		Date:  "2020-09-17",
		From:  "USD",
		To:    "IDR",
		Rates: 1212112,
	}

	t.Run("success-from-usd", func(t *testing.T) {
		//mockExChangeRateRepo.On("GetByDate", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(&mockExChangeRate, nil).Once()
		u := usecase.NewCurrencyUsecase(mockExChangeRateRepo, timeoutContext)

		a, err := u.ExchangeRatesWithApi(context.TODO(), mockExChangeRate.From,mockExChangeRate.To)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockExChangeRateRepo.AssertExpectations(t)
	})
	t.Run("success-from-idr", func(t *testing.T) {
		//mockExChangeRateRepo.On("GetByDate", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(&mockExChangeRate, nil).Once()
		u := usecase.NewCurrencyUsecase(mockExChangeRateRepo, timeoutContext)

		a, err := u.ExchangeRatesWithApi(context.TODO(), mockExChangeRate.To,mockExChangeRate.From)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockExChangeRateRepo.AssertExpectations(t)
	})
	//t.Run("error-failed", func(t *testing.T) {
	//	//mockExChangeRateRepo.On("GetByDate", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(nil, errors.New("UnExpected")).Once()
	//	u := usecase.NewCurrencyUsecase(mockExChangeRateRepo, timeoutContext)
	//
	//	_, err := u.ExchangeRatesWithApi(context.TODO(), mockExChangeRate.From,mockExChangeRate.To)
	//
	//	assert.Error(t, err)
	//	//assert.Nil(t, a)
	//
	//	mockExChangeRateRepo.AssertExpectations(t)
	//})
}

func TestExchangeFreeCurrconv(t *testing.T) {
	mockExChangeRateRepo := new(mocks.Repository)
	mockExChangeRate := models.ExChangeRate{
		Id:    1,
		Date:  "2020-09-17",
		From:  "USD",
		To:    "IDR",
		Rates: 1212112,
	}

	t.Run("success-from-usd", func(t *testing.T) {
		//mockExChangeRateRepo.On("GetByDate", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(&mockExChangeRate, nil).Once()
		u := usecase.NewCurrencyUsecase(mockExChangeRateRepo, timeoutContext)
		exchangeKey := mockExChangeRate.From + "_" + mockExChangeRate.To
		a, err := u.ExchangeFreeCurrconv(context.TODO(), exchangeKey)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockExChangeRateRepo.AssertExpectations(t)
	})
	t.Run("success-from-idr", func(t *testing.T) {
		//mockExChangeRateRepo.On("GetByDate", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(&mockExChangeRate, nil).Once()
		u := usecase.NewCurrencyUsecase(mockExChangeRateRepo, timeoutContext)
		exchangeKey := mockExChangeRate.To + "_" + mockExChangeRate.From
		a, err := u.ExchangeFreeCurrconv(context.TODO(), exchangeKey)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockExChangeRateRepo.AssertExpectations(t)
	})
	//t.Run("error-failed", func(t *testing.T) {
	//	//mockExChangeRateRepo.On("GetByDate", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(nil, errors.New("UnExpected")).Once()
	//	u := usecase.NewCurrencyUsecase(mockExChangeRateRepo, timeoutContext)
	//
	//	_, err := u.ExchangeRatesWithApi(context.TODO(), mockExChangeRate.From,mockExChangeRate.To)
	//
	//	assert.Error(t, err)
	//	//assert.Nil(t, a)
	//
	//	mockExChangeRateRepo.AssertExpectations(t)
	//})
}

func TestExchange(t *testing.T) {
	mockExChangeRateRepo := new(mocks.Repository)
	mockExChangeRate := models.ExChangeRate{
		Id:    1,
		Date:  "2020-09-17",
		From:  "USD",
		To:    "IDR",
		Rates: 1212112,
	}

	t.Run("success-from-usd", func(t *testing.T) {
		//mockExChangeRateRepo.On("GetByDate", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(&mockExChangeRate, nil).Once()
		u := usecase.NewCurrencyUsecase(mockExChangeRateRepo, timeoutContext)
		exchangeKey := mockExChangeRate.From + "_" + mockExChangeRate.To
		a, err := u.Exchange(context.TODO(),exchangeKey)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockExChangeRateRepo.AssertExpectations(t)
	})
	t.Run("success-from-idr", func(t *testing.T) {
		//mockExChangeRateRepo.On("GetByDate", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(&mockExChangeRate, nil).Once()
		u := usecase.NewCurrencyUsecase(mockExChangeRateRepo, timeoutContext)
		exchangeKey := mockExChangeRate.To + "_" + mockExChangeRate.From
		a, err := u.Exchange(context.TODO(), exchangeKey)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockExChangeRateRepo.AssertExpectations(t)
	})
	//t.Run("error-failed", func(t *testing.T) {
	//	//mockExChangeRateRepo.On("GetByDate", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(nil, errors.New("UnExpected")).Once()
	//	u := usecase.NewCurrencyUsecase(mockExChangeRateRepo, timeoutContext)
	//
	//	_, err := u.ExchangeRatesWithApi(context.TODO(), mockExChangeRate.From,mockExChangeRate.To)
	//
	//	assert.Error(t, err)
	//	//assert.Nil(t, a)
	//
	//	mockExChangeRateRepo.AssertExpectations(t)
	//})
}

func TestInsert(t *testing.T) {
	mockExChangeRateRepo := new(mocks.Repository)
	mockExChangeRateUsecase := new(mocks.Usecase)
	t.Run("success", func(t *testing.T) {
		mockExChangeRate := models.ExChangeRate{
			Id:    1,
			Date:  "2020-09-17",
			From:  "USD",
			To:    "IDR",
			Rates: 1212112,
		}
		mockExChangeRateUsecase.On("ExchangeRatesWithApi", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(mockExChangeRate, nil).Once()
		mockExChangeRateRepo.On("Insert", mock.Anything, mock.AnythingOfType("*models.ExChangeRate")).Return(nil).Once()
		mockExChangeRateUsecase.On("ExchangeRatesWithApi", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(mockExChangeRate, nil).Once()
		mockExChangeRateRepo.On("Insert", mock.Anything, mock.AnythingOfType("*models.ExChangeRate")).Return(nil).Once()

		u := usecase.NewCurrencyUsecase(mockExChangeRateRepo,timeoutContext)

		 err := u.Insert(context.TODO())

		assert.NoError(t, err)
		//assert.Equal(t, mockFacilities.FacilityName, tempMockFacilities.FacilityName)
		mockExChangeRateRepo.AssertExpectations(t)
	})
	//t.Run("error-Convert", func(t *testing.T) {
	//	mockExChangeRate := models.ExChangeRate{
	//		Id:    1,
	//		Date:  "2020-09-17",
	//		From:  "USD",
	//		To:    "IDR",
	//		Rates: 1212112,
	//	}
	//	mockExChangeRateUsecase.On("ExchangeRatesWithApi", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(nil,errors.New("Unexpected")).Once()
	//	mockExChangeRateRepo.On("Insert", mock.Anything, mock.AnythingOfType("*models.ExChangeRate")).Return(nil).Once()
	//	mockExChangeRateUsecase.On("ExchangeRatesWithApi", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(mockExChangeRate, nil).Once()
	//	mockExChangeRateRepo.On("Insert", mock.Anything, mock.AnythingOfType("*models.ExChangeRate")).Return(nil).Once()
	//
	//	u := usecase.NewCurrencyUsecase(mockExChangeRateRepo,timeoutContext)
	//
	//	err := u.Insert(context.TODO())
	//
	//	assert.Error(t, err)
	//	//assert.Equal(t, mockFacilities.FacilityName, tempMockFacilities.FacilityName)
	//	mockExChangeRateRepo.AssertExpectations(t)
	//})

}