package facilities

import (
	"context"

	"github.com/models"
)

type Usecase interface {
	List(ctx context.Context) ([]*models.FacilityDto, error)
	GetAll(ctx context.Context ,page,limit,offset int)(*models.FacilityDtoWithPagination,error)
	GetById(ctx context.Context,id int)(*models.FacilityDto,error)
	Create(ctx context.Context, f *models.NewCommandFacilities,token string)(*models.ResponseDelete,error)
	Update(ctx context.Context, f *models.NewCommandFacilities,token string)(*models.ResponseDelete,error)
	Delete(ctx context.Context, id int,token string)(*models.ResponseDelete,error)
}
