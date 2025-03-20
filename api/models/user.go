package models

import "time"

type User struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Name        string       `json:"name"`
	CompanyID   uint         `json:"company_id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Attendances []Attendance `json:"attendances,omitempty" gorm:"foreignKey:UserID"`
}
