package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ronishg27/rfid_attendance/config"
	"github.com/ronishg27/rfid_attendance/constants"
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
		ScannedAt:  time.Now(),
	})

	if result.RowsAffected > 0 {
		// registered user
		// 3. check if the user has already checked in today
		var attendanceRecord models.AttendanceRecord
		result = db.Where("user_id = ? AND day = ?", user.ID, time.Now().Format("2006-01-02")).First(&attendanceRecord)

		if result.RowsAffected > 0 {
			// user has already checked in today
			// checking duplicate entry for 5 mins to skip double entry

			if time.Since(*attendanceRecord.CheckIn).Minutes() < constants.EntryDuplicationDelay {
				c.JSON(200, gin.H{"message": "User has already checked in today"})
			}

		}
	}
}
