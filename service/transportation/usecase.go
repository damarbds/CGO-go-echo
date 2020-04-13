package transportation

import (
	"context"

	"github.com/models"
)

type Usecase interface {
	CreateTransportation(ctx context.Context, newCommandTransportation models.NewCommandTransportation, token string) (*models.ResponseCreateExperience, error)
	UpdateTransportation(ctx context.Context, newCommandTransportation models.NewCommandTransportation, token string) (*models.ResponseCreateExperience, error)
	PublishTransportation(ctx context.Context, newCommandTransportation models.NewCommandTransportation, token string) (*models.ResponseCreateExperience, error)
	TimeOptions(ctx context.Context) ([]*models.TimeOptionDto, error)
	FilterSearchTrans(ctx context.Context, sortBy, harborSourceId, harborDestId, depDate, class string, isReturn bool, depTimeOptions, arrTimeOptions, guest, page, limit, offset int) (*models.FilterSearchTransWithPagination, error)
}
