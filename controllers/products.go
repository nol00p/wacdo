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

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the provided details
// @Tags Products
// @Accept json
// @Produce json
// @Param product body models.Products true "Product details"
// @Success 200 {object} models.Products
// @Failure 400 {object} map[string]string "Invalid data or product already exists"
// @Security BearerAuth
// @Router /products [post]
func CreateProduct(c *gin.Context) {
	var product models.Products

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check if the category exists
	var category models.Category
	if err := config.DB.First(&category, product.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
		return
	}

	// Check if the product already exists
	var existingProduct models.Products
	if err := config.DB.Where("name = ?", product.Name).First(&existingProduct).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product already exists"})
		return
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product could not be created"})
		return
	}

	// Load the category for response
	config.DB.Preload("Category").First(&product, product.ID)

	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]string "Product deleted"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Product not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	// Check if the id exists
	var product models.Products
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	// Delete the product from DB
	if err := config.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

// GetProducts godoc
// @Summary Get all products
// @Description Retrieve a list of all products with their categories
// @Tags Products
// @Produce json
// @Success 200 {array} models.Products
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /products [get]
func GetProducts(c *gin.Context) {
	var products []models.Products

	// Preload the Category to get the details with the product
	if err := config.DB.Preload("Category").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't get products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProduct godoc
// @Summary Get a product by ID
// @Description Retrieve a single product by its ID
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Products
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Product not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /products/{id} [get]
func GetProduct(c *gin.Context) {
	var product models.Products

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	// Preload the Category to get the details with the product
	if err := config.DB.Preload("Category").First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update an existing product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.Products true "Updated product details"
// @Success 200 {object} models.Products
// @Failure 400 {object} map[string]string "Invalid ID or data"
// @Failure 404 {object} map[string]string "Product not found"
// @Failure 409 {object} map[string]string "Product name already exists"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /products/{id} [put]
func UpdateProduct(c *gin.Context) {
	var product models.Products

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var input models.Products
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check if the new name conflicts with another product
	var existingProduct models.Products
	if err := config.DB.Where("name = ? AND id != ?", input.Name, id).First(&existingProduct).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Product name already exists"})
		return
	}

	// If category is being changed, verify it exists
	if input.CategoryID != 0 && input.CategoryID != product.CategoryID {
		var category models.Category
		if err := config.DB.First(&category, input.CategoryID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
			return
		}
	}

	if err := config.DB.Model(&product).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	// Load the category for response
	config.DB.Preload("Category").First(&product, id)

	c.JSON(http.StatusOK, product)
}

// GetProductsByCategory godoc
// @Summary Get products by category
// @Description Retrieve all products belonging to a specific category
// @Tags Products
// @Produce json
// @Param category_id path int true "Category ID"
// @Success 200 {array} models.Products
// @Failure 400 {object} map[string]string "Invalid category ID"
// @Failure 404 {object} map[string]string "Category not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /products/category/{category_id} [get]
func GetProductsByCategory(c *gin.Context) {
	idParam := c.Param("category_id")
	categoryID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Verify category exists
	var category models.Category
	if err := config.DB.First(&category, categoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var products []models.Products
	if err := config.DB.Preload("Category").Where("category_id = ?", categoryID).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't get products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// ToggleProductAvailability godoc
// @Summary Toggle product availability
// @Description Quick toggle for is_available field
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Products
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Product not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /products/{id}/availability [patch]
func ToggleProductAvailability(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var product models.Products
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	product.IsAvailable = !product.IsAvailable

	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update availability"})
		return
	}

	config.DB.Preload("Category").First(&product, id)

	c.JSON(http.StatusOK, product)
}

// UpdateProductStock godoc
// @Summary Update product stock
// @Description Update the stock quantity for a product
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param stock body object true "Stock update" example({"stock_quantity": 100})
// @Success 200 {object} models.Products
// @Failure 400 {object} map[string]string "Invalid ID or data"
// @Failure 404 {object} map[string]string "Product not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /products/{id}/stock [patch]
func UpdateProductStock(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var product models.Products
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var input struct {
		StockQuantity uint `json:"stock_quantity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	product.StockQuantity = input.StockQuantity

	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
		return
	}

	config.DB.Preload("Category").First(&product, id)

	c.JSON(http.StatusOK, product)
}
