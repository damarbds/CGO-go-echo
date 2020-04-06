package usecase

import (
	"github.com/models"
	"github.com/transaction/experience_payment_type"
	"golang.org/x/net/context"
	"time"
)

type experiencePaymentTypeUsecase struct {
	experiencePaymentTypeRepo experience_payment_type.Repostiory
	contextTimeout time.Duration
}


// NewPaymentUsecase will create new an paymentUsecase object representation of payment.Usecase interface
func NewexperiencePaymentTypeUsecase(p experience_payment_type.Repostiory, timeout time.Duration) experience_payment_type.Usecase {
	return &experiencePaymentTypeUsecase{
		experiencePaymentTypeRepo:p,
		contextTimeout: timeout,
	}
}


func (e experiencePaymentTypeUsecase) GetAll(c context.Context, page *int, size *int) ([]*models.ExperiencePaymentTypeDto, error) {
	ctx, cancel := context.WithTimeout(c, e.contextTimeout)
	defer cancel()

	var result []*models.ExperiencePaymentTypeDto
	if page != nil && size != nil {
		query ,error := e.experiencePaymentTypeRepo.GetAll(ctx,page,size)
		if error != nil {
			return nil,error
		}
		for _,element := range query{
			dto := models.ExperiencePaymentTypeDto{
				Id:   element.Id,
				Name: element.ExpPaymentTypeName,
				Desc: element.ExpPaymentTypeDesc,
			}
			result = append(result,&dto)
		}
	}else {
		query ,error := e.experiencePaymentTypeRepo.GetAll(ctx,nil,nil)
		if error != nil {
			return nil,error
		}
		for _,element := range query{
			dto := models.ExperiencePaymentTypeDto{
				Id:   element.Id,
				Name: element.ExpPaymentTypeName,
				Desc: element.ExpPaymentTypeDesc,
			}
			result = append(result,&dto)
		}
	}
	return result,nil

}