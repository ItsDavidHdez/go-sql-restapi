package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnect() {
	var host = os.Getenv("HOST")
	var user = os.Getenv("USER")
	var password = os.Getenv("PASSWORD")
	var dbName = os.Getenv("DATABASE_NAME")
	var dbPort = os.Getenv("DATABASE_PORT")

	var DSN = "host=" + host + " user=" + string(user) + " password=" + string(password) + " dbname=" + string(dbName) + " port=" + string(dbPort)

	var err error
	DB, err = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("DB Connected!! :)")
	}
}