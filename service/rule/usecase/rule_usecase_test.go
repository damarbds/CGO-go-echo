package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	_adminUsecaseMock "github.com/auth/admin/mocks"
	"github.com/models"
	"github.com/service/rule/mocks"
	"github.com/service/rule/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	timeoutContext = time.Second * 30
)
var (
	imagePath        = "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Rule/8941695193938718058.jpg"
	mockRuleRepo     = new(mocks.Repository)
	mockAdminUsecase = new(_adminUsecaseMock.Usecase)
	mockRule         = models.Rule{
		Id:           1,
		CreatedBy:    "test",
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		RuleName:     "Test Rule 2",
		RuleIcon:     imagePath,
	}
)

func TestGetAll(t *testing.T) {

	var mockListRule []*models.Rule
	mockListRule = append(mockListRule, &mockRule)
	page := 1
	limit := 1
	offset := 0
	count := 1
	t.Run("success", func(t *testing.T) {
		mockRuleRepo.On("Fetch", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(mockListRule, nil).Once()
		mockRuleRepo.On("GetCount", mock.Anything).Return(count, nil).Once()

		u := usecase.NewRuleUsecase(nil, mockRuleRepo, timeoutContext)

		_, err := u.GetAll(context.TODO(), page, limit, offset)

		assert.NoError(t, err)
		//assert.Len(t, list, len(mockListRule))

		mockRuleRepo.AssertExpectations(t)
		mockRuleRepo.AssertExpectations(t)
	})

}

func TestList(t *testing.T) {

	var mockListRule []*models.Rule
	mockListRule = append(mockListRule, &mockRule)

	t.Run("success", func(t *testing.T) {
		mockRuleRepo.On("List", mock.Anything).Return(mockListRule, nil).Once()

		u := usecase.NewRuleUsecase(nil, mockRuleRepo, timeoutContext)

		list, err := u.List(context.TODO())

		assert.NoError(t, err)
		assert.Len(t, list, len(mockListRule))

		mockRuleRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockRuleRepo.On("List", mock.Anything).Return(nil, errors.New("Unexpexted Error")).Once()

		u := usecase.NewRuleUsecase(nil, mockRuleRepo, timeoutContext)

		list, err := u.List(context.TODO())

		assert.Error(t, err)
		assert.Len(t, list, 0)
		mockRuleRepo.AssertExpectations(t)
	})

}

func TestGetById(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		mockRuleRepo.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(&mockRule, nil).Once()
		u := usecase.NewRuleUsecase(nil, mockRuleRepo, timeoutContext)

		a, err := u.GetById(context.TODO(), mockRule.Id)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockRuleRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockRuleRepo.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(nil, errors.New("Unexpected")).Once()

		u := usecase.NewRuleUsecase(nil, mockRuleRepo, timeoutContext)

		a, err := u.GetById(context.TODO(), mockRule.Id)

		assert.Error(t, err)
		assert.Nil(t, a)

		mockRuleRepo.AssertExpectations(t)
	})

}

func TestCreate(t *testing.T) {

	mockAdmin := &models.AdminDto{
		Id:    "972fe09f-93e9-4798-b642-14e0ac77c6be",
		Name:  "adminCGO",
		Email: "admin1234@gmail.com",
	}
	t.Run("success", func(t *testing.T) {
		tempMockRule := models.NewCommandRule{
			Id:          1,
			RuleName:    mockRule.RuleName,
			RuleIcon:    mockRule.RuleIcon,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockRule.Id = 0
		mockRuleRepo.On("Insert", mock.Anything, mock.AnythingOfType("*models.Rule")).Return(&mockRule.Id, nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin, nil).Once()

		u := usecase.NewRuleUsecase(mockAdminUsecase, mockRuleRepo, timeoutContext)

		_, err := u.Create(context.TODO(), &tempMockRule, token)

		assert.NoError(t, err)
		assert.Equal(t, mockRule.RuleName, tempMockRule.RuleName)
		mockRuleRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {
		tempMockRule := models.NewCommandRule{
			Id:          1,
			RuleName:    mockRule.RuleName,
			RuleIcon:    mockRule.RuleIcon,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockRule.Id = 0
		mockRuleRepo.On("Insert", mock.Anything, mock.AnythingOfType("*models.Rule")).Return(&mockRule.Id, nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("UnAuthorize")).Once()

		u := usecase.NewRuleUsecase(mockAdminUsecase, mockRuleRepo, timeoutContext)

		_, err := u.Create(context.TODO(), &tempMockRule, token)

		assert.Error(t, err)
		//assert.Equal(t, mockRule.RuleName, tempMockRule.RuleName)
		//mockRuleRepo.AssertExpectations(t)
	})

}

func TestUpdate(t *testing.T) {
	now := time.Now()

	mockAdmin := &models.AdminDto{
		Id:    "972fe09f-93e9-4798-b642-14e0ac77c6be",
		Name:  "adminCGO",
		Email: "admin1234@gmail.com",
	}
	mockRule.ModifiedBy = &mockAdmin.Name
	mockRule.ModifiedDate = &now
	t.Run("success", func(t *testing.T) {
		tempMockRule := models.NewCommandRule{
			Id:          1,
			RuleName:    mockRule.RuleName,
			RuleIcon:    mockRule.RuleIcon,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockRule.Id = 0
		mockRuleRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Rule")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin, nil).Once()

		u := usecase.NewRuleUsecase(mockAdminUsecase, mockRuleRepo, timeoutContext)

		_, err := u.Update(context.TODO(), &tempMockRule, token)

		assert.NoError(t, err)
		assert.Equal(t, mockRule.RuleName, tempMockRule.RuleName)
		//mockRuleRepo.AssertExpectations(t)
	})
	t.Run("success-without-image", func(t *testing.T) {
		tempMockRule := models.NewCommandRule{
			Id:          1,
			RuleName:    mockRule.RuleName,
			RuleIcon:    mockRule.RuleIcon,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockRule.Id = 0
		tempMockRule.RuleIcon = ""
		mockRuleRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Rule")).Return(nil).Once()
		mockRuleRepo.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(&mockRule, nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin, nil).Once()

		u := usecase.NewRuleUsecase(mockAdminUsecase, mockRuleRepo, timeoutContext)

		_, err := u.Update(context.TODO(), &tempMockRule, token)

		assert.NoError(t, err)
		assert.Equal(t, mockRule.RuleName, tempMockRule.RuleName)
		//mockRuleRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {
		tempMockRule := models.NewCommandRule{
			Id:          1,
			RuleName:    mockRule.RuleName,
			RuleIcon:    mockRule.RuleIcon,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockRule.Id = 0
		mockRuleRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Rule")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("unAuthorize")).Once()

		u := usecase.NewRuleUsecase(mockAdminUsecase, mockRuleRepo, timeoutContext)

		_, err := u.Update(context.TODO(), &tempMockRule, token)

		assert.Error(t, err)
		//assert.Equal(t, mockRule.RuleName, tempMockRule.RuleName)
		//mockRuleRepo.AssertExpectations(t)
	})
}
func TestDelete(t *testing.T) {

	now := time.Now()

	mockAdmin := &models.AdminDto{
		Id:    "972fe09f-93e9-4798-b642-14e0ac77c6be",
		Name:  "adminCGO",
		Email: "admin1234@gmail.com",
	}
	mockRule.ModifiedBy = &mockAdmin.Name
	mockRule.ModifiedDate = &now
	t.Run("success", func(t *testing.T) {

		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

		mockRuleRepo.On("Delete", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin, nil).Once()

		u := usecase.NewRuleUsecase(mockAdminUsecase, mockRuleRepo, timeoutContext)

		_, err := u.Delete(context.TODO(), mockRule.Id, token)

		assert.NoError(t, err)
		//assert.Equal(t, mockRule.RuleName, tempMockRule.RuleName)
		//mockRuleRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {

		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

		mockRuleRepo.On("Delete", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("unAuthorize")).Once()

		u := usecase.NewRuleUsecase(mockAdminUsecase, mockRuleRepo, timeoutContext)

		_, err := u.Delete(context.TODO(), mockRule.Id, token)

		assert.Error(t, err)
		//assert.Equal(t, mockRule.RuleName, tempMockRule.RuleName)
		//mockRuleRepo.AssertExpectations(t)
	})
}

//func TestDelete(t *testing.T) {
//	mockArticleRepo := new(mocks.Repository)
//	mockArticle := models.Article{
//		Title:   "Hello",
//		Content: "Content",
//	}
//
//	t.Run("success", func(t *testing.T) {
//		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(&mockArticle, nil).Once()
//
//		mockArticleRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()
//
//		mockAuthorrepo := new(_authorMock.Repository)
//		u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)
//
//		err := u.Delete(context.TODO(), mockArticle.ID)
//
//		assert.NoError(t, err)
//		mockArticleRepo.AssertExpectations(t)
//		mockAuthorrepo.AssertExpectations(t)
//	})
//	t.Run("article-is-not-exist", func(t *testing.T) {
//		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(nil, nil).Once()
//
//		mockAuthorrepo := new(_authorMock.Repository)
//		u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)
//
//		err := u.Delete(context.TODO(), mockArticle.ID)
//
//		assert.Error(t, err)
//		mockArticleRepo.AssertExpectations(t)
//		mockAuthorrepo.AssertExpectations(t)
//	})
//	t.Run("error-happens-in-db", func(t *testing.T) {
//		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(nil, errors.New("Unexpected Error")).Once()
//
//		mockAuthorrepo := new(_authorMock.Repository)
//		u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)
//
//		err := u.Delete(context.TODO(), mockArticle.ID)
//
//		assert.Error(t, err)
//		mockArticleRepo.AssertExpectations(t)
//		mockAuthorrepo.AssertExpectations(t)
//	})
//
//}
