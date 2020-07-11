package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/models"
	FacilitiesRepo "github.com/service/facilities/repository"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)
var (
	imagePath = "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Facilities/8941695193938718058.jpg"
	mockFacilities = []models.Facilities{
		models.Facilities{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			FacilityName: "test facilities",
			IsNumerable:  1,
			FacilityIcon: &imagePath,
		},
		models.Facilities{
			Id:           2,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			FacilityName: "test facilities 2",
			IsNumerable:  2,
			FacilityIcon: &imagePath,
		},
	}
)
func TestCount(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(len(mockFacilities))

	query := `SELECT count\(\*\) AS count FROM facilities WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := FacilitiesRepo.NewFacilityRepository(db)

	res, err := a.GetCount(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, res, 2, "")
}
func TestCountFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow("test")

	query := `SELECT count\(\*\) AS count FROM facilities WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := FacilitiesRepo.NewFacilityRepository(db)

	_, err = a.GetCount(context.TODO())
	assert.Error(t, err)
}
func TestList(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","facility_name","is_numerable","facility_icon"}).
		AddRow(mockFacilities[0].Id, mockFacilities[0].CreatedBy,mockFacilities[0].CreatedDate,mockFacilities[0].ModifiedBy,
			mockFacilities[0].ModifiedDate,mockFacilities[0].DeletedBy,mockFacilities[0].DeletedDate,mockFacilities[0].IsDeleted,
			mockFacilities[0].IsActive,mockFacilities[0].FacilityName,mockFacilities[0].IsNumerable,mockFacilities[0].FacilityIcon).
		AddRow(mockFacilities[1].Id, mockFacilities[1].CreatedBy,mockFacilities[1].CreatedDate,mockFacilities[1].ModifiedBy,
			mockFacilities[1].ModifiedDate,mockFacilities[1].DeletedBy,mockFacilities[1].DeletedDate,mockFacilities[1].IsDeleted,
			mockFacilities[1].IsActive,mockFacilities[1].FacilityName,mockFacilities[1].IsNumerable,mockFacilities[1].FacilityIcon)

	query := `SELECT                                         \*\                                               FROM facilities WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := FacilitiesRepo.NewFacilityRepository(db)

	anArticle, err := a.List(context.TODO())
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestListError(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","facility_name","is_numerable","facility_icon"}).
		AddRow(mockFacilities[0].Id, mockFacilities[0].CreatedBy,mockFacilities[0].CreatedDate,mockFacilities[0].ModifiedBy,
			mockFacilities[0].ModifiedDate,mockFacilities[0].DeletedBy,mockFacilities[0].DeletedDate,mockFacilities[0].IsDeleted,
			mockFacilities[0].IsActive,mockFacilities[0].FacilityName,mockFacilities[0].IsNumerable,mockFacilities[0].FacilityIcon).
		AddRow(mockFacilities[1].Id, mockFacilities[1].CreatedBy,mockFacilities[1].CreatedDate,mockFacilities[1].ModifiedBy,
			mockFacilities[1].ModifiedDate,mockFacilities[1].DeletedBy,mockFacilities[1].DeletedDate,mockFacilities[1].IsDeleted,
			mockFacilities[1].ModifiedBy,mockFacilities[1].FacilityName,mockFacilities[1].IsNumerable,mockFacilities[1].FacilityIcon)


	query := `SELECT                                         \*\                                               FROM facilities WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := FacilitiesRepo.NewFacilityRepository(db)

	anArticle, err := a.List(context.TODO())
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestFetchWithPagination(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","facility_name","is_numerable","facility_icon"}).
		AddRow(mockFacilities[0].Id, mockFacilities[0].CreatedBy,mockFacilities[0].CreatedDate,mockFacilities[0].ModifiedBy,
			mockFacilities[0].ModifiedDate,mockFacilities[0].DeletedBy,mockFacilities[0].DeletedDate,mockFacilities[0].IsDeleted,
			mockFacilities[0].IsActive,mockFacilities[0].FacilityName,mockFacilities[0].IsNumerable,mockFacilities[0].FacilityIcon).
		AddRow(mockFacilities[1].Id, mockFacilities[1].CreatedBy,mockFacilities[1].CreatedDate,mockFacilities[1].ModifiedBy,
			mockFacilities[1].ModifiedDate,mockFacilities[1].DeletedBy,mockFacilities[1].DeletedDate,mockFacilities[1].IsDeleted,
			mockFacilities[1].IsActive,mockFacilities[1].FacilityName,mockFacilities[1].IsNumerable,mockFacilities[1].FacilityIcon)


	query := `SELECT \*\ FROM facilities where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := FacilitiesRepo.NewFacilityRepository(db)

	limit := 10
	offset := 0
	anArticle, err := a.Fetch(context.TODO(), limit,offset)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestFetchWithoutPagination(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","facility_name","is_numerable","facility_icon"}).
		AddRow(mockFacilities[0].Id, mockFacilities[0].CreatedBy,mockFacilities[0].CreatedDate,mockFacilities[0].ModifiedBy,
			mockFacilities[0].ModifiedDate,mockFacilities[0].DeletedBy,mockFacilities[0].DeletedDate,mockFacilities[0].IsDeleted,
			mockFacilities[0].IsActive,mockFacilities[0].FacilityName,mockFacilities[0].IsNumerable,mockFacilities[0].FacilityIcon).
		AddRow(mockFacilities[1].Id, mockFacilities[1].CreatedBy,mockFacilities[1].CreatedDate,mockFacilities[1].ModifiedBy,
			mockFacilities[1].ModifiedDate,mockFacilities[1].DeletedBy,mockFacilities[1].DeletedDate,mockFacilities[1].IsDeleted,
			mockFacilities[1].IsActive,mockFacilities[1].FacilityName,mockFacilities[1].IsNumerable,mockFacilities[1].FacilityIcon)


	query := `SELECT \*\ FROM facilities where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := FacilitiesRepo.NewFacilityRepository(db)

	//limit := 10
	//offset := 0
	anArticle, err := a.Fetch(context.TODO(), 0,0)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestFetchWithPaginationErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","facility_name","is_numerable","facility_icon"}).
		AddRow(mockFacilities[0].Id, mockFacilities[0].CreatedBy,mockFacilities[0].CreatedDate,mockFacilities[0].ModifiedBy,
			mockFacilities[0].ModifiedDate,mockFacilities[0].DeletedBy,mockFacilities[0].DeletedDate,mockFacilities[0].IsDeleted,
			mockFacilities[0].IsActive,mockFacilities[0].FacilityName,mockFacilities[0].IsNumerable,mockFacilities[0].FacilityIcon).
		AddRow(mockFacilities[1].Id, mockFacilities[1].CreatedBy,mockFacilities[1].CreatedDate,mockFacilities[1].ModifiedBy,
			mockFacilities[1].ModifiedDate,mockFacilities[1].DeletedBy,mockFacilities[1].DeletedDate,mockFacilities[1].IsDeleted,
			mockFacilities[1].ModifiedBy,mockFacilities[1].FacilityName,mockFacilities[1].IsNumerable,mockFacilities[1].FacilityIcon)


	query := `SELECT \*\ FROM facilities where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := FacilitiesRepo.NewFacilityRepository(db)

	limit := 10
	offset := 0
	anArticle, err := a.Fetch(context.TODO(), limit,offset)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t,anArticle)
	//assert.Len(t, anArticle, 2)
}
func TestFetchWithoutPaginationErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","facility_name","is_numerable","facility_icon"}).
		AddRow(mockFacilities[0].Id, mockFacilities[0].CreatedBy,mockFacilities[0].CreatedDate,mockFacilities[0].ModifiedBy,
			mockFacilities[0].ModifiedDate,mockFacilities[0].DeletedBy,mockFacilities[0].DeletedDate,mockFacilities[0].IsDeleted,
			mockFacilities[0].IsActive,mockFacilities[0].FacilityName,mockFacilities[0].IsNumerable,mockFacilities[0].FacilityIcon).
		AddRow(mockFacilities[1].Id, mockFacilities[1].CreatedBy,mockFacilities[1].CreatedDate,mockFacilities[1].ModifiedBy,
			mockFacilities[1].ModifiedDate,mockFacilities[1].DeletedBy,mockFacilities[1].DeletedDate,mockFacilities[1].IsDeleted,
			mockFacilities[1].ModifiedBy,mockFacilities[1].FacilityName,mockFacilities[1].IsNumerable,mockFacilities[1].FacilityIcon)

	query := `SELECT \*\ FROM facilities where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := FacilitiesRepo.NewFacilityRepository(db)

	//limit := 10
	//offset := 0
	anArticle, err := a.Fetch(context.TODO(), 0,0)
	assert.Error(t, err)
	assert.Nil(t,anArticle)
}
func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()


	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","facility_name","is_numerable","facility_icon"}).
		AddRow(mockFacilities[0].Id, mockFacilities[0].CreatedBy,mockFacilities[0].CreatedDate,mockFacilities[0].ModifiedBy,
			mockFacilities[0].ModifiedDate,mockFacilities[0].DeletedBy,mockFacilities[0].DeletedDate,mockFacilities[0].IsDeleted,
			mockFacilities[0].IsActive,mockFacilities[0].FacilityName,mockFacilities[0].IsNumerable,mockFacilities[0].FacilityIcon)

	query := `SELECT \*\ FROM facilities WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := FacilitiesRepo.NewFacilityRepository(db)

	num := 1
	anArticle, err := a.GetById(context.TODO(), num)
	assert.NoError(t, err)
	assert.NotNil(t, anArticle)
}
func TestGetByIDNotfound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","facility_name","is_numerable","facility_icon"})

	query := `SELECT \*\ FROM facilities WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := FacilitiesRepo.NewFacilityRepository(db)

	num := 4
	anArticle, err := a.GetById(context.TODO(), num)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestGetByIDErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","facility_name","is_numerable","facility_icon"}).
		AddRow(mockFacilities[0].Id, mockFacilities[0].CreatedBy,mockFacilities[0].CreatedDate,mockFacilities[0].ModifiedBy,
			mockFacilities[0].ModifiedDate,mockFacilities[0].DeletedBy,mockFacilities[0].DeletedDate,mockFacilities[0].IsDeleted,
			mockFacilities[0].ModifiedBy,mockFacilities[0].FacilityName,mockFacilities[0].IsNumerable,mockFacilities[0].FacilityIcon)

	query := `SELECT \*\ FROM facilities WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := FacilitiesRepo.NewFacilityRepository(db)

	num := 1
	anArticle, err := a.GetById(context.TODO(), num)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
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

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","facility_name","is_numerable","facility_icon"}).
		AddRow(mockFacilities[0].Id, mockFacilities[0].CreatedBy,mockFacilities[0].CreatedDate,mockFacilities[0].ModifiedBy,
			mockFacilities[0].ModifiedDate,mockFacilities[0].DeletedBy,mockFacilities[0].DeletedDate,mockFacilities[0].IsDeleted,
			mockFacilities[0].IsActive,mockFacilities[0].FacilityName,mockFacilities[0].IsNumerable,mockFacilities[0].FacilityIcon).
		AddRow(mockFacilities[1].Id, mockFacilities[1].CreatedBy,mockFacilities[1].CreatedDate,mockFacilities[1].ModifiedBy,
			mockFacilities[1].ModifiedDate,mockFacilities[1].DeletedBy,mockFacilities[1].DeletedDate,mockFacilities[1].IsDeleted,
			mockFacilities[1].IsActive,mockFacilities[1].FacilityName,mockFacilities[1].IsNumerable,mockFacilities[1].FacilityIcon)

	query := `SELECT \*\ FROM facilities WHERE facility_name = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := FacilitiesRepo.NewFacilityRepository(db)

	FacilitiesName := "Test Facilities 2"
	anArticle, err := a.GetByName(context.TODO(), FacilitiesName)
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

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","facility_name","is_numerable","facility_icon"})

	query := `SELECT \*\ FROM facilities WHERE facility_name = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := FacilitiesRepo.NewFacilityRepository(db)

	FacilitiesName := "Test Facilities 2"
	anArticle, err := a.GetByName(context.TODO(), FacilitiesName)
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

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","facility_name","is_numerable","facility_icon"}).
		AddRow(mockFacilities[0].Id, mockFacilities[0].CreatedBy,mockFacilities[0].CreatedDate,mockFacilities[0].ModifiedBy,
			mockFacilities[0].ModifiedDate,mockFacilities[0].DeletedBy,mockFacilities[0].DeletedDate,mockFacilities[0].IsDeleted,
			mockFacilities[0].IsActive,mockFacilities[0].FacilityName,mockFacilities[0].IsNumerable,mockFacilities[0].FacilityIcon).
		AddRow(mockFacilities[1].Id, mockFacilities[1].CreatedBy,mockFacilities[1].CreatedDate,mockFacilities[1].ModifiedBy,
			mockFacilities[1].ModifiedDate,mockFacilities[1].DeletedBy,mockFacilities[1].DeletedDate,mockFacilities[1].IsDeleted,
			mockFacilities[1].IsActive,mockFacilities[1].FacilityName,mockFacilities[1].IsNumerable,mockFacilities[1].FacilityIcon)

	query := `SELECT \*\ FROM facilities WHERE Facilities_name = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := FacilitiesRepo.NewFacilityRepository(db)

	FacilitiesName := "Test Facilities 2"
	anArticle, err := a.GetByName(context.TODO(), FacilitiesName)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
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

	query := "UPDATE facilities SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := FacilitiesRepo.NewFacilityRepository(db)

	err = a.Delete(context.TODO(), id,deletedBy)
	assert.NoError(t, err)
}
func TestDeleteErrorExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "UPDATE facilities SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,id,id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := FacilitiesRepo.NewFacilityRepository(db)

	err = a.Delete(context.TODO(), id,deletedBy)
	assert.Error(t, err)
}
func TestInsert(t *testing.T) {
	user := "test"
	now := time.Now()
	a := mockFacilities[0]
	a.ModifiedBy = &user
	a.ModifiedDate = &now
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "INSERT facilities SET created_by=\\? , created_date=\\? , modified_by=\\?, modified_date=\\? , 				deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? , facility_name=\\?,  is_numerable=\\?	,facility_icon=\\? "
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.FacilityName,
		a.IsNumerable,a.FacilityIcon).WillReturnResult(sqlmock.NewResult(1, 1))

	i := FacilitiesRepo.NewFacilityRepository(db)

	id, err := i.Insert(context.TODO(), &a)
	assert.NoError(t, err)
	assert.Equal(t, *id, a.Id)
}
func TestInsertErrorExec(t *testing.T) {
	user := "test"
	now := time.Now()
	a := mockFacilities[0]
	a.ModifiedBy = &user
	a.ModifiedDate = &now
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "INSERT facilities SET created_by=\\? , created_date=\\? , modified_by=\\?, modified_date=\\? , 				deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? , facility_name=\\?,  is_numerable=\\?	,facility_icon=\\? "
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.FacilityName,
		a.IsNumerable,a.FacilityIcon,a.FacilityIcon).WillReturnResult(sqlmock.NewResult(1, 1))

	i := FacilitiesRepo.NewFacilityRepository(db)

	_, err = i.Insert(context.TODO(), &a)
	assert.Error(t, err)
}
func TestUpdate(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	ar := mockFacilities[0]
	ar.ModifiedBy = &modifyBy
	ar.ModifiedDate = &now

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE facilities set modified_by=\?, modified_date=\? ,facility_name=\?,is_numerable=\?,facility_icon=\? WHERE id = \?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.ModifiedBy, ar.ModifiedDate, ar.FacilityName, ar.IsNumerable,ar.FacilityIcon, ar.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	a := FacilitiesRepo.NewFacilityRepository(db)

	err = a.Update(context.TODO(), &ar)
	assert.NoError(t, err)
	assert.Nil(t,err)
}
func TestUpdateErrorExec(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	ar := mockFacilities[0]
	ar.ModifiedDate = &now
	ar.ModifiedBy = &modifyBy

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE facilities set modified_by=\?, modified_date=\? ,facility_name=\?,is_numerable=\?,facility_icon=\? WHERE id = \?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.ModifiedBy, ar.ModifiedDate, ar.FacilityName, ar.IsNumerable,ar.FacilityIcon, ar.Id,ar.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	a := FacilitiesRepo.NewFacilityRepository(db)

	err = a.Update(context.TODO(), &ar)
	assert.Error(t, err)
}
