package promo_experience_transport

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetById(ctx context.Context,promoId string)([]*models.PromoExperienceTransport,error)
	Insert(ctx context.Context,expId string, transId models.PromoExperienceTransport)error
	DeleteById(ctx context.Context,serviceId string,promoId string)error
}