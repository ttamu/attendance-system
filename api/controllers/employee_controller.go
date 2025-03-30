package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/attendance-system.git/db"
	"github.com/t2469/attendance-system.git/models"
	"github.com/t2469/attendance-system.git/services"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

func formatEmployee(emp models.Employee) gin.H {
	formattedTimeClocks := []gin.H{}
	for _, tc := range emp.TimeClocks {
		formattedTimeClocks = append(formattedTimeClocks, gin.H{
			"employee_id": tc.EmployeeID,
			"type":        tc.Type,
			"timestamp":   tc.Timestamp.In(time.Local).Format("2006/1/2 15:04:05"),
		})
	}

	return gin.H{
		"id":             emp.ID,
		"name":           emp.Name,
		"monthly_salary": emp.MonthlySalary,
		"date_of_birth":  emp.DateOfBirth.In(time.Local).Format("2006/1/2"),
		"time_clocks":    formattedTimeClocks,
	}
}

func GetEmployees(c *gin.Context) {
	companyID, exists := c.Get("company_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証情報がありません"})
		return
	}

	var employees []models.Employee
	if err := db.DB.Where("company_id = ?", companyID).Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employees)
}

func GetEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee
	if err := db.DB.Preload("TimeClocks", func(db *gorm.DB) *gorm.DB { return db.Order("timestamp ASC") }).First(&employee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, formatEmployee(employee))
}

func CreateEmployee(c *gin.Context) {
	var employee models.Employee

	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&employee).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, formatEmployee(employee))
}

func CalculateEmployeeInsurance(c *gin.Context) {
	id := c.Param("id")
	employeeID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	resp, err := services.CalculateInsurance(db.DB, uint(employeeID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, resp)
}

func CalculateEmployeePension(c *gin.Context) {
	id := c.Param("id")
	employeeID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	resp, err := services.CalculatePension(db.DB, uint(employeeID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func CalculateEmployeePayroll(c *gin.Context) {
	id := c.Param("id")
	employeeID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	yearParam := c.Query("year")
	monthParam := c.Query("month")
	if yearParam == "" || monthParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year and month parameters are required"})
		return
	}

	year, err := strconv.Atoi(yearParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
	}
	month, err := strconv.Atoi(monthParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month"})
	}

	resp, err := services.CalculatePayroll(db.DB, uint(employeeID), year, month)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
