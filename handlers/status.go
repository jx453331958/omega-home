package handlers

import (
	"net/http"

	"omega-home/services"

	"github.com/gin-gonic/gin"
)

func GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, services.GetAllStatus())
}
