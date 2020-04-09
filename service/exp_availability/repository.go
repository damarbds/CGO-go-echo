package exp_availability

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetByExpId(ctx context.Context,expId string)([]*models.ExpAvailability ,error)
	Insert(ctx context.Context,availability models.ExpAvailability) (string,error)
	Update(ctx context.Context,availability models.ExpAvailability) error
	Deletes(ctx context.Context, ids []string,expId string,deletedBy string) error
}