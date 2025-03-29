package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/t2469/attendance-system.git/db"
	"github.com/t2469/attendance-system.git/helpers"
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

	companyID, err := helpers.GetCompanyID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var employee models.Employee
	if err := db.DB.First(&employee, input.EmployeeID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee not found"})
		return
	}

	if employee.CompanyID != companyID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not allowed to record time for this employee"})
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

	companyID, err := helpers.GetCompanyID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var emp models.Employee
	if err := db.DB.First(&emp, timeClock.EmployeeID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve employee"})
		return
	}

	if emp.CompanyID != companyID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access to this record is forbidden"})
		return
	}

	c.JSON(http.StatusOK, timeClock)
}
