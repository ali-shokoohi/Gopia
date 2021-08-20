package models

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"gitlab.com/greenly/go-rest-api/database"
	"gorm.io/gorm"
)

// DB database client
var DB *gorm.DB = new(database.Database).GetDatabase()

// AppCache cache client
var AppCache *cache.Cache = new(database.Database).GetCache()

var models map[string]interface{}

func autoMigrate(models map[string]interface{}) {
	for index, model := range models {
		fmt.Printf("%s: %v\n", index, model)
		DB.AutoMigrate(model)
	}
}

// Model type
type Model struct{}

// PerpareModels as startup
func (model *Model) PerpareModels() {
	models = make(map[string]interface{})
	models["article"] = Article{}
	models["user"] = User{}
	models["comment"] = Comment{}
	models["like"] = Like{}
	DB.Preload("Articles").Preload("Comments").Find(&Users)
	DB.Preload("Comments").Find(&Articles)
	DB.Preload("Replies").Find(&Comments)
	DB.Find(&Likes)
	AppCache.Set("users", Users, 24*time.Hour)
	AppCache.Set("articles", Articles, 24*time.Hour)
	AppCache.Set("comments", Comments, 24*time.Hour)
	AppCache.Set("likes", Likes, 24*time.Hour)
	autoMigrate(models)
}
