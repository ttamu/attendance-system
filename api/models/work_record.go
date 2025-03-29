package models

import "time"

type WorkRecord struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	EmployeeID   uint      `gorm:"not null" json:"employee_id"`
	Date         time.Time `gorm:"type:date;not null" json:"date"`
	ClockIn      time.Time `json:"clock_in"`
	ClockOut     time.Time `json:"clock_out"`
	BreakMinutes int64     `json:"break_minutes"`
	WorkMinutes  int64     `json:"work_minutes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
