package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	model "github.com/models"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	db, err := gorm.Open("mysql", "AdminCgo@api-blog-cgo-mysqldbserver:Standar123.@(api-blog-cgo-mysqldbserver.mysql.database.azure.com)/cgo_indonesia?charset=utf8&parseTime=True&loc=Local")
	if err != nil{
		fmt.Println(err)
	}
	user := model.User{}
	db.AutoMigrate(&user)

	fmt.Println("test")
	db.Close()

}
