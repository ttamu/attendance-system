package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/attendance-system.git/db"
	"github.com/t2469/attendance-system.git/models"
	"net/http"
	"strconv"
	"time"
)

type WorkRecordResponse struct {
	ID           uint      `json:"id"`
	EmployeeID   uint      `json:"employee_id"`
	Date         time.Time `json:"date"`
	ClockIn      time.Time `json:"clock_in"`
	ClockOut     time.Time `json:"clock_out"`
	BreakMinutes int64     `json:"break_minutes"`
	WorkMinutes  int64     `json:"work_minutes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func GetWorkRecords(c *gin.Context) {
	empIDStr := c.Query("employee_id")
	if empIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id is required"})
		return
	}

	empID, err := strconv.Atoi(empIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee_id"})
		return
	}

	yearStr := c.Query("year")
	monthStr := c.Query("month")

	var from, to time.Time
	useFilter := false

	if yearStr != "" && monthStr != "" {
		year, err1 := strconv.Atoi(yearStr)
		month, err2 := strconv.Atoi(monthStr)

		if err1 != nil || err2 != nil || month < 1 || month > 12 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year or month"})
			return
		}

		from = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
		to = from.AddDate(0, 1, 0)
		useFilter = true
	}

	var records []models.WorkRecord

	q := db.DB.Where("employee_id = ?", empID)
	if useFilter {
		q = q.Where("date >= ? AND date < ?", from, to)
	}

	if err := q.Order("date ASC").Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []WorkRecordResponse
	for _, r := range records {
		response = append(response, WorkRecordResponse{
			ID:           r.ID,
			EmployeeID:   r.EmployeeID,
			Date:         r.Date,
			ClockIn:      r.ClockIn,
			ClockOut:     r.ClockOut,
			BreakMinutes: r.BreakMinutes,
			WorkMinutes:  r.WorkMinutes,
			CreatedAt:    r.CreatedAt,
			UpdatedAt:    r.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}
