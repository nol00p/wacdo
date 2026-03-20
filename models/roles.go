package models

import "time"

// Roles defines the access level for a user.
// Three roles are expected: "admin", "accueil" (front desk), and "preparation" (kitchen).
type Roles struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	RoleName    string    `gorm:"size:50;unique;not null" json:"role_name"` // Unique role identifier used in authorization checks
	Description string    `gorm:"size:255" json:"description"`             // Human-readable description of the role
	Permissions string    `gorm:"type:text" json:"permissions"`            // Free-text permissions field (reserved for future use)
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
