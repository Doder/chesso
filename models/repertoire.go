package models

import (
	"time"
)

type Repertoire struct {
	ID          uint       `gorm:"primaryKey"`
  CreatedAt   time.Time
  UpdatedAt   time.Time

	Name string `gorm:"not null;uniqueIndex:idx_user_repertoire_name" json:"name"`
	UserID uint `gorm:"not null;uniqueIndex:idx_user_repertoire_name" json:"user_id"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
}
