package repository_test
//
//import (
//	"context"
//	"github.com/models"
//	AdminRepo "github.com/auth/admin/repository"
//	"testing"
//	"time"
//
//	"github.com/stretchr/testify/assert"
//	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
//)
//var(
//	mockAdmin = []models.Admin{
//		models.Admin{
//			Id:           "asdasdasdsa",
//			CreatedBy:    "Test",
//			CreatedDate:  time.Now(),
//			ModifiedBy:   nil,
//			ModifiedDate: nil,
//			DeletedBy:    nil,
//			DeletedDate:  nil,
//			IsDeleted:    0,
//			IsActive:     1,
//			Name:         "AdminTest",
//			Email:        "AdminTest@gmail.com",
//		},
//		models.Admin{
//			Id:           "asdasdasdsa",
//			CreatedBy:    "Test",
//			CreatedDate:  time.Now(),
//			ModifiedBy:   nil,
//			ModifiedDate: nil,
//			DeletedBy:    nil,
//			DeletedDate:  nil,
//			IsDeleted:    0,
//			IsActive:     1,
//			Name:         "AdminTest2",
//			Email:        "AdminTest2@gmail.com",
//		},
//	}
//)
//func TestGetByID(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//
//	//defer func() {
//	//	err = db.Close()
//	//	require.NoError(t, err)
//	//}()
//
//	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
//		"deleted_date","is_deleted","is_active","name","email"}).
//		AddRow(mockAdmin[0].Id, mockAdmin[0].CreatedBy,mockAdmin[0].CreatedDate,mockAdmin[0].ModifiedBy,
//			mockAdmin[0].ModifiedDate,mockAdmin[0].DeletedBy,mockAdmin[0].DeletedDate,mockAdmin[0].IsDeleted,
//			mockAdmin[0].IsActive,mockAdmin[0].Name,mockAdmin[0].AdminIcon)
//
//	query := `SELECT \*\ FROM Admins WHERE id = \\?`
//
//	mock.ExpectQuery(query).WillReturnRows(rows)
//	a := AdminRepo.NewadminRepository(db)
//
//	num := 1
//	anArticle, err := a.GetByID(context.TODO(), mockAdmin[0].Id)
//	assert.NoError(t, err)
//	assert.NotNil(t, anArticle)
//}
//func TestGetByIDNotfound(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//
//	//defer func() {
//	//	err = db.Close()
//	//	require.NoError(t, err)
//	//}()
//
//	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
//		"deleted_date","is_deleted","is_active","Admin_name","Admin_icon"})
//
//	query := `SELECT \*\ FROM Admins WHERE id = \\?`
//
//	mock.ExpectQuery(query).WillReturnRows(rows)
//	a := AdminRepo.NewAdminRepository(db)
//
//	num := 4
//	anArticle, err := a.GetById(context.TODO(), num)
//	assert.Error(t, err)
//	assert.Nil(t, anArticle)
//}
//func TestGetByIDErrorFetch(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//
//	//defer func() {
//	//	err = db.Close()
//	//	require.NoError(t, err)
//	//}()
//
//	mockAdmin := []models.Admin{
//		models.Admin{
//			Id:           1,
//			CreatedBy:    "test",
//			CreatedDate:  time.Now(),
//			ModifiedBy:   nil,
//			ModifiedDate: nil,
//			DeletedBy:    nil,
//			DeletedDate:  nil,
//			IsDeleted:    0,
//			IsActive:     1,
//			AdminName:  "Test Admin 1",
//			AdminIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Admin/8941695193938718058.jpg",
//		},
//		models.Admin{
//			Id:           2,
//			CreatedBy:    "test",
//			CreatedDate:  time.Now(),
//			ModifiedBy:   nil,
//			ModifiedDate: nil,
//			DeletedBy:    nil,
//			DeletedDate:  nil,
//			IsDeleted:    0,
//			IsActive:     1,
//			AdminName:  "Test Admin 2",
//			AdminIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Admin/8941695193938718058.jpg",
//		},
//	}
//	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
//		"deleted_date","is_deleted","is_active","Admin_name","Admin_icon"}).
//		AddRow(mockAdmin[0].Id, mockAdmin[0].CreatedBy,mockAdmin[0].CreatedDate,mockAdmin[0].ModifiedBy,
//			mockAdmin[0].ModifiedDate,mockAdmin[0].DeletedBy,mockAdmin[0].DeletedDate,mockAdmin[0].IsDeleted,
//			mockAdmin[0].AdminName,mockAdmin[0].AdminName,mockAdmin[0].AdminIcon)
//
//	query := `SELECT \*\ FROM Admins WHERE id = \\?`
//
//	mock.ExpectQuery(query).WillReturnRows(rows)
//	a := AdminRepo.NewAdminRepository(db)
//
//	num := 1
//	anArticle, err := a.GetById(context.TODO(), num)
//	assert.Error(t, err)
//	assert.Nil(t, anArticle)
//}
//func TestGetByAdminEmail(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//
//	//defer func() {
//	//	err = db.Close()
//	//	require.NoError(t, err)
//	//}()
//
//	mockAdmin := []models.Admin{
//		models.Admin{
//			Id:           1,
//			CreatedBy:    "test",
//			CreatedDate:  time.Now(),
//			ModifiedBy:   nil,
//			ModifiedDate: nil,
//			DeletedBy:    nil,
//			DeletedDate:  nil,
//			IsDeleted:    0,
//			IsActive:     1,
//			AdminName:  "Test Admin 1",
//			AdminIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Admin/8941695193938718058.jpg",
//		},
//		models.Admin{
//			Id:           1,
//			CreatedBy:    "test",
//			CreatedDate:  time.Now(),
//			ModifiedBy:   nil,
//			ModifiedDate: nil,
//			DeletedBy:    nil,
//			DeletedDate:  nil,
//			IsDeleted:    0,
//			IsActive:     1,
//			AdminName:  "Test Admin 2",
//			AdminIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Admin/8941695193938718058.jpg",
//		},
//	}
//	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
//		"deleted_date","is_deleted","is_active","Admin_name","Admin_icon"}).
//		AddRow(mockAdmin[0].Id, mockAdmin[0].CreatedBy,mockAdmin[0].CreatedDate,mockAdmin[0].ModifiedBy,
//			mockAdmin[0].ModifiedDate,mockAdmin[0].DeletedBy,mockAdmin[0].DeletedDate,mockAdmin[0].IsDeleted,
//			mockAdmin[0].IsActive,mockAdmin[0].AdminName,mockAdmin[0].AdminIcon).
//		AddRow(mockAdmin[1].Id, mockAdmin[1].CreatedBy,mockAdmin[1].CreatedDate,mockAdmin[1].ModifiedBy,
//			mockAdmin[1].ModifiedDate,mockAdmin[1].DeletedBy,mockAdmin[1].DeletedDate,mockAdmin[1].IsDeleted,
//			mockAdmin[1].IsActive,mockAdmin[1].AdminName,mockAdmin[1].AdminIcon)
//
//	query := `SELECT \*\ FROM Admins WHERE Admin_name = \\?`
//
//	mock.ExpectQuery(query).WillReturnRows(rows)
//	a := AdminRepo.NewAdminRepository(db)
//
//	AdminName := "Test Admin 2"
//	anArticle, err := a.GetByName(context.TODO(), AdminName)
//	assert.NoError(t, err)
//	assert.NotNil(t, anArticle)
//}
//func TestGetByAdminEmailNotFound(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//
//	//defer func() {
//	//	err = db.Close()
//	//	require.NoError(t, err)
//	//}()
//
//	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
//		"deleted_date","is_deleted","is_active","Admin_name","Admin_icon"})
//
//	query := `SELECT \*\ FROM Admins WHERE Admin_name = \\?`
//
//	mock.ExpectQuery(query).WillReturnRows(rows)
//	a := AdminRepo.NewAdminRepository(db)
//
//	AdminName := "Test Admin 2"
//	anArticle, err := a.GetByName(context.TODO(), AdminName)
//	assert.Error(t, err)
//	assert.Nil(t, anArticle)
//}
//func TestGetByAdminEmailErrorFetch(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//
//	//defer func() {
//	//	err = db.Close()
//	//	require.NoError(t, err)
//	//}()
//
//	mockAdmin := []models.Admin{
//		models.Admin{
//			Id:           1,
//			CreatedBy:    "test",
//			CreatedDate:  time.Now(),
//			ModifiedBy:   nil,
//			ModifiedDate: nil,
//			DeletedBy:    nil,
//			DeletedDate:  nil,
//			IsDeleted:    0,
//			IsActive:     1,
//			AdminName:  "Test Admin 1",
//			AdminIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Admin/8941695193938718058.jpg",
//		},
//		models.Admin{
//			Id:           1,
//			CreatedBy:    "test",
//			CreatedDate:  time.Now(),
//			ModifiedBy:   nil,
//			ModifiedDate: nil,
//			DeletedBy:    nil,
//			DeletedDate:  nil,
//			IsDeleted:    0,
//			IsActive:     1,
//			AdminName:  "Test Admin 2",
//			AdminIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Admin/8941695193938718058.jpg",
//		},
//	}
//	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
//		"deleted_date","is_deleted","is_active","Admin_name","Admin_icon"}).
//		AddRow(mockAdmin[0].Id, mockAdmin[0].CreatedBy,mockAdmin[0].CreatedDate,mockAdmin[0].ModifiedBy,
//			mockAdmin[0].ModifiedDate,mockAdmin[0].DeletedBy,mockAdmin[0].DeletedDate,mockAdmin[0].IsDeleted,
//			mockAdmin[0].AdminName,mockAdmin[0].AdminName,mockAdmin[0].AdminIcon).
//		AddRow(mockAdmin[1].Id, mockAdmin[1].CreatedBy,mockAdmin[1].CreatedDate,mockAdmin[1].ModifiedBy,
//			mockAdmin[1].ModifiedDate,mockAdmin[1].DeletedBy,mockAdmin[1].DeletedDate,mockAdmin[1].IsDeleted,
//			mockAdmin[1].AdminName,mockAdmin[1].AdminName,mockAdmin[1].AdminIcon)
//
//	query := `SELECT \*\ FROM Admins WHERE Admin_name = \\?`
//
//	mock.ExpectQuery(query).WillReturnRows(rows)
//	a := AdminRepo.NewAdminRepository(db)
//
//	AdminName := "Test Admin 2"
//	anArticle, err := a.GetByName(context.TODO(), AdminName)
//	assert.Error(t, err)
//	assert.Nil(t, anArticle)
//}
//func TestDelete(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	//defer func() {
//	//	err = db.Close()
//	//	require.NoError(t, err)
//	//}()
//
//	query := "UPDATE Admins SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
//	id := 2
//	deletedBy := "test"
//	prep := mock.ExpectPrepare(query)
//	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,id).WillReturnResult(sqlmock.NewResult(2, 1))
//
//	a := AdminRepo.NewAdminRepository(db)
//
//	err = a.Delete(context.TODO(), id,deletedBy)
//	assert.NoError(t, err)
//}
//func TestDeleteErrorExec(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	//defer func() {
//	//	err = db.Close()
//	//	require.NoError(t, err)
//	//}()
//
//	query := "UPDATE Admins SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
//	id := 2
//	deletedBy := "test"
//	prep := mock.ExpectPrepare(query)
//	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,id,id).WillReturnResult(sqlmock.NewResult(2, 1))
//
//	a := AdminRepo.NewAdminRepository(db)
//
//	err = a.Delete(context.TODO(), id,deletedBy)
//	assert.Error(t, err)
//}
//func TestInsert(t *testing.T) {
//	user := "test"
//	now := time.Now()
//	a := models.Admin{
//		Id:           1,
//		CreatedBy:    user,
//		CreatedDate:  now,
//		ModifiedBy:   &user,
//		ModifiedDate: &now,
//		DeletedBy:    &user,
//		DeletedDate:  &now,
//		IsDeleted:    0,
//		IsActive:     0,
//		AdminName:  "test Admin 1",
//		AdminIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Admin/8941695193938718058.jpg",
//	}
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	//defer func() {
//	//	err = db.Close()
//	//	require.NoError(t, err)
//	//}()
//
//	query := "INSERT Admins SET created_by=\\? , created_date=\\? , modified_by=\\?, modified_date=\\? , 				deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? , Admin_name=\\?,  				Admin_icon=\\? "
//	prep := mock.ExpectPrepare(query)
//	prep.ExpectExec().WithArgs(a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.AdminName,
//		a.AdminIcon).WillReturnResult(sqlmock.NewResult(1, 1))
//
//	i := AdminRepo.NewAdminRepository(db)
//
//	id, err := i.Insert(context.TODO(), &a)
//	assert.NoError(t, err)
//	assert.Equal(t, *id, a.Id)
//}
//func TestInsertErrorExec(t *testing.T) {
//	user := "test"
//	now := time.Now()
//	a := models.Admin{
//		Id:           1,
//		CreatedBy:    user,
//		CreatedDate:  now,
//		ModifiedBy:   &user,
//		ModifiedDate: &now,
//		DeletedBy:    &user,
//		DeletedDate:  &now,
//		IsDeleted:    0,
//		IsActive:     0,
//		AdminName:  "test Admin 1",
//		AdminIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Admin/8941695193938718058.jpg",
//	}
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	//defer func() {
//	//	err = db.Close()
//	//	require.NoError(t, err)
//	//}()
//
//	query := "INSERT Admins SET created_by=\\? , created_date=\\? , modified_by=\\?, modified_date=\\? , 				deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? , Admin_name=\\?,  				Admin_icon=\\? "
//	prep := mock.ExpectPrepare(query)
//	prep.ExpectExec().WithArgs(a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.AdminName,
//		a.AdminIcon,a.Id).WillReturnResult(sqlmock.NewResult(1, 1))
//
//	i := AdminRepo.NewAdminRepository(db)
//
//	_, err = i.Insert(context.TODO(), &a)
//	assert.Error(t, err)
//}
//func TestUpdate(t *testing.T) {
//	now := time.Now()
//	modifyBy := "test"
//	ar := models.Admin{
//		Id:           1,
//		CreatedBy:    "",
//		CreatedDate:  time.Time{},
//		ModifiedBy:   &modifyBy,
//		ModifiedDate: &now,
//		DeletedBy:    nil,
//		DeletedDate:  nil,
//		IsDeleted:    0,
//		IsActive:     0,
//		AdminName:  "test Admin 1",
//		AdminIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Admin/8941695193938718058.jpg",
//	}
//
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	//defer func() {
//	//	err = db.Close()
//	//	require.NoError(t, err)
//	//}()
//
//	query := `UPDATE Admins set modified_by=\?, modified_date=\? ,Admin_name=\?,Admin_icon=\? WHERE id = \?`
//
//	prep := mock.ExpectPrepare(query)
//	prep.ExpectExec().WithArgs(ar.ModifiedBy, ar.ModifiedDate, ar.AdminName, ar.AdminIcon, ar.Id).
//		WillReturnResult(sqlmock.NewResult(1, 1))
//
//	a := AdminRepo.NewAdminRepository(db)
//
//	err = a.Update(context.TODO(), &ar)
//	assert.NoError(t, err)
//	assert.Nil(t,err)
//}
//func TestUpdateErrorExec(t *testing.T) {
//	now := time.Now()
//	modifyBy := "test"
//	ar := models.Admin{
//		Id:           1,
//		CreatedBy:    "",
//		CreatedDate:  time.Time{},
//		ModifiedBy:   &modifyBy,
//		ModifiedDate: &now,
//		DeletedBy:    nil,
//		DeletedDate:  nil,
//		IsDeleted:    0,
//		IsActive:     0,
//		AdminName:  "test Admin 1",
//		AdminIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Admin/8941695193938718058.jpg",
//	}
//
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	//defer func() {
//	//	err = db.Close()
//	//	require.NoError(t, err)
//	//}()
//
//	query := `UPDATE Admins set modified_by=\?, modified_date=\? ,Admin_name=\?,Admin_icon=\? WHERE id = \?`
//
//	prep := mock.ExpectPrepare(query)
//	prep.ExpectExec().WithArgs(ar.ModifiedBy, ar.ModifiedDate, ar.AdminName, ar.AdminIcon, ar.Id,ar.Id).
//		WillReturnResult(sqlmock.NewResult(1, 1))
//
//	a := AdminRepo.NewAdminRepository(db)
//
//	err = a.Update(context.TODO(), &ar)
//	assert.Error(t, err)
//}
