package http

import (
	"context"
	"net/http"
	"strconv"

	excelize "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"github.com/transactions/balance_history"
	"gopkg.in/go-playground/validator.v9"
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

	e.PUT("/transaction/withdraw-accept-decline", handler.Confirm)
	e.PUT("/transaction/withdraw-amount", handler.UpdateAmount)
	e.POST("/transaction/withdraw", handler.CreateBalanceHistory)
	e.GET("/transaction/withdraw-history", handler.GetBalanceHistory)
	e.GET("/transaction/export-excel", handler.ExportExcel)
}

func isRequestValid(m *models.NewBalanceHistoryCommand) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
func PrepareAndReturnExcel(balanceHistory *models.BalanceHistoryDtoWithPagination) *excelize.File {
	f := excelize.NewFile()

	f.SetCellValue("Sheet1", "A1", "Status")
	f.SetCellValue("Sheet1", "B1", "Amount")
	f.SetCellValue("Sheet1", "C1", "Date Of Request")
	f.SetCellValue("Sheet1", "D1", "Date Of Payment")
	f.SetCellValue("Sheet1", "E1", "Remarks")
	for index,element := range balanceHistory.Data{
		if index == 0 {
			index = index + 3
		}else {
			index = index + 2
		}
		f.SetCellValue("Sheet1", "A" + strconv.Itoa(index), element.Status)
		f.SetCellValue("Sheet1", "B" + strconv.Itoa(index), element.Amount)
		f.SetCellValue("Sheet1", "C"+ strconv.Itoa(index), element.DateOfRequest.Format("2006-01-02"))
		f.SetCellValue("Sheet1", "D"+ strconv.Itoa(index), element.DateOfPayment.Format("2006-01-02"))
		f.SetCellValue("Sheet1", "E"+ strconv.Itoa(index), element.Remarks)
	}
	return f
}
func (t *balanceHistoryHandler) ExportExcel(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	merchantId := c.QueryParam("merchant_id")
	month := c.QueryParam("month")
	year := c.QueryParam("year")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	getresult, err := t.balanceHistoryUsecase.List(ctx, merchantId, "", 1, nil, nil, month, year,token,false)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	file := PrepareAndReturnExcel(getresult)
	// Set the headers necessary to get browsers to interpret the downloadable file
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	c.Response().Header().Set("Content-Disposition", "attachment;filename=balance.xlsx")
	c.Response().Header().Set("File-Name", "userInputData.xlsx")
	c.Response().Header().Set("Content-Transfer-Encoding", "binary")
	c.Response().Header().Set("Expires", "0")
	result,err := file.WriteTo(c.Response())

	//result, err := t.TransUsecase.CountSuccess(ctx)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}
func (p *balanceHistoryHandler) GetBalanceHistory(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}
	merchantId := c.QueryParam("merchant_id")
	status := c.QueryParam("status")
	qpage := c.QueryParam("page")
	qperPage := c.QueryParam("size")

	month := c.QueryParam("month")
	year := c.QueryParam("year")
	isAdmin := c.QueryParam("isAdmin")
	var admin bool
	if isAdmin != ""{
		admin = true
	}else {
		admin = false
	}
	var limit *int
	var page = 1
	var offset *int
	if qpage != "" && qperPage != "" {
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

	res, err := p.balanceHistoryUsecase.List(ctx, merchantId, status, page, limit, offset, month, year,token,admin)
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
	//t := new(models.NewBalanceHistoryCommand)
	//if err := c.Bind(t); err != nil {
	//	return c.JSON(http.StatusBadRequest, models.ErrBadParamInput)
	//}
	//
	//if ok, err := isRequestValid(t); !ok {
	//	return c.JSON(http.StatusBadRequest, err.Error())
	//}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	status, _ := strconv.Atoi(c.FormValue("status"))
	amount, _ := strconv.ParseFloat(c.FormValue("amount"), 64)
	t := models.NewBalanceHistoryCommand{
		Id:            c.FormValue("id"),
		Status:        status,
		Amount:        amount,
		DateOfPayment: c.FormValue("date_of_payment"),
		Remarks:       c.FormValue("remarks"),
	}
	res, err := p.balanceHistoryUsecase.Create(ctx, t, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}
func (p *balanceHistoryHandler) Confirm(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}
	t := new(models.NewBalanceHistoryConfirmCommand)
	if err := c.Bind(t); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrBadParamInput)
	}
	//
	//if ok, err := isRequestValid(t); !ok {
	//	return c.JSON(http.StatusBadRequest, err.Error())
	//}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := p.balanceHistoryUsecase.ConfirmWithdraw(ctx, *t, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}
func (p *balanceHistoryHandler) UpdateAmount(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}
	t := new(models.NewBalanceHistoryAmountCommand)
	if err := c.Bind(t); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrBadParamInput)
	}

	//if ok, err := isRequestValid(t); !ok {
	//	return c.JSON(http.StatusBadRequest, err.Error())
	//}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	res, err := p.balanceHistoryUsecase.UpdateAmount(ctx, *t, token)
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
