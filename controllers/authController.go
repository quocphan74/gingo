package controllers

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/quocphan74/gingo.git/common"
	"github.com/quocphan74/gingo.git/database"
	"github.com/quocphan74/gingo.git/models"
	"github.com/quocphan74/gingo.git/utils"
)

func Register(c *gin.Context) {

	var dataUser models.User
	en_password := c.PostForm("en_password")
	if err := c.Bind(&dataUser); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]utils.ApiError, len(ve))
			for i, fe := range ve {
				out[i] = utils.ApiError{fe.Field(), utils.MsgForTag(fe.Tag())}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": out})
		}
		return
	}

	if dataUser.Password != en_password {
		c.JSON(http.StatusFound, gin.H{"error": "Passwords do not match."})
		return
	}

	path, _ := common.UploadFile(c)
	if database.DB.First(&dataUser, "email=?", dataUser.Email); dataUser.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already taken"})
		return
	}
	user := models.User{
		FirstName: dataUser.FirstName,
		LastName:  dataUser.LastName,
		Phone:     dataUser.Phone,
		Email:     dataUser.Email,
		Avatar:    path,
		RoleUser: models.RoleUser{
			UserID: dataUser.ID,
			RoleID: 2,
		},
	}
	user.SetPassword(dataUser.Password)

	transaction := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			transaction.Rollback()
		}
	}()

	if err := transaction.Preload("RoleUser.Role").Create(&user); err.Error != nil {
		transaction.Rollback()
		log.Println(err)
	}

	transaction.Commit()

	res := models.UserStructure{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		RoleID:    user.RoleUser.RoleID,
		RoleName:  user.RoleUser.Role.Name,
	}

	c.Status(200)
	c.JSON(http.StatusOK, gin.H{
		"user":    res,
		"message": "Account created successfully",
	})

}

func Login(c *gin.Context) {
	var data map[string]string
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	var user models.User
	database.DB.Preload("RoleUser.Role").Where("email=?", data["email"]).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email Address does'n exit.",
		})
		return
	}

	if err := user.ComparePassword(data["password"]); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorret password.",
		})
		return
	}

	token, err := utils.GenerateToken(strconv.Itoa(int(user.ID)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error.",
		})
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	expirationSeconds := int(expirationTime.Unix())
	c.SetCookie("jwt", token, expirationSeconds, "/", "localhost", false, true)
	fmt.Println(user.RoleUser.Role)
	res := models.UserStructure{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		RoleID:    user.RoleUser.RoleID,
		RoleName:  user.RoleUser.Role.Name,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successfully.",
		"user":    res,
	})
	return
}

func ChangePassword(c *gin.Context) {

	cookie, _ := c.Cookie("jwt")
	var user models.User
	userID, _ := utils.Parsejwt(cookie)
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.Status(400)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Record not found!",
		})
		return
	}
	if err := c.ShouldBind(&user); err != nil {
		c.Status(400)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Record not found!",
		})
		return
	}
	user.SetPassword(user.Password)

	transition := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			transition.Rollback()
		}
	}()

	if err := transition.Model(&user).Updates(user); err.Error != nil {
		transition.Rollback()
		c.JSON(500, gin.H{"error": "Failed to update"})
		return
	}

	transition.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Updated password successfully.",
		"data":    user,
	})
	return

}

func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(1000000)
	codeStr := strconv.Itoa(code)
	paddedCode := fmt.Sprintf("%06s", codeStr)
	return paddedCode
}

func ResetPass(c *gin.Context) {

	code := c.Query("code")
	var user models.User
	var codeN models.Code

	database.DB.Where("code=?", code).First(&codeN)

	if codeN.ID == 0 {
		c.JSON(http.StatusFound, gin.H{
			"message": "Record not found.",
		})
		return
	}
	database.DB.Where("id=?", codeN.UserID).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusFound, gin.H{
			"message": "Record not found.",
		})
		return
	}
	pass := generateVerificationCode()
	user.SetPassword(pass)
	trnsition := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			trnsition.Rollback()
		}
	}()

	if err := trnsition.Model(&user).Updates(user); err.Error != nil {

		trnsition.Rollback()
		c.JSON(http.StatusFound, gin.H{
			"message": "Reset password failed.",
		})
		return
	}

	if err := trnsition.Delete(&codeN); err.Error != nil {

		trnsition.Rollback()
		c.JSON(http.StatusFound, gin.H{
			"message": "Deleted code failed.",
		})
		return
	}

	trnsition.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message":      "Change password successfully.",
		"new password": pass,
	})
	return
}

func CheckEmail(c *gin.Context) {

	var user models.User

	email := c.Query("email")

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	database.DB.Where("email=?", email).First(&user)
	if user.ID == 0 {
		c.Status(http.StatusFound)
		c.JSON(http.StatusFound, gin.H{
			"message": "Email Address does'n exit.",
		})
	}
	code := generateVerificationCode()

	codeN := models.Code{
		Code:   code,
		UserID: int(user.ID),
		// User:   user,
	}
	transition := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			transition.Rollback()
		}
	}()
	if errdb := transition.Create(&codeN); errdb.Error != nil {
		transition.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error.",
		})
		return
	}
	transition.Commit()
	from := os.Getenv("FROMEMAIL")
	password := os.Getenv("PASSWORDEMAIL")
	to := c.PostForm("email")
	subject := "Verification Code"
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Email Content</title>
		</head>
		<body>
			<h1>Change Password Request</h1>
			<p>Hello,</p>
			<p>We have received a request to change your password. If you did not make this request, please ignore this email.</p>
			<p>Your code is <strong>%s</strong></p>
		</body>
		</html>
	`, code)
	_smtp := os.Getenv("SMTP")
	port := os.Getenv("PORTSMTP")
	auth := smtp.PlainAuth("", from, password, _smtp)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=utf-8\r\n\r\n%s", from, to, subject, body)
	err := smtp.SendMail(_smtp+":"+port, auth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Println("Error sending email:", err)
		c.JSON(500, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(200, gin.H{"message": "Email sent"})
}

func GetRole(c *gin.Context) {
	var role []models.RoleUser

	database.DB.Preload("Role").Find(&role)

	c.Status(200)
	c.JSON(http.StatusOK, gin.H{

		"data":    role,
		"message": "Get All user Successfully",
	})
	return

}
