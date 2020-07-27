package usecase_test

import (
	"context"
	"errors"
	_adminUsecaseMock "github.com/auth/admin/mocks"
	_expPhotosMock "github.com/service/exp_photos/mocks"
	_experienceMock "github.com/service/experience/mocks"
	_userMock "github.com/auth/user/mocks"
	_expPayment "github.com/service/exp_payment/mocks"
	"github.com/models"
	"github.com/profile/wishlists/mocks"
	"github.com/profile/wishlists/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)
var(
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
	mockExpPhotos = new()
)
func TestList(t *testing.T) {
	mockWishlistRepo := new(mocks.Repository)
	mockWishlist := mockWishlist[0]
	var mockListWishlist []*models.Wishlist
	mockListWishlist = append(mockListWishlist, &mockWishlist)

	t.Run("success", func(t *testing.T) {
		mockWishlistRepo.On("List", mock.Anything,mock.AnythingOfType("string"),mock.AnythingOfType("int"),
			mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(mockListWishlist, nil).Once()

		u := usecase.NewWishlistUsecase(nil, mockWishlistRepo, timeoutContext)

		list, err := u.List(context.TODO())

		assert.NoError(t, err)
		assert.Len(t, list, len(mockListWishlist))

		mockWishlistRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockWishlistRepo.On("List", mock.Anything).Return(nil,errors.New("Unexpexted Error")).Once()

		u := usecase.NewWishlistUsecase(nil, mockWishlistRepo, timeoutContext)

		list, err := u.List(context.TODO())

		assert.Error(t, err)
		assert.Len(t, list, 0)
		mockWishlistRepo.AssertExpectations(t)
	})

}

func TestInsert(t *testing.T) {
	mockWishlistRepo := new(mocks.Repository)
	mockAdminUsecase := new(_adminUsecaseMock.Usecase)
	mockWishlist := 	models.Wishlist{
		Id:           0,
		CreatedBy:    "adminCGO",
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		WishlistName:  "Test Wishlist 2",
		WishlistIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Wishlist/8941695193938718058.jpg",
	}
	mockAdmin := &models.AdminDto{
		Id:    "972fe09f-93e9-4798-b642-14e0ac77c6be",
		Name:  "adminCGO",
		Email: "admin1234@gmail.com",
	}
	t.Run("success", func(t *testing.T) {
		tempMockWishlist := models.NewCommandWishlist{
			Id:          0,
			WishlistName: mockWishlist.WishlistName,
			WishlistIcon: mockWishlist.WishlistIcon,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockWishlist.Id = 0
		mockWishlistRepo.On("Insert", mock.Anything, mock.AnythingOfType("*models.Wishlist")).Return(&mockWishlist.Id,nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()

		u := usecase.NewWishlistUsecase(mockAdminUsecase, mockWishlistRepo, timeoutContext)

		_,err := u.Create(context.TODO(), &tempMockWishlist,token)

		assert.NoError(t, err)
		assert.Equal(t, mockWishlist.WishlistName, tempMockWishlist.WishlistName)
		mockWishlistRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {
		tempMockWishlist := models.NewCommandWishlist{
			Id:          0,
			WishlistName: mockWishlist.WishlistName,
			WishlistIcon: mockWishlist.WishlistIcon,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockWishlist.Id = 0
		mockWishlistRepo.On("Insert", mock.Anything, mock.AnythingOfType("*models.Wishlist")).Return(&mockWishlist.Id,nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil,errors.New("UnAuthorize")).Once()

		u := usecase.NewWishlistUsecase(mockAdminUsecase, mockWishlistRepo, timeoutContext)

		_,err := u.Create(context.TODO(), &tempMockWishlist,token)

		assert.Error(t, err)
		//assert.Equal(t, mockWishlist.WishlistName, tempMockWishlist.WishlistName)
		//mockWishlistRepo.AssertExpectations(t)
	})

}



