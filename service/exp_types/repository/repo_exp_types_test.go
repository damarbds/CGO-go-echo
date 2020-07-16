package repository_test

import (
	"context"
	"testing"

	"github.com/models"
	ExpTypeRepo "github.com/service/exp_types/repository"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	mockExpType = []models.ExpTypeObject{
		models.ExpTypeObject{
			ExpTypeID:   1,
			ExpTypeName: "test1",
			ExpTypeIcon: "dasdasdasd",
		},
		models.ExpTypeObject{
			ExpTypeID:   2,
			ExpTypeName: "test2",
			ExpTypeIcon: "dasdasdasd",
		},
	}
)

func TestGetExpTypes(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"exp_type_id", "exp_type_name", "exp_type_icon"}).
		AddRow(mockExpType[0].ExpTypeID, mockExpType[0].ExpTypeName, mockExpType[0].ExpTypeIcon).
		AddRow(mockExpType[1].ExpTypeID, mockExpType[1].ExpTypeName, mockExpType[1].ExpTypeIcon)

	query := `
	SELECT
		id AS exp_type_id,
		exp_type_name,
		COALESCE\(exp_type_icon,""\) AS exp_type_icon
	FROM
		experience_types
	WHERE
		is_active = 1
		AND is_deleted = 0`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExpTypeRepo.NewExpTypeRepository(db)
	//
	//expId := mockExperienceIncludeJoin[0].ExpId
	anArticle, err := a.GetExpTypes(context.TODO())
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestGetExpTypesErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"exp_type_id", "exp_type_name", "exp_type_icon"}).
		AddRow(mockExpType[0].ExpTypeID, mockExpType[0].ExpTypeName, mockExpType[0].ExpTypeIcon).
		AddRow(mockExpType[1].ExpTypeID, mockExpType[1].ExpTypeName, mockExpType[1].ExpTypeIcon)

	query := `
	SELECT
		id AS exp_type_id,
		exp_type_name,
		COALESCE\(exp_type_icon,""\) AS exp_type_icon
	FROM
		experience_types
	WHERE
		is_active = 1
		AND is_deleted = 0asdasdasd`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExpTypeRepo.NewExpTypeRepository(db)
	//
	//expId := mockExperienceIncludeJoin[0].ExpId
	_, err = a.GetExpTypes(context.TODO())
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	//assert.Nil(t, anArticle)
	//assert.Len(t, anArticle, 2)
}
func TestGetByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"exp_type_id", "exp_type_name", "exp_type_icon"}).
		AddRow(mockExpType[0].ExpTypeID, mockExpType[0].ExpTypeName, mockExpType[0].ExpTypeIcon).
		AddRow(mockExpType[1].ExpTypeID, mockExpType[1].ExpTypeName, mockExpType[1].ExpTypeIcon)

	query := `
	SELECT
		id AS exp_type_id,
		exp_type_name,
		COALESCE\(exp_type_icon,""\) AS exp_type_icon
	FROM
		experience_types
	WHERE
		is_active = 1
		AND is_deleted = 0 AND exp_type_name = \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExpTypeRepo.NewExpTypeRepository(db)

	expTypeName := "test1"
	anArticle, err := a.GetByName(context.TODO(), expTypeName)
	assert.NoError(t, err)
	assert.NotNil(t, anArticle)
}
func TestGetByNameNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"exp_type_id", "exp_type_name", "exp_type_icon"})

	query := `
	SELECT
		id AS exp_type_id,
		exp_type_name,
		COALESCE\(exp_type_icon,""\) AS exp_type_icon
	FROM
		experience_types
	WHERE
		is_active = 1
		AND is_deleted = 0 AND exp_type_name = \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExpTypeRepo.NewExpTypeRepository(db)

	expTypeName := "test1"
	anArticle, err := a.GetByName(context.TODO(), expTypeName)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestGetByNameErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"exp_type_id", "exp_type_name", "exp_type_icon"}).
		AddRow(mockExpType[0].ExpTypeName, mockExpType[0].ExpTypeName, mockExpType[0].ExpTypeIcon).
		AddRow(mockExpType[1].ExpTypeName, mockExpType[1].ExpTypeName, mockExpType[1].ExpTypeIcon)

	query := `
	SELECT
		id AS exp_type_id,
		exp_type_name,
		COALESCE\(exp_type_icon,""\) AS exp_type_icon
	FROM
		experience_types
	WHERE
		is_active = 1
		AND is_deleted = 0 AND exp_type_name = \?adadssad`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExpTypeRepo.NewExpTypeRepository(db)

	expTypeName := "test1"
	anArticle, err := a.GetByName(context.TODO(), expTypeName)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
