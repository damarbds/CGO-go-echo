package usecase

import (
	"context"
	"encoding/json"
	"github.com/service/exp_photos"
	"math"
	"time"

	"github.com/auth/user"
	"github.com/models"
	"github.com/product/reviews"
	"github.com/profile/wishlists"
	"github.com/service/exp_payment"
	"github.com/service/experience"
)

type wishListUsecase struct {
	expPhotos   exp_photos.Repository
	wlRepo      wishlists.Repository
	userUsecase user.Usecase
	expRepo     experience.Repository
	paymentRepo exp_payment.Repository
	reviewRepo  reviews.Repository
	ctxTimeout  time.Duration
}

func NewWishlistUsecase(
	expPhotos exp_photos.Repository,
	w wishlists.Repository,
	u user.Usecase,
	e experience.Repository,
	p exp_payment.Repository,
	r reviews.Repository,
	timeout time.Duration,
) wishlists.Usecase {
	return &wishListUsecase{
		expPhotos:expPhotos,
		wlRepo:      w,
		userUsecase: u,
		expRepo:     e,
		paymentRepo: p,
		reviewRepo:  r,
		ctxTimeout:  timeout,
	}
}

func (w wishListUsecase) List(ctx context.Context, token string,page int, limit int, offset int) (*models.WishlistOutWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, w.ctxTimeout)
	defer cancel()

	currentUser, err := w.userUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return nil, err
	}

	wLists, err := w.wlRepo.List(ctx, currentUser.Id,limit,offset)
	if err != nil {
		return nil, err
	}

	results := make([]*models.WishlistOut, len(wLists))
	for i, wl := range wLists {

		exp, err := w.expRepo.GetByID(ctx, wl.ExpId.String)
		if err != nil {
			return nil, err
		}

		var expType []string
		if errUnmarshal := json.Unmarshal([]byte(exp.ExpType), &expType); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}

		expPayment, err := w.paymentRepo.GetByExpID(ctx, exp.Id)
		if err != nil {
			return nil, err
		}

		var currency string
		if expPayment[0].Currency == 1 {
			currency = "USD"
		} else {
			currency = "IDR"
		}

		var priceItemType string
		if expPayment[0].PriceItemType == 1 {
			priceItemType = "Per Pax"
		} else {
			priceItemType = "Per Trip"
		}

		countRating, err := w.reviewRepo.CountRating(ctx, 0, exp.Id)
		if err != nil {
			return nil, err
		}

			wtype := "EXPERIENCE"
		if wl.TransId.String != "" {
			wtype = "TRANSPORTATION"
		}
		listPhotos := make([]models.ExpPhotosObj,0)
		expPhotoQuery, errorQuery := w.expPhotos.GetByExperienceID(ctx, exp.Id)
		if errorQuery != nil {
			return nil, errorQuery
		}
		if expPhotoQuery != nil {
			for _, element := range expPhotoQuery {
				expPhoto := models.ExpPhotosObj{
					Folder:        element.ExpPhotoFolder,
					ExpPhotoImage: nil,
				}
				var expPhotoImage []models.CoverPhotosObj
				errObject := json.Unmarshal([]byte(element.ExpPhotoImage), &expPhotoImage)
				if errObject != nil {
					return nil, models.ErrInternalServerError
				}
				expPhoto.ExpPhotoImage = expPhotoImage
				listPhotos = append(listPhotos, expPhoto)
			}
		}
		results[i] = &models.WishlistOut{
			WishlistID:  wl.Id,
			Type:        wtype,
			ExpID:       exp.Id,
			ExpTitle:    exp.ExpTitle,
			ExpType:     expType,
			Rating:      exp.Rating,
			CountRating: countRating,
			Currency:    currency,
			Price:       expPayment[0].Price,
			PaymentType: priceItemType,
			CoverPhoto:  *exp.ExpCoverPhoto,
			ListPhoto:listPhotos,
		}
	}
	totalRecords, _ := w.wlRepo.Count(ctx,currentUser.Id)
	totalPage := int(math.Ceil(float64(totalRecords) / float64(limit)))
	prev := page
	next := page
	if page != 1 {
		prev = page - 1
	}

	if page != totalPage {
		next = page + 1
	}

	meta := &models.MetaPagination{
		Page:          page,
		Total:         totalPage,
		TotalRecords:  totalRecords,
		Prev:          prev,
		Next:          next,
		RecordPerPage: len(results),
	}

	response := &models.WishlistOutWithPagination{
		Data: results,
		Meta: meta,
	}
	return response, nil
}

func (w wishListUsecase) Insert(ctx context.Context, wl *models.WishlistIn, token string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, w.ctxTimeout)
	defer cancel()

	currentUser, err := w.userUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return "", err
	}
	checkWhislist, err := w.wlRepo.GetByUserAndExpId(ctx, currentUser.Id, wl.ExpID,wl.TransID)
	if err != nil {
		return "", err
	}
	if len(checkWhislist) != 0{
		errDelete := w.wlRepo.DeleteByUserIdAndExpIdORTransId(ctx, currentUser.Id, wl.ExpID,wl.TransID,currentUser.UserEmail)
		if errDelete != nil {
			return "", errDelete
		}
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
		ExpId:        wl.ExpID,
		TransId:      wl.TransID,
	}

	res, err := w.wlRepo.Insert(ctx, newData)
	if err != nil {
		return "", err
	}

	return res.Id, nil
}
