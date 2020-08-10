package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/misc/notif"
	"github.com/models"
	"github.com/sirupsen/logrus"
)

type ResponseError struct {
	Message string `json:"message"`
}

type NotifHandler struct {
	NotifUsecase notif.Usecase
}

func NewNotifHandler(e *echo.Echo, us notif.Usecase) {
	handler := &NotifHandler{
		NotifUsecase: us,
	}
	e.GET("misc/notif", handler.Get)
	e.POST("misc/push-notif", handler.PushNotifFCM)
}
func (a *NotifHandler) PushNotifFCM(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	var pushNotif models.FCMPushNotif
	err := c.Bind(&pushNotif)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	response, error := a.NotifUsecase.FCMPushNotification(ctx, pushNotif)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, response)
}
func (a *NotifHandler) Get(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}
	qpage := c.QueryParam("page")
	qperPage := c.QueryParam("size")

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

	res, err := a.NotifUsecase.GetByMerchantID(ctx, token,page,limit,offset)
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
	//case models.ErrInternalServerError:
	//	return http.StatusInternalServerError
	//case models.ErrNotFound:
	//	return http.StatusNotFound
	case models.ErrUnAuthorize:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
