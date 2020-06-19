package exp_facilities

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetJoin(ctx context.Context,expId string,transId string)([]*models.ExperienceFacilitiesJoin ,error)
	Insert(ctx context.Context,a *models.ExperienceFacilities)error
	Delete(ctx context.Context,expId string,transId string)error
}
