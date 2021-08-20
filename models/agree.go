package models

import "gorm.io/gorm"

// Agree - gorm database model type
type Agree struct {
	gorm.Model
	UserID    uint `gorm:"not null"`
	CommentID uint
}

// Agrees - List of all Agrees
var Agrees []Agree
