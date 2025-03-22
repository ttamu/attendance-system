package services

import (
	"errors"
	"github.com/t2469/labor-management-system.git/models"
	"gorm.io/gorm"
	"sort"
)

type HealthInsuranceResponse struct {
	EmployeeName            string  `json:"employee_name"`
	CompanyName             string  `json:"company_name"`
	PrefectureName          string  `json:"prefecture_name"`
	Grade                   string  `json:"grade"`
	CalculatedMonthlyAmount int     `json:"calculated_monthly_amount"`
	HealthTotal             float64 `json:"health_total"`
	EmployeeHealth          float64 `json:"employee_health"`
	EmployerHealth          float64 `json:"employer_health"`
	WithCare                bool    `json:"with_care"`
	Age                     int     `json:"age"`
}

func CalculateInsurance(db *gorm.DB, employeeID uint) (HealthInsuranceResponse, error) {
	var employee models.Employee
	if err := db.
		Preload("Company.Prefecture.HealthInsuranceRates").
		First(&employee, employeeID).Error; err != nil {
		return HealthInsuranceResponse{}, err
	}

	age := calculateAge(employee.DateOfBirth)
	withCare := age >= 40 && age < 65

	pref := employee.Company.Prefecture
	if pref.ID == 0 {
		return HealthInsuranceResponse{}, errors.New("prefecture not found for employee's company")
	}

	rates := pref.HealthInsuranceRates
	if len(rates) == 0 {
		return HealthInsuranceResponse{}, errors.New("no health insurance rates configured for the prefecture")
	}

	sort.Slice(rates, func(i, j int) bool {
		return rates[i].MinMonthlyAmount < rates[j].MinMonthlyAmount
	})

	var selectedRate *models.HealthInsuranceRate
	for i, rate := range rates {
		if rate.MinMonthlyAmount <= employee.MonthlySalary && employee.MonthlySalary <= rate.MaxMonthlyAmount {
			selectedRate = &rates[i]
			break
		}
	}

	if selectedRate == nil {
		return HealthInsuranceResponse{}, errors.New("no matching rates found for employee's company")
	}

	var totalHealth, employeeHealth float64
	if withCare {
		totalHealth = selectedRate.HealthTotalWithCare
		employeeHealth = selectedRate.HealthHalfWithCare
	} else {
		totalHealth = selectedRate.HealthTotalNonCare
		employeeHealth = selectedRate.HealthHalfNonCare
	}
	employerHealth := totalHealth - employeeHealth

	resp := HealthInsuranceResponse{
		EmployeeName:            employee.Name,
		CompanyName:             employee.Company.Name,
		PrefectureName:          pref.Name,
		Grade:                   selectedRate.Grade,
		CalculatedMonthlyAmount: selectedRate.MonthlyAmount,
		HealthTotal:             totalHealth,
		EmployeeHealth:          employeeHealth,
		EmployerHealth:          employerHealth,
		Age:                     age,
		WithCare:                withCare,
	}
	return resp, nil
}
