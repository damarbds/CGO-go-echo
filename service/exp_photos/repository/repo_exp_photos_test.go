package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/models"
	ExpPhotosRepo "github.com/service/exp_photos/repository"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)
var (
	imagePath = `[{"original":"https://cgostorage.blob.core.windows.net/cgo-storage/Experience/6569483590502428383.jpg","thumbnail":""}]`
	mockExpPhotos = []models.ExpPhotos{
		models.ExpPhotos{
			Id:             "adsasdasda",
			CreatedBy:      "Test 1",
			CreatedDate:    time.Now(),
			ModifiedBy:     nil,
			ModifiedDate:   nil,
			DeletedBy:      nil,
			DeletedDate:    nil,
			IsDeleted:      0,
			IsActive:       1,
			ExpPhotoFolder: "Facilities",
			ExpPhotoImage:  imagePath,
			ExpId:          "qweqwewq",
		},
		models.ExpPhotos{
			Id:             "zxcxzxczx",
			CreatedBy:      "Test 1",
			CreatedDate:    time.Now(),
			ModifiedBy:     nil,
			ModifiedDate:   nil,
			DeletedBy:      nil,
			DeletedDate:    nil,
			IsDeleted:      0,
			IsActive:       1,
			ExpPhotoFolder: "Include",
			ExpPhotoImage:  imagePath,
			ExpId:          "qweqwewq",
		},
	}
)

func TestGetByExperienceID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()


	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","exp_photo_folder","exp_photo_image","exp_id"}).
		AddRow(mockExpPhotos[0].Id, mockExpPhotos[0].CreatedBy,mockExpPhotos[0].CreatedDate,mockExpPhotos[0].ModifiedBy,
			mockExpPhotos[0].ModifiedDate,mockExpPhotos[0].DeletedBy,mockExpPhotos[0].DeletedDate,mockExpPhotos[0].IsDeleted,
			mockExpPhotos[0].IsActive,mockExpPhotos[0].ExpPhotoFolder,mockExpPhotos[0].ExpPhotoImage,mockExpPhotos[0].ExpId)

	query := `SELECT \*\ FROM exp_photos WHERE exp_id = \? AND is_deleted = 0 AND is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExpPhotosRepo.Newexp_photosRepository(db)

	num := mockExpPhotos[0].ExpId
	anArticle, err := a.GetByExperienceID(context.TODO(), num)
	assert.NoError(t, err)
	assert.NotNil(t, anArticle)
}
func TestGetByExperienceIDNotfound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","exp_photo_folder","exp_photo_image","exp_id"})

	query := `SELECT \*\ FROM exp_photos WHERE exp_id = \? AND is_deleted = 0 AND is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExpPhotosRepo.Newexp_photosRepository(db)

	num := mockExpPhotos[0].ExpId
	anArticle, err := a.GetByExperienceID(context.TODO(), num)
	//assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestGetByExperienceIDErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","exp_photo_folder","exp_photo_image","exp_id"}).
		AddRow(mockExpPhotos[0].Id, mockExpPhotos[0].CreatedBy,mockExpPhotos[0].CreatedDate,mockExpPhotos[0].ModifiedBy,
			mockExpPhotos[0].ModifiedDate,mockExpPhotos[0].DeletedBy,mockExpPhotos[0].DeletedDate,mockExpPhotos[0].IsDeleted,
			mockExpPhotos[0].ModifiedDate,mockExpPhotos[0].ExpPhotoFolder,mockExpPhotos[0].ExpPhotoImage,mockExpPhotos[0].ExpId)

	query := `SELECT \*\ FROM exp_photos WHERE exp_id = \? AND is_deleted = 0 AND is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExpPhotosRepo.Newexp_photosRepository(db)

	num := mockExpPhotos[0].ExpId
	anArticle, err := a.GetByExperienceID(context.TODO(), num)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestDeletes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	var ids = []string{"qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw"}
	query := `UPDATE exp_photos SET deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? WHERE exp_id=\?`
	for index, id := range ids {
		if index == 0 && index != (len(ids)-1) {
			query = query + ` AND \(id !=` + id
		} else if index == 0 && index == (len(ids)-1) {
			query = query + ` AND \(id !=` + id + ` \) `
		} else if index == (len(ids) - 1) {
			query = query + ` OR id !=` + id + ` \) `
		} else {
			query = query + ` OR id !=` + id
		}
	}

	expId := mockExpPhotos[0].ExpId
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,expId).WillReturnResult(sqlmock.NewResult(2, 1))

	a := ExpPhotosRepo.Newexp_photosRepository(db)

	err = a.Deletes(context.TODO(), ids,expId,deletedBy)
	assert.NoError(t, err)
}
func TestDeletesErrorExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	var ids = []string{"qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw"}
	query := `UPDATE exp_photos SET deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? WHERE exp_id=\?`
	for index, id := range ids {
		if index == 0 && index != (len(ids)-1) {
			query = query + ` AND \(id !=` + id
		} else if index == 0 && index == (len(ids)-1) {
			query = query + ` AND \(id !=` + id + ` \) `
		} else if index == (len(ids) - 1) {
			query = query + ` OR id !=` + id + ` \) `
		} else {
			query = query + ` OR id !=` + id
		}
	}

	expId := mockExpPhotos[0].ExpId
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,expId,expId).WillReturnResult(sqlmock.NewResult(2, 1))

	a := ExpPhotosRepo.Newexp_photosRepository(db)

	err = a.Deletes(context.TODO(), ids,expId,deletedBy)
	assert.Error(t, err)
}
func TestDeleteByExpId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE exp_photos SET deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? WHERE exp_id=\?`

	expId := mockExpPhotos[0].ExpId
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,expId).WillReturnResult(sqlmock.NewResult(2, 1))

	a := ExpPhotosRepo.Newexp_photosRepository(db)

	err = a.DeleteByExpId(context.TODO(), expId,deletedBy)
	assert.NoError(t, err)
}
func TestDeleteByExpIdErrorExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE exp_photos SET deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? WHERE exp_id=\?`

	expId := mockExpPhotos[0].ExpId
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,expId,expId).WillReturnResult(sqlmock.NewResult(2, 1))

	a := ExpPhotosRepo.Newexp_photosRepository(db)

	err = a.DeleteByExpId(context.TODO(), expId,deletedBy)
	assert.Error(t, err)
}
func TestInsert(t *testing.T) {
	user := "test"
	now := time.Now()
	a := mockExpPhotos[0]
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

	query := `INSERT exp_photos SET id=\? , created_by=\? , created_date=\? , modified_by=\?, modified_date=\? , 
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , exp_photo_folder=\?,exp_photo_image=\?,
				exp_id=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.ExpPhotoFolder,
		a.ExpPhotoImage, a.ExpId).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ExpPhotosRepo.Newexp_photosRepository(db)

	id, err := i.Insert(context.TODO(), &a)
	assert.NoError(t, err)
	assert.Equal(t, *id, a.Id)
}
func TestInsertErrorExec(t *testing.T) {
	user := "test"
	now := time.Now()
	a := mockExpPhotos[0]
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

	query := `INSERT exp_photos SET id=\? , created_by=\? , created_date=\? , modified_by=\?, modified_date=\? , 
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , exp_photo_folder=\?,exp_photo_image=\?,
				exp_id=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.ExpPhotoFolder,
		a.ExpPhotoImage, a.ExpId,a.ExpId).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ExpPhotosRepo.Newexp_photosRepository(db)

	_, err = i.Insert(context.TODO(), &a)
	assert.Error(t, err)
}
func TestUpdate(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	a := mockExpPhotos[0]
	a.ModifiedBy = &modifyBy
	a.ModifiedDate = &now

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE exp_photos SET modified_by=\?, modified_date=\? , 
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , exp_photo_folder=\?,exp_photo_image=\?,
				exp_id=\? WHERE id=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.ModifiedBy, time.Now(), nil, nil, 0, 1, a.ExpPhotoFolder,
		a.ExpPhotoImage, a.ExpId, a.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	u := ExpPhotosRepo.Newexp_photosRepository(db)

	_,err = u.Update(context.TODO(), &a)
	assert.NoError(t, err)
	assert.Nil(t,err)
}
func TestUpdateErrorExec(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	a := mockExpPhotos[0]
	a.ModifiedDate = &now
	a.ModifiedBy = &modifyBy

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE exp_photos SET modified_by=\?, modified_date=\? , 
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , exp_photo_folder=\?,exp_photo_image=\?,
				exp_id=\? WHERE id=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.ModifiedBy, time.Now(), nil, nil, 0, 1, a.ExpPhotoFolder,
		a.ExpPhotoImage, a.ExpId, a.Id,a.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	u := ExpPhotosRepo.Newexp_photosRepository(db)

	_,err = u.Update(context.TODO(), &a)
	assert.Error(t, err)
}
