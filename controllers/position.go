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
		LastMove   string `json:"last_move" binding:"required"`
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
								LastMove: input.LastMove,
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
        fen := c.Query("fen")
        if fen == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "fen query param required"})
            return
        }

        hashedFen := utils.NormalizeHashFEN(fen) // Implement this function to hash FEN consistently

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
				candidatePositions := positionWithSameFen.NextPositions 
        if len(candidatePositions) == 0 {
            c.JSON(http.StatusOK, []interface{}{})
            return
        }

        c.JSON(http.StatusOK, candidatePositions)
    }
	}
