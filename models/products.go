package models

import (
	"time"
)

// Category groups products for display and filtering (e.g. "Burgers", "Drinks", "Desserts").
type Category struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `json:"name"`                // Unique display name
	Description  string    `json:"description"`         // Short description for the kiosk UI
	DisplayOrder uint      `json:"display_order"`       // Controls the display order in the frontend
	ImageURL     string    `json:"image_url"`           // URL to the category image
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Products represents a single orderable item (e.g. "Big Mac", "Coca-Cola").
type Products struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	CategoryID      uint      `json:"category_id"`                         // FK to Category
	Category        Category  `gorm:"foreignKey:CategoryID" json:"category"` // Preloaded category
	Name            string    `json:"name"`                                // Unique product name
	Description     string    `json:"description"`
	Price           float64   `json:"price"`                               // Unit price in euros, used for order price calculation
	StockQuantity   uint      `json:"stock_quantity"`                      // Available stock count
	IsAvailable     bool      `json:"is_available"`                        // Unavailable products cannot be ordered
	ImageURL        string    `json:"image_url"`                           // URL to the product image
	PreparationTime uint      `json:"preparation_time"`                    // Estimated prep time in minutes
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ProductOptions defines a customization group for a product (e.g. "Size", "Toppings").
type ProductOptions struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	ProductID  uint     `json:"product_id"`                  // FK to Products
	Product    Products `gorm:"foreignKey:ProductID" json:"-"`
	Name       string   `json:"name"`                        // Option group name (e.g. "Size")
	IsUnique   string   `json:"is_unique"`                   // "single" = pick one, "multiple" = pick many
	IsRequired bool     `json:"is_required"`                 // Whether the customer must select a value
}

// OptionValues represents one selectable value within a ProductOption (e.g. "Large" for "Size").
type OptionValues struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	OptionID    uint           `json:"option_id"`                    // FK to ProductOptions
	Option      ProductOptions `gorm:"foreignKey:OptionID" json:"-"`
	Value       string         `json:"value"`                        // Display label (e.g. "Large")
	OptionPrice float64        `json:"option_price"`                 // Additional cost added to the product price
}
