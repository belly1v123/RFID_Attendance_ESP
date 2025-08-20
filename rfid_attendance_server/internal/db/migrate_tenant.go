package db

import (
	"fmt"

	"github.com/ronishg27/rfid_attendance/internal/models/public"
	"github.com/ronishg27/rfid_attendance/internal/models/tenant"
	"gorm.io/gorm"
)

func ProvisionTenant(db *gorm.DB, org *public.Organization) error {
	// 1. Insert org into public.organizations
	if err := db.Create(org).Error; err != nil {
		return err
	}

	// 2. Create tenant schema
	schemaName := org.Schema
	createSchemaSQL := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName)
	if err := db.Exec(createSchemaSQL).Error; err != nil {
		return err
	}

	// 3. Set search_path and migrate tenant models
	tenantDB := db.Session(&gorm.Session{NewDB: true})
	if err := tenantDB.Exec(fmt.Sprintf("SET search_path TO %s", schemaName)).Error; err != nil {
		return err
	}

	if err := tenantDB.AutoMigrate(
		&tenant.Member{},
		&tenant.AttendanceRecord{},
		&tenant.ScanLog{},
		&tenant.Device{},
		&tenant.Role{},
		&tenant.Shift{},
		&tenant.Admin{},
	); err != nil {
		return err
	}

	return nil
}
