package http

import (
	"golang.org/x/net/context"
	"net/http"
	"strconv"

	//"strconv"

	"github.com/labstack/echo"
	"github.com/misc/faq"
	"github.com/models"
	"github.com/sirupsen/logrus"
	//"golang.org/x/net/context"
	validator "gopkg.in/go-playground/validator.v9"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// faqHandler  represent the httphandler for faq
type faqHandler struct {
	faqUsecase faq.Usecase
}

// NewfaqHandler will initialize the faqs/ resources endpoint
func NewfaqHandler(e *echo.Echo, us faq.Usecase) {
	handler := &faqHandler{
		faqUsecase: us,
	}
	//e.POST("/faqs", handler.Createfaq)
	//e.PUT("/faqs/:id", handler.Updatefaq)
	e.GET("misc/faq", handler.GetByType)
	//e.DELETE("/faqs/:id", handler.Delete)
}

func isRequestValid(m *models.NewCommandMerchant) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
//
//// GetByID will get article by given id
func (a *faqHandler) GetByType(c echo.Context) error {
	qtypes := c.QueryParam("type")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
		types , _:= strconv.Atoi(qtypes)
		art, err := a.faqUsecase.GetByType(ctx,types)
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
