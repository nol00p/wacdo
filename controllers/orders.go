package controllers

import (
	"errors"
	"net/http"
	"slices"
	"strconv"
	"time"
	"wacdo/config"
	"wacdo/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderItemOptionInput struct {
	OptionValueID uint `json:"option_value_id"`
}

type OrderItemInput struct {
	ProductID *uint                  `json:"product_id"`
	MenuID    *uint                  `json:"menu_id"`
	Quantity  uint                   `json:"quantity"`
	Options   []OrderItemOptionInput `json:"options"`
}

type OrderInput struct {
	CustomerID    *uint            `json:"customer_id"`
	OrderType     string           `json:"order_type"`
	Notes         string           `json:"notes"`
	ScheduledTime *time.Time       `json:"scheduled_time"`
	Items         []OrderItemInput `json:"items"`
}

type StatusInput struct {
	Status string `json:"status"`
}

// orderPreloads applies the standard set of preloads for order queries.
func orderPreloads(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Customer").
		Preload("CreatedBy").
		Preload("OrderItems.Product").
		Preload("OrderItems.Menu").
		Preload("OrderItems.OrderItemOptions.OptionValue")
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create an order with items and options. Prices are computed server-side.
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body OrderInput true "Order details"
// @Success 201 {object} models.Order
// @Failure 400 {object} map[string]string "Invalid data"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /orders [post]
func CreateOrder(c *gin.Context) {
	var input OrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Validate order type
	if input.OrderType != "counter" && input.OrderType != "phone" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order type must be 'counter' or 'phone'"})
		return
	}

	// Validate customer exists if provided
	if input.CustomerID != nil {
		var customer models.Customer
		if err := config.DB.First(&customer, *input.CustomerID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Customer not found"})
			return
		}
	}

	// Require at least one item
	if len(input.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order must have at least one item"})
		return
	}

	var createdOrder models.Order

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// Create order record
		// TODO: CreatedByID should come from JWT context once auth middleware is active
		order := models.Order{
			CustomerID:    input.CustomerID,
			CreatedByID:   1,
			OrderType:     input.OrderType,
			Status:        "pending",
			Notes:         input.Notes,
			ScheduledTime: input.ScheduledTime,
			TotalPrice:    0,
		}

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		var totalPrice float64

		for _, itemInput := range input.Items {
			// Validate exactly one of ProductID or MenuID
			hasProduct := itemInput.ProductID != nil
			hasMenu := itemInput.MenuID != nil
			if hasProduct == hasMenu {
				return errors.New("each item must have exactly one of product_id or menu_id")
			}

			if itemInput.Quantity == 0 {
				return errors.New("item quantity must be at least 1")
			}

			var unitPrice float64

			if hasProduct {
				var product models.Products
				if err := tx.First(&product, *itemInput.ProductID).Error; err != nil {
					return errors.New("product not found")
				}
				if !product.IsAvailable {
					return errors.New("product '" + product.Name + "' is not available")
				}
				unitPrice = product.Price
			} else {
				var menu models.Menu
				if err := tx.First(&menu, *itemInput.MenuID).Error; err != nil {
					return errors.New("menu not found")
				}
				if !menu.IsAvailable {
					return errors.New("menu '" + menu.Name + "' is not available")
				}
				unitPrice = menu.Price
			}

			// Process options and compute option price sum
			var optionPriceSum float64
			var optionRecords []models.OrderItemOption

			for _, optInput := range itemInput.Options {
				var optionValue models.OptionValues
				if err := tx.Preload("Option").First(&optionValue, optInput.OptionValueID).Error; err != nil {
					return errors.New("option value not found")
				}

				// Verify the option belongs to the product (only for product items)
				if hasProduct {
					if optionValue.Option.ProductID != *itemInput.ProductID {
						return errors.New("option value does not belong to the selected product")
					}
				}

				optionPriceSum += optionValue.OptionPrice
				optionRecords = append(optionRecords, models.OrderItemOption{
					OptionValueID: optInput.OptionValueID,
					PriceApplied:  optionValue.OptionPrice,
				})
			}

			itemTotal := (unitPrice + optionPriceSum) * float64(itemInput.Quantity)

			orderItem := models.OrderItem{
				OrderID:   order.ID,
				ProductID: itemInput.ProductID,
				MenuID:    itemInput.MenuID,
				Quantity:  itemInput.Quantity,
				UnitPrice: unitPrice,
				ItemTotal: itemTotal,
			}

			if err := tx.Create(&orderItem).Error; err != nil {
				return err
			}

			// Create option records
			for i := range optionRecords {
				optionRecords[i].OrderItemID = orderItem.ID
				if err := tx.Create(&optionRecords[i]).Error; err != nil {
					return err
				}
			}

			totalPrice += itemTotal
		}

		// Update total price
		if err := tx.Model(&order).Update("total_price", totalPrice).Error; err != nil {
			return err
		}

		createdOrder = order
		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Reload with preloads
	var result models.Order
	if err := orderPreloads(config.DB).First(&result, createdOrder.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load created order"})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetOrders godoc
// @Summary Get all orders
// @Description Retrieve all orders with optional status filter
// @Tags Orders
// @Produce json
// @Param status query string false "Filter by status"
// @Success 200 {array} models.Order
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /orders [get]
func GetOrders(c *gin.Context) {
	var orders []models.Order

	query := orderPreloads(config.DB)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrder godoc
// @Summary Get an order by ID
// @Description Retrieve a single order with all its items and options
// @Tags Orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Order not found"
// @Security BearerAuth
// @Router /orders/{id} [get]
func GetOrder(c *gin.Context) {
	var order models.Order

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := orderPreloads(config.DB).First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// UpdateOrderStatus godoc
// @Summary Update order status
// @Description Update the status of an order with valid transition enforcement
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param status body StatusInput true "New status"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string "Invalid transition"
// @Failure 404 {object} map[string]string "Order not found"
// @Security BearerAuth
// @Router /orders/{id}/status [patch]
func UpdateOrderStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var order models.Order
	if err := config.DB.First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var input StatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Enforce valid transitions
	validTransitions := map[string][]string{
		"pending":   {"preparing", "cancelled"},
		"preparing": {"prepared"},
		"prepared":  {"delivered"},
	}

	allowed, exists := validTransitions[order.Status]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No transitions allowed from status '" + order.Status + "'"})
		return
	}

	if !slices.Contains(allowed, input.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transition from '" + order.Status + "' to '" + input.Status + "'"})
		return
	}

	if err := config.DB.Model(&order).Update("status", input.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	// Reload with preloads
	var result models.Order
	if err := orderPreloads(config.DB).First(&result, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load order"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// CancelOrder godoc
// @Summary Cancel an order
// @Description Cancel an order (only if status is pending)
// @Tags Orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string "Cannot cancel"
// @Failure 404 {object} map[string]string "Order not found"
// @Security BearerAuth
// @Router /orders/{id}/cancel [patch]
func CancelOrder(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var order models.Order
	if err := config.DB.First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if order.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only pending orders can be cancelled"})
		return
	}

	if err := config.DB.Model(&order).Update("status", "cancelled").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order"})
		return
	}

	// Reload with preloads
	var result models.Order
	if err := orderPreloads(config.DB).First(&result, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load order"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetOrdersByCustomer godoc
// @Summary Get orders by customer
// @Description Retrieve all orders for a specific customer
// @Tags Orders
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {array} models.Order
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Customer not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Security BearerAuth
// @Router /customers/{id}/orders [get]
func GetOrdersByCustomer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Check customer exists
	var customer models.Customer
	if err := config.DB.First(&customer, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var orders []models.Order
	if err := orderPreloads(config.DB).Where("customer_id = ?", id).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
