package repository

import (
	"crud_api_us/internal/models"

	"gorm.io/gorm"
)

func MigrateAndSeed(db *gorm.DB, seed []models.User) error {
	if err := db.AutoMigrate(&models.User{}, &models.RefreshToken{}); err != nil {
		return err
	}
	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 && len(seed) > 0 {
		if err := db.Create(&seed).Error; err != nil {
			return err
		}
	}
	return nil
}
