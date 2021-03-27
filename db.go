package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "go_user"
	password = "pag%5!%zhQ*cjGr^orjZfKC*V65HhPb5"
	dbname   = "go_api"
)

func createTable(db *sql.DB, tableName string, columns string) (*sql.Rows, error) {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", tableName, columns)
	row, err := db.Query(query)
	if err != nil {
		fmt.Printf("Error: '%v' !\n", err)
		return nil, err
	}
	defer row.Close()
	fmt.Printf("Create table: '%v' !\n", row)
	return row, nil
}

func getDatabase() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to database!")

	return db
}
