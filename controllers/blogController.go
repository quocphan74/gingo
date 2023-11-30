package controllers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quocphan74/gingo.git/database"
	"github.com/quocphan74/gingo.git/models"
	"github.com/quocphan74/gingo.git/utils"
)

func CreatePost(c *gin.Context) {
	var blog models.Blog

	if err := c.ShouldBind(&blog); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	transition := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			transition.Rollback()
		}
	}()

	if err := transition.Create(&blog); err.Error != nil {
		transition.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Created post failed.",
		})
		return
	}

	transition.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Created post successfully.",
		"data":    blog,
	})
	return

}

func GetAllPost(c *gin.Context) {
	var blogs []models.Blog
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "4"))
	offset := (page - 1) * pageSize
	var total int64
	database.DB.Preload("User", "id <> ?", 0).Preload("Comment").Preload("Like").Offset(offset).Limit(pageSize).Find(&blogs)
	database.DB.Model(&models.Blog{}).Count(&total)
	c.Status(200)
	c.JSON(http.StatusOK, gin.H{
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"last_page": int(math.Ceil(float64(total) / float64(pageSize))),
		},
		"data":    blogs,
		"message": "Get All post Successfully",
	})
	return
}

func GetPost(c *gin.Context) {
	var blog models.Blog
	postID := c.Param("id")
	database.DB.First(&blog, "id=?", postID)
	if blog.ID == 0 {
		c.JSON(http.StatusFound, gin.H{
			"message": "Record not found.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get post successfully.",
		"data":    blog,
	})
	return
}

func UpdatePost(c *gin.Context) {
	var blog models.Blog
	postID := c.Param("id")

	if err := database.DB.First(&blog, "id=?", postID); err.Error != nil {
		c.JSON(http.StatusFound, gin.H{
			"message": "Record not found.",
		})
		return
	}

	if err := c.ShouldBind(&blog); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	transition := database.DB.Begin()

	if err := transition.Model(&blog).Updates(blog); err.Error != nil {
		transition.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	transition.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Updated data successfully.",
		"data":    blog,
	})
	return
}

func DeletePost(c *gin.Context) {
	var blog models.Blog

	postID := c.Param("id")

	if err := database.DB.First(&blog, "id=?", postID); err.Error != nil {
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
	if err := transition.Delete(&blog); err.Error != nil {
		transition.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error,
		})
		return
	}

	transition.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Deleted successfully.",
	})
	return
}

func UniquePost(c *gin.Context) {

	cookie, _ := c.Cookie("jwt")
	id, _ := utils.Parsejwt(cookie)
	var blog []models.Blog
	database.DB.Preload("User", "id <> ?", 0).Preload("Comment").Model(&blog).Find(&blog, "user_id=?", id)

	c.JSON(http.StatusOK, gin.H{
		"data": blog,
	})
}
