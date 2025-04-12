package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/t2469/attendance-system.git/db"
	"github.com/t2469/attendance-system.git/helpers"
	"github.com/t2469/attendance-system.git/models"
)

type EmployeeAllowanceResponse struct {
	ID                uint     `json:"id"`
	EmployeeID        uint     `json:"employee_id"`
	AllowanceTypeID   uint     `json:"allowance_type_id"`
	Amount            int      `json:"amount"`
	CommissionRate    *float64 `json:"commission_rate,omitempty"`
	Year              int      `json:"year"`
	Month             int      `json:"month"`
	EmployeeName      string   `json:"employee_name"`
	AllowanceTypeName string   `json:"allowance_type_name"`
	AllowanceType     string   `json:"allowance_type"` // "fixed" or "commission"
}

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

	var response EmployeeAllowanceResponse
	if err := db.DB.Table("employee_allowances").
		Select("employee_allowances.*, employees.name as employee_name, allowance_types.name as allowance_type_name, allowance_types.type as allowance_type").
		Joins("JOIN employees ON employee_allowances.employee_id = employees.id").
		Joins("JOIN allowance_types ON employee_allowances.allowance_type_id = allowance_types.id").
		Where("employee_allowances.id = ?", ea.ID).
		Scan(&response).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func GetEmployeeAllowances(c *gin.Context) {
	companyID, err := helpers.GetCompanyID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var responses []EmployeeAllowanceResponse
	if err := db.DB.Table("employee_allowances").
		Select("employee_allowances.*, employees.name as employee_name, allowance_types.name as allowance_type_name, allowance_types.type as allowance_type").
		Joins("JOIN employees ON employee_allowances.employee_id = employees.id").
		Joins("JOIN allowance_types ON employee_allowances.allowance_type_id = allowance_types.id").
		Where("employees.company_id = ?", companyID).
		Scan(&responses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if responses == nil {
		responses = []EmployeeAllowanceResponse{}
	}
	c.JSON(http.StatusOK, responses)
}

func GetEmployeeAllowance(c *gin.Context) {
	id := c.Param("id")
	companyID, err := helpers.GetCompanyID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var response EmployeeAllowanceResponse
	if err := db.DB.Table("employee_allowances").
		Select("employee_allowances.*, employees.name as employee_name, allowance_types.name as allowance_type_name, allowance_types.type as allowance_type").
		Joins("JOIN employees ON employee_allowances.employee_id = employees.id").
		Joins("JOIN allowance_types ON employee_allowances.allowance_type_id = allowance_types.id").
		Where("employee_allowances.id = ? AND employees.company_id = ?", id, companyID).
		Scan(&response).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee allowance not found"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func UpdateEmployeeAllowance(c *gin.Context) {
	id := c.Param("id")
	var ea models.EmployeeAllowance
	if err := db.DB.First(&ea, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee allowance not found"})
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
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed to update this allowance"})
		return
	}

	if err := c.ShouldBindJSON(&ea); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newEmp models.Employee
	if err := db.DB.First(&newEmp, ea.EmployeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}
	if newEmp.CompanyID != companyID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed to update allowance to an employee outside your company"})
		return
	}

	if err := db.DB.Save(&ea).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response EmployeeAllowanceResponse
	if err := db.DB.Table("employee_allowances").
		Select("employee_allowances.*, employees.name as employee_name, allowance_types.name as allowance_type_name, allowance_types.type as allowance_type").
		Joins("JOIN employees ON employee_allowances.employee_id = employees.id").
		Joins("JOIN allowance_types ON employee_allowances.allowance_type_id = allowance_types.id").
		Where("employee_allowances.id = ?", ea.ID).
		Scan(&response).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func DeleteEmployeeAllowance(c *gin.Context) {
	id := c.Param("id")
	var ea models.EmployeeAllowance

	if err := db.DB.First(&ea, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee allowance not found"})
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
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed to delete this allowance"})
		return
	}

	if err := db.DB.Delete(&ea).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee allowance deleted"})
}
