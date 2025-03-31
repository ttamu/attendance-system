package seed

import (
	"github.com/t2469/attendance-system.git/models"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAccounts(db *gorm.DB) error {
	var tokyoCompany models.Company
	if err := db.Where("name = ?", "株式会社東京").First(&tokyoCompany).Error; err != nil {
		log.Printf("failed to get company 株式会社東京: %v", err)
		return err
	}

	var hokkaidoCompany models.Company
	if err := db.Where("name = ?", "株式会社北海道").First(&hokkaidoCompany).Error; err != nil {
		log.Printf("failed to get company 株式会社北海道: %v", err)
		return err
	}

	rawPassword := "password"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash password: %v", err)
		return err
	}

	accounts := []models.Account{
		{
			CompanyID: tokyoCompany.ID,
			Email:     "tokyo_admin@example.com",
			Password:  string(hashedPassword),
			IsAdmin:   true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			CompanyID: tokyoCompany.ID,
			Email:     "tokyo_user@example.com",
			Password:  string(hashedPassword),
			IsAdmin:   false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			CompanyID: hokkaidoCompany.ID,
			Email:     "hokkaido_admin@example.com",
			Password:  string(hashedPassword),
			IsAdmin:   true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			CompanyID: hokkaidoCompany.ID,
			Email:     "hokkaido_user@example.com",
			Password:  string(hashedPassword),
			IsAdmin:   false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, acc := range accounts {
		var account models.Account
		result := db.Where(&models.Account{
			CompanyID: acc.CompanyID,
			Email:     acc.Email,
		}).Attrs(models.Account{
			Password:  acc.Password,
			IsAdmin:   acc.IsAdmin,
			CreatedAt: acc.CreatedAt,
			UpdatedAt: acc.UpdatedAt,
		}).FirstOrCreate(&account)
		if result.Error != nil {
			log.Printf("failed to create account %s: %v", acc.Email, result.Error)
			return result.Error
		}
		if result.RowsAffected > 0 {
			log.Printf("Created account %s", acc.Email)
		} else {
			log.Printf("Account already exists: %s", acc.Email)
		}
	}
	return nil
}
