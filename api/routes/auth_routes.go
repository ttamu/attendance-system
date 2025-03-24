package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/labor-management-system.git/controllers"
	"github.com/t2469/labor-management-system.git/middleware"
)

func addAuthRoutes(router *gin.Engine) {
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)

	auth := router.Group("/", middleware.AuthMiddleware())
	{
		auth.GET("/current_account", controllers.CurrentAccount)
		auth.GET("/admin/employees", controllers.GetEmployees)
	}
}
