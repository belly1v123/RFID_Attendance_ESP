package db

import (
	"github.com/ronishg27/rfid_attendance/internal/models/public"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPublicDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// migrate public schema models
	if err := db.AutoMigrate(
		&public.Organization{},
		&public.SuperAdmin{},
		&public.ActionLog{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
