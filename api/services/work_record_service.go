package services

import (
	"errors"
	"time"

	"github.com/t2469/attendance-system.git/db"
	"github.com/t2469/attendance-system.git/models"
	"gorm.io/gorm"
)

func UpsertWorkRecord(empID uint, date time.Time) error {
	from := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	to := from.AddDate(0, 0, 1)

	// 対象日の TimeClock レコードを取得
	var clocks []models.TimeClock
	if err := db.DB.
		Where("employee_id = ? AND timestamp >= ? AND timestamp < ?", empID, from, to).
		Order("timestamp ASC").
		Find(&clocks).Error; err != nil {
		return err
	}

	var clockIn, clockOut time.Time
	var breakDuration time.Duration
	var breakStart *time.Time

	for _, clock := range clocks {
		switch clock.Type {
		case models.ClockIn:
			if clockIn.IsZero() {
				clockIn = clock.Timestamp
			}
		case models.ClockOut:
			clockOut = clock.Timestamp
		case models.BreakBegin:
			breakStart = &clock.Timestamp
		case models.BreakEnd:
			if breakStart != nil {
				breakDuration += clock.Timestamp.Sub(*breakStart)
				breakStart = nil
			}
		}
	}

	// clockIn/clockOut が揃っていれば、実働時間を計算。なければ 0
	workDuration := time.Duration(0)
	if !clockIn.IsZero() && !clockOut.IsZero() && clockOut.After(clockIn) {
		workDuration = clockOut.Sub(clockIn) - breakDuration
	}

	// 該当日の WorkRecord を検索し、なければ作成、あれば更新
	var wr models.WorkRecord
	err := db.DB.Where("employee_id = ? AND date = ?", empID, from).First(&wr).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		wr = models.WorkRecord{
			EmployeeID:    empID,
			Date:          from,
			ClockIn:       clockIn,
			ClockOut:      clockOut,
			BreakDuration: breakDuration,
			WorkDuration:  workDuration,
		}
		return db.DB.Create(&wr).Error
	} else if err != nil {
		return err
	}

	wr.ClockIn = clockIn
	wr.ClockOut = clockOut
	wr.BreakDuration = breakDuration
	wr.WorkDuration = workDuration

	return db.DB.Save(&wr).Error
}
