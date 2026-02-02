package models

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(databaseURL string) {
	var err error
	var dialector gorm.Dialector

	if strings.HasPrefix(databaseURL, "postgres://") || strings.HasPrefix(databaseURL, "postgresql://") {
		dialector = postgres.Open(databaseURL)
	} else {
		// sqlite:///data/omega.db -> /data/omega.db
		dbPath := strings.TrimPrefix(databaseURL, "sqlite://")
		dir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create database directory: %v", err)
		}
		dialector = sqlite.Open(dbPath)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := DB.AutoMigrate(&Group{}, &Service{}, &Setting{}, &Bookmark{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	seed()
}

func seed() {
	// Only seed if no settings exist
	var count int64
	DB.Model(&Setting{}).Count(&count)
	if count > 0 {
		return
	}

	// Default settings
	settings := []Setting{
		{Key: "title", Value: "Omega Home"},
		{Key: "subtitle", Value: "Your Personal Portal"},
		{Key: "theme", Value: "indigo-purple"},
		{Key: "search_enabled", Value: "true"},
		{Key: "clock_enabled", Value: "true"},
		{Key: "greeting_enabled", Value: "true"},
		{Key: "weather_enabled", Value: "false"},
		{Key: "footer_text", Value: "Powered by Omega Home"},
		{Key: "logo", Value: ""},
		{Key: "background_image", Value: ""},
	}
	DB.Create(&settings)

	// Example group
	group := Group{Name: "常用工具", Icon: "🛠️", SortOrder: 0, Columns: 3}
	DB.Create(&group)

	// Example services
	services := []Service{
		{Name: "Google", URL: "https://www.google.com", Icon: "🔍", Description: "Search engine", GroupID: group.ID, SortOrder: 0},
		{Name: "GitHub", URL: "https://github.com", Icon: "🐙", Description: "Code hosting", GroupID: group.ID, SortOrder: 1},
		{Name: "Wikipedia", URL: "https://wikipedia.org", Icon: "📚", Description: "Encyclopedia", GroupID: group.ID, SortOrder: 2},
	}
	DB.Create(&services)
}
