package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/labor-management-system.git/controllers"
)

func addAllowanceRoutes(router *gin.Engine) {
	allowances := router.Group("/allowances")
	{
		allowances.POST("type", controllers.CrateAllowanceType)
		allowances.POST("", controllers.CreateEmployeeAllowance)
	}
}
