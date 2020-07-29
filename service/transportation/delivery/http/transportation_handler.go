package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/models"
	"github.com/service/transportation"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// transportationHandler  represent the httphandler for transportation
type transportationHandler struct {
	transportationUsecase transportation.Usecase
}

// NewtransportationHandler will initialize the transportations/ resources endpoint
func NewtransportationHandler(e *echo.Echo, us transportation.Usecase) {
	handler := &transportationHandler{
		transportationUsecase: us,
	}
	e.POST("service/transportation/create", handler.CreateTransportation)
	e.GET("/service/transportation/time-options", handler.TimeOptions)
	e.GET("/service/transportation/filter-search", handler.FilterSearchTrans)
	e.GET("/service/transportation/:id", handler.GetByID)
	e.GET("service/transportation/published-count", handler.GetPublishedTransCount)
	//e.PUT("/transportations/:id", handler.Updatetransportation)
	//e.GET("service/special-transportation", handler.GetAlltransportation)
	//e.GET("service/special-transportation/:code", handler.GettransportationByCode)
	//e.DELETE("/transportations/:id", handler.Delete)
	e.GET("/service/transportation/master-data-transport", handler.GetAllTransport)
}
func (a *transportationHandler) GetPublishedTransCount(c echo.Context) error {
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

	result, err := a.transportationUsecase.GetPublishedTransCount(ctx, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}
func (a *transportationHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	//if err != nil {
	//	return c.JSON(http.StatusNotFound, models.ErrNotFound.Error())
	//}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	art, err := a.transportationUsecase.GetDetail(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, art)
}
func (t *transportationHandler) FilterSearchTrans(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")
	isMerchant := c.QueryParam("isMerchant")
	returnTransId := c.QueryParam("return_trans_id")
	harborSourceId := c.QueryParam("harbor_source_id")
	harborDestId := c.QueryParam("harbor_dest_id")
	guest := c.QueryParam("guest")
	depDate := c.QueryParam("departure_date")
	class := c.QueryParam("class")
	isReturn := c.QueryParam("isReturn")
	sortBy := c.QueryParam("sortBy")
	depTimeOptions := c.QueryParam("dep_timeoption_id")
	arrTimeOptions := c.QueryParam("arr_timeoption_id")
	notReturn := c.QueryParam("not_return")
	qpage := c.QueryParam("page")
	qperPage := c.QueryParam("size")
	currency := c.QueryParam("currency")
	search := c.QueryParam("search")
	status := c.QueryParam("status")

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

	needMerchantAuth := false
	if isMerchant != "" {
		needMerchantAuth = true
	}
	var isReturnTrip bool
	if isReturn == "0" {
		isReturnTrip = false
	} else if isReturn == "1" {
		isReturnTrip = true
	}
	guestTrip, _ := strconv.Atoi(guest)
	depTimeOp, _ := strconv.Atoi(depTimeOptions)
	arrTimeOp, _ := strconv.Atoi(arrTimeOptions)


	results, err := t.transportationUsecase.FilterSearchTrans(ctx, needMerchantAuth, token, search, status, sortBy,
		harborSourceId, harborDestId, depDate, class, isReturnTrip, depTimeOp, arrTimeOp, guestTrip, page, limit,
		offset,returnTransId,notReturn,currency)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}

func (t *transportationHandler) TimeOptions(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := t.transportationUsecase.TimeOptions(ctx)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func isRequestValid(m *models.NewCommandTransportation) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (a *transportationHandler) CreateTransportation(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}
	var transporationsCommand models.NewCommandTransportation
	err := c.Bind(&transporationsCommand)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	if ok, err := isRequestValid(&transporationsCommand); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	response, error := a.transportationUsecase.PublishTransportation(ctx, transporationsCommand, token)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, response)
}
func (t *transportationHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := t.transportationUsecase.TimeOptions(ctx)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (a *transportationHandler) GetAllTransport(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := a.transportationUsecase.GetAllTransport(ctx)
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
