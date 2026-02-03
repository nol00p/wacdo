package models

import (
	"time"
)

type Category struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	DisplayOrder uint      `json:"display_order" `
	ImageURL     string    `json:"image_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Products struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	CategoryID      uint      `json:"category_id"`
	Category        Category  `gorm:"foreignKey:CategoryID" json:"category"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Price           float64   `json:"price"`
	StockQuantity   uint      `json:"stock_quantity"`
	IsAvailable     bool      `json:"is_available"`
	PreparationTime uint      `json:"preparation_time"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ProductOptions struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	ProductID  uint     `json:"product_id"`
	Product    Products `gorm:"foreignKey:ProductID" json:"-"`
	Name       string   `json:"name"`
	IsUnique   string   `json:"is_unique"` // define if the option is single or multiple choise
	IsRequired bool     `json:"is_required"`
}

type OptionValues struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	OptionID    uint           `json:"option_id"`
	Option      ProductOptions `gorm:"foreignKey:OptionID" json:"-"`
	Value       string         `json:"value"`
	OptionPrice float64        `json:"option_price"`
}
