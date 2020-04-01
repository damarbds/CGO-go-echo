package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

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

	_experienceHttpDeliver "github.com/service/experience/delivery/http"
	_experienceRepo "github.com/service/experience/repository"
	_experienceUcase "github.com/service/experience/usecase"

	_harborsHttpDeliver "github.com/service/harbors/delivery/http"
	_harborsRepo "github.com/service/harbors/repository"
	_harborsUcase "github.com/service/harbors/usecase"
	//"github.com/bxcodec/go-clean-arch/middleware"
	_isHttpDeliver "github.com/auth/identityserver/delivery/http"
	_isUcase "github.com/auth/identityserver/usecase"

	_merchantHttpDeliver "github.com/auth/merchant/delivery/http"
	_merchantRepo "github.com/auth/merchant/repository"
	_merchantUcase "github.com/auth/merchant/usecase"

	_promoHttpDeliver "github.com/service/promo/delivery/http"
	_promoRepo "github.com/service/promo/repository"
	_promoUcase "github.com/service/promo/usecase"

	_userHttpDeliver "github.com/auth/user/delivery/http"
	_userRepo "github.com/auth/user/repository"
	_userUcase "github.com/auth/user/usecase"
	_paymentRepo "github.com/service/exp_payment/repository"
	_reviewsRepo "github.com/service/reviews/repository"
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
	baseUrlis := "https://identity-server-cgo-indonesia.azurewebsites.net"
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
		AllowMethods: []string{"GET", "POST","PUT","DELETE","OPTIONS"},
	}))
	//e.Use(_echoMiddleware.CORS())
	exp_photos := _expPhotosRepo.Newexp_photosRepository(dbConn)
	harborsRepo := _harborsRepo.NewharborsRepository(dbConn)
	cpcRepo := _cpcRepo.NewcpcRepository(dbConn)
	experienceRepo := _experienceRepo.NewexperienceRepository(dbConn)
	merchantRepo := _merchantRepo.NewmerchantRepository(dbConn)
	userRepo 	:= _userRepo.NewuserRepository(dbConn)
	authorRepo := _authorRepo.NewMysqlAuthorRepository(dbConn)
	ar := _articleRepo.NewMysqlArticleRepository(dbConn)
	paymentRepo := _paymentRepo.NewExpPaymentRepository(dbConn)
	reviewsRepo := _reviewsRepo.NewReviewRepository(dbConn)
	promoRepo := _promoRepo.NewpromoRepository(dbConn)

	timeoutContext := time.Duration(30) * time.Second

	promoUsecase := _promoUcase.NewArticleUsecase(promoRepo,timeoutContext)
	harborsUsecase := _harborsUcase.NewharborsUsecase(harborsRepo,timeoutContext)
	exp_photosUsecase := _expPhotosUcase.Newexp_photosUsecase(exp_photos,timeoutContext)
	experienceUsecase := _experienceUcase.NewexperienceUsecase(experienceRepo,harborsRepo,cpcRepo,paymentRepo,reviewsRepo,timeoutContext)
	isUsecase := _isUcase.NewidentityserverUsecase(baseUrlis, basicAuth,accountStorage,accessKeyStorage)
	userUsecase := _userUcase.NewuserUsecase(userRepo,isUsecase,timeoutContext)
	merchantUsecase := _merchantUcase.NewmerchantUsecase(merchantRepo, isUsecase, timeoutContext)
	au := _articleUcase.NewArticleUsecase(ar, authorRepo, timeoutContext)

	_harborsHttpDeliver.NewharborsHandler(e,harborsUsecase)
	_expPhotosHttpDeliver.Newexp_photosHandler(e,exp_photosUsecase)
	_experienceHttpDeliver.NewexperienceHandler(e,experienceUsecase)
	_isHttpDeliver.NewisHandler(e,merchantUsecase,userUsecase)
	_userHttpDeliver.NewuserHandler(e,userUsecase)
	_merchantHttpDeliver.NewmerchantHandler(e, merchantUsecase)
	_articleHttpDeliver.NewArticleHandler(e, au)
	_promoHttpDeliver.NewpromoHandler(e,promoUsecase)

	log.Fatal(e.Start(":9090"))
}


