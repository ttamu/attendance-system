package services

import (
	"gorm.io/gorm"
)

type PayrollCalculationResponse struct {
	EmployeeName    string  `json:"employee_name"`
	GrossSalary     int     `json:"gross_salary"`
	HealthInsurance float64 `json:"health_insurance"`
	Pension         float64 `json:"pension"`
	TotalDeductions float64 `json:"total_deductions"`
	NetSalary       float64 `json:"net_salary"`
}

func CalculatePayroll(db *gorm.DB, employeeID uint) (PayrollCalculationResponse, error) {
	insuranceResp, err := CalculateInsurance(db, employeeID)
	if err != nil {
		return PayrollCalculationResponse{}, err
	}

	pensionResp, err := CalculatePension(db, employeeID)
	if err != nil {
		return PayrollCalculationResponse{}, err
	}

	grossSalary := insuranceResp.CalculatedMonthlyAmount

	totalDeductions := insuranceResp.EmployeePremium + pensionResp.EmployeePension

	netSalary := float64(grossSalary) - totalDeductions

	resp := PayrollCalculationResponse{
		EmployeeName:    insuranceResp.EmployeeName,
		GrossSalary:     grossSalary,
		HealthInsurance: insuranceResp.EmployeePremium,
		Pension:         pensionResp.EmployeePension,
		TotalDeductions: totalDeductions,
		NetSalary:       netSalary,
	}
	return resp, nil
}
