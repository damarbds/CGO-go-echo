package temp_user_preferences

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetAll(ctx context.Context,page *int,size *int)([]*models.TempUserPreference,error)
}
