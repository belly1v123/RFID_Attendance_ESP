package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ronishg27/rfid_attendance/config"
	"github.com/ronishg27/rfid_attendance/handlers"
	"github.com/ronishg27/rfid_attendance/internal/db"
	"github.com/ronishg27/rfid_attendance/internal/middleware"
	"github.com/ronishg27/rfid_attendance/utils"
)

func main() {

	config.InitDB()

	r := gin.Default()
	queries := db.NewQueries(config.DB)

	err := config.SeedSuperAdmin(queries.Public)
	utils.HandleError(err, false)

	r.Use(func(c *gin.Context) {
		c.Set("db", queries)
		c.Next()
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/api/login", handlers.SuperAdminLoginHandler(queries))

	r.POST("/api/organizations", middleware.JWTAuthMiddleware(), handlers.RegisterOrganizationHandler)

	// routes.SetupRoutes(r)

	err = r.Run(":3000")
	utils.HandleError(err, true)
}
