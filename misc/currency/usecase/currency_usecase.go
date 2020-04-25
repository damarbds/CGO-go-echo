package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/misc/currency"
)

var (
	currencyApiKey = "cffe1d21f8a298007dd2"
	currencyUrl    = "https://free.currconv.com/api/v7/convert"
)

type currencyUsecase struct {
	contextTimeout time.Duration
}

func NewCurrencyUsecase(timeout time.Duration) currency.Usecase {
	return &currencyUsecase{
		contextTimeout: timeout,
	}
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
