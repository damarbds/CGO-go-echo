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
	mockExperienceInclude = []models.ExperienceInclude{
		models.ExperienceInclude{
			Id:        1,
			ExpId:     "asdasdadqwdqw",
			IncludeId: 1,
		},
		models.ExperienceInclude{
			Id:        2,
			ExpId:     "asdasdadqwdqwqrqereq",
			IncludeId: 1,
		},
	}
	mockExperienceIncludeJoin = []models.ExperienceIncludeJoin{
		models.ExperienceIncludeJoin{
			Id:          1,
			ExpId:       "asdafafasdsa",
			IncludeId:   1,
			IncludeName: "qeqeqe",
			IncludeIcon: "dadqweqwe",
		},
		models.ExperienceIncludeJoin{
			Id:          2,
			ExpId:       "aadsadsdafafasdsa",
			IncludeId:   2,
			IncludeName: "qeqeqe",
			IncludeIcon: "dadqweqwe",
		},
	}
)

func TestGetByExpIdJoin(t *testing.T) {
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
			mockExperienceIncludeJoin[0].IncludeIcon).
		AddRow(mockExperienceIncludeJoin[1].Id, mockExperienceIncludeJoin[1].ExpId, mockExperienceIncludeJoin[1].IncludeId, mockExperienceIncludeJoin[1].IncludeName,
			mockExperienceIncludeJoin[1].IncludeIcon)

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
func TestGetByExpIdJoinErrorFetch(t *testing.T) {
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
func TestInsert(t *testing.T) {
	a := mockExperienceInclude[0]

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "INSERT experience_includes SET exp_id=\\?,include_id=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.ExpId, a.IncludeId).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ExperienceIncludeRepo.NewExpIncludeRepository(db)

	err = i.Insert(context.TODO(), &a)
	assert.NoError(t, err)
	//assert.Equal(t, *id, a.Id)
}
func TestInsertErrorExec(t *testing.T) {

	a := mockExperienceInclude[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "INSERT experience_Includes SET exp_id=\\?,Include_id=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.ExpId, a.IncludeId, a.Id).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ExperienceIncludeRepo.NewExpIncludeRepository(db)

	err = i.Insert(context.TODO(), &a)
	assert.Error(t, err)
}
func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "DELETE FROM experience_includes WHERE exp_id = \\?"

	num := mockExperienceInclude[0].ExpId

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(num).WillReturnResult(sqlmock.NewResult(12, 1))

	a := ExperienceIncludeRepo.NewExpIncludeRepository(db)

	err = a.Delete(context.TODO(), num)
	assert.NoError(t, err)
}
func TestDeleteErrorExecQueryString(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	num := mockExperienceInclude[0].ExpId

	query := "DELETE FROM experience_Includes WHERE exp_id = \\? asdsad"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(num).WillReturnResult(sqlmock.NewResult(12, 1))

	a := ExperienceIncludeRepo.NewExpIncludeRepository(db)

	err = a.Delete(context.TODO(), num)
	assert.Error(t, err)
}
