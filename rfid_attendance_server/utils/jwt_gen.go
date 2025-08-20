package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(userID string, role string) string {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})

	signedToken, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return signedToken
}
