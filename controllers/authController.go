package controllers

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
		Id:        user.Id,
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
	if user.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email Address does'n exit.",
		})
		return
	}

	// if err := user.ComparePassword(data["password"]); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"message": "Incorret password.",
	// 	})
	// 	return
	// }

	res := models.UserResponse{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.FirstName,
		Email:     user.Email,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
	}

	token, err := utils.GenerateToken(strconv.Itoa(int(user.Id)))
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
