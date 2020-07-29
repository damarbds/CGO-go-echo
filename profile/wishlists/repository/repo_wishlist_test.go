package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/models"
	WishlistRepo "github.com/profile/wishlists/repository"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)
var(
	mockWishlist = []models.Wishlist{
		models.Wishlist{
			Id:           "asdasdasdsad",
			CreatedBy:    "Test 1",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			TransId:      "asdsdas",
			ExpId:        "asdasd",
			UserId:       "zxczxcxzc",
		},
		models.Wishlist{
			Id:           "asdasdasdsadfddfdf",
			CreatedBy:    "Test 2",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			TransId:      "asdsdas",
			ExpId:        "asdasd",
			UserId:       "zxczxcxzc",
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
		AddRow(len(mockWishlist))

	query := `SELECT COUNT\(\*\) as count FROM wishlists WHERE user_id = \? AND is_deleted = 0 AND is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := WishlistRepo.NewWishListRepository(db)

	res, err := a.Count(context.TODO(),mockWishlist[0].UserId)
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

	query := `SELECT COUNT\(\*\) as count FROM wishlists WHERE user_id = \? AND is_deleted = 0 AND is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := WishlistRepo.NewWishListRepository(db)

	_, err = a.Count(context.TODO(),mockWishlist[0].UserId)
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
		"deleted_date", "is_deleted", "is_active", "trans_id", "exp_id","user_id"}).
		AddRow(mockWishlist[0].Id, mockWishlist[0].CreatedBy, mockWishlist[0].CreatedDate, mockWishlist[0].ModifiedBy,
			mockWishlist[0].ModifiedDate, mockWishlist[0].DeletedBy, mockWishlist[0].DeletedDate, mockWishlist[0].IsDeleted,
			mockWishlist[0].IsActive, mockWishlist[0].TransId, mockWishlist[0].ExpId,mockWishlist[0].UserId).
		AddRow(mockWishlist[1].Id, mockWishlist[1].CreatedBy, mockWishlist[1].CreatedDate, mockWishlist[1].ModifiedBy,
			mockWishlist[1].ModifiedDate, mockWishlist[1].DeletedBy, mockWishlist[1].DeletedDate, mockWishlist[1].IsDeleted,
			mockWishlist[1].IsActive, mockWishlist[1].TransId, mockWishlist[1].ExpId,mockWishlist[0].UserId)

	expId := mockWishlist[0].ExpId
	query := `SELECT \*\ FROM wishlists WHERE user_id = \? AND is_deleted = 0 AND is_active = 1`
	if expId != ""{
		query = query + ` AND exp_id = '` + expId + `' `
	}
	mock.ExpectQuery(query).WillReturnRows(rows)
	a := WishlistRepo.NewWishListRepository(db)

	anArticle, err := a.List(context.TODO(),mockWishlist[0].UserId,0,0,mockWishlist[0].ExpId)
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
		"deleted_date", "is_deleted", "is_active", "trans_id", "exp_id","user_id"}).
		AddRow(mockWishlist[0].Id, mockWishlist[0].CreatedBy, mockWishlist[0].CreatedDate, mockWishlist[0].ModifiedBy,
			mockWishlist[0].ModifiedDate, mockWishlist[0].DeletedBy, mockWishlist[0].DeletedDate, mockWishlist[0].IsDeleted,
			mockWishlist[0].IsActive, mockWishlist[0].TransId, mockWishlist[0].ExpId,mockWishlist[0].UserId).
		AddRow(mockWishlist[1].Id, mockWishlist[1].CreatedBy, mockWishlist[1].CreatedDate, mockWishlist[1].ModifiedBy,
			mockWishlist[1].ModifiedDate, mockWishlist[1].DeletedBy, mockWishlist[1].DeletedDate, mockWishlist[1].IsDeleted,
			mockWishlist[1].ModifiedDate, mockWishlist[1].TransId, mockWishlist[1].ExpId,mockWishlist[0].UserId)

	expId := mockWishlist[0].ExpId
	query := `SELECT \*\ FROM wishlists WHERE user_id = \? AND is_deleted = 0 AND is_active = 1`
	if expId != ""{
		query = query + ` AND exp_id = '` + expId + `' `
	}
	mock.ExpectQuery(query).WillReturnRows(rows)
	a := WishlistRepo.NewWishListRepository(db)

	anArticle, err := a.List(context.TODO(),mockWishlist[0].UserId,0,0,mockWishlist[0].ExpId)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestListWithPagination(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "trans_id", "exp_id","user_id"}).
		AddRow(mockWishlist[0].Id, mockWishlist[0].CreatedBy, mockWishlist[0].CreatedDate, mockWishlist[0].ModifiedBy,
			mockWishlist[0].ModifiedDate, mockWishlist[0].DeletedBy, mockWishlist[0].DeletedDate, mockWishlist[0].IsDeleted,
			mockWishlist[0].IsActive, mockWishlist[0].TransId, mockWishlist[0].ExpId,mockWishlist[0].UserId).
		AddRow(mockWishlist[1].Id, mockWishlist[1].CreatedBy, mockWishlist[1].CreatedDate, mockWishlist[1].ModifiedBy,
			mockWishlist[1].ModifiedDate, mockWishlist[1].DeletedBy, mockWishlist[1].DeletedDate, mockWishlist[1].IsDeleted,
			mockWishlist[1].IsActive, mockWishlist[1].TransId, mockWishlist[1].ExpId,mockWishlist[0].UserId)

	expId := mockWishlist[0].ExpId
	query := `SELECT \*\ FROM wishlists WHERE user_id = \? AND is_deleted = 0 AND is_active = 1`
	if expId != ""{
		query = query + ` AND exp_id = '` + expId + `' `
	}
	mock.ExpectQuery(query).WillReturnRows(rows)
	a := WishlistRepo.NewWishListRepository(db)

	anArticle, err := a.List(context.TODO(),mockWishlist[0].UserId,2,0,mockWishlist[0].ExpId)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestListWithPaginationError(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "trans_id", "exp_id","user_id"}).
		AddRow(mockWishlist[0].Id, mockWishlist[0].CreatedBy, mockWishlist[0].CreatedDate, mockWishlist[0].ModifiedBy,
			mockWishlist[0].ModifiedDate, mockWishlist[0].DeletedBy, mockWishlist[0].DeletedDate, mockWishlist[0].IsDeleted,
			mockWishlist[0].IsActive, mockWishlist[0].TransId, mockWishlist[0].ExpId,mockWishlist[0].UserId).
		AddRow(mockWishlist[1].Id, mockWishlist[1].CreatedBy, mockWishlist[1].CreatedDate, mockWishlist[1].ModifiedBy,
			mockWishlist[1].ModifiedDate, mockWishlist[1].DeletedBy, mockWishlist[1].DeletedDate, mockWishlist[1].IsDeleted,
			mockWishlist[1].ModifiedDate, mockWishlist[1].TransId, mockWishlist[1].ExpId,mockWishlist[0].UserId)

	expId := mockWishlist[0].ExpId
	query := `SELECT \*\ FROM wishlists WHERE user_id = \? AND is_deleted = 0 AND is_active = 1`
	if expId != ""{
		query = query + ` AND exp_id = '` + expId + `' `
	}
	mock.ExpectQuery(query).WillReturnRows(rows)
	a := WishlistRepo.NewWishListRepository(db)

	anArticle, err := a.List(context.TODO(),mockWishlist[0].UserId,2,0,mockWishlist[0].ExpId)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestGetByUserAndExpId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()


	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "trans_id", "exp_id","user_id"}).
		AddRow(mockWishlist[0].Id, mockWishlist[0].CreatedBy, mockWishlist[0].CreatedDate, mockWishlist[0].ModifiedBy,
			mockWishlist[0].ModifiedDate, mockWishlist[0].DeletedBy, mockWishlist[0].DeletedDate, mockWishlist[0].IsDeleted,
			mockWishlist[0].IsActive, mockWishlist[0].TransId, mockWishlist[0].ExpId,mockWishlist[0].UserId).
		AddRow(mockWishlist[1].Id, mockWishlist[1].CreatedBy, mockWishlist[1].CreatedDate, mockWishlist[1].ModifiedBy,
			mockWishlist[1].ModifiedDate, mockWishlist[1].DeletedBy, mockWishlist[1].DeletedDate, mockWishlist[1].IsDeleted,
			mockWishlist[1].IsActive, mockWishlist[1].TransId, mockWishlist[1].ExpId,mockWishlist[0].UserId)

	query := `SELECT \*\ FROM wishlists WHERE user_id = \? AND exp_id = \? AND is_deleted = 0 AND is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := WishlistRepo.NewWishListRepository(db)

	anArticle, err := a.GetByUserAndExpId(context.TODO(), mockWishlist[0].UserId,mockWishlist[0].ExpId,mockWishlist[0].TransId)
	assert.NoError(t, err)
	assert.NotNil(t, anArticle)
}
func TestGetByUserAndExpIdErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "trans_id", "exp_id","user_id"}).
		AddRow(mockWishlist[0].Id, mockWishlist[0].CreatedBy, mockWishlist[0].CreatedDate, mockWishlist[0].ModifiedBy,
			mockWishlist[0].ModifiedDate, mockWishlist[0].DeletedBy, mockWishlist[0].DeletedDate, mockWishlist[0].IsDeleted,
			mockWishlist[0].IsActive, mockWishlist[0].TransId, mockWishlist[0].ExpId,mockWishlist[0].UserId).
		AddRow(mockWishlist[1].Id, mockWishlist[1].CreatedBy, mockWishlist[1].CreatedDate, mockWishlist[1].ModifiedBy,
			mockWishlist[1].ModifiedDate, mockWishlist[1].DeletedBy, mockWishlist[1].DeletedDate, mockWishlist[1].IsDeleted,
			mockWishlist[1].ModifiedDate, mockWishlist[1].TransId, mockWishlist[1].ExpId,mockWishlist[0].UserId)

	query := `SELECT \*\ FROM wishlists WHERE user_id = \? AND exp_id = \? AND is_deleted = 0 AND is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := WishlistRepo.NewWishListRepository(db)

	anArticle, err := a.GetByUserAndExpId(context.TODO(), mockWishlist[0].UserId,mockWishlist[0].ExpId,mockWishlist[0].TransId)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestGetByUserAndTransId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()


	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "trans_id", "exp_id","user_id"}).
		AddRow(mockWishlist[0].Id, mockWishlist[0].CreatedBy, mockWishlist[0].CreatedDate, mockWishlist[0].ModifiedBy,
			mockWishlist[0].ModifiedDate, mockWishlist[0].DeletedBy, mockWishlist[0].DeletedDate, mockWishlist[0].IsDeleted,
			mockWishlist[0].IsActive, mockWishlist[0].TransId, mockWishlist[0].ExpId,mockWishlist[0].UserId).
		AddRow(mockWishlist[1].Id, mockWishlist[1].CreatedBy, mockWishlist[1].CreatedDate, mockWishlist[1].ModifiedBy,
			mockWishlist[1].ModifiedDate, mockWishlist[1].DeletedBy, mockWishlist[1].DeletedDate, mockWishlist[1].IsDeleted,
			mockWishlist[1].IsActive, mockWishlist[1].TransId, mockWishlist[1].ExpId,mockWishlist[0].UserId)

	query := `SELECT \*\ FROM wishlists WHERE user_id = \? AND trans_id = \? AND is_deleted = 0 AND is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := WishlistRepo.NewWishListRepository(db)

	anArticle, err := a.GetByUserAndExpId(context.TODO(), mockWishlist[0].UserId,"",mockWishlist[0].TransId)
	assert.NoError(t, err)
	assert.NotNil(t, anArticle)
}
func TestGetByUserAndTransIdErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "trans_id", "exp_id","user_id"}).
		AddRow(mockWishlist[0].Id, mockWishlist[0].CreatedBy, mockWishlist[0].CreatedDate, mockWishlist[0].ModifiedBy,
			mockWishlist[0].ModifiedDate, mockWishlist[0].DeletedBy, mockWishlist[0].DeletedDate, mockWishlist[0].IsDeleted,
			mockWishlist[0].IsActive, mockWishlist[0].TransId, mockWishlist[0].ExpId,mockWishlist[0].UserId).
		AddRow(mockWishlist[1].Id, mockWishlist[1].CreatedBy, mockWishlist[1].CreatedDate, mockWishlist[1].ModifiedBy,
			mockWishlist[1].ModifiedDate, mockWishlist[1].DeletedBy, mockWishlist[1].DeletedDate, mockWishlist[1].IsDeleted,
			mockWishlist[1].ModifiedDate, mockWishlist[1].TransId, mockWishlist[1].ExpId,mockWishlist[0].UserId)

	query := `SELECT \*\ FROM wishlists WHERE user_id = \? AND trans_id = \? AND is_deleted = 0 AND is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := WishlistRepo.NewWishListRepository(db)

	anArticle, err := a.GetByUserAndExpId(context.TODO(), mockWishlist[0].UserId,"",mockWishlist[0].TransId)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestDeleteByUserIdAndExpId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE wishlists SET deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? WHERE user_id =\? AND exp_id=\?`

	//id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,mockWishlist[0].UserId,mockWishlist[0].ExpId).WillReturnResult(sqlmock.NewResult(2, 1))

	a := WishlistRepo.NewWishListRepository(db)

	err = a.DeleteByUserIdAndExpIdORTransId(context.TODO(), mockWishlist[0].UserId,mockWishlist[0].ExpId,"", deletedBy)
	assert.NoError(t, err)
}
func TestDeleteByUserIdAndExpIdErrorExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE wishlists SET deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? WHERE user_id =\? AND exp_id=\?`

	//id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,mockWishlist[0].UserId,mockWishlist[0].ExpId,mockWishlist[0].ExpId).WillReturnResult(sqlmock.NewResult(2, 1))

	a := WishlistRepo.NewWishListRepository(db)

	err = a.DeleteByUserIdAndExpIdORTransId(context.TODO(), mockWishlist[0].UserId,mockWishlist[0].ExpId,"", deletedBy)

	assert.Error(t, err)
}
func TestDeleteByUserIdAndTransId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE wishlists SET deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? WHERE user_id =\? AND trans_id =\?`
	//id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,mockWishlist[0].UserId,mockWishlist[0].TransId).WillReturnResult(sqlmock.NewResult(2, 1))

	a := WishlistRepo.NewWishListRepository(db)

	err = a.DeleteByUserIdAndExpIdORTransId(context.TODO(), mockWishlist[0].UserId,"",mockWishlist[0].TransId, deletedBy)
	assert.NoError(t, err)
}
func TestDeleteByUserIdAndTransIdErrorExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE wishlists SET deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? WHERE user_id =\? AND trans_id =\?`
	//id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,mockWishlist[0].UserId,mockWishlist[0].TransId,mockWishlist[0].Id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := WishlistRepo.NewWishListRepository(db)

	err = a.DeleteByUserIdAndExpIdORTransId(context.TODO(), mockWishlist[0].UserId,"",mockWishlist[0].TransId, deletedBy)

	assert.Error(t, err)
}
func TestInsertByTransId(t *testing.T) {
	//user := "test"
	//now := time.Now()
	wl := mockWishlist[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	q := `INSERT wishlists SET id = \?, created_by = \?, created_date = \?, modified_by = \?, 
	modified_date = \?, deleted_by = \?, deleted_date = \?, is_deleted = \?, is_active = \?, user_id = \?, `

	if wl.TransId != "" {
		q = q + `trans_id = \?`
	}
	wl.ExpId = ""
	prep := mock.ExpectPrepare(q)
	prep.ExpectExec().WithArgs(wl.Id,
		wl.CreatedBy,
		wl.CreatedDate,
		wl.ModifiedBy,
		wl.ModifiedDate,
		wl.DeletedBy,
		wl.DeletedDate,
		wl.IsDeleted,
		wl.IsActive,
		wl.UserId,
		wl.TransId).WillReturnResult(sqlmock.NewResult(1, 1))

	i := WishlistRepo.NewWishListRepository(db)

	id, err := i.Insert(context.TODO(), &wl)
	assert.NoError(t, err)
	assert.Equal(t, id.Id, wl.Id)
}
func TestInsertByTransIdErrorExec(t *testing.T) {
	//user := "test"
	//now := time.Now()
	wl := mockWishlist[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	q := `INSERT wishlists SET id = \?, created_by = \?, created_date = \?, modified_by = \?, 
	modified_date = \?, deleted_by = \?, deleted_date = \?, is_deleted = \?, is_active = \?, user_id = \?, `

	if wl.TransId != "" {
		q = q + `trans_id = \?`
	}
	wl.ExpId = ""
	//if wl.ExpId != "" {
	//	q = q + `exp_id = \?`
	//
	//}
	prep := mock.ExpectPrepare(q)
	prep.ExpectExec().WithArgs(wl.Id,
		wl.CreatedBy,
		wl.CreatedDate,
		wl.ModifiedBy,
		wl.ModifiedDate,
		wl.DeletedBy,
		wl.DeletedDate,
		wl.IsDeleted,
		wl.IsActive,
		wl.UserId,
		wl.TransId,
		wl.TransId).WillReturnResult(sqlmock.NewResult(1, 1))

	i := WishlistRepo.NewWishListRepository(db)

	_, err = i.Insert(context.TODO(), &wl)
	assert.Error(t, err)
}
func TestInsertByExpId(t *testing.T) {
	//user := "test"
	//now := time.Now()
	wl := mockWishlist[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	q := `INSERT wishlists SET id = \?, created_by = \?, created_date = \?, modified_by = \?, 
	modified_date = \?, deleted_by = \?, deleted_date = \?, is_deleted = \?, is_active = \?, user_id = \?, `

	//if wl.TransId != "" {
	//	q = q + `trans_id = \?`
	//}
	wl.TransId = ""
	if wl.ExpId != "" {
		q = q + `exp_id = \?`

	}
	prep := mock.ExpectPrepare(q)
	prep.ExpectExec().WithArgs(wl.Id,
		wl.CreatedBy,
		wl.CreatedDate,
		wl.ModifiedBy,
		wl.ModifiedDate,
		wl.DeletedBy,
		wl.DeletedDate,
		wl.IsDeleted,
		wl.IsActive,
		wl.UserId,
		wl.ExpId).WillReturnResult(sqlmock.NewResult(1, 1))

	i := WishlistRepo.NewWishListRepository(db)

	id, err := i.Insert(context.TODO(), &wl)
	assert.NoError(t, err)
	assert.Equal(t, id.Id, wl.Id)
}
func TestInsertByExpIdErrorExec(t *testing.T) {
	//user := "test"
	//now := time.Now()
	wl := mockWishlist[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	q := `INSERT wishlists SET id = \?, created_by = \?, created_date = \?, modified_by = \?, 
	modified_date = \?, deleted_by = \?, deleted_date = \?, is_deleted = \?, is_active = \?, user_id = \?, `

	//if wl.TransId != "" {
	//	q = q + `trans_id = \?`
	//}
	wl.TransId = ""
	if wl.ExpId != "" {
		q = q + `exp_id = \?`

	}
	prep := mock.ExpectPrepare(q)
	prep.ExpectExec().WithArgs(wl.Id,
		wl.CreatedBy,
		wl.CreatedDate,
		wl.ModifiedBy,
		wl.ModifiedDate,
		wl.DeletedBy,
		wl.DeletedDate,
		wl.IsDeleted,
		wl.IsActive,
		wl.UserId,
		wl.ExpId,
		wl.ExpId).WillReturnResult(sqlmock.NewResult(1, 1))

	i := WishlistRepo.NewWishListRepository(db)

	_, err = i.Insert(context.TODO(), &wl)
	assert.Error(t, err)
}
