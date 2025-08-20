package utils

import (
	"github.com/gin-gonic/gin"
)

func GetAuthContext(c *gin.Context) (userID, email, role string) {
	uid, _ := c.Get("user_id")
	em, _ := c.Get("email")
	rl, _ := c.Get("role")
	return uid.(string), em.(string), rl.(string)
}

func GetUserID(c *gin.Context) string {
	if userID, ok := c.Get("user_id"); ok {
		if idStr, ok := userID.(string); ok {
			return idStr
		}
	}
	return ""
}
