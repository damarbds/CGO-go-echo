package experience

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	GetByID(ctx context.Context, id string) (*models.ExperienceDto, error)
}
