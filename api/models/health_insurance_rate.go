package models

import "time"

type HealthInsuranceRate struct {
	ID                  uint `gorm:"primaryKey" json:"id"`
	CompanyID           uint `json:"company_id"`
	PrefectureID        uint
	Grade               string
	MonthlyAmount       int
	MinMonthlyAmount    int
	MaxMonthlyAmount    int
	HealthTotalNonCare  float64
	HealthHalfNonCare   float64
	HealthTotalWithCare float64
	HealthHalfWithCare  float64
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
