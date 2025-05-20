package models

import "gorm.io/gorm"

type Position struct {
	gorm.Model	

	FEN string `gorm:"not null" json:"fen"`
	HashedFEN string `gorm:"index;not null" json:"hashed_fen` 
	LastMove string `json: "last_move"`
	MoveNumber uint `gorm:"not null" json:"move_number"` 
	OpeningID uint `gorm:"not null" json:"opening_id"` 	

	Opening Opening `gorm:"foreignKey:OpeningID json:"-"`
}
