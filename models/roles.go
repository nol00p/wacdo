package models

import "time"

type Roles struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	RoleName    string    `gorm:"size:50;unique;not null" json:"role_name"`
	Description string    `gorm:"size:255" json:"description"`
	Permissions string    `gorm:"type:text" json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
