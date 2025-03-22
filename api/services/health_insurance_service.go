package services

import (
	"errors"
	"github.com/t2469/labor-management-system.git/models"
	"gorm.io/gorm"
	"sort"
)

type InsuranceCalculationResponse struct {
	UserName                string  `json:"user_name"`
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

func CalculateInsuranceForUser(db *gorm.DB, userID uint) (InsuranceCalculationResponse, error) {
	var user models.User
	if err := db.Preload("Company.Prefecture.HealthInsuranceRates").First(&user, userID).Error; err != nil {
		return InsuranceCalculationResponse{}, err
	}

	age := calculateAge(user.DateOfBirth)
	withCare := 40 <= age && age < 65

	pref := user.Company.Prefecture
	if pref.ID == 0 {
		return InsuranceCalculationResponse{}, errors.New("prefecture not found for user's company")
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
		if rate.MinMonthlyAmount <= user.MonthlySalary && user.MonthlySalary <= rate.MaxMonthlyAmount {
			selectedRate = &rates[i]
			break
		}
	}

	if selectedRate == nil {
		return InsuranceCalculationResponse{}, errors.New("no matching rates found for user's company")
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
		UserName:                user.Name,
		CompanyName:             user.Company.Name,
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
