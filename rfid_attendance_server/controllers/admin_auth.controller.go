package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ronishg27/rfid_attendance/config"
	"github.com/ronishg27/rfid_attendance/constants"
	models "github.com/ronishg27/rfid_attendance/internal/models/temp"

	"github.com/ronishg27/rfid_attendance/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	UserId         string `json:"user_id"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	AdminLevel     string `json:"admin_level"`
	PhoneNumber    string `json:"phone_number"`
	OrganizationID string `json:"organization_id"`
}

func AdminSignIn(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var admin models.Admin
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
	token := utils.GenerateJWT(admin.ID, "super_admin")
	admin.AuthToken = token

	config.DB.Save(&admin)
	// Development setup - no domain, insecure (HTTP)
	c.SetCookie("auth_token", token, 3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func AdminLogout(c *gin.Context) {
	// email, _ := c.Get("email")
	userId, _ := c.Get("user_id")

	var admin models.Admin

	config.DB.Where("id = ?", userId).First(&admin)

	admin.AuthToken = ""
	config.DB.Save(&admin)

	c.SetCookie("auth_token", "", 0, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func AdminRegister(c *gin.Context) {
	userId := utils.GetUserID(c)
	var req models.Admin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		log.Print(err.Error())
		return
	}

	req.Password = string(hashedPassword)
	req.CreatedByID = &userId
	if err := config.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin"})
		log.Print(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin created successfully"})
}

func AdminUpdate(c *gin.Context) {
	var request AuthRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, _, role := utils.GetAuthContext(c)
	// if utils.IsSystemAdmin(c) {

	// }

	if role != string(constants.SuperAdmin) && role != string(constants.OrgAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if role == string(constants.OrgAdmin) && request.UserId != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var admin models.Admin
	if err := config.DB.Where("id = ?", request.UserId).First(&admin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}
}
