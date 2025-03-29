package models

import "time"

type AllowanceType struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CompanyID      uint      `json:"company_id"`
	Name           string    `gorm:"not null" json:"name"`
	Type           string    `gorm:"not null;check:type IN ('commission','fixed')" json:"type"`
	Description    string    `json:"description"`
	CommissionRate float64   `json:"commission_rate"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
