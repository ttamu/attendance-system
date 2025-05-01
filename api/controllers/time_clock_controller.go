package controllers

import (
	"errors"
	"fmt"
	"github.com/t2469/attendance-system.git/services"
	"log"
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
	Notify     bool                 `json:"notify"` // 通知するかのフラグ
	NotifyAt   string               `json:"notify_at" binding:"omitempty"`
}

func formatTimeClock(tc models.TimeClock) gin.H {
	return gin.H{
		"id":          tc.ID,
		"employee_id": tc.EmployeeID,
		"type":        tc.Type,
		"timestamp":   tc.Timestamp.In(time.Local).Format("2006/1/2 15:04:05"),
	}
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

	timeClock, err := services.RecordTimeClock(input.EmployeeID, input.Type, eventTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if input.Notify && input.Type == models.ClockIn {
		go func(empId uint, notifyAt string, clockIn time.Time) {
			now := time.Now().In(time.Local)

			t, err := time.Parse("15:04", notifyAt)
			if err != nil {
				log.Println("invalid notify_at format:", err)
				return
			}
			notifyTime := time.Date(
				now.Year(), now.Month(), now.Day(),
				t.Hour(), t.Minute(), 0, 0,
				now.Location(),
			)

			if notifyTime.Before(now) {
				return
			}

			wait := time.Until(notifyTime)
			if wait > 0 {
				time.Sleep(wait)
			}

			var cnt int64
			db.DB.Model(&models.TimeClock{}).Where("employee_id = ? AND type = ? AND timestamp > ?", empId, models.ClockOut, clockIn).Count(&cnt)

			if cnt == 0 {
				var emp models.Employee
				if err := db.DB.First(&emp, empId).Error; err != nil {
					log.Println(err)
					return
				}
				if emp.LineUserID != nil {
					msg := fmt.Sprintf(
						"%sさん\n退勤を忘れていませんか？\n出勤時刻: %s\n現在時刻: %s",
						emp.Name,
						clockIn.In(time.Local).Format("15:04"),
						time.Now().In(time.Local).Format("15:04"),
					)
					if err := services.SendMessage(*emp.LineUserID, msg); err != nil {
						log.Println(err)
					}
				}
			}
		}(input.EmployeeID, input.NotifyAt, eventTime)
	}

	c.JSON(http.StatusCreated, formatTimeClock(timeClock))
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

	c.JSON(http.StatusOK, formatTimeClock(timeClock))
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

	if err := q.Order("time_clocks.timestamp ASC").Find(&timeClocks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resp []gin.H
	for _, tc := range timeClocks {
		resp = append(resp, formatTimeClock(tc))
	}
	c.JSON(http.StatusOK, resp)
}
