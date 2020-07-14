package repository_test

import (
	"context"
	"testing"

	"github.com/models"
	ExperienceFacilitiesRepo "github.com/service/exp_facilities/repository"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	imagePath = "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Facilities/8941695193938718058.jpg"
	transId = "dalsjdkalsdjlksdj"
	expId = "werewrwer"
	mockExperienceFacilities = []models.ExperienceFacilities{
		models.ExperienceFacilities{
			Id:           1,
			ExpId:        &expId,
			//TransId:	&transId,
			FacilitiesId: 1,
			Amount:1234,
		},
		models.ExperienceFacilities{
			Id:           2,
			TransId:&transId,
			FacilitiesId: 1,
			Amount:890789,
		},
	}
	mockExperienceFacilitiesJoin = []models.ExperienceFacilitiesJoin{
		models.ExperienceFacilitiesJoin{
			Id:           1,
			ExpId:        &expId,
			TransId:      nil,
			FacilitiesId: 2,
			Amount:       123123,
			FacilityName: "Test Facilities 1",
			IsNumerable:  1,
			FacilityIcon: &imagePath,
		},
		models.ExperienceFacilitiesJoin{
			Id:           1,
			ExpId:        nil,
			TransId:      &transId,
			FacilitiesId: 2,
			Amount:       123123,
			FacilityName: "Test Facilities 2",
			IsNumerable:  1,
			FacilityIcon: &imagePath,
		},
	}
)

func TestGetJoinByExpID(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "exp_id", "trans_id","facilities_id", "amount","facility_name","is_numerable", "facility_icon"}).
		AddRow(mockExperienceFacilitiesJoin[0].Id, mockExperienceFacilitiesJoin[0].ExpId, mockExperienceFacilitiesJoin[0].TransId,
			mockExperienceFacilitiesJoin[0].FacilitiesId,mockExperienceFacilitiesJoin[0].Amount,mockExperienceFacilitiesJoin[0].FacilityName,mockExperienceFacilitiesJoin[0].IsNumerable,
			mockExperienceFacilitiesJoin[0].FacilityIcon).
		AddRow(mockExperienceFacilitiesJoin[1].Id, mockExperienceFacilitiesJoin[1].ExpId, mockExperienceFacilitiesJoin[1].TransId,
			mockExperienceFacilitiesJoin[1].FacilitiesId,mockExperienceFacilitiesJoin[1].Amount,mockExperienceFacilitiesJoin[1].FacilityName,mockExperienceFacilitiesJoin[1].IsNumerable,
			mockExperienceFacilitiesJoin[1].FacilityIcon)

	query := `SELECT ef.* , f.facility_name,f.is_numerable,f.facility_icon
				FROM experience_facilities ef 
				JOIN facilities f ON ef.facilities_id = f.id`

	query = query + " WHERE ef.exp_id = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExperienceFacilitiesRepo.NewExpFacilitiesRepository(db)

	expId := mockExperienceFacilitiesJoin[0].ExpId
	anArticle, err := a.GetJoin(context.TODO(), *expId,"")
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestGetJoinByTransID(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "exp_id", "trans_id","facilities_id", "amount","facility_name","is_numerable", "facility_icon"}).
		AddRow(mockExperienceFacilitiesJoin[0].Id, mockExperienceFacilitiesJoin[0].ExpId, mockExperienceFacilitiesJoin[0].TransId,
			mockExperienceFacilitiesJoin[0].FacilitiesId,mockExperienceFacilitiesJoin[0].Amount,mockExperienceFacilitiesJoin[0].FacilityName,mockExperienceFacilitiesJoin[0].IsNumerable,
			mockExperienceFacilitiesJoin[0].FacilityIcon).
		AddRow(mockExperienceFacilitiesJoin[1].Id, mockExperienceFacilitiesJoin[1].ExpId, mockExperienceFacilitiesJoin[1].TransId,
			mockExperienceFacilitiesJoin[1].FacilitiesId,mockExperienceFacilitiesJoin[1].Amount,mockExperienceFacilitiesJoin[1].FacilityName,mockExperienceFacilitiesJoin[1].IsNumerable,
			mockExperienceFacilitiesJoin[1].FacilityIcon)

	query := `SELECT ef.* , f.facility_name,f.is_numerable,f.facility_icon
				FROM experience_facilities ef 
				JOIN facilities f ON ef.facilities_id = f.id`

	query = query + " WHERE ef.trans_id = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExperienceFacilitiesRepo.NewExpFacilitiesRepository(db)

	transId := mockExperienceFacilitiesJoin[1].TransId
	anArticle, err := a.GetJoin(context.TODO(), "",*transId)
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
	rows := sqlmock.NewRows([]string{"id", "exp_id", "trans_id","facilities_id", "amount","facility_name","is_numerable", "facility_icon"}).
		AddRow(mockExperienceFacilitiesJoin[0].Id, mockExperienceFacilitiesJoin[0].ExpId, mockExperienceFacilitiesJoin[0].TransId,
			mockExperienceFacilitiesJoin[0].FacilitiesId,mockExperienceFacilitiesJoin[0].Amount,mockExperienceFacilitiesJoin[0].FacilityName,mockExperienceFacilitiesJoin[0].IsNumerable,
			mockExperienceFacilitiesJoin[0].FacilityIcon).
		AddRow(mockExperienceFacilitiesJoin[1].Id, mockExperienceFacilitiesJoin[1].ExpId, mockExperienceFacilitiesJoin[1].TransId,
			mockExperienceFacilitiesJoin[1].FacilitiesId,mockExperienceFacilitiesJoin[1].Amount,mockExperienceFacilitiesJoin[1].FacilityName,mockExperienceFacilitiesJoin[1].IsNumerable,
			mockExperienceFacilitiesJoin[1].FacilityIcon)

	query := `SELECT ef.* , f.facility_name,f.is_numerable,f.facility_icon
				FROM experience_facilities ef 
				JOIN facilities f ON ef.facilities_id = f.id adasd`

	query = query + " WHERE ef.exp_id = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExperienceFacilitiesRepo.NewExpFacilitiesRepository(db)

	expId := mockExperienceFacilitiesJoin[0].ExpId
	_, err = a.GetJoin(context.TODO(), *expId,"")
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	//assert.Len(t, anArticle, 2)
}
func TestGetJoinByTransIDErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "exp_id", "trans_id","facilities_id", "amount","facility_name","is_numerable", "facility_icon"}).
		AddRow(mockExperienceFacilitiesJoin[0].Id, mockExperienceFacilitiesJoin[0].ExpId, mockExperienceFacilitiesJoin[0].TransId,
			mockExperienceFacilitiesJoin[0].FacilitiesId,mockExperienceFacilitiesJoin[0].Amount,mockExperienceFacilitiesJoin[0].FacilityName,mockExperienceFacilitiesJoin[0].IsNumerable,
			mockExperienceFacilitiesJoin[0].FacilityIcon).
		AddRow(mockExperienceFacilitiesJoin[1].Id, mockExperienceFacilitiesJoin[1].ExpId, mockExperienceFacilitiesJoin[1].TransId,
			mockExperienceFacilitiesJoin[1].FacilitiesId,mockExperienceFacilitiesJoin[1].Amount,mockExperienceFacilitiesJoin[1].FacilityName,mockExperienceFacilitiesJoin[1].IsNumerable,
			mockExperienceFacilitiesJoin[1].FacilityIcon)

	query := `SELECT ef.* , f.facility_name,f.is_numerable,f.facility_icon
				FROM experience_facilities ef 
				JOIN facilities f ON ef.facilities_id = f.id adasdas`

	query = query + " WHERE ef.trans_id = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExperienceFacilitiesRepo.NewExpFacilitiesRepository(db)

	transId := mockExperienceFacilitiesJoin[1].TransId
	_, err = a.GetJoin(context.TODO(), "",*transId)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	//assert.Len(t, anArticle, 2)
}
func TestInsert(t *testing.T) {
	a := mockExperienceFacilities[0]

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "INSERT experience_facilities SET exp_id=\\?,trans_id=\\?,facilities_id=\\?,amount=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.ExpId,a.TransId,a.FacilitiesId,a.Amount).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ExperienceFacilitiesRepo.NewExpFacilitiesRepository(db)

	err = i.Insert(context.TODO(), &a)
	assert.NoError(t, err)
	//assert.Equal(t, *id, a.Id)
}
func TestInsertErrorExec(t *testing.T) {

	a := mockExperienceFacilities[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "INSERT experience_facilities SET exp_id=\\?,trans_id=\\?,facilities_id=\\?,amount=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.ExpId,a.TransId,a.FacilitiesId,a.Amount,a.TransId).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ExperienceFacilitiesRepo.NewExpFacilitiesRepository(db)

	err = i.Insert(context.TODO(), &a)
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

	query := "DELETE FROM experience_facilities WHERE "
	query = query + "exp_id = ?"
	num := mockExperienceFacilities[0].ExpId

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(num).WillReturnResult(sqlmock.NewResult(12, 1))

	a := ExperienceFacilitiesRepo.NewExpFacilitiesRepository(db)

	err = a.Delete(context.TODO(), *num,"")
	assert.NoError(t, err)
}
func TestDeleteByExpIdErrorExecQueryString(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	query := "DELETE FROM experience_facilities WHERE "
	query = query + "exp_id = ?"
	num := mockExperienceFacilities[0].ExpId

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(num,num).WillReturnResult(sqlmock.NewResult(12, 1))

	a := ExperienceFacilitiesRepo.NewExpFacilitiesRepository(db)

	err = a.Delete(context.TODO(), *num,"")
	assert.Error(t, err)
}
func TestDeleteByTransId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "DELETE FROM experience_facilities WHERE "
	query = query + "trans_id = ?"
	num := mockExperienceFacilities[1].TransId

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(num).WillReturnResult(sqlmock.NewResult(12, 1))

	a := ExperienceFacilitiesRepo.NewExpFacilitiesRepository(db)

	err = a.Delete(context.TODO(), "",*num)
	assert.NoError(t, err)
}
func TestDeleteByTransIdErrorExecQueryString(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	query := "DELETE FROM experience_facilities WHERE "
	query = query + "trans_id = ?"
	num := mockExperienceFacilities[1].TransId

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(num,num).WillReturnResult(sqlmock.NewResult(12, 1))

	a := ExperienceFacilitiesRepo.NewExpFacilitiesRepository(db)

	err = a.Delete(context.TODO(), "",*num)
	assert.Error(t, err)
}
