package models

import (
	"time"
	// "gorm.io/gorm"
)

type Site struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SiteName  string    `gorm:"unique;not null" json:"site_name"` // Ensure the site name is unique
	Password  string    `gorm:"not null" json:"password"`         // Store password
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
