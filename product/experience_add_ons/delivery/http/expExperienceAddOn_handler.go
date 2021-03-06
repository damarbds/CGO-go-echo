package http

import (
	"github.com/labstack/echo"
	"github.com/product/experience_add_ons"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"net/http"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// experience_add_onsHandler  represent the httphandler for experience_add_ons
type Experience_add_onsHandler struct {
	Experience_add_onsUsecase experience_add_ons.Usecase
}

// Newexperience_add_onsHandler will initialize the experience_add_onss/ resources endpoint
func Newexperience_add_onsHandler(e *echo.Echo, us experience_add_ons.Usecase) {
	handler := &Experience_add_onsHandler{
		Experience_add_onsUsecase: us,
	}
	//e.POST("/experience_add_onss", handler.Createexperience_add_ons)
	//e.PUT("/experience_add_onss/:id", handler.Updateexperience_add_ons)
	e.GET("product/product-add-ons", handler.GetallexperienceAddOns)
	//e.DELETE("/experience_add_onss/:id", handler.Delete)
}

//func isRequestValid(m *models.NewCommandMerchant) (bool, error) {
//	validate := validator.New()
//	err := validate.Struct(m)
//	if err != nil {
//		return false, err
//	}
//	return true, nil
//}

// GetByID will get article by given id
func (a *Experience_add_onsHandler) GetallexperienceAddOns(c echo.Context) error {
	expId := c.QueryParam("exp_id")
	currency := c.QueryParam("currency")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	art, err := a.Experience_add_onsUsecase.GetByExpId(ctx, expId,currency)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, art)
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
