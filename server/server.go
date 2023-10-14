package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ugokalp/controllers"
	"github.com/ugokalp/middleware"
	"gorm.io/gorm"
)

func InitServer(db *gorm.DB) {
	router := gin.Default()
	router.Use(middleware.InitCors())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.GET("/healthcheck", controllers.HealthCheck)

	v1 := router.Group("/v1")
	uc := controllers.UserController{
		Db: db,
	}
	v1.GET("username", uc.GetUserByUsername)

	userRouter := v1.Group("user")
	userRouter.POST("signup", uc.SignUp)
	userRouter.POST("login", uc.Login)
	userRouter.GET("me", middleware.Authenticated, uc.GetUser)
	userRouter.PUT("me", middleware.Authenticated, uc.UpdateUser)
	userRouter.DELETE("", uc.DeleteUser)

	router.Run()
}
