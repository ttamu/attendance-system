package seed

import (
	"log"
	"time"

	"github.com/t2469/attendance-system.git/models"
	"gorm.io/gorm"
)

func SeedEmployees(db *gorm.DB) error {
	var companyTokyo models.Company
	if err := db.Where("name = ?", "株式会社東京").First(&companyTokyo).Error; err != nil {
		log.Printf("failed to get company 株式会社東京: %v", err)
		return err
	}

	var companyHokkaido models.Company
	if err := db.Where("name = ?", "株式会社北海道").First(&companyHokkaido).Error; err != nil {
		log.Printf("failed to get company 株式会社北海道: %v", err)
		return err
	}

	employees := []models.Employee{
		{
			CompanyID:     companyTokyo.ID,
			Name:          "山田太郎",
			MonthlySalary: 300000,
			DateOfBirth:   time.Date(1985, time.April, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			CompanyID:     companyHokkaido.ID,
			Name:          "佐藤花子",
			MonthlySalary: 280000,
			DateOfBirth:   time.Date(1990, time.June, 10, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, emp := range employees {
		var employee models.Employee
		result := db.Where(&models.Employee{
			CompanyID: emp.CompanyID,
			Name:      emp.Name,
		}).
			Attrs(models.Employee{
				MonthlySalary: emp.MonthlySalary,
				DateOfBirth:   emp.DateOfBirth,
			}).
			FirstOrCreate(&employee)
		if result.Error != nil {
			log.Printf("failed to create employee %s: %v", emp.Name, result.Error)
			return result.Error
		}
		if result.RowsAffected > 0 {
			log.Println("Created employee", emp.Name)
		} else {
			log.Println("Employee already exists:", emp.Name)
		}
	}
	return nil
}
