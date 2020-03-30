package http

import (
	"context"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"

	"github.com/models"
	"github.com/user"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// userHandler  represent the httphandler for user
type userHandler struct {
	userUsecase user.Usecase
}

// NewuserHandler will initialize the users/ resources endpoint
func NewuserHandler(e *echo.Echo, us user.Usecase) {
	handler := &userHandler{
		userUsecase: us,
	}
	e.POST("/users", handler.CreateUser)
	e.PUT("/users/:id", handler.UpdateUser)
}

func isRequestValid(m *models.NewCommandUser) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store will store the user by given request body
func (a *userHandler) CreateUser(c echo.Context) error {

	//name := c.FormValue("name")
	//var userCommand models.NewCommandUser
	//err := c.Bind(&userCommand)
	//if err != nil {
	//	return c.JSON(http.StatusUnprocessableEntity, err.Error())
	//}
	phoneNumber, _ := strconv.Atoi(c.FormValue("phone_number"))
	verificationCode, _ := strconv.Atoi(c.FormValue("verification_code"))
	gender, _ := strconv.Atoi(c.FormValue("gender"))
	idType, _ := strconv.Atoi(c.FormValue("id_type"))
	referralCode, _ := strconv.Atoi(c.FormValue("referral_code"))
	points, _ := strconv.Atoi(c.FormValue("points"))
	userCommand:= models.NewCommandUser{
		Id:                   c.FormValue("id"),
		UserEmail:            c.FormValue("user_email"),
		Password:             c.FormValue("password"),
		FullName:             c.FormValue("full_name"),
		PhoneNumber:          phoneNumber,
		VerificationSendDate: c.FormValue("verification_send_date"),
		VerificationCode:     verificationCode,
		ProfilePictUrl:       "#",
		Address:              c.FormValue("address"),
		Dob:                  c.FormValue("dob"),
		Gender:               gender,
		IdType:               idType,
		IdNumber:             c.FormValue("id_number"),
		ReferralCode:         referralCode,
		Points:               points,
	}
	if ok, err := isRequestValid(&userCommand); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	error := a.userUsecase.Create(ctx, &userCommand, "admin")

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, userCommand)
}

func (a *userHandler) UpdateUser(c echo.Context) error {
	//var userCommand models.NewCommandUser
	//err := c.Bind(&userCommand)
	//if err != nil {
	//	return c.JSON(http.StatusUnprocessableEntity, err.Error())
	//}
	phoneNumber, _ := strconv.Atoi(c.FormValue("phone_number"))
	verificationCode, _ := strconv.Atoi(c.FormValue("verification_code"))
	gender, _ := strconv.Atoi(c.FormValue("gender"))
	idType, _ := strconv.Atoi(c.FormValue("id_type"))
	referralCode, _ := strconv.Atoi(c.FormValue("referral_code"))
	points, _ := strconv.Atoi(c.FormValue("points"))
	userCommand:= models.NewCommandUser{
		Id:                   c.FormValue("id"),
		UserEmail:            c.FormValue("user_email"),
		Password:             c.FormValue("password"),
		FullName:             c.FormValue("full_name"),
		PhoneNumber:          phoneNumber,
		VerificationSendDate: c.FormValue("verification_send_date"),
		VerificationCode:     verificationCode,
		ProfilePictUrl:       "#",
		Address:              c.FormValue("address"),
		Dob:                  c.FormValue("dob"),
		Gender:               gender,
		IdType:               idType,
		IdNumber:             c.FormValue("id_number"),
		ReferralCode:         referralCode,
		Points:               points,
	}
	if ok, err := isRequestValid(&userCommand); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	error := a.userUsecase.Update(ctx, &userCommand, "admin")

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, userCommand)
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
