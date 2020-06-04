package http

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/models"
	"github.com/service/exp_photos"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	validator "gopkg.in/go-playground/validator.v9"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// exp_photosHandler  represent the httphandler for exp_photos
type exp_photosHandler struct {
	exp_photosUsecase exp_photos.Usecase
}

// Newexp_photosHandler will initialize the exp_photoss/ resources endpoint
func Newexp_photosHandler(e *echo.Echo, us exp_photos.Usecase) {
	handler := &exp_photosHandler{
		exp_photosUsecase: us,
	}
	//e.POST("/exp_photoss", handler.Createexp_photos)
	//e.PUT("/exp_photoss/:id", handler.Updateexp_photos)
	e.GET("service/exp_photos/:id", handler.GetByID)
	//e.DELETE("/exp_photoss/:id", handler.Delete)
}

func isRequestValid(m *models.NewCommandMerchant) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetByID will get article by given id
func (a *exp_photosHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	//if err != nil {
	//	return c.JSON(http.StatusNotFound, models.ErrNotFound.Error())
	//}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	art, err := a.exp_photosUsecase.GetByExperienceID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, art)
}
func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	logrus.Error(err)
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrUnAuthorize:
		return http.StatusUnauthorized
	case models.ErrConflict:
		return http.StatusBadRequest
	case models.ErrBadParamInput:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
