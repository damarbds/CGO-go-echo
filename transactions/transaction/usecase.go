package transaction

import (
	"context"

	"github.com/models"
)

type Usecase interface {
	CountSuccess(ctx context.Context) (*models.Count, error)
}
