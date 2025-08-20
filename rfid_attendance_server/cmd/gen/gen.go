package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ronishg27/rfid_attendance/internal/models/public"
	"github.com/ronishg27/rfid_attendance/internal/models/tenant"
	"gorm.io/driver/postgres"
	"gorm.io/gen"

	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:       "./internal/my_queries",
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
		FieldNullable: true,
	})

	// g.UseDB(db)

	// register public schema models
	g.ApplyBasic(
		new(public.Organization),
		new(public.SuperAdmin),
		new(public.ActionLog),
	)

	// register tenant schema models
	g.ApplyBasic(
		new(tenant.Member),
		new(tenant.AttendanceRecord),
		new(tenant.ScanLog),
		new(tenant.Device),
		new(tenant.Role),
		new(tenant.Shift),
	)

	g.Execute()
	log.Println("✅ Query code generated successfully")

	log.Println("Migrating the public schema")
	if err := db.AutoMigrate(
		&public.Organization{},
		&public.SuperAdmin{},
		&public.ActionLog{},
	); err != nil {
		log.Fatalf("Failed to migrate public schema: %v", err)
	}
	log.Println("✅ Public schema migrated successfully")

	// Example tenant migration for existing tenants
	// tenantSchemas := []string{"tenant_test"} // optionally fetch from public.organizations
	// for _, schema := range tenantSchemas {
	// 	db.Exec("CREATE SCHEMA IF NOT EXISTS " + schema)
	// 	db.Exec("SET search_path TO " + schema)

	// 	if err := db.AutoMigrate(
	// 		&tenant.Member{},
	// 		&tenant.AttendanceRecord{},
	// 		&tenant.ScanLog{},
	// 		&tenant.Device{},
	// 		&tenant.Role{},
	// 		&tenant.Shift{},
	// 		&tenant.Admin{},
	// 	); err != nil {
	// 		log.Fatalf("Failed to migrate tenant schema %s: %v", schema, err)
	// 	}
	// 	log.Printf("✅ Tenant schema %s migrated successfully\n", schema)
	// }

	// log.Println("All migrations completed!")

}
