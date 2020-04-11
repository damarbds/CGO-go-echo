package transportation

import (
	"context"
	"github.com/models"
)

type Usecase interface {
	List(ctx context.Context) ([]*models.TimeOptionDto, error)
}
