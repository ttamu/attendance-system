package main

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/labor-management-system.git/models"
	"log"
	"net/http"
)

func main() {
	models.InitDB()
	if err := models.DB.AutoMigrate(&models.User{}, &models.Attendance{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from Gin!"})
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
