package http

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/service/exp_photos"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// exp_photosHandler  represent the httphandler for exp_photos
type Exp_photosHandler struct {
	Exp_photosUsecase exp_photos.Usecase
}

// Newexp_photosHandler will initialize the exp_photoss/ resources endpoint
func Newexp_photosHandler(e *echo.Echo, us exp_photos.Usecase) {
	handler := &Exp_photosHandler{
		Exp_photosUsecase: us,
	}
	//e.POST("/exp_photoss", handler.Createexp_photos)
	//e.PUT("/exp_photoss/:id", handler.Updateexp_photos)
	e.GET("service/exp_photos/:id", handler.GetByID)
	//e.DELETE("/exp_photoss/:id", handler.Delete)
}
// GetByID will get article by given id
func (a *Exp_photosHandler) GetByID(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	art, err := a.Exp_photosUsecase.GetByExperienceID(ctx, id)
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
	default:
		return http.StatusInternalServerError
	}
}
