package main

import (
	"github.com/t2469/labor-management-system.git/config"
	"github.com/t2469/labor-management-system.git/db"
	"github.com/t2469/labor-management-system.git/models"
	"github.com/t2469/labor-management-system.git/routes"
	"log"
	"time"
)

func main() {
	cfg := config.LoadEnv()

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic("failed to load location")
	}
	time.Local = loc

	db.InitDB()
	if err := db.DB.AutoMigrate(
		&models.Prefecture{},
		&models.Company{},
		&models.Employee{},
		&models.Attendance{},
		&models.HealthInsuranceRate{},
		&models.PensionInsuranceRate{},
		&models.AllowanceType{},
		&models.EmployeeAllowance{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	routes.Run(cfg)
}
