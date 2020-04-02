package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	model "github.com/models"
	"time"
)

func main() {
	db, err := gorm.Open("mysql", "AdminCgo@api-blog-cgo-mysqldbserver:Standar123.@(api-blog-cgo-mysqldbserver.mysql.database.azure.com)/cgo_indonesia?charset=utf8&parseTime=True&loc=Local")
	if err != nil{
		fmt.Println(err)
	}
	//minimumBooking := model.MinimumBooking{}
	//merchant := model.Merchant{}
	user := model.FAQ{}
	error := db.AutoMigrate(&user)
	if error != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add_table_FAQ",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}
	pointRules := model.BookingExp{}
	errorpointRules := db.AutoMigrate(&pointRules).AddForeignKey("exp_id","experiences(id)","RESTRICT", "RESTRICT")
	if errorpointRules != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add_table_booking_exp",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}
	pointRule := model.BookingExp{}
	errorpointRule := db.Model(&pointRule).AddForeignKey("user_id","users(id)","RESTRICT", "RESTRICT")
	if errorpointRule != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add_Foregn_key_userId_booking_exp",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}
	facilities := model.PaymentMethod{}
	errorfacilities := db.AutoMigrate(&facilities)
	if errorfacilities != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add_table_Payment_method",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}
	exlusionService := model.Payment{}
	errorexlusionService := db.AutoMigrate(&exlusionService).AddForeignKey("booking_exp_id","booking_exps(id)","RESTRICT", "RESTRICT")
	if errorexlusionService != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add_table_Payment",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}
	exlusionServices := model.Payment{}
	errorexlusionServices := db.Model(&exlusionServices).AddForeignKey("promo_id","promos(id)","RESTRICT", "RESTRICT")
	if errorexlusionServices != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add_Foregn_key_promo_id_Payment",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}
	exlusionServicess := model.Payment{}
	errorexlusionServicess := db.Model(&exlusionServicess).AddForeignKey("payment_method_id","payment_methods(id)","RESTRICT", "RESTRICT")
	if errorexlusionServicess != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add_Foregn_key_payment_method_id_Payment",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}

	exlusionServicesss := model.Payment{}
	errorexlusionServicesss := db.Model(&exlusionServicesss).AddForeignKey("experience_payment_id","experience_payments(id)","RESTRICT", "RESTRICT")
	if errorexlusionServicesss != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add_Foregn_key_experience_payment_id_Payment",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}
	db.Close()

}
