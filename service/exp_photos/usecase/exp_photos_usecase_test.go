package usecase_test

import (
	"context"
	"errors"
	"github.com/models"
	"github.com/service/exp_photos/mocks"
	"github.com/service/exp_photos/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)
var(
	timeoutContext = time.Second*30
)
func TestGetByExperienceID(t *testing.T) {
	imagePath := `[{"original":"https://cgostorage.blob.core.windows.net/cgo-storage/Experience/6569483590502428383.jpg","thumbnail":""}]`
	mockExpPhotosRepo := new(mocks.Repository)
	var mockExpPhotos []*models.ExpPhotos
	mockExpPhoto := models.ExpPhotos{
		Id:             "adsasdasda",
		CreatedBy:      "Test 1",
		CreatedDate:    time.Now(),
		ModifiedBy:     nil,
		ModifiedDate:   nil,
		DeletedBy:      nil,
		DeletedDate:    nil,
		IsDeleted:      0,
		IsActive:       1,
		ExpPhotoFolder: "Facilities",
		ExpPhotoImage:  imagePath,
		ExpId:          "qweqwewq",
	}
	mockExpPhotos = append(mockExpPhotos,&mockExpPhoto)
	t.Run("success", func(t *testing.T) {
		mockExpPhotosRepo.On("GetByExperienceID", mock.Anything, mock.AnythingOfType("string")).Return(mockExpPhotos, nil).Once()
		u := usecase.Newexp_photosUsecase(mockExpPhotosRepo, timeoutContext)

		a, err := u.GetByExperienceID(context.TODO(), mockExpPhotos[0].ExpId)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockExpPhotosRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockExpPhotosRepo.On("GetByExperienceID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpected")).Once()

		u := usecase.Newexp_photosUsecase( mockExpPhotosRepo, timeoutContext)

		a, err := u.GetByExperienceID(context.TODO(),  mockExpPhotos[0].ExpId)

		assert.Error(t, err)
		assert.Nil(t, a)

		mockExpPhotosRepo.AssertExpectations(t)
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

