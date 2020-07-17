package repository_test

import (
	"context"
	"testing"

	TempUserPreferencesRepo "github.com/service/temp_user_preferences/repository"

	"github.com/models"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	harborsId                  = "asdasdasdasd"
	cityId 	= 1
	province = 1
	mockTempUserPreference = []models.TempUserPreference{
		models.TempUserPreference{
			Id:         1,
			ProvinceId: &province,
			CityId:     &cityId,
			HarborsId:  &harborsId,
		},
		models.TempUserPreference{
			Id:           2,
			ProvinceId: &province,
			CityId:     &cityId,
			HarborsId:  &harborsId,
		},
	}
)

func TestTimeOptionsWithPagination(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "province_id", "city_id", "harbors_id"}).
		AddRow(mockTempUserPreference[0].Id, mockTempUserPreference[0].ProvinceId, mockTempUserPreference[0].CityId,
			mockTempUserPreference[0].HarborsId).
		AddRow(mockTempUserPreference[1].Id, mockTempUserPreference[1].ProvinceId, mockTempUserPreference[1].CityId,
			mockTempUserPreference[1].HarborsId)

	query := `Select \*\ FROM temp_user_preferences `

	query = query + ` LIMIT \? OFFSET \? `

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := TempUserPreferencesRepo.NewtempUserPreferencesRepository(db)
	page := 0
	size := 1
	anArticle, err := a.GetAll(context.TODO(),&page,&size)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestTimeOptionsWithPaginationErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "province_id", "city_id", "harbors_id"}).
		AddRow(mockTempUserPreference[0].Id, mockTempUserPreference[0].ProvinceId, mockTempUserPreference[0].CityId,
			mockTempUserPreference[0].HarborsId).
		AddRow(mockTempUserPreference[1].Id, mockTempUserPreference[1].ProvinceId, mockTempUserPreference[1].CityId,
			mockTempUserPreference[1].ProvinceId)

	query := `Select \*\ FROM temp_user_preferences `

	query = query + ` LIMIT \? OFFSET \? `

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := TempUserPreferencesRepo.NewtempUserPreferencesRepository(db)
	page := 0
	size := 1
	_, err = a.GetAll(context.TODO(),&page,&size)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	//assert.Len(t, anArticle, 2)
}
func TestTimeOptionsWithoutPagination(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "province_id", "city_id", "harbors_id"}).
		AddRow(mockTempUserPreference[0].Id, mockTempUserPreference[0].ProvinceId, mockTempUserPreference[0].CityId,
			mockTempUserPreference[0].HarborsId).
		AddRow(mockTempUserPreference[1].Id, mockTempUserPreference[1].ProvinceId, mockTempUserPreference[1].CityId,
			mockTempUserPreference[1].HarborsId)

	query := `Select \*\ FROM temp_user_preferences `

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := TempUserPreferencesRepo.NewtempUserPreferencesRepository(db)

	anArticle, err := a.GetAll(context.TODO(),nil,nil)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestTimeOptionsWithoutPaginationErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "province_id", "city_id", "harbors_id"}).
		AddRow(mockTempUserPreference[0].Id, mockTempUserPreference[0].ProvinceId, mockTempUserPreference[0].CityId,
			mockTempUserPreference[0].HarborsId).
		AddRow(mockTempUserPreference[1].Id, mockTempUserPreference[1].ProvinceId, mockTempUserPreference[1].CityId,
			mockTempUserPreference[1].ProvinceId)

	query := `Select \*\ FROM temp_user_preferences `

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := TempUserPreferencesRepo.NewtempUserPreferencesRepository(db)

	_, err = a.GetAll(context.TODO(),nil,nil)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	//assert.Len(t, anArticle, 2)
}
