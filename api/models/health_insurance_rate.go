package models

import "time"

type HealthInsuranceRate struct {
	ID                  uint `gorm:"primaryKey" json:"id"`
	PrefectureID        uint
	Grade               string
	MonthlyAmount       int
	MinMonthlyAmount    int
	MaxMonthlyAmount    int
	HealthTotalNonCare  float64
	HealthHalfNonCare   float64
	HealthTotalWithCare float64
	HealthHalfWithCare  float64
	Year                int
	Month               int
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
