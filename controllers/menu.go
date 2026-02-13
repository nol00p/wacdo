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

// CreateMenu godoc
// @Summary Create a new menu
// @Description Create a new menu with the provided details
// @Tags Menus
// @Accept json
// @Produce json
// @Param menu body models.Menu true "Menu details"
// @Success 200 {object} models.Menu
// @Failure 400 {object} map[string]string "Invalid data or menu already exists"
// @Security BearerAuth
// @Router /menus [post]
func CreateMenu(c *gin.Context) {
	var menu models.Menu

	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check if the menu already exists
	var existingMenu models.Menu
	if err := config.DB.Where("name = ?", menu.Name).First(&existingMenu).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Menu already exists"})
		return
	}

	if err := config.DB.Create(&menu).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Menu could not be created"})
		return
	}

	c.JSON(http.StatusOK, menu)
}

// GetMenus godoc
// @Summary Get all menus
// @Description Retrieve a list of all menus with their products
// @Tags Menus
// @Produce json
// @Success 200 {array} models.Menu
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /menus [get]
func GetMenus(c *gin.Context) {
	var menus []models.Menu

	// pre load menu product to get the list of products in the menu.
	if err := config.DB.Preload("MenuProducts").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't get menus"})
		return
	}

	c.JSON(http.StatusOK, menus)
}

// GetMenu godoc
// @Summary Get a menu by ID
// @Description Retrieve a single menu by its ID with associated products
// @Tags Menus
// @Produce json
// @Param id path int true "Menu ID"
// @Success 200 {object} models.Menu
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Menu not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /menus/{id} [get]
func GetMenu(c *gin.Context) {
	var menu models.Menu

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	//Check if the Id exists
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// get menu it's associated products
	if err := config.DB.Preload("MenuProducts").First(&menu, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, menu)
}

// UpdateMenu godoc
// @Summary Update a menu
// @Description Update an existing menu by ID
// @Tags Menus
// @Accept json
// @Produce json
// @Param id path int true "Menu ID"
// @Param menu body models.Menu true "Updated menu details"
// @Success 200 {object} models.Menu
// @Failure 400 {object} map[string]string "Invalid ID or data"
// @Failure 404 {object} map[string]string "Menu not found"
// @Failure 409 {object} map[string]string "Menu name already exists"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /menus/{id} [put]
func UpdateMenu(c *gin.Context) {
	var menu models.Menu

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := config.DB.First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
		return
	}

	var input models.Menu
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	//Check if the menu name allready exists
	var existingMenu models.Menu
	if err := config.DB.Where("name = ? AND id != ?", input.Name, id).First(&existingMenu).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Menu name already exists"})
		return
	}

	if err := config.DB.Model(&menu).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update menu"})
		return
	}

	config.DB.Preload("MenuProducts").First(&menu, id)

	c.JSON(http.StatusOK, menu)
}

// DeleteMenu godoc
// @Summary Delete a menu
// @Description Delete a menu by ID
// @Tags Menus
// @Produce json
// @Param id path int true "Menu ID"
// @Success 200 {object} map[string]string "Menu deleted"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Menu not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /menus/{id} [delete]
func DeleteMenu(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var menu models.Menu
	if err := config.DB.First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
		return
	}

	if err := config.DB.Delete(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete menu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu deleted"})
}

// ToggleMenuAvailability godoc
// @Summary Toggle menu availability
// @Description Quick toggle for is_available field
// @Tags Menus
// @Produce json
// @Param id path int true "Menu ID"
// @Success 200 {object} models.Menu
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Menu not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /menus/{id}/availability [patch]
func ToggleMenuAvailability(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var menu models.Menu
	if err := config.DB.First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
		return
	}

	menu.IsAvailable = !menu.IsAvailable

	if err := config.DB.Save(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update availability"})
		return
	}

	config.DB.Preload("MenuProducts").First(&menu, id)

	c.JSON(http.StatusOK, menu)
}

// AddProductToMenu godoc
// @Summary Add a product to a menu
// @Description Add a product to a specific menu
// @Tags Menus
// @Accept json
// @Produce json
// @Param id path int true "Menu ID"
// @Param menu_product body models.MenuProduct true "Menu product details"
// @Success 200 {object} models.MenuProduct
// @Failure 400 {object} map[string]string "Invalid ID or data"
// @Failure 404 {object} map[string]string "Menu or product not found"
// @Failure 409 {object} map[string]string "Product already in menu"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /menus/{id}/products [post]
func AddProductToMenu(c *gin.Context) {
	idParam := c.Param("id")
	menuID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Check for menu
	var menu models.Menu
	if err := config.DB.First(&menu, menuID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
		return
	}

	// validate input data
	var input models.MenuProduct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check if the product exists
	var product models.Products
	if err := config.DB.First(&product, input.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Check if the product is already in the menu
	var existing models.MenuProduct
	if err := config.DB.Where("menu_id = ? AND product_id = ?", menuID, input.ProductID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Product already in menu"})
		return
	}

	input.MenuID = uint(menuID)

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product to menu"})
		return
	}

	c.JSON(http.StatusOK, input)
}

// GetMenuProducts godoc
// @Summary Get products in a menu
// @Description Retrieve all products belonging to a specific menu
// @Tags Menus
// @Produce json
// @Param id path int true "Menu ID"
// @Success 200 {array} models.MenuProduct
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Menu not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /menus/{id}/products [get]
func GetMenuProducts(c *gin.Context) {
	idParam := c.Param("id")
	menuID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var menu models.Menu
	if err := config.DB.First(&menu, menuID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
		return
	}

	var menuProducts []models.MenuProduct
	if err := config.DB.Where("menu_id = ?", menuID).Find(&menuProducts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't get menu products"})
		return
	}

	c.JSON(http.StatusOK, menuProducts)
}

// RemoveProductFromMenu godoc
// @Summary Remove a product from a menu
// @Description Remove a product from a menu by menu product ID
// @Tags Menus
// @Produce json
// @Param id path int true "Menu Product ID"
// @Success 200 {object} map[string]string "Product removed from menu"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Menu product not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /menus/products/{id} [delete]
func RemoveProductFromMenu(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var menuProduct models.MenuProduct
	if err := config.DB.First(&menuProduct, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu product not found"})
		return
	}

	if err := config.DB.Delete(&menuProduct).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove product from menu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product removed from menu"})
}
