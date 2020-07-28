package repository_test

import (
	"context"
	"testing"
	"time"

	NotificationRepo "github.com/misc/notif/repository"
	"github.com/models"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetByMerchantID(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	mockNotification := []models.Notification{
		models.Notification{
			Id:           "asdsadasdqqweq",
			CreatedBy:    "Test ",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			MerchantId:   "zxcxzcxzc",
			Type:         1,
			Title:        "qweqweqw",
			Desc:         "zxczxccxz",
		},
		models.Notification{
			Id:           "qweqweqweqwe",
			CreatedBy:    "Test ",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			MerchantId:   "zxcxzcxzc",
			Type:         1,
			Title:        "qweqweqw",
			Desc:         "zxczxccxz",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "merchant_id", "type","title","desc"}).
		AddRow(mockNotification[0].Id, mockNotification[0].CreatedBy, mockNotification[0].CreatedDate, mockNotification[0].ModifiedBy,
			mockNotification[0].ModifiedDate, mockNotification[0].DeletedBy, mockNotification[0].DeletedDate, mockNotification[0].IsDeleted,
			mockNotification[0].IsActive, mockNotification[0].MerchantId, mockNotification[0].Type, mockNotification[0].Title,
		 	mockNotification[0].Desc).
		AddRow(mockNotification[0].Id, mockNotification[0].CreatedBy, mockNotification[0].CreatedDate, mockNotification[0].ModifiedBy,
		mockNotification[0].ModifiedDate, mockNotification[0].DeletedBy, mockNotification[0].DeletedDate, mockNotification[0].IsDeleted,
		mockNotification[0].IsActive, mockNotification[0].MerchantId, mockNotification[0].Type, mockNotification[0].Title,
		mockNotification[0].Desc)

	query := `
	SELECT
		\*\
	FROM
		notifications
	WHERE
		merchant_id = \?
		AND is_deleted = 0
		AND is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := NotificationRepo.NewNotifRepository(db)

	anArticle, err := a.GetByMerchantID(context.TODO(),mockNotification[0].MerchantId)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestGetByMerchantIDError(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	mockNotification := []models.Notification{
		models.Notification{
			Id:           "asdsadasdqqweq",
			CreatedBy:    "Test ",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			MerchantId:   "zxcxzcxzc",
			Type:         1,
			Title:        "qweqweqw",
			Desc:         "zxczxccxz",
		},
		models.Notification{
			Id:           "qweqweqweqwe",
			CreatedBy:    "Test ",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			MerchantId:   "zxcxzcxzc",
			Type:         1,
			Title:        "qweqweqw",
			Desc:         "zxczxccxz",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "merchant_id", "type","title","desc"}).
		AddRow(mockNotification[0].Id, mockNotification[0].CreatedBy, mockNotification[0].CreatedDate, mockNotification[0].ModifiedBy,
			mockNotification[0].ModifiedDate, mockNotification[0].DeletedBy, mockNotification[0].DeletedDate, mockNotification[0].IsDeleted,
			mockNotification[0].ModifiedDate, mockNotification[0].MerchantId, mockNotification[0].Type, mockNotification[0].Title,
			mockNotification[0].Desc).
		AddRow(mockNotification[0].Id, mockNotification[0].CreatedBy, mockNotification[0].CreatedDate, mockNotification[0].ModifiedBy,
			mockNotification[0].ModifiedDate, mockNotification[0].DeletedBy, mockNotification[0].DeletedDate, mockNotification[0].IsDeleted,
			mockNotification[0].IsActive, mockNotification[0].MerchantId, mockNotification[0].Type, mockNotification[0].Title,
			mockNotification[0].Desc)

	query := `
	SELECT
		\*\
	FROM
		notifications
	WHERE
		merchant_id = \?
		AND is_deleted = 0
		AND is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := NotificationRepo.NewNotifRepository(db)

	anArticle, err := a.GetByMerchantID(context.TODO(),mockNotification[0].MerchantId)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestInsert(t *testing.T) {
	//user := "test"
	//now := time.Now()
	a := models.Notification{
		Id:           "qweqweqweqwe",
		CreatedBy:    "Test ",
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		MerchantId:   "zxcxzcxzc",
		Type:         1,
		Title:        "qweqweqw",
		Desc:         "zxczxccxz",
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT notifications SET id=\? , created_by=\? , created_date=\? , modified_by=\?, modified_date=\? ,
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , merchant_id=\?,type=\? , title=\? ,notifications.desc=\?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.MerchantId,a.Type,a.Title,a.Desc).WillReturnResult(sqlmock.NewResult(1, 1))

	i := NotificationRepo.NewNotifRepository(db)

	 err = i.Insert(context.TODO(), a)
	assert.NoError(t, err)
	//assert.Equal(t, *id, a.Id)
}
func TestInsertErrorExec(t *testing.T) {
	a := models.Notification{
		Id:           "qweqweqweqwe",
		CreatedBy:    "Test ",
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		MerchantId:   "zxcxzcxzc",
		Type:         1,
		Title:        "qweqweqw",
		Desc:         "zxczxccxz",
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT notifications SET id=\? , created_by=\? , created_date=\? , modified_by=\?, modified_date=\? ,
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , merchant_id=\?,type=\? , title=\? ,notifications.desc=\?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.MerchantId,a.Type,a.Title,a.Desc,a.Id).WillReturnResult(sqlmock.NewResult(1, 1))

	i := NotificationRepo.NewNotifRepository(db)

	err = i.Insert(context.TODO(), a)
	assert.Error(t, err)
}
func TestInsertErrorExecQuery(t *testing.T) {
	a := models.Notification{
		Id:           "qweqweqweqwe",
		CreatedBy:    "Test ",
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		MerchantId:   "zxcxzcxzc",
		Type:         1,
		Title:        "qweqweqw",
		Desc:         "zxczxccxz",
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT notifications SET id=\? , cdasdasdassdasdasdasdsaated_by=\? , created_date=\? , modified_by=\?, modified_date=\? ,
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , merchant_id=\?,type=\? , title=\? ,notifications.desc=\?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.MerchantId,a.Type,a.Title,a.Desc,a.Id).WillReturnResult(sqlmock.NewResult(1, 1))

	i := NotificationRepo.NewNotifRepository(db)

	err = i.Insert(context.TODO(), a)
	assert.Error(t, err)
}

