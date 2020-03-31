package exp_photos

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	GetByExperienceID(ctx context.Context, id string) ([]models.ExpPhotosDto, error)
}
