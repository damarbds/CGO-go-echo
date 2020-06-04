package faq

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	GetByType(context context.Context, types int) ([]*models.FAQDto, error)
}
