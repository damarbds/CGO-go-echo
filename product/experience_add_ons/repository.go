package experience_add_ons

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	GetByExpId(ctx context.Context, exp_id string) ([]*models.ExperienceAddOn, error)
	GetById(ctx context.Context, id string) ([]*models.ExperienceAddOn, error)
	Insert(ctx context.Context, addOns models.ExperienceAddOn) (string, error)
	Update(ctx context.Context, addOns models.ExperienceAddOn) error
	Deletes(ctx context.Context, ids []string, expId string, deletedBy string) error
	DeleteByExpId(ctx context.Context, expId string, deletedBy string) error
}
