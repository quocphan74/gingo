package controllers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/quocphan74/gingo.git/database"
	"github.com/quocphan74/gingo.git/models"
	"github.com/quocphan74/gingo.git/utils"
)

func Register(c *gin.Context) {

	var dataUser models.User

	if err := c.ShouldBind(&dataUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// upload file

	file, _ := c.FormFile("file")
	filename := filepath.Join("./uploads", file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
		return
	}

	user := models.User{
		FirstName: dataUser.FirstName,
		LastName:  dataUser.LastName,
		Phone:     dataUser.Phone,
		Email:     dataUser.Email,
		Avatar:    "/upload/" + file.Filename,
	}
	user.SetPassword(dataUser.Password)
	err := database.DB.Create(&user)
	if err != nil {
		log.Println(err)
	}
	res := models.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.FirstName,
		Email:     user.Email,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
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
	database.DB.Where("email=?", data["email"]).First(&user)
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

	res := models.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.FirstName,
		Email:     user.Email,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
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
	database.DB.Model(&user).Updates(user)
	res := models.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.FirstName,
		Email:     user.Email,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Updated password successfully.",
		"data":    res,
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

	if err := database.DB.Model(&user).Updates(user); err.Error != nil {
		c.JSON(http.StatusFound, gin.H{
			"message": "Reset password failed.",
		})
		return
	}

	if err := database.DB.Delete(&codeN); err.Error != nil {
		c.JSON(http.StatusFound, gin.H{
			"message": "Deleted code failed.",
		})
		return
	}
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

	if errdb := database.DB.Create(&codeN); errdb.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error.",
		})
		return
	}
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
