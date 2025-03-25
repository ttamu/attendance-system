package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/labor-management-system.git/controllers"
	"github.com/t2469/labor-management-system.git/middleware"
)

func addCompanyRoutes(router *gin.Engine) {
	companies := router.Group("/companies", middleware.AuthMiddleware())
	{
		companies.POST("", controllers.CreateCompany)
	}
}
