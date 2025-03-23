package services

import (
	"github.com/t2469/labor-management-system.git/models"
	"gorm.io/gorm"
	"log"
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

func CalculatePayroll(db *gorm.DB, employeeID uint, year, month int) (PayrollCalculationResponse, error) {
	var emp models.Employee
	if err := db.
		Preload("Company.Prefecture.HealthInsuranceRates").
		Preload("Company.Prefecture.PensionInsuranceRates").
		Preload("Allowances", "year = ? AND month = ?", year, month).
		Preload("Allowances.AllowanceType").
		First(&emp, employeeID).Error; err != nil {
		return PayrollCalculationResponse{}, err
	}

	totalAllowance := calculateTotalAllowance(emp.Allowances)

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

func calculateTotalAllowance(allowances []models.EmployeeAllowance) float64 {
	var total float64
	for _, ea := range allowances {
		switch ea.AllowanceType.Type {
		case "commission":
			// 従業員ごとに設定された割合があればそれを使用、なければデフォルトの値を使用
			rate := ea.CommissionRate
			if rate == 0 {
				rate = ea.AllowanceType.CommissionRate
			}
			total += float64(ea.Amount) * rate
		case "fixed":
			total += float64(ea.Amount)
		default:
			log.Printf("unknown allowance type: %s", ea.AllowanceType.Type)
		}
	}
	return total
}
