package http

import (
	"github.com/auth/identityserver"
	"github.com/service/cpc"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	validator "gopkg.in/go-playground/validator.v9"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// promoHandler  represent the httphandler for promo
type cPCHandler struct {
	isUsecase identityserver.Usecase
	cPCUsecase cpc.Usecase
}

// NewpromoHandler will initialize the promos/ resources endpoint
func NewCPCHandler(e *echo.Echo, cPCUsecase cpc.Usecase,is identityserver.Usecase) {
	handler := &cPCHandler{
		isUsecase:is,
		cPCUsecase: cPCUsecase,
	}
	e.POST("master/city", handler.CreateCity)
	e.PUT("master/city/:id", handler.UpdateCity)
	e.GET("master/city", handler.ListCity)
	e.GET("master/city/:id", handler.GetDetailCityID)
	e.DELETE("master/city/:id", handler.DeleteCity)
	e.POST("master/province", handler.CreateProvince)
	e.PUT("master/province/:id", handler.UpdateProvince)
	e.GET("master/province", handler.ListProvince)
	e.GET("master/province/:id", handler.GetDetailProvinceID)
	e.DELETE("master/province/:id", handler.DeleteProvince)
	e.POST("master/country", handler.CreateCountry)
	e.PUT("master/country/:id", handler.UpdateCountry)
	e.GET("master/country", handler.ListCountry)
	e.GET("master/country/:id", handler.GetDetailCountryID)
	e.DELETE("master/country/:id", handler.DeleteCountry)
}

func isRequestValid(m *models.NewCommandPromo) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (a *cPCHandler) DeleteCity(c echo.Context) error {
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
	cityId,_:= strconv.Atoi(id)
	result, err := a.cPCUsecase.DeleteCity(ctx, cityId, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

// GetByID will get article by given id
func (a *cPCHandler) ListCity(c echo.Context) error {
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
		art, err := a.cPCUsecase.GetAllCity(ctx, page,limit,offset)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, art)

}
// Store will store the user by given request body
func (a *cPCHandler) CreateCity(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	filupload, image, _ := c.Request().FormFile("city_image")
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
		imagePat, _ := a.isUsecase.UploadFileToBlob(fileLocation, "Master/City")
		imagePath = imagePat
		targetFile.Close()
		errRemove := os.Remove(fileLocation)
		if errRemove != nil {
			return models.ErrInternalServerError
		}
	}

	id, _ := strconv.Atoi(c.FormValue("id"))
	provinceId, _ := strconv.Atoi(c.FormValue("province_id"))
	coverPhotos := make([]models.CoverPhotosObj,1)
	if imagePath != ""{
		coverPhoto := models.CoverPhotosObj{
			Original:  imagePath,
			Thumbnail: "",
		}
		coverPhotos[0] = coverPhoto
	}
	cityCommand := models.NewCommandCity{
		Id:         id,
		CityName:   c.FormValue("city_name"),
		CityDesc:   c.FormValue("city_desc"),
		CityPhotos: coverPhotos,
		ProvinceId: provinceId,
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	city,error := a.cPCUsecase.CreateCity(ctx, &cityCommand, token)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, city)
}

func (a *cPCHandler) UpdateCity(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	filupload, image, _ := c.Request().FormFile("city_image")
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
		imagePat, _ := a.isUsecase.UploadFileToBlob(fileLocation, "Master/City")
		imagePath = imagePat
		targetFile.Close()
		errRemove := os.Remove(fileLocation)
		if errRemove != nil {
			return models.ErrInternalServerError
		}
	}

	id, _ := strconv.Atoi(c.FormValue("id"))
	provinceId, _ := strconv.Atoi(c.FormValue("province_id"))
	coverPhotos := make([]models.CoverPhotosObj,0)
	if imagePath != ""{
		coverPhoto := models.CoverPhotosObj{
			Original:  imagePath,
			Thumbnail: "",
		}
		coverPhotos[0] = coverPhoto
	}
	cityCommand := models.NewCommandCity{
		Id:         id,
		CityName:   c.FormValue("city_name"),
		CityDesc:   c.FormValue("city_desc"),
		CityPhotos: coverPhotos,
		ProvinceId: provinceId,
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	city,error := a.cPCUsecase.UpdateCity(ctx, &cityCommand, token)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, city)
}

func (a *cPCHandler) GetDetailCityID(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	cityId ,_:= strconv.Atoi(id)
	result, err := a.cPCUsecase.GetCityById(ctx,cityId)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}


func (a *cPCHandler) DeleteProvince(c echo.Context) error {
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
	provinceId,_:= strconv.Atoi(id)
	result, err := a.cPCUsecase.DeleteProvince(ctx, provinceId, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

// GetByID will get article by given id
func (a *cPCHandler) ListProvince(c echo.Context) error {
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
	art, err := a.cPCUsecase.GetAllProvince(ctx, page,limit,offset)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, art)

}
// Store will store the user by given request body
func (a *cPCHandler) CreateProvince(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	id, _ := strconv.Atoi(c.FormValue("id"))
	countryId, _ := strconv.Atoi(c.FormValue("country_id"))
	provinceCommand := models.NewCommandProvince{
		Id:           id,
		ProvinceName: c.FormValue("province_name"),
		CountryId:    countryId,
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	province,error := a.cPCUsecase.CreateProvince(ctx, &provinceCommand, token)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, province)
}

func (a *cPCHandler) UpdateProvince(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	id, _ := strconv.Atoi(c.FormValue("id"))
	countryId, _ := strconv.Atoi(c.FormValue("country_id"))
	provinceCommand := models.NewCommandProvince{
		Id:           id,
		ProvinceName: c.FormValue("province_name"),
		CountryId:    countryId,
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	province,error := a.cPCUsecase.UpdateProvince(ctx, &provinceCommand, token)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, province)
}

func (a *cPCHandler) GetDetailProvinceID(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	provinceId ,_:= strconv.Atoi(id)
	result, err := a.cPCUsecase.GetProvinceById(ctx,provinceId)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}


func (a *cPCHandler) DeleteCountry(c echo.Context) error {
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
	countryId,_:= strconv.Atoi(id)
	result, err := a.cPCUsecase.DeleteCountry(ctx, countryId, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

// GetByID will get article by given id
func (a *cPCHandler) ListCountry(c echo.Context) error {
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
	art, err := a.cPCUsecase.GetAllCountry(ctx, page,limit,offset)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, art)

}
// Store will store the user by given request body
func (a *cPCHandler) CreateCountry(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	id, _ := strconv.Atoi(c.FormValue("id"))
	countryCommand := models.NewCommandCountry{
		Id:          id,
		CountryName: c.FormValue("country_name"),
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	country,error := a.cPCUsecase.CreateCountry(ctx, &countryCommand, token)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, country)
}

func (a *cPCHandler) UpdateCountry(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	id, _ := strconv.Atoi(c.FormValue("id"))
	countryCommand := models.NewCommandCountry{
		Id:          id,
		CountryName: c.FormValue("country_name"),
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	country,error := a.cPCUsecase.UpdateCountry(ctx, &countryCommand, token)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, country)
}

func (a *cPCHandler) GetDetailCountryID(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	countryId ,_:= strconv.Atoi(id)
	result, err := a.cPCUsecase.GetCountryById(ctx,countryId)
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
