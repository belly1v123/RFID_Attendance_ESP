package config

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ronishg27/rfid_attendance/internal/models/public"
	"github.com/ronishg27/rfid_attendance/internal/my_queries"
	"github.com/ronishg27/rfid_attendance/utils"
	"golang.org/x/crypto/bcrypt"
)

func SeedSuperAdmin(query *my_queries.Query) error {

	sAdminQuery := query.SuperAdmin

	count, _ := sAdminQuery.Count()
	if count > 0 {
		log.Println("Super admin already exists, skipping seed.")
		return nil
	}

	err := godotenv.Load()
	utils.HandleError(err, true)

	password := os.Getenv("SUPER_ADMIN_PASSWORD")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &public.SuperAdmin{
		Name:     "Super Admin",
		Email:    "a@admin.com",
		Password: string(passwordHash),
	}

	log.Println("Seeding super admin...")
	return sAdminQuery.WithContext(context.Background()).Create(admin)

}
