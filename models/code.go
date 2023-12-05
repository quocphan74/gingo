package models

import (
	"gorm.io/gorm"
)

type Code struct {
	gorm.Model
	Code   string
	UserID int
}
