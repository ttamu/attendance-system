package models

import "time"

type EmployeeAllowance struct {
	ID              uint          `gorm:"primaryKey" json:"id"`
	EmployeeID      uint          `json:"employee_id"`
	AllowanceTypeID uint          `json:"allowance_type_id"`
	Amount          int           `json:"amount"`
	CommissionRate  float64       `json:"commission_rate,omitempty"`
	Year            int           `json:"year"`
	Month           int           `json:"month"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	AllowanceType   AllowanceType `json:"allowance_type" gorm:"foreignKey:AllowanceTypeID"`
}
