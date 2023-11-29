package models

import "gorm.io/gorm"

type Reply struct {
	gorm.Model
	CommentId uint   `json:"commentId"`
	Content   string `json:"content"`
	UserID    uint   `json:"userId"`
}
