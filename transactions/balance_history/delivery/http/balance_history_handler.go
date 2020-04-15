package http

import (
	"context"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"github.com/transactions/balance_history"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

type balanceHistoryHandler struct {
	balanceHistoryUsecase balance_history.Usecase
}

func NewBalanceHistoryHandler(e *echo.Echo, bh balance_history.Usecase) {
	handler := &balanceHistoryHandler{
		balanceHistoryUsecase: bh,
	}
	e.POST("/transaction/withdraw", handler.CreateBalanceHistory)
	e.GET("/transaction/withdraw-history", handler.GetBalanceHistory)
}

func isRequestValid(m *models.NewBalanceHistoryCommand) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (p *balanceHistoryHandler) GetBalanceHistory(c echo.Context) error {
	merchantId := c.QueryParam("merchant_id")
	status := c.QueryParam("status")
	qpage := c.QueryParam("page")
	qperPage := c.QueryParam("size")

	var limit *int
	var page = 1
	var offset *int
	if qpage != "" && qperPage != ""{
		page, _ = strconv.Atoi(qpage)
		limits, _ := strconv.Atoi(qperPage)
		limit = &limits
		offsets := (page - 1) * *limit
		offset = &offsets
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res,err := p.balanceHistoryUsecase.List(ctx,merchantId,status,page,limit,offset)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (p *balanceHistoryHandler) CreateBalanceHistory(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}
	t := new(models.NewBalanceHistoryCommand)
	if err := c.Bind(t); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrBadParamInput)
	}

	if ok, err := isRequestValid(t); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := p.balanceHistoryUsecase.Create(ctx,*t,token)
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
	case models.BookingTypeRequired:
		return http.StatusBadRequest
	case models.BookingExpIdRequired:
		return http.StatusBadRequest
	case models.PromoIdRequired:
		return http.StatusBadRequest
	case models.PaymentMethodIdRequired:
		return http.StatusBadRequest
	case models.ExpPaymentIdRequired:
		return http.StatusBadRequest
	case models.StatusRequired:
		return http.StatusBadRequest
	case models.TotalPriceRequired:
		return http.StatusBadRequest
	case models.CurrencyRequired:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
