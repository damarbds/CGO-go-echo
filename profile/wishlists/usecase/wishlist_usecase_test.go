package usecase_test

import (
	"context"
	"database/sql"
	"errors"
	//"errors"
	//_adminUsecaseMock "github.com/auth/admin/mocks"
	_userMock "github.com/auth/user/mocks"
	"github.com/models"
	_reviewsMock "github.com/product/reviews/mocks"
	"github.com/profile/wishlists/mocks"
	"github.com/profile/wishlists/usecase"
	_expPaymentMock "github.com/service/exp_payment/mocks"
	_expPhotosMock "github.com/service/exp_photos/mocks"
	_experienceMock "github.com/service/experience/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)
var(
	expId = "asdasdsad"
	transId = "adasd"
	timeoutContext = time.Second*30
	mockWishlist = []models.Wishlist{
		models.Wishlist{
			Id:           "asdasdasdsad",
			CreatedBy:    "Test 1",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			TransId:      "asdsdas",
			ExpId:        "asdasd",
			UserId:       "zxczxcxzc",
		},
		models.Wishlist{
			Id:           "asdasdasdsadfddfdf",
			CreatedBy:    "Test 2",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			TransId:      "asdsdas",
			ExpId:        "asdasd",
			UserId:       "zxczxcxzc",
		},
	}
	mockWishlistObj = []models.WishlistObj{
		models.WishlistObj{
			Id:           "asdasdasdsad",
			CreatedBy:    "Test 1",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			TransId:     sql.NullString{String:transId},
			ExpId:       sql.NullString{String:expId},
			UserId:       "zxczxcxzc",
		},
	}
	mockExpPhotosRepo = new(_expPhotosMock.Repository)
	mockExperienceRepo = new(_experienceMock.Repository)
	mockExpPaymentRepo = new(_expPaymentMock.Repository)
	mockUserUsecase = new(_userMock.Usecase)
	mockReviewsRepo = new(_reviewsMock.Repository)
	userInfoDto = &models.UserInfoDto{
		Id:             "asdasdsa",
		CreatedDate:    time.Time{},
		UpdatedDate:    nil,
		IsActive:       0,
		UserEmail:      "test@gmai.com",
		FullName:       "test1231",
		PhoneNumber:    "1231",
		ProfilePictUrl: "",
		Address:        "",
		Dob:            time.Time{},
		Gender:         0,
		IdType:         0,
		IdNumber:       "",
		ReferralCode:   "",
		Points:         0,
	}
	token = "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	coverPhoto = "https://cgostorage.blob.core.windows.net/cgo-storage/folder%2011111/738864447603174415.jpg"
	experience = &models.ExperienceJoinForegnKey{
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
		ExpType:                  `["Diving","Fishing","Snorkeling"]`,
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
		ExpCoverPhoto:            &coverPhoto,
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
	experiencePaymentType = []*models.ExperiencePaymentJoinType{
		&models.ExperiencePaymentJoinType{
			Id:                 "",
			CreatedBy:          "",
			CreatedDate:        time.Time{},
			ModifiedBy:         nil,
			ModifiedDate:       nil,
			DeletedBy:          nil,
			DeletedDate:        nil,
			IsDeleted:          0,
			IsActive:           0,
			ExpPaymentTypeId:   "",
			ExpId:              "",
			PriceItemType:      0,
			Currency:           0,
			Price:              0,
			CustomPrice:        nil,
			ExpPaymentTypeName: "",
			ExpPaymentTypeDesc: "",
		},
	}
	expPhotos = []*models.ExpPhotos{
		&models.ExpPhotos{
			Id:             "",
			CreatedBy:      "",
			CreatedDate:    time.Time{},
			ModifiedBy:     nil,
			ModifiedDate:   nil,
			DeletedBy:      nil,
			DeletedDate:    nil,
			IsDeleted:      0,
			IsActive:       0,
			ExpPhotoFolder: "Cover Photo",
			ExpPhotoImage:  `[{"original":"https://cgostorage.blob.core.windows.net/cgo-storage/aaa/386680487554187410.jpg","thumbnail":""},{"original":"https://cgostorage.blob.core.windows.net/cgo-storage/aaa/8029052665198754438.jpg","thumbnail":""}]`,
			ExpId:          "",
		},
	}
	)
func TestList(t *testing.T) {
	mockWishlistRepo := new(mocks.Repository)
	mockWishlistObj := mockWishlistObj[0]
	var mockListWishlistObj []*models.WishlistObj
	mockListWishlistObj = append(mockListWishlistObj, &mockWishlistObj)

	t.Run("success", func(t *testing.T) {
		mockWishlistRepo.On("List", mock.Anything,mock.AnythingOfType("string"),mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(mockListWishlistObj, nil).Once()
		mockUserUsecase.On("ValidateTokenUser", mock.Anything, mock.AnythingOfType("string")).Return(userInfoDto,nil).Once()
		mockExperienceRepo.On("GetByID",mock.Anything,mock.AnythingOfType("string")).Return(experience,nil)
		mockExpPaymentRepo.On("GetByExpID",mock.Anything,mock.AnythingOfType("string")).Return(experiencePaymentType,nil)
		mockReviewsRepo.On("CountRating",mock.Anything,mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(1,nil)
		mockExpPhotosRepo.On("GetByExperienceID",mock.Anything,mock.AnythingOfType("string")).Return(expPhotos,nil)
		mockWishlistRepo.On("Count", mock.Anything,mock.AnythingOfType("string")).Return(1, nil).Once()
		u := usecase.NewWishlistUsecase(mockExpPhotosRepo, mockWishlistRepo, mockUserUsecase,mockExperienceRepo,
			mockExpPaymentRepo,mockReviewsRepo,timeoutContext)

		_, err := u.List(context.TODO(),token,1,1,0,mockWishlist[0].ExpId)

		assert.NoError(t, err)
		//assert.Len(t, list, len(mockListWishlistObj))

		mockWishlistRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockWishlistRepo.On("List", mock.Anything,mock.AnythingOfType("string"),mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(nil, errors.New("UnExpected")).Once()
		mockUserUsecase.On("ValidateTokenUser", mock.Anything, mock.AnythingOfType("string")).Return(userInfoDto,nil).Once()
		mockExperienceRepo.On("GetByID",mock.Anything,mock.AnythingOfType("string")).Return(experience,nil)
		mockExpPaymentRepo.On("GetByExpID",mock.Anything,mock.AnythingOfType("string")).Return(experiencePaymentType,nil)
		mockReviewsRepo.On("CountRating",mock.Anything,mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(1,nil)
		mockExpPhotosRepo.On("GetByExperienceID",mock.Anything,mock.AnythingOfType("string")).Return(expPhotos,nil)
		mockWishlistRepo.On("Count", mock.Anything,mock.AnythingOfType("string")).Return(1, nil).Once()
		u := usecase.NewWishlistUsecase(mockExpPhotosRepo, mockWishlistRepo, mockUserUsecase,mockExperienceRepo,
			mockExpPaymentRepo,mockReviewsRepo,timeoutContext)

		_, err := u.List(context.TODO(),token,1,1,0,mockWishlist[0].ExpId)

		assert.Error(t, err)
		//assert.Len(t, list, len(mockListWishlistObj))

		//mockWishlistRepo.AssertExpectations(t)
	})

}

func TestInsert(t *testing.T) {
	mockWishlistRepo := new(mocks.Repository)
	mockWishlistIn := &models.WishlistIn{
		IsDeleted: false,
		TransID:   mockWishlist[0].TransId,
		ExpID:     mockWishlist[0].ExpId,
	}
	mockWishlistObj := mockWishlistObj[0]
	var mockListWishlistObj []*models.WishlistObj
	mockListWishlistObj = append(mockListWishlistObj, &mockWishlistObj)
	t.Run("success-isDeleted-false", func(t *testing.T) {
		var mockListWishlistObj []*models.WishlistObj
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		mockUserUsecase.On("ValidateTokenUser", mock.Anything, mock.AnythingOfType("string")).Return(userInfoDto,nil).Once()
		mockWishlistRepo.On("Insert", mock.Anything, mock.AnythingOfType("*models.Wishlist")).Return(&mockWishlist[0],nil).Once()
		mockWishlistRepo.On("GetByUserAndExpId", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockListWishlistObj,nil).Once()

		u := usecase.NewWishlistUsecase(mockExpPhotosRepo, mockWishlistRepo, mockUserUsecase,mockExperienceRepo,
			mockExpPaymentRepo,mockReviewsRepo,timeoutContext)

		_,err := u.Insert(context.TODO(), mockWishlistIn,token)

		assert.NoError(t, err)
		//assert.Equal(t, mockWishlist.WishlistName, tempMockWishlist.WishlistName)
		mockWishlistRepo.AssertExpectations(t)
	})
	t.Run("error-isDeleted-false-validate-duplicate", func(t *testing.T) {
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		mockUserUsecase.On("ValidateTokenUser", mock.Anything, mock.AnythingOfType("string")).Return(userInfoDto,nil).Once()
		mockWishlistRepo.On("Insert", mock.Anything, mock.AnythingOfType("*models.Wishlist")).Return(&mockWishlist[0],nil).Once()
		mockWishlistRepo.On("GetByUserAndExpId", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockListWishlistObj,nil).Once()

		u := usecase.NewWishlistUsecase(mockExpPhotosRepo, mockWishlistRepo, mockUserUsecase,mockExperienceRepo,
			mockExpPaymentRepo,mockReviewsRepo,timeoutContext)

		_,err := u.Insert(context.TODO(), mockWishlistIn,token)

		assert.Error(t, err)
		//assert.Equal(t, mockWishlist.WishlistName, tempMockWishlist.WishlistName)
		//mockWishlistRepo.AssertExpectations(t)
	})
	t.Run("success-isDeleted-true", func(t *testing.T) {
		mockWishlistIn.IsDeleted = true
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		mockUserUsecase.On("ValidateTokenUser", mock.Anything, mock.AnythingOfType("string")).Return(userInfoDto,nil).Once()
		mockWishlistRepo.On("GetByUserAndExpId", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockListWishlistObj,nil).Once()
		mockWishlistRepo.On("DeleteByUserIdAndExpIdORTransId", mock.Anything,mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(nil).Once()
		u := usecase.NewWishlistUsecase(mockExpPhotosRepo, mockWishlistRepo, mockUserUsecase,mockExperienceRepo,
			mockExpPaymentRepo,mockReviewsRepo,timeoutContext)

		_,err := u.Insert(context.TODO(), mockWishlistIn,token)

		assert.NoError(t, err)
		//assert.Equal(t, mockWishlist.WishlistName, tempMockWishlist.WishlistName)
		//mockWishlistRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {
		var mockListWishlistObj []*models.WishlistObj
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		mockUserUsecase.On("ValidateTokenUser", mock.Anything, mock.AnythingOfType("string")).Return(nil,errors.New("UnAuthorize")).Once()
		mockWishlistRepo.On("Insert", mock.Anything, mock.AnythingOfType("*models.Wishlist")).Return(&mockWishlist[0],nil).Once()
		mockWishlistRepo.On("GetByUserAndExpId", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockListWishlistObj,nil).Once()

		u := usecase.NewWishlistUsecase(mockExpPhotosRepo, mockWishlistRepo, mockUserUsecase,mockExperienceRepo,
			mockExpPaymentRepo,mockReviewsRepo,timeoutContext)

		_,err := u.Insert(context.TODO(), mockWishlistIn,token)

		assert.Error(t, err)
		//assert.Equal(t, mockWishlist.WishlistName, tempMockWishlist.WishlistName)
		//mockWishlistRepo.AssertExpectations(t)
	})

}



