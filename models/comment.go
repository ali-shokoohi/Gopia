package models

import "gorm.io/gorm"

// Comment gorm database model type
type Comment struct {
	gorm.Model
	UserID    uint
	ArticleID uint
	Message   string     `gorm:"not null" json:"message"`
	Replies   []*Comment `gorm:"many2many:comment_replies" json:"replies"`
}

// Comments List of all comments
var Comments []Comment
