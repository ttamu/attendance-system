package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/t2469/attendance-system.git/db"
	"github.com/t2469/attendance-system.git/models"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type CreateTimeClockInput struct {
	EmployeeID uint                 `json:"employee_id" binding:"required"`
	Type       models.TimeClockType `json:"type" binding:"required"`
	Timestamp  *time.Time           `json:"timestamp"`
}

func CreateTimeClock(c *gin.Context) {
	var input CreateTimeClockInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Timestampが未指定の場合、現在時刻を設定
	eventTime := time.Now()
	if input.Timestamp != nil {
		eventTime = *input.Timestamp
	}

	timeClock := models.TimeClock{
		EmployeeID: input.EmployeeID,
		Type:       input.Type,
		Timestamp:  eventTime,
	}

	if err := db.DB.Create(&timeClock).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, timeClock)
}

func GetTimeClock(c *gin.Context) {
	id := c.Param("id")

	var timeClock models.TimeClock
	if err := db.DB.First(&timeClock, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "time clock record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timeClock)
}
