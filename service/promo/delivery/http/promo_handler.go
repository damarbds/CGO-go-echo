package http

import (
	"encoding/json"
	"github.com/auth/identityserver"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo"
	"github.com/models"
	"github.com/service/promo"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	validator "gopkg.in/go-playground/validator.v9"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// promoHandler  represent the httphandler for promo
type promoHandler struct {
	isUsecase identityserver.Usecase
	promoUsecase promo.Usecase
}

// NewpromoHandler will initialize the promos/ resources endpoint
func NewpromoHandler(e *echo.Echo, us promo.Usecase,is identityserver.Usecase) {
	handler := &promoHandler{
		isUsecase:is,
		promoUsecase: us,
	}
	e.POST("admin/promo", handler.CreatePromo)
	e.PUT("admin/promo/:id", handler.UpdatePromo)
	e.GET("admin/promo", handler.List)
	e.GET("admin/promo/:id", handler.GetDetailID)
	e.DELETE("admin/promo/:id", handler.Delete)
	e.GET("service/special-promo", handler.GetAllPromo)
	e.GET("service/special-promo/:code", handler.GetPromoByCode)
}
func (a *promoHandler) Delete(c echo.Context) error {
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

	result, err := a.promoUsecase.Delete(ctx, id, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func isRequestValid(m *models.NewCommandPromo) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetByID will get article by given id
func (a *promoHandler) GetAllPromo(c echo.Context) error {
	qpage := c.QueryParam("page")
	qsize := c.QueryParam("size")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	if qpage != "" && qsize != "" {
		page, _ := strconv.Atoi(qpage)
		size, _ := strconv.Atoi(qsize)
		art, err := a.promoUsecase.Fetch(ctx, &page, &size)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, art)
	} else {
		art, err := a.promoUsecase.Fetch(ctx, nil, nil)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, art)
	}
}
// Store will store the user by given request body
func (a *promoHandler) CreatePromo(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	filupload, image, _ := c.Request().FormFile("promo_image")
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
		imagePat, _ := a.isUsecase.UploadFileToBlob(fileLocation, "Promo")
		imagePath = imagePat
		targetFile.Close()
		errRemove := os.Remove(fileLocation)
		if errRemove != nil {
			return models.ErrInternalServerError
		}
	}

	maxUsage, _ := strconv.Atoi(c.FormValue("max_usage"))
	currency, _ := strconv.Atoi(c.FormValue("currency"))
	promoValue, _ := strconv.ParseFloat(c.FormValue("promo_value"),64)
	promoType, _ := strconv.Atoi(c.FormValue("promo_type"))
	productionCapacity , _ := strconv.Atoi(c.FormValue("production_capacity"))
	promoProductType , _ := strconv.Atoi(c.FormValue("promo_product_type"))
	isAnyTripPeriod , _ := strconv.Atoi(c.FormValue("is_any_trip_period"))
	maxDiscount , _ := strconv.ParseFloat(c.FormValue("max_discount"),32)
	maxDiscount32 := float32(maxDiscount)
	merchants := c.FormValue("merchant_id")
	merchantId := make([]string,0)
	if merchants != ""{
		if errUnmarshal := json.Unmarshal([]byte(merchants), &merchantId); errUnmarshal != nil {
			return models.ErrInternalServerError
		}
	}
	promoCommand := models.NewCommandPromo{
		Id:                     "",
		PromoCode:              c.FormValue("promo_code"),
		PromoName:              c.FormValue("promo_name"),
		PromoDesc:              c.FormValue("promo_desc"),
		PromoValue:             promoValue,
		PromoType:              promoType,
		PromoImage:             imagePath,
		StartDate:              c.FormValue("start_date"),
		EndDate:                c.FormValue("end_date"),
		Currency:               currency,
		MaxUsage:               maxUsage,
		ProductionCapacity: productionCapacity,
		MerchantId:merchantId,
		PromoProductType:&promoProductType,
		StartTripPeriod: c.FormValue("start_trip_period"),
		EndTripPeriod: c.FormValue("end_trip_period"),
		IsAnyTripPeriod: isAnyTripPeriod,
		Disclaimer: c.FormValue("disclaimer"),
		MaxDiscount: maxDiscount32,
		TermCondition: c.FormValue("term_condition"),
		HowToGet: c.FormValue("how_to_get"),
		HowToUse: c.FormValue("how_to_use"),

	}
	if ok, err := isRequestValid(&promoCommand); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	promo,error := a.promoUsecase.Create(ctx, promoCommand, token)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, promo)
}

func (a *promoHandler) UpdatePromo(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	filupload, image, _ := c.Request().FormFile("promo_image")
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
		imagePath, _ = a.isUsecase.UploadFileToBlob(fileLocation, "Promo")
		targetFile.Close()
		errRemove := os.Remove(fileLocation)
		if errRemove != nil {
			return models.ErrInternalServerError
		}
	}

	maxUsage, _ := strconv.Atoi(c.FormValue("max_usage"))
	currency, _ := strconv.Atoi(c.FormValue("currency"))
	promoValue, _ := strconv.ParseFloat(c.FormValue("promo_value"),64)
	promoType, _ := strconv.Atoi(c.FormValue("promo_type"))
	productionCapacity , _ := strconv.Atoi(c.FormValue("production_capacity"))
	promoProductType , _ := strconv.Atoi(c.FormValue("promo_product_type"))
	isAnyTripPeriod , _ := strconv.Atoi(c.FormValue("is_any_trip_period"))
	maxDiscount , _ := strconv.ParseFloat(c.FormValue("max_discount"),32)
	maxDiscount32 := float32(maxDiscount)
	merchants := c.FormValue("merchant_id")
	merchantId := make([]string,0)
	if merchants != ""{
		if errUnmarshal := json.Unmarshal([]byte(merchants), &merchantId); errUnmarshal != nil {
			return models.ErrInternalServerError
		}
	}
	promoCommand := models.NewCommandPromo{
		Id:                     c.FormValue("id"),
		PromoCode:              c.FormValue("promo_code"),
		PromoName:              c.FormValue("promo_name"),
		PromoDesc:              c.FormValue("promo_desc"),
		PromoValue:             promoValue,
		PromoType:              promoType,
		PromoImage:             imagePath,
		StartDate:              c.FormValue("start_date"),
		EndDate:                c.FormValue("end_date"),
		Currency:               currency,
		MaxUsage:               maxUsage,
		ProductionCapacity: productionCapacity,
		MerchantId:merchantId,
		PromoProductType:&promoProductType,
		StartTripPeriod: c.FormValue("start_trip_period"),
		EndTripPeriod: c.FormValue("end_trip_period"),
		IsAnyTripPeriod: isAnyTripPeriod,
		Disclaimer: c.FormValue("disclaimer"),
		MaxDiscount: maxDiscount32,
		TermCondition: c.FormValue("term_condition"),
		HowToGet: c.FormValue("how_to_get"),
		HowToUse: c.FormValue("how_to_use"),
	}
	if ok, err := isRequestValid(&promoCommand); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	response,error := a.promoUsecase.Update(ctx, promoCommand, token)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, response)
}
// GetByID will get article by given id
func (a *promoHandler) List(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	qpage := c.QueryParam("page")
	qperPage := c.QueryParam("size")
	search := c.QueryParam("search")
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
		art, err := a.promoUsecase.List(ctx, page, limit,offset,search,token)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, art)

}

func (a *promoHandler) GetDetailID(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := a.promoUsecase.GetDetail(ctx,id,token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (a *promoHandler) GetPromoByCode(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	code := c.Param("code")

	promoProductType := c.QueryParam("promo_type")
	merchantId := c.QueryParam("merchant_id")

	promoType , _ := strconv.Atoi(promoProductType)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	results, err := a.promoUsecase.GetByCode(ctx, code,promoType,merchantId,token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, results)
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
