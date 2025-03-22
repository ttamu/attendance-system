package models

import "gorm.io/gorm"

type Prefecture struct {
	gorm.Model
	Name                  string `gorm:"unique;not null"`
	HealthRateNoCare      float64
	HealthRateWithCare    float64
	PensionRate           float64
	HealthInsuranceRates  []HealthInsuranceRate
	PensionInsuranceRates []PensionInsuranceRate
}
