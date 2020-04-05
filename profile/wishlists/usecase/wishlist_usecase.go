package usecase

import (
	"context"
	"github.com/auth/user"
	"github.com/models"
	"github.com/profile/wishlists"
	"time"
)

type wishListUsecase struct {
	wlRepo wishlists.Repository
	userUsercase user.Usecase
	ctxTimeout time.Duration
}

func NewWishlistUsecase(w wishlists.Repository, u user.Usecase, timeout time.Duration) wishlists.Usecase {
	return &wishListUsecase{
		wlRepo:     w,
		userUsercase: u,
		ctxTimeout: timeout,
	}
}

func (w wishListUsecase) Insert(ctx context.Context, wl *models.WishlistIn, token string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, w.ctxTimeout)
	defer cancel()

	currentUser, err := w.userUsercase.ValidateTokenUser(ctx, token)
	if err != nil {
		return "", err
	}

	newData := &models.Wishlist{
		Id:           "",
		CreatedBy:    currentUser.UserEmail,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		UserId:       currentUser.Id,
		ExpId: wl.ExpID,
		TransId: wl.TransID,
	}

	res, err := w.wlRepo.Insert(ctx, newData)
	if err != nil {
		return "", err
	}

	return res.Id, nil
}