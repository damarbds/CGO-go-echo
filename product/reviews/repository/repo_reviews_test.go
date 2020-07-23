package repository_test

import (
	"context"
	"github.com/models"
	ReviewRepo "github.com/product/reviews/repository"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)
var(
	rating = 12
	userId = "adasdasdqweqwe"
	reviewValue float64 = 12
	mockReview = []models.Review{
		models.Review{
			Id:                "asdasdasdsadqeq",
			CreatedBy:         "Test 1",
			CreatedDate:       time.Now(),
			ModifiedBy:        nil,
			ModifiedDate:      nil,
			DeletedBy:         nil,
			DeletedDate:       nil,
			IsDeleted:         0,
			IsActive:          1,
			Values:            12,
			Desc:              "Buruk",
			ExpId:             "qewqweqwe",
			UserId:            &userId,
			GuideReview:       &reviewValue,
			ActivitiesReview:  &reviewValue,
			ServiceReview:     &reviewValue,
			CleanlinessReview: &reviewValue,
			ValueReview:       &reviewValue,
		},
		models.Review{
			Id:                "jlkjlkjlkjkl",
			CreatedBy:         "Test 1",
			CreatedDate:       time.Now(),
			ModifiedBy:        nil,
			ModifiedDate:      nil,
			DeletedBy:         nil,
			DeletedDate:       nil,
			IsDeleted:         0,
			IsActive:          1,
			Values:            12,
			Desc:              "Buruk",
			ExpId:             "qewqweqwe",
			UserId:            &userId,
			GuideReview:       &reviewValue,
			ActivitiesReview:  &reviewValue,
			ServiceReview:     &reviewValue,
			CleanlinessReview: &reviewValue,
			ValueReview:       &reviewValue,
		},
	}
)
func TestCountRating(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(len(mockReview))

	query := `SELECT COUNT\(\*\) as count FROM reviews r WHERE exp_id = \? AND is_deleted = 0 AND is_active = 1`
	if rating != 0 {
		query = query + ` AND r.values = ` + strconv.Itoa(rating)
	}
	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ReviewRepo.NewReviewRepository(db)

	res, err := a.CountRating(context.TODO(),rating,mockReview[0].ExpId)
	assert.NoError(t, err)
	assert.Equal(t, res, 2, "")
}
func TestCountRatingFetch(t *testing.T) {
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

	query := `SELECT COUNT\(\*\) as count FROM reviews r WHERE exp_id = \? AND is_deleted = 0 AND is_active = 1`
	if rating != 0 {
		query = query + ` AND r.values = ` + strconv.Itoa(rating)
	}
	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ReviewRepo.NewReviewRepository(db)

	_, err = a.CountRating(context.TODO(),rating,mockReview[0].ExpId)
	assert.Error(t, err)
}
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
		"deleted_date","is_deleted","is_active","values","desc","exp_id","user_id","guide_review","activities_review",
		"service_review","cleanlinessreview","value_review"}).
		AddRow(mockReview[0].Id, mockReview[0].CreatedBy,mockReview[0].CreatedDate,mockReview[0].ModifiedBy,
			mockReview[0].ModifiedDate,mockReview[0].DeletedBy,mockReview[0].DeletedDate,mockReview[0].IsDeleted,
			mockReview[0].IsActive,mockReview[0].Values,mockReview[0].Desc,mockReview[0].ExpId,mockReview[0].UserId,
			mockReview[0].GuideReview,mockReview[0].ActivitiesReview,mockReview[0].ServiceReview,mockReview[0].CleanlinessReview,
			mockReview[0].ValueReview).
		AddRow(mockReview[1].Id, mockReview[1].CreatedBy,mockReview[1].CreatedDate,mockReview[1].ModifiedBy,
			mockReview[1].ModifiedDate,mockReview[1].DeletedBy,mockReview[1].DeletedDate,mockReview[1].IsDeleted,
			mockReview[1].IsActive,mockReview[1].Values,mockReview[1].Desc,mockReview[1].ExpId,mockReview[1].UserId,
			mockReview[1].GuideReview,mockReview[1].ActivitiesReview,mockReview[1].ServiceReview,mockReview[1].CleanlinessReview,
			mockReview[1].ValueReview)

	query := `select \*\ from reviews r where exp_id = \? AND is_deleted = 0 AND is_active = 1`
	if rating != 0 {
		query = query + ` AND r.values = ` + strconv.Itoa(rating)
	}
	if userId != "" {
		query = query + ` AND r.user_id = '` + userId + `' `
	}
	sortBy := "oldestdate"
	if sortBy != "" {
		if sortBy == "ratingup" {
			query = query + ` ORDER BY r.values DESC`
		} else if sortBy == "ratingdown" {
			query = query + ` ORDER BY r.values ASC`
		} else if sortBy == "latestdate" {
			query = query + ` ORDER BY created_date DESC`
		} else if sortBy == "oldestdate" {
			query = query + ` ORDER BY created_date ASC`
		}
	}
	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ReviewRepo.NewReviewRepository(db)

	anArticle, err := a.GetByExpId(context.TODO(),mockReview[0].ExpId,sortBy,rating,0,0,userId)
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
		"deleted_date","is_deleted","is_active","values","desc","exp_id","user_id","guide_review","activities_review",
		"service_review","cleanlinessreview","value_review"}).
		AddRow(mockReview[0].Id, mockReview[0].CreatedBy,mockReview[0].CreatedDate,mockReview[0].ModifiedBy,
			mockReview[0].ModifiedDate,mockReview[0].DeletedBy,mockReview[0].DeletedDate,mockReview[0].IsDeleted,
			mockReview[0].IsActive,mockReview[0].Values,mockReview[0].Desc,mockReview[0].ExpId,mockReview[0].UserId,
			mockReview[0].GuideReview,mockReview[0].ActivitiesReview,mockReview[0].ServiceReview,mockReview[0].CleanlinessReview,
			mockReview[0].ValueReview).
		AddRow(mockReview[1].Id, mockReview[1].CreatedBy,mockReview[1].CreatedDate,mockReview[1].ModifiedBy,
			mockReview[1].ModifiedDate,mockReview[1].DeletedBy,mockReview[1].DeletedDate,mockReview[1].IsDeleted,
			mockReview[1].ModifiedDate,mockReview[1].Values,mockReview[1].Desc,mockReview[1].ExpId,mockReview[1].UserId,
			mockReview[1].GuideReview,mockReview[1].ActivitiesReview,mockReview[1].ServiceReview,mockReview[1].CleanlinessReview,
			mockReview[1].ValueReview)

	query := `select \*\ from reviews r where exp_id = \? AND is_deleted = 0 AND is_active = 1`
	if rating != 0 {
		query = query + ` AND r.values = ` + strconv.Itoa(rating)
	}
	if userId != "" {
		query = query + ` AND r.user_id = '` + userId + `' `
	}
	sortBy := "oldestdate"
	if sortBy != "" {
		if sortBy == "ratingup" {
			query = query + ` ORDER BY r.values DESC`
		} else if sortBy == "ratingdown" {
			query = query + ` ORDER BY r.values ASC`
		} else if sortBy == "latestdate" {
			query = query + ` ORDER BY created_date DESC`
		} else if sortBy == "oldestdate" {
			query = query + ` ORDER BY created_date ASC`
		}
	}
	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ReviewRepo.NewReviewRepository(db)

	anArticle, err := a.GetByExpId(context.TODO(),mockReview[0].ExpId,sortBy,rating,0,0,userId)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestGetByExpIdWithPagination(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","values","desc","exp_id","user_id","guide_review","activities_review",
		"service_review","cleanlinessreview","value_review"}).
		AddRow(mockReview[0].Id, mockReview[0].CreatedBy,mockReview[0].CreatedDate,mockReview[0].ModifiedBy,
			mockReview[0].ModifiedDate,mockReview[0].DeletedBy,mockReview[0].DeletedDate,mockReview[0].IsDeleted,
			mockReview[0].IsActive,mockReview[0].Values,mockReview[0].Desc,mockReview[0].ExpId,mockReview[0].UserId,
			mockReview[0].GuideReview,mockReview[0].ActivitiesReview,mockReview[0].ServiceReview,mockReview[0].CleanlinessReview,
			mockReview[0].ValueReview).
		AddRow(mockReview[1].Id, mockReview[1].CreatedBy,mockReview[1].CreatedDate,mockReview[1].ModifiedBy,
			mockReview[1].ModifiedDate,mockReview[1].DeletedBy,mockReview[1].DeletedDate,mockReview[1].IsDeleted,
			mockReview[1].IsActive,mockReview[1].Values,mockReview[1].Desc,mockReview[1].ExpId,mockReview[1].UserId,
			mockReview[1].GuideReview,mockReview[1].ActivitiesReview,mockReview[1].ServiceReview,mockReview[1].CleanlinessReview,
			mockReview[1].ValueReview)

	query := `select \*\ from reviews r where exp_id = \? AND is_deleted = 0 AND is_active = 1`
	if rating != 0 {
		query = query + ` AND r.values = ` + strconv.Itoa(rating)
	}
	if userId != "" {
		query = query + ` AND r.user_id = '` + userId + `' `
	}
	sortBy := "oldestdate"
	if sortBy != "" {
		if sortBy == "ratingup" {
			query = query + ` ORDER BY r.values DESC`
		} else if sortBy == "ratingdown" {
			query = query + ` ORDER BY r.values ASC`
		} else if sortBy == "latestdate" {
			query = query + ` ORDER BY created_date DESC`
		} else if sortBy == "oldestdate" {
			query = query + ` ORDER BY created_date ASC`
		}
	}
	query = query + ` LIMIT \? OFFSET \?`
	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ReviewRepo.NewReviewRepository(db)

	anArticle, err := a.GetByExpId(context.TODO(),mockReview[0].ExpId,sortBy,rating,2,0,userId)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestGetByExpIdWithPaginationErrorQuery(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","values","desc","exp_id","user_id","guide_review","activities_review",
		"service_review","cleanlinessreview","value_review"}).
		AddRow(mockReview[0].Id, mockReview[0].CreatedBy,mockReview[0].CreatedDate,mockReview[0].ModifiedBy,
			mockReview[0].ModifiedDate,mockReview[0].DeletedBy,mockReview[0].DeletedDate,mockReview[0].IsDeleted,
			mockReview[0].IsActive,mockReview[0].Values,mockReview[0].Desc,mockReview[0].ExpId,mockReview[0].UserId,
			mockReview[0].GuideReview,mockReview[0].ActivitiesReview,mockReview[0].ServiceReview,mockReview[0].CleanlinessReview,
			mockReview[0].ValueReview).
		AddRow(mockReview[1].Id, mockReview[1].CreatedBy,mockReview[1].CreatedDate,mockReview[1].ModifiedBy,
			mockReview[1].ModifiedDate,mockReview[1].DeletedBy,mockReview[1].DeletedDate,mockReview[1].IsDeleted,
			mockReview[1].IsActive,mockReview[1].Values,mockReview[1].Desc,mockReview[1].ExpId,mockReview[1].UserId,
			mockReview[1].GuideReview,mockReview[1].ActivitiesReview,mockReview[1].ServiceReview,mockReview[1].CleanlinessReview,
			mockReview[1].ValueReview)

	query := `select \*\ from reviewsadasdasdas r where exp_id = \? AND is_deleted = 0 AND is_active = 1`
	if rating != 0 {
		query = query + ` AND r.values = ` + strconv.Itoa(rating)
	}
	if userId != "" {
		query = query + ` AND r.user_id = '` + userId + `' `
	}
	sortBy := "oldestdate"
	if sortBy != "" {
		if sortBy == "ratingup" {
			query = query + ` ORDER BY r.values DESC`
		} else if sortBy == "ratingdown" {
			query = query + ` ORDER BY r.values ASC`
		} else if sortBy == "latestdate" {
			query = query + ` ORDER BY created_date DESC`
		} else if sortBy == "oldestdate" {
			query = query + ` ORDER BY created_date ASC`
		}
	}
	query = query + ` LIMIT \? OFFSET \?`
	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ReviewRepo.NewReviewRepository(db)

	anArticle, err := a.GetByExpId(context.TODO(),mockReview[0].ExpId,sortBy,rating,2,0,userId)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestInsert(t *testing.T) {
	//user := "test"
	//now := time.Now()
	a := mockReview[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT reviews SET id=\? , created_by=\? , created_date=\? , modified_by=\?, modified_date=\? ,
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , reviews.values=\?,reviews.desc=\?,exp_id=\?,user_id=\?,
				guide_review=\?,activities_review=\?,service_review=\?,cleanliness_review=\?,value_review=\?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs( a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.Values, a.Desc,
		a.ExpId, a.UserId, a.GuideReview, a.ActivitiesReview, a.ServiceReview, a.CleanlinessReview, a.ValueReview).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ReviewRepo.NewReviewRepository(db)

	id, err := i.Insert(context.TODO(), a)
	assert.NoError(t, err)
	assert.Equal(t, id, a.Id)
}
func TestInsertErrorExec(t *testing.T) {
	//user := "test"
	//now := time.Now()
	a := mockReview[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT reviews SET id=\? , created_by=\? , created_date=\? , modified_by=\?, modified_date=\? ,
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , reviews.values=\?,reviews.desc=\?,exp_id=\?,user_id=\?,
				guide_review=\?,activities_review=\?,service_review=\?,cleanliness_review=\?,value_review=\?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs( a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.Values, a.Desc,
		a.ExpId, a.UserId, a.GuideReview, a.ActivitiesReview, a.ServiceReview, a.CleanlinessReview, a.ValueReview,a.ValueReview).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ReviewRepo.NewReviewRepository(db)

	_, err = i.Insert(context.TODO(), a)
	assert.Error(t, err)
}
