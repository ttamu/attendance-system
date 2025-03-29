package services

import (
	"errors"
	"github.com/t2469/attendance-system.git/models"
	"gorm.io/gorm"
	"sort"
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

func CalculatePension(db *gorm.DB, employeeID uint) (PensionInsuranceResponse, error) {
	var employee models.Employee
	if err := db.Preload("Company.Prefecture.PensionInsuranceRates").First(&employee, employeeID).Error; err != nil {
		return PensionInsuranceResponse{}, err
	}

	age := calculateAge(employee.DateOfBirth)

	pref := employee.Company.Prefecture
	if pref.ID == 0 {
		return PensionInsuranceResponse{}, errors.New("prefecture not found for employee's company")
	}

	rates := pref.PensionInsuranceRates
	if len(rates) == 0 {
		return PensionInsuranceResponse{}, errors.New("no pension insurance rates configured for the prefecture")
	}

	sort.Slice(rates, func(i, j int) bool {
		return rates[i].MinMonthlyAmount < rates[j].MinMonthlyAmount
	})

	var selectedRate *models.PensionInsuranceRate
	for i, rate := range rates {
		if employee.MonthlySalary >= rate.MinMonthlyAmount && employee.MonthlySalary <= rate.MaxMonthlyAmount {
			selectedRate = &rates[i]
			break
		}
	}
	if selectedRate == nil {
		return PensionInsuranceResponse{}, errors.New("no matching rates found for employee's company")
	}

	total := selectedRate.PensionTotal
	employeePension := selectedRate.PensionHalf
	employerPension := total - employeePension

	resp := PensionInsuranceResponse{
		EmployeeName:            employee.Name,
		CompanyName:             employee.Company.Name,
		PrefectureName:          pref.Name,
		Grade:                   selectedRate.Grade,
		CalculatedMonthlyAmount: selectedRate.MonthlyAmount,
		PensionTotal:            total,
		EmployeePension:         employeePension,
		EmployerPension:         employerPension,
		Age:                     age,
	}
	return resp, nil
}
