package seed

import (
	"log"

	"github.com/t2469/attendance-system.git/models"
	"gorm.io/gorm"
)

func SeedCompanies(db *gorm.DB) error {
	var tokyo models.Prefecture
	if err := db.Where("name = ?", "東京").First(&tokyo).Error; err != nil {
		log.Printf("failed to get prefecture 東京: %v", err)
		return err
	}

	var hokkaido models.Prefecture
	if err := db.Where("name = ?", "北海道").First(&hokkaido).Error; err != nil {
		log.Printf("failed to get prefecture 北海道: %v", err)
		return err
	}

	companies := []models.Company{
		{
			Name:         "株式会社東京",
			PrefectureID: tokyo.ID,
		},
		{
			Name:         "株式会社北海道",
			PrefectureID: hokkaido.ID,
		},
	}

	for _, comp := range companies {
		var company models.Company
		result := db.Where(&models.Company{Name: comp.Name}).
			Attrs(models.Company{PrefectureID: comp.PrefectureID}).
			FirstOrCreate(&company)
		if result.Error != nil {
			log.Printf("failed to create company %s: %v", comp.Name, result.Error)
			return result.Error
		}
		if result.RowsAffected > 0 {
			log.Println("Created company", comp.Name)
		} else {
			log.Println("Company already exists:", comp.Name)
		}
	}
	return nil
}
