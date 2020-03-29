package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/models"
	"github.com/user"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// userHandler  represent the httphandler for user
type userHandler struct {
	userUsecase user.Usecase
}

// NewuserHandler will initialize the users/ resources endpoint
func NewuserHandler(e *echo.Echo, us user.Usecase) {
	handler := &userHandler{
		userUsecase: us,
	}
	e.POST("/users", handler.CreateUser)
	e.PUT("/users/:id", handler.UpdateUser)
}

func isRequestValid(m *models.NewCommandUser) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store will store the user by given request body
func (a *userHandler) CreateUser(c echo.Context) error {
	test:= time.Now()
	fmt.Println(test)
	var userCommand models.NewCommandUser
	err := c.Bind(&userCommand)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&userCommand); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	err = a.userUsecase.Create(ctx, &userCommand, "admin")

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, userCommand)
}

func (a *userHandler) UpdateUser(c echo.Context) error {
	var userCommand models.NewCommandUser
	err := c.Bind(&userCommand)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&userCommand); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = a.userUsecase.Update(ctx, &userCommand, "admin")

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, userCommand)
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
		return http.StatusConflict
	case models.ErrBadParamInput:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
