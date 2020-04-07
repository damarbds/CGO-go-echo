package facilities

import (
	"context"

	"github.com/models"
)

type Usecase interface {
	List(ctx context.Context) ([]*models.FacilityDto, error)
}
