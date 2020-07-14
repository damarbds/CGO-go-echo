package repository_test

import (
	"context"
	"testing"

	"github.com/models"
	ExpInspirationRepo "github.com/service/exp_inspiration/repository"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	mockExpInspiration = []models.ExpInspirationObject{
		models.ExpInspirationObject{
			ExpInspirationID: "qeqwsadasdasd",
			ExpId:            "sadsadasdasdasdqwewq",
			ExpTitle:         "Test 1",
			ExpDesc:          "asdadasdqqw",
			ExpCoverPhoto:    "adasdqqweqw",
			ExpType:          "asdadas",
			Rating:           1,
		},
		models.ExpInspirationObject{
			ExpInspirationID: "qeqwsadasdasdasdasdas",
			ExpId:            "sadsadaqeqweqwasdasdasdqwewq",
			ExpTitle:         "Test 2",
			ExpDesc:          "asdadasdqqwqesad",
			ExpCoverPhoto:    "adasdqqweqw",
			ExpType:          "asdadas",
			Rating:           1,
		},
	}
)

func TestGetExpInspirations(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"exp_inspiration_id", "exp_id", "exp_title", "exp_desc", "exp_cover_photo",
		"exp_type", "rating"}).
		AddRow(mockExpInspiration[0].ExpInspirationID, mockExpInspiration[0].ExpId, mockExpInspiration[0].ExpTitle,
		mockExpInspiration[0].ExpDesc, mockExpInspiration[0].ExpCoverPhoto,mockExpInspiration[0].ExpType,
		mockExpInspiration[0].Rating).
		AddRow(mockExpInspiration[1].ExpInspirationID, mockExpInspiration[1].ExpId, mockExpInspiration[1].ExpTitle,
		mockExpInspiration[1].ExpDesc, mockExpInspiration[1].ExpCoverPhoto,mockExpInspiration[1].ExpType,
		mockExpInspiration[1].Rating)

	query := `SELECT
		ei.id as exp_inspiration_id,
		ei.exp_id,
		ei.exp_title,
		ei.exp_desc,
		ei.exp_cover_photo,
		e.exp_type,
		e.rating
	FROM
		exp_inspirations ei
	JOIN 
		experiences e on e.id = ei.exp_id
	WHERE
		ei.is_deleted = 0
		AND ei.is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExpInspirationRepo.NewExpInspirationRepository(db)
	//
	//expId := mockExperienceIncludeJoin[0].ExpId
	anArticle, err := a.GetExpInspirations(context.TODO())
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestGetExpInspirationsErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"exp_inspiration_id", "exp_id", "exp_title", "exp_desc", "exp_cover_photo",
		"exp_type", "rating"}).
		AddRow(mockExpInspiration[0].ExpInspirationID, mockExpInspiration[0].ExpId, mockExpInspiration[0].ExpTitle,
			mockExpInspiration[0].ExpDesc, mockExpInspiration[0].ExpCoverPhoto,mockExpInspiration[0].ExpType,
			mockExpInspiration[0].Rating).
		AddRow(mockExpInspiration[1].ExpInspirationID, mockExpInspiration[1].ExpId, mockExpInspiration[1].ExpTitle,
			mockExpInspiration[1].ExpDesc, mockExpInspiration[1].ExpCoverPhoto,mockExpInspiration[1].ExpType,
			mockExpInspiration[1].Rating)

	query := `SELECT
		ei.id as exp_inspiration_id,
		ei.exp_id,
		ei.exp_title,
		ei.exp_desc,
		ei.exp_cover_photo,
		e.exp_type,
		e.rating
	FROM
		exp_inspirations ei
	JOIN 
		experiences e on e.id = ei.exp_id
	WHERE
		ei.is_deleted = 0
		AND ei.is_active = 1asdasds`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExpInspirationRepo.NewExpInspirationRepository(db)

	//expId := mockExperienceIncludeJoin[0].ExpId
	_, err = a.GetExpInspirations(context.TODO())
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	//assert.Nil(t, anArticle)
	//assert.Len(t, anArticle, 2)
}
