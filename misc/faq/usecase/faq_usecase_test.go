package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/misc/faq/mocks"
	"github.com/misc/faq/usecase"
	"github.com/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	timeoutContext = time.Second * 30
)

func TestGetByType(t *testing.T) {
	mockFAQRepo := new(mocks.Repository)
	mockFAQ := []*models.FAQ{
		&models.FAQ{
			Id:           2,
			CreatedBy:    "Test 2",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			Type:         1,
			Title:        "Test Faq",
			Desc:         "Faq Desc",
		},
	}
	t.Run("success", func(t *testing.T) {
		mockFAQRepo.On("GetByType", mock.Anything, mock.AnythingOfType("int")).Return(mockFAQ, nil).Once()
		u := usecase.NewfaqUsecase(mockFAQRepo, timeoutContext)

		a, err := u.GetByType(context.TODO(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockFAQRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockFAQRepo.On("GetByType", mock.Anything, mock.AnythingOfType("int")).Return(nil, models.ErrNotFound).Once()
		u := usecase.NewfaqUsecase(mockFAQRepo, timeoutContext)

		a, err := u.GetByType(context.TODO(), 1)

		assert.Error(t, err)
		assert.Nil(t, a)

		mockFAQRepo.AssertExpectations(t)
	})

}
