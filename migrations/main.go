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
	user := model.ExperienceRules{}
	error := db.AutoMigrate(&user)
	if error != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add_table_ExpRules",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}
	pointRules := model.PointRules{}
	errorpointRules := db.AutoMigrate(&pointRules)
	if errorpointRules != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add_table_pontRules",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}
	facilities := model.Facilities{}
	errorfacilities := db.AutoMigrate(&facilities)
	if errorfacilities != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add_table_facilities",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}
	exlusionService := model.ExclusionService{}
	errorexlusionService := db.AutoMigrate(&exlusionService)
	if errorexlusionService != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add_table_exlusionService",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}
	experienceType := model.ExperienceType{}
	errorexperienceType := db.AutoMigrate(&experienceType)
	if errorexperienceType != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add_table_experienceType",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}
	id_type := model.IdType{}
	errorid_type := db.AutoMigrate(&id_type)
	if errorid_type != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add_table_id_type",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}
	db.Close()

}
