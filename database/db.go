package database

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
	sslmode    = os.Getenv("sslmode")
)

// Database type
type Database struct{}

// GetDatabase () *gorm.DB {...} Return a valid database client
func (database *Database) GetDatabase() *gorm.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s TimeZone=Asia/Tehran",
		dbhost, dbport, dbuser, dbpassword, dbname, sslmode)

	//db, err := sql.Open("postgres", psqlInfo)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to database!")

	return db
}
