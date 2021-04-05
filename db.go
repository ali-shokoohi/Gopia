package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	_          = godotenv.Load()
	dbhost     = os.Getenv("dbhost")
	dbport     = os.Getenv("dbport")
	dbuser     = os.Getenv("dbuser")
	dbpassword = os.Getenv("dbpassword")
	dbname     = os.Getenv("dbname")
)

func getDatabase() *gorm.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=require TimeZone=Asia/Tehran",
		dbhost, dbport, dbuser, dbpassword, dbname)

	//db, err := sql.Open("postgres", psqlInfo)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to database!")

	return db
}
