package models

import "time"

type Attendance struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	EmployeeID uint      `json:"employee_id"`
	CheckIn    time.Time `json:"check_in"`
	CheckOut   time.Time `json:"check_out"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
