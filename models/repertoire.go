package models

import "gorm.io/gorm"


type Repertoire struct {
	gorm.Model
	Name string          `gorm:"not null;uniqueIndex:idx_user_repertoire_name" json:"name"`
	UserID uint          `gorm:"not null;uniqueIndex:idx_user_repertoire_name" json:"user_id"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}
