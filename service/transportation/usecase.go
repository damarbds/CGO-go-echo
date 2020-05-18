package transportation

import (
	"context"

	"github.com/models"
)

type Usecase interface {
	GetDetail(ctx context.Context,id string)(*models.TransportationDto, error)
	UpdateStatus(ctx context.Context,status int,id string,token string)(*models.NewCommandChangeStatus,error)
	CreateTransportation(ctx context.Context, newCommandTransportation models.NewCommandTransportation, token string) (*models.ResponseCreateExperience, error)
	UpdateTransportation(ctx context.Context, newCommandTransportation models.NewCommandTransportation, token string) (*models.ResponseCreateExperience, error)
	PublishTransportation(ctx context.Context, newCommandTransportation models.NewCommandTransportation, token string) (*models.ResponseCreateExperience, error)
	TimeOptions(ctx context.Context) ([]*models.TimeOptionDto, error)
	FilterSearchTrans(ctx context.Context, isMerchant bool, token, search, status, sortBy, harborSourceId, harborDestId, depDate, class string, isReturn bool, depTimeOptions, arrTimeOptions, guest, page, limit, offset int,returnTransId string) (*models.FilterSearchTransWithPagination, error)
}
