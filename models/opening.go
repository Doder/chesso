package models

import "gorm.io/gorm"

type Opening struct {
	gorm.Model

	Name string `gorm:"not null;uniqueIndex:idx-user_repertoire_opening" json:"name"`
	Side string `gorm:"not null;type:char(1);uniqueIndex:idx-user_repertoire_opening;check:side IN ('w', 'b')" json:"side"` 
	RepertoireID uint `gorm:"not null;uniqueIndex:idx-user_repertoire_opening" json:"repertoire_id"`

	Repertoire Repertoire `gorm:"foreignKey:RepertoireID" json:"-"`
}
