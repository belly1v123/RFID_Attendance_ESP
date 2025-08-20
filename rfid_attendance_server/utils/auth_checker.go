package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ronishg27/rfid_attendance/constants"
)

func IsSystemAdmin(c *gin.Context) bool {
	_, _, role := GetAuthContext(c)

	if role != string(constants.SuperAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return false
	}
	return true
}

func IsOrgAdmin(c *gin.Context) bool {
	_, _, role := GetAuthContext(c)

	if role != string(constants.OrgAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return false
	}
	return true
}
