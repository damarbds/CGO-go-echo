package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/models"
	"github.com/profile/wishlists"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

type WishlistHandler struct {
	WlUsecase wishlists.Usecase
}

func NewWishlistHandler(e *echo.Echo, wus wishlists.Usecase) {
	handler := &WishlistHandler{
		WlUsecase: wus,
	}
	e.POST("/profile/wishlists", handler.Create)
	e.GET("/profile/wishlists", handler.List)
	e.GET("/profile/check-wishlists", handler.CheckWishList)
}

func isRequestValid(m *models.WishlistIn) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (w *WishlistHandler) CheckWishList(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}
	expId := c.QueryParam("exp_id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := w.WlUsecase.List(ctx, token, 1, 1, 0, expId)
	var response bool
	if len(res.Data) != 0 {
		response = true
	} else {
		response = false
	}
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, response)
}

func (w *WishlistHandler) List(c echo.Context) error {
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

	res, err := w.WlUsecase.List(ctx, token, page, limit, offset, "")
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (w *WishlistHandler) Create(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

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

	res, err := w.WlUsecase.Insert(ctx, wi, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	response := map[string]interface{}{
		"status":      http.StatusOK,
		"message":     "Create Wishlist Succeeds",
		"wishlist_id": res,
	}
	if wi.IsDeleted == true {
		response = map[string]interface{}{
			"status":      http.StatusOK,
			"message":     "Deleted Wishlist Succeeds",
			"wishlist_id": res,
		}
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
