package repository

import (
	"database/sql"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"github.com/transactions/experience_payment_type"
	"golang.org/x/net/context"
)

type experiencePaymentTypeRepository struct {
	Conn *sql.DB
}


// NewPaymentRepository will create an object that represent the article.Repository interface
func NewExperiencePaymentTypeRepository(Conn *sql.DB) experience_payment_type.Repostiory {
	return &experiencePaymentTypeRepository{Conn}
}

func (m *experiencePaymentTypeRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.ExperiencePaymentType, error) {
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

	result := make([]*models.ExperiencePaymentType, 0)
	for rows.Next() {
		t := new(models.ExperiencePaymentType)
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
			&t.ExpPaymentTypeName ,
			&t.ExpPaymentTypeDesc,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (e experiencePaymentTypeRepository) GetAll(ctx context.Context, page *int, size *int) ([]*models.ExperiencePaymentType, error) {
		//var res []*models.ExperiencePaymentType
		if page != nil && size != nil{

			query := `SELECT * from experience_payment_types 
				WHERE is_deleted = 0  AND is_active = 1 
 				ORDER BY created_date desc LIMIT ? OFFSET ? `

			res, err := e.fetch(ctx, query, size, page)
			if err != nil {
				return nil, err
			}

			return res, err
		}else {
			query := `SELECT * from experience_payment_types 
				WHERE is_deleted = 0  AND is_active = 1 
 				ORDER BY created_date desc LIMIT ? OFFSET ? `

			res, err := e.fetch(ctx, query, size, page)
			if err != nil {
				return nil, err
			}
			return res, err
		}

	return nil, nil
}
