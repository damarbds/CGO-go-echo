package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

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
	e.POST("booking/checkout", handler.CreateBooking)
	e.GET("booking/detail/:id", handler.GetDetail)
	e.GET("booking/my", handler.GetMyBooking)
	e.GET("booking/history-user", handler.GetHistoryBookingByUser)
	e.GET("booking/growth", handler.GetGrowth)
	e.GET("booking/count-month", handler.CountThisMonth)
	e.GET("booking/check-experience", handler.CheckBookingCountGuest)
}

func (a *booking_expHandler) CountThisMonth(c echo.Context) error {
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

	result, err := a.booking_expUsecase.CountThisMonth(ctx)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (a *booking_expHandler) GetGrowth(c echo.Context) error {
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

	res, err := a.booking_expUsecase.GetGrowthByMerchantID(ctx, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func isRequestValid(m *models.NewBookingExpCommand) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *booking_expHandler) GetMyBooking(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	status := c.QueryParam("status")
	qpage := c.QueryParam("page")
	qperPage := c.QueryParam("size")

	var limit = 20
	var page = 1
	var offset = 0

	page, _ = strconv.Atoi(qpage)
	limit, _ = strconv.Atoi(qperPage)
	offset = (page - 1) * limit

	var transactionStatus, bookingStatus int
	if status == "confirm" {
		transactionStatus = 3
		bookingStatus = 1
	} else if status == "waiting" {
		transactionStatus = 1
		bookingStatus = 1
	} else if status == "pending" {
		transactionStatus = 0
		bookingStatus = 1
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := a.booking_expUsecase.GetByUserID(ctx, transactionStatus, bookingStatus, token, page, limit, offset)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (a *booking_expHandler) CheckBookingCountGuest(c echo.Context) error {

	expId := c.QueryParam("exp_id")
	date := c.QueryParam("date")
	qguest := c.QueryParam("guest")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	guest, _ := strconv.Atoi(qguest)
	result, err := a.booking_expUsecase.GetByGuestCount(ctx, expId, date, guest)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}
func (a *booking_expHandler) GetHistoryBookingByUser(c echo.Context) error {

	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}
	monthType := c.QueryParam("month_type")
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

	result, err := a.booking_expUsecase.GetHistoryBookingByUserId(ctx, token, monthType, page, limit, offset)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (a *booking_expHandler) GetDetail(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := a.booking_expUsecase.GetDetailBookingID(ctx, id, id)
	if err != nil {
		if err == models.ErrNotFound {
			result, err = a.booking_expUsecase.GetDetailTransportBookingID(ctx, id, id)
			if err != nil {
				return c.JSON(getStatusCode(err), ResponseError{Message: "Get Booking Trans Detail Failed"})
			}
			return c.JSON(http.StatusOK, result)
		}
		return c.JSON(getStatusCode(err), ResponseError{Message: "Get Booking Exp Detail Failed"})
	} else {
		return c.JSON(http.StatusOK, result)
	}
}

// Store will store the booking_exp by given request body
func (a *booking_expHandler) CreateBooking(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")
	//filupload, image, _ := c.Request().FormFile("ticket_qr_code")

	var bookingExpcommand models.NewBookingExpCommand
	user_id := c.FormValue("user_id")
	exp_add_ons := c.FormValue("experience_add_on_id")
	transId := c.FormValue("trans_id")
	scheduleId := c.FormValue("schedule_id")
	transReturnId := c.FormValue("trans_return_id")
	scheduleReturnId := c.FormValue("schedule_return_id")
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
		TransId:           &transId,
		ScheduleId:        &scheduleId,
	}

	if ok, err := isRequestValid(&bookingExpcommand); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	res, err, errorValidation := a.booking_expUsecase.Insert(ctx, &bookingExpcommand, transReturnId, scheduleReturnId, token)
	if errorValidation != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: errorValidation.Error()})
	}
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
