package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/attendance-system.git/db"
	"github.com/t2469/attendance-system.git/models"
	"net/http"
	"time"
)

func CreateAttendance(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee
	if err := db.DB.First(&employee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var attendance models.Attendance
	if err := c.ShouldBindJSON(&attendance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	attendance.EmployeeID = employee.ID
	if attendance.CheckIn.IsZero() {
		attendance.CheckIn = time.Now()
	}

	if err := db.DB.Create(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, attendance)
}
