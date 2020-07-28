package http

import (
	"context"
	"github.com/misc/version_app"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type ResponseError struct {
	Message string `json:"message"`
}

type VersionAPPHandler struct {
	VersionAPPUsecase version_app.Usecase
}

func NewVersionAPPHandler(e *echo.Echo, versionAPPUsecase version_app.Usecase) {
	handler := &VersionAPPHandler{
		VersionAPPUsecase: versionAPPUsecase,
	}
	e.GET("version", handler.Get)
}

func (a *VersionAPPHandler) Get(c echo.Context) error {
	appType := c.QueryParam("type")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}


	typeApp, _ := strconv.Atoi(appType)

	res, err := a.VersionAPPUsecase.GetAllVersion(ctx,typeApp)
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
	//case models.ErrUnAuthorize:
	//	return http.StatusUnauthorized
	//case models.ErrConflict:
	//	return http.StatusBadRequest
	//case models.ErrBadParamInput:
	//	return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
