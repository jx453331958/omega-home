package handlers

import (
	"net/http"

	"omega-home/models"

	"github.com/gin-gonic/gin"
)

func AdminPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", nil)
}

// Services CRUD
func ListServices(c *gin.Context) {
	var services []models.Service
	models.DB.Order("sort_order asc").Find(&services)
	c.JSON(http.StatusOK, services)
}

func CreateService(c *gin.Context) {
	var s models.Service
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Create(&s)
	c.JSON(http.StatusCreated, s)
}

func UpdateService(c *gin.Context) {
	id := c.Param("id")
	var s models.Service
	if err := models.DB.First(&s, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Save(&s)
	c.JSON(http.StatusOK, s)
}

func DeleteService(c *gin.Context) {
	id := c.Param("id")
	models.DB.Delete(&models.Service{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func ReorderServices(c *gin.Context) {
	var items []struct {
		ID        uint `json:"id"`
		SortOrder int  `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, item := range items {
		models.DB.Model(&models.Service{}).Where("id = ?", item.ID).Update("sort_order", item.SortOrder)
	}
	c.JSON(http.StatusOK, gin.H{"message": "reordered"})
}

// Groups CRUD
func ListGroups(c *gin.Context) {
	var groups []models.Group
	models.DB.Order("sort_order asc").Find(&groups)
	c.JSON(http.StatusOK, groups)
}

func CreateGroup(c *gin.Context) {
	var g models.Group
	if err := c.ShouldBindJSON(&g); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Create(&g)
	c.JSON(http.StatusCreated, g)
}

func UpdateGroup(c *gin.Context) {
	id := c.Param("id")
	var g models.Group
	if err := models.DB.First(&g, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err := c.ShouldBindJSON(&g); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Save(&g)
	c.JSON(http.StatusOK, g)
}

func DeleteGroup(c *gin.Context) {
	id := c.Param("id")
	models.DB.Delete(&models.Group{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func ReorderGroups(c *gin.Context) {
	var items []struct {
		ID        uint `json:"id"`
		SortOrder int  `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, item := range items {
		models.DB.Model(&models.Group{}).Where("id = ?", item.ID).Update("sort_order", item.SortOrder)
	}
	c.JSON(http.StatusOK, gin.H{"message": "reordered"})
}

// Bookmarks CRUD
func ListBookmarks(c *gin.Context) {
	var bookmarks []models.Bookmark
	models.DB.Order("sort_order asc").Find(&bookmarks)
	c.JSON(http.StatusOK, bookmarks)
}

func CreateBookmark(c *gin.Context) {
	var b models.Bookmark
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Create(&b)
	c.JSON(http.StatusCreated, b)
}

func UpdateBookmark(c *gin.Context) {
	id := c.Param("id")
	var b models.Bookmark
	if err := models.DB.First(&b, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.DB.Save(&b)
	c.JSON(http.StatusOK, b)
}

func DeleteBookmark(c *gin.Context) {
	id := c.Param("id")
	models.DB.Delete(&models.Bookmark{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Settings
func GetSettings(c *gin.Context) {
	var settings []models.Setting
	models.DB.Find(&settings)
	settingsMap := make(map[string]string)
	for _, s := range settings {
		settingsMap[s.Key] = s.Value
	}
	c.JSON(http.StatusOK, settingsMap)
}

func UpdateSettings(c *gin.Context) {
	var settingsMap map[string]string
	if err := c.ShouldBindJSON(&settingsMap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for key, value := range settingsMap {
		models.DB.Where("key = ?", key).Assign(models.Setting{Value: value}).FirstOrCreate(&models.Setting{Key: key})
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}
