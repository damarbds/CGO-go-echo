package experience_add_ons

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	GetByExpId(ctx context.Context, exp_id string) ([]*models.ExperienceAddOnDto, error)
}
