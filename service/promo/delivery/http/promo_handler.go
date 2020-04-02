package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/models"
	"github.com/service/promo"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	validator "gopkg.in/go-playground/validator.v9"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// promoHandler  represent the httphandler for promo
type promoHandler struct {
	promoUsecase promo.Usecase
}

// NewpromoHandler will initialize the promos/ resources endpoint
func NewpromoHandler(e *echo.Echo, us promo.Usecase) {
	handler := &promoHandler{
		promoUsecase: us,
	}
	//e.POST("/promos", handler.Createpromo)
	//e.PUT("/promos/:id", handler.Updatepromo)
	e.GET("service/special-promo", handler.GetAllPromo)
	e.GET("service/special-promo/:code", handler.GetPromoByCode)
	//e.DELETE("/promos/:id", handler.Delete)
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
func (a *promoHandler) GetAllPromo(c echo.Context) error {
	qpage := c.QueryParam("page")
	qsize := c.QueryParam("size")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	if qpage != "" && qsize != "" {
		page, _ := strconv.Atoi(qpage)
		size, _ := strconv.Atoi(qsize)
		art, err := a.promoUsecase.Fetch(ctx, &page,&size)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, art)
	} else {
		art, err := a.promoUsecase.Fetch(ctx, nil, nil)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, art)
	}
}

func (a *promoHandler) GetPromoByCode(c echo.Context) error {
	code := c.Param("code")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	results, err := a.promoUsecase.GetByCode(ctx, code)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, results)
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
