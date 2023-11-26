package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quocphan74/gingo.git/controllers"
	"github.com/quocphan74/gingo.git/middleware"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1/user")
	{
		v1.GET("/:id", controllers.GetUser)
		v1.GET("/", controllers.GetAllUser)
		v1.POST("/register", controllers.Register)
		v1.POST("/login", controllers.Login)
	}
	authorized := router.Group("/api/v2/post")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("/home", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "hello work",
			})
		})
	}
	return router
}
