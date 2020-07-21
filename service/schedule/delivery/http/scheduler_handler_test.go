package http_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/models"
	"github.com/service/schedule/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ScheduleHttp "github.com/service/schedule/delivery/http"
)

func TestGetSchedule(t *testing.T) {
	merchantId := "sakldjlkasdjsalkdj"
	date := "2020-09"
	mockSchedule := models.ScheduleDtoObj{
		MerchantId:   merchantId,
		ScheduleDate: []models.ScheduleObjDate{
			models.ScheduleObjDate{
				Date:                "2020-09-01",
				TransportationCount: 2,
				ExperienceCount:     2,
			},
		},
	}
	//err := faker.FakeData(&mockSchedule)
	//assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)


	mockUCase.On("GetScheduleByMerchantId", mock.Anything, merchantId,date).Return(&mockSchedule, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/service/schedule?merchant_id="+merchantId+"&date="+date, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := ScheduleHttp.ScheduleHandler{
		ScheduleUsecase: mockUCase,
	}
	err = handler.GetSchedule(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetScheduleErrorNotFound(t *testing.T) {
	merchantId := "sakldjlkasdjsalkdj"
	date := "2020-09"

	//err := faker.FakeData(&mockSchedule)
	//assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)


	mockUCase.On("GetScheduleByMerchantId", mock.Anything, merchantId,date).Return(nil, models.ErrNotFound)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/service/schedule?merchant_id="+merchantId+"&date="+date, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := ScheduleHttp.ScheduleHandler{
		ScheduleUsecase: mockUCase,
	}
	err = handler.GetSchedule(c)
	//assert.Error(t, err)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestCreateSchedule(t *testing.T) {
	mockSchedule := models.NewCommandSchedule{
		TimeObj:   []models.TimeObj{
			models.TimeObj{
				DepartureTime: "12:00:00",
				ArrivalTime:   "13:00:00",
			},
		},
		CreatedBy: "Admin",
		TransId:   "qweqweqweqwe",
		Schedule:  []models.YearObj{
			models.YearObj{
				Year:  2020,
				Month: []models.MonthObj{
					models.MonthObj{
						Month:    "June",
						DayPrice: []models.DayPriceObj{
							models.DayPriceObj{
								DepartureDate: "2021-04-18",
								Day:           "Sunday",
								AdultPrice:    10,
								ChildrenPrice: 10,
								Currency:      "IDR",
							},
						},
					},
				},
			},
		},
	}

	//tempMockSchedule := mockSchedule
	mockUCase := new(mocks.Usecase)
	//j, err := json.Marshal(tempMockSchedule)
	//assert.NoError(t, err)
	j, err := json.Marshal(mockSchedule)
	assert.NoError(t, err)
	mockUCase.On("InsertSchedule", mock.Anything, mock.AnythingOfType("*models.NewCommandSchedule")).Return(&mockSchedule, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/service/schedule", strings.NewReader(string(j)))
	//req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/service/schedule")
	//c.Request().Header.Add("Authorization", token)
	//c.Request().ParseForm()
	handler := ScheduleHttp.ScheduleHandler{
		ScheduleUsecase: mockUCase,
		//IsUsecase:       mockIsUsecase,
	}
	err = handler.CreateSchedule(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestCreateScheduleInternalServerError(t *testing.T) {
	mockSchedule := models.NewCommandSchedule{
		TimeObj:   []models.TimeObj{
			models.TimeObj{
				DepartureTime: "12:00:00",
				ArrivalTime:   "13:00:00",
			},
		},
		CreatedBy: "Admin",
		TransId:   "qweqweqweqwe",
		Schedule:  []models.YearObj{
			models.YearObj{
				Year:  2020,
				Month: []models.MonthObj{
					models.MonthObj{
						Month:    "June",
						DayPrice: []models.DayPriceObj{
							models.DayPriceObj{
								DepartureDate: "2021-04-18",
								Day:           "Sunday",
								AdultPrice:    10,
								ChildrenPrice: 10,
								Currency:      "IDR",
							},
						},
					},
				},
			},
		},
	}

	//tempMockSchedule := mockSchedule
	mockUCase := new(mocks.Usecase)
	//j, err := json.Marshal(tempMockSchedule)
	//assert.NoError(t, err)
	j, err := json.Marshal(mockSchedule)
	assert.NoError(t, err)
	mockUCase.On("InsertSchedule", mock.Anything, mock.AnythingOfType("*models.NewCommandSchedule")).Return(nil, errors.New("UnExpected"))

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/service/schedule", strings.NewReader(string(j)))
	//req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/service/schedule")
	//c.Request().Header.Add("Authorization", token)
	//c.Request().ParseForm()
	handler := ScheduleHttp.ScheduleHandler{
		ScheduleUsecase: mockUCase,
		//IsUsecase:       mockIsUsecase,
	}
	err = handler.CreateSchedule(c)
	//require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}
