package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ronishg27/rfid_attendance/config"
	"github.com/ronishg27/rfid_attendance/models"
)

type ScannedRfidRequest struct {
	RFID_UID string `json:"uid"`
}

func HandleRFIDScan(c *gin.Context) {
	var req ScannedRfidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	if req.RFID_UID == "" {
		c.JSON(400, gin.H{"error": "RFID UID is required"})
		return
	}

	db := config.DB

	// 1. check if 	rfid exits in users
	var user models.User
	result := db.Where("rfid_uid = ?", req.RFID_UID).First(&user)

	// 2. log the scan attempts
	db.Create(&models.Log{
		RFIDUID:    req.RFID_UID,
		Recognized: result.RowsAffected > 0,
	})

}
