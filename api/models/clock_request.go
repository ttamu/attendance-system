package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type RequestStatus string

const (
	Pending  RequestStatus = "pending"
	Approved RequestStatus = "approved"
	Rejected RequestStatus = "rejected"
)

var validStatuses = map[RequestStatus]bool{
	Pending:  true,
	Approved: true,
	Rejected: true,
}

type ClockRequest struct {
	ID         uint          `gorm:"primaryKey" json:"id"`
	EmployeeID uint          `gorm:"not null" json:"employee_id"`
	ClockID    uint          `gorm:"not null" json:"clock_id"`
	Type       TimeClockType `gorm:"type:varchar(20);not null" json:"type"`
	Time       time.Time     `gorm:"not null" json:"time"`
	Status     RequestStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Reason     string        `gorm:"type:text" json:"reason"`
	ReviewedBy *uint         `json:"reviewed_by,omitempty"`
	Reviewer   *Account      `gorm:"foreignKey:ReviewedBy" json:"reviewer,omitempty"`
	ReviewedAt *time.Time    `json:"reviewed_at,omitempty"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

func (r *ClockRequest) BeforeCreate(tx *gorm.DB) error {
	return r.validate()
}

func (r *ClockRequest) BeforeUpdate(tx *gorm.DB) error {
	return r.validate()
}

func (r *ClockRequest) validate() error {
	// status
	if !validStatuses[r.Status] {
		return errors.New("invalid request status: " + string(r.Status))
	}

	// type
	switch r.Type {
	case ClockIn, ClockOut, BreakBegin, BreakEnd:
	default:
		return errors.New("invalid time clock type: " + string(r.Type))
	}

	return nil
}
