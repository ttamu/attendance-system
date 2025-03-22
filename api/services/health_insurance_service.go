package services

import (
	"errors"
	"github.com/t2469/labor-management-system.git/models"
	"gorm.io/gorm"
	"sort"
)

type InsuranceCalculationResponse struct {
	EmployeeName            string  `json:"employee_name"`
	CompanyName             string  `json:"company_name"`
	PrefectureName          string  `json:"prefecture_name"`
	Grade                   string  `json:"grade"`
	CalculatedMonthlyAmount int     `json:"calculated_monthly_amount"`
	TotalPremium            float64 `json:"total_premium"`
	EmployeePremium         float64 `json:"employee_premium"`
	EmployerPremium         float64 `json:"employer_premium"`
	Age                     int     `json:"age"`
	WithCare                bool    `json:"with_care"`
}

func CalculateInsurance(db *gorm.DB, employeeID uint) (InsuranceCalculationResponse, error) {
	var employee models.Employee
	if err := db.Preload("Company.Prefecture.HealthInsuranceRates").First(&employee, employeeID).Error; err != nil {
		return InsuranceCalculationResponse{}, err
	}

	age := calculateAge(employee.DateOfBirth)
	withCare := age >= 40 && age < 65

	pref := employee.Company.Prefecture
	if pref.ID == 0 {
		return InsuranceCalculationResponse{}, errors.New("prefecture not found for employee's company")
	}

	rates := pref.HealthInsuranceRates
	if len(rates) == 0 {
		return InsuranceCalculationResponse{}, errors.New("no health insurance rates configured for the prefecture")
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
		return InsuranceCalculationResponse{}, errors.New("no matching rates found for employee's company")
	}

	var totalPremium, employeePremium float64
	if withCare {
		totalPremium = selectedRate.HealthTotalWithCare
		employeePremium = selectedRate.HealthHalfWithCare
	} else {
		totalPremium = selectedRate.HealthTotalNonCare
		employeePremium = selectedRate.HealthHalfNonCare
	}
	employerPremium := totalPremium - employeePremium

	resp := InsuranceCalculationResponse{
		EmployeeName:            employee.Name,
		CompanyName:             employee.Company.Name,
		PrefectureName:          pref.Name,
		Grade:                   selectedRate.Grade,
		CalculatedMonthlyAmount: selectedRate.MonthlyAmount,
		TotalPremium:            totalPremium,
		EmployeePremium:         employeePremium,
		EmployerPremium:         employerPremium,
		Age:                     age,
		WithCare:                withCare,
	}
	return resp, nil
}
