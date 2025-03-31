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

	tokyoEmployees := []models.Employee{
		{
			CompanyID:     companyTokyo.ID,
			Name:          "東京なこ",
			MonthlySalary: 300000,
			DateOfBirth:   time.Date(1985, time.April, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			CompanyID:     companyTokyo.ID,
			Name:          "東京じろう",
			MonthlySalary: 320000,
			DateOfBirth:   time.Date(1988, time.May, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			CompanyID:     companyTokyo.ID,
			Name:          "東京さぶろう",
			MonthlySalary: 310000,
			DateOfBirth:   time.Date(1990, time.January, 20, 0, 0, 0, 0, time.UTC),
		},
	}

	hokkaidoEmployees := []models.Employee{
		{
			CompanyID:     companyHokkaido.ID,
			Name:          "北海道太郎",
			MonthlySalary: 280000,
			DateOfBirth:   time.Date(1985, time.February, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			CompanyID:     companyHokkaido.ID,
			Name:          "北海道花子",
			MonthlySalary: 270000,
			DateOfBirth:   time.Date(1987, time.June, 18, 0, 0, 0, 0, time.UTC),
		},
		{
			CompanyID:     companyHokkaido.ID,
			Name:          "北海道次郎",
			MonthlySalary: 290000,
			DateOfBirth:   time.Date(1992, time.December, 25, 0, 0, 0, 0, time.UTC),
		},
	}

	employees := append(tokyoEmployees, hokkaidoEmployees...)
	for _, emp := range employees {
		var employee models.Employee
		result := db.Where(&models.Employee{
			CompanyID: emp.CompanyID,
			Name:      emp.Name,
		}).
			Attrs(models.Employee{
				MonthlySalary: emp.MonthlySalary,
				DateOfBirth:   emp.DateOfBirth,
			}).FirstOrCreate(&employee)
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
