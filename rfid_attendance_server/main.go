package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ronishg27/rfid_attendance/config"
	"github.com/ronishg27/rfid_attendance/models"
)

func main() {
	config.InitDB()

	// Auto-migrate user table
	err := config.DB.AutoMigrate(&models.User{})
	checkPanicError(err)
	r := gin.Default()

	r.GET("/test-db", func(c *gin.Context) {
		// Insert test user (skip if exists)
		testUser := models.User{
			Name:       "Test User",
			RFIDUID:    "TEST123456",
			Department: "CSIT",
			Status:     "active",
		}
		config.DB.FirstOrCreate(&testUser, models.User{RFIDUID: "TEST123456"})

		// Fetch all users
		var users []models.User
		result := config.DB.Find(&users)
		if result.Error != nil {
			c.JSON(500, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(200, users)
	})

	err = r.Run(":8080")
	checkPanicError(err)
}

func checkPanicError(err error) {
	if err != nil {
		panic(err)
	}
}
