package models

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	UserID    uint  `json:"userId,omitempty" `
	BlogId    *uint `json:"blogId,omitempty" `
	CommentId *uint `json:"commentId,omitempty" `
}
