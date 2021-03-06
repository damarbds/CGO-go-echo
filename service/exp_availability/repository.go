package exp_availability

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetCountDate(ctx context.Context,date string,merchantId string)(int,error)
	GetByExpIds(ctx context.Context, expId []*string,year int,month string) ([]*models.ExpAvailability, error)
	GetByExpId(ctx context.Context, expId string) ([]*models.ExpAvailability, error)
	Insert(ctx context.Context, availability models.ExpAvailability) (string, error)
	Update(ctx context.Context, availability models.ExpAvailability) error
	Deletes(ctx context.Context, ids []string, expId string, deletedBy string) error
	DeleteByExpId(ctx context.Context,expId string, deletedBy string) error
}
