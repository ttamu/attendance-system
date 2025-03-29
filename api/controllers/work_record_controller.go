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

	var records []models.WorkRecord
	if err := db.DB.
		Where("employee_id = ?", empID).
		Order("date ASC").
		Find(&records).Error; err != nil {
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
