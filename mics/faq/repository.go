package faq

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetByType(context context.Context,types int)([]*models.FAQ,error)
}
