package main

import (
	"github.com/t2469/attendance-system.git/config"
	"github.com/t2469/attendance-system.git/db"
	"github.com/t2469/attendance-system.git/models"
	"github.com/t2469/attendance-system.git/routes"
	"github.com/t2469/attendance-system.git/services"
	"log"
	"time"
)

func main() {
	cfg := config.LoadEnv()

	// LineBotの初期化
	services.InitLineBot(cfg)

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
		&models.Account{},
		&models.Attendance{},
		&models.TimeClock{},
		&models.WorkRecord{},
		&models.ClockRequest{},
		&models.HealthInsuranceRate{},
		&models.PensionInsuranceRate{},
		&models.AllowanceType{},
		&models.EmployeeAllowance{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	routes.Run(cfg)
}
