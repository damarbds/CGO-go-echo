package usecase

import (
	"context"
	"github.com/auth/merchant"
	"github.com/misc/notif"
	"github.com/models"
	"time"
)

type notifUsecase struct {
	notifRepo notif.Repository
	merchantUsecase merchant.Usecase
	contextTimeout time.Duration
}

func NewNotifUsecase(n notif.Repository, u merchant.Usecase, timeout time.Duration) notif.Usecase {
	return &notifUsecase{
		notifRepo:   n,
		merchantUsecase: u,
		contextTimeout: timeout,
	}
}


func (n notifUsecase) GetByMerchantID(ctx context.Context, token string) ([]*models.NotifDto, error) {
	ctx, cancel := context.WithTimeout(ctx, n.contextTimeout)
	defer cancel()

	currentMerchant, err := n.merchantUsecase.ValidateTokenMerchant(ctx, token)
	if err != nil {
		return nil, err
	}

	res, err := n.notifRepo.GetByMerchantID(ctx, currentMerchant.Id)
	if err != nil {
		return nil, err
	}


	notifs := make([]*models.NotifDto, len(res))
	for i, n := range res {
		notifType := "General"
		if n.Type == 2 {
			notifType = "Corporation"
		}
		notifs[i] = &models.NotifDto{
			Id:    n.Id,
			Type:  notifType,
			Title: n.Title,
			Desc:  n.Desc,
			Date:  n.CreatedDate,
		}
	}

	return notifs, nil
}

