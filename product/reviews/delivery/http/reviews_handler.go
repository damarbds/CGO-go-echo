package http

import (
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/product/reviews"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	validator "gopkg.in/go-playground/validator.v9"
	"net/http"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// reviewsHandler  represent the httphandler for reviews
type reviewsHandler struct {
	reviewsUsecase reviews.Usecase
}

// NewreviewsHandler will initialize the reviewss/ resources endpoint
func NewreviewsHandler(e *echo.Echo, us reviews.Usecase) {
	handler := &reviewsHandler{
		reviewsUsecase: us,
	}
	//e.POST("/reviewss", handler.Createreviews)
	//e.PUT("/reviewss/:id", handler.Updatereviews)
	e.GET("product/exp-reviews", handler.GetReviewsByExpId)
	//e.DELETE("/reviewss/:id", handler.Delete)
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
func (a *reviewsHandler) GetReviewsByExpId(c echo.Context) error {
	expId := c.QueryParam("exp_id")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	art, err := a.reviewsUsecase.GetReviewsByExpId(ctx, expId)
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
