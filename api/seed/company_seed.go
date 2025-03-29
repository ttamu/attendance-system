package seed

import (
	"log"

	"github.com/t2469/attendance-system.git/models"
	"gorm.io/gorm"
)

func SeedCompanies(db *gorm.DB) error {
	companies := []models.Company{
		{
			Name:         "株式会社サンプル",
			PrefectureID: 13,
		},
		{
			Name:         "株式会社テスト",
			PrefectureID: 1,
		},
	}
	for _, comp := range companies {
		if err := db.Create(&comp).Error; err != nil {
			log.Printf("failed to create company %s: %v", comp.Name, err)
			return err
		}
		log.Println("Created company", comp.Name)
	}
	return nil
}
