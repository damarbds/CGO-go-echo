package usecase

import (
	"github.com/misc/faq"
	"github.com/models"
	"golang.org/x/net/context"
	"time"
)

type faqUsecase struct {
	faqUsecase     faq.Repository
	contextTimeout time.Duration
}

// NewharborsUsecase will create new an harborsUsecase object representation of harbors.Usecase interface
func NewfaqUsecase(a faq.Repository, timeout time.Duration) faq.Usecase {
	return &faqUsecase{
		faqUsecase:     a,
		contextTimeout: timeout,
	}
}

func (f faqUsecase) GetByType(c context.Context, types int) ([]*models.FAQDto, error) {
	ctx, cancel := context.WithTimeout(c, f.contextTimeout)
	defer cancel()

	res, err := f.faqUsecase.GetByType(ctx, types)
	if err != nil {
		return nil, err
	}
	var faqs []*models.FAQDto
	for _, element := range res {
		faq := models.FAQDto{
			Id:    element.Id,
			Type:  element.Type,
			Title: element.Title,
			Desc:  element.Desc,
		}
		faqs = append(faqs, &faq)
	}

	return faqs, nil
}
