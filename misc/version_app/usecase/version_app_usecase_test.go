package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/misc/version_app/mocks"
	"github.com/misc/version_app/usecase"
	"github.com/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	timeoutContext = time.Second * 30
)

func TestGetAllVersion(t *testing.T) {
	mockVersionAppRepo := new(mocks.Repository)
	mockVersionApp := []*models.VersionApp{
		&models.VersionApp{
			Id:          1,
			VersionCode: 1,
			VersionName: "ljdlkadj",
			Type:        1,
		},
	}
	t.Run("success", func(t *testing.T) {
		mockVersionAppRepo.On("GetAll", mock.Anything, mock.AnythingOfType("int")).Return(mockVersionApp, nil).Once()
		u := usecase.NewVersionAPPUsecase(mockVersionAppRepo, timeoutContext)

		a, err := u.GetAllVersion(context.TODO(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockVersionAppRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockVersionAppRepo.On("GetAll", mock.Anything, mock.AnythingOfType("int")).Return(nil, models.ErrNotFound).Once()
		u := usecase.NewVersionAPPUsecase(mockVersionAppRepo, timeoutContext)

		a, err := u.GetAllVersion(context.TODO(), 1)

		assert.Error(t, err)
		assert.Nil(t, a)

		mockVersionAppRepo.AssertExpectations(t)
	})

}
