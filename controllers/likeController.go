package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quocphan74/gingo.git/database"
	"github.com/quocphan74/gingo.git/models"
)

func Like(c *gin.Context) {
	var like models.Like

	if err := c.ShouldBindJSON(&like); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error",
		})
		return
	}

	fmt.Println(like)

	transition := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			transition.Rollback()
		}
	}()

	if err := transition.Create(&like).Error; err != nil {
		transition.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Like error.",
		})
		return
	}

	transition.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Like successfully.",
	})
	return
}

func UnLike(c *gin.Context) {
	var like models.Like

	likeID := c.Param("id")

	transition := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			transition.Rollback()
		}
	}()

	if err := transition.First(&like, "id=?", likeID).Delete(&like).Error; err != nil {
		transition.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "UnLike not success.",
		})
		return
	}

	transition.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Unlike successfully.",
	})
	return
}
