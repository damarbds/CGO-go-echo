package http

import (
	"context"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/profile/wishlists"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

type wishlistHandler struct {
	wlUsecase wishlists.Usecase
}

func NewWishlistHandler(e *echo.Echo, wus wishlists.Usecase) {
	handler := &wishlistHandler{
		wlUsecase: wus,
	}
	e.POST("/profile/wishlists", handler.CreatePayment)
}

func isRequestValid(m *models.WishlistIn) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (w *wishlistHandler) CreatePayment(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	wi := new(models.WishlistIn)
	if err := c.Bind(wi); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrBadParamInput)
	}

	if ok, err := isRequestValid(wi); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := w.wlUsecase.Insert(ctx, wi, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	response := map[string]interface{}{
		"status": http.StatusOK,
		"message": "Create Wishlist Succeeds",
		"wishlist_id": res,
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
