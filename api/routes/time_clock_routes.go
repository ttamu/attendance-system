package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/attendance-system.git/controllers"
	"github.com/t2469/attendance-system.git/middleware"
)

func addTimeClockRoutes(router *gin.Engine) {
	timeClocks := router.Group("/time_clocks", middleware.AuthMiddleware())
	{
		timeClocks.POST("", controllers.CreateTimeClock)
		timeClocks.GET("", controllers.GetTimeClocks)
		timeClocks.GET("/:id", controllers.GetTimeClock)

		timeClocks.POST("/:id/requests", controllers.CreateClockRequest)
	}
}
