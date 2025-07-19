package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/ronishg27/rfid_attendance/models"
	"github.com/ronishg27/rfid_attendance/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedSuperAdmin(db *gorm.DB) error {
	var count int64
	db.Model(&models.SystemAdmin{}).Count(&count)
	if count > 0 {
		log.Println("Super admin already exists, skipping seed.")
		return nil
	}

	err := godotenv.Load()
	utils.HandleError(err, true)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("SuperSecret123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &models.SystemAdmin{
		Name:        "Super Admin",
		Email:       "ronishunofficial@gmail.com",
		Password:    string(passwordHash),
		PhoneNumber: "9800000000",
		Role:        "SuperAdmin",
	}

	log.Println("Seeding super admin...")
	return db.Create(&admin).Error
}
