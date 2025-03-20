package models

import "gorm.io/gorm"

type HealthInsuranceRate struct {
	gorm.Model
	PrefectureID        uint
	CompanyID           uint `json:"company_id"`
	Grade               string
	MonthlyAmount       int
	MinMonthlyAmount    int
	MaxMonthlyAmount    int
	HealthTotalNonCare  float64
	HealthHalfNonCare   float64
	HealthTotalWithCare float64
	HealthHalfWithCare  float64
}
