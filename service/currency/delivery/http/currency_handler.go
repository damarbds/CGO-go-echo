package http

import (
	"context"
	"github.com/service/currency"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/models"
	"github.com/sirupsen/logrus"
)

type ResponseError struct {
	Message string `json:"message"`
}

type currencyHandler struct {
	currencyUsecase currency.Usecase
}

func NewCurrencyHandler(e *echo.Echo, currencyUsecase currency.Usecase) {
	handler := &currencyHandler{
		currencyUsecase: currencyUsecase,
	}
	e.POST("master/currency", handler.CreateCurrency)
	e.PUT("master/currency/:id", handler.UpdateCurrency)
	e.GET("master/currency", handler.GetAllCurrency)
	e.GET("master/currency/:id", handler.GetDetailCurrencyID)
	e.DELETE("master/currency/:id", handler.DeleteCurrency)
}

func (a *currencyHandler) DeleteCurrency(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	id := c.Param("id")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	facilitiesId,_:= strconv.Atoi(id)
	result, err := a.currencyUsecase.Delete(ctx, facilitiesId, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

// GetByID will get article by given id
func (a *currencyHandler) GetAllCurrency(c echo.Context) error {
	qpage := c.QueryParam("page")
	qsize := c.QueryParam("size")

	var limit = 20
	var page = 1
	var offset = 0

	page, _ = strconv.Atoi(qpage)
	limit, _ = strconv.Atoi(qsize)
	offset = (page - 1) * limit
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	art, err := a.currencyUsecase.GetAll(ctx, page,limit,offset)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, art)

}
// Store will store the user by given request body
func (a *currencyHandler) CreateCurrency(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	var currencyCommand models.NewCommandCurrency
	err := c.Bind(&currencyCommand)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	result,error := a.currencyUsecase.Create(ctx,&currencyCommand,token)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (a *currencyHandler) UpdateCurrency(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	var currencyCommand models.NewCommandCurrency
	err := c.Bind(&currencyCommand)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	result,error := a.currencyUsecase.Update(ctx,&currencyCommand,token)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (a *currencyHandler) GetDetailCurrencyID(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	facilitiesID ,_:= strconv.Atoi(id)
	result, err := a.currencyUsecase.GetById(ctx,facilitiesID)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
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
