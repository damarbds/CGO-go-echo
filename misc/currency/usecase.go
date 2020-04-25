package currency

import "context"

type Usecase interface {
	Exchange(ctx context.Context, exchangeKey string) (map[string]interface{}, error)
}
