package services

import (
	"errors"
	"github.com/t2469/attendance-system.git/models"
	"gorm.io/gorm"
)

type PensionInsuranceResponse struct {
	EmployeeName            string  `json:"employee_name"`
	CompanyName             string  `json:"company_name"`
	PrefectureName          string  `json:"prefecture_name"`
	Grade                   string  `json:"grade"`
	CalculatedMonthlyAmount int     `json:"calculated_monthly_amount"`
	PensionTotal            float64 `json:"pension_total"`
	EmployeePension         float64 `json:"employee_pension"`
	EmployerPension         float64 `json:"employer_pension"`
	Age                     int     `json:"age"`
}

// CalculatePension は、指定された年・月をもとに年金保険料を計算する関数
func CalculatePension(db *gorm.DB, employeeID uint, calcYear, calcMonth int) (PensionInsuranceResponse, error) {
	var employee models.Employee
	if err := db.Preload("Company").First(&employee, employeeID).Error; err != nil {
		return PensionInsuranceResponse{}, err
	}

	age := calculateAge(employee.DateOfBirth)
	prefID := employee.Company.PrefectureID
	if prefID == 0 {
		return PensionInsuranceResponse{}, errors.New("prefecture ID not set for employee's company")
	}

	var rate models.PensionInsuranceRate
	err := db.Where(
		"prefecture_id = ? AND min_monthly_amount <= ? AND max_monthly_amount >= ? "+
			"AND ((from_year < ? OR (from_year = ? AND from_month <= ?)) "+
			"AND (to_year > ? OR (to_year = ? AND to_month >= ?)))",
		prefID, employee.MonthlySalary, employee.MonthlySalary,
		calcYear, calcYear, calcMonth,
		calcYear, calcYear, calcMonth,
	).Order("from_year desc, from_month desc").
		First(&rate).Error
	if err != nil {
		return PensionInsuranceResponse{}, errors.New("no matching rate found for employee's company for the specified calculation date")
	}

	var pref models.Prefecture
	if err := db.First(&pref, prefID).Error; err != nil {
		return PensionInsuranceResponse{}, errors.New("prefecture not found")
	}

	total := rate.PensionTotal
	employeePension := rate.PensionHalf
	employerPension := total - employeePension

	resp := PensionInsuranceResponse{
		EmployeeName:            employee.Name,
		CompanyName:             employee.Company.Name,
		PrefectureName:          pref.Name,
		Grade:                   rate.Grade,
		CalculatedMonthlyAmount: rate.MonthlyAmount,
		PensionTotal:            total,
		EmployeePension:         employeePension,
		EmployerPension:         employerPension,
		Age:                     age,
	}
	return resp, nil
}
