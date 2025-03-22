package services

import (
	"errors"
	"github.com/t2469/labor-management-system.git/models"
	"gorm.io/gorm"
	"sort"
)

// PensionCalculationResponse は、厚生年金保険料計算結果のレスポンス形式です。
type PensionCalculationResponse struct {
	UserName                string  `json:"user_name"`
	CompanyName             string  `json:"company_name"`
	PrefectureName          string  `json:"prefecture_name"`
	Grade                   string  `json:"grade"`
	CalculatedMonthlyAmount int     `json:"calculated_monthly_amount"`
	PensionTotal            float64 `json:"pension_total"`
	EmployeePension         float64 `json:"employee_pension"`
	EmployerPension         float64 `json:"employer_pension"`
	Age                     int     `json:"age"`
}

func CalculatePensionForUser(db *gorm.DB, userID uint) (PensionCalculationResponse, error) {
	var user models.User
	if err := db.Preload("Company.Prefecture.PensionInsuranceRates").First(&user, userID).Error; err != nil {
		return PensionCalculationResponse{}, err
	}

	age := calculateAge(user.DateOfBirth)

	pref := user.Company.Prefecture
	if pref.ID == 0 {
		return PensionCalculationResponse{}, errors.New("prefecture not found for user's company")
	}

	rates := pref.PensionInsuranceRates
	if len(rates) == 0 {
		return PensionCalculationResponse{}, errors.New("no pension insurance rates configured for the prefecture")
	}

	sort.Slice(rates, func(i, j int) bool {
		return rates[i].MinMonthlyAmount < rates[j].MinMonthlyAmount
	})

	var selectedRate *models.PensionInsuranceRate
	for i, rate := range rates {
		if user.MonthlySalary >= rate.MinMonthlyAmount && user.MonthlySalary <= rate.MaxMonthlyAmount {
			selectedRate = &rates[i]
			break
		}
	}
	if selectedRate == nil {
		return PensionCalculationResponse{}, errors.New("no matching rates found for user's company")
	}

	total := selectedRate.PensionTotal
	employee := selectedRate.PensionHalf
	employer := total - employee

	resp := PensionCalculationResponse{
		UserName:                user.Name,
		CompanyName:             user.Company.Name,
		PrefectureName:          pref.Name,
		Grade:                   selectedRate.Grade,
		CalculatedMonthlyAmount: selectedRate.MonthlyAmount,
		PensionTotal:            total,
		EmployeePension:         employee,
		EmployerPension:         employer,
		Age:                     age,
	}
	return resp, nil
}
