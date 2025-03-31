package models

import "gorm.io/gorm"

type PensionInsuranceRate struct {
	gorm.Model
	PrefectureID     uint
	Grade            string
	MonthlyAmount    int
	MinMonthlyAmount int
	MaxMonthlyAmount int
	PensionTotal     float64
	PensionHalf      float64
	FromYear         int
	FromMonth        int
	ToYear           int
	ToMonth          int
}
