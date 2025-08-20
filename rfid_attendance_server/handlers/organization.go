package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ronishg27/rfid_attendance/config"
	"github.com/ronishg27/rfid_attendance/internal/db"
	"github.com/ronishg27/rfid_attendance/internal/models/public"
	"github.com/ronishg27/rfid_attendance/utils"
)

type OrgReq struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

func RegisterOrganizationHandler(c *gin.Context) {

	// queries := c.MustGet("queries").(*db.Queries)

	var req OrgReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org := &public.Organization{
		Name:   req.Name,
		Email:  req.Email,
		Schema: utils.GetTenantName(req.Name),
	}

	// ✅ Delegate schema creation + migration to ProvisionTenant
	if err := db.ProvisionTenant(config.DB, org); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// later you can use tenantQ.Member.WithContext(c).Create(&tenant.Member{...})
	c.JSON(http.StatusCreated, gin.H{
		"message": "Organization provisioned successfully",
		"org":     org,
	})
}
