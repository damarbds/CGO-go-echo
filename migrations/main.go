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
	//user := model.Country{}
	//error := db.AutoMigrate(&user)
	//if error != nil{
	//	migration := model.MigrationHistory{
	//		DescMigration:"Add_table_Country",
	//		Date:  time.Now(),
	//	}
	//
	//	db.Create(&migration)
	//}
	pointRules := model.Experience{}
	errorpointRules := db.Model(&pointRules).AddForeignKey("harbors_id","harbors(id)","RESTRICT", "RESTRICT")
	if errorpointRules != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add_ForegnKey_Harbors_Id_in_experience",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}
	//facilities := model.City{}
	//errorfacilities := db.AutoMigrate(&facilities).AddForeignKey("province_id","provinces(id)","RESTRICT", "RESTRICT")
	//if errorfacilities != nil{
	//	migration := model.MigrationHistory{
	//		DescMigration:"Add_table_City",
	//		Date:  time.Now(),
	//	}
	//
	//	db.Create(&migration)
	//}
	//exlusionService := model.Harbors{}
	//errorexlusionService := db.Model(&exlusionService).AddForeignKey("city_id","cities(id)","RESTRICT", "RESTRICT")
	//if errorexlusionService != nil{
	//	migration := model.MigrationHistory{
	//		DescMigration:"Add_table_Harbors",
	//		Date:  time.Now(),
	//	}
	//
	//	db.Create(&migration)
	//}
	db.Close()

}
