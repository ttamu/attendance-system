package models

import "time"

type Company struct {
	ID                   uint                  `gorm:"primaryKey" json:"id"`
	Name                 string                `json:"name"`
	HealthInsuranceRates []HealthInsuranceRate `json:"health_insurance_rates" gorm:"foreignKey:CompanyID"`
	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at"`
}
