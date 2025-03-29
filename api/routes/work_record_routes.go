package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/attendance-system.git/controllers"
	"github.com/t2469/attendance-system.git/middleware"
)

func addWorkRecordRoutes(router *gin.Engine) {
	workRecords := router.Group("/work_records", middleware.AuthMiddleware())
	{
		workRecords.GET("", controllers.GetWorkRecords)
	}
}
