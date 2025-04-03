package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/t2469/attendance-system.git/config"
	"log"
	"net/http"
	"time"
)

func Run(cfg *config.Config, bot *messaging_api.MessagingApiAPI) {
	router := setupRouter(cfg, bot)
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func setupRouter(cfg *config.Config, bot *messaging_api.MessagingApiAPI) *gin.Engine {
	router := gin.Default()
	setCors(router, cfg)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from Gin!"})
	})

	addEmployeeRoutes(router)
	addCompanyRoutes(router)
	addAllowanceRoutes(router)
	addAuthRoutes(router)
	addTimeClockRoutes(router)
	addWorkRecordRoutes(router)
	addClockRequestRoutes(router)
	addLineWebhookRoutes(router, cfg, bot)

	return router
}

func setCors(router *gin.Engine, cfg *config.Config) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.AllowedOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}
