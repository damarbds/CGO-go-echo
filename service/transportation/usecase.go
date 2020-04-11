package transportation

import (
	"context"

	"github.com/models"
)

type Usecase interface {
	CreateTransportation(ctx context.Context, newCommandTransportation models.NewCommandTransportation, token string) (*models.ResponseCreateExperience, error)
	UpdateTransportation(ctx context.Context, newCommandTransportation models.NewCommandTransportation, token string) (*models.ResponseCreateExperience, error)
	PublishTransportation(ctx context.Context, newCommandTransportation models.NewCommandTransportation, token string) (*models.ResponseCreateExperience, error)
	List(ctx context.Context) ([]*models.TimeOptionDto, error)
}
