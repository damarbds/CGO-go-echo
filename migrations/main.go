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
	//user := model.Include{}
	//error := db.AutoMigrate(&user)
	//if error != nil {
	//	migration := model.MigrationHistory{
	//		DescMigration: "Add Table include",
	//		Date:          time.Now(),
	//	}
	//
	//	db.Create(&migration)
	//}
	////transportationdestid := model.Transaction{}
	//errortransportationdestid := db.AutoMigrate(&transportationdestid)
	//if errortransportationdestid != nil {
	//	migration := model.MigrationHistory{
	//		DescMigration: "alter table currency add column exchange_rates and exchange_currency ",
	//		Date:          time.Now(),
	//	}
	//
	//	db.Create(&migration)
	//}
	pointRules := model.Notification{}
	errorpointRules := db.AutoMigrate(&pointRules)
	if errorpointRules != nil{
		migration := model.MigrationHistory{
			DescMigration:"Alter table country add columns some notification",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}
	pointRuless := model.Notification{}
	errorpointRuless := db.Model(&pointRuless).AddForeignKey("exp_id","experiences(id)","RESTRICT", "RESTRICT")
	if errorpointRuless != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add_foregn_key_exp_id_notification",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}
	pointRulesss := model.Notification{}
	errorpointRulesss := db.Model(&pointRulesss).AddForeignKey("schedule_id","schedules(id)","RESTRICT", "RESTRICT")
	if errorpointRulesss != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add_foregn_key_schedule_id_notification",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}

	facilities := model.Notification{}
	errorfacilities := db.Model(&facilities).AddForeignKey("booking_exp_id","booking_exps(id)","RESTRICT", "RESTRICT")
	if errorfacilities != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add_foregn_key_booking_exp_id_notification",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}
	exlusionService := model.UserMerchant{}
	errorexlusionService := db.AutoMigrate(&exlusionService)
	if errorexlusionService != nil{
		migration := model.MigrationHistory{
			DescMigration:"Alter table add column fcm token in table user merchant",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}
	exlusionServices := model.User{}
	errorexlusionServices := db.AutoMigrate(&exlusionServices)
	if errorexlusionServices != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add_Column_fcm_token in table user",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}
	exlusionServicesss := model.Guide{}
	errorexlusionServicesss := db.AutoMigrate(&exlusionServicesss)
	if errorexlusionServicesss != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add table guide",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}

	exlusionServicessss := model.GuideLanguage{}
	errorexlusionServicessss := db.AutoMigrate(&exlusionServicessss).AddForeignKey("guide_id","guides(id)","RESTRICT", "RESTRICT")
	if errorexlusionServicessss != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add table guide language",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}

	guideLanguage := model.GuideLanguage{}
	guideLanguageErr := db.Model(&guideLanguage).AddForeignKey("experience_id","experiences(id)","RESTRICT", "RESTRICT")
	if guideLanguageErr != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add Foregn Key Experience Id guide language",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}

	guideReview := model.GuideReviews{}
	guideReviewErr := db.AutoMigrate(&guideReview).AddForeignKey("guide_id","guides(id)","RESTRICT", "RESTRICT")
	if guideReviewErr != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add Table Guide Review",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}

	guideExperience := model.GuideExperience{}
	guideExperienceErr := db.AutoMigrate(&guideExperience).AddForeignKey("guide_id","guides(id)","RESTRICT", "RESTRICT")
	if guideExperienceErr != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add Table Guide Experience",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}

	language := model.Language{}
	languageErr := db.AutoMigrate(&language)
	if languageErr != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add Table Language",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}

	guideExperiencee := model.GuideExperience{}
	guideExperienceeErr := db.Model(&guideExperiencee).AddForeignKey("language_id","languages(id)","RESTRICT", "RESTRICT")
	if guideExperienceeErr != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add Foregn Key language Id Guide Experience",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}

	expLanguage := model.ExpLanguage{}
	expLanguageErr := db.AutoMigrate(&expLanguage).AddForeignKey("language_id","languages(id)","RESTRICT", "RESTRICT")
	if expLanguageErr != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add Table Exp language",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}

	expLanguages := model.ExpLanguage{}
	expLanguagesErr := db.Model(&expLanguages).AddForeignKey("experience_id","experiences(id)","RESTRICT", "RESTRICT")
	if expLanguagesErr != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add Foregn Key Experience ID Table Exp language",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}

	accomodation := model.Accomodation{}
	accomodationErr := db.AutoMigrate(&accomodation)
	if accomodationErr != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add Table Exp Accomodation",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}

	accomodationExperience := model.AccomodationExperience{}
	accomodationExperienceErr := db.AutoMigrate(&accomodationExperience).AddForeignKey("experience_id","experiences(id)","RESTRICT", "RESTRICT")
	if accomodationExperienceErr != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add Table Exp AccomodationExperience",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}

	accomodationExperiences := model.AccomodationExperience{}
	accomodationExperiencesErr := db.Model(&accomodationExperiences).AddForeignKey("accomodation_id","accomodations(id)","RESTRICT", "RESTRICT")
	if accomodationExperiencesErr != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add Foregn Key Accomodation Experience Table Exp AccomodationExperience",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}

	packages := model.Package{}
	packagesErr := db.AutoMigrate(&packages)
	if packagesErr != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add  Table Exp packages",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}

	experiencePayment := model.ExperiencePayment{}
	experiencePaymentErr := db.AutoMigrate(&experiencePayment)
	if experiencePaymentErr != nil{
		migration := model.MigrationHistory{
			DescMigration:"Add Column in Experience Payment",
			Date:  time.Now(),
		}

		db.Create(&migration)
	}

	db.Close()

}
