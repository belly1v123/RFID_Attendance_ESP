package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ronishg27/rfid_attendance/internal/db"
	"github.com/ronishg27/rfid_attendance/internal/models/tenant"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")

	// Initialize public DB + global connection
	publicDB, err := db.InitPublicDB(dsn)

	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	DB = publicDB
	fmt.Println("✅ Database connected")
}

// to create schema or tenant for an organization that has been just created
func ProvisionTenant(db *gorm.DB, schema string) error {
	// 1️⃣ Create schema
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS " + schema).Error; err != nil {
		return err
	}

	// 2️⃣ Switch search_path
	if err := db.Exec("SET search_path TO " + schema).Error; err != nil {
		return err
	}

	// 3️⃣ Migrate tenant tables
	return db.AutoMigrate(
		&tenant.Member{},
		&tenant.AttendanceRecord{},
		&tenant.ScanLog{},
		&tenant.Device{},
		&tenant.Role{},
		&tenant.Shift{},
		&tenant.Admin{},
	)
}
