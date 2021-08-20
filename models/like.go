package models

import "gorm.io/gorm"

// Like - gorm database model type
type Like struct {
	gorm.Model
	UserID    uint `gorm:"not null"`
	ArticleID uint
}

// Likes - List of all likes
var Likes []Like
