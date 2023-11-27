package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quocphan74/gingo.git/controllers"
	"github.com/quocphan74/gingo.git/middleware"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	router.POST("/api/v1/register", controllers.Register)
	router.POST("/api/v1/login", controllers.Login)

	v1 := router.Group("/api/v1/")
	v1.Use(middleware.AuthMiddleware())
	{
		v1.GET("user/:id", controllers.GetUser)
		v1.GET("users", controllers.GetAllUser)
		v1.DELETE("user/:id", controllers.DeleteUser)
		v1.PUT("user/:id", controllers.UpdateUser)

		v1.PUT("user/change-password", controllers.ChangePassword)

		v1.GET("send-mail", controllers.CheckEmail)

		v1.GET("rest-password", controllers.ResetPass)
		v1.GET("post/home", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "hello work",
			})
		})
	}
	return router
}
