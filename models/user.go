package models

type User struct {
    ID       uint   `gorm:"primaryKey"`
    Username string `gorm:"uniqueIndex"`
		Email    string `gorm:"uniqueIndex"`
		Password string
    Rating   int
}

