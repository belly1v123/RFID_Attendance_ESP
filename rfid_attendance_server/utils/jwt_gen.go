package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/ronishg27/rfid_attendance/constants"
	"github.com/ronishg27/rfid_attendance/models"
)

func GenerateJWT(systemAdmin *models.SystemAdmin, orgAdmin *models.OrganizationAdmin) string {
	var userID, email string

	if systemAdmin != nil {
		userID = systemAdmin.ID
		email = systemAdmin.Email
	} else if orgAdmin != nil {
		userID = orgAdmin.ID
		email = orgAdmin.Email
	} else {
		return "" // or return an error
	}

	role := constants.OrgAdmin
	if systemAdmin != nil {
		role = constants.SystemAdmin
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})

	signedToken, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return signedToken
}
