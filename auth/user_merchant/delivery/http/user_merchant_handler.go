package http

import (
	"context"
	"github.com/auth/user_merchant"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/models"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// merchantHandler  represent the httphandler for merchant
type userMerchantHandler struct {
	userMerchantUsecase user_merchant.Usecase
}

// NewmerchantHandler will initialize the merchants/ resources endpoint
func NewuserMerchantHandler(e *echo.Echo, us user_merchant.Usecase) {
	handler := &userMerchantHandler{
		userMerchantUsecase: us,
	}
	e.POST("/user-merchants", handler.CreateUserMerchant)
	e.PUT("/user-merchants/:id", handler.UpdateUserMerchant)
	e.GET("/user-merchants", handler.List)
	e.DELETE("/user-merchants/:id", handler.Delete)
	e.GET("/user-merchants/:id", handler.GetDetailID)
	e.GET("/roles-merchants", handler.GetRoles)
	e.POST("/assign-roles-merchants", handler.AssignRoles)
	//e.GET("/merchants/:id", handler.GetByID)
	//e.DELETE("/merchants/:id", handler.Delete)
}
func (a *userMerchantHandler) AssignRoles(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}
	var isAdmin bool
	currentUser := c.QueryParam("isAdmin")
	if currentUser != ""{
		isAdmin = true
	}else {
		isAdmin = false
	}
	cp := new(models.NewCommandAssignRoleUserMerchant)
	if err := c.Bind(cp); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrBadParamInput)
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := a.userMerchantUsecase.AssignRoles(ctx, token,isAdmin,cp)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}
func (a *userMerchantHandler) GetRoles(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}
	var isAdmin bool
	currentUser := c.QueryParam("isAdmin")
	if currentUser != ""{
		isAdmin = true
	}else {
		isAdmin = false
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := a.userMerchantUsecase.GetRoles(ctx, token,isAdmin)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}
func (a *userMerchantHandler) GetDetailID(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := a.userMerchantUsecase.GetUserDetailById(ctx, id,token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (a *userMerchantHandler) Delete(c echo.Context) error {
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

	result, err := a.userMerchantUsecase.Delete(ctx, id, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (a *userMerchantHandler) List(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	qpage := c.QueryParam("page")
	qperPage := c.QueryParam("size")
	search := c.QueryParam("search")
	merchantId := c.QueryParam("merchant_id")
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
	if merchantId != ""{

		result, err := a.userMerchantUsecase.GetUserByMerchantId(ctx, merchantId,token)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, result)
	}else {
		result, err := a.userMerchantUsecase.List(ctx, page, limit, offset, token,search)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, result)
	}
	return nil
}

func isRequestValid(m *models.NewCommandUserMerchant) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store will store the merchant by given request body
func (a *userMerchantHandler) CreateUserMerchant(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}
	merchantCommand := models.NewCommandUserMerchant{
		Id:               c.FormValue("id"),
		FullName:     c.FormValue("full_name"),
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
		PhoneNumber : c.FormValue("phone_number"),
		MerchantId:c.FormValue("merchant_id"),
	}
	if ok, err := isRequestValid(&merchantCommand); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	response,error := a.userMerchantUsecase.Create(ctx, &merchantCommand, token)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, response)
}

func (a *userMerchantHandler) UpdateUserMerchant(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}
	var isAdmin bool
	currentUser := c.FormValue("isAdmin")
	if currentUser != ""{
		isAdmin = true
	}else {
		isAdmin = false
	}
	merchantCommand := models.NewCommandUserMerchant{
		Id:               c.FormValue("id"),
		FullName:     c.FormValue("full_name"),
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
		PhoneNumber : c.FormValue("phone_number"),
		MerchantId:c.FormValue("merchant_id"),
	}
	if ok, err := isRequestValid(&merchantCommand); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := a.userMerchantUsecase.Update(ctx, &merchantCommand, isAdmin,token)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, merchantCommand)
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
