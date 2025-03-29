package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/labor-management-system.git/controllers"
)

func addCompanyRoutes(router *gin.Engine) {
	companies := router.Group("/companies")
	{
		companies.POST("", controllers.CreateCompany)
	}
}
