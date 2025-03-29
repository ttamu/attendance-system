package helpers

import (
	"errors"
	"github.com/t2469/attendance-system.git/db"
	"github.com/t2469/attendance-system.git/models"
)

// CheckEmployeeAccess 指定された従業員が会社に属しているかを確認するためのメソッド
func CheckEmployeeAccess(employeeID uint, companyID uint) error {
	var employee models.Employee
	if err := db.DB.First(&employee, employeeID).Error; err != nil {
		return errors.New("employee not found")
	}

	if employee.CompanyID != companyID {
		return errors.New("you are not allowed to access this employee")
	}

	return nil
}
