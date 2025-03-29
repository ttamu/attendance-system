package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func InitDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	env := os.Getenv("GO_ENV")
	sslmode := "disable"
	if env == "production" {
		sslmode = "require"
	}

	var logLv logger.LogLevel
	if env == "production" {
		logLv = logger.Warn
	} else {
		logLv = logger.Info
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Tokyo",
		host, user, password, dbname, port, sslmode)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLv,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  env != "production",
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("failed to connect database")
	}

	DB = db
}
