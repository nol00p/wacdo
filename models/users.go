package models

import "time"

type Users struct {
	ID        uint `gorm:"primaryKey"`
	Username  string
	Email     string `gorm:"unique"`
	Password  string
	RolesID   uint  `gorm:"not null"`
	Role      Roles `gorm:"foreignKey:RolesID"`
	IsActive  bool  `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
