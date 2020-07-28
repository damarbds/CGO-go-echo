package usecase

import (
	"context"
	"github.com/auth/admin"
	"github.com/models"
	"github.com/service/rule"
	"math"
	"strconv"
	"time"
)

type ruleUsecase struct {
	adminUsecase	admin.Usecase
	ruleRepo		rule.Repository
	contextTimeout time.Duration
}


func NewRuleUsecase(adminUsecase admin.Usecase,f rule.Repository, timeout time.Duration) rule.Usecase {
	return &ruleUsecase{
		adminUsecase:adminUsecase,
		ruleRepo:   f,
		contextTimeout: timeout,
	}
}

func (f ruleUsecase) List(ctx context.Context) ([]*models.RuleDto, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	res, err := f.ruleRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]*models.RuleDto, len(res))
	for i, n := range res {
		results[i] = &models.RuleDto{
			Id:           n.Id,
			RuleIcon: n.RuleIcon,
			RuleName: n.RuleName,
		}
	}

	return results, nil
}

func (f ruleUsecase) GetAll(ctx context.Context, page, limit, offset int) (*models.RuleDtoWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	getRules,err := f.ruleRepo.Fetch(ctx,limit,offset)
	if err != nil {
		return nil,err
	}

	rulesDtos := make([]*models.RuleDto,0)

	for _,element := range getRules{
		dto := models.RuleDto{
			Id:          element.Id,
			RuleName: element.RuleName,
			RuleIcon: element.RuleIcon,
		}
		rulesDtos = append(rulesDtos,&dto)
	}
	totalRecords, _ := f.ruleRepo.GetCount(ctx)

	totalPage := int(math.Ceil(float64(totalRecords) / float64(limit)))
	prev := page
	next := page
	if page != 1 {
		prev = page - 1
	}

	if page != totalPage {
		next = page + 1
	}

	meta := &models.MetaPagination{
		Page:          page,
		Total:         totalPage,
		TotalRecords:  totalRecords,
		Prev:          prev,
		Next:          next,
		RecordPerPage: len(rulesDtos),
	}

	response := &models.RuleDtoWithPagination{
		Data: rulesDtos,
		Meta: meta,
	}
	return response, nil
}

func (f ruleUsecase) GetById(ctx context.Context, id int) (*models.RuleDto, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	getById,err := f.ruleRepo.GetById(ctx,id)
	if err != nil {
		return nil,err
	}
	result := models.RuleDto{
		Id:          getById.Id,
		RuleName: getById.RuleName,
		RuleIcon: getById.RuleIcon,
	}
	return &result,nil

}

func (f ruleUsecase) Create(ctx context.Context, inc *models.NewCommandRule, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	currentUser, err := f.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}

	rules := models.Rule{
		Id:           0,
		CreatedBy:    currentUser.Name,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		RuleName:  inc.RuleName,
		RuleIcon:  inc.RuleIcon,
	}
	id ,err := f.ruleRepo.Insert(ctx,&rules)
	if err != nil {
		return nil,err
	}

	result := models.ResponseDelete{
		Id:      strconv.Itoa(*id),
		Message: "Success Create Rule",
	}
	return &result,nil
}

func (f ruleUsecase) Update(ctx context.Context, inc *models.NewCommandRule, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	currentUser, err := f.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	rules := models.Rule{
		Id:           inc.Id,
		CreatedBy:    currentUser.Name,
		CreatedDate:  time.Time{},
		ModifiedBy:   &currentUser.Name,
		ModifiedDate: &now,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		RuleName:  inc.RuleName,
		RuleIcon:  inc.RuleIcon,
	}
	err = f.ruleRepo.Update(ctx,&rules)
	if err != nil {
		return nil,err
	}

	result := models.ResponseDelete{
		Id:      strconv.Itoa(inc.Id),
		Message: "Success Update Rule",
	}
	return &result,nil
}

func (f ruleUsecase) Delete(ctx context.Context, id int, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	currentUser, err := f.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}

	err = f.ruleRepo.Delete(ctx,id,currentUser.Name)
	if err != nil {
		return nil,err
	}

	result := models.ResponseDelete{
		Id:      strconv.Itoa(id),
		Message: "Success Delete Rule",
	}
	return &result,nil
}