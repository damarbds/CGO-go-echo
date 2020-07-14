package repository_test

import (
	"context"
	"testing"

	"github.com/models"
	ExperienceIncludeRepo "github.com/service/exp_Include/repository"

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
		mockExpInspiration[1].ExpDesc, mockExpInspiration[0].ExpCoverPhoto,mockExpInspiration[0].ExpType,
		mockExpInspiration[0].Rating)

	query := `SELECT ei.*,i.include_name,i.include_icon
				FROM experience_includes ei 
				JOIN includes i ON ei.include_id = i.id
				WHERE ei.exp_id = \\? `

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExperienceIncludeRepo.NewExpIncludeRepository(db)

	expId := mockExperienceIncludeJoin[0].ExpId
	anArticle, err := a.GetByExpIdJoin(context.TODO(), expId)
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

	rows := sqlmock.NewRows([]string{"id", "exp_id", "include_id", "include_name", "include_icon"}).
		AddRow(mockExperienceIncludeJoin[0].Id, mockExperienceIncludeJoin[0].ExpId, mockExperienceIncludeJoin[0].IncludeId, mockExperienceIncludeJoin[0].IncludeName,
			mockExperienceIncludeJoin[0].Id).
		AddRow(mockExperienceIncludeJoin[1].Id, mockExperienceIncludeJoin[1].ExpId, mockExperienceIncludeJoin[1].IncludeId, mockExperienceIncludeJoin[1].IncludeName,
			mockExperienceIncludeJoin[1].Id)

	query := `SELECT ei.*,i.include_name,i.include_icon
				FROM experience_includes ei 
				JOIN includes i ON ei.include_id = i.id
				WHERE ei.exp_id = \\? asdsadasd`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExperienceIncludeRepo.NewExpIncludeRepository(db)

	expId := mockExperienceIncludeJoin[0].ExpId
	_, err = a.GetByExpIdJoin(context.TODO(), expId)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	//assert.Nil(t, anArticle)
	//assert.Len(t, anArticle, 2)
}
