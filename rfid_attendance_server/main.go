package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ronishg27/rfid_attendance/config"
	"github.com/ronishg27/rfid_attendance/controllers"
	"github.com/ronishg27/rfid_attendance/models"
	"github.com/ronishg27/rfid_attendance/utils"
)

func main() {
	config.InitDB()

	// Auto-migrate user table
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.AttendanceRecord{},
		&models.Log{},
	)

	utils.HandlePanicError(err)
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/api/scan", controllers.HandleRFIDScan)

	err = r.Run(":3000")
	utils.HandlePanicError(err)
}
