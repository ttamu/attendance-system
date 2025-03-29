package seed

import (
	"log"
	"time"

	"github.com/t2469/attendance-system.git/models"
	"gorm.io/gorm"
)

func SeedEmployees(db *gorm.DB) error {
	employees := []models.Employee{
		{
			CompanyID:     1,
			Name:          "山田太郎",
			MonthlySalary: 300000,
			DateOfBirth:   time.Date(1985, time.April, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			CompanyID:     2,
			Name:          "佐藤花子",
			MonthlySalary: 280000,
			DateOfBirth:   time.Date(1990, time.June, 10, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, emp := range employees {
		if err := db.Create(&emp).Error; err != nil {
			log.Printf("failed to create employee %s: %v", emp.Name, err)
			return err
		}
		log.Println("Created employee", emp.Name)
	}
	return nil
}
