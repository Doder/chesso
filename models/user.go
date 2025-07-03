package models

import "time"

type User struct {
    ID          uint      `gorm:"primaryKey"`
    Username    string    `gorm:"uniqueIndex"`
		Email       string    `gorm:"uniqueIndex"`
		Password    string
    Rating      int
    LastLoggedIn *time.Time `gorm:"autoUpdateTime"`
}

