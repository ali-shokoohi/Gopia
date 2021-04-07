package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

var db = getDatabase()

type Article struct {
	gorm.Model
	Title   string `gorm:"not null" json:"Title"`
	Desc    string `gorm:"not null" json:"Descriptions"`
	Content string `gorm:"not null" json:"Content"`
	UserID  uint   `gorm:"default:1"`
}

var Articles []Article

type User struct {
	gorm.Model
	FirstName string    `gorm:"not null" json:"first_name"`
	LastName  string    `gorm:"not null" json:"last_name"`
	Email     string    `gorm:"not null;unique" json:"email"`
	Age       string    `gorm:"not null" json:"age"`
	Username  string    `gorm:"not null;unique"`
	Password  string    `gorm:"not null"`
	Articles  []Article `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

var Users []User

var models map[string]interface{}
var objects map[string]interface{}
var objectsJsonMap map[string][]interface{}

var objectsJson []byte

//Validate incoming user details...
func (user *User) Validate() (string, bool) {
	if !strings.Contains(user.Email, "@") {
		return "Email address is required", false
	}
	if len(user.Password) < 6 {
		return "Strong Password is required", false
	}
	//Email and Username must be unique
	temp := &User{}
	//check for errors and duplicate emails
	err := db.Table("users").Where("email = ?", user.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "Connection error. Please retry", false
	}
	if temp.Email != "" {
		return "Email address already in use by another user.", false
	}
	errUser := db.Table("users").Where("Username = ?", user.Username).First(temp).Error
	if errUser != nil && errUser != gorm.ErrRecordNotFound {
		return "Connection error. Please retry", false
	}
	if temp.Username != "" {
		return "Username already in use by another user.", false
	}
	return "Requirement passed", true
}

func (user *User) Create() (string, bool) {
	if resp, ok := user.Validate(); !ok {
		return resp, false
	}
	hasher := md5.New()
	hasher.Write([]byte(user.Password))
	user.Password = hex.EncodeToString(hasher.Sum(nil))
	db.Create(user)
	return "Account has been created", true
}

func autoMigrate(models map[string]interface{}) {
	for index, model := range models {
		fmt.Printf("%s: %v\n", index, model)
		db.AutoMigrate(model)
	}
}

func reloadObjects() {
	// Convert objects map to a []byte map
	objectsJson, _ = json.Marshal(objects)
	// Again convert to a string map
	json.Unmarshal(objectsJson, &objectsJsonMap)
}

func perpareModels() (map[string]interface{}, map[string]interface{}) {
	models = make(map[string]interface{})
	objects = make(map[string]interface{})
	models["article"] = Article{}
	models["user"] = User{}
	autoMigrate(models)
	db.Find(&Articles)
	db.Preload("Articles").Find(&Users)
	objects["articles"] = Articles
	objects["users"] = Users
	reloadObjects()
	return models, objects
}
