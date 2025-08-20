package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ronishg27/rfid_attendance/config"
	"github.com/ronishg27/rfid_attendance/internal/db"
)

// TenantMiddleware injects tenant-specific queries into context
func TenantMiddleware(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminRole, _ := c.Get("role")
		tenantID := c.GetHeader("X-Tenant-ID")

		if tenantID == "" && adminRole != "super_admin" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "X-Tenant-ID header required"})
			c.Abort()
			return
		}

		// Look up tenant schema
		org, err := queries.Public.Organization.WithContext(c).
			Where(queries.Public.Organization.ID.Eq(tenantID)).
			First()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
			c.Abort()
			return
		}

		// Create tenant query object
		tenantQ, err := queries.Tenant(config.DB, org.Schema)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant schema"})
			c.Abort()
			return
		}

		// Inject into context
		c.Set("tenantQ", tenantQ)

		c.Next()
	}
}

/* in handlers

tenantQ := c.MustGet("tenantQ").(*db.TenantQueries)
tenantQ.Member.Create(&tenant.Member{...})


*/
