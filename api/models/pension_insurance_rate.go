package models

import "gorm.io/gorm"

type PensionInsuranceRate struct {
	gorm.Model
	PrefectureID     uint
	Grade            string
	MinMonthlyAmount int
	MaxMonthlyAmount int
	PensionTotal     float64
	PensionHalf      float64
}
