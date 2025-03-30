package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/attendance-system.git/controllers"
	"github.com/t2469/attendance-system.git/middleware"
)

func addClockRequestRoutes(router *gin.Engine) {
	requests := router.Group("/clock_requests", middleware.AuthMiddleware())
	{
		requests.GET("", controllers.GetClockRequests)
		requests.POST("/:id/approve", controllers.ApproveClockRequest)
	}
}
