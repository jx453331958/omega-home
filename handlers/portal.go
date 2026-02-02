package handlers

import (
	"net/http"

	"omega-home/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PortalPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func GetConfig(c *gin.Context) {
	var groups []models.Group
	models.DB.Order("sort_order asc").Preload("Services", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order asc")
	}).Find(&groups)

	var settings []models.Setting
	models.DB.Find(&settings)
	settingsMap := make(map[string]string)
	for _, s := range settings {
		settingsMap[s.Key] = s.Value
	}

	var bookmarks []models.Bookmark
	models.DB.Order("sort_order asc").Find(&bookmarks)

	c.JSON(http.StatusOK, gin.H{
		"groups":    groups,
		"settings":  settingsMap,
		"bookmarks": bookmarks,
	})
}
