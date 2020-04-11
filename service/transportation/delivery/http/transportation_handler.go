package http

import (
	"context"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/service/transportation"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
	"net/http"
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
	//e.PUT("/transportations/:id", handler.Updatetransportation)
	//e.GET("service/special-transportation", handler.GetAlltransportation)
	//e.GET("service/special-transportation/:code", handler.GettransportationByCode)
	//e.DELETE("/transportations/:id", handler.Delete)
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
	response,error := a.transportationUsecase.PublishTransportation(ctx,transporationsCommand,token)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, response)
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
