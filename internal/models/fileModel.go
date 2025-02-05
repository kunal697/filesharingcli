package models

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	SiteName  string     `gorm:"not null" json:"site_name"`
	FileName  string     `gorm:"not null" json:"file_name"`
	FileURL   string     `gorm:"not null" json:"file_url"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (f *File) BeforeCreate(tx *gorm.DB) (err error) {
	// Any pre-create logic, if needed
	return
}
