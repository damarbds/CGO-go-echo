package currency

import (
	"context"
	"github.com/models"
)

type Usecase interface {
	Exchange(ctx context.Context, exchangeKey string) (map[string]interface{}, error)
	ExchangeRatesApi(ctx context.Context, base string , symbols string) (models.CurrencyExChangeRate, error)
	ExchangeRatesWithApi(ctx context.Context, base string , symbols string) (models.CurrencyExChangeRate, error)
	ExchangeFreeCurrconv(ctx context.Context, exchangeKey string) (map[string]interface{}, error)
}
