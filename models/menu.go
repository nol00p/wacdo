package models

import (
	"time"
)

type Menu struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	Name         string        `gorm:"not null;size:100" json:"name" binding:"required"`
	Description  string        `gorm:"size:255" json:"description"`
	Price        float64       `gorm:"not null" json:"price" binding:"required"`
	IsAvailable  bool          `gorm:"default:true" json:"is_available"`
	MenuProducts []MenuProduct `gorm:"foreignKey:MenuID" json:"menu_products"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type MenuProduct struct {
	ID           uint     `gorm:"primaryKey" json:"id"`
	MenuID       uint     `gorm:"not null;constraint:OnDelete:CASCADE" json:"menu_id"` //this is delete the MenuProducts when a menu is deleted
	Menu         Menu     `gorm:"foreignKey:MenuID" json:"-"`
	ProductID    uint     `gorm:"not null" json:"product_id"`
	Product      Products `gorm:"foreignKey:ProductID" json:"-"`
	Quantity     uint     `gorm:"not null;default:1" json:"quantity"`
	IsOptional   bool     `gorm:"default:false" json:"is_optional"`
	DisplayOrder uint     `gorm:"default:0" json:"display_order"`
}
