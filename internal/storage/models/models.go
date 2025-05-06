package models

import (
	//"gorm.io/gorm"
	"time"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	Password string `gorm:"unique;not null"`
	Tasks    []Task `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

type Task struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"not null;size:255"`
	Description string    `gorm:"size:1024"`
	Completed   bool      `gorm:"default:false"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UserID      uint      `gorm:"not null;index"`
}
