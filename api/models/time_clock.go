package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type TimeClockType string

const (
	ClockIn    TimeClockType = "clock_in"
	ClockOut   TimeClockType = "clock_out"
	BreakBegin TimeClockType = "break_begin"
	BreakEnd   TimeClockType = "break_end"
)

type TimeClock struct {
	ID         uint          `gorm:"primaryKey" json:"id"`
	EmployeeID uint          `gorm:"not null" json:"employee_id"`
	Type       TimeClockType `gorm:"not null;check:type IN ('clock_in','clock_out','break_begin','break_end')" json:"type"`
	Timestamp  time.Time     `gorm:"not null" json:"timestamp"`
	CreatedAt  time.Time     `gorm:"autoCreateTime" json:"created_at"`
}

func (tc *TimeClock) BeforeCreate(tx *gorm.DB) (err error) {
	return tc.validate()
}

func (tc *TimeClock) BeforeUpdate(tx *gorm.DB) (err error) {
	return tc.validate()
}

// 想定されたType以外をDBに保存しないためのバリデーション
func (tc *TimeClock) validate() error {
	switch tc.Type {
	case ClockIn, ClockOut, BreakBegin, BreakEnd:
		return nil
	default:
		return errors.New("invalid time clock type: " + string(tc.Type))
	}
}
