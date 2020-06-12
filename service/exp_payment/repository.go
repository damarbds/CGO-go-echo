package exp_payment

import (
	"context"

	"github.com/models"
)

type Repository interface {
	GetById(ctx context.Context,id string)([]*models.ExperiencePaymentJoinType,error)
	GetByExpID(ctx context.Context, expID string) ([]*models.ExperiencePaymentJoinType, error)
	Insert(ctx context.Context, payment models.ExperiencePayment) (string, error)
	Update(ctx context.Context, payment models.ExperiencePayment) error
	Deletes(ctx context.Context, ids []string, expId string, deletedBy string) error
	DeleteByExpId(ctx context.Context,expId string,deletedBy string)error
}
