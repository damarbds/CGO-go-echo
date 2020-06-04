package exclusion_service

import (
	"context"
	"github.com/models"
)

type Usecase interface {
	List(ctx context.Context,page, limit, offset int, token string)(*models.ExclusionServiceWithPagination,error)
}
