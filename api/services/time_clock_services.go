package services

import (
	"github.com/t2469/attendance-system.git/db"
	"github.com/t2469/attendance-system.git/models"
	"time"
)

func RecordTimeClock(empId uint, clockType models.TimeClockType, timestamp time.Time) (models.TimeClock, error) {
	timeClock := models.TimeClock{
		EmployeeID: empId,
		Type:       clockType,
		Timestamp:  timestamp,
	}

	if err := db.DB.Create(&timeClock).Error; err != nil {
		return timeClock, err
	}

	day := timestamp.Truncate(24 * time.Hour)
	if err := UpsertWorkRecord(empId, day); err != nil {
		return timeClock, err
	}

	return timeClock, nil
}
