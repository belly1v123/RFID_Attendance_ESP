package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ronishg27/rfid_attendance/config"
	"github.com/ronishg27/rfid_attendance/constants"
	"github.com/ronishg27/rfid_attendance/models"
	"gorm.io/gorm"
)

type ScannedRfidRequest struct {
	RFID_UID string `json:"uid"`
}

func HandleRFIDScan(c *gin.Context) {
	var req ScannedRfidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if req.RFID_UID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "RFID UID is required"})
		return
	}

	db := config.DB
	now := time.Now()
	today := now.Format("2006-01-02")
	tolerance := time.Duration(constants.EntryDuplicationDelay) * time.Minute

	// 1. check if user exists with rfid
	var user models.OrganizationMember
	if err := db.Where("rfid_uid = ?", req.RFID_UID).First(&user).Error; err != nil {
		// unrecognized user - logging the scan
		logUser(db, now, false, req.RFID_UID)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "unrecognized user",
			"status":  "unrecognized",
		})
		return
	}

	// log recognized user scan
	logUser(db, now, true, req.RFID_UID)

	var attendance models.AttendanceRecord
	tx := db.Where("user_id = ? AND day = ?", user.ID, today).First(&attendance)

	if tx.Error != nil {
		newRecord := models.AttendanceRecord{
			UserID:     user.ID,
			PresentDay: now,
			CheckIn:    &now,
			CheckOut:   nil,
			Status:     "checked_in",
		}
		db.Create(&newRecord)

		c.JSON(http.StatusOK, gin.H{
			"message": "user checked in",
			"status":  "checked_in", //green
		})
		return
	}
	if attendance.CheckOut != nil {
		// logUser(db, now, true, req.RFID_UID)

		c.JSON(http.StatusOK, gin.H{
			"message": "user already checked in and checked out today",
			"status":  "duplicate", //yellow
		})
		return
	}
	if attendance.CheckIn != nil && now.Sub(*attendance.CheckIn) < tolerance {
		// logUser(db, now, true, req.RFID_UID)
		c.JSON(http.StatusOK, gin.H{
			"message": "user already checked in recently",
			"status":  "duplicate", //yellow
		})
		return
	}

	//else checkout
	attendance.CheckOut = &now
	attendance.Status = "checked_out"
	db.Save(&attendance)
	c.JSON(http.StatusOK, gin.H{
		"message": "user checked out",
		"status":  "checked_out", //blue
	})
}

func logUser(db *gorm.DB, time time.Time, recognized bool, rfidUID string) {
	db.Create(&models.ScanLog{
		RFIDUID:    rfidUID,
		ScannedAt:  time,
		Recognized: recognized,
	})
}
