package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/labor-management-system.git/controllers"
)

func addEmployeeRoutes(router *gin.Engine) {
	employees := router.Group("/employees")
	{
		employees.GET("", controllers.GetEmployees)
		employees.GET("/:id", controllers.GetEmployee)
		employees.POST("", controllers.CreateEmployee)
		employees.POST("/:id/attendances", controllers.CreateAttendance)
		employees.GET("/:id/insurance", controllers.CalculateEmployeeInsurance)
		employees.GET("/:id/pension", controllers.CalculateEmployeePension)
		employees.GET("/:id/payroll", controllers.CalculateEmployeePayroll)
	}
}
