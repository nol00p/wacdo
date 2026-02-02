package models

import "time"

type Users struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `json:"username"`
	Email     string    `gorm:"unique" json:"email" binding:"required,email"`
	Password  string    `json:"password" binding:"required,min=6"`
	RolesID   uint      `gorm:"not null" json:"roles_id"` // only 1 role per user
	Role      Roles     `gorm:"foreignKey:RolesID"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
