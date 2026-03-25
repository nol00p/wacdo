// Package models defines the GORM data models and database schema for the application.
package models

import (
	"time"

	"gorm.io/gorm"
)

// Users represents an internal back-office user (staff member).
// Each user has exactly one role that determines their permissions.
// Uses GORM soft delete — deleted users are hidden from queries but preserved for order audit trails.
// Email uniqueness is enforced at the application level (not DB unique constraint) so that
// a deleted user's email can be reused when creating a new account.
type Users struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `json:"username"`                          // Display name
	Email     string         `gorm:"uniqueIndex:idx_users_email_active,where:deleted_at IS NULL" json:"email" binding:"required,email"` // Unique among non-deleted users
	Password  string         `json:"-" binding:"required,min=6"`        // Bcrypt-hashed password, hidden from JSON output
	RolesID   uint           `gorm:"not null" json:"roles_id"`          // FK to Roles — each user has exactly one role
	Role      Roles          `gorm:"foreignKey:RolesID"`                // Preloaded role relationship
	IsActive  bool           `gorm:"default:true" json:"is_active"`     // Deactivated users cannot log in
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                    // Soft delete timestamp — nil means active
}

// UserInput is the request body for creating a user.
// Separates input validation from the database model to avoid leaking fields like ID or timestamps.
type UserInput struct {
	Username string `json:"username"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"` // Plain-text password, hashed before storage
	RolesID  uint   `json:"roles_id"`                          // Must reference an existing role
}
