package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/duke/config"
	"github.com/duke/model"
	"github.com/joho/godotenv"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var err error
var db *gorm.DB

func DbConn() (db *sql.DB) {
	er := godotenv.Load(".env")
	if er != nil {
		panic(er.Error())
	}
	dbDriver := os.Getenv("DB_Driver")
	dbUser := os.Getenv("DB_User")
	dbPass := os.Getenv("DB_Password")
	dbName := os.Getenv("DB_Name")
	dbHost := os.Getenv("DB_HOST")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+")/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func Init() {
	configuration := config.GetConfig()
	connect_string := fmt.Sprintf(configuration.DB_USERNAME + ":" + configuration.DB_PASSWORD + "@tcp(" + configuration.DB_HOST + ")/" + configuration.DB_NAME)
	sqlDB, err := sql.Open("mysql", connect_string)
	db, err = gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		panic("DB Connection Error")
	}
	db.AutoMigrate(&model.Artist{}, &model.File{}, &model.User{}, &model.Work{}, &model.Nft{}, &model.Like{}, &model.Exhibition{}, &model.Fantalk{})
}

func DbManager() *gorm.DB { return db }
