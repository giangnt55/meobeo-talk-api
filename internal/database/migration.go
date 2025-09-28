package database

import (
	"meobeo-talk-api/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		// Add other models here
	)
}
