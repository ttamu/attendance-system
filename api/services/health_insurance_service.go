package services

import (
	"errors"
	"github.com/t2469/attendance-system.git/models"
	"gorm.io/gorm"
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

// CalculateInsurance 指定された年・月をもとに健康保険料を計算する関数
func CalculateInsurance(db *gorm.DB, employeeID uint, year, month int) (HealthInsuranceResponse, error) {
	var employee models.Employee
	if err := db.Preload("Company").First(&employee, employeeID).Error; err != nil {
		return HealthInsuranceResponse{}, err
	}

	age := calculateAge(employee.DateOfBirth)
	withCare := age >= 40 && age < 65

	prefectureID := employee.Company.PrefectureID
	if prefectureID == 0 {
		return HealthInsuranceResponse{}, errors.New("prefecture ID not set for employee's company")
	}

	var rate models.HealthInsuranceRate
	err := db.Where(
		"prefecture_id = ? AND min_monthly_amount <= ? AND max_monthly_amount >= ? "+
			"AND ((? > from_year) OR (? = from_year AND ? >= from_month)) "+
			"AND ((? < to_year) OR (? = to_year AND ? <= to_month))",
		prefectureID, employee.MonthlySalary, employee.MonthlySalary,
		year, year, month,
		year, year, month,
	).Order("from_year desc, from_month desc").
		First(&rate).Error
	if err != nil {
		return HealthInsuranceResponse{}, errors.New("no matching rate found for employee's company for the specified calculation date")
	}

	var pref models.Prefecture
	if err := db.First(&pref, prefectureID).Error; err != nil {
		return HealthInsuranceResponse{}, errors.New("prefecture not found")
	}

	var totalHealth, employeeHealth float64
	if withCare {
		totalHealth = rate.HealthTotalWithCare
		employeeHealth = rate.HealthHalfWithCare
	} else {
		totalHealth = rate.HealthTotalNonCare
		employeeHealth = rate.HealthHalfNonCare
	}
	employerHealth := totalHealth - employeeHealth

	resp := HealthInsuranceResponse{
		EmployeeName:            employee.Name,
		CompanyName:             employee.Company.Name,
		PrefectureName:          pref.Name,
		Grade:                   rate.Grade,
		CalculatedMonthlyAmount: rate.MonthlyAmount,
		HealthTotal:             totalHealth,
		EmployeeHealth:          employeeHealth,
		EmployerHealth:          employerHealth,
		Age:                     age,
		WithCare:                withCare,
	}
	return resp, nil
}
