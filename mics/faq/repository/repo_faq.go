package repository

import (
	"database/sql"
	"github.com/mics/faq"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type faqRepository struct {
	Conn *sql.DB
}

// NewReviewRepository will create an object that represent the exp_payment.Repository interface
func NewReviewRepository(Conn *sql.DB) faq.Repository {
	return &faqRepository{Conn}
}


func (m *faqRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.FAQ, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]*models.FAQ, 0)
	for rows.Next() {
		t := new(models.FAQ)
		err = rows.Scan(
			&t.Id,
			&t.CreatedBy,
			&t.CreatedDate,
			&t.ModifiedBy,
			&t.ModifiedDate,
			&t.DeletedBy,
			&t.DeletedDate,
			&t.IsDeleted,
			&t.IsActive,
			&t.Type,
			&t.Title ,
			&t.Desc ,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (f faqRepository) GetByType(context context.Context, types int) ([]*models.FAQ, error) {
	query := `select * from faqs
			where type = ?`
	ex,err := f.fetch(context,query,types)
	if err != nil{
		logrus.Error(err)
		return nil, err
	}
	return ex,nil
}
