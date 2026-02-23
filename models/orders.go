package models

import "time"

// there is no point of sale option. the users are taking the orders on behalf of the customer.

type Order struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	CustomerID    *uint       `json:"customer_id"`
	Customer      Customer    `gorm:"foreignKey:CustomerID" json:"customer"`
	CreatedByID   uint        `gorm:"not null" json:"created_by_id"`
	CreatedBy     Users       `gorm:"foreignKey:CreatedByID" json:"created_by"`
	OrderType     string      `gorm:"not null" json:"order_type"`             // "counter" or "phone"
	Status        string      `gorm:"not null;default:pending" json:"status"` // pending, preparing, prepared, delivered, cancelled
	Notes         string      `json:"notes"`
	ScheduledTime *time.Time  `json:"scheduled_time"`
	TotalPrice    float64     `gorm:"not null;default:0" json:"total_price"`
	OrderItems    []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order_items"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID               uint              `gorm:"primaryKey" json:"id"`
	OrderID          uint              `gorm:"not null" json:"order_id"`
	ProductID        *uint             `json:"product_id"`
	Product          Products          `gorm:"foreignKey:ProductID" json:"product"`
	MenuID           *uint             `json:"menu_id"`
	Menu             Menu              `gorm:"foreignKey:MenuID" json:"menu"`
	Quantity         uint              `gorm:"not null;default:1" json:"quantity"`
	UnitPrice        float64           `gorm:"not null" json:"unit_price"`
	ItemTotal        float64           `gorm:"not null" json:"item_total"`
	OrderItemOptions []OrderItemOption `gorm:"foreignKey:OrderItemID;constraint:OnDelete:CASCADE" json:"order_item_options"`
}

type OrderItemOption struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	OrderItemID   uint         `gorm:"not null" json:"order_item_id"`
	OptionValueID uint         `gorm:"not null" json:"option_value_id"`
	OptionValue   OptionValues `gorm:"foreignKey:OptionValueID" json:"option_value"`
	PriceApplied  float64      `gorm:"not null" json:"price_applied"`
}
