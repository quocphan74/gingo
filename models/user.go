package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `json:"first_name" form:"first_name" gorm:"type:varchar(100);" binding:"required,min=3,max=40"`
	LastName  string `json:"last_name" form:"last_name" gorm:"type:varchar(100); not null;"  binding:"required,min=3,max=40"`
	Email     string `json:"email" form:"email" gorm:"type:varchar(100); not null;" binding:"required,min=3,max=100,email"`
	Phone     string `json:"phone" form:"phone" gorm:"type:varchar(11); not null;" binding:"required,min=3,max=11"`
	Password  string `json:"password" form:"password" gorm:"type:varchar(100); not null;" binding:"required,min=4,max=20"`
	Avatar    string `gorm:"type:varchar(255)"`
	RoleUser  RoleUser
}

type UserStructure struct {
	ID        uint
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Avatar    string
	RoleID    uint
	RoleName  string
}

func (user *User) SetPassword(password string) {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user.Password = string(hashPassword)
}

func (user *User) ComparePassword(password string) error {
	fmt.Println(user.Password)
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

}
