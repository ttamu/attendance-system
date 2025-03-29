package controllers

import (
	"errors"
	"github.com/t2469/attendance-system.git/services"
	"net/http"
	"strconv"
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

	day := eventTime.In(time.Local).Truncate(24 * time.Hour)
	if err := services.UpsertWorkRecord(input.EmployeeID, day); err != nil {
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

	y := c.Query("year")
	m := c.Query("month")

	var from, to time.Time
	useFilter := false

	if y != "" && m != "" {
		year, err1 := strconv.Atoi(y)
		month, err2 := strconv.Atoi(m)

		if err1 != nil || err2 != nil || month < 1 || month > 12 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year or month"})
			return
		}

		from = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		to = from.AddDate(0, 1, 0)
		useFilter = true
	}

	var timeClocks []models.TimeClock

	q := db.DB.
		Joins("JOIN employees ON employees.id = time_clocks.employee_id").
		Where("employees.company_id = ?", companyID)

	if useFilter {
		q = q.Where("time_clocks.timestamp >= ? AND time_clocks.timestamp < ?", from, to)
	}

	if err := q.Order("time_clocks.timestamp DESC").Find(&timeClocks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timeClocks)
}
