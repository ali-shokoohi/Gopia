package main

import (
	"fmt"

	"gorm.io/gorm"
)

var db = getDatabase()

type Article struct {
	gorm.Model
	Title   string `gorm:"not null" json:"Title"`
	Desc    string `gorm:"not null" json:"Descriptions"`
	Content string `gorm:"not null" json:"Content"`
}

var Articles []Article

type User struct {
	gorm.Model
	FirstName string `gorm:"not null" json:"first_name"`
	LastName  string `gorm:"not null" json:"last_name"`
	Email     string `gorm:"not null" json:"email"`
	Age       string `gorm:"not null" json:"age"`
	Username  string `gorm:"not null" json:"username"`
	Password  string `gorm:"not null" json:"password"`
}

var Users []User

var models map[string]interface{}
var objects map[string]interface{}

func autoMigrate(models map[string]interface{}) {
	for index, model := range models {
		fmt.Printf("%s: %v\n", index, model)
		db.AutoMigrate(model)
	}
}

func perpareModels() (map[string]interface{}, map[string]interface{}) {
	models = make(map[string]interface{})
	objects = make(map[string]interface{})
	models["article"] = Article{}
	models["user"] = User{}
	autoMigrate(models)
	db.Find(&Articles)
	db.Find(&Users)
	objects["articles"] = Articles
	objects["users"] = Users
	return models, objects
}
