package repository_test

import (
	"context"
	"github.com/models"
	ExperienceExcludeRepo "github.com/service/exp_exclude/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	mockExperienceExclude = []models.ExperienceExclude{
		models.ExperienceExclude{
			Id:        1,
			ExpId:     "asdasdadqwdqw",
			ExcludeId: 1,
		},
		models.ExperienceExclude{
			Id:        2,
			ExpId:     "asdasdadqwdqwqrqereq",
			ExcludeId: 1,
		},
	}
	mockExperienceExcludeJoin = []models.ExperienceExcludeJoin{
		models.ExperienceExcludeJoin{
			Id:          1,
			ExpId:       "asdafafasdsa",
			ExcludeId:   1,
			ExcludeName: "qeqeqe",
			ExcludeIcon: "dadqweqwe",
		},
		models.ExperienceExcludeJoin{
			Id:          2,
			ExpId:       "aadsadsdafafasdsa",
			ExcludeId:   2,
			ExcludeName: "qeqeqe",
			ExcludeIcon: "dadqweqwe",
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
rows := sqlmock.NewRows([]string{"id", "exp_id", "exclude_id", "exclude_name", "exclude_icon"}).
		AddRow(mockExperienceExcludeJoin[0].Id, mockExperienceExcludeJoin[0].ExpId, mockExperienceExcludeJoin[0].ExcludeId, mockExperienceExcludeJoin[0].ExcludeName,
		mockExperienceExcludeJoin[0].ExcludeIcon).
		AddRow(mockExperienceExcludeJoin[1].Id, mockExperienceExcludeJoin[1].ExpId, mockExperienceExcludeJoin[1].ExcludeId, mockExperienceExcludeJoin[1].ExcludeName,
		mockExperienceExcludeJoin[1].ExcludeIcon)

	query := `SELECT ee.*,e.exclude_name, e.exclude_icon
				FROM experience_excludes ee 
				JOIN excludes e ON ee.exclude_id = e.id
				WHERE ee.exp_id = \\? `

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExperienceExcludeRepo.NewExpExcludeRepository(db)

	expId := mockExperienceExcludeJoin[0].ExpId
	anArticle, err := a.GetByExpIdJoin(context.TODO(),expId)
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

rows := sqlmock.NewRows([]string{"id", "exp_id", "exclude_id", "exclude_name", "exclude_icon"}).
			AddRow(mockExperienceExcludeJoin[0].Id, mockExperienceExcludeJoin[0].ExpId, mockExperienceExcludeJoin[0].ExcludeId, mockExperienceExcludeJoin[0].ExcludeName,
		mockExperienceExcludeJoin[0].Id).
		AddRow(mockExperienceExcludeJoin[1].Id, mockExperienceExcludeJoin[1].ExpId, mockExperienceExcludeJoin[1].ExcludeId, mockExperienceExcludeJoin[1].ExcludeName,
		mockExperienceExcludeJoin[1].Id)

	query := `SELECT ee.*,e.exclude_name, e.exclude_icon
				FROM experience_excludes ee 
				JOIN excludes e ON ee.exclude_id = e.id
				WHERE ee.exp_id = \\? `

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExperienceExcludeRepo.NewExpExcludeRepository(db)

	expId := mockExperienceExcludeJoin[0].ExpId
	_, err = a.GetByExpIdJoin(context.TODO(),expId)
	//assert.NotEmpty(t, nextCursor)
	//assert.Error(t, err)
	//assert.Nil(t, anArticle)
	//assert.Len(t, anArticle, 2)
}
func TestInsert(t *testing.T) {
	a := mockExperienceExclude[0]

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "INSERT experience_excludes SET exp_id=\\?,exclude_id=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.ExpId,a.ExcludeId).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ExperienceExcludeRepo.NewExpExcludeRepository(db)

	err = i.Insert(context.TODO(), &a)
	assert.NoError(t, err)
	//assert.Equal(t, *id, a.Id)
}
func TestInsertErrorExec(t *testing.T) {

	a := mockExperienceExclude[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()


	query := "INSERT experience_excludes SET exp_id=\\?,exclude_id=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.ExpId,a.ExcludeId,a.Id).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ExperienceExcludeRepo.NewExpExcludeRepository(db)

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


	query := "DELETE FROM experience_excludes WHERE exp_id = \\?"

	num := mockExperienceExclude[0].ExpId

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(num).WillReturnResult(sqlmock.NewResult(12, 1))

	a := ExperienceExcludeRepo.NewExpExcludeRepository(db)

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
	num := mockExperienceExclude[0].ExpId

	query := "DELETE FROM experience_excludes WHERE exp_id = \\? asdsad"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(num).WillReturnResult(sqlmock.NewResult(12, 1))

	a := ExperienceExcludeRepo.NewExpExcludeRepository(db)


	err = a.Delete(context.TODO(), num)
	assert.Error(t, err)
}
