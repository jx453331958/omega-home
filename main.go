package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"omega-home/config"
	"omega-home/handlers"
	"omega-home/middleware"
	"omega-home/models"
	"omega-home/services"

	"github.com/gin-gonic/gin"
)

//go:embed templates/*
var templatesFS embed.FS

func main() {
	cfg := config.Load()

	models.InitDB(cfg.DatabaseURL)
	services.StartChecker(cfg.CheckInterval)

	r := gin.Default()

	tmpl := template.Must(template.ParseFS(templatesFS, "templates/*.html"))
	r.SetHTMLTemplate(tmpl)

	r.Static("/static", "./static")

	// Public routes
	r.GET("/", handlers.PortalPage)
	r.GET("/api/config", handlers.GetConfig)
	r.GET("/api/status", handlers.GetStatus)

	// Admin page (no auth needed to load page, auth is for API)
	r.GET("/admin", handlers.AdminPage)

	// Auth
	authHandler := &handlers.AuthHandler{Password: cfg.AdminPassword, Secret: cfg.SecretKey}
	r.POST("/api/admin/login", authHandler.Login)

	// Protected admin API
	admin := r.Group("/api/admin")
	admin.Use(middleware.JWTAuth(cfg.SecretKey))
	{
		admin.GET("/services", handlers.ListServices)
		admin.POST("/services", handlers.CreateService)
		admin.PUT("/services/:id", handlers.UpdateService)
		admin.DELETE("/services/:id", handlers.DeleteService)
		admin.PUT("/services/reorder", handlers.ReorderServices)

		admin.GET("/groups", handlers.ListGroups)
		admin.POST("/groups", handlers.CreateGroup)
		admin.PUT("/groups/:id", handlers.UpdateGroup)
		admin.DELETE("/groups/:id", handlers.DeleteGroup)
		admin.PUT("/groups/reorder", handlers.ReorderGroups)

		admin.GET("/bookmarks", handlers.ListBookmarks)
		admin.POST("/bookmarks", handlers.CreateBookmark)
		admin.PUT("/bookmarks/:id", handlers.UpdateBookmark)
		admin.DELETE("/bookmarks/:id", handlers.DeleteBookmark)

		admin.GET("/settings", handlers.GetSettings)
		admin.PUT("/settings", handlers.UpdateSettings)

		admin.POST("/upload", handlers.UploadImage)
	}

	log.Printf("Omega Home starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
