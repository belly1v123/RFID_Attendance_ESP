package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ronishg27/rfid_attendance/config"
	models "github.com/ronishg27/rfid_attendance/internal/models/temp"
)

func RegisterUser(c *gin.Context) {
	var req models.OrganizationMember
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.RFIDUID == "" && req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "RFID UID and Name is required",
		})
	}

	db := config.DB

	var user models.OrganizationMember

	tx := db.Where("rfid_uid = ?", req.RFIDUID).First(&user)
	if tx.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "RFID UID already exists",
		})
		return
	}

	user = models.OrganizationMember{
		RFIDUID:     req.RFIDUID,
		Name:        req.Name,
		Department:  req.Department,
		Status:      req.Status,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}

	res := db.Create(&user).First(&user)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register user", "error": res.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "user": user})
}
