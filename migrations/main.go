package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	model "github.com/models"
)

func main() {
	db, err := gorm.Open("mysql", "AdminCgo@api-blog-cgo-mysqldbserver:Standar123.@(api-blog-cgo-mysqldbserver.mysql.database.azure.com)/cgo_indonesia?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	//minimumBooking := model.MinimumBooking{}
	//merchant := model.Merchant{}
	user := model.BookingExp{}
	error := db.Model(&user).AddForeignKey("trans_id", "transportations(id)", "RESTRICT", "RESTRICT")
	if error != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add_foregn_key_trans_id_bookingExp",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}
	transportationdestid := model.Schedule{}
	errortransportationdestid := db.AutoMigrate(&transportationdestid).AddForeignKey("arrival_timeoption_id", "times_options(id)", "RESTRICT", "RESTRICT")
	if errortransportationdestid != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add_foregn_key_arrival_time_option_Schedule",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}
	//merchantId := model.Transportation{}
	//errormerchantId := db.Model(&merchantId).AddForeignKey("merchant_id","merchants(id)","RESTRICT", "RESTRICT")
	//if errormerchantId != nil{
	//	migration := model.MigrationHistory{
	//		DescMigration:"Add_foregn_key_merchant_id_Transportation",
	//		Date:  time.Now(),
	//	}
	//
	//	db.Create(&migration)
	//}
	//returnId := model.Transportation{}
	//errorreturnId := db.Model(&returnId).AddForeignKey("return_trans_id","transportations(id)","RESTRICT", "RESTRICT")
	//if errorreturnId != nil{
	//	migration := model.MigrationHistory{
	//		DescMigration:"Add_foregn_key_return_id_Transportation",
	//		Date:  time.Now(),
	//	}
	//
	//	db.Create(&migration)
	//}
	//pointRules := model.Wishlist{}
	//errorpointRules := db.AutoMigrate(&pointRules).AddForeignKey("trans_id","transportations(id)","RESTRICT", "RESTRICT")
	//if errorpointRules != nil{
	//	migration := model.MigrationHistory{
	//		DescMigration:"Add_table_wishlist",
	//		Date:  time.Now(),
	//	}
	//
	//	db.Create(&migration)
	//}
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
