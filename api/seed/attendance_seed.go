package seed

import (
	"log"
	"time"

	"github.com/t2469/labor-management-system.git/models"
	"gorm.io/gorm"
)

func SeedAttendances(db *gorm.DB) error {
	baseDate := time.Date(2025, time.March, 1, 9, 0, 0, 0, time.Local)
	for i := 0; i < 5; i++ {
		day := baseDate.AddDate(0, 0, i)
		attendance := models.Attendance{
			EmployeeID: 1,
			CheckIn:    day,
			CheckOut:   day.Add(9 * time.Hour),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		if err := db.Create(&attendance).Error; err != nil {
			log.Printf("failed to create attendance for employee 1 on %v: %v", day.Format("2006-01-02"), err)
			return err
		}
		log.Printf("Created attendance for employee 1 on %v", day.Format("2006-01-02"))
	}
	return nil
}
