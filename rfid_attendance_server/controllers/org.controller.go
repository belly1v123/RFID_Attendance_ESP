package controllers

// TODO: complete the UpdateOrganization controller,
// TODO: check if the user is authorized to update the organization i.e. if user in orgAdmin check if its belong to org or not --  done ✅
// TODO: complete the DeleteOrganization controller -- ✅

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ronishg27/rfid_attendance/config"
	"github.com/ronishg27/rfid_attendance/constants"
	"github.com/ronishg27/rfid_attendance/models"
	"github.com/ronishg27/rfid_attendance/utils"
)

type CreateOrganizationRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type UpdateOrganizationRequest struct {
	Name                  *string         `json:"name"`
	Email                 *string         `json:"email"`
	Phone                 *string         `json:"phone"`
	Address               *string         `json:"address"`
	AvailableRoles        *map[string]any `json:"available_roles"`
	AvailableShifts       *map[string]any `json:"available_shifts"`
	EntryDuplicationDelay *int            `json:"entry_duplication_delay"`
}

func CreateOrganization(c *gin.Context) {

	adminID, _, role := utils.GetAuthContext(c)

	if role != constants.SystemAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden access"})
		return
	}

	var req CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	org := models.Organization{
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		Address:     req.Address,
		CreatedByID: adminID,
	}

	if err := config.DB.WithContext(c).Create(&org).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{"error": "Organization already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Organization"})
		log.Print(err.Error())
		return
	}

	c.JSON(http.StatusCreated, org)

}

func GetOrganization(c *gin.Context) {
	adminID, _, role := utils.GetAuthContext(c)
	orgID := c.Param("orgID")

	var org models.Organization
	if err := config.DB.WithContext(c).First(&org, "id = ?", orgID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}
	if role == constants.SystemAdmin {
		c.JSON(http.StatusOK, org)
		return
	}

	var orgAdmin models.OrganizationAdmin
	if err := config.DB.WithContext(c).First(&orgAdmin, "id=?", adminID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify admin organization"})
		log.Print("OrgAdmin fetch error:", err)
		return
	}

	if orgAdmin.OrganizationID != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this organization"})
		return
	}

	c.JSON(http.StatusOK, org)
}

func GetAllOrganizations(c *gin.Context) {

	if _, _, role := utils.GetAuthContext(c); role != constants.SystemAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Forbidden access"})
		return
	}

	var orgs []models.Organization
	if err := config.DB.WithContext(c).Find(&orgs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		log.Print(err.Error())
		return
	}
	c.JSON(http.StatusOK, orgs)
}

func UpdateOrganization(c *gin.Context) {
	// adminID, _, role := utils.GetAuthContext(c)

	//1. adminId, role extract
	// 2. if not system admin and not org admin return

	// 3. if system admin, save the updates and return
	// 4. if org admin, check if its belong to org or not
	// 5. if belong to org, save the updates and return

	adminID, _, role := utils.GetAuthContext(c)

	if role != constants.SystemAdmin && role != constants.OrgAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	orgID := c.Param("orgID")
	var req UpdateOrganizationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var org models.Organization
	if err := config.DB.WithContext(c).First(&org, "id = ?", orgID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		log.Print("Org fetch error:", err.Error())
		return
	}

	if role == constants.OrgAdmin {
		var orgAdmin models.OrganizationAdmin

		if err := config.DB.WithContext(c).First(&orgAdmin, "id=?", adminID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify admin organization"})
			log.Print("OrgAdmin fetch error:", err)
			return
		}
		if orgAdmin.OrganizationID != org.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden access"})
			return
		}
	}

	if req.Name != nil {
		org.Name = *req.Name
	}
	if req.Email != nil {
		org.Email = *req.Email
	}
	if req.Phone != nil {
		org.Phone = *req.Phone
	}
	if req.Address != nil {
		org.Address = *req.Address
	}

	if req.AvailableRoles != nil {
		org.AvailableRoles = *req.AvailableRoles
	}
	if req.AvailableShifts != nil {
		org.AvailableShifts = *req.AvailableShifts
	}
	if req.EntryDuplicationDelay != nil {
		org.EntryDuplicationDelay = *req.EntryDuplicationDelay
	}

	org.UpdatedBy = adminID
	org.UpdatedByType = role

	if err := config.DB.WithContext(c).Save(&org).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organization"})
		log.Print(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Organization updated successfully", "organization": org})

}

func DeleteOrganization(c *gin.Context) {
	orgID := c.Param("orgID")
	// soft delete
	adminID, _, role := utils.GetAuthContext(c)
	if role != constants.SystemAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var org models.Organization
	if err := config.DB.WithContext(c).First(&org, "id = ?", orgID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	org.DeletedByID = adminID
	if err := config.DB.WithContext(c).Save(&org).Delete(&org).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete organization"})
		log.Print(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Organization deleted"})

}
