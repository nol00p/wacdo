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

// CreateCategory adds a new product category.
// Category names must be unique.
//
// @Summary Create a new category
// @Description Create a new product category with the provided details
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body models.Category true "Category details"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]string "Invalid data or category already exists"
// @Security BearerAuth
// @Router /categories [post]
func CreateCategory(c *gin.Context) {
	var category models.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check if the category already exist
	var existingCategory models.Category
	if err := config.DB.Where("name = ?", category.Name).First(&existingCategory).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category already exists"})
		return
	}

	// Create
	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category couldn't be created"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory removes a category from the system.
// A category cannot be deleted if products are still assigned to it — reassign or delete products first.
//
// @Summary Delete a category
// @Description Delete a category by ID (only if not in use by any products)
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]string "Category deleted"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Category not found"
// @Failure 409 {object} map[string]string "Category still in use"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /categories/{id} [delete]
func DeleteCategory(c *gin.Context) {

	// ensure the id of the Category is valid
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Check if category exists
	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Check if any products are using this category
	var count int64
	config.DB.Model(&models.Products{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete category: still in use by products"})
		return
	}

	// Safe to delete
	if err := config.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}

// GetCategories returns all product categories.
//
// @Summary Get all categories
// @Description Retrieve a list of all product categories
// @Tags Categories
// @Produce json
// @Success 200 {array} models.Category
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /categories [get]
func GetCategories(c *gin.Context) {

	var category []models.Category

	if err := config.DB.Find(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// GetCategory returns a single category by ID.
//
// @Summary Get a category by ID
// @Description Retrieve a single category by its ID
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Category not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /categories/{id} [get]
func GetCategory(c *gin.Context) {

	var category models.Category

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := config.DB.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// UpdateCategory modifies an existing category.
// Validates that the new name doesn't conflict with another category.
//
// @Summary Update a category
// @Description Update an existing category by ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body models.Category true "Updated category details"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]string "Invalid ID or data"
// @Failure 404 {object} map[string]string "Category not found"
// @Failure 409 {object} map[string]string "Category name already exists"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /categories/{id} [put]
func UpdateCategory(c *gin.Context) {
	var category models.Category

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	// Check if the category exists
	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	// Bind the update data
	var input models.Category
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check if the new name conflicts with another category
	var existingCategory models.Category
	if err := config.DB.Where("name = ? AND id != ?", input.Name, id).First(&existingCategory).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Category name already exists"})
		return
	}

	// Update the category
	if err := config.DB.Model(&category).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, category)
}
