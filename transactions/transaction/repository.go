package transaction

import (
	"context"
)

type Repository interface {
	CountSuccess(ctx context.Context) (int, error)
}
