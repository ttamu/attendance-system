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
		allowanceTypes.GET("", controllers.GetAllowanceTypes)
		allowanceTypes.GET("/:id", controllers.GetAllowanceType)
		allowanceTypes.PUT("/:id", controllers.UpdateAllowanceType)
		allowanceTypes.DELETE("/:id", controllers.DeleteAllowanceType)
	}

	employeeAllowances := router.Group("/employee_allowances", middleware.AuthMiddleware())
	{
		employeeAllowances.POST("", controllers.CreateEmployeeAllowance)
		employeeAllowances.GET("", controllers.GetEmployeeAllowances)
		employeeAllowances.GET("/:id", controllers.GetEmployeeAllowance)
		employeeAllowances.PUT("/:id", controllers.UpdateEmployeeAllowance)
	}
}
