package usecase

import (
	"encoding/json"
	guuid "github.com/google/uuid"
	"github.com/service/experience"
	"math"
	"time"

	"github.com/auth/user"
	"github.com/models"
	"github.com/product/reviews"
	"golang.org/x/net/context"
)

type reviewsUsecase struct {
	experienceRepo experience.Repository
	userUsecase 	user.Usecase
	userRepo       user.Repository
	reviewsRepo    reviews.Repository
	contextTimeout time.Duration
}

func NewreviewsUsecase(exp experience.Repository,usu user.Usecase,a reviews.Repository, us user.Repository, timeout time.Duration) reviews.Usecase {
	return &reviewsUsecase{
		experienceRepo:exp,
		userUsecase:usu,
		userRepo:       us,
		reviewsRepo:    a,
		contextTimeout: timeout,
	}
}

func (r reviewsUsecase) CreateReviews(c context.Context, command models.NewReviewCommand,token string) (*models.NewReviewCommand, error) {
	ctx, cancel := context.WithTimeout(c, r.contextTimeout)
	defer cancel()

	validateToken ,err := r.userUsecase.ValidateTokenUser(ctx,token)
	if err != nil {
		return nil,models.ErrUnAuthorize
	}
	 var name string
	if validateToken.FullName != ""{
		name = validateToken.FullName
	}else {
		name = validateToken.UserEmail
	}
	 desc :=  models.ReviewDtoObject{
		 Name:   name,
		 UserId: validateToken.Id ,
		 Desc:   command.Desc,
	 }
	 descJson , _ := json.Marshal(desc)
	 currentValue := (command.GuideReview + command.ActivitiesReview + command.ServiceReview + command.CleanlinessReview +
	 	command.ValueReview)/5

	 review := models.Review{
		 Id:                guuid.New().String(),
		 CreatedBy:         validateToken.UserEmail,
		 CreatedDate:       time.Now(),
		 ModifiedBy:        nil,
		 ModifiedDate:      nil,
		 DeletedBy:         nil,
		 DeletedDate:       nil,
		 IsDeleted:         0,
		 IsActive:          0,
		 Values:            currentValue,
		 Desc:              string(descJson),
		 ExpId:             command.ExpId,
		 UserId:            &validateToken.Id,
		 GuideReview:       &command.GuideReview,
		 ActivitiesReview:  &command.ActivitiesReview,
		 ServiceReview:     &command.ServiceReview,
		 CleanlinessReview: &command.CleanlinessReview,
		 ValueReview:       &command.ValueReview,
	 }
	 insert,err := r.reviewsRepo.Insert(ctx,review)
	 if err != nil {
	 	return  nil,err
	 }
	 command.Id = insert
	 countRating , err := r.reviewsRepo.CountRating(ctx,0,command.ExpId)

	 getExperience, err := r.experienceRepo.GetByID(ctx,command.ExpId)
	if getExperience.GuideReview == nil {
		var initial float64 = 0
		getExperience.GuideReview = &initial
		getExperience.ActivitiesReview = &initial
		getExperience.ServiceReview = &initial
		getExperience.CleanlinessReview = &initial
		getExperience.ValueReview = &initial
	}
	rating := (getExperience.Rating + currentValue)/float64(countRating)
	guideReview := (*getExperience.GuideReview + command.GuideReview)/float64(countRating)
	activitiesReview := (*getExperience.ActivitiesReview + command.ActivitiesReview)/float64(countRating)
	serviceReview := (*getExperience.ServiceReview + command.ServiceReview)/float64(countRating)
	cleanlinessReview := (*getExperience.CleanlinessReview + command.CleanlinessReview)/float64(countRating)
	valueReview := (*getExperience.ValueReview + command.ValueReview)/float64(countRating)

	experience := models.Experience{
		Id:                      command.ExpId,
		ModifiedBy:              &validateToken.UserEmail,
		Rating:                  rating,
		GuideReview:             &guideReview,
		ActivitiesReview:        &activitiesReview,
		ServiceReview:           &serviceReview,
		CleanlinessReview:       &cleanlinessReview,
		ValueReview:             &valueReview,
	}

	err = r.experienceRepo.UpdateRating(ctx,experience)
	if err != nil {
		return nil,err
	}

	command.Id = insert

	return &command,nil
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

	res, err := r.reviewsRepo.GetByExpId(ctx, exp_id, sortBy, rating, limit, offset,"")
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
			UserId   : element.UserId,
			GuideReview : element.GuideReview,
			ActivitiesReview : element.ActivitiesReview,
			ServiceReview : element.ServiceReview,
			CleanlinessReview : element.CleanlinessReview,
			ValueReview : element.ValueReview,
		}
		reviewDtos = append(reviewDtos, &reviewDto)
	}

	return reviewDtos, nil
}
