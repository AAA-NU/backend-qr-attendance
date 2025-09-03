package database

import (
	"github.com/aaanu/backend-qr-attendance/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Initialize(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Автомиграция
	err = db.AutoMigrate(&models.QRCode{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
