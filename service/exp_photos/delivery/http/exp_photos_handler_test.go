package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/service/exp_photos/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ExpPhotosHttp "github.com/service/exp_photos/delivery/http"
)

var (
	imagePath = `[{"original":"https://cgostorage.blob.core.windows.net/cgo-storage/Experience/6569483590502428383.jpg","thumbnail":""}]`
	mockExpPhotos = models.ExpPhotos{
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
)

func TestGetByID(t *testing.T) {
	var mockExpPhotos []models.ExpPhotosDto
	err := faker.FakeData(&mockExpPhotos)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)

	num := mockExpPhotos[0].ExpId

	mockUCase.On("GetByExperienceID", mock.Anything, num).Return(mockExpPhotos, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/exp_photos/"+num, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("master/exp_photos/:id")
	c.SetParamNames("id")
	c.SetParamValues(num)
	handler := ExpPhotosHttp.Exp_photosHandler{
		Exp_photosUsecase: mockUCase,
	}
	err = handler.GetByID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestGetByIDEInternalServerError(t *testing.T) {
	var mockExpPhotos []models.ExpPhotosDto
	err := faker.FakeData(&mockExpPhotos)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)

	num := mockExpPhotos[0].Id

	mockUCase.On("GetByExperienceID", mock.Anything, num).Return(nil, errors.New("unExpected"))

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/exp_photos/"+num, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("master/exp_photos/:id")
	c.SetParamNames("id")
	c.SetParamValues(num)
	handler := ExpPhotosHttp.Exp_photosHandler{
		Exp_photosUsecase: mockUCase,
	}
	err = handler.GetByID(c)
	//assert.Error(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}
