package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/labor-management-system.git/controllers"
)

func addAuthRoutes(router *gin.Engine) {
	router.POST("/register", controllers.Register)
}
