package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/t2469/labor-management-system.git/db"
	"github.com/t2469/labor-management-system.git/helpers"
	"github.com/t2469/labor-management-system.git/models"
	"net/http"
)

func CreateEmployeeAllowance(c *gin.Context) {
	var ea models.EmployeeAllowance
	if err := c.ShouldBindJSON(&ea); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyID, err := helpers.GetCompanyID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var emp models.Employee
	if err := db.DB.First(&emp, ea.EmployeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}
	if emp.CompanyID != companyID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed to create allowance for this employee"})
		return
	}

	if err := db.DB.Create(&ea).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ea)
}

func GetEmployeeAllowances(c *gin.Context) {
	companyID, err := helpers.GetCompanyID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var eas []models.EmployeeAllowance
	if err := db.DB.
		Joins("JOIN employees ON employee_allowances.employee_id = employees.id").
		Where("employees.company_id = ?", companyID).
		Find(&eas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, eas)
}

func GetEmployeeAllowance(c *gin.Context) {
	id := c.Param("id")
	companyID, err := helpers.GetCompanyID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var ea models.EmployeeAllowance
	if err := db.DB.First(&ea, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee allowance not found"})
		return
	}

	var emp models.Employee
	if err := db.DB.First(&emp, ea.EmployeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}
	if emp.CompanyID != companyID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed to access this allowance"})
		return
	}

	c.JSON(http.StatusOK, ea)
}
