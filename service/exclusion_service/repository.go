package exclusion_service

import (
	"context"
	"github.com/models"
)

type Repository interface {
	Fetch(ctx context.Context,limit,offset int)([]*models.ExclusionService,error)
	GetCount(ctx context.Context)(int,error)
}
