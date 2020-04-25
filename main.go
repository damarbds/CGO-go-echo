package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	_merchantUcase "github.com/auth/merchant/usecase"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	_echoMiddleware "github.com/labstack/echo/middleware"

	_articleHttpDeliver "github.com/bxcodec/go-clean-arch/article/delivery/http"
	_articleRepo "github.com/bxcodec/go-clean-arch/article/repository"
	_articleUcase "github.com/bxcodec/go-clean-arch/article/usecase"
	_authorRepo "github.com/bxcodec/go-clean-arch/author/repository"

	_cpcRepo "github.com/service/cpc/repository"
	_expPhotosHttpDeliver "github.com/service/exp_photos/delivery/http"
	_expPhotosRepo "github.com/service/exp_photos/repository"
	_expPhotosUcase "github.com/service/exp_photos/usecase"

	_expAvailabilityRepo "github.com/service/exp_availability/repository"
	_experienceHttpDeliver "github.com/service/experience/delivery/http"
	_experienceRepo "github.com/service/experience/repository"
	_experienceUcase "github.com/service/experience/usecase"

	_harborsHttpDeliver "github.com/service/harbors/delivery/http"
	_harborsRepo "github.com/service/harbors/repository"
	_harborsUcase "github.com/service/harbors/usecase"

	_expPaymentTypeHttpDeliver "github.com/transactions/experience_payment_type/delivery/http"
	_expPaymentTypeRepo "github.com/transactions/experience_payment_type/repository"
	_expPaymentTypeUcase "github.com/transactions/experience_payment_type/usecase"

	//"github.com/bxcodec/go-clean-arch/middleware"
	_isHttpDeliver "github.com/auth/identityserver/delivery/http"
	_isUcase "github.com/auth/identityserver/usecase"

	_adminHttpDeliver "github.com/auth/admin/delivery/http"
	_adminRepo "github.com/auth/admin/repository"
	_adminUcase "github.com/auth/admin/usecase"
	_merchantHttpDeliver "github.com/auth/merchant/delivery/http"
	_merchantRepo "github.com/auth/merchant/repository"
	_userMerchantHttpDeliver "github.com/auth/user_merchant/delivery/http"
	_userMerchantRepo "github.com/auth/user_merchant/repository"
	_userMerchantUcase "github.com/auth/user_merchant/usecase"

	_promoHttpDeliver "github.com/service/promo/delivery/http"
	_promoRepo "github.com/service/promo/repository"
	_promoUcase "github.com/service/promo/usecase"

	_experienceAddOnHttpDeliver "github.com/product/experience_add_ons/delivery/http"
	_experienceAddOnRepo "github.com/product/experience_add_ons/repository"
	_experienceAddOnUcase "github.com/product/experience_add_ons/usecase"

	_reviewsHttpDeliver "github.com/product/reviews/delivery/http"
	_reviewsRepo "github.com/product/reviews/repository"
	_reviewsUcase "github.com/product/reviews/usecase"

	_userHttpDeliver "github.com/auth/user/delivery/http"
	_userRepo "github.com/auth/user/repository"
	_userUcase "github.com/auth/user/usecase"
	_paymentRepo "github.com/service/exp_payment/repository"

	_fAQHttpDeliver "github.com/misc/faq/delivery/handler/http"
	_fAQRepo "github.com/misc/faq/repository"
	_fAQUcase "github.com/misc/faq/usecase"

	_bookingExpHttpDeliver "github.com/booking/booking_exp/delivery/http"
	_bookingExpRepo "github.com/booking/booking_exp/repository"
	_bookingExpUcase "github.com/booking/booking_exp/usecase"

	_inspirationRepo "github.com/service/exp_inspiration/repository"
	_typesRepo "github.com/service/exp_types/repository"

	_paymentMethodHttpDeliver "github.com/transactions/payment_methods/delivery/http"
	_paymentMethodRepo "github.com/transactions/payment_methods/repository"
	_paymentMethodUcase "github.com/transactions/payment_methods/usecase"

	_paymentHttpDeliver "github.com/transactions/payment/delivery/http"
	_paymentTrRepo "github.com/transactions/payment/repository"
	_paymentUcase "github.com/transactions/payment/usecase"

	_wishlistHttpHandler "github.com/profile/wishlists/delivery/http"
	_wishlistRepo "github.com/profile/wishlists/repository"
	_wishlistUcase "github.com/profile/wishlists/usecase"

	_notifHttpHandler "github.com/misc/notif/delivery/http"
	_notifRepo "github.com/misc/notif/repository"
	_notifUcase "github.com/misc/notif/usecase"

	_facilityHttpHandler "github.com/service/facilities/delivery/http"
	_facilityRepo "github.com/service/facilities/repository"
	_facilityUcase "github.com/service/facilities/usecase"

	_transportationHttpHandler "github.com/service/transportation/delivery/http"
	_transportationRepo "github.com/service/transportation/repository"
	_transportationUcase "github.com/service/transportation/usecase"

	_scheduleHttpHandler "github.com/service/schedule/delivery/http"
	_scheduleRepo "github.com/service/schedule/repository"
	_scheduleUsecase "github.com/service/schedule/usecase"

	_timeOptionsRepo "github.com/service/time_options/repository"

	_transactionHttpHandler "github.com/transactions/transaction/delivery/http"
	_transactionRepo "github.com/transactions/transaction/repository"
	_transactionUcase "github.com/transactions/transaction/usecase"

	_balanceHistoryHttpHandler "github.com/transactions/balance_history/delivery/http"
	_balanceHistoryRepo "github.com/transactions/balance_history/repository"
	_balanceHistoryUcase "github.com/transactions/balance_history/usecase"

	_midtransHttpHandler "github.com/third-party/midtrans/delivery/http"

	_filterActivityTypeRepo "github.com/service/filter_activity_type/repository"

	_currencyHttpHandler "github.com/misc/currency/delivery/http"
	_currencyUsecase "github.com/misc/currency/usecase"
)

// func init() {
// 	viper.SetConfigFile(`config.json`)
// 	err := viper.ReadInConfig()
// 	if err != nil {
// 		panic(err)
// 	}

// 	if viper.GetBool(`debug`) {
// 		fmt.Println("Service RUN on DEBUG mode")
// 	}
// }

func main() {
	// dbHost := viper.GetString(`database.host`)
	// dbPort := viper.GetString(`database.port`)
	// dbUser := viper.GetString(`database.user`)
	// dbPass := viper.GetString(`database.pass`)
	// dbName := viper.GetString(`database.name`)
	// baseUrlis := viper.GetString(`identityServer.baseUrl`)
	// basicAuth := viper.GetString(`identityServer.basicAuth`)
	dbHost := "api-blog-cgo-mysqldbserver.mysql.database.azure.com"
	dbPort := "3306"
	dbUser := "AdminCgo@api-blog-cgo-mysqldbserver"
	dbPass := "Standar123."
	dbName := "cgo_indonesia"
	baseUrlis := "http://identity-server-cgo-indonesia.azurewebsites.net"
	basicAuth := "cm9jbGllbnQ6c2VjcmV0"
	accountStorage := "cgostorage"
	accessKeyStorage := "OwvEOlzf6e7QwVoV0H75GuSZHpqHxwhYnYL9QbpVPgBRJn+26K26aRJxtZn7Ip5AhaiIkw9kH11xrZSscavXfQ=="
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)
	// if err != nil && viper.GetBool("debug") {
	// 	fmt.Println(err)
	// }
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	//middL := middleware.InitMiddleware()

	//e.Use(middL.CORS)
	e.Use(_echoMiddleware.CORSWithConfig(_echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))
	//e.Use(_echoMiddleware.CORS())
	filterActivityTypeRepo := _filterActivityTypeRepo.NewFilterActivityTypeRepository(dbConn)
	expAvailabilityRepo := _expAvailabilityRepo.NewExpavailabilityRepository(dbConn)
	bookingExpRepo := _bookingExpRepo.NewbookingExpRepository(dbConn)
	fAQRepo := _fAQRepo.NewReviewRepository(dbConn)
	experienceAddOnRepo := _experienceAddOnRepo.NewexperienceRepository(dbConn)
	exp_photos := _expPhotosRepo.Newexp_photosRepository(dbConn)
	harborsRepo := _harborsRepo.NewharborsRepository(dbConn)
	cpcRepo := _cpcRepo.NewcpcRepository(dbConn)
	experienceRepo := _experienceRepo.NewexperienceRepository(dbConn)
	merchantRepo := _merchantRepo.NewmerchantRepository(dbConn)
	userRepo := _userRepo.NewuserRepository(dbConn)
	authorRepo := _authorRepo.NewMysqlAuthorRepository(dbConn)
	ar := _articleRepo.NewMysqlArticleRepository(dbConn)
	paymentRepo := _paymentRepo.NewExpPaymentRepository(dbConn)
	reviewsRepo := _reviewsRepo.NewReviewRepository(dbConn)
	promoRepo := _promoRepo.NewpromoRepository(dbConn)
	typesRepo := _typesRepo.NewExpTypeRepository(dbConn)
	inspirationRepo := _inspirationRepo.NewExpInspirationRepository(dbConn)
	paymentMethodRepo := _paymentMethodRepo.NewPaymentMethodRepository(dbConn)
	paymentTrRepo := _paymentTrRepo.NewPaymentRepository(dbConn)
	wlRepo := _wishlistRepo.NewWishListRepository(dbConn)
	expPaymentTypeRepo := _expPaymentTypeRepo.NewExperiencePaymentTypeRepository(dbConn)
	notifRepo := _notifRepo.NewNotifRepository(dbConn)
	facilityRepo := _facilityRepo.NewFacilityRepository(dbConn)
	adminRepo := _adminRepo.NewadminRepository(dbConn)
	transportationRepo := _transportationRepo.NewTransportationRepository(dbConn)
	timeOptionsRepo := _timeOptionsRepo.NewTimeOptionsRepository(dbConn)
	schedulerRepo := _scheduleRepo.NewScheduleRepository(dbConn)
	transactionRepo := _transactionRepo.NewTransactionRepository(dbConn)
	balanceHistoryRepo := _balanceHistoryRepo.NewbalanceHistoryRepository(dbConn)
	userMerchantRepo := _userMerchantRepo.NewuserMerchantRepository(dbConn)

	timeoutContext := time.Duration(30) * time.Second

	expPaymentTypeUsecase := _expPaymentTypeUcase.NewexperiencePaymentTypeUsecase(expPaymentTypeRepo, timeoutContext)
	fAQUsecase := _fAQUcase.NewfaqUsecase(fAQRepo, timeoutContext)
	reivewsUsecase := _reviewsUcase.NewreviewsUsecase(reviewsRepo, userRepo, timeoutContext)
	experienceAddOnUsecase := _experienceAddOnUcase.NewharborsUsecase(experienceAddOnRepo, timeoutContext)
	harborsUsecase := _harborsUcase.NewharborsUsecase(harborsRepo, timeoutContext)
	exp_photosUsecase := _expPhotosUcase.Newexp_photosUsecase(exp_photos, timeoutContext)
	isUsecase := _isUcase.NewidentityserverUsecase(baseUrlis, basicAuth, accountStorage, accessKeyStorage)

	adminUsecase := _adminUcase.NewadminUsecase(adminRepo, isUsecase, timeoutContext)
	merchantUsecase := _merchantUcase.NewmerchantUsecase(userMerchantRepo, merchantRepo, experienceRepo, transportationRepo, isUsecase, adminUsecase, timeoutContext)

	promoUsecase := _promoUcase.NewPromoUsecase(promoRepo, adminUsecase, timeoutContext)
	experienceUsecase := _experienceUcase.NewexperienceUsecase(
		filterActivityTypeRepo,
		experienceAddOnRepo,
		expAvailabilityRepo,
		exp_photos,
		experienceRepo,
		harborsRepo,
		cpcRepo,
		paymentRepo,
		reviewsRepo,
		typesRepo,
		inspirationRepo,
		merchantUsecase,
		timeoutContext,
	)
	userUsecase := _userUcase.NewuserUsecase(userRepo, isUsecase, adminUsecase, timeoutContext)
	au := _articleUcase.NewArticleUsecase(ar, authorRepo, timeoutContext)
	pmUsecase := _paymentMethodUcase.NewPaymentMethodUsecase(paymentMethodRepo, timeoutContext)
	paymentUsecase := _paymentUcase.NewPaymentUsecase(paymentTrRepo, userUsecase, bookingExpRepo, timeoutContext)
	bookingExpUcase := _bookingExpUcase.NewbookingExpUsecase(bookingExpRepo, userUsecase, merchantUsecase, isUsecase, experienceRepo, transactionRepo, timeoutContext)
	wlUcase := _wishlistUcase.NewWishlistUsecase(wlRepo, userUsecase, experienceRepo, paymentRepo, reviewsRepo, timeoutContext)
	notifUcase := _notifUcase.NewNotifUsecase(notifRepo, merchantUsecase, timeoutContext)
	facilityUcase := _facilityUcase.NewFacilityUsecase(facilityRepo, timeoutContext)
	transportationUcase := _transportationUcase.NewTransportationUsecase(transportationRepo, merchantUsecase, schedulerRepo, timeOptionsRepo, timeoutContext)
	transactionUcase := _transactionUcase.NewTransactionUsecase(transactionRepo, timeoutContext)
	scheduleUcase := _scheduleUsecase.NewScheduleUsecase(transportationRepo, merchantUsecase, schedulerRepo, timeOptionsRepo, experienceRepo, expAvailabilityRepo, timeoutContext)
	balanceHistoryUcase := _balanceHistoryUcase.NewBalanceHistoryUsecase(balanceHistoryRepo, merchantUsecase, timeoutContext)
	userMerchantUcase := _userMerchantUcase.NewuserMerchantUsecase(userMerchantRepo, merchantUsecase, isUsecase, adminUsecase, timeoutContext)
	currencyUcase := _currencyUsecase.NewCurrencyUsecase(timeoutContext)

	_userMerchantHttpDeliver.NewuserMerchantHandler(e, userMerchantUcase)
	_scheduleHttpHandler.NewScheduleHandler(e, promoUsecase, scheduleUcase)
	_adminHttpDeliver.NewadminHandler(e, adminUsecase)
	_expPaymentTypeHttpDeliver.NewexpPaymentTypeHandlerHandler(e, expPaymentTypeUsecase)
	_bookingExpHttpDeliver.Newbooking_expHandler(e, bookingExpUcase)
	_fAQHttpDeliver.NewfaqHandler(e, fAQUsecase)
	_reviewsHttpDeliver.NewreviewsHandler(e, reivewsUsecase)
	_experienceAddOnHttpDeliver.Newexperience_add_onsHandler(e, experienceAddOnUsecase)
	_harborsHttpDeliver.NewharborsHandler(e, harborsUsecase)
	_expPhotosHttpDeliver.Newexp_photosHandler(e, exp_photosUsecase)
	_experienceHttpDeliver.NewexperienceHandler(e, experienceUsecase, isUsecase)
	_isHttpDeliver.NewisHandler(e, merchantUsecase, userUsecase, adminUsecase, isUsecase)
	_userHttpDeliver.NewuserHandler(e, userUsecase, isUsecase)
	_merchantHttpDeliver.NewmerchantHandler(e, merchantUsecase, isUsecase)
	_articleHttpDeliver.NewArticleHandler(e, au)
	_promoHttpDeliver.NewpromoHandler(e, promoUsecase, isUsecase)
	_paymentMethodHttpDeliver.NewPaymentMethodHandler(e, pmUsecase)
	_paymentHttpDeliver.NewPaymentHandler(e, paymentUsecase, bookingExpUcase, bookingExpRepo, paymentMethodRepo)
	_wishlistHttpHandler.NewWishlistHandler(e, wlUcase)
	_notifHttpHandler.NewNotifHandler(e, notifUcase)
	_facilityHttpHandler.NewFacilityHandler(e, facilityUcase)
	_transportationHttpHandler.NewtransportationHandler(e, transportationUcase)
	_transactionHttpHandler.NewTransactionHandler(e, transactionUcase)
	_balanceHistoryHttpHandler.NewBalanceHistoryHandler(e, balanceHistoryUcase)
	_midtransHttpHandler.NewMidtransHandler(e, bookingExpRepo, experienceRepo, transactionRepo, isUsecase)
	_currencyHttpHandler.NewCurrencyHandler(e, currencyUcase)

	log.Fatal(e.Start(":9090"))
}
