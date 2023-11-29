package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content string `json:"content" form:"content" gorm:"type:varchar(255);null"`
	BlogID  uint   `json:"blogId" form:"blogId" gorm:"type:int(2);null"`
	Blog    Blog   `gorm:"foreignKey:BlogID"`
	UserID  uint   `json:"userId" form:"userId" gorm:"type:int(2);null"`
	User    User   `gorm:"foreignKey:UserID"`
}
