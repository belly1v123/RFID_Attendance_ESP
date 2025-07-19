package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ronishg27/rfid_attendance/config"
	"github.com/ronishg27/rfid_attendance/models"
	"github.com/ronishg27/rfid_attendance/utils"
	"golang.org/x/crypto/bcrypt"
)

func AdminSignIn(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var admin models.SystemAdmin
	err := config.DB.Where("email = ?", req.Email).First(&admin).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid credentials"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create JWT token, store and send it back to the client
	token := utils.GenerateJWT(&admin, nil)
	admin.AuthToken = token

	config.DB.Save(&admin)
	c.JSON(http.StatusOK, gin.H{"token": token})

}
