package models

import "gorm.io/gorm"

// Article gorm database model type
type Article struct {
	gorm.Model
	Title    string    `gorm:"not null" json:"Title"`
	Desc     string    `gorm:"not null" json:"Descriptions"`
	Content  string    `gorm:"not null" json:"Content"`
	UserID   uint      `gorm:"default:1"`
	Comments []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Likes    []Like    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Articles List of all articles
var Articles []Article
