package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/t2469/labor-management-system.git/config"
	"github.com/t2469/labor-management-system.git/controllers"
	"net/http"
	"time"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()
	setCors(router, cfg)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from Gin!"})
	})

	router.GET("/users", controllers.GetUsers)
	router.POST("/users", controllers.CreateUser)
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
