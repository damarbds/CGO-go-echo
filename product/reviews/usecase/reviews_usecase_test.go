package usecase_test

import (
	"context"
	"errors"
	//_adminUsecaseMock "github.com/auth/admin/mocks"
	_userMock "github.com/auth/user/mocks"
	"github.com/models"
	"github.com/product/reviews/mocks"
	"github.com/product/reviews/usecase"
	//_adminUsecaseMock "github.com/auth/admin/mocks"
	_experienceUsecaseMock "github.com/service/experience/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)
var(
	timeoutContext = time.Second*30
	rating = 12
	userId = "adasdasdqweqwe"
	reviewValue float64 = 12
	mockReview = []*models.Review{
		&models.Review{
			Id:                "asdasdasdsadqeq",
			CreatedBy:         "Test 1",
			CreatedDate:       time.Now(),
			ModifiedBy:        nil,
			ModifiedDate:      nil,
			DeletedBy:         nil,
			DeletedDate:       nil,
			IsDeleted:         0,
			IsActive:          1,
			Values:            12,
			Desc:               `{"name":"radyaupdate12345","userid":"radyaupdate12345","desc":"good exp"}`,
			ExpId:             "qewqweqwe",
			UserId:            &userId,
			GuideReview:       &reviewValue,
			ActivitiesReview:  &reviewValue,
			ServiceReview:     &reviewValue,
			CleanlinessReview: &reviewValue,
			ValueReview:       &reviewValue,
		},
		&models.Review{
			Id:                "jlkjlkjlkjkl",
			CreatedBy:         "Test 1",
			CreatedDate:       time.Now(),
			ModifiedBy:        nil,
			ModifiedDate:      nil,
			DeletedBy:         nil,
			DeletedDate:       nil,
			IsDeleted:         0,
			IsActive:          1,
			Values:            12,
			Desc:               `{"name":"radyaupdate12345","userid":"radyaupdate12345","desc":"good exp"}`,
			ExpId:             "qewqweqwe",
			UserId:            &userId,
			GuideReview:       &reviewValue,
			ActivitiesReview:  &reviewValue,
			ServiceReview:     &reviewValue,
			CleanlinessReview: &reviewValue,
			ValueReview:       &reviewValue,
		},
	}
	mockReviewDto = []*models.ReviewDto{
		&models.ReviewDto{
			Name:              "Test",
			Image:             "adasdasd",
			Desc:              `{"name":"radyaupdate12345","userid":"radyaupdate12345","desc":"good exp"}`,
			Values:            12121,
			Date:              time.Now(),
			UserId:            &userId,
			GuideReview:       &reviewValue,
			ActivitiesReview:  &reviewValue,
			ServiceReview:     &reviewValue,
			CleanlinessReview: &reviewValue,
			ValueReview:       &reviewValue,
		},
	}
)
func TestGetReviewsByExpId(t *testing.T) {
	mockReviewRepo := new(mocks.Repository)
	mockUserRepo := new(_userMock.Repository)
	mockUserUsecase := new(_userMock.Usecase)
	mockExperienceRepo := new(_experienceUsecaseMock.Repository)

	user := models.User{
		Id:                   "",
		CreatedBy:            "",
		CreatedDate:          time.Time{},
		ModifiedBy:           nil,
		ModifiedDate:         nil,
		DeletedBy:            nil,
		DeletedDate:          nil,
		IsDeleted:            0,
		IsActive:             0,
		UserEmail:            "",
		FullName:             "",
		PhoneNumber:          "",
		VerificationSendDate: time.Time{},
		VerificationCode:     "",
		ProfilePictUrl:       "",
		Address:              "",
		Dob:                  time.Time{},
		Gender:               0,
		IdType:               0,
		IdNumber:             "",
		ReferralCode:         "",
		Points:               0,
	}
	//mockReview := mockReview
	var mockListReview []*models.Review
	mockListReview = mockReview

	//page := 1
	limit := 1
	offset := 0
	//count := 1
	t.Run("success", func(t *testing.T) {
		mockReviewRepo.On("GetByExpId",
			mock.Anything,mock.AnythingOfType("string"),mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(mockListReview, nil).Once()
		//mockReviewRepo.On("GetCount", mock.Anything).Return(count, nil).Once()
		mockUserRepo.On("GetByID", mock.Anything,mock.AnythingOfType("string")).Return(&user, nil).Once()
		mockUserRepo.On("GetByID", mock.Anything,mock.AnythingOfType("string")).Return(&user, nil).Once()
		u := usecase.NewreviewsUsecase(mockExperienceRepo,mockUserUsecase, mockReviewRepo,mockUserRepo, timeoutContext)

		_, err := u.GetReviewsByExpId(context.TODO(),mockReview[0].ExpId,"oldestdate",2,limit,offset)

		assert.NoError(t, err)
		//assert.Len(t, list, len(mockListReview))

		mockReviewRepo.AssertExpectations(t)
		mockReviewRepo.AssertExpectations(t)
	})

}

func TestGetReviewsByExpIdWithPagination(t *testing.T) {
	mockReviewRepo := new(mocks.Repository)
	mockReviewUsecase := new(mocks.Usecase)
	mockUserRepo := new(_userMock.Repository)
	mockUserUsecase := new(_userMock.Usecase)
	mockExperienceRepo := new(_experienceUsecaseMock.Repository)
	user := models.User{
		Id:                   "",
		CreatedBy:            "",
		CreatedDate:          time.Time{},
		ModifiedBy:           nil,
		ModifiedDate:         nil,
		DeletedBy:            nil,
		DeletedDate:          nil,
		IsDeleted:            0,
		IsActive:             0,
		UserEmail:            "",
		FullName:             "",
		PhoneNumber:          "",
		VerificationSendDate: time.Time{},
		VerificationCode:     "",
		ProfilePictUrl:       "",
		Address:              "",
		Dob:                  time.Time{},
		Gender:               0,
		IdType:               0,
		IdNumber:             "",
		ReferralCode:         "",
		Points:               0,
	}
	//mockReview := mockReview
	var mockListReview []*models.Review
	mockListReview = mockReview

	mockReviewPagination := &models.ReviewsWithPagination{
		Data: mockReviewDto,
		Meta: &models.MetaPagination{
			Page:          1,
			Total:         1,
			TotalRecords:  1,
			Prev:          1,
			Next:          2,
			RecordPerPage: 1,
		},
	}
	t.Run("success", func(t *testing.T) {
		mockReviewUsecase.On("GetReviewsByExpId",mock.Anything,
			mock.AnythingOfType("string"),mock.AnythingOfType("string"),mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),mock.AnythingOfType("int")).Return(mockListReview, nil).Once()

		mockReviewRepo.On("GetByExpId",
			mock.Anything,mock.AnythingOfType("string"),mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(mockListReview, nil).Once()
		//mockReviewRepo.On("GetCount", mock.Anything).Return(count, nil).Once()
		mockUserRepo.On("GetByID", mock.Anything,mock.AnythingOfType("string")).Return(&user, nil).Once()
		mockUserRepo.On("GetByID", mock.Anything,mock.AnythingOfType("string")).Return(&user, nil).Once()

		mockReviewRepo.On("CountRating",mock.Anything,
			mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(2, nil).Once()


		u := usecase.NewreviewsUsecase(mockExperienceRepo,mockUserUsecase, mockReviewRepo,mockUserRepo, timeoutContext)

		a, err := u.GetReviewsByExpIdWithPagination(context.TODO(), mockReviewPagination.Meta.Page,2,
			0,2,"asdasd",mockReview[0].ExpId)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockReviewRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockReviewUsecase.On("GetReviewsByExpId",mock.Anything,
			mock.AnythingOfType("string"),mock.AnythingOfType("string"),mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),mock.AnythingOfType("int")).Return(mockListReview, nil).Once()

		mockReviewRepo.On("GetByExpId",
			mock.Anything,mock.AnythingOfType("string"),mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(nil, errors.New("UnExpected")).Once()
		//mockReviewRepo.On("GetCount", mock.Anything).Return(count, nil).Once()
		mockUserRepo.On("GetByID", mock.Anything,mock.AnythingOfType("string")).Return(&user, nil).Once()
		mockUserRepo.On("GetByID", mock.Anything,mock.AnythingOfType("string")).Return(&user, nil).Once()

		mockReviewRepo.On("CountRating",mock.Anything,
			mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(2, nil).Once()


		u := usecase.NewreviewsUsecase(mockExperienceRepo,mockUserUsecase, mockReviewRepo,mockUserRepo, timeoutContext)

		a, err := u.GetReviewsByExpIdWithPagination(context.TODO(), mockReviewPagination.Meta.Page,2,
			0,2,"asdasd",mockReview[0].ExpId)


		assert.Error(t, err)
		assert.Nil(t, a)

		//mockReviewRepo.AssertExpectations(t)
	})

}

func TestCreateReviews(t *testing.T) {
	mockReviewRepo := new(mocks.Repository)
	//mockReviewUsecase := new(mocks.Usecase)
	mockUserRepo := new(_userMock.Repository)
	mockUserUsecase := new(_userMock.Usecase)
	mockExperienceRepo := new(_experienceUsecaseMock.Repository)

	mockAdmin := &models.UserInfoDto{
		Id:             "",
		CreatedDate:    time.Time{},
		UpdatedDate:    nil,
		IsActive:       0,
		UserEmail:      "",
		FullName:       "",
		PhoneNumber:    "",
		ProfilePictUrl: "",
		Address:        "",
		Dob:            time.Time{},
		Gender:         0,
		IdType:         0,
		IdNumber:       "",
		ReferralCode:   "",
		Points:         0,
	}
	experience := &models.ExperienceJoinForegnKey{
		Id:                       "",
		CreatedBy:                "",
		CreatedDate:              time.Time{},
		ModifiedBy:               nil,
		ModifiedDate:             nil,
		DeletedBy:                nil,
		DeletedDate:              nil,
		IsDeleted:                0,
		IsActive:                 0,
		ExpTitle:                 "",
		ExpType:                  "",
		ExpTripType:              "",
		ExpBookingType:           "",
		ExpDesc:                  "",
		ExpMaxGuest:              0,
		ExpPickupPlace:           "",
		ExpPickupTime:            "",
		ExpPickupPlaceLongitude:  0,
		ExpPickupPlaceLatitude:   0,
		ExpPickupPlaceMapsName:   "",
		ExpInternary:             "",
		ExpFacilities:            "",
		ExpInclusion:             "",
		ExpRules:                 "",
		Status:                   0,
		Rating:                   0,
		ExpLocationLatitude:      0,
		ExpLocationLongitude:     0,
		ExpLocationName:          "",
		ExpCoverPhoto:            nil,
		ExpDuration:              0,
		MinimumBookingId:         "",
		MerchantId:               "",
		HarborsId:                "",
		GuideReview:              nil,
		ActivitiesReview:         nil,
		ServiceReview:            nil,
		CleanlinessReview:        nil,
		ValueReview:              nil,
		ExpPaymentDeadlineAmount: nil,
		ExpPaymentDeadlineType:   nil,
		IsCustomisedByUser:       nil,
		ExpLocationMapName:       nil,
		ExpLatitudeMap:           nil,
		ExpLongitudeMap:          nil,
		ExpMaximumBookingAmount:  nil,
		ExpMaximumBookingType:    nil,
		MinimumBookingAmount:     nil,
		MinimumBookingDesc:       "",
	}
	t.Run("success", func(t *testing.T) {
		tempMockReview := models.NewReviewCommand{
			Id:                mockReview[0].Id,
			ExpId:             mockReview[0].ExpId,
			Desc:              mockReview[0].Desc,
			GuideReview:       reviewValue,
			ActivitiesReview:  reviewValue,
			ServiceReview:     reviewValue,
			CleanlinessReview: reviewValue,
			ValueReview:       reviewValue,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockReview.Id = ""
		mockReviewRepo.On("Insert", mock.Anything, mock.AnythingOfType("models.Review")).Return(mockReview[0].Id,nil).Once()
		mockUserUsecase.On("ValidateTokenUser", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()
		mockReviewRepo.On("CountRating",mock.Anything,
			mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(2, nil).Once()
		mockExperienceRepo.On("GetByID",mock.Anything,mock.AnythingOfType("string")).Return(experience, nil).Once()
		mockExperienceRepo.On("UpdateRating",mock.Anything,mock.AnythingOfType("models.Experience")).Return(nil).Once()

		u := usecase.NewreviewsUsecase(mockExperienceRepo,mockUserUsecase, mockReviewRepo,mockUserRepo, timeoutContext)

		_,err := u.CreateReviews(context.TODO(), tempMockReview,token)

		assert.NoError(t, err)
		//assert.Equal(t, mockReview[]0, tempMockReview.ReviewName)
		mockReviewRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {
		tempMockReview := models.NewReviewCommand{
			Id:                mockReview[0].Id,
			ExpId:             mockReview[0].ExpId,
			Desc:              mockReview[0].Desc,
			GuideReview:       reviewValue,
			ActivitiesReview:  reviewValue,
			ServiceReview:     reviewValue,
			CleanlinessReview: reviewValue,
			ValueReview:       reviewValue,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockReview.Id = ""
		mockReviewRepo.On("Insert", mock.Anything, mock.AnythingOfType("models.Review")).Return("",errors.New("UnExpected")).Once()
		mockUserUsecase.On("ValidateTokenUser", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()
		mockReviewRepo.On("CountRating",mock.Anything,
			mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(2, nil).Once()
		mockExperienceRepo.On("GetByID",mock.Anything,mock.AnythingOfType("string")).Return(experience, nil).Once()
		mockExperienceRepo.On("UpdateRating",mock.Anything,mock.AnythingOfType("models.Experience")).Return(nil).Once()

		u := usecase.NewreviewsUsecase(mockExperienceRepo,mockUserUsecase, mockReviewRepo,mockUserRepo, timeoutContext)

		_,err := u.CreateReviews(context.TODO(), tempMockReview,token)

		assert.Error(t, err)
		//assert.Equal(t, mockReview[]0, tempMockReview.ReviewName)
		//mockReviewRepo.AssertExpectations(t)
	})

}

