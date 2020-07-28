package repository_test

import (
	"context"
	"testing"
	"time"

	FAQRepo "github.com/misc/faq/repository"
	"github.com/models"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetByType(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	mockInclude := []models.FAQ{
		models.FAQ{
			Id:           1,
			CreatedBy:    "Test 1",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			Type:         1,
			Title:        "Test Faq",
			Desc:         "Faq Desc",
		},
		models.FAQ{
			Id:           2,
			CreatedBy:    "Test 2",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			Type:         1,
			Title:        "Test Faq",
			Desc:         "Faq Desc",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","type","title","desc"}).
		AddRow(mockInclude[0].Id, mockInclude[0].CreatedBy,mockInclude[0].CreatedDate,mockInclude[0].ModifiedBy,
			mockInclude[0].ModifiedDate,mockInclude[0].DeletedBy,mockInclude[0].DeletedDate,mockInclude[0].IsDeleted,
			mockInclude[0].IsActive,mockInclude[0].Type,mockInclude[0].Title,mockInclude[0].Desc).
		AddRow(mockInclude[1].Id, mockInclude[1].CreatedBy,mockInclude[1].CreatedDate,mockInclude[1].ModifiedBy,
			mockInclude[1].ModifiedDate,mockInclude[1].DeletedBy,mockInclude[1].DeletedDate,mockInclude[1].IsDeleted,
			mockInclude[1].IsActive,mockInclude[1].Type,mockInclude[1].Title,mockInclude[1].Desc)

	query := `select \*\ from faqs
			where type = \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := FAQRepo.NewReviewRepository(db)

	anArticle, err := a.GetByType(context.TODO(), 1)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestGetByTypeErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	mockInclude := []models.FAQ{
		models.FAQ{
			Id:           1,
			CreatedBy:    "Test 1",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			Type:         1,
			Title:        "Test Faq",
			Desc:         "Faq Desc",
		},
		models.FAQ{
			Id:           2,
			CreatedBy:    "Test 2",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			Type:         1,
			Title:        "Test Faq",
			Desc:         "Faq Desc",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","type","title","desc"}).
		AddRow(mockInclude[0].Id, mockInclude[0].CreatedBy,mockInclude[0].CreatedDate,mockInclude[0].ModifiedBy,
			mockInclude[0].ModifiedDate,mockInclude[0].DeletedBy,mockInclude[0].DeletedDate,mockInclude[0].IsDeleted,
			mockInclude[0].IsActive,mockInclude[0].Type,mockInclude[0].Title,mockInclude[0].Desc).
		AddRow(mockInclude[1].Id, mockInclude[1].CreatedBy,mockInclude[1].CreatedDate,mockInclude[1].ModifiedBy,
			mockInclude[1].ModifiedDate,mockInclude[1].DeletedBy,mockInclude[1].DeletedDate,mockInclude[1].IsDeleted,
			mockInclude[1].ModifiedDate,mockInclude[1].Type,mockInclude[1].Title,mockInclude[1].Desc)

	query := `select \*\ from faqs
			where type = \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := FAQRepo.NewReviewRepository(db)

	anArticle, err := a.GetByType(context.TODO(), 1)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
