// Package models defines the GORM data models and database schema for the application.
package models

import "time"

// Users represents an internal back-office user (staff member).
// Each user has exactly one role that determines their permissions.
type Users struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `json:"username"`                          // Display name
	Email     string    `gorm:"unique" json:"email" binding:"required,email"` // Unique login identifier
	Password  string    `json:"-" binding:"required,min=6"`        // Bcrypt-hashed password, hidden from JSON output
	RolesID   uint      `gorm:"not null" json:"roles_id"`          // FK to Roles — each user has exactly one role
	Role      Roles     `gorm:"foreignKey:RolesID"`                // Preloaded role relationship
	IsActive  bool      `gorm:"default:true" json:"is_active"`     // Deactivated users cannot log in
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserInput is the request body for creating a user.
// Separates input validation from the database model to avoid leaking fields like ID or timestamps.
type UserInput struct {
	Username string `json:"username"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"` // Plain-text password, hashed before storage
	RolesID  uint   `json:"roles_id"`                          // Must reference an existing role
}
