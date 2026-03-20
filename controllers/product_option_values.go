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

// CreateOptionValue creates one or more values for a product option in batch.
// Each value must be unique within its option. The option must exist.
// Example: for a "Size" option, create values ["Small", "Medium", "Large"].
//
// @Summary Create option values
// @Description Create one or more values for a product option (e.g. "Small", "Medium", "Large" for Size)
// @Tags Option Values
// @Accept json
// @Produce json
// @Param option_id path int true "Option ID"
// @Param values body []models.OptionValues true "Array of option value details"
// @Success 200 {array} models.OptionValues
// @Failure 400 {object} map[string]string "Invalid data or option not found"
// @Security BearerAuth
// @Router /options/{option_id}/values [post]
func CreateOptionValue(c *gin.Context) {
	optionIDParam := c.Param("id")
	optionID, err := strconv.Atoi(optionIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid option ID"})
		return
	}

	// Check if the option exists
	var option models.ProductOptions
	if err := config.DB.First(&option, optionID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Option not found"})
		return
	}

	var optionValues []models.OptionValues
	if err := c.ShouldBindJSON(&optionValues); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data, expected an array of values"})
		return
	}

	// Validate all values before inserting
	for i := range optionValues {
		optionValues[i].OptionID = uint(optionID)

		var existing models.OptionValues
		if err := config.DB.Where("option_id = ? AND value = ?", optionValues[i].OptionID, optionValues[i].Value).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Value '" + optionValues[i].Value + "' already exists for this option"})
			return
		}
	}

	if err := config.DB.Create(&optionValues).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Option values could not be created"})
		return
	}

	c.JSON(http.StatusOK, optionValues)
}

// DeleteOptionValue removes a single option value.
//
// @Summary Delete an option value
// @Description Delete an option value by ID
// @Tags Option Values
// @Produce json
// @Param id path int true "Option Value ID"
// @Success 200 {object} map[string]string "Option value deleted"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Option value not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /options/values/{id} [delete]
func DeleteOptionValue(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var optionValue models.OptionValues
	if err := config.DB.First(&optionValue, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Option value not found"})
		return
	}

	if err := config.DB.Delete(&optionValue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete option value"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Option value deleted"})
}

// GetOptionValue returns a single option value by ID.
//
// @Summary Get an option value by ID
// @Description Retrieve a single option value by its ID
// @Tags Option Values
// @Produce json
// @Param id path int true "Option Value ID"
// @Success 200 {object} models.OptionValues
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Option value not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /options/values/{id} [get]
func GetOptionValue(c *gin.Context) {
	var optionValue models.OptionValues

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := config.DB.First(&optionValue, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Option value not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, optionValue)
}

// UpdateOptionValue modifies an existing option value.
// Validates that the new value doesn't conflict with another value for the same option.
//
// @Summary Update an option value
// @Description Update an existing option value by ID
// @Tags Option Values
// @Accept json
// @Produce json
// @Param id path int true "Option Value ID"
// @Param value body models.OptionValues true "Updated option value details"
// @Success 200 {object} models.OptionValues
// @Failure 400 {object} map[string]string "Invalid ID or data"
// @Failure 404 {object} map[string]string "Option value not found"
// @Failure 409 {object} map[string]string "Value already exists for this option"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /options/values/{id} [put]
func UpdateOptionValue(c *gin.Context) {
	var optionValue models.OptionValues

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := config.DB.First(&optionValue, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Option value not found"})
		return
	}

	var input models.OptionValues
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check if the new value conflicts with another value for the same option
	optionID := optionValue.OptionID
	if input.OptionID != 0 {
		optionID = input.OptionID
	}
	var existing models.OptionValues
	if err := config.DB.Where("option_id = ? AND value = ? AND id != ?", optionID, input.Value, id).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Value already exists for this option"})
		return
	}

	if err := config.DB.Model(&optionValue).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update option value"})
		return
	}

	c.JSON(http.StatusOK, optionValue)
}

// GetValuesByOption returns all values for a given option.
// The option must exist.
//
// @Summary Get values by option
// @Description Retrieve all values belonging to a specific option
// @Tags Option Values
// @Produce json
// @Param option_id path int true "Option ID"
// @Success 200 {array} models.OptionValues
// @Failure 400 {object} map[string]string "Invalid option ID"
// @Failure 404 {object} map[string]string "Option not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /options/{option_id}/values [get]
func GetValuesByOption(c *gin.Context) {
	idParam := c.Param("id")
	optionID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid option ID"})
		return
	}

	// Verify option exists
	var option models.ProductOptions
	if err := config.DB.First(&option, optionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Option not found"})
		return
	}

	var values []models.OptionValues
	if err := config.DB.Where("option_id = ?", optionID).Find(&values).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve option values"})
		return
	}

	c.JSON(http.StatusOK, values)
}
