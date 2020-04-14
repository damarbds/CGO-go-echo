package http

import (
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/product/reviews"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	validator "gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
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

func (a *reviewsHandler) GetReviewsByExpId(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")
	isMerchant := c.QueryParam("isMerchant")

	if isMerchant != "" {
		if token != "" {
			return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
		}
	}

	expId := c.QueryParam("exp_id")
	qpage := c.QueryParam("page")
	qperPage := c.QueryParam("size")
	sortBy := c.QueryParam("sortBy")
	qRating := c.QueryParam("rating")

	var limit = 20
	var page = 1
	var offset = 0

	page, _ = strconv.Atoi(qpage)
	limit, _ = strconv.Atoi(qperPage)
	offset = (page - 1) * limit

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var rating int
	if qRating != "" {
		rating, _ = strconv.Atoi(qRating)
	}
	res, err := a.reviewsUsecase.GetReviewsByExpIdWithPagination(ctx, page, limit, offset, rating, sortBy, expId)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
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
