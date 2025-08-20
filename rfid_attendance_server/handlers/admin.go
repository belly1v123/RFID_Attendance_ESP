package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ronishg27/rfid_attendance/internal/db"
	"github.com/ronishg27/rfid_attendance/utils"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	// Token string `json:"token"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"` // super_admin or org_admin
}

func SuperAdminLoginHandler(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		// fetch admin from public schema
		admin, err := queries.Public.SuperAdmin.WithContext(c).Where(
			queries.Public.SuperAdmin.Email.Eq(req.Email)).First()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credentials"})
		}

		// verify password

		if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credentials"})
			return
		}

		token := utils.GenerateJWT(admin.ID, "Super_Admin")
		// setting token as authorization header
		// c.Set("Authorization", "Bearer "+token)
		// c.Header("Authorization", "Bearer "+token)
		c.SetCookie(
			"auth_token", // cookie name
			token,        // your JWT token
			3600,         // max age in seconds (e.g., 1 hour)
			"/",          // path
			"",           // domain (empty means current domain)
			true,         // secure (HTTPS only)
			true,         // HttpOnly (no JavaScript access)
		)
		c.JSON(http.StatusOK, LoginResponse{
			Name:  admin.Name,
			Email: admin.Email,
			Role:  "Super_Admin",
		})
	}
}
