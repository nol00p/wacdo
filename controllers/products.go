package controllers

import (
	"net/http"
	"strconv"
	"wacdo/config"
	"wacdo/models"

	"github.com/gin-gonic/gin"
)

func CreateCategory(c *gin.Context) {
	var category models.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data"})
		return
	}

	// Check if the category already exist
	var existingCategory models.Category
	if err := config.DB.Where("category_name = ?", category.Name).First(&existingCategory).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category already exists"})
		return
	}

	// Create
	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category couln't not be reated"})
		return
	}

	c.JSON(http.StatusOK, category)
}

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
		c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete category: still in use by product"})
		return
	}

	// Safe to delete
	if err := config.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}

func GetCategories(c *gin.Context) {

	var category []models.Category

	if err := config.DB.Find(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't get categories"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// CreateProduct
// DeleteProduct
// ListProducts

//CreateOption
//DeleteOption
//ListOptions

//CreateOptionValue
//DeleteOptionValue
//ListOptionvalues
