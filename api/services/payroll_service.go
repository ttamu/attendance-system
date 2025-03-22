package services

import (
	"github.com/t2469/labor-management-system.git/models"
	"gorm.io/gorm"
)

type PayrollCalculationResponse struct {
	EmployeeName    string  `json:"employee_name"`
	GrossSalary     float64 `json:"gross_salary"`
	TotalAllowance  float64 `json:"total_allowance"`
	HealthInsurance float64 `json:"health_insurance"`
	Pension         float64 `json:"pension"`
	TotalDeductions float64 `json:"total_deductions"`
	NetSalary       float64 `json:"net_salary"`
}

func CalculatePayroll(db *gorm.DB, employeeID uint) (PayrollCalculationResponse, error) {
	var emp models.Employee
	if err := db.
		Preload("Company.Prefecture.HealthInsuranceRates").
		Preload("Company.Prefecture.PensionInsuranceRates").
		Preload("Allowances.AllowanceType").
		First(&emp, employeeID).Error; err != nil {
		return PayrollCalculationResponse{}, err
	}

	var totalAllowance float64
	for _, ea := range emp.Allowances {
		switch ea.AllowanceType.Type {
		case "commission":
			// 従業員ごとに設定された割合があればそれを使い、なければ AllowanceType の値を使用
			rate := ea.CommissionRate
			if rate == 0 {
				rate = ea.AllowanceType.CommissionRate
			}
			commissionAmount := float64(ea.Amount) * rate
			totalAllowance += commissionAmount
		case "fixed":
			totalAllowance += float64(ea.Amount)
		}
	}

	healthResp, err := CalculateInsurance(db, employeeID)
	if err != nil {
		return PayrollCalculationResponse{}, err
	}

	pensionResp, err := CalculatePension(db, employeeID)
	if err != nil {
		return PayrollCalculationResponse{}, err
	}

	// 支給額(基本給+手当合計)
	baseSalary := emp.MonthlySalary
	grossSalary := float64(baseSalary) + totalAllowance

	// 控除額(健康保険,介護保険,厚生年金)
	totalDeductions := healthResp.EmployeeHealth + pensionResp.EmployeePension

	// 手取り給与
	netSalary := grossSalary - totalDeductions

	resp := PayrollCalculationResponse{
		EmployeeName:    healthResp.EmployeeName,
		GrossSalary:     grossSalary,
		TotalAllowance:  totalAllowance,
		HealthInsurance: healthResp.EmployeeHealth,
		Pension:         pensionResp.EmployeePension,
		TotalDeductions: totalDeductions,
		NetSalary:       netSalary,
	}
	return resp, nil
}
