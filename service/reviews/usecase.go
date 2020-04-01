package reviews

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	GetReviewsByExpId(ctx context.Context,exp_id string) ([]*models.PromoDto, error)
}
