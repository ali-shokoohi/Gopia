package models

import "gorm.io/gorm"

// Comment gorm database model type
type Comment struct {
	gorm.Model
	UserID    uint
	ArticleID uint
	Message   string     `gorm:"not null" json:"message"`
	Replies   []*Comment `gorm:"many2many:comment_replies;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"replies"`
	Likes     []Like     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Comments List of all comments
var Comments []Comment
