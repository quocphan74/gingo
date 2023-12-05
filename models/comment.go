package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content  string     `json:"content" form:"content" gorm:"type:varchar(255)"`
	BlogID   uint       `json:"blogId" form:"blogId" gorm:"type:int"`
	UserID   uint       `json:"userId" form:"userId" gorm:"type:int"`
	ParentID *uint      `json:"parentID,omitempty"`
	Parent   *Comment   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Replies  []*Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
	Like     []Like     `gorm:"foreignKey:CommentId"`
}
