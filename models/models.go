package models

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/patrickmn/go-cache"
	"gitlab.com/greenly/go-rest-api/database"
	"gorm.io/gorm"
)

// DB database client
var DB *gorm.DB = new(database.Database).GetDatabase()

// AppCache cache client
var AppCache *cache.Cache = new(database.Database).GetCache()

// Article type
type Article struct {
	gorm.Model
	Title    string    `gorm:"not null" json:"Title"`
	Desc     string    `gorm:"not null" json:"Descriptions"`
	Content  string    `gorm:"not null" json:"Content"`
	UserID   uint      `gorm:"default:1"`
	Comments []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Articles List of all articles
var Articles []Article

// User type
type User struct {
	gorm.Model
	FirstName string    `gorm:"not null" json:"first_name"`
	LastName  string    `gorm:"not null" json:"last_name"`
	Email     string    `gorm:"not null;unique" json:"email"`
	Age       string    `gorm:"not null" json:"age"`
	Username  string    `gorm:"not null;unique"`
	Password  string    `gorm:"not null"`
	Articles  []Article `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Comments  []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Token     string    `gorm:"-" sql:"-" json:"token"`
	Admin     bool      `gorm:"not null; default:false" json:"admin"`
}

// Users List of all users
var Users []User

// Comment type
type Comment struct {
	gorm.Model
	UserID    uint
	ArticleID uint
	Message   string     `gorm:"not null" json:"message"`
	Replies   []*Comment `gorm:"many2many:comment_replies" json:"replies"`
}

// Comments List of all comments
var Comments []Comment

/*
Token JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

var models map[string]interface{}

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
	err := DB.Table("users").Where("email = ?", user.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "Connection error. Please retry", false
	}
	if temp.Email != "" {
		return "Email address already in use by another user.", false
	}
	errUser := DB.Table("users").Where("Username = ?", user.Username).First(temp).Error
	if errUser != nil && errUser != gorm.ErrRecordNotFound {
		return "Connection error. Please retry", false
	}
	if temp.Username != "" {
		return "Username already in use by another user.", false
	}
	return "Requirement passed", true
}

// Create new user
func (user *User) Create() (string, bool) {
	if resp, ok := user.Validate(); !ok {
		return resp, false
	}
	hasher := md5.New()
	hasher.Write([]byte(user.Password))
	user.Password = hex.EncodeToString(hasher.Sum(nil))
	DB.Create(&user)
	if user.ID <= 0 {
		return "Failed to create account, connection error.", false
	}
	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString
	return "Account has been created", true
}

// Update one user
func (user *User) Update() (string, bool) {
	hasher := md5.New()
	hasher.Write([]byte(user.Password))
	user.Password = hex.EncodeToString(hasher.Sum(nil))
	DB.Save(&user)
	return "Account has been updated", true
}

// MarshalJSON user as safe
func (user User) MarshalJSON() ([]byte, error) {
	var tmp struct {
		ID        uint
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Email     string    `json:"email"`
		Age       string    `json:"age"`
		Admin     bool      `json:"admin"`
		Articles  []Article `json:"articles"`
		Comments  []Comment `json:"comments"`
	}
	tmp.ID = user.ID
	tmp.FirstName = user.FirstName
	tmp.LastName = user.LastName
	tmp.Email = user.Email
	tmp.Age = user.Age
	tmp.Admin = user.Admin
	tmp.Articles = user.Articles
	return json.Marshal(&tmp)
}

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
	DB.Preload("Articles").Preload("Comments").Find(&Users)
	DB.Preload("Comments").Find(&Articles)
	DB.Preload("Replies").Find(&Comments)
	AppCache.Set("users", Users, 24*time.Hour)
	AppCache.Set("articles", Articles, 24*time.Hour)
	AppCache.Set("comments", Comments, 24*time.Hour)
	autoMigrate(models)
}