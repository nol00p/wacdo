package controllers

import (
	"net/http"
	"wacdo/config"
	"wacdo/models"

	"github.com/gin-gonic/gin"
)

func CreateRole(c *gin.Context) {
	var role models.Roles

	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data"})
		return
	}

	// Check if the role already exist
	var existingRole models.Roles
	if err := config.DB.Where("role_name = ?", role.RoleName).First(&existingRole).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role allready exists"})
		return
	}

	// Create
	if err := config.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role couln't not be reated"})
		return
	}

	c.JSON(http.StatusOK, role)
}
