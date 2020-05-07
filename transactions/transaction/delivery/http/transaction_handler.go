package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	excelize "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"github.com/transactions/transaction"
)

type ResponseError struct {
	Message string `json:"message"`
}

type transactionHandler struct {
	TransUsecase transaction.Usecase
}

func NewTransactionHandler(e *echo.Echo, us transaction.Usecase) {
	handler := &transactionHandler{
		TransUsecase: us,
	}
	e.GET("/transaction/export-excel", handler.ExportExcel)
	e.GET("/transaction/count-success", handler.CountSuccess)
	e.GET("/transaction", handler.List)
	e.GET("/transaction/count-this-month", handler.CountThisMonth)
}

func (t *transactionHandler) CountThisMonth(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := t.TransUsecase.CountThisMonth(ctx)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (t *transactionHandler) List(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	qpage := c.QueryParam("page")
	qperPage := c.QueryParam("size")
	qStatus := c.QueryParam("status")
	qSearch := c.QueryParam("search")
	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")

	isAdmin := c.QueryParam("isAdmin")
	var admin bool
	if isAdmin != ""{
		admin = true
	}else {
		admin = false
	}

	isTrans := c.QueryParam("isTransportation")
	var trans bool
	if isTrans != ""{
		trans = true
	}else {
		trans = false
	}

	isExp := c.QueryParam("isExperience")
	var exp bool
	if isExp != ""{
		exp = true
	}else {
		exp = false
	}
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
	if qpage != "" && qperPage != ""{
		result, err := t.TransUsecase.List(ctx, startDate, endDate, qSearch, qStatus, &page, &limit, &offset,token,admin,trans,exp)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, result)
	}else {
		result, err := t.TransUsecase.List(ctx, startDate, endDate, qSearch, qStatus, nil, nil, nil,token,admin,trans,exp)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, result)
	}
}

func (t *transactionHandler) CountSuccess(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := t.TransUsecase.CountSuccess(ctx)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}
func PrepareAndReturnExcel(transaction *models.TransactionWithPagination) *excelize.File {
	f := excelize.NewFile()

	f.SetCellValue("Sheet1", "A1", "Item")
	f.SetCellValue("Sheet1", "B1", "Type")
	f.SetCellValue("Sheet1", "C1", "Order Id")
	f.SetCellValue("Sheet1", "D1", "Booking Date")
	f.SetCellValue("Sheet1", "E1", "Check In Date")
	f.SetCellValue("Sheet1", "F1", "Booked By")
	f.SetCellValue("Sheet1", "G1", "Guest")
	for index,element := range transaction.Data{
		if index == 0 {
			index = index + 3
		}else {
			index = index + 2
		}
		var expType string
		var bookedBy string
		if len(element.BookedBy) != 0 {
			bookedBy = element.BookedBy[0].FullName
		}
		if len(element.ExpType) != 0 {
			expTypeConvert, _ := json.Marshal(element.ExpType)
			expType = string(expTypeConvert)
		}
		f.SetCellValue("Sheet1", "A" + strconv.Itoa(index), element.ExpTitle)
		f.SetCellValue("Sheet1", "B" + strconv.Itoa(index), expType)
		f.SetCellValue("Sheet1", "C"+ strconv.Itoa(index), *element.OrderId)
		f.SetCellValue("Sheet1", "D"+ strconv.Itoa(index), element.BookingDate.Format("2006-01-02"))
		f.SetCellValue("Sheet1", "E"+ strconv.Itoa(index), element.CheckInDate.Format("2006-01-02"))
		f.SetCellValue("Sheet1", "F"+ strconv.Itoa(index), bookedBy)
		f.SetCellValue("Sheet1", "G"+ strconv.Itoa(index), strconv.Itoa(element.Guest))
	}
	return f
}
func (t *transactionHandler) ExportExcel(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}
	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")


	isAdmin := c.QueryParam("isAdmin")
	var admin bool
	if isAdmin != ""{
		admin = true
	}else {
		admin = false
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	getresult, err := t.TransUsecase.List(ctx, startDate, endDate, "", "", nil, nil, nil,token,admin,false,false)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	file := PrepareAndReturnExcel(getresult)
	// Set the headers necessary to get browsers to interpret the downloadable file
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	c.Response().Header().Set("Content-Disposition", "attachment;filename=transaction.xlsx")
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
