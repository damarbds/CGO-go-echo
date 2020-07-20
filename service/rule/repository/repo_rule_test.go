package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/models"
	RuleRepo "github.com/service/rule/repository"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	imagePath = "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Rule/8941695193938718058.jpg"
	mockRule  = []models.Rule{
		models.Rule{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			RuleName:     "test Rule",
			RuleIcon:     imagePath,
		},
		models.Rule{
			Id:           2,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			RuleName:     "test Rule 2",
			RuleIcon:     imagePath,
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
		AddRow(len(mockRule))

	query := `SELECT count\(\*\) AS count FROM rules WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := RuleRepo.NewRuleRepository(db)

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

	query := `SELECT count\(\*\) AS count FROM rules WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := RuleRepo.NewRuleRepository(db)

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
		"deleted_date", "is_deleted", "is_active", "rule_name", "rule_icon"}).
		AddRow(mockRule[0].Id, mockRule[0].CreatedBy, mockRule[0].CreatedDate, mockRule[0].ModifiedBy,
			mockRule[0].ModifiedDate, mockRule[0].DeletedBy, mockRule[0].DeletedDate, mockRule[0].IsDeleted,
			mockRule[0].IsActive, mockRule[0].RuleName, mockRule[0].RuleIcon).
		AddRow(mockRule[1].Id, mockRule[1].CreatedBy, mockRule[1].CreatedDate, mockRule[1].ModifiedBy,
			mockRule[1].ModifiedDate, mockRule[1].DeletedBy, mockRule[1].DeletedDate, mockRule[1].IsDeleted,
			mockRule[1].IsActive, mockRule[1].RuleName, mockRule[1].RuleIcon)

	query := `SELECT                                         \*\                                               FROM rules WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := RuleRepo.NewRuleRepository(db)

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
		"deleted_date", "is_deleted", "is_active", "rule_name", "rule_icon"}).
		AddRow(mockRule[0].Id, mockRule[0].CreatedBy, mockRule[0].CreatedDate, mockRule[0].ModifiedBy,
			mockRule[0].ModifiedDate, mockRule[0].DeletedBy, mockRule[0].DeletedDate, mockRule[0].IsDeleted,
			mockRule[0].IsActive, mockRule[0].RuleName,mockRule[0].RuleIcon).
		AddRow(mockRule[1].Id, mockRule[1].CreatedBy, mockRule[1].CreatedDate, mockRule[1].ModifiedBy,
			mockRule[1].ModifiedDate, mockRule[1].DeletedBy, mockRule[1].DeletedDate, mockRule[1].IsDeleted,
			mockRule[1].ModifiedBy, mockRule[1].RuleName, mockRule[1].RuleIcon)

	query := `SELECT                                         \*\                                               FROM rules WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := RuleRepo.NewRuleRepository(db)

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
		"deleted_date", "is_deleted", "is_active", "rule_name","rule_icon"}).
		AddRow(mockRule[0].Id, mockRule[0].CreatedBy, mockRule[0].CreatedDate, mockRule[0].ModifiedBy,
			mockRule[0].ModifiedDate, mockRule[0].DeletedBy, mockRule[0].DeletedDate, mockRule[0].IsDeleted,
			mockRule[0].IsActive, mockRule[0].RuleName, mockRule[0].RuleIcon).
		AddRow(mockRule[1].Id, mockRule[1].CreatedBy, mockRule[1].CreatedDate, mockRule[1].ModifiedBy,
			mockRule[1].ModifiedDate, mockRule[1].DeletedBy, mockRule[1].DeletedDate, mockRule[1].IsDeleted,
			mockRule[1].IsActive, mockRule[1].RuleName, mockRule[1].RuleIcon)

	query := `SELECT \*\ FROM rules where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := RuleRepo.NewRuleRepository(db)

	limit := 10
	offset := 0
	anArticle, err := a.Fetch(context.TODO(), limit, offset)
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
		"deleted_date", "is_deleted", "is_active", "rule_name",  "rule_icon"}).
		AddRow(mockRule[0].Id, mockRule[0].CreatedBy, mockRule[0].CreatedDate, mockRule[0].ModifiedBy,
			mockRule[0].ModifiedDate, mockRule[0].DeletedBy, mockRule[0].DeletedDate, mockRule[0].IsDeleted,
			mockRule[0].IsActive, mockRule[0].RuleName, mockRule[0].RuleIcon).
		AddRow(mockRule[1].Id, mockRule[1].CreatedBy, mockRule[1].CreatedDate, mockRule[1].ModifiedBy,
			mockRule[1].ModifiedDate, mockRule[1].DeletedBy, mockRule[1].DeletedDate, mockRule[1].IsDeleted,
			mockRule[1].IsActive, mockRule[1].RuleName, mockRule[1].RuleIcon)

	query := `SELECT \*\ FROM rules where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := RuleRepo.NewRuleRepository(db)

	//limit := 10
	//offset := 0
	anArticle, err := a.Fetch(context.TODO(), 0, 0)
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
		"deleted_date", "is_deleted", "is_active", "rule_name","rule_icon"}).
		AddRow(mockRule[0].Id, mockRule[0].CreatedBy, mockRule[0].CreatedDate, mockRule[0].ModifiedBy,
			mockRule[0].ModifiedDate, mockRule[0].DeletedBy, mockRule[0].DeletedDate, mockRule[0].IsDeleted,
			mockRule[0].IsActive, mockRule[0].RuleName,  mockRule[0].RuleIcon).
		AddRow(mockRule[1].Id, mockRule[1].CreatedBy, mockRule[1].CreatedDate, mockRule[1].ModifiedBy,
			mockRule[1].ModifiedDate, mockRule[1].DeletedBy, mockRule[1].DeletedDate, mockRule[1].IsDeleted,
			mockRule[1].ModifiedBy, mockRule[1].RuleName, mockRule[1].RuleIcon)

	query := `SELECT \*\ FROM rules where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := RuleRepo.NewRuleRepository(db)

	limit := 10
	offset := 0
	anArticle, err := a.Fetch(context.TODO(), limit, offset)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
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
		"deleted_date", "is_deleted", "is_active", "rule_name","rule_icon"}).
		AddRow(mockRule[0].Id, mockRule[0].CreatedBy, mockRule[0].CreatedDate, mockRule[0].ModifiedBy,
			mockRule[0].ModifiedDate, mockRule[0].DeletedBy, mockRule[0].DeletedDate, mockRule[0].IsDeleted,
			mockRule[0].IsActive, mockRule[0].RuleName, mockRule[0].RuleIcon).
		AddRow(mockRule[1].Id, mockRule[1].CreatedBy, mockRule[1].CreatedDate, mockRule[1].ModifiedBy,
			mockRule[1].ModifiedDate, mockRule[1].DeletedBy, mockRule[1].DeletedDate, mockRule[1].IsDeleted,
			mockRule[1].ModifiedBy, mockRule[1].RuleName, mockRule[1].RuleIcon)

	query := `SELECT \*\ FROM rules where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := RuleRepo.NewRuleRepository(db)

	//limit := 10
	//offset := 0
	anArticle, err := a.Fetch(context.TODO(), 0, 0)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
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
		"deleted_date", "is_deleted", "is_active", "rule_name", "rule_icon"}).
		AddRow(mockRule[0].Id, mockRule[0].CreatedBy, mockRule[0].CreatedDate, mockRule[0].ModifiedBy,
			mockRule[0].ModifiedDate, mockRule[0].DeletedBy, mockRule[0].DeletedDate, mockRule[0].IsDeleted,
			mockRule[0].IsActive, mockRule[0].RuleName,mockRule[0].RuleIcon)

	query := `SELECT \*\ FROM rules WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := RuleRepo.NewRuleRepository(db)

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
		"deleted_date", "is_deleted", "is_active", "rule_name",  "rule_icon"})

	query := `SELECT \*\ FROM Rule WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := RuleRepo.NewRuleRepository(db)

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
		"deleted_date", "is_deleted", "is_active", "rule_name", "rule_icon"}).
		AddRow(mockRule[0].Id, mockRule[0].CreatedBy, mockRule[0].CreatedDate, mockRule[0].ModifiedBy,
			mockRule[0].ModifiedDate, mockRule[0].DeletedBy, mockRule[0].DeletedDate, mockRule[0].IsDeleted,
			mockRule[0].ModifiedBy, mockRule[0].RuleName, mockRule[0].RuleIcon)

	query := `SELECT \*\ FROM rules WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := RuleRepo.NewRuleRepository(db)

	num := 1
	anArticle, err := a.GetById(context.TODO(), num)
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

	query := "UPDATE rules SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := RuleRepo.NewRuleRepository(db)

	err = a.Delete(context.TODO(), id, deletedBy)
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

	query := "UPDATE rules SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, id, id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := RuleRepo.NewRuleRepository(db)

	err = a.Delete(context.TODO(), id, deletedBy)
	assert.Error(t, err)
}
func TestInsert(t *testing.T) {
	user := "test"
	now := time.Now()
	a := mockRule[0]
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

	query := "INSERT rules SET created_by=\\? , created_date=\\? , modified_by=\\?, modified_date=\\? , 				deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? , rule_name=\\?,  rule_icon=\\? "
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.RuleName,
	 a.RuleIcon).WillReturnResult(sqlmock.NewResult(1, 1))

	i := RuleRepo.NewRuleRepository(db)

	id, err := i.Insert(context.TODO(), &a)
	assert.NoError(t, err)
	assert.Equal(t, *id, a.Id)
}
func TestInsertErrorExec(t *testing.T) {
	user := "test"
	now := time.Now()
	a := mockRule[0]
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

	query := "INSERT rules SET created_by=\\? , created_date=\\? , modified_by=\\?, modified_date=\\? , 				deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? , rule_name=\\?	,rule_icon=\\? "
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.RuleName,
		 a.RuleIcon, a.RuleIcon).WillReturnResult(sqlmock.NewResult(1, 1))

	i := RuleRepo.NewRuleRepository(db)

	_, err = i.Insert(context.TODO(), &a)
	assert.Error(t, err)
}
func TestUpdate(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	ar := mockRule[0]
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

	query := `UPDATE rules set modified_by=\?, modified_date=\? ,rule_name=\?,rule_icon=\? WHERE id = \?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.ModifiedBy, ar.ModifiedDate, ar.RuleName,  ar.RuleIcon, ar.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	a := RuleRepo.NewRuleRepository(db)

	err = a.Update(context.TODO(), &ar)
	assert.NoError(t, err)
	assert.Nil(t, err)
}
func TestUpdateErrorExec(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	ar := mockRule[0]
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

	query := `UPDATE rules set modified_by=\?, modified_date=\? ,rule_name=\?,rule_icon=\? WHERE id = \?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.ModifiedBy, ar.ModifiedDate, ar.RuleName,  ar.RuleIcon, ar.Id, ar.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	a := RuleRepo.NewRuleRepository(db)

	err = a.Update(context.TODO(), &ar)
	assert.Error(t, err)
}
