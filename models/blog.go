package models

import "gorm.io/gorm"

type Blog struct {
	gorm.Model
	Title   string `json:"title" form:"title" gorm:"type:varchar(255);not null"`
	Desc    string `json:"desc" form:"desc" gorm:"type:varchar(255)"`
	Content string `json:"content" form:"content" gorm:"type:varchar(255)"`
	Image   string `json:"image" form:"image" gorm:"type:varchar(255)"`
	UserID  uint   `json:"userId" form:"userId" gorm:"type:int(2)"`
	Comment []*Comment
	User    User    `gorm:"foreignKey:UserID"`
	Like    []*Like `gorm:"polymorphic:Target;polymorphicValue:blog"`
}
