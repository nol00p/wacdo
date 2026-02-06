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

// CreateOption godoc
// @Summary Create a new product option
// @Description Create a new option for a product (e.g. Size, Toppings)
// @Tags Product Options
// @Accept json
// @Produce json
// @Param option body models.ProductOptions true "Option details"
// @Success 200 {object} models.ProductOptions
// @Failure 400 {object} map[string]string "Invalid data or product not found"
// @Security BearerAuth
// @Router /options [post]
func CreateOption(c *gin.Context) {
	var option models.ProductOptions

	if err := c.ShouldBindJSON(&option); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check if the product exists
	var product models.Products
	if err := config.DB.First(&product, option.ProductID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
		return
	}

	// Validate is_unique value
	if option.IsUnique != "single" && option.IsUnique != "multiple" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "is_unique must be 'single' or 'multiple'"})
		return
	}

	// Check if the same option name already exists for this product
	var existing models.ProductOptions
	if err := config.DB.Where("product_id = ? AND name = ?", option.ProductID, option.Name).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Option already exists for this product"})
		return
	}

	if err := config.DB.Create(&option).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Option could not be created"})
		return
	}

	c.JSON(http.StatusOK, option)
}

// DeleteOption godoc
// @Summary Delete a product option
// @Description Delete a product option by ID (also deletes its values)
// @Tags Product Options
// @Produce json
// @Param id path int true "Option ID"
// @Success 200 {object} map[string]string "Option deleted"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Option not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /options/{id} [delete]
func DeleteOption(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var option models.ProductOptions
	if err := config.DB.First(&option, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Option not found"})
		return
	}

	// Delete associated option values first
	if err := config.DB.Where("option_id = ?", id).Delete(&models.OptionValues{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete option values"})
		return
	}

	if err := config.DB.Delete(&option).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete option"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Option deleted"})
}

// GetOptions godoc
// @Summary Get all product options
// @Description Retrieve a list of all product options
// @Tags Product Options
// @Produce json
// @Success 200 {array} models.ProductOptions
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /options [get]
func GetOptions(c *gin.Context) {
	var options []models.ProductOptions

	if err := config.DB.Find(&options).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't get options"})
		return
	}

	c.JSON(http.StatusOK, options)
}

// GetOption godoc
// @Summary Get a product option by ID
// @Description Retrieve a single product option by its ID
// @Tags Product Options
// @Produce json
// @Param id path int true "Option ID"
// @Success 200 {object} models.ProductOptions
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Option not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /options/{id} [get]
func GetOption(c *gin.Context) {
	var option models.ProductOptions

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := config.DB.First(&option, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Option not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, option)
}

// UpdateOption godoc
// @Summary Update a product option
// @Description Update an existing product option by ID
// @Tags Product Options
// @Accept json
// @Produce json
// @Param id path int true "Option ID"
// @Param option body models.ProductOptions true "Updated option details"
// @Success 200 {object} models.ProductOptions
// @Failure 400 {object} map[string]string "Invalid ID or data"
// @Failure 404 {object} map[string]string "Option not found"
// @Failure 409 {object} map[string]string "Option name already exists for this product"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /options/{id} [put]
func UpdateOption(c *gin.Context) {
	var option models.ProductOptions

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := config.DB.First(&option, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Option not found"})
		return
	}

	var input models.ProductOptions
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Validate is_unique value if provided
	if input.IsUnique != "" && input.IsUnique != "single" && input.IsUnique != "multiple" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "is_unique must be 'single' or 'multiple'"})
		return
	}

	// Check if the new name conflicts with another option for the same product
	productID := option.ProductID
	if input.ProductID != 0 {
		productID = input.ProductID
	}
	var existing models.ProductOptions
	if err := config.DB.Where("product_id = ? AND name = ? AND id != ?", productID, input.Name, id).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Option name already exists for this product"})
		return
	}

	if err := config.DB.Model(&option).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update option"})
		return
	}

	c.JSON(http.StatusOK, option)
}

// GetOptionsByProduct godoc
// @Summary Get options by product
// @Description Retrieve all options belonging to a specific product
// @Tags Product Options
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {array} models.ProductOptions
// @Failure 400 {object} map[string]string "Invalid product ID"
// @Failure 404 {object} map[string]string "Product not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /options/product/{product_id} [get]
func GetOptionsByProduct(c *gin.Context) {
	idParam := c.Param("product_id")
	productID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Verify product exists
	var product models.Products
	if err := config.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var options []models.ProductOptions
	if err := config.DB.Where("product_id = ?", productID).Find(&options).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't get options"})
		return
	}

	c.JSON(http.StatusOK, options)
}
