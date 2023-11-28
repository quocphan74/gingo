package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `form:"first_name" json:"first_name" binding:"required"`
	LastName  string `form:"last_name" json:"last_name" binding:"required"`
	Email     string `form:"email" json:"email" binding:"required"`
	Phone     string `form:"phone" json:"phone" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	Avatar    string
	Blog      []Blog
}

type UserResponse struct {
	ID        uint   `json:"id"`
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
	fmt.Println(user.Password)
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

}
