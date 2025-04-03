package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/t2469/attendance-system.git/config"
	"github.com/t2469/attendance-system.git/controllers"
)

func addLineWebhookRoutes(router *gin.Engine, cfg *config.Config, bot *messaging_api.MessagingApiAPI) {
	router.POST("/webhook/line", controllers.HandleLineWebhook(cfg.LineChannelSecret, bot))
}
