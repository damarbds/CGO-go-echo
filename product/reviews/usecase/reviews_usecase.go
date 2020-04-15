package usecase

import (
	"encoding/json"
	"math"
	"time"

	"github.com/auth/user"
	"github.com/models"
	"github.com/product/reviews"
	"golang.org/x/net/context"
)

type reviewsUsecase struct {
	userRepo       user.Repository
	reviewsRepo    reviews.Repository
	contextTimeout time.Duration
}

func NewreviewsUsecase(a reviews.Repository, us user.Repository, timeout time.Duration) reviews.Usecase {
	return &reviewsUsecase{
		userRepo:       us,
		reviewsRepo:    a,
		contextTimeout: timeout,
	}
}

func (r reviewsUsecase) GetReviewsByExpIdWithPagination(
	ctx context.Context,
	page int,
	limit int,
	offset int,
	rating int,
	sortBy,
	exp_id string,
) (*models.ReviewsWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, r.contextTimeout)
	defer cancel()

	results, err := r.GetReviewsByExpId(ctx, exp_id, sortBy, rating, limit, offset)
	if err != nil {
		return nil, err
	}

	totalRecords, _ := r.reviewsRepo.CountRating(ctx, rating, exp_id)
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

	response := &models.ReviewsWithPagination{
		Data: results,
		Meta: meta,
	}

	return response, nil
}

func (r reviewsUsecase) GetReviewsByExpId(c context.Context, exp_id, sortBy string, rating, limit, offset int) ([]*models.ReviewDto, error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	res, err := r.reviewsRepo.GetByExpId(ctx, exp_id, sortBy, rating, limit, offset)
	if err != nil {
		return nil, err
	}
	var reviewDtos []*models.ReviewDto
	for _, element := range res {
		reviewDtoObject := models.ReviewDtoObject{}
		errObject := json.Unmarshal([]byte(element.Desc), &reviewDtoObject)
		if errObject != nil {
			return nil, models.ErrInternalServerError
		}
		var imageUrl string
		if reviewDtoObject.UserId != "" {
			getUser, _ := r.userRepo.GetByID(ctx, reviewDtoObject.UserId)
			if getUser != nil {
				imageUrl = getUser.ProfilePictUrl
			}
		}
		reviewDto := models.ReviewDto{
			Name:   reviewDtoObject.Name,
			Image:  imageUrl,
			Desc:   reviewDtoObject.Desc,
			Values: element.Values,
			Date:   element.CreatedDate,
		}
		reviewDtos = append(reviewDtos, &reviewDto)
	}

	return reviewDtos, nil
}
