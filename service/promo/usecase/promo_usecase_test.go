package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	_adminUsecaseMock "github.com/auth/admin/mocks"
	_userUsecaseMock "github.com/auth/user/mocks"
	"github.com/models"
	"github.com/service/promo/mocks"
	"github.com/service/promo/usecase"
	_promoMerchantRepoMock "github.com/service/promo_merchant/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	_transactionRepoMock "github.com/transactions/transaction/mocks"
)

var (
	mockPromoRepo = new(mocks.Repository)
	mockAdminUsecase = new(_adminUsecaseMock.Usecase)
	mockPromoMerchantRepo = new(_promoMerchantRepoMock.Repository)
	mockUserUsecase = new(_userUsecaseMock.Usecase)
	mockTransactionRepo = new(_transactionRepoMock.Repository)
	timeoutContext = time.Second * 30
	date = time.Now().String()
	isAnyTripPeriod = 1
	capacity = 20
	maxUsage = 20
	desc = "ini description"
	disc float32 = 12.2
	mockPromo = models.Promo{
		Id:                 "asdqeqrqrasdsad",
		CreatedBy:          "test",
		CreatedDate:        time.Now(),
		ModifiedBy:         nil,
		ModifiedDate:       nil,
		DeletedBy:          nil,
		DeletedDate:        nil,
		IsDeleted:          0,
		IsActive:           1,
		PromoCode:          "Test1",
		PromoName:          "Test 1",
		PromoDesc:          "Test 1",
		PromoValue:         1,
		PromoType:          1,
		PromoImage:         "asdasdasdas",
		StartDate:          &date,
		EndDate:            &date,
		StartTripPeriod:    &date,
		EndTripPeriod:      &date,
		IsAnyTripPeriod:    &isAnyTripPeriod,
		HowToGet:           &desc,
		HowToUse:           &desc,
		TermCondition:      &desc,
		Disclaimer:         &desc,
		MaxDiscount:        &disc,
		MaxUsage:           &maxUsage,
		ProductionCapacity: &capacity,
		CurrencyId:         &isAnyTripPeriod,
		PromoProductType:   &isAnyTripPeriod,
	}
	mockPromoMerchantList = []*models.PromoMerchant{
		&models.PromoMerchant{
			Id:         1,
			PromoId:    "qeqwewq",
			MerchantId: "adsadsa",
		},
	}
	tempMockPromo = models.NewCommandPromo{
		Id:                 "asdqeqrqrasdsad",
		PromoCode:          "Test1",
		PromoName:          "Test 1",
		PromoDesc:          "Test 1",
		PromoValue:         1,
		PromoType:          1,
		PromoImage:         "asdasdasdas",
		StartDate:          date,
		EndDate:            date,
		StartTripPeriod:    date,
		EndTripPeriod:      date,
		IsAnyTripPeriod:    isAnyTripPeriod,
		HowToGet:           desc,
		HowToUse:           desc,
		TermCondition:      desc,
		Disclaimer:         desc,
		MaxDiscount:        disc,
		MaxUsage:           isAnyTripPeriod,
		ProductionCapacity: isAnyTripPeriod,
		PromoProductType:   &isAnyTripPeriod,
		MerchantId:[]string{
			"adsadsa",// You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
	}
)

func TestList(t *testing.T) {

	var mockListPromo []*models.Promo
	mockListPromo = append(mockListPromo, &mockPromo)
	page := 1
	limit := 1
	offset := 0
	count := 1
	mockAdmin := &models.AdminDto{
		Id:    "972fe09f-93e9-4798-b642-14e0ac77c6be",
		Name:  "adminCGO",
		Email: "admin1234@gmail.com",
	}
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	t.Run("success", func(t *testing.T) {
		mockPromoRepo.On("Fetch", mock.Anything, mock.AnythingOfType("*int"), mock.AnythingOfType("*int"),mock.AnythingOfType("string")).Return(mockListPromo, nil).Once()
		mockPromoRepo.On("GetCount", mock.Anything).Return(count, nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin, nil).Once()
		mockPromoMerchantRepo.On("GetByMerchantId",mock.Anything,mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(mockPromoMerchantList,nil)
		u := usecase.NewPromoUsecase(mockUserUsecase,mockTransactionRepo,mockPromoMerchantRepo, mockPromoRepo,mockAdminUsecase, timeoutContext)

		_, err := u.List(context.TODO(), page, limit, offset,"",token,true,true,make([]string,0))

		assert.NoError(t, err)
		//assert.Len(t, list, len(mockListPromo))

		mockPromoRepo.AssertExpectations(t)
		mockPromoRepo.AssertExpectations(t)
	})

}
func TestFetch(t *testing.T) {
	var mockListPromo []*models.Promo
	mockListPromo = append(mockListPromo, &mockPromo)
	//page := 1
	limit := 1
	offset := 0
	//count := 1
	//mockAdmin := &models.AdminDto{
	//	Id:    "972fe09f-93e9-4798-b642-14e0ac77c6be",
	//	Name:  "adminCGO",
	//	Email: "admin1234@gmail.com",
	//}
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	t.Run("success", func(t *testing.T) {
		mockPromoRepo.On("Fetch", mock.Anything, mock.AnythingOfType("*int"), mock.AnythingOfType("*int"),mock.AnythingOfType("string")).Return(mockListPromo, nil).Once()
		mockPromoMerchantRepo.On("GetByMerchantId",mock.Anything,mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(mockPromoMerchantList,nil)
		u := usecase.NewPromoUsecase(mockUserUsecase,mockTransactionRepo,mockPromoMerchantRepo, mockPromoRepo,mockAdminUsecase, timeoutContext)

		_, err := u.Fetch(context.TODO(), &offset,&limit)

		assert.NoError(t, err)
		//assert.Len(t, list, len(mockListPromo))

		mockPromoRepo.AssertExpectations(t)
		mockPromoRepo.AssertExpectations(t)
	})

}

func TestGetDetail(t *testing.T) {
	mockAdmin := &models.AdminDto{
		Id:    "972fe09f-93e9-4798-b642-14e0ac77c6be",
		Name:  "adminCGO",
		Email: "admin1234@gmail.com",
	}
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	t.Run("success", func(t *testing.T) {
		mockPromoRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(&mockPromo, nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin, nil).Once()
		mockPromoMerchantRepo.On("GetByMerchantId",mock.Anything,mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(mockPromoMerchantList,nil)

		u := usecase.NewPromoUsecase(mockUserUsecase,mockTransactionRepo,mockPromoMerchantRepo, mockPromoRepo,mockAdminUsecase, timeoutContext)

		a, err := u.GetDetail(context.TODO(), mockPromo.Id,token)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockPromoRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockPromoRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpected")).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin, nil).Once()
		mockPromoMerchantRepo.On("GetByMerchantId",mock.Anything,mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(mockPromoMerchantList,nil)

		u := usecase.NewPromoUsecase(mockUserUsecase,mockTransactionRepo,mockPromoMerchantRepo, mockPromoRepo,mockAdminUsecase, timeoutContext)

		a, err := u.GetDetail(context.TODO(), mockPromo.Id,token)

		assert.Error(t, err)
		assert.Nil(t, a)

		mockPromoRepo.AssertExpectations(t)
	})

}
func TestGetByCode(t *testing.T) {
	var mockListPromo []*models.Promo
	mockListPromo = append(mockListPromo, &mockPromo)
	mockUser := &models.UserInfoDto{
		Id:             "kjkljljlkjlk",
		CreatedDate:    time.Now(),
		UpdatedDate:    nil,
		IsActive:       1,
		UserEmail:      "test@gmail.com",
		FullName:       "asdasdas",
		PhoneNumber:    "098907097",
		ProfilePictUrl: "dasdasdsa",
		Address:        "",
		Dob:            time.Now(),
		Gender:         0,
		IdType:         0,
		IdNumber:       "",
		ReferralCode:   "",
		Points:         0,
	}
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	t.Run("success", func(t *testing.T) {
		mockPromoRepo.On("GetByCode", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("*int"),mock.AnythingOfType("string")).Return(mockListPromo, nil).Once()
		mockUserUsecase.On("ValidateTokenUser", mock.Anything, mock.AnythingOfType("string")).Return(mockUser, nil).Once()
		mockTransactionRepo.On("GetCountTransactionByPromoId",mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(1,nil)
		u := usecase.NewPromoUsecase(mockUserUsecase,mockTransactionRepo,mockPromoMerchantRepo, mockPromoRepo,mockAdminUsecase, timeoutContext)

		a, err := u.GetByCode(context.TODO(), mockPromo.PromoCode,*mockPromo.PromoProductType,tempMockPromo.MerchantId[0],token)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockPromoRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockPromoRepo.On("GetByCode", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("*int"),mock.AnythingOfType("string")).Return(nil, errors.New("Unexpected")).Once()
		mockUserUsecase.On("ValidateTokenUser", mock.Anything, mock.AnythingOfType("string")).Return(mockUser, nil).Once()
		mockTransactionRepo.On("GetCountTransactionByPromoId",mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(1,nil)
		u := usecase.NewPromoUsecase(mockUserUsecase,mockTransactionRepo,mockPromoMerchantRepo, mockPromoRepo,mockAdminUsecase, timeoutContext)

		a, err := u.GetByCode(context.TODO(), mockPromo.PromoCode,*mockPromo.PromoProductType,tempMockPromo.MerchantId[0],token)


		assert.Error(t, err)
		assert.Nil(t, a)

		mockPromoRepo.AssertExpectations(t)
	})

}

func TestCreate(t *testing.T) {

	mockAdmin := &models.AdminDto{
		Id:    "972fe09f-93e9-4798-b642-14e0ac77c6be",
		Name:  "adminCGO",
		Email: "admin1234@gmail.com",
	}
	t.Run("success", func(t *testing.T) {

		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockPromo.Id = ""
		mockPromoRepo.On("Insert", mock.Anything, mock.AnythingOfType("*models.Promo")).Return(mockPromo.Id, nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin, nil).Once()
		mockPromoMerchantRepo.On("Insert", mock.Anything, mock.AnythingOfType("models.PromoMerchant")).Return(nil).Once()

		u := usecase.NewPromoUsecase(mockUserUsecase,mockTransactionRepo,mockPromoMerchantRepo, mockPromoRepo,mockAdminUsecase, timeoutContext)

		_, err := u.Create(context.TODO(), tempMockPromo, token)

		assert.NoError(t, err)
		assert.Equal(t, mockPromo.PromoName, tempMockPromo.PromoName)
		mockPromoRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {

		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockPromo.Id = ""
		mockPromoRepo.On("Insert", mock.Anything, mock.AnythingOfType("*models.Promo")).Return(mockPromo.Id, nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("UnAuthorize")).Once()
		mockPromoMerchantRepo.On("Insert", mock.Anything, mock.AnythingOfType("models.PromoMerchant")).Return(nil).Once()

		u := usecase.NewPromoUsecase(mockUserUsecase,mockTransactionRepo,mockPromoMerchantRepo, mockPromoRepo,mockAdminUsecase, timeoutContext)

		_, err := u.Create(context.TODO(), tempMockPromo, token)

		assert.Error(t, err)
		//assert.Equal(t, mockPromo.PromoName, tempMockPromo.PromoName)
		//mockPromoRepo.AssertExpectations(t)
	})

}

func TestUpdate(t *testing.T) {

	now := time.Now()

	mockAdmin := &models.AdminDto{
		Id:    "972fe09f-93e9-4798-b642-14e0ac77c6be",
		Name:  "adminCGO",
		Email: "admin1234@gmail.com",
	}
	mockPromo.ModifiedDate = &now
	mockPromo.ModifiedBy = &mockAdmin.Name
	t.Run("success", func(t *testing.T) {

		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockPromo.Id = ""
		mockPromoRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Promo")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin, nil).Once()
		mockPromoMerchantRepo.On("Insert", mock.Anything, mock.AnythingOfType("models.PromoMerchant")).Return(nil).Once()
		mockPromoMerchantRepo.On("DeleteByMerchantId", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(nil).Once()

		u := usecase.NewPromoUsecase(mockUserUsecase,mockTransactionRepo,mockPromoMerchantRepo, mockPromoRepo,mockAdminUsecase, timeoutContext)

		_, err := u.Update(context.TODO(), tempMockPromo, token)

		assert.NoError(t, err)
		assert.Equal(t, mockPromo.PromoName, tempMockPromo.PromoName)
		//mockPromoRepo.AssertExpectations(t)
	})
	t.Run("success-without-image", func(t *testing.T) {

		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockPromo.Id = ""
		tempMockPromo.PromoImage = ""
		mockPromoRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Promo")).Return(nil).Once()
		mockPromoRepo.On("GetById", mock.Anything, mock.AnythingOfType("string")).Return(&mockPromo, nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin, nil).Once()
		mockPromoMerchantRepo.On("Insert", mock.Anything, mock.AnythingOfType("models.PromoMerchant")).Return(nil).Once()
		mockPromoMerchantRepo.On("DeleteByMerchantId", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(nil).Once()

		u := usecase.NewPromoUsecase(mockUserUsecase,mockTransactionRepo,mockPromoMerchantRepo, mockPromoRepo,mockAdminUsecase, timeoutContext)

		_, err := u.Update(context.TODO(), tempMockPromo, token)

		assert.NoError(t, err)
		assert.Equal(t, mockPromo.PromoName, tempMockPromo.PromoName)
		//mockPromoRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {

		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockPromo.Id = ""
		mockPromoRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Promo")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("unAuthorize")).Once()
		mockPromoMerchantRepo.On("Insert", mock.Anything, mock.AnythingOfType("models.PromoMerchant")).Return(nil).Once()
		mockPromoMerchantRepo.On("DeleteByMerchantId", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(nil).Once()

		u := usecase.NewPromoUsecase(mockUserUsecase,mockTransactionRepo,mockPromoMerchantRepo, mockPromoRepo,mockAdminUsecase, timeoutContext)

		_, err := u.Update(context.TODO(), tempMockPromo, token)

		assert.Error(t, err)
		//assert.Equal(t, mockPromo.PromoName, tempMockPromo.PromoName)
		//mockPromoRepo.AssertExpectations(t)
	})
}
func TestDelete(t *testing.T) {

	now := time.Now()

	mockAdmin := &models.AdminDto{
		Id:    "972fe09f-93e9-4798-b642-14e0ac77c6be",
		Name:  "adminCGO",
		Email: "admin1234@gmail.com",
	}
	mockPromo.DeletedDate = &now
	mockPromo.DeletedBy = &mockAdmin.Name
	t.Run("success", func(t *testing.T) {

		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

		mockPromoRepo.On("Delete", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin, nil).Once()

		u := usecase.NewPromoUsecase(mockUserUsecase,mockTransactionRepo,mockPromoMerchantRepo, mockPromoRepo,mockAdminUsecase, timeoutContext)

		_, err := u.Delete(context.TODO(), mockPromo.Id, token)

		assert.NoError(t, err)
		//assert.Equal(t, mockPromo.PromoName, tempMockPromo.PromoName)
		//mockPromoRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {

		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

		mockPromoRepo.On("Delete", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("unAuthorize")).Once()

		u := usecase.NewPromoUsecase(mockUserUsecase,mockTransactionRepo,mockPromoMerchantRepo, mockPromoRepo,mockAdminUsecase, timeoutContext)

		_, err := u.Delete(context.TODO(), mockPromo.Id, token)

		assert.Error(t, err)
		//assert.Equal(t, mockPromo.PromoName, tempMockPromo.PromoName)
		//mockPromoRepo.AssertExpectations(t)
	})
}

//func TestDelete(t *testing.T) {
//	mockArticleRepo := new(mocks.Repository)
//	mockArticle := models.Article{
//		Title:   "Hello",
//		Content: "Content",
//	}
//
//	t.Run("success", func(t *testing.T) {
//		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(&mockArticle, nil).Once()
//
//		mockArticleRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()
//
//		mockAuthorrepo := new(_authorMock.Repository)
//		u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)
//
//		err := u.Delete(context.TODO(), mockArticle.ID)
//
//		assert.NoError(t, err)
//		mockArticleRepo.AssertExpectations(t)
//		mockAuthorrepo.AssertExpectations(t)
//	})
//	t.Run("article-is-not-exist", func(t *testing.T) {
//		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(nil, nil).Once()
//
//		mockAuthorrepo := new(_authorMock.Repository)
//		u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)
//
//		err := u.Delete(context.TODO(), mockArticle.ID)
//
//		assert.Error(t, err)
//		mockArticleRepo.AssertExpectations(t)
//		mockAuthorrepo.AssertExpectations(t)
//	})
//	t.Run("error-happens-in-db", func(t *testing.T) {
//		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(nil, errors.New("Unexpected Error")).Once()
//
//		mockAuthorrepo := new(_authorMock.Repository)
//		u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)
//
//		err := u.Delete(context.TODO(), mockArticle.ID)
//
//		assert.Error(t, err)
//		mockArticleRepo.AssertExpectations(t)
//		mockAuthorrepo.AssertExpectations(t)
//	})
//
//}
