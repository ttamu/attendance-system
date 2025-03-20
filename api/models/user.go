package models

import "time"

type User struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	CompanyID     uint         `json:"company_id"`
	Name          string       `json:"name"`
	MonthlySalary int          `json:"monthly_salary"`
	DateOfBirth   time.Time    `json:"date_of_birth"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	Attendances   []Attendance `json:"attendances,omitempty" gorm:"foreignKey:UserID"`
}
