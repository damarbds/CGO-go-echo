package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/models"
	"github.com/service/harbors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	validator "gopkg.in/go-playground/validator.v9"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// harborsHandler  represent the httphandler for harbors
type harborsHandler struct {
	harborsUsecase harbors.Usecase
}

// NewharborsHandler will initialize the harborss/ resources endpoint
func NewharborsHandler(e *echo.Echo, us harbors.Usecase) {
	handler := &harborsHandler{
		harborsUsecase: us,
	}
	//e.POST("/harborss", handler.Createharbors)
	//e.PUT("/harborss/:id", handler.Updateharbors)
	e.GET("service/exp-destination", handler.GetAllHarbors)
	//e.DELETE("/harborss/:id", handler.Delete)
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
func (a *harborsHandler) GetAllHarbors(c echo.Context) error {
	qpage := c.QueryParam("page")
	qsize := c.QueryParam("size")
	search := c.QueryParam("search")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	if qpage != "" && qsize != "" {
		page, _ := strconv.Atoi(qpage)
		size, _ := strconv.Atoi(qsize)
		art, err := a.harborsUsecase.GetAllWithJoinCPC(ctx, &size, &page, search)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, art)
	} else {
		art, err := a.harborsUsecase.GetAllWithJoinCPC(ctx, nil, nil, search)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, art)
	}

	return nil
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
