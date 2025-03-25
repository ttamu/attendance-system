package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/labor-management-system.git/controllers"
	"github.com/t2469/labor-management-system.git/middleware"
)

func addAllowanceRoutes(router *gin.Engine) {
	allowanceTypes := router.Group("/allowance_types", middleware.AuthMiddleware())
	{
		allowanceTypes.POST("", controllers.CreateAllowanceType)
	}

	employeeAllowances := router.Group("/employee_allowances", middleware.AuthMiddleware())
	{
		employeeAllowances.POST("", controllers.CreateEmployeeAllowance)
	}
}
