package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/labor-management-system.git/db"
	"github.com/t2469/labor-management-system.git/models"
	"github.com/t2469/labor-management-system.git/services"
	"net/http"
	"strconv"
)

func GetEmployees(c *gin.Context) {
	var employee []models.Employee
	if err := db.DB.Find(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, employee)
}

func GetEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee
	if err := db.DB.Preload("Attendances").First(&employee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employee)
}

func CreateEmployee(c *gin.Context) {
	var employee models.Employee

	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&employee).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, employee)
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

	//year, err := strconv.Atoi(yearParam)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
	//}
	//month, err := strconv.Atoi(monthParam)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month"})
	//}

	resp, err := services.CalculatePayroll(db.DB, uint(employeeID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
