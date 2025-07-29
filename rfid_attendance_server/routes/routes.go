package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ronishg27/rfid_attendance/controllers"
	"github.com/ronishg27/rfid_attendance/middleware"
)





func SetupRoutes(r *gin.Engine) {
	apiRoutes := r.Group("/api")
	adminRoutes := r.Group("/api/admin")
	adminOrgRoutes := adminRoutes.Group("/org")

	apiRoutes.POST("/scan", controllers.HandleRFIDScan)
	apiRoutes.POST("/register", middleware.JWTAuthMiddleware(), controllers.RegisterUser)

	adminRoutes.POST("/login", controllers.SystemAdminSignIn)
	adminRoutes.POST("/logout", middleware.JWTAuthMiddleware(), controllers.SystemAdminLogout)

	adminOrgRoutes.POST("/", middleware.JWTAuthMiddleware(), controllers.CreateOrganization)
	adminOrgRoutes.GET("/", middleware.JWTAuthMiddleware(), controllers.GetAllOrganizations)
	adminOrgRoutes.GET("/:orgID", middleware.JWTAuthMiddleware(), controllers.GetOrganization)
	adminOrgRoutes.PUT("/:orgID", middleware.JWTAuthMiddleware(), controllers.UpdateOrganization)
	adminOrgRoutes.DELETE("/:orgID", middleware.JWTAuthMiddleware(), controllers.DeleteOrganization)
}
