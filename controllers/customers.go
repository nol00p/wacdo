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

// CreateCustomer registers a new customer.
// If a phone number is provided, it must be unique to prevent duplicate customer records.
//
// @Summary Create a new customer
// @Description Create a new customer with the provided details
// @Tags Customers
// @Accept json
// @Produce json
// @Param customer body models.Customer true "Customer details"
// @Success 200 {object} models.Customer
// @Failure 400 {object} map[string]string "Invalid data"
// @Security BearerAuth
// @Router /customers [post]
func CreateCustomer(c *gin.Context) {
	var customer models.Customer

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check if a customer with the same phone already exists (if phone provided)
	if customer.Phone != "" {
		var existing models.Customer
		if err := config.DB.Where("phone = ?", customer.Phone).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Customer with this phone number already exists", "existing_id": existing.ID})
			return
		}
	}

	if err := config.DB.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// GetCustomers returns all customers. Supports GDPR right of consultation.
//
// @Summary Get all customers
// @Description Retrieve a list of all customers
// @Tags Customers
// @Produce json
// @Success 200 {array} models.Customer
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /customers [get]
func GetCustomers(c *gin.Context) {
	var customers []models.Customer

	if err := config.DB.Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customers"})
		return
	}

	c.JSON(http.StatusOK, customers)
}

// GetCustomer returns a single customer by ID.
//
// @Summary Get a customer by ID
// @Description Retrieve a single customer by their ID
// @Tags Customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {object} models.Customer
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Customer not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /customers/{id} [get]
func GetCustomer(c *gin.Context) {
	var customer models.Customer

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := config.DB.First(&customer, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// UpdateCustomer modifies an existing customer's details.
// Validates that the new phone number doesn't conflict with another customer.
// Supports GDPR right of modification.
//
// @Summary Update a customer
// @Description Update an existing customer by ID
// @Tags Customers
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Param customer body models.Customer true "Updated customer details"
// @Success 200 {object} models.Customer
// @Failure 400 {object} map[string]string "Invalid ID or data"
// @Failure 404 {object} map[string]string "Customer not found"
// @Failure 409 {object} map[string]string "Phone number already in use"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /customers/{id} [put]
func UpdateCustomer(c *gin.Context) {
	var customer models.Customer

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Check if the customer exists
	if err := config.DB.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	// Bind the update data
	var input models.Customer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check if the new phone conflicts with another customer
	if input.Phone != "" {
		var existing models.Customer
		if err := config.DB.Where("phone = ? AND id != ?", input.Phone, id).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Phone number already in use by another customer"})
			return
		}
	}

	if err := config.DB.Model(&customer).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// DeleteCustomer permanently removes a customer record.
// Supports GDPR right of deletion (droit de suppression).
//
// @Summary Delete a customer
// @Description Delete a customer by ID
// @Tags Customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {object} map[string]string "Customer deleted"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Customer not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /customers/{id} [delete]
func DeleteCustomer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var customer models.Customer
	if err := config.DB.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	if err := config.DB.Delete(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted"})
}
