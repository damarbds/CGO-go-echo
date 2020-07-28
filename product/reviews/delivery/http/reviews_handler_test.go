package http_test

import (
	"encoding/json"
	"errors"
	"github.com/bxcodec/faker"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/product/reviews/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	ReviewHttp "github.com/product/reviews/delivery/http"
)
var(
	timeoutContext = time.Second*30
	rating = 12
	userId = "adasdasdqweqwe"
	reviewValue float64 = 12
	mockReview = []*models.Review{
		&models.Review{
			Id:                "asdasdasdsadqeq",
			CreatedBy:         "Test 1",
			CreatedDate:       time.Now(),
			ModifiedBy:        nil,
			ModifiedDate:      nil,
			DeletedBy:         nil,
			DeletedDate:       nil,
			IsDeleted:         0,
			IsActive:          1,
			Values:            12,
			Desc:               `{"name":"radyaupdate12345","userid":"radyaupdate12345","desc":"good exp"}`,
			ExpId:             "qewqweqwe",
			UserId:            &userId,
			GuideReview:       &reviewValue,
			ActivitiesReview:  &reviewValue,
			ServiceReview:     &reviewValue,
			CleanlinessReview: &reviewValue,
			ValueReview:       &reviewValue,
		},
		&models.Review{
			Id:                "jlkjlkjlkjkl",
			CreatedBy:         "Test 1",
			CreatedDate:       time.Now(),
			ModifiedBy:        nil,
			ModifiedDate:      nil,
			DeletedBy:         nil,
			DeletedDate:       nil,
			IsDeleted:         0,
			IsActive:          1,
			Values:            12,
			Desc:               `{"name":"radyaupdate12345","userid":"radyaupdate12345","desc":"good exp"}`,
			ExpId:             "qewqweqwe",
			UserId:            &userId,
			GuideReview:       &reviewValue,
			ActivitiesReview:  &reviewValue,
			ServiceReview:     &reviewValue,
			CleanlinessReview: &reviewValue,
			ValueReview:       &reviewValue,
		},
	}
	mockReviewDto = []*models.ReviewDto{
		&models.ReviewDto{
			Name:              "Test",
			Image:             "adasdasd",
			Desc:              `{"name":"radyaupdate12345","userid":"radyaupdate12345","desc":"good exp"}`,
			Values:            12121,
			Date:              time.Now(),
			UserId:            &userId,
			GuideReview:       &reviewValue,
			ActivitiesReview:  &reviewValue,
			ServiceReview:     &reviewValue,
			CleanlinessReview: &reviewValue,
			ValueReview:       &reviewValue,
		},
	}
)
func TestGetAllReview(t *testing.T) {

	mockUCase := new(mocks.Usecase)

	var mockReviewPagination = models.ReviewsWithPagination{}

	err := faker.FakeData(&mockReviewPagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/product/exp-reviews?isMerchant=true&page=" + strconv.Itoa(mockReviewPagination.Meta.Page) + "&size=" + strconv.Itoa(mockReviewPagination.Meta.RecordPerPage) , strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetReviewsByExpIdWithPagination", mock.Anything,mock.AnythingOfType("int"),
		mock.AnythingOfType("int"),mock.AnythingOfType("int"),mock.AnythingOfType("int"),
		mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(&mockReviewPagination,nil)

	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Request().Header.Add("Authorization",token)

	handler := ReviewHttp.ReviewsHandler{
		ReviewsUsecase:mockUCase,
	}
	err = handler.GetReviewsByExpId(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetAllReviewErrorInternalServer(t *testing.T) {


	mockUCase := new(mocks.Usecase)

	var mockReviewPagination = models.ReviewsWithPagination{}

	err := faker.FakeData(&mockReviewPagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/product/exp-reviews?isMerchant=true&page=" + strconv.Itoa(mockReviewPagination.Meta.Page) + "&size=" + strconv.Itoa(mockReviewPagination.Meta.RecordPerPage) , strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetReviewsByExpIdWithPagination", mock.Anything,mock.AnythingOfType("int"),
		mock.AnythingOfType("int"),mock.AnythingOfType("int"),mock.AnythingOfType("int"),
		mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(nil,errors.New("UnExpected"))

	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Request().Header.Add("Authorization",token)

	handler := ReviewHttp.ReviewsHandler{
		ReviewsUsecase:mockUCase,
	}
	err = handler.GetReviewsByExpId(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestCreateReview(t *testing.T) {
	mockReview := models.NewReviewCommand{
		Id:                mockReview[0].Id,
		ExpId:             mockReview[0].ExpId,
		Desc:              mockReview[0].Desc,
		GuideReview:       reviewValue,
		ActivitiesReview:  reviewValue,
		ServiceReview:     reviewValue,
		CleanlinessReview: reviewValue,
		ValueReview:       reviewValue,
	}
	tempMockReview := mockReview
	mockUCase := new(mocks.Usecase)
	j, err := json.Marshal(tempMockReview)
	assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("CreateReviews", mock.Anything, mock.AnythingOfType("models.NewReviewCommand"),mock.AnythingOfType("string")).Return(&mockReview,nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/product/exp-reviews", strings.NewReader(string(j)))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/product/exp-reviews")
	c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	handler := ReviewHttp.ReviewsHandler{
		ReviewsUsecase:mockUCase,
	}
	err = handler.CreateReview(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestCreateReviewWithoutToken(t *testing.T) {
	mockReview := models.NewReviewCommand{
		Id:                mockReview[0].Id,
		ExpId:             mockReview[0].ExpId,
		Desc:              mockReview[0].Desc,
		GuideReview:       reviewValue,
		ActivitiesReview:  reviewValue,
		ServiceReview:     reviewValue,
		CleanlinessReview: reviewValue,
		ValueReview:       reviewValue,
	}
	tempMockReview := mockReview
	mockUCase := new(mocks.Usecase)
	j, err := json.Marshal(tempMockReview)
	assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("CreateReviews", mock.Anything, mock.AnythingOfType("models.NewReviewCommand"),mock.AnythingOfType("string")).Return(&mockReview,nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/product/exp-reviews", strings.NewReader(string(j)))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/product/exp-reviews")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	handler := ReviewHttp.ReviewsHandler{
		ReviewsUsecase:mockUCase,
	}
	err = handler.CreateReview(c)
	//require.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestCreateReviewConflict(t *testing.T) {
	mockReview := models.NewReviewCommand{
		Id:                mockReview[0].Id,
		ExpId:             mockReview[0].ExpId,
		Desc:              mockReview[0].Desc,
		GuideReview:       reviewValue,
		ActivitiesReview:  reviewValue,
		ServiceReview:     reviewValue,
		CleanlinessReview: reviewValue,
		ValueReview:       reviewValue,
	}
	tempMockReview := mockReview
	mockUCase := new(mocks.Usecase)
	j, err := json.Marshal(tempMockReview)
	assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("CreateReviews", mock.Anything, mock.AnythingOfType("models.NewReviewCommand"),mock.AnythingOfType("string")).Return(nil,models.ErrConflict)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/product/exp-reviews", strings.NewReader(string(j)))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/product/exp-reviews")
	c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	handler := ReviewHttp.ReviewsHandler{
		ReviewsUsecase:mockUCase,
	}
	err = handler.CreateReview(c)
	//require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//mockUCase.AssertExpectations(t)
}
