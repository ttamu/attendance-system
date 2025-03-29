package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/t2469/attendance-system.git/db"
	"github.com/t2469/attendance-system.git/helpers"
	"github.com/t2469/attendance-system.git/models"
	"gorm.io/gorm"
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

	companyID, err := helpers.GetCompanyID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := helpers.CheckEmployeeAccess(input.EmployeeID, companyID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

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

	companyID, err := helpers.GetCompanyID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := helpers.CheckEmployeeAccess(timeClock.EmployeeID, companyID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timeClock)
}

func GetTimeClocks(c *gin.Context) {
	companyID, err := helpers.GetCompanyID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var timeClocks []models.TimeClock
	if err := db.DB.Joins("JOIN employees ON employees.id = time_clocks.employee_id").Where("employees.company_id = ?", companyID).Order("time_clocks.timestamp DESC").Find(&timeClocks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timeClocks)
}
