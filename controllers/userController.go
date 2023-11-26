package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quocphan74/gingo.git/database"
	"github.com/quocphan74/gingo.git/models"
)

func GetAllUser(c *gin.Context) {
	var user []models.User
	data := database.DB.Find(&user)
	c.Status(200)
	c.JSON(http.StatusOK, gin.H{
		"data":    data,
		"message": "Get All user Successfully",
	})
	return
}

func GetUser(c *gin.Context) {
	var user models.User
	userID := c.Param("id")
	database.DB.Where("id=?", userID).First(&user)
	if user.Id == 0 {
		c.Status(http.StatusFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "Successfully",
		"data":    user,
	})
	return

}
