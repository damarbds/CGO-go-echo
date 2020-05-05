package http

import (
	"context"
	"github.com/auth/identityserver"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo"
	"github.com/models"
	"github.com/service/facilities"
	"github.com/sirupsen/logrus"
)

type ResponseError struct {
	Message string `json:"message"`
}

type facilityHandler struct {
	isUsecase identityserver.Usecase
	facilityUsecase facilities.Usecase
}

func NewFacilityHandler(e *echo.Echo, us facilities.Usecase,isUsecase identityserver.Usecase) {
	handler := &facilityHandler{
		facilityUsecase: us,
		isUsecase:isUsecase,
	}
	e.POST("master/facilities", handler.CreateFacilities)
	e.PUT("master/facilities/:id", handler.UpdateFacilities)
	e.GET("master/facilities", handler.GetAllFacilities)
	e.GET("master/facilities/:id", handler.GetDetailFacilitiesID)
	e.DELETE("master/facilities/:id", handler.DeleteFacilities)
	e.GET("service/facilities", handler.List)
}

func (f *facilityHandler) List(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	list, err := f.facilityUsecase.List(ctx)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}
func (a *facilityHandler) DeleteFacilities(c echo.Context) error {
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
	facilitiesId,_:= strconv.Atoi(id)
	result, err := a.facilityUsecase.Delete(ctx, facilitiesId, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

// GetByID will get article by given id
func (a *facilityHandler) GetAllFacilities(c echo.Context) error {
	qpage := c.QueryParam("page")
	qsize := c.QueryParam("size")

	var limit = 20
	var page = 1
	var offset = 0

	page, _ = strconv.Atoi(qpage)
	limit, _ = strconv.Atoi(qsize)
	offset = (page - 1) * limit
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	art, err := a.facilityUsecase.GetAll(ctx, page,limit,offset)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, art)

}
// Store will store the user by given request body
func (a *facilityHandler) CreateFacilities(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	filupload, image, _ := c.Request().FormFile("facilities_icon")
	dir, err := os.Getwd()
	if err != nil {
		return models.ErrInternalServerError
	}
	var imagePath string
	if filupload != nil {
		fileLocation := filepath.Join(dir, "files", image.Filename)
		targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			os.MkdirAll(filepath.Join(dir, "files"), os.ModePerm)
			return models.ErrInternalServerError
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, filupload); err != nil {
			return models.ErrInternalServerError
		}

		//w.Write([]byte("done"))
		imagePat, _ := a.isUsecase.UploadFileToBlob(fileLocation, "Master/Facilities")
		imagePath = imagePat
		targetFile.Close()
		errRemove := os.Remove(fileLocation)
		if errRemove != nil {
			return models.ErrInternalServerError
		}
	}

	id, _ := strconv.Atoi(c.FormValue("id"))
	isNumerable, _ := strconv.Atoi(c.FormValue("isNumerable"))
	facilitesCommand := models.NewCommandFacilities{
		Id:           id,
		FacilityName: c.FormValue("facilities_name"),
		FacilityIcon: imagePath,
		IsNumerable:  isNumerable,
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	facilities,error := a.facilityUsecase.Create(ctx,&facilitesCommand,token)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, facilities)
}

func (a *facilityHandler) UpdateFacilities(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	filupload, image, _ := c.Request().FormFile("facilities_icon")
	dir, err := os.Getwd()
	if err != nil {
		return models.ErrInternalServerError
	}
	var imagePath string
	if filupload != nil {
		fileLocation := filepath.Join(dir, "files", image.Filename)
		targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			os.MkdirAll(filepath.Join(dir, "files"), os.ModePerm)
			return models.ErrInternalServerError
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, filupload); err != nil {
			return models.ErrInternalServerError
		}

		//w.Write([]byte("done"))
		imagePat, _ := a.isUsecase.UploadFileToBlob(fileLocation, "Master/Facilities")
		imagePath = imagePat
		targetFile.Close()
		errRemove := os.Remove(fileLocation)
		if errRemove != nil {
			return models.ErrInternalServerError
		}
	}

	id, _ := strconv.Atoi(c.FormValue("id"))
	isNumerable, _ := strconv.Atoi(c.FormValue("isNumerable"))
	facilitesCommand := models.NewCommandFacilities{
		Id:           id,
		FacilityName: c.FormValue("facilities_name"),
		FacilityIcon: imagePath,
		IsNumerable:  isNumerable,
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	facilities,error := a.facilityUsecase.Create(ctx,&facilitesCommand,token)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, facilities)
}

func (a *facilityHandler) GetDetailFacilitiesID(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	facilitiesID ,_:= strconv.Atoi(id)
	result, err := a.facilityUsecase.GetById(ctx,facilitiesID)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
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
