package models

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex"`
	Email        string `gorm:"uniqueIndex"`
	Password     string
	Rating       int
	LastLoggedIn *time.Time `gorm:"autoUpdateTime"`
	LoginCount   uint       `gorm:"default:0"`
}

type PasswordReset struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index"`
	Token     string `gorm:"uniqueIndex"`
	ExpiresAt time.Time
	CreatedAt time.Time
	User      User `gorm:"foreignKey:UserID"`
}
