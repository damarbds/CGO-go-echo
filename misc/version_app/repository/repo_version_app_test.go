package repository_test

import (
	"context"
	"testing"


	VersionAppRepo "github.com/misc/version_app/repository"
	"github.com/models"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	mockVersionApp := []models.VersionApp{
		models.VersionApp{
			Id:          1,
			VersionCode: 1,
			VersionName: "ljdlkadj",
			Type:        1,
		},
	}
	rows := sqlmock.NewRows([]string{"id", "version_code", "version_name", "type"}).
		AddRow(mockVersionApp[0].Id, mockVersionApp[0].VersionCode, mockVersionApp[0].VersionName, mockVersionApp[0].Type)

	query := `
	SELECT
		\*\
	FROM
		version_apps
	WHERE
		version_apps.type = \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := VersionAppRepo.NewVersionAPPRepositoryRepository(db)

	anArticle, err := a.GetAll(context.TODO(), 1)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 1)
}
func TestGetAllErrorQuery(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	mockVersionApp := []models.VersionApp{
		models.VersionApp{
			Id:          1,
			VersionCode: 1,
			VersionName: "ljdlkadj",
			Type:        1,
		},
	}
	rows := sqlmock.NewRows([]string{"id", "version_code", "version_name", "type"}).
		AddRow(mockVersionApp[0].Id, mockVersionApp[0].VersionCode, mockVersionApp[0].VersionName, mockVersionApp[0].Type)

	query := `
	SELECT
		\*\
	FROM
		version_apps
	WHERE
		version_apps.type = \?asdasd`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := VersionAppRepo.NewVersionAPPRepositoryRepository(db)

	anArticle, err := a.GetAll(context.TODO(), 1)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
