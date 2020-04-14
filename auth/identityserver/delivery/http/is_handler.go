package http

import (
	"context"
	"github.com/auth/admin"
	"github.com/auth/identityserver"
	"github.com/auth/merchant"
	"github.com/auth/user"
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
	isUsecase 		identityserver.Usecase
	merchantUsecase merchant.Usecase
	userUsecase     user.Usecase
	adminUsecase    admin.Usecase
}

// NewisHandler will initialize the iss/ resources endpoint
func NewisHandler(e *echo.Echo, m merchant.Usecase, u user.Usecase, a admin.Usecase,is identityserver.Usecase) {
	handler := &isHandler{
		merchantUsecase: m,
		userUsecase:     u,
		adminUsecase:    a,
		isUsecase:is,
	}
	e.GET("/account/info", handler.GetInfo)
	e.POST("/account/login", handler.Login)
	e.POST("/account/request-otp", handler.RequestOTP)
	e.GET("/account/verified-email", handler.VerifiedEmail)
}

func (a *isHandler) RequestOTP(c echo.Context) error {
	var requestOTP models.RequestOTPNumber
	err := c.Bind(&requestOTP)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	requestOTP.PhoneNumber = c.Request().Form.Get("phone_number")
	responseOTP , err:= a.userUsecase.RequestOTP(ctx,requestOTP.PhoneNumber)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, responseOTP)
}
func (a *isHandler) Login(c echo.Context) error {
	var isLogin models.Login
	err := c.Bind(&isLogin)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	isLogin.Email = c.Request().Form.Get("email")
	isLogin.Password = c.Request().Form.Get("password")
	isLogin.Type = c.Request().Form.Get("type")
	isLogin.Scope = c.Request().Form.Get("scope")
	var responseToken *models.GetToken
	if isLogin.Type == "user" {

		token, err := a.userUsecase.Login(ctx, &isLogin)

		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		responseToken = token
	} else if isLogin.Type == "merchant" {

		token, err := a.merchantUsecase.Login(ctx, &isLogin)

		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		responseToken = token
	} else if isLogin.Type == "admin" {

		token, err := a.adminUsecase.Login(ctx, &isLogin)

		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		responseToken = token
	} else {
		return c.JSON(http.StatusBadRequest, "Bad Request")
	}

	return c.JSON(http.StatusOK, responseToken)
}

func (a *isHandler) VerifiedEmail(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	token := c.Request().Header.Get("Authorization")
	typeUser := c.QueryParam("type")
	otpCode := c.QueryParam("otp")
	if typeUser == "user" {
		response, err := a.userUsecase.VerifiedEmail(ctx, token, otpCode)

		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, response)
	} else if typeUser == "merchant" {
		return c.JSON(http.StatusNotFound, "Not Implemented")
	} else {
		return c.JSON(http.StatusBadRequest, "Bad Request")
	}

	return c.JSON(http.StatusBadRequest, "Bad Request")
}

func (a *isHandler) GetInfo(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	token := c.Request().Header.Get("Authorization")
	typeUser := c.QueryParam("type")
	if typeUser == "user" {
		response, err := a.userUsecase.GetUserInfo(ctx, token)

		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, response)
	} else if typeUser == "merchant" {
		response, err := a.merchantUsecase.GetMerchantInfo(ctx, token)

		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, response)
	} else if typeUser == "admin" {
		response, err := a.adminUsecase.GetAdminInfo(ctx, token)

		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, response)
	} else {
		return c.JSON(http.StatusBadRequest, "Bad Request")
	}

	return c.JSON(http.StatusBadRequest, "Bad Request")
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
	case models.ErrUsernamePassword:
		return http.StatusUnauthorized
	case models.ErrInvalidOTP:
		return http.StatusUnauthorized
	case models.ErrNotYetRegister:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
