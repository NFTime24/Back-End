package db

import (
	"database/sql"
	"fmt"

	"github.com/duke/config"
	"github.com/duke/model"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func ConnectDB() *gorm.DB {
	configuration := config.GetConfig()
	// fmt.Println(configuration.DB_HOST)
	connect_string := fmt.Sprintf(configuration.DB_USERNAME + ":" + configuration.DB_PASSWORD + "@tcp(" + configuration.DB_HOST + ")/" + configuration.DB_NAME)
	sqlDB, err := sql.Open("mysql", connect_string)
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		panic("DB Connection Error")
	}
	return db
}
func Init() {
	configuration := config.GetConfig()
	connect_string := fmt.Sprintf(configuration.DB_USERNAME + ":" + configuration.DB_PASSWORD + "@tcp(" + configuration.DB_HOST + ")/" + configuration.DB_NAME)
	sqlDB, err := sql.Open("mysql", connect_string)
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		panic("DB Connection Error")
	}
	db.AutoMigrate(&model.Artist{}, &model.File{}, &model.User{}, &model.Work{}, &model.Nft{})

}
