package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/models"
	ExperiencePaymentRepo "github.com/service/exp_payment/repository"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	customPrice           = `[{"currency":"USD","price":300,"guest":4},{"currency":"USD","price":350,"guest":6}]`
	mockExperiencePayment = []models.ExperiencePayment{
		models.ExperiencePayment{
			Id:               "dasdsadasdasd",
			CreatedBy:        "test",
			CreatedDate:      time.Now(),
			ModifiedBy:       nil,
			ModifiedDate:     nil,
			DeletedBy:        nil,
			DeletedDate:      nil,
			IsDeleted:        0,
			IsActive:         1,
			ExpPaymentTypeId: "qewewqeqw",
			ExpId:            "zxczxczxczx",
			PriceItemType:    1,
			Currency:         1,
			Price:            123213,
			CustomPrice:      &customPrice,
		},
		models.ExperiencePayment{
			Id:               "kljlkjljl",
			CreatedBy:        "test",
			CreatedDate:      time.Now(),
			ModifiedBy:       nil,
			ModifiedDate:     nil,
			DeletedBy:        nil,
			DeletedDate:      nil,
			IsDeleted:        0,
			IsActive:         1,
			ExpPaymentTypeId: "qewewqeqw",
			ExpId:            "zxczxczxczx",
			PriceItemType:    1,
			Currency:         1,
			Price:            123213,
			CustomPrice:      &customPrice,
		},
	}
	mockExperiencePaymentJoinType = []models.ExperiencePaymentJoinType{
		models.ExperiencePaymentJoinType{
			Id:                 "kljlkjljl",
			CreatedBy:          "test",
			CreatedDate:        time.Now(),
			ModifiedBy:         nil,
			ModifiedDate:       nil,
			DeletedBy:          nil,
			DeletedDate:        nil,
			IsDeleted:          0,
			IsActive:           1,
			ExpPaymentTypeId:   "qewewqeqw",
			ExpId:              "zxczxczxczx",
			PriceItemType:      1,
			Currency:           1,
			Price:              123213,
			CustomPrice:        &customPrice,
			ExpPaymentTypeName: "Full Payment",
			ExpPaymentTypeDesc: "100%",
		},
		models.ExperiencePaymentJoinType{
			Id:                 "kljlkjljl",
			CreatedBy:          "test",
			CreatedDate:        time.Now(),
			ModifiedBy:         nil,
			ModifiedDate:       nil,
			DeletedBy:          nil,
			DeletedDate:        nil,
			IsDeleted:          0,
			IsActive:           1,
			ExpPaymentTypeId:   "qewewqeqw",
			ExpId:              "zxczxczxczx",
			PriceItemType:      1,
			Currency:           1,
			Price:              123213,
			CustomPrice:        &customPrice,
			ExpPaymentTypeName: "Down Payment",
			ExpPaymentTypeDesc: "30%",
		},
	}
)

func TestGetById(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "exp_payment_type_id", "exp_id","price_item_type","currency","price",
		"custom_price","exp_payment_type_name","exp_payment_type_desc"}).
		AddRow(mockExperiencePaymentJoinType[0].Id, mockExperiencePaymentJoinType[0].CreatedBy, mockExperiencePaymentJoinType[0].CreatedDate, mockExperiencePaymentJoinType[0].ModifiedBy,
			mockExperiencePaymentJoinType[0].ModifiedDate, mockExperiencePaymentJoinType[0].DeletedBy, mockExperiencePaymentJoinType[0].DeletedDate, mockExperiencePaymentJoinType[0].IsDeleted,
			mockExperiencePaymentJoinType[0].IsActive, mockExperiencePaymentJoinType[0].ExpPaymentTypeId, mockExperiencePaymentJoinType[0].ExpId,
		mockExperiencePaymentJoinType[0].PriceItemType,mockExperiencePaymentJoinType[0].Currency,mockExperiencePaymentJoinType[0].Price,
		mockExperiencePaymentJoinType[0].CustomPrice,mockExperiencePaymentJoinType[0].ExpPaymentTypeName,mockExperiencePaymentJoinType[0].ExpPaymentTypeDesc).
		AddRow(mockExperiencePaymentJoinType[1].Id, mockExperiencePaymentJoinType[1].CreatedBy, mockExperiencePaymentJoinType[1].CreatedDate, mockExperiencePaymentJoinType[1].ModifiedBy,
			mockExperiencePaymentJoinType[1].ModifiedDate, mockExperiencePaymentJoinType[1].DeletedBy, mockExperiencePaymentJoinType[1].DeletedDate, mockExperiencePaymentJoinType[1].IsDeleted,
			mockExperiencePaymentJoinType[1].IsActive, mockExperiencePaymentJoinType[1].ExpPaymentTypeId, mockExperiencePaymentJoinType[1].ExpId,
			mockExperiencePaymentJoinType[1].PriceItemType,mockExperiencePaymentJoinType[1].Currency,mockExperiencePaymentJoinType[1].Price,
			mockExperiencePaymentJoinType[1].CustomPrice,mockExperiencePaymentJoinType[1].ExpPaymentTypeName,mockExperiencePaymentJoinType[1].ExpPaymentTypeDesc)

	query := `SELECT ep.\*\,ept.exp_payment_type_name as payment_type_name ,ept.exp_payment_type_desc as payment_type_desc 
			FROM experience_payments ep
			JOIN experience_payment_types ept on ept.id = ep.exp_payment_type_id
			WHERE ep.id = \? `

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExperiencePaymentRepo.NewExpPaymentRepository(db)

	id := mockExperiencePaymentJoinType[1].Id
	anArticle, err := a.GetById(context.TODO(), id)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestGetByIdErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "exp_payment_type_id", "exp_id","price_item_type","currency","price",
		"custom_price","exp_payment_type_name","exp_payment_type_desc"}).
		AddRow(mockExperiencePaymentJoinType[0].Id, mockExperiencePaymentJoinType[0].CreatedBy, mockExperiencePaymentJoinType[0].CreatedDate, mockExperiencePaymentJoinType[0].ModifiedBy,
			mockExperiencePaymentJoinType[0].ModifiedDate, mockExperiencePaymentJoinType[0].DeletedBy, mockExperiencePaymentJoinType[0].DeletedDate, mockExperiencePaymentJoinType[0].IsDeleted,
			mockExperiencePaymentJoinType[0].IsActive, mockExperiencePaymentJoinType[0].ExpPaymentTypeId, mockExperiencePaymentJoinType[0].ExpId,
			mockExperiencePaymentJoinType[0].PriceItemType,mockExperiencePaymentJoinType[0].Currency,mockExperiencePaymentJoinType[0].Price,
			mockExperiencePaymentJoinType[0].CustomPrice,mockExperiencePaymentJoinType[0].ExpPaymentTypeName,mockExperiencePaymentJoinType[0].ExpPaymentTypeDesc).
		AddRow(mockExperiencePaymentJoinType[1].Id, mockExperiencePaymentJoinType[1].CreatedBy, mockExperiencePaymentJoinType[1].CreatedDate, mockExperiencePaymentJoinType[1].ModifiedBy,
			mockExperiencePaymentJoinType[1].ModifiedDate, mockExperiencePaymentJoinType[1].DeletedBy, mockExperiencePaymentJoinType[1].DeletedDate, mockExperiencePaymentJoinType[1].IsDeleted,
			mockExperiencePaymentJoinType[1].ModifiedDate, mockExperiencePaymentJoinType[1].ExpPaymentTypeId, mockExperiencePaymentJoinType[1].ExpId,
			mockExperiencePaymentJoinType[1].PriceItemType,mockExperiencePaymentJoinType[1].Currency,mockExperiencePaymentJoinType[1].Price,
			mockExperiencePaymentJoinType[1].CustomPrice,mockExperiencePaymentJoinType[1].ExpPaymentTypeName,mockExperiencePaymentJoinType[1].ExpPaymentTypeDesc)

	query := `SELECT ep.\*\,ept.exp_payment_type_name as payment_type_name ,ept.exp_payment_type_desc as payment_type_desc 
			FROM experience_payments ep
			JOIN experience_payment_types ept on ept.id = ep.exp_payment_type_id
			WHERE ep.id = \? `

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExperiencePaymentRepo.NewExpPaymentRepository(db)

	id := mockExperiencePaymentJoinType[1].Id
	_, err = a.GetById(context.TODO(), id)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	//assert.Len(t, anArticle, 2)
}
func TestGetJoinByExpID(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "exp_payment_type_id", "exp_id","price_item_type","currency","price",
		"custom_price","exp_payment_type_name","exp_payment_type_desc"}).
		AddRow(mockExperiencePaymentJoinType[0].Id, mockExperiencePaymentJoinType[0].CreatedBy, mockExperiencePaymentJoinType[0].CreatedDate, mockExperiencePaymentJoinType[0].ModifiedBy,
			mockExperiencePaymentJoinType[0].ModifiedDate, mockExperiencePaymentJoinType[0].DeletedBy, mockExperiencePaymentJoinType[0].DeletedDate, mockExperiencePaymentJoinType[0].IsDeleted,
			mockExperiencePaymentJoinType[0].IsActive, mockExperiencePaymentJoinType[0].ExpPaymentTypeId, mockExperiencePaymentJoinType[0].ExpId,
			mockExperiencePaymentJoinType[0].PriceItemType,mockExperiencePaymentJoinType[0].Currency,mockExperiencePaymentJoinType[0].Price,
			mockExperiencePaymentJoinType[0].CustomPrice,mockExperiencePaymentJoinType[0].ExpPaymentTypeName,mockExperiencePaymentJoinType[0].ExpPaymentTypeDesc).
		AddRow(mockExperiencePaymentJoinType[1].Id, mockExperiencePaymentJoinType[1].CreatedBy, mockExperiencePaymentJoinType[1].CreatedDate, mockExperiencePaymentJoinType[1].ModifiedBy,
			mockExperiencePaymentJoinType[1].ModifiedDate, mockExperiencePaymentJoinType[1].DeletedBy, mockExperiencePaymentJoinType[1].DeletedDate, mockExperiencePaymentJoinType[1].IsDeleted,
			mockExperiencePaymentJoinType[1].IsActive, mockExperiencePaymentJoinType[1].ExpPaymentTypeId, mockExperiencePaymentJoinType[1].ExpId,
			mockExperiencePaymentJoinType[1].PriceItemType,mockExperiencePaymentJoinType[1].Currency,mockExperiencePaymentJoinType[1].Price,
			mockExperiencePaymentJoinType[1].CustomPrice,mockExperiencePaymentJoinType[1].ExpPaymentTypeName,mockExperiencePaymentJoinType[1].ExpPaymentTypeDesc)

	query := `SELECT ep.\*\,ept.exp_payment_type_name as payment_type_name ,ept.exp_payment_type_desc as payment_type_desc 
			FROM experience_payments ep
			JOIN experience_payment_types ept on ept.id = ep.exp_payment_type_id
			WHERE ep.exp_id = \? AND ep.is_deleted = 0 AND ep.is_active = 1
			ORDER BY exp_payment_type_id DESC`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExperiencePaymentRepo.NewExpPaymentRepository(db)

	expId := mockExperiencePaymentJoinType[0].ExpId
	anArticle, err := a.GetByExpID(context.TODO(), expId)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestGetJoinByExpIDErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "exp_payment_type_id", "exp_id","price_item_type","currency","price",
		"custom_price","exp_payment_type_name","exp_payment_type_desc"}).
		AddRow(mockExperiencePaymentJoinType[0].Id, mockExperiencePaymentJoinType[0].CreatedBy, mockExperiencePaymentJoinType[0].CreatedDate, mockExperiencePaymentJoinType[0].ModifiedBy,
			mockExperiencePaymentJoinType[0].ModifiedDate, mockExperiencePaymentJoinType[0].DeletedBy, mockExperiencePaymentJoinType[0].DeletedDate, mockExperiencePaymentJoinType[0].IsDeleted,
			mockExperiencePaymentJoinType[0].IsActive, mockExperiencePaymentJoinType[0].ExpPaymentTypeId, mockExperiencePaymentJoinType[0].ExpId,
			mockExperiencePaymentJoinType[0].PriceItemType,mockExperiencePaymentJoinType[0].Currency,mockExperiencePaymentJoinType[0].Price,
			mockExperiencePaymentJoinType[0].CustomPrice,mockExperiencePaymentJoinType[0].ExpPaymentTypeName,mockExperiencePaymentJoinType[0].ExpPaymentTypeDesc).
		AddRow(mockExperiencePaymentJoinType[1].Id, mockExperiencePaymentJoinType[1].CreatedBy, mockExperiencePaymentJoinType[1].CreatedDate, mockExperiencePaymentJoinType[1].ModifiedBy,
			mockExperiencePaymentJoinType[1].ModifiedDate, mockExperiencePaymentJoinType[1].DeletedBy, mockExperiencePaymentJoinType[1].DeletedDate, mockExperiencePaymentJoinType[1].IsDeleted,
			mockExperiencePaymentJoinType[1].ModifiedDate, mockExperiencePaymentJoinType[1].ExpPaymentTypeId, mockExperiencePaymentJoinType[1].ExpId,
			mockExperiencePaymentJoinType[1].PriceItemType,mockExperiencePaymentJoinType[1].Currency,mockExperiencePaymentJoinType[1].Price,
			mockExperiencePaymentJoinType[1].CustomPrice,mockExperiencePaymentJoinType[1].ExpPaymentTypeName,mockExperiencePaymentJoinType[1].ExpPaymentTypeDesc)

	query := `SELECT ep.\*\,ept.exp_payment_type_name as payment_type_name ,ept.exp_payment_type_desc as payment_type_desc 
			FROM experience_payments ep
			JOIN experience_payment_types ept on ept.id = ep.exp_payment_type_id
			WHERE ep.exp_id = \? AND ep.is_deleted = 0 AND ep.is_active = 1
			ORDER BY exp_payment_type_id DESC`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExperiencePaymentRepo.NewExpPaymentRepository(db)

	expId := mockExperiencePaymentJoinType[0].ExpId
	_, err = a.GetByExpID(context.TODO(), expId)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	//assert.Len(t, anArticle, 2)
}
func TestInsert(t *testing.T) {
	a := mockExperiencePayment[0]

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT experience_payments SET id=\? , created_by=\? , created_date=\? , modified_by=\?, modified_date=\? , 
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , exp_payment_type_id=\?,exp_id=\?,
				price_item_type=\?,currency=\?,price=\?,custom_price=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.Id, a.CreatedBy, a.CreatedBy, nil, nil, nil, nil, 0, 1, a.ExpPaymentTypeId, a.ExpId,
		a.PriceItemType, a.Currency, a.Price, a.CustomPrice).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ExperiencePaymentRepo.NewExpPaymentRepository(db)

	_,err = i.Insert(context.TODO(), a)
	assert.NoError(t, err)
	//assert.Equal(t, *id, a.Id)
}
func TestInsertErrorExec(t *testing.T) {

	a := mockExperiencePayment[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT experience_payments SET id=\? , created_by=\? , created_date=\? , modified_by=\?, modified_date=\? , 
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , exp_payment_type_id=\?,exp_id=\?,
				price_item_type=\?,currency=\?,price=\?,custom_price=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.Id, a.CreatedBy, a.CreatedBy, nil, nil, nil, nil, 0, 1, a.ExpPaymentTypeId, a.ExpId,
		a.PriceItemType, a.Currency, a.Price, a.CustomPrice,a.CustomPrice).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ExperiencePaymentRepo.NewExpPaymentRepository(db)

	_,err = i.Insert(context.TODO(), a)
	assert.Error(t, err)
}
func TestUpdate(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	a := mockExperiencePayment[0]
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

	query := `UPDATE experience_payments SET modified_by=\?, modified_date=\? , 
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , exp_payment_type_id=\?,exp_id=\?,
				price_item_type=\?,currency=\?,price=\?,custom_price=\? 
				WHERE id =\?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.ModifiedBy, a.ModifiedDate, nil, nil, 0, 1, a.ExpPaymentTypeId, a.ExpId,
		a.PriceItemType, a.Currency, a.Price, a.CustomPrice, a.Id).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ExperiencePaymentRepo.NewExpPaymentRepository(db)

	err = i.Update(context.TODO(), a)
	assert.NoError(t, err)
	//assert.Equal(t, *id, a.Id)
}
func TestUpdateErrorExec(t *testing.T) {

	now := time.Now()
	modifyBy := "test"
	a := mockExperiencePayment[0]
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

	query := `UPDATE experience_payments SET modified_by=\?, modified_date=\? , 
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , exp_payment_type_id=\?,exp_id=\?,
				price_item_type=\?,currency=\?,price=\?,custom_price=\? 
				WHERE id =\?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.ModifiedBy, a.ModifiedDate, nil, nil, 0, 1, a.ExpPaymentTypeId, a.ExpId,
		a.PriceItemType, a.Currency, a.Price, a.CustomPrice, a.Id,a.Id).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ExperiencePaymentRepo.NewExpPaymentRepository(db)

	err = i.Update(context.TODO(), a)
	assert.Error(t, err)
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
	ids := []string{"adsasfjalkfja","adalksdjalkfjla","ijaodijasdoiajsd","ijaodijasdoiajsd","ijaodijasdoiajsd"}
	query := `UPDATE  experience_payments SET deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? WHERE exp_id=\?`
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
	deletedBy := mockExperiencePayment[0].CreatedBy
	expId := mockExperiencePayment[0].ExpId
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, expId).WillReturnResult(sqlmock.NewResult(1, 1))

	a := ExperiencePaymentRepo.NewExpPaymentRepository(db)

	err = a.Deletes(context.TODO(), ids, expId,deletedBy)
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
	ids := []string{"adsasfjalkfja","adalksdjalkfjla","ijaodijasdoiajsd","ijaodijasdoiajsd","ijaodijasdoiajsd"}
	query := `UPDATE  experience_payments SET deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? WHERE exp_id=\?`
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
	deletedBy := mockExperiencePayment[0].CreatedBy
	expId := mockExperiencePayment[0].ExpId
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, expId,expId).WillReturnResult(sqlmock.NewResult(1, 1))

	a := ExperiencePaymentRepo.NewExpPaymentRepository(db)

	err = a.Deletes(context.TODO(), ids, expId,deletedBy)
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

	query := `UPDATE experience_payments SET deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? WHERE exp_id=\?`

	expId := mockExperiencePayment[0].ExpId
	deletedBy := mockExperiencePayment[0].CreatedBy
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, expId).WillReturnResult(sqlmock.NewResult(1, 1))

	a := ExperiencePaymentRepo.NewExpPaymentRepository(db)

	err = a.DeleteByExpId(context.TODO(), expId, deletedBy)
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
	query := `UPDATE experience_payments SET deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? WHERE exp_id=\?`

	expId := mockExperiencePayment[0].ExpId
	deletedBy := mockExperiencePayment[0].CreatedBy
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, expId,expId).WillReturnResult(sqlmock.NewResult(1, 1))

	a := ExperiencePaymentRepo.NewExpPaymentRepository(db)

	err = a.DeleteByExpId(context.TODO(), expId, deletedBy)
	assert.Error(t, err)
}
