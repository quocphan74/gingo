package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `json:"first_name" form:"first_name" gorm:"type:varchar(255);not null"`
	LastName  string `json:"last_name" form:"last_name"  gorm:"type:varchar(255);not null"`
	Email     string `json:"email" form:"email" gorm:"type:varchar(255);unique;not null"`
	Phone     string `json:"phone" form:"phone" gorm:"type:varchar(255);unique;not null"`
	Password  string `json:"password" form:"password" gorm:"type:varchar(255);not null"`
	Avatar    string `gorm:"type:varchar(255);null"`
	Blog      []Blog
}

type UserResponse struct {
	ID        uint
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
