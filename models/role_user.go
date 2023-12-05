package models

import "gorm.io/gorm"

type RoleUser struct {
	gorm.Model
	UserID uint
	RoleID uint
	Role   Role `gorm:"foreignKey:RoleID"`
}
