package models

import (
	"gorm.io/gorm"
	"github.com/Doder/chesso/utils"
)

type Position struct {
	gorm.Model	

	FEN string `gorm:"not null" json:"fen"`
	HashedFEN string `gorm:"not null;unique" json:"hashed_fen"` 
	LastMove string `json:"last_move"`
	MoveNumber uint `gorm:"not null" json:"move_number"` 
	OpeningID uint `gorm:"not null" json:"opening_id"` 	

	OpeningName string `json:"opening_name" gorm:"-"`
	Opening Opening `gorm:"foreignKey:OpeningID" json:"-"`
}

func (p *Position) BeforeSave(tx *gorm.DB) (err error) {
    normalized := utils.NormalizeFEN(p.FEN)
    p.HashedFEN = utils.HashFEN(normalized)
    return nil
}
