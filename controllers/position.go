package controllers

import (
    "net/http"
		"errors"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Doder/chesso/models"
	  "github.com/Doder/chesso/utils"
)

type PositionInput struct {
		FromFEN    string `json:"from_fen" binding:"required"`
		ToFEN      string `json:"to_fen" binding:"required"`
		OpeningID  uint   `json:"opening_id" binding:"required"`
		RepertoireID uint `json:"repertoire_id" binding:"required"`
}

// POST /positions
func CreateOrUpdatePosition(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input PositionInput
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

				var prevPos models.Position
				fromFENHashed := utils.NormalizeHashFEN(input.FromFEN)
				if err := db.Joins("JOIN openings ON openings.id = positions.opening_id").Where("hashed_fen = ? AND openings.repertoire_id = ?", fromFENHashed, input.RepertoireID).First(&prevPos).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Related position not found"})
					return
				}
				hashedFEN := utils.HashFEN(input.ToFEN)
				var pos models.Position 
				var position models.Position
				if err := db.
    Joins("JOIN openings ON openings.id = positions.opening_id").
		Where("positions.hashed_fen = ? AND openings.repertoire_id = ?", hashedFEN, input.RepertoireID).
    Preload("Opening.Repertoire").
    First(&pos).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						position = models.Position{
								FEN:        input.ToFEN,
								OpeningID:  input.OpeningID,
								HashedFEN:  hashedFEN,
						}
						if err := db.Create(&position).Error; err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
						}
					} else {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
				}
				if err := db.Model(&prevPos).Association("NextPositions").Append(&position); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

        c.JSON(http.StatusOK, position)
    }
}

// GET /positions/search?fen=...
func FindPositionsByFEN(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        fen := c.Query("fen")
        if fen == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "FEN is required"})
            return
        }

        hashed := utils.NormalizeHashFEN(fen)

        var positions []models.Position
        if err := db.Preload("Opening").Where("hashed_fen = ?", hashed).Find(&positions).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, positions)
    }
}

// DeletePosition deletes position by ID
func DeletePosition(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")

        if err := db.Delete(&models.Position{}, id).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.Status(http.StatusNoContent)
    }
}

// SearchCandidatePositions finds candidate positions based on fen
func SearchCandidatePositions(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
			  rep_id := c.Query("repertoire_id")
        to_fen := c.Query("to_fen")
				from_fen := c.Query("from_fen")
        if to_fen == "" || rep_id == "" || to_fen == from_fen{
            c.JSON(http.StatusBadRequest, gin.H{"error": "to_fen, repertoire_id query params are required"})
            return
        }

        hashedFromFen := utils.NormalizeHashFEN(from_fen) 
				hashedToFen := utils.NormalizeHashFEN(to_fen)

        var position1WithSameFen models.Position
        var position2WithSameFen models.Position

        err1 := db.
    Joins("JOIN openings ON openings.id = positions.opening_id").
    Where("positions.hashed_fen = ? AND openings.repertoire_id = ?", hashedFromFen, rep_id).
    Preload("Opening.Repertoire").
		Preload("NextPositions").
    First(&position1WithSameFen).Error

        err2 := db.
    Joins("JOIN openings ON openings.id = positions.opening_id").
    Where("positions.hashed_fen = ? AND openings.repertoire_id = ?", hashedToFen, rep_id).
    Preload("Opening.Repertoire").
		Preload("NextPositions").
    First(&position2WithSameFen).Error

				position := models.Position{
						FEN:        to_fen,
						OpeningID:  position1WithSameFen.OpeningID,
						HashedFEN:  hashedToFen,
				}
				if err1 != nil && err2 != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
						return
				} else if err1 != nil{
					if !errors.Is(err1, gorm.ErrRecordNotFound){
								c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
								return
					}
				} else if err2 != nil {
						if errors.Is(err2, gorm.ErrRecordNotFound){
							//create new positions regularly
							if err := db.Create(&position).Error; err != nil {
									c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
									return
							}
							if err := db.Model(&position1WithSameFen).Association("NextPositions").Append(&position); err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}
						} else {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
								return
						}
				} else if err1 == nil && err2 == nil {
							//append to to_pos new positions
							if err := db.Model(&position2WithSameFen).Association("PrevPositions").Append(&position1WithSameFen); err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}
					}
				candidatePositions := position2WithSameFen.NextPositions 
				if len(candidatePositions) == 0 {
						c.JSON(http.StatusOK, []interface{}{})
						return
				}

				c.JSON(http.StatusOK, candidatePositions)
    }
	}
	
func FindPrevMoves(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
			  rep_id := c.Query("repertoire_id")
        fen := c.Query("fen")
        if fen == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "fen query param required"})
            return
        }

        hashedFen := utils.NormalizeHashFEN(fen) 

        var positionWithSameFen models.Position
        if err := db.
    Joins("JOIN openings ON openings.id = positions.opening_id").
    Where("positions.hashed_fen = ? AND openings.repertoire_id = ?", hashedFen, rep_id).
    Preload("Opening.Repertoire").
		Preload("NextPositions").
    First(&positionWithSameFen).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
							return
        }
				candidatePositions := positionWithSameFen.PrevPositions 
        if len(candidatePositions) == 0 {
            c.JSON(http.StatusOK, []interface{}{})
            return
        }

        c.JSON(http.StatusOK, candidatePositions)
    }
	}

