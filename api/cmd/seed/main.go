package main

import (
	"github.com/t2469/attendance-system.git/config"
	"github.com/t2469/attendance-system.git/db"
	"github.com/t2469/attendance-system.git/seed"
	"log"
	"time"
)

func main() {
	_ = config.LoadEnv()
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	time.Local = loc
	db.InitDB()

	if err := seed.SeedPrefectures(db.DB); err != nil {
		log.Fatalf("seedPrefectures failed: %v", err)
	}

	if err := seed.SeedCompanies(db.DB); err != nil {
		log.Fatalf("SeedCompanies failed: %v", err)
	}

	if err := seed.SeedEmployees(db.DB); err != nil {
		log.Fatalf("SeedEmployees failed: %v", err)
	}
	
	if err := seed.SeedInsuranceRates(db.DB, "seed/insurance_rates.xlsx"); err != nil {
		log.Fatalf("seedInsuranceRates failed: %v", err)
	}

	log.Println("Seeding completed successfully.")
}
