package models

import "time"

type RequestStatus string

const (
	Pending  RequestStatus = "pending"
	Approved RequestStatus = "approved"
	Rejected RequestStatus = "rejected"
)

type ClockRequest struct {
	ID         uint          `gorm:"primaryKey"`
	EmployeeID uint          `gorm:"not null"`
	ClockID    uint          `gorm:"not null"`
	Type       TimeClockType `gorm:"type:varchar(20);not null"`
	Time       time.Time     `gorm:"not null"`
	Status     RequestStatus `gorm:"type:varchar(20);default:'pending'"`
	Reason     string        `gorm:"type:text"`
	ReviewedBy *uint
	ReviewedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
