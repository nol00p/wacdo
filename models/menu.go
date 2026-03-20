package models

import (
	"time"
)

// Menu represents a combo meal that bundles multiple products at a fixed price (e.g. "Big Mac Menu").
// Menus are orderable items just like products, and their price is independent of the individual product prices.
type Menu struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	Name         string        `gorm:"not null;size:100" json:"name" binding:"required"` // Unique menu name
	Description  string        `gorm:"size:255" json:"description"`
	Price        float64       `gorm:"not null" json:"price" binding:"required"`         // Fixed combo price in euros
	IsAvailable  bool          `gorm:"default:true" json:"is_available"`                 // Unavailable menus cannot be ordered
	MenuProducts []MenuProduct `gorm:"foreignKey:MenuID" json:"menu_products"`           // Products included in this menu
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// MenuProduct is the join table between Menu and Products.
// It defines which products are part of a menu and in what quantity.
type MenuProduct struct {
	ID           uint     `gorm:"primaryKey" json:"id"`
	MenuID       uint     `gorm:"not null;constraint:OnDelete:CASCADE" json:"menu_id"` // FK to Menu — cascade deletes when menu is removed
	Menu         Menu     `gorm:"foreignKey:MenuID" json:"-"`
	ProductID    uint     `gorm:"not null" json:"product_id"`                          // FK to Products
	Product      Products `gorm:"foreignKey:ProductID" json:"-"`
	Quantity     uint     `gorm:"not null;default:1" json:"quantity"`                  // How many of this product are in the menu
	IsOptional   bool     `gorm:"default:false" json:"is_optional"`                   // Whether the product can be swapped out
	DisplayOrder uint     `gorm:"default:0" json:"display_order"`                     // Controls display order in the frontend
}
