package repository_test

import (
	"context"
	"testing"

	"github.com/models"
	PromoMerchantRepo "github.com/service/promo_merchant/repository"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	imagePath         = "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Facilities/8941695193938718058.jpg"
	transId           = "dalsjdkalsdjlksdj"
	expId             = "werewrwer"
	mockPromoMerchant = []models.PromoMerchant{
		models.PromoMerchant{
			Id:         1,
			PromoId:    "adklasjdklas",
			MerchantId: "kjhjhkjhkj",
		},
		models.PromoMerchant{
			Id:         2,
			PromoId:    "zxcxzcxz",
			MerchantId: "kjhjhkjhkj",
		},
	}

)

func TestGetByMerchantId(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "promo_id", "merchant_id"}).
		AddRow(mockPromoMerchant[0].Id, mockPromoMerchant[0].PromoId, mockPromoMerchant[0].MerchantId).
		AddRow(mockPromoMerchant[1].Id, mockPromoMerchant[1].PromoId, mockPromoMerchant[1].MerchantId)

	query := `SELECT \*\ FROM promo_merchants WHERE promo_id= \?`
	merchantId := mockPromoMerchant[0].MerchantId
	if merchantId != "" {
		query = query + ` and merchant_id = '` + merchantId + `' `
	}

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := PromoMerchantRepo.NewpromoMerchantRepository(db)

	anArticle, err := a.GetByMerchantId(context.TODO(), merchantId, mockPromoMerchant[0].PromoId)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestGetByMerchantIdErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "promo_id", "merchant_id"}).
		AddRow(mockPromoMerchant[0].Id, mockPromoMerchant[0].PromoId, mockPromoMerchant[0].MerchantId).
		AddRow(mockPromoMerchant[1].Id, mockPromoMerchant[1].PromoId, mockPromoMerchant[1].MerchantId)

	query := `SELECT \*\ FROM promo_merchants WHERE promo_id= \?asdasdsa`
	merchantId := mockPromoMerchant[0].MerchantId
	if merchantId != "" {
		query = query + ` and merchant_id = '` + merchantId + `' `
	}

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := PromoMerchantRepo.NewpromoMerchantRepository(db)

	_, err = a.GetByMerchantId(context.TODO(), merchantId, mockPromoMerchant[0].PromoId)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	//assert.Len(t, anArticle, 2)
}
func TestInsert(t *testing.T) {
	a := mockPromoMerchant[0]

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT promo_merchants SET promo_id=\?,merchant_id=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.PromoId, a.MerchantId).WillReturnResult(sqlmock.NewResult(1, 1))

	i := PromoMerchantRepo.NewpromoMerchantRepository(db)

	err = i.Insert(context.TODO(), a)
	assert.NoError(t, err)
	//assert.Equal(t, *id, a.Id)
}
func TestInsertErrorExec(t *testing.T) {

	a := mockPromoMerchant[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT promo_merchants SET promo_id=\?,merchant_id=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.PromoId, a.MerchantId,a.PromoId).WillReturnResult(sqlmock.NewResult(1, 1))

	i := PromoMerchantRepo.NewpromoMerchantRepository(db)

	err = i.Insert(context.TODO(), a)
	assert.Error(t, err)
}
func TestDeleteByMerchantId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "DELETE FROM promo_merchants WHERE merchant_id = \\? AND promo_id=\\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(mockPromoMerchant[0].MerchantId,mockPromoMerchant[0].PromoId).WillReturnResult(sqlmock.NewResult(12, 1))

	a := PromoMerchantRepo.NewpromoMerchantRepository(db)

	err = a.DeleteByMerchantId(context.TODO(), mockPromoMerchant[0].MerchantId, mockPromoMerchant[0].PromoId)
	assert.NoError(t, err)
}
func TestDeleteByMerchantIdErrorExecQueryString(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "DELETE FROM promo_merchants WHERE merchant_id = \\? AND promo_id=\\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(mockPromoMerchant[0].MerchantId,mockPromoMerchant[0].PromoId,mockPromoMerchant).WillReturnResult(sqlmock.NewResult(12, 1))

	a := PromoMerchantRepo.NewpromoMerchantRepository(db)

	err = a.DeleteByMerchantId(context.TODO(), mockPromoMerchant[0].MerchantId, mockPromoMerchant[0].PromoId)
	assert.Error(t, err)
}
