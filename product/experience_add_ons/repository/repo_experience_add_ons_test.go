package repository_test

import (
	"context"
	"github.com/models"
	ExperienceAddOnRepo "github.com/product/experience_add_ons/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)
var(
	mockExperienceAddOn = []models.ExperienceAddOn{
		models.ExperienceAddOn{
			Id:           "asdasdasd",
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			Name:         "qeqwe",
			Desc:         "adasd",
			Currency:     1,
			Amount:       11212,
			ExpId:        "sfdsf",
		},
		models.ExperienceAddOn{
			Id:           "asdasdasd",
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			Name:         "qeqwe",
			Desc:         "adasd",
			Currency:     1,
			Amount:       11212,
			ExpId:        "sfdsf",
		},
	}
)
func TestGetByExpId(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","name","desc","currency","amount","exp_id"}).
		AddRow(mockExperienceAddOn[0].Id, mockExperienceAddOn[0].CreatedBy,mockExperienceAddOn[0].CreatedDate,mockExperienceAddOn[0].ModifiedBy,
			mockExperienceAddOn[0].ModifiedDate,mockExperienceAddOn[0].DeletedBy,mockExperienceAddOn[0].DeletedDate,mockExperienceAddOn[0].IsDeleted,
			mockExperienceAddOn[0].IsActive,mockExperienceAddOn[0].Name,mockExperienceAddOn[0].Desc,
		mockExperienceAddOn[0].Currency,mockExperienceAddOn[0].Amount,mockExperienceAddOn[0].ExpId).
		AddRow(mockExperienceAddOn[1].Id, mockExperienceAddOn[1].CreatedBy,mockExperienceAddOn[1].CreatedDate,mockExperienceAddOn[1].ModifiedBy,
			mockExperienceAddOn[1].ModifiedDate,mockExperienceAddOn[1].DeletedBy,mockExperienceAddOn[1].DeletedDate,mockExperienceAddOn[1].IsDeleted,
			mockExperienceAddOn[1].IsActive,mockExperienceAddOn[0].Name,mockExperienceAddOn[0].Desc,
		mockExperienceAddOn[0].Currency,mockExperienceAddOn[0].Amount,mockExperienceAddOn[0].ExpId)

	query := `select \*\ from experience_add_ons where exp_id =\? AND is_deleted = 0 AND is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExperienceAddOnRepo.NewexperienceRepository(db)

	anArticle, err := a.GetByExpId(context.TODO(),mockExperienceAddOn[0].ExpId)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestGetByExpIdError(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","name","desc","currency","amount","exp_id"}).
		AddRow(mockExperienceAddOn[0].Id, mockExperienceAddOn[0].CreatedBy,mockExperienceAddOn[0].CreatedDate,mockExperienceAddOn[0].ModifiedBy,
			mockExperienceAddOn[0].ModifiedDate,mockExperienceAddOn[0].DeletedBy,mockExperienceAddOn[0].DeletedDate,mockExperienceAddOn[0].IsDeleted,
			mockExperienceAddOn[0].IsActive,mockExperienceAddOn[0].Name,mockExperienceAddOn[0].Desc,
			mockExperienceAddOn[0].Currency,mockExperienceAddOn[0].Amount,mockExperienceAddOn[0].ExpId).
		AddRow(mockExperienceAddOn[1].Id, mockExperienceAddOn[1].CreatedBy,mockExperienceAddOn[1].CreatedDate,mockExperienceAddOn[1].ModifiedBy,
			mockExperienceAddOn[1].ModifiedDate,mockExperienceAddOn[1].DeletedBy,mockExperienceAddOn[1].DeletedDate,mockExperienceAddOn[1].IsDeleted,
			mockExperienceAddOn[1].ModifiedDate,mockExperienceAddOn[0].Name,mockExperienceAddOn[0].Desc,
			mockExperienceAddOn[0].Currency,mockExperienceAddOn[0].Amount,mockExperienceAddOn[0].ExpId)

	query := `select \*\ from experience_add_ons where exp_id =\? AND is_deleted = 0 AND is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExperienceAddOnRepo.NewexperienceRepository(db)

	anArticle, err := a.GetByExpId(context.TODO(),mockExperienceAddOn[0].ExpId)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
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
		"deleted_date","is_deleted","is_active","name","desc","currency","amount","exp_id"}).
		AddRow(mockExperienceAddOn[0].Id, mockExperienceAddOn[0].CreatedBy,mockExperienceAddOn[0].CreatedDate,mockExperienceAddOn[0].ModifiedBy,
			mockExperienceAddOn[0].ModifiedDate,mockExperienceAddOn[0].DeletedBy,mockExperienceAddOn[0].DeletedDate,mockExperienceAddOn[0].IsDeleted,
			mockExperienceAddOn[0].IsActive,mockExperienceAddOn[0].Name,mockExperienceAddOn[0].Desc,
			mockExperienceAddOn[0].Currency,mockExperienceAddOn[0].Amount,mockExperienceAddOn[0].ExpId).
		AddRow(mockExperienceAddOn[1].Id, mockExperienceAddOn[1].CreatedBy,mockExperienceAddOn[1].CreatedDate,mockExperienceAddOn[1].ModifiedBy,
			mockExperienceAddOn[1].ModifiedDate,mockExperienceAddOn[1].DeletedBy,mockExperienceAddOn[1].DeletedDate,mockExperienceAddOn[1].IsDeleted,
			mockExperienceAddOn[1].IsActive,mockExperienceAddOn[0].Name,mockExperienceAddOn[0].Desc,
			mockExperienceAddOn[0].Currency,mockExperienceAddOn[0].Amount,mockExperienceAddOn[0].ExpId)

	query := `select \*\ from experience_add_ons where id =\? `

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExperienceAddOnRepo.NewexperienceRepository(db)

	anArticle, err := a.GetById(context.TODO(), mockExperienceAddOn[0].Id)
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
		"deleted_date","is_deleted","is_active","name","desc","currency","amount","exp_id"}).
		AddRow(mockExperienceAddOn[0].Id, mockExperienceAddOn[0].CreatedBy,mockExperienceAddOn[0].CreatedDate,mockExperienceAddOn[0].ModifiedBy,
			mockExperienceAddOn[0].ModifiedDate,mockExperienceAddOn[0].DeletedBy,mockExperienceAddOn[0].DeletedDate,mockExperienceAddOn[0].IsDeleted,
			mockExperienceAddOn[0].IsActive,mockExperienceAddOn[0].Name,mockExperienceAddOn[0].Desc,
			mockExperienceAddOn[0].Currency,mockExperienceAddOn[0].Amount,mockExperienceAddOn[0].ExpId).
		AddRow(mockExperienceAddOn[1].Id, mockExperienceAddOn[1].CreatedBy,mockExperienceAddOn[1].CreatedDate,mockExperienceAddOn[1].ModifiedBy,
			mockExperienceAddOn[1].ModifiedDate,mockExperienceAddOn[1].DeletedBy,mockExperienceAddOn[1].DeletedDate,mockExperienceAddOn[1].IsDeleted,
			mockExperienceAddOn[1].ModifiedDate,mockExperienceAddOn[0].Name,mockExperienceAddOn[0].Desc,
			mockExperienceAddOn[0].Currency,mockExperienceAddOn[0].Amount,mockExperienceAddOn[0].ExpId)

	query := `select \*\ from experience_add_ons where id =\? `

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExperienceAddOnRepo.NewexperienceRepository(db)

	anArticle, err := a.GetById(context.TODO(), mockExperienceAddOn[0].Id)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t,anArticle)
	//assert.Len(t, anArticle, 2)
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
	query := `UPDATE  experience_add_ons SET deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? WHERE exp_id=\?`
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
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, mockExperienceAddOn[0].ExpId).WillReturnResult(sqlmock.NewResult(2, 1))

	a := ExperienceAddOnRepo.NewexperienceRepository(db)

	err = a.Deletes(context.TODO(), ids,mockExperienceAddOn[0].ExpId,deletedBy)
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
	query := `UPDATE  experience_add_ons SET deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? WHERE exp_id=\?`
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
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, mockExperienceAddOn[0].ExpId,ids).WillReturnResult(sqlmock.NewResult(2, 1))

	a := ExperienceAddOnRepo.NewexperienceRepository(db)

	err = a.Deletes(context.TODO(), ids,mockExperienceAddOn[0].ExpId,deletedBy)
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

	query := `UPDATE experience_add_ons SET deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? WHERE exp_id=\?`

	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, mockExperienceAddOn[0].ExpId).WillReturnResult(sqlmock.NewResult(2, 1))

	a := ExperienceAddOnRepo.NewexperienceRepository(db)

	err = a.DeleteByExpId(context.TODO(), mockExperienceAddOn[0].ExpId,deletedBy)
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


	query := `UPDATE experience_add_ons SET deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? WHERE exp_id=\?`

	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, mockExperienceAddOn[0].ExpId,deletedBy).WillReturnResult(sqlmock.NewResult(2, 1))

	a := ExperienceAddOnRepo.NewexperienceRepository(db)

	err = a.DeleteByExpId(context.TODO(), mockExperienceAddOn[0].ExpId,deletedBy)
	assert.Error(t, err)
}
func TestInsert(t *testing.T) {
	//user := "test"
	//now := time.Now()
	a := mockExperienceAddOn[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT experience_add_ons SET id=\? , created_by=\? , created_date=\? , modified_by=\?, modified_date=\? , 
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? ,experience_add_ons.name=\? , 
				experience_add_ons.desc = \? , currency=\? , amount=\?,exp_id=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.Name, a.Desc, a.Currency,
		a.Amount, a.ExpId).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ExperienceAddOnRepo.NewexperienceRepository(db)

	id, err := i.Insert(context.TODO(), a)
	assert.NoError(t, err)
	assert.Equal(t, id, a.Id)
}
func TestInsertErrorExec(t *testing.T) {
	//user := "test"
	//now := time.Now()
	a := mockExperienceAddOn[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT experience_add_ons SET id=\? , created_by=\? , created_date=\? , modified_by=\?, modified_date=\? , 
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? ,experience_add_ons.name=\? , 
				experience_add_ons.desc = \? , currency=\? , amount=\?,exp_id=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.Name, a.Desc, a.Currency,
		a.Amount, a.ExpId,a.ExpId).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ExperienceAddOnRepo.NewexperienceRepository(db)

	_, err = i.Insert(context.TODO(), a)
	assert.Error(t, err)
}
func TestUpdate(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	a := mockExperienceAddOn[0]
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

	query := `UPDATE experience_add_ons SET modified_by=\?, modified_date=\? , 
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? ,experience_add_ons.name=\? , 
				experience_add_ons.desc = \?,currency=\?,amount=\?,exp_id=\?
				WHERE id=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.ModifiedBy, time.Now(), nil, nil, 0, 1, a.Name, a.Desc, a.Currency,
		a.Amount, a.ExpId, a.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	u := ExperienceAddOnRepo.NewexperienceRepository(db)

	err = u.Update(context.TODO(), a)
	assert.NoError(t, err)
	assert.Nil(t,err)
}
func TestUpdateErrorExec(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	a := mockExperienceAddOn[0]
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

	query := `UPDATE experience_add_ons SET modified_by=\?, modified_date=\? , 
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? ,experience_add_ons.name=\? , 
				experience_add_ons.desc = \?,currency=\?,amount=\?,exp_id=\?
				WHERE id=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.ModifiedBy, time.Now(), nil, nil, 0, 1, a.Name, a.Desc, a.Currency,
		a.Amount, a.ExpId, a.Id,a.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	u := ExperienceAddOnRepo.NewexperienceRepository(db)

	err = u.Update(context.TODO(), a)
	assert.Error(t, err)
}
