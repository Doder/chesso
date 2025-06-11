package models

import (
	"time"
)

type Opening struct {
	ID          uint       `gorm:"primaryKey"`
  CreatedAt   time.Time
  UpdatedAt   time.Time

	Name string `gorm:"not null;uniqueIndex:idx-user_repertoire_opening" json:"name"`
	Side string `gorm:"not null;type:char(1);uniqueIndex:idx-user_repertoire_opening;check:side IN ('w', 'b')" json:"side"` 
	RepertoireID uint `gorm:"not null;uniqueIndex:idx-user_repertoire_opening" json:"repertoire_id"`

	Repertoire Repertoire `gorm:"foreignKey:RepertoireID;constraint:OnDelete:CASCADE;" json:"-"`
}
