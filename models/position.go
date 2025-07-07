package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/Doder/chesso/utils"
)

type Position struct {
	ID          uint       `gorm:"primaryKey"`
  CreatedAt   time.Time
  UpdatedAt   time.Time

	FEN string `gorm:"not null" json:"fen"`
	HashedFEN string `gorm:"not null;uniqueIndex:idx-fen-opening" json:"hashed_fen"` 
	OpeningID uint `gorm:"not null;uniqueIndex:idx-fen-opening" json:"opening_id"` 	
	Eval string `gorm:"default:'=';check: eval IN ('+', '+=', '=', '-=', '-')" json:"eval"`
	Comment string `json:"comment"`
	Order uint `gorm:"default:0" json:"order"`
	
	// Spaced repetition fields
	LastCorrectGuess *time.Time `json:"last_correct_guess"`
	RepetitionCount uint `gorm:"default:0" json:"repetition_count"`

	OpeningName string `json:"opening_name" gorm:"-"`
	Opening Opening `gorm:"foreignKey:OpeningID;constraint:OnDelete:CASCADE;" json:"-"`
	PrevPositions []*Position `gorm:"many2many:position_prevposition;joinForeignKey:PositionID;joinReferences:PrevPositionID;constraint:OnDelete:CASCADE;"`
	NextPositions []*Position `gorm:"many2many:position_prevposition;joinForeignKey:PrevPositionID;joinReferences:PositionID;constraint:OnDelete:CASCADE;"`
}

func (p *Position) BeforeSave(tx *gorm.DB) (err error) {
    normalized := utils.NormalizeFEN(p.FEN)
    p.HashedFEN = utils.HashFEN(normalized)
    return nil
}
