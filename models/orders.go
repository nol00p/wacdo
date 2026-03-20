package models

import "time"

// Order represents a customer order taken by a staff member (accueil or admin).
// Staff members create orders on behalf of customers — there is no self-service kiosk mode.
// Status follows a state machine: pending → preparing → prepared → delivered (cancel only from pending).
type Order struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	CustomerID    *uint       `json:"customer_id"`                                                       // Optional FK to Customer — counter orders may have no customer
	Customer      Customer    `gorm:"foreignKey:CustomerID" json:"customer"`                              // Preloaded customer
	CreatedByID   uint        `gorm:"not null" json:"created_by_id"`                                     // FK to Users — the staff member who created the order
	CreatedBy     Users       `gorm:"foreignKey:CreatedByID" json:"created_by"`                           // Preloaded staff user
	OrderType     string      `gorm:"not null" json:"order_type"`                                        // "counter" (walk-in) or "phone" (call-in)
	Status        string      `gorm:"not null;default:pending" json:"status"`                             // pending, preparing, prepared, delivered, cancelled
	Notes         string      `json:"notes"`                                                              // Free-text notes for the kitchen
	ScheduledTime *time.Time  `json:"scheduled_time"`                                                    // Requested delivery time, used for preparation sorting
	TotalPrice    float64     `gorm:"not null;default:0" json:"total_price"`                              // Server-computed total (sum of all item totals)
	OrderItems    []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order_items"` // Line items in this order
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

// OrderItem is a single line in an order. Each item references either a Product or a Menu (exactly one, never both).
// Prices are captured at order time so they are immutable even if the product price changes later.
type OrderItem struct {
	ID               uint              `gorm:"primaryKey" json:"id"`
	OrderID          uint              `gorm:"not null" json:"order_id"`                                              // FK to Order
	ProductID        *uint             `json:"product_id"`                                                            // FK to Products (nil if this item is a menu)
	Product          Products          `gorm:"foreignKey:ProductID" json:"product"`
	MenuID           *uint             `json:"menu_id"`                                                               // FK to Menu (nil if this item is a product)
	Menu             Menu              `gorm:"foreignKey:MenuID" json:"menu"`
	Quantity         uint              `gorm:"not null;default:1" json:"quantity"`                                    // Number of this item ordered
	UnitPrice        float64           `gorm:"not null" json:"unit_price"`                                           // Price per unit at order time (product price or menu price)
	ItemTotal        float64           `gorm:"not null" json:"item_total"`                                           // (UnitPrice + option prices) * Quantity
	OrderItemOptions []OrderItemOption `gorm:"foreignKey:OrderItemID;constraint:OnDelete:CASCADE" json:"order_item_options"` // Selected options for this item
}

// OrderItemOption records a selected option value for an order item (e.g. "Large" size).
// The price is captured at order time to preserve the price even if the option value changes later.
type OrderItemOption struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	OrderItemID   uint         `gorm:"not null" json:"order_item_id"`                    // FK to OrderItem
	OptionValueID uint         `gorm:"not null" json:"option_value_id"`                  // FK to OptionValues
	OptionValue   OptionValues `gorm:"foreignKey:OptionValueID" json:"option_value"`
	PriceApplied  float64      `gorm:"not null" json:"price_applied"`                    // Option price snapshot at order time
}
