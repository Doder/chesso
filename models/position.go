package models

import (
	"gorm.io/gorm"
	"github.com/Doder/chesso/utils"
)

type Position struct {
	gorm.Model	

	FEN string `gorm:"not null" json:"fen"`
	HashedFEN string `gorm:"not null;uniqueIndex:idx-fen-opening" json:"hashed_fen"` 
	OpeningID uint `gorm:"not null;uniqueIndex:idx-fen-opening" json:"opening_id"` 	

	OpeningName string `json:"opening_name" gorm:"-"`
	Opening Opening `gorm:"foreignKey:OpeningID" json:"-"`
	PrevPositions []*Position `gorm:"many2many:position_prevposition;joinForeignKey:PositionID;joinReferences:PrevPositionID"`
	NextPositions []*Position `gorm:"many2many:position_prevposition;joinForeignKey:PrevPositionID;joinReferences:PositionID"`
}

func (p *Position) BeforeSave(tx *gorm.DB) (err error) {
    normalized := utils.NormalizeFEN(p.FEN)
    p.HashedFEN = utils.HashFEN(normalized)
    return nil
}
