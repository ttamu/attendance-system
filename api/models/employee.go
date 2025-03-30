package models

import "time"

type Employee struct {
	ID            uint                `gorm:"primaryKey" json:"id"`
	CompanyID     uint                `json:"company_id"`
	Company       Company             `json:"company" gorm:"foreignKey:CompanyID"`
	Name          string              `json:"name"`
	MonthlySalary int                 `json:"monthly_salary"`
	DateOfBirth   time.Time           `json:"date_of_birth"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	TimeClocks    []TimeClock         `json:"time_clocks,omitempty" gorm:"foreignKey:EmployeeID"`
	Allowances    []EmployeeAllowance `json:"allowances,omitempty" gorm:"foreignKey:EmployeeID"`
}
