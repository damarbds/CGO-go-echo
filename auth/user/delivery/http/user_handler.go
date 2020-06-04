package http

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/auth/identityserver"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/auth/user"
	"github.com/models"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// userHandler  represent the http handler for user
type userHandler struct {
	userUsecase user.Usecase
	isUsecase   identityserver.Usecase
}

// NewuserHandler will initialize the users/ resources endpoint
func NewuserHandler(e *echo.Echo, us user.Usecase, is identityserver.Usecase) {
	handler := &userHandler{
		userUsecase: us,
		isUsecase:   is,
	}
	e.POST("/users", handler.CreateUser)
	e.PUT("/users/:id", handler.UpdateUser)
	e.DELETE("/users/:id", handler.Delete)
	e.GET("/users/:id/credit", handler.GetCreditByID)
	e.GET("/users/:id", handler.GetDetailID)
	e.GET("/users", handler.List)
}

func isRequestValid(m *models.NewCommandUser) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (a *userHandler) Delete(c echo.Context) error {
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

	result, err := a.userUsecase.Delete(ctx, id, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (a *userHandler) List(c echo.Context) error {
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

	result, err := a.userUsecase.List(ctx, page, limit, offset,search)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (a *userHandler) GetCreditByID(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := a.userUsecase.GetCreditByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}
func (a *userHandler) GetDetailID(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := a.userUsecase.GetUserDetailById(ctx, id,token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}
// Store will store the user by given request body
func (a *userHandler) CreateUser(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")
	isAdmin :=  c.FormValue("is_admin")
	needAdminAuth := false
	if isAdmin != "" {
		needAdminAuth = true
	}
	filupload, image, _ := c.Request().FormFile("profile_pict_url")
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
		imagePat, _ := a.isUsecase.UploadFileToBlob(fileLocation, "UserProfile")
		imagePath = imagePat
		targetFile.Close()
		errRemove := os.Remove(fileLocation)
		if errRemove != nil {
			return models.ErrInternalServerError
		}
	}

	//phoneNumber, _ := strconv.Atoi(c.FormValue("phone_number"))
	verificationCode, _ := strconv.Atoi(c.FormValue("verification_code"))
	gender, _ := strconv.Atoi(c.FormValue("gender"))
	idType, _ := strconv.Atoi(c.FormValue("id_type"))
	referralCode, _ := strconv.Atoi(c.FormValue("referral_code"))
	points, _ := strconv.Atoi(c.FormValue("points"))
	userCommand := models.NewCommandUser{
		Id:                   c.FormValue("id"),
		UserEmail:            c.FormValue("user_email"),
		Password:             c.FormValue("password"),
		FullName:             c.FormValue("full_name"),
		PhoneNumber:          c.FormValue("phone_number"),
		VerificationSendDate: c.FormValue("verification_send_date"),
		VerificationCode:     verificationCode,
		ProfilePictUrl:       imagePath,
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
	user,error := a.userUsecase.Create(ctx, &userCommand, needAdminAuth,token)

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (a *userHandler) UpdateUser(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	filupload, image, _ := c.Request().FormFile("profile_pict_url")
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
		imagePath, _ = a.isUsecase.UploadFileToBlob(fileLocation, "UserProfile")
		targetFile.Close()
		errRemove := os.Remove(fileLocation)
		if errRemove != nil {
			return models.ErrInternalServerError
		}
	}
	//phoneNumber, _ := strconv.Atoi(c.FormValue("phone_number"))
	verificationCode, _ := strconv.Atoi(c.FormValue("verification_code"))
	gender, _ := strconv.Atoi(c.FormValue("gender"))
	idType, _ := strconv.Atoi(c.FormValue("id_type"))
	referralCode, _ := strconv.Atoi(c.FormValue("referral_code"))
	points, _ := strconv.Atoi(c.FormValue("points"))
	userCommand := models.NewCommandUser{
		Id:                   c.FormValue("id"),
		UserEmail:            c.FormValue("user_email"),
		Password:             c.FormValue("password"),
		FullName:             c.FormValue("full_name"),
		PhoneNumber:          c.FormValue("phone_number"),
		VerificationSendDate: c.FormValue("verification_send_date"),
		VerificationCode:     verificationCode,
		ProfilePictUrl:       imagePath,
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
	var isAdmin bool
	currentUser := c.FormValue("isAdmin")
	if currentUser != ""{
		isAdmin = true
	}else {
		isAdmin = false
	}
	error := a.userUsecase.Update(ctx, &userCommand, isAdmin,token)

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
