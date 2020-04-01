package usecase

import (
	"encoding/json"
	"github.com/models"
	"github.com/product/reviews"
	"golang.org/x/net/context"
	"time"
)

type reviewsUsecase struct {
	reviewsRepo    reviews.Repository
	contextTimeout time.Duration
}
// NewharborsUsecase will create new an harborsUsecase object representation of harbors.Usecase interface
func NewreviewsUsecase(a  reviews.Repository, timeout time.Duration) reviews.Usecase {
	return &reviewsUsecase{
		reviewsRepo:    a,
		contextTimeout: timeout,
	}
}

func (r reviewsUsecase) GetReviewsByExpId(c context.Context, exp_id string) ([]*models.ReviewDto, error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	res, err := r.reviewsRepo.GetByExpId(ctx,exp_id)
	if err != nil {
		return nil, err
	}
	var reviewDtos []*models.ReviewDto
	for _, element := range res {

		reviewDtoObject := models.ReviewDtoObject{}
		errObject := json.Unmarshal([]byte(element.Desc), &reviewDtoObject)
		if errObject != nil {
			//fmt.Println("Error : ",err.Error())
			return nil,models.ErrInternalServerError
		}
		reviewDto := models.ReviewDto{
			Name:   reviewDtoObject.Name,
			Desc:   reviewDtoObject.Desc,
			Values: element.Values,
		}
		reviewDtos = append(reviewDtos, &reviewDto)
	}

	return reviewDtos, nil
}