package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint         `gorm:"primarykey" json:"id"`
	Email     string         `gorm:"uniqueIndex;size:255" json:"email"`
	Username  string         `gorm:"size:100" json:"username"`
	Password  string         `gorm:"size:255" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type File struct {
	ID          uint         `gorm:"primarykey" json:"id"`
	Filename    string         `gorm:"size:512" json:"filename"`
	StoragePath string         `gorm:"size:1024" json:"-"`
	Size        int64          `json:"size"`
	Description string         `gorm:"size:1024" json:"description"`
	UploaderID  string         `json:"uploader_id"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
