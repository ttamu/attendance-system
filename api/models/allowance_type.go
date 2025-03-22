package models

import "time"

type AllowanceType struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"unique;not null" json:"name"`
	Type           string    `json:"type"` // "commission" or "fixed"
	Description    string    `json:"description"`
	CommissionRate float64   `json:"commission_rate"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
