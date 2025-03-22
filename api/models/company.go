package models

import "time"

type Company struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Name         string     `json:"name"`
	PrefectureID uint       `json:"prefecture_id"`
	Prefecture   Prefecture `json:"prefecture" gorm:"foreignKey:PrefectureID"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
