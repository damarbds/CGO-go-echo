package http

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/auth/identityserver"

	"github.com/labstack/echo"
	"github.com/models"
	"github.com/service/harbors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// harborsHandler  represent the httphandler for harbors
type HarborsHandler struct {
	IsUsecase      identityserver.Usecase
	HarborsUsecase harbors.Usecase
}

// NewharborsHandler will initialize the harborss/ resources endpoint
func NewharborsHandler(e *echo.Echo, us harbors.Usecase, isUsecase identityserver.Usecase) {
	handler := &HarborsHandler{
		HarborsUsecase: us,
		IsUsecase:      isUsecase,
	}
	e.GET("service/exp-destination", handler.GetAllHarbors)
	e.POST("master/harbors", handler.CreateHarbors)
	e.PUT("master/harbors/:id", handler.UpdateHarbors)
	e.GET("master/harbors", handler.ListHarbors)
	e.GET("master/harbors/:id", handler.GetDetailHarborsID)
	e.DELETE("master/harbors/:id", handler.DeleteHarbors)
}
// GetByID will get article by given id
func (a *HarborsHandler) GetAllHarbors(c echo.Context) error {
	qpage := c.QueryParam("page")
	qsize := c.QueryParam("size")
	search := c.QueryParam("search")
	harborsType := c.QueryParam("harbors_type")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	if qpage != "" && qsize != "" {
		page, _ := strconv.Atoi(qpage)
		size, _ := strconv.Atoi(qsize)
		art, err := a.HarborsUsecase.GetAllWithJoinCPC(ctx, &size, &page, search, harborsType)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, art)
	} else {
		art, err := a.HarborsUsecase.GetAllWithJoinCPC(ctx, nil, nil, search, harborsType)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, art)
	}

	return nil
}
func (a *HarborsHandler) DeleteHarbors(c echo.Context) error {
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
	result, err := a.HarborsUsecase.Delete(ctx, id, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

// GetByID will get article by given id
func (a *HarborsHandler) ListHarbors(c echo.Context) error {
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
	art, err := a.HarborsUsecase.GetAll(ctx, page, limit, offset)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, art)

}

// Store will store the user by given request body
func (a *HarborsHandler) CreateHarbors(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	filupload, image, _ := c.Request().FormFile("harbors_image")
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
		imagePat, _ := a.IsUsecase.UploadFileToBlob(fileLocation, "Master/Harbors")
		imagePath = imagePat
		targetFile.Close()
		errRemove := os.Remove(fileLocation)
		if errRemove != nil {
			return models.ErrInternalServerError
		}
	}

	cityId, _ := strconv.Atoi(c.FormValue("city_id"))
	harborsLongitude, _ := strconv.ParseFloat(c.FormValue("harbors_longitude"), 64)
	harborsLatitude, _ := strconv.ParseFloat(c.FormValue("harbors_latitude"), 64)
	harborsType, _ := strconv.Atoi(c.FormValue("harbors_type"))
	harborsCommand := models.NewCommandHarbors{
		Id:               c.FormValue("id"),
		HarborsName:      c.FormValue("harbors_name"),
		HarborsLongitude: harborsLongitude,
		HarborsLatitude:  harborsLatitude,
		HarborsImage:     imagePath,
		CityId:           cityId,
		HarborsType:      harborsType,
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	harbors, error := a.HarborsUsecase.Create(ctx, &harborsCommand, token)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, harbors)
}

func (a *HarborsHandler) UpdateHarbors(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	filupload, image, _ := c.Request().FormFile("harbors_image")
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
		imagePat, _ := a.IsUsecase.UploadFileToBlob(fileLocation, "Master/Harbors")
		imagePath = imagePat
		targetFile.Close()
		errRemove := os.Remove(fileLocation)
		if errRemove != nil {
			return models.ErrInternalServerError
		}
	}

	cityId, _ := strconv.Atoi(c.FormValue("city_id"))
	harborsLongitude, _ := strconv.ParseFloat(c.FormValue("harbors_longitude"), 64)
	harborsLatitude, _ := strconv.ParseFloat(c.FormValue("harbors_latitude"), 64)
	harborsType, _ := strconv.Atoi(c.FormValue("harbors_type"))
	harborsCommand := models.NewCommandHarbors{
		Id:               c.FormValue("id"),
		HarborsName:      c.FormValue("harbors_name"),
		HarborsLongitude: harborsLongitude,
		HarborsLatitude:  harborsLatitude,
		HarborsImage:     imagePath,
		CityId:           cityId,
		HarborsType:      harborsType,
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	harbors, error := a.HarborsUsecase.Update(ctx, &harborsCommand, token)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, harbors)
}

func (a *HarborsHandler) GetDetailHarborsID(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	result, err := a.HarborsUsecase.GetById(ctx, id)
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
