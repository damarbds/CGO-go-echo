package http

import (
	"context"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
	"net/http"

	"github.com/booking/booking_exp"
	"github.com/models"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// booking_expHandler  represent the httphandler for booking_exp
type booking_expHandler struct {
	booking_expUsecase booking_exp.Usecase
}

// Newbooking_expHandler will initialize the booking_exps/ resources endpoint
func Newbooking_expHandler(e *echo.Echo, us booking_exp.Usecase) {
	handler := &booking_expHandler{
		booking_expUsecase: us,
	}
	e.POST("booking/checkout", handler.Createbooking_exp)
	e.GET("booking/detail/:id", handler.GetDetail)
}

func isRequestValid(m *models.NewBookingExpCommand) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (a *booking_expHandler) GetDetail(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := a.booking_expUsecase.GetDetailBookingID(ctx,id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}
// Store will store the booking_exp by given request body
func (a *booking_expHandler) Createbooking_exp(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")
	//filupload, image, _ := c.Request().FormFile("ticket_qr_code")

	var bookingExpcommand models.NewBookingExpCommand
	user_id := c.FormValue("user_id")
	exp_add_ons := c.FormValue("experience_add_on_id")
	bookingExpcommand = models.NewBookingExpCommand{
		Id:                c.FormValue("id"),
		ExpId:             c.FormValue("exp_id"),
		GuestDesc:         c.FormValue("guest_desc"),
		BookedBy:          c.FormValue("booked_by"),
		BookedByEmail:     c.FormValue("booked_by_email"),
		BookingDate:       c.FormValue("booked_date"),
		UserId:            &user_id,
		Status:            c.FormValue("status"),
		TicketCode:        c.FormValue("ticket_code"),
		TicketQRCode:      "#",
		ExperienceAddOnId: &exp_add_ons,
	}

	if ok, err := isRequestValid(&bookingExpcommand); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	res, error, errorValidation := a.booking_expUsecase.Insert(ctx, &bookingExpcommand, token)
	if errorValidation != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: error.Error()})
	}
	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}

	return c.JSON(http.StatusOK, res)
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
	case models.ValidationBookedDate:
		return http.StatusBadRequest
	case models.ValidationStatus:
		return http.StatusBadRequest
	case models.ValidationBookedBy:
		return http.StatusBadRequest
	case models.ValidationExpId:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
