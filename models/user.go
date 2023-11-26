package models

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Model struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	gorm.Model
	Id        uint   `form:"id"`
	FirstName string `form:"first_name" binding:"required"`
	LastName  string `form:"last_name" binding:"required"`
	Email     string `form:"email" binding:"required"`
	Phone     string `form:"phone" binding:"required"`
	Password  string `form:"-"`
	Avatar    string
}

type UserResponse struct {
	Id        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Avatar    string `json:"avatar"`
}

func (user *User) SetPassword(password string) {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user.Password = string(hashPassword)
}

func (user *User) ComparePassword(password string) error {
	match := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	fmt.Println("Match:   ", match)
	return match

}
