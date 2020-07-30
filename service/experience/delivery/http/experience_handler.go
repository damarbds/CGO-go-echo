package http

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/service/transportation"

	"github.com/auth/identityserver"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/service/experience"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	validator "gopkg.in/go-playground/validator.v9"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// experienceHandler  represent the httphandler for experience
type experienceHandler struct {
	transportationsUsecase transportation.Usecase
	experienceUsecase      experience.Usecase
	isUsecase              identityserver.Usecase
}

// NewexperienceHandler will initialize the experiences/ resources endpoint
func NewexperienceHandler(e *echo.Echo, t transportation.Usecase, us experience.Usecase, is identityserver.Usecase) {
	handler := &experienceHandler{
		isUsecase:              is,
		experienceUsecase:      us,
		transportationsUsecase: t,
	}
	//e.POST("/experiences", handler.Createexperience)
	e.PUT("service/change-status", handler.UpdateStatusExpORTrans)
	e.POST("media/upload", handler.UploadFile)
	e.POST("service/experience/create", handler.CreateExperiences)
	e.GET("service/experience/:id", handler.GetByID)
	e.GET("service/experience/search", handler.SearchExp)
	e.GET("service/experience/filter-search", handler.FilterSearchExp)
	e.GET("service/experience/get-user-discover-preference", handler.GetUserDiscoverPreference)
	e.GET("service/experience/categories", handler.GetExpTypes)
	e.GET("service/experience/inspirations", handler.GetExpInspirations)
	e.GET("service/experience/categories/:id", handler.GetByCategoryID)
	e.GET("service/experience/success-book-count", handler.GetSuccessBookCount)
	e.GET("service/experience/published-count", handler.GetPublishedExpCount)
	e.GET("service/experience/transaction-count", handler.GetExpTransactionCount)
	e.GET("service/experience/master-data-experience", handler.GetAllExperience)
	//e.DELETE("/experiences/:id", handler.Delete)
}

func isRequestValid(m *models.NewCommandExperience) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (a *experienceHandler) UploadFile(c echo.Context) error {
	err := c.Request().ParseMultipartForm(200000) // grab the multipart form
	if err != nil {
		fmt.Fprintln(c.Response(), err)
		return nil
	}

	formdata := c.Request().MultipartForm

	files := formdata.File["image"]
	//filupload, image, _ := c.Request().FormFile("image")
	folder := c.QueryParam("folder")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	dir, err := os.Getwd()
	if err != nil {
		return models.ErrInternalServerError
	}
	var imagePath []string

	for i, _ := range files { // loop through the files one by one
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			fmt.Fprintln(c.Response(), err)
			return err
		}
		//out, err := os.Create("/tmp/" + files[i].Filename)
		fileLocation := filepath.Join(dir, "files", files[i].Filename)
		targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			os.MkdirAll(filepath.Join(dir, "files"), os.ModePerm)
			return models.ErrInternalServerError
		}
		defer targetFile.Close()
		//out.Close()
		if _, err := io.Copy(targetFile, file); err != nil {
			return models.ErrInternalServerError
		}

		//w.Write([]byte("done"))
		imagePat, error := a.isUsecase.UploadFileToBlob(fileLocation, folder)
		if error != nil {
			return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
		}
		imagePath = append(imagePath, imagePat)
		targetFile.Close()

		//out.Close()
		//os.Remove(out.Name())
		errRemove := os.Remove(fileLocation)
		if errRemove != nil {
			return models.ErrInternalServerError
		}

	}
	return c.JSON(http.StatusOK, imagePath)
}
func (a *experienceHandler) UpdateStatusExpORTrans(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}
	var experienceCommand models.NewCommandChangeStatus
	err := c.Bind(&experienceCommand)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	//if ok, err := isRequestValid(&experienceCommand); !ok {
	//	return c.JSON(http.StatusBadRequest, err.Error())
	//}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	if experienceCommand.ExpId != "" {

		response, error := a.experienceUsecase.UpdateStatus(ctx, experienceCommand.Status, experienceCommand.ExpId, token)

		if error != nil {
			return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
		}
		return c.JSON(http.StatusOK, response)
	} else if experienceCommand.TransId != "" {
		response, error := a.transportationsUsecase.UpdateStatus(ctx, experienceCommand.Status, experienceCommand.TransId, token)

		if error != nil {
			return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
		}
		return c.JSON(http.StatusOK, response)
	}
	return nil
}
func (a *experienceHandler) CreateExperiences(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}
	var experienceCommand models.NewCommandExperience
	err := c.Bind(&experienceCommand)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	//if ok, err := isRequestValid(&experienceCommand); !ok {
	//	return c.JSON(http.StatusBadRequest, err.Error())
	//}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	response, error := a.experienceUsecase.PublishExperience(ctx, experienceCommand, token)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, response)
}
func (a *experienceHandler) GetExpTransactionCount(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	transactionType := c.QueryParam("type")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := a.experienceUsecase.GetExpPendingTransactionCount(ctx, token)
	if transactionType == "failed" {
		result, err = a.experienceUsecase.GetExpFailedTransactionCount(ctx, token)
	}
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (a *experienceHandler) GetPublishedExpCount(c echo.Context) error {
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

	result, err := a.experienceUsecase.GetPublishedExpCount(ctx, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (a *experienceHandler) GetSuccessBookCount(c echo.Context) error {
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

	result, err := a.experienceUsecase.GetSuccessBookCount(ctx, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

// GetByID will get article by given id
func (a *experienceHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	//if err != nil {
	//	return c.JSON(http.StatusNotFound, models.ErrNotFound.Error())
	//}
	currency := c.QueryParam("currency")
	isMerchant := c.QueryParam("is_merchant")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	art, err := a.experienceUsecase.GetByID(ctx, id,currency,isMerchant)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, art)
}

func (a *experienceHandler) SearchExp(c echo.Context) error {
	harborID := c.QueryParam("harbor_id")
	cityID := c.QueryParam("city_id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	searchResult, err := a.experienceUsecase.SearchExp(ctx, harborID, cityID)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, searchResult)
}

func (a *experienceHandler) GetByCategoryID(c echo.Context) error {
	cid := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	categoryID, _ := strconv.Atoi(cid)
	results, err := a.experienceUsecase.GetByCategoryID(ctx, categoryID)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, results)
}

func (a *experienceHandler) FilterSearchExp(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")
	isMerchant := c.QueryParam("isMerchant")

	provinceID := c.QueryParam("province_id")
	harborID := c.QueryParam("harbor_id")
	cityID := c.QueryParam("city_id")
	qtype := c.QueryParam("type")
	guest := c.QueryParam("guest")
	trip := c.QueryParam("trip")
	bottomprice := c.QueryParam("bottomprice")
	upprice := c.QueryParam("upprice")
	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")
	sortby := c.QueryParam("sortby")
	qpage := c.QueryParam("page")
	qperPage := c.QueryParam("size")
	status := c.QueryParam("status")
	search := c.QueryParam("search")
	currency := c.QueryParam("currency")
	bookingType := c.QueryParam("booking_type")
	paymentType := c.QueryParam("payment_type")

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

	needMerchantAuth := false
	if isMerchant != "" {
		needMerchantAuth = true
	}

	searchResult, err := a.experienceUsecase.FilterSearchExp(ctx, needMerchantAuth, search, token, status, cityID, harborID, qtype, startDate, endDate, guest, trip,
		bottomprice, upprice, sortby, page, limit, offset, provinceID,currency,bookingType,paymentType)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, searchResult)
}

func (a *experienceHandler) GetUserDiscoverPreference(c echo.Context) error {
	qpage := c.QueryParam("page")
	qsize := c.QueryParam("size")
	currency := c.QueryParam("currency")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	if qpage != "" && qsize != "" {
		page, _ := strconv.Atoi(qpage)
		size, _ := strconv.Atoi(qsize)
		art, err := a.experienceUsecase.GetUserDiscoverPreference(ctx, &page, &size,currency)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, art)
	} else {
		art, err := a.experienceUsecase.GetUserDiscoverPreference(ctx, nil, nil,currency)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, art)
	}
}

func (a *experienceHandler) GetExpTypes(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	expTypeResults, err := a.experienceUsecase.GetExpTypes(ctx)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, expTypeResults)
}

func (a *experienceHandler) GetAllExperience(c echo.Context) error {

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}


	searchResult, err := a.experienceUsecase.GetExperienceMasterData(ctx)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, searchResult)
}


func (a *experienceHandler) GetExpInspirations(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	expInspirationResults, err := a.experienceUsecase.GetExpInspirations(ctx)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, expInspirationResults)
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
