package http

import (
	"context"
	"net/http"

	"github.com/auth/admin"
	"github.com/auth/identityserver"
	"github.com/auth/merchant"
	"github.com/auth/user"
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
	isUsecase       identityserver.Usecase
	merchantUsecase merchant.Usecase
	userUsecase     user.Usecase
	adminUsecase    admin.Usecase
}

// NewisHandler will initialize the iss/ resources endpoint
func NewisHandler(e *echo.Echo, m merchant.Usecase, u user.Usecase, a admin.Usecase, is identityserver.Usecase) {
	handler := &isHandler{
		merchantUsecase: m,
		userUsecase:     u,
		adminUsecase:    a,
		isUsecase:       is,
	}
	e.POST("/account/auto-login", handler.AutoLogin)
	e.GET("/account/info", handler.GetInfo)
	e.POST("/account/login", handler.Login)
	e.POST("/account/refresh-token", handler.RefreshToken)
	e.POST("/account/request-otp", handler.RequestOTP)
	e.POST("/account/request-otp-tmp", handler.RequestOTPTmp)
	e.GET("/account/verified-email", handler.VerifiedEmail)
	e.GET("/account/callback", handler.CallBack)
	e.GET("/account/forgot-password", handler.ForgotPassword)
	e.GET("/account/change-password",handler.ChangePassword)
}
func (a *isHandler) AutoLogin(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	token := c.Request().Header.Get("Authorization")
	var isLogin models.AutoLogin
	err := c.Bind(&isLogin)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	isLogin.MerchantId = c.Request().Form.Get("merchant_id")

	responseToken, err := a.merchantUsecase.AutoLoginByCMSAdmin(ctx, isLogin.MerchantId, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, responseToken)
}
func (a *isHandler) ChangePassword(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	password := c.QueryParam("password")
	responseOTP, err := a.userUsecase.ChangePassword(ctx,token,password)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, responseOTP)
}
func (a *isHandler) ForgotPassword(c echo.Context) error {

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	email := c.QueryParam("email")
	getUser,err := a.userUsecase.GetUserByEmail(ctx,email)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrNotFound)
	}
	getUserDetail , err := a.isUsecase.GetDetailUserById(getUser.Id,"","true")
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrNotFound)
	}
	login, err := a.isUsecase.GetToken(email,getUserDetail.Password,"")
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrNotFound)
	}
	responseOTP, err := a.isUsecase.ForgotPassword(email,login.RefreshToken)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, responseOTP)
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
	responseOTP, err := a.userUsecase.RequestOTP(ctx, requestOTP.PhoneNumber)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, responseOTP)
}
func (a *isHandler) RequestOTPTmp(c echo.Context) error {
	var requestOTP models.RequestOTPTmpNumber
	err := c.Bind(&requestOTP)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	requestOTP.PhoneNumber = c.Request().Form.Get("phone_number")
	requestOTP.Email = c.Request().Form.Get("email")
	if requestOTP.Email != "" {
		checkEmail, _ := a.userUsecase.GetUserByEmail(ctx, requestOTP.Email)
		if checkEmail != nil {
			return c.JSON(getStatusCode(models.ErrConflict), ResponseError{Message: models.ErrConflict.Error()})
		}
	}
	responseOTP, err := a.isUsecase.RequestOTPTmp(requestOTP.PhoneNumber, requestOTP.Email)
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
	isLogin.XMode = c.Request().Form.Get("x_mode")
	var responseToken *models.GetToken
	if isLogin.Type == "user" {

		token, err := a.userUsecase.Login(ctx, &isLogin)

		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		responseToken = token
	} else if isLogin.Type == "merchant" && isLogin.XMode == "mobile"{

		token, err := a.merchantUsecase.LoginMobile(ctx, &isLogin)

		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, token)
		//responseToken = token
	}else if isLogin.Type == "merchant" {

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
func (a *isHandler) RefreshToken(c echo.Context) error {
	var isLogin models.RefreshTokenLogin
	err := c.Bind(&isLogin)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	isLogin.RefreshToken = c.Request().Form.Get("refresh_token")
	responseToken, err := a.isUsecase.RefreshToken(isLogin.RefreshToken)
	return c.JSON(http.StatusOK, responseToken)
}
func (a *isHandler) VerifiedEmail(c echo.Context) error {
	//c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	//c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	//token := c.Request().Header.Get("Authorization")
	email := c.QueryParam("email")
	otpCode := c.QueryParam("otp")
	verified := models.VerifiedEmail{
		Email:   email,
		CodeOTP: otpCode,
	}
	//if typeUser == "user" {
	response, err := a.isUsecase.VerifiedEmail(&verified)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, response)
	//} else if typeUser == "merchant" {
	//	return c.JSON(http.StatusNotFound, "Not Implemented")
	//} else {
	//	return c.JSON(http.StatusBadRequest, "Bad Request")
	//}

	return c.JSON(http.StatusBadRequest, "Bad Request")
}
func (a *isHandler) CallBack(c echo.Context) error {
	code := c.QueryParam("code")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	token, err := a.userUsecase.LoginByGoogle(ctx, code)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusBadRequest, token)
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
