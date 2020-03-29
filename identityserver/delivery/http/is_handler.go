package http

import (
	"context"
	"github.com/merchant"
	"github.com/user"
	"net/http"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"

	"github.com/models"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// isHandler  represent the httphandler for is
type isHandler struct {
	merchantUsecase merchant.Usecase
	userUsecase 	user.Usecase
}

// NewisHandler will initialize the iss/ resources endpoint
func NewisHandler(e *echo.Echo, m merchant.Usecase,u user.Usecase) {
	handler := &isHandler{
		merchantUsecase:m,
		userUsecase:u,
	}
	e.GET("/account/info", handler.GetInfo)
	e.POST("/account/login", handler.Login)
}

func (a *isHandler) Login(c echo.Context) error {
	var isLogin models.Login
	err := c.Bind(&isLogin)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	isLogin.Email = c.Request().Form.Get("email")
	isLogin.Password = c.Request().Form.Get("password")
	isLogin.Type = c.Request().Form.Get("type")
	var responseToken *models.GetToken
	if isLogin.Type == "user" {

		token, err := a.userUsecase.Login(ctx, &isLogin)

		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		responseToken = token
	}else if isLogin.Type == "merchant" {

		token, err := a.merchantUsecase.Login(ctx, &isLogin)

		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		responseToken = token
	}else {
		return c.JSON(http.StatusBadRequest,"Bad Request")
	}

	return c.JSON(http.StatusCreated, responseToken)
}

func (a *isHandler) GetInfo(c echo.Context) error {

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	token := c.Request().Header.Get("Authorization")
	typeUser := c.QueryParam("type")
	if typeUser == "user"{
		response, err := a.userUsecase.GetUserInfo(ctx, token)

		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}

		return c.JSON(http.StatusCreated, response)
	} else if typeUser == "merchant"{
		response, err := a.merchantUsecase.GetMerchantInfo(ctx, token)

		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusCreated, response)
	}else {
		return c.JSON(http.StatusBadRequest,"Bad Request")
	}

	return c.JSON(http.StatusBadRequest,"Bad Request")
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
	case models.ErrConflict:
		return http.StatusConflict
	case models.ErrUnAuthorize:
		return http.StatusUnauthorized
	case models.ErrBadParamInput:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
