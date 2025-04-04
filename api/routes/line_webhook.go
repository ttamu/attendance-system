package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/attendance-system.git/config"
	"github.com/t2469/attendance-system.git/controllers"
)

func addLineWebhookRoutes(router *gin.Engine, cfg *config.Config) {
	router.POST("/webhook/line", controllers.HandleLineWebhook(cfg.LineChannelSecret))
}
