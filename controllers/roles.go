package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"wacdo/config"
	"wacdo/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateRole godoc
// @Summary Create a new role
// @Description Create a new role with the provided details
// @Tags Roles
// @Accept json
// @Produce json
// @Param role body models.Roles true "Role details"
// @Success 200 {object} models.Roles
// @Failure 400 {object} map[string]string "Invalid data or role already exists"
// @Security BearerAuth
// @Router /roles [post]
func CreateRole(c *gin.Context) {
	var role models.Roles

	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data"})
		return
	}

	// Check if the role already exist
	var existingRole models.Roles
	if err := config.DB.Where("role_name = ?", role.RoleName).First(&existingRole).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role already exists"})
		return
	}

	// Create
	if err := config.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role couln't not be reated"})
		return
	}

	c.JSON(http.StatusOK, role)
}

// DeleteRole godoc
// @Summary Delete a role
// @Description Delete a role by ID (only if not in use by any users)
// @Tags Roles
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} map[string]string "Role deleted"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Role not found"
// @Failure 409 {object} map[string]string "Role still in use"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /roles/{id} [delete]
func DeleteRole(c *gin.Context) {

	// ensure the id of the role is valid
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Check if role exists
	var role models.Roles
	if err := config.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	// Check if any users are using this role
	var count int64
	config.DB.Model(&models.Users{}).Where("roles_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete role: still in use by users"})
		return
	}

	// Safe to delete
	if err := config.DB.Delete(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted"})
}

// GetRoles godoc
// @Summary Get all roles
// @Description Retrieve a list of all roles
// @Tags Roles
// @Produce json
// @Success 200 {array} models.Roles
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /roles [get]
func GetRoles(c *gin.Context) {

	var roles []models.Roles

	if err := config.DB.Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't get Role data"})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// GetRole godoc
// @Summary Get a role by ID
// @Description Retrieve a single role by its ID
// @Tags Roles
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} models.Roles
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Role not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /roles/{id} [get]
func GetRole(c *gin.Context) {

	var role models.Roles

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := config.DB.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role can't be found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, role)
}
