package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ronishg27/rfid_attendance/config"
	"github.com/ronishg27/rfid_attendance/controllers"
	"github.com/ronishg27/rfid_attendance/middleware"
	"github.com/ronishg27/rfid_attendance/models"
	"github.com/ronishg27/rfid_attendance/utils"
)

func main() {
	config.InitDB()

	// Auto-migrate user table
	err := config.DB.AutoMigrate(
		// &models.User{},
		&models.SystemAdmin{},
		&models.Organization{},
		&models.OrganizationAdmin{},
		&models.OrganizationMember{},
		&models.ScanLog{},
		&models.AttendanceRecord{},
		&models.ActionLog{},
	)

	utils.HandleError(err, true)
	err = config.SeedSuperAdmin(config.DB)
	utils.HandleError(err, false)
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/api/scan", controllers.HandleRFIDScan)
	r.POST("/api/register", controllers.RegisterUser)
	r.POST("/api/admin/login", controllers.SystemAdminSignIn)
	r.POST("/api/admin/logout", middleware.JWTAuthMiddleware(), controllers.SystemAdminLogout)

	r.POST("/api/admin/org", middleware.JWTAuthMiddleware(), controllers.CreateOrganization)

	err = r.Run(":3000")
	utils.HandleError(err, true)
}
