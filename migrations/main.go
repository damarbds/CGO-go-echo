package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	model "github.com/models"
)

func main() {
	////prd
	//db, err := gorm.Open("mysql", "admincgo@cgo-indonesia-prod:k_)V/p53u9z.V{C,@(cgo-indonesia-prod.mysql.database.azure.com)/cgo_indonesia?charset=utf8&parseTime=True&loc=Local")
	//if err != nil {
	//	fmt.Println(err)
	//}

	//dev
	db, err := gorm.Open("mysql", "AdminCgo@api-blog-cgo-mysqldbserver:Standar123.@(api-blog-cgo-mysqldbserver.mysql.database.azure.com)/cgo_indonesia?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	//minimumBooking := model.MinimumBooking{}
	//merchant := model.Merchant{}
	user := model.Promo{}
	error := db.AutoMigrate(&user)
	if error != nil {
		migration := model.MigrationHistory{
			DescMigration: "Alter_table_Promo Add column PromoProductType",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}
	transportationdestid := model.Transaction{}
	errortransportationdestid := db.AutoMigrate(&transportationdestid)
	if errortransportationdestid != nil {
		migration := model.MigrationHistory{
			DescMigration: "alter table currency add column exchange_rates and exchange_currency ",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}
	pointRules := model.Country{}
	errorpointRules := db.AutoMigrate(&pointRules)
	if errorpointRules != nil{
		migration := model.MigrationHistory{
			DescMigration:"Alter table country add columns",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}
	//pointRuless := model.Wishlist{}
	//errorpointRuless := db.Model(&pointRuless).AddForeignKey("exp_id","experiences(id)","RESTRICT", "RESTRICT")
	//if errorpointRuless != nil{
	//	migration := model.MigrationHistory{
	//		DescMigration:"Add_foregn_key_exp_id_wishlist",
	//		Date:  time.Now(),
	//	}
	//
	//	db.Create(&migration)
	//}
	//pointRulesss := model.Wishlist{}
	//errorpointRulesss := db.Model(&pointRulesss).AddForeignKey("user_id","users(id)","RESTRICT", "RESTRICT")
	//if errorpointRulesss != nil{
	//	migration := model.MigrationHistory{
	//		DescMigration:"Add_foregn_key_user_id_wishlist",
	//		Date:  time.Now(),
	//	}
	//
	//	db.Create(&migration)
	//}
	//facilities := model.Transaction{}
	//errorfacilities := db.AutoMigrate(&facilities).AddForeignKey("booking_exp_id","booking_exps(id)","RESTRICT", "RESTRICT")
	//if errorfacilities != nil{
	//	migration := model.MigrationHistory{
	//		DescMigration:"Add_table_Transaction",
	//		Date:  time.Now(),
	//	}
	//
	//	db.Create(&migration)
	//}
	//exlusionService := model.Transaction{}
	//errorexlusionService := db.Model(&exlusionService).AddForeignKey("promo_id","promos(id)","RESTRICT", "RESTRICT")
	//if errorexlusionService != nil{
	//	migration := model.MigrationHistory{
	//		DescMigration:"Add_forefn_key_promo_id_Transaction",
	//		Date:  time.Now(),
	//	}
	//
	//	db.Create(&migration)
	//}
	//exlusionServices := model.Transaction{}
	//errorexlusionServices := db.Model(&exlusionServices).AddForeignKey("payment_method_id","payment_methods(id)","RESTRICT", "RESTRICT")
	//if errorexlusionServices != nil{
	//	migration := model.MigrationHistory{
	//		DescMigration:"Add_Foregn_key_payment_methodId_Transaction",
	//		Date:  time.Now(),
	//	}
	//
	//	db.Create(&migration)
	//}
	//exlusionServicess := model.Transaction{}
	//errorexlusionServicess := db.Model(&exlusionServicess).AddForeignKey("experience_payment_id","experience_payments(id)","RESTRICT", "RESTRICT")
	//if errorexlusionServicess != nil{
	//	migration := model.MigrationHistory{
	//		DescMigration:"Add_Foregn_key_experience_payment_id_Transaction",
	//		Date:  time.Now(),
	//	}
	//
	//	db.Create(&migration)
	//}
	//
	//exlusionServicesss := model.Payment{}
	//errorexlusionServicesss := db.Model(&exlusionServicesss).AddForeignKey("experience_payment_id","experience_payments(id)","RESTRICT", "RESTRICT")
	//if errorexlusionServicesss != nil{
	//	migration := model.MigrationHistory{
	//		DescMigration:"Add_Foregn_key_experience_payment_id_Payment",
	//		Date:  time.Now(),
	//	}
	//
	//	db.Create(&migration)
	//}
	//db.Close()

}
