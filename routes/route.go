package routes

import (
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
		// User
		v1.GET("user/:id", controllers.GetUser)
		v1.GET("users", controllers.GetAllUser)
		v1.DELETE("user/:id", controllers.DeleteUser)
		v1.PUT("user/:id", controllers.UpdateUser)
		v1.PUT("user/change-password", controllers.ChangePassword)
		v1.GET("send-mail", controllers.CheckEmail)
		v1.GET("rest-password", controllers.ResetPass)
		// Post
		v1.GET("posts", controllers.GetAllPost)
		v1.GET("post/:id", controllers.GetPost)
		v1.POST("post", controllers.CreatePost)
		v1.PUT("post/:id", controllers.UpdatePost)
		v1.DELETE("post/:id", controllers.DeletePost)
		v1.GET("unique-post", controllers.UniquePost)

		// Comment
		v1.POST("comment", controllers.CreateComment)
		v1.PUT("comment/:id", controllers.UpdateComment)
		v1.DELETE("comment/:id", controllers.DeleteComment)
	}
	return router
}
