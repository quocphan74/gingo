package controllers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quocphan74/gingo.git/database"
	"github.com/quocphan74/gingo.git/models"
)

func GetAllUser(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "4"))
	offset := (page - 1) * pageSize
	var total int64
	var user []models.User
	database.DB.Preload("Blog").Offset(offset).Limit(pageSize).Find(&user)
	database.DB.Model(&models.User{}).Count(&total)
	c.Status(200)
	c.JSON(http.StatusOK, gin.H{
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"last_page": int(math.Ceil(float64(total) / float64(pageSize))),
		},
		"data":    user,
		"message": "Get All user Successfully",
	})
	return
}

func GetUser(c *gin.Context) {
	var user models.User
	userID := c.Param("id")
	database.DB.Where("id=?", userID).First(&user)
	if user.ID == 0 {
		c.Status(http.StatusFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "Successfully",
		"data":    user,
	})
	return

}

func DeleteUser(c *gin.Context) {
	var user models.User
	var users []models.User
	userID, _ := strconv.Atoi(c.Param("id"))
	if err := database.DB.Where("id=?", userID).First(&user).Error; err != nil {
		c.Status(400)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Record not found!",
		})
		return
	}
	transition := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			transition.Rollback()
		}
	}()

	if err := transition.Delete(&user); err != nil {
		transition.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Deleted to update"})
	}
	database.DB.Find(&users)

	transition.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Deleted successfully.",
		"data":    users,
	})
	return
}

func UpdateUser(c *gin.Context) {
	var user models.User
	userID, _ := strconv.Atoi(c.Param("id"))
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
	transition := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			transition.Rollback()
		}
	}()

	if err := transition.Model(&user).Updates(user); err.Error != nil {
		transition.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update!",
		})
	}

	transition.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Updated successfully.",
		"data":    user,
	})
	return
}
