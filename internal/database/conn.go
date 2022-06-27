package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDatabaseConn() *gorm.DB {
	dsn := "root:@tcp(localhost:3306)/course2?parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("can't connect to database")
	}
	fmt.Println("Succes Connection")
	return db
}
