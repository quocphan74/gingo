package models

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	UserID     uint   `json:"userId,omitempty" `
	TargetID   uint   `json:"targetId,omitempty" `
	TargetType string `json:"targetType,omitempty" `
}
