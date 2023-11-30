package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quocphan74/gingo.git/database"
	"github.com/quocphan74/gingo.git/models"
	"gorm.io/gorm/clause"
)

func CreateComment(c *gin.Context) {
	var comment models.Comment

	if err := c.ShouldBindJSON(&comment); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error Server.",
		})
		return
	}

	transition := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			transition.Rollback()
		}
	}()

	if err := transition.Create(&comment); err.Error != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Created not success.",
		})
		return
	}

	transition.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Created comment successfully.",
		"data":    comment,
	})
	return
}

func UpdateComment(c *gin.Context) {
	var comment models.Comment

	commnetID := c.Param("id")

	if err := database.DB.First(&comment, "id=?", commnetID); err.Error != nil {
		c.JSON(http.StatusFound, gin.H{
			"message": "Record not found.",
		})
		return
	}

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error Server.",
		})
		return
	}

	transition := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			transition.Rollback()
		}
	}()

	if err := transition.Model(&comment).Updates(comment); err.Error != nil {
		transition.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Updated not success.",
		})
		return
	}
	transition.Commit()

	c.JSON(http.StatusOK, gin.H{
		"Message": "Updated comment success",
		"data":    comment,
	})
	return

}

func DeleteComment(c *gin.Context) {
	var comment models.Comment

	commentID := c.Param("id")

	if err := database.DB.First(&comment, "id=?", commentID); err.Error != nil {
		c.JSON(http.StatusFound, gin.H{
			"message": "Record not found.",
		})
		return
	}

	transition := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			transition.Rollback()
		}
	}()

	if err := transition.Delete(&comment); err.Error != nil {
		transition.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Deleted comment not success.",
		})
	}

	transition.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Deleted successfully.",
	})
	return
}

func ReplyComment(c *gin.Context) {
	var comment models.Comment

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error Server.",
		})
		return
	}

	transition := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			transition.Rollback()
		}
	}()

	if err := transition.Create(&comment); err.Error != nil {
		transition.Rollback()
		c.JSON(http.StatusFound, gin.H{
			"message": "Comment not found.",
		})
		return
	}
	transition.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Reply comment successfully.",
	})
	return
}

func GetComment(c *gin.Context) {
	var comments []models.Comment
	commnentID := c.Param("id")
	if err := database.DB.Preload(clause.Associations).Find(&comments, commnentID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get comment reply successfully.",
		"data":    comments,
		"id":      commnentID,
	})
	return
}
