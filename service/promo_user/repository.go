package promo_user

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetByUserId(ctx context.Context,userId string,promoId string)([]*models.PromoUser,error)
	Insert(ctx context.Context,pm models.PromoUser)error
	DeleteByUserId(ctx context.Context,userId string,promoId string)error
}