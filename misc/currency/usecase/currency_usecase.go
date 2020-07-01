package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/models"
	"net/http"
	"time"

	"github.com/misc/currency"
)

var (
	currencyApiKey = "cffe1d21f8a298007dd2"
	currencyUrl    = "https://free.currconv.com/api/v7/convert"
	currencyUrlExchangeRatesApi = "https://api.exchangeratesapi.io/latest"
)

type currencyUsecase struct {
	contextTimeout time.Duration
	currencyRepo currency.Repository
}


func NewCurrencyUsecase(currencyRepo currency.Repository,timeout time.Duration) currency.Usecase {
	return &currencyUsecase{
		currencyRepo:currencyRepo,
		contextTimeout: timeout,
	}
}

func (c currencyUsecase) Insert(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	getCurrencyNowFromIDR ,err := c.ExchangeRatesWithApi(ctx,"IDR","USD")
	if err != nil {
		return err
	}
	idr := models.ExChangeRate{
		Id:    0,
		Date:  getCurrencyNowFromIDR.Date,
		From:  getCurrencyNowFromIDR.Base,
		To:    "USD",
		Rates: getCurrencyNowFromIDR.Rates.USD,
	}
	_ = c.currencyRepo.Insert(ctx, &idr)

	getCurrencyNowFromUSD ,err := c.ExchangeRatesWithApi(ctx,"USD","IDR")
	if err != nil {
		return err
	}
	usd := models.ExChangeRate{
		Id:    0,
		Date:  getCurrencyNowFromUSD.Date,
		From:  getCurrencyNowFromUSD.Base,
		To:    "IDR",
		Rates: getCurrencyNowFromUSD.Rates.IDR,
	}
	_ = c.currencyRepo.Insert(ctx, &usd)
	return nil
}
func (c currencyUsecase) Exchange(ctx context.Context, exchangeKey string) (map[string]interface{}, error) {
	client := &http.Client{}

	var data map[string]interface{}

	req, err := http.NewRequest(http.MethodGet, currencyUrl, nil)
	if err != nil {
		return data, err
	}

	q := req.URL.Query()
	q.Add("q", exchangeKey)
	q.Add("compact", "ultra")
	q.Add("apiKey", currencyApiKey)
	req.URL.RawQuery = q.Encode()

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err != nil {
			return data, err
		}
		return data, nil
	} else {
		err := errors.New("Currency Exchange Error : " + resp.Status)
		return data, err
	}
}


func (c currencyUsecase) ExchangeRatesApi(ctx context.Context, base string , symbols string) (models.CurrencyExChangeRate, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	var data models.CurrencyExChangeRate
	getExChangeRates ,err:=c.currencyRepo.GetByDate(ctx,base,symbols)
	if err != nil {
		return data, err
	}
	if symbols == "IDR" {
		data.Rates.IDR = getExChangeRates.Rates
		data.Date = getExChangeRates.Date
		data.Base = getExChangeRates.From
	}else if symbols == "USD" {
		data.Rates.USD = getExChangeRates.Rates
		data.Date = getExChangeRates.Date
		data.Base = getExChangeRates.From
	}
	return data,nil
}
func (c currencyUsecase) ExchangeRatesWithApi(ctx context.Context, base string , symbols string) (models.CurrencyExChangeRate, error) {
	client := &http.Client{}

	var data models.CurrencyExChangeRate

	req, err := http.NewRequest(http.MethodGet, currencyUrlExchangeRatesApi, nil)
	if err != nil {
		return data, err
	}

	q := req.URL.Query()
	q.Add("base", base)
	q.Add("symbols", symbols)
	req.URL.RawQuery = q.Encode()

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err != nil {
			return data, err
		}
		return data, nil
	} else {
		err := errors.New("Currency Exchange Error : " + resp.Status)
		return data, err
	}
}
func (c currencyUsecase) ExchangeFreeCurrconv(ctx context.Context, exchangeKey string) (map[string]interface{}, error) {
	client := &http.Client{}

	var data map[string]interface{}

	req, err := http.NewRequest(http.MethodGet, currencyUrl, nil)
	if err != nil {
		return data, err
	}

	q := req.URL.Query()
	q.Add("q", exchangeKey)
	q.Add("compact", "ultra")
	q.Add("apiKey", currencyApiKey)
	req.URL.RawQuery = q.Encode()

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err != nil {
			return data, err
		}
		return data, nil
	} else {
		err := errors.New("Currency Exchange Error : " + resp.Status)
		return data, err
	}
}
