package models

import "gorm.io/gorm"

type Opening struct {
	gorm.Model

	Name string `gorm:"not null;uniqueIndex:idx-user_repertoire_opening" json:"name"`
	RepertoireID uint `gorm:"not null;uniqueIndex:idx-user_repertoire_opening" json:"repertoire_id"`

	Repertoire Repertoire `gorm:"foreignKey:RepertoireID" json: "-"`
}
