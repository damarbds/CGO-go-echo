package http

import (
	"context"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
	"net/http"

	"github.com/auth/admin"
	"github.com/models"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// adminHandler  represent the httphandler for admin
type adminHandler struct {
	adminUsecase admin.Usecase
}

// NewadminHandler will initialize the admins/ resources endpoint
func NewadminHandler(e *echo.Echo, us admin.Usecase) {
	handler := &adminHandler{
		adminUsecase: us,
	}
	e.POST("/admin", handler.Createadmin)
	e.PUT("/admin/:id", handler.Updateadmin)
	//e.GET("/admins/:id", handler.GetByID)
	//e.DELETE("/admins/:id", handler.Delete)
}

func isRequestValid(m *models.NewCommandAdmin) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store will store the admin by given request body
func (a *adminHandler) Createadmin(c echo.Context) error {
	//var adminCommand models.NewCommandadmin
	//err := c.Bind(&adminCommand)
	//if err != nil {
	//	return c.JSON(http.StatusUnprocessableEntity, err.Error())
	//}
	adminCommand := models.NewCommandAdmin{
		Id:               c.FormValue("id"),
		Name:     c.FormValue("name"),
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}
	if ok, err := isRequestValid(&adminCommand); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	error := a.adminUsecase.Create(ctx, &adminCommand,"admin")

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, adminCommand)
}

func (a *adminHandler) Updateadmin(c echo.Context) error {
	//var adminCommand models.NewCommandadmin
	//err := c.Bind(&adminCommand)
	//if err != nil {
	//	return c.JSON(http.StatusUnprocessableEntity, err.Error())
	//}
	adminCommand := models.NewCommandAdmin{
		Id:               c.FormValue("id"),
		Name:     c.FormValue("name"),
		Email:     c.FormValue("email"),
		Password: c.FormValue("password"),
	}
	if ok, err := isRequestValid(&adminCommand); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := a.adminUsecase.Update(ctx, &adminCommand,"admin")

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, adminCommand)
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
