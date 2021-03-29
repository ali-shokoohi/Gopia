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

var models map[string]interface{}

func autoMigrate(models map[string]interface{}) {
	for index, model := range models {
		fmt.Printf("%s: %v\n", index, model)
		db.AutoMigrate(model)
	}
}

func perpareModels() map[string]interface{} {
	models = make(map[string]interface{})
	models["article"] = Article{}
	autoMigrate(models)
	return models
}
