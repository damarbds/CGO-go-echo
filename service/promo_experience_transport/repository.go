package promo_experience_transport

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetByExperienceTransportId(ctx context.Context,expId string,transportId string,promoId string)([]*models.PromoExperienceTransport,error)
	Insert(ctx context.Context,pet models.PromoExperienceTransport)error
	DeleteById(ctx context.Context,serviceId string,promoId string)error
}