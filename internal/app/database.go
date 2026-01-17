package app

import (
	"database/sql"
	"errors"
	"os"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDb() (*gorm.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("No database url were provided")
	}

	sqlDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	defer sqlDB.Close()

	if err := goose.Up(sqlDB, "./migrations"); err != nil {
		return nil, err
	}
	gormDB, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}
