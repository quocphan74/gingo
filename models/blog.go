package models

import "gorm.io/gorm"

type Blog struct {
	gorm.Model
	Title   string `json:"title" form:"title" binding:"required"`
	Desc    string `json:"desc" form:"desc" binding:"required"`
	Content string `json:"content" form:"content" binding:"required"`
	Image   string `json:"image" form:"image" binding:"required"`
	UserID  uint   `json:"userId"`
	// User    User   `gorm:"foreignKey:UserID"`
}
