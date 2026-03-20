package models

import "time"

// Customer represents an external customer who places orders.
// Customers are optional on orders — counter orders may not have a customer attached.
type Customer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name" binding:"required"` // Customer's full name
	Phone     string    `json:"phone"`                                   // Phone number, used for phone orders and duplicate detection
	Email     string    `json:"email" binding:"omitempty,email"`         // Optional email for contact
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
