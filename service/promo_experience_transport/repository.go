package promo_experience_transport

import (
	"context"
	"github.com/models"
)

type Repository interface {
	CountByPromoId(ctx context.Context,promoId string)(int,error)
	GetByExperienceTransportId(ctx context.Context,expId string,transportId string,promoId string)([]*models.PromoExperienceTransport,error)
	Insert(ctx context.Context,pet models.PromoExperienceTransport)error
	DeleteById(ctx context.Context,serviceId string,promoId string)error
}