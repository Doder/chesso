package controllers

import (
    "net/http"
    "strings"
    "crypto/sha256"
    "encoding/hex"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Doder/chesso/models"
	  "github.com/Doder/chesso/utils"
)

type PositionInput struct {
    FEN        string `json:"fen" binding:"required"`
    MoveNumber int    `json:"move_number"`
    OpeningID  uint   `json:"opening_id"`
}

// POST /positions
func CreatePosition(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input PositionInput
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        position := models.Position{
            FEN:        input.FEN,
            MoveNumber: uint(input.MoveNumber),
            OpeningID:  input.OpeningID,
            HashedFEN:  utils.HashFEN(input.FEN),
        }

        if err := db.Create(&position).Error; err != nil {
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

        hashed := utils.HashFEN(fen)

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
        fen := c.Query("fen")
        if fen == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "fen query param required"})
            return
        }

        hashedFen := hashFen(fen) // Implement this function to hash FEN consistently

        var positionsWithSameFen []models.Position
        // Step 1: Find all positions with the same hashedFen, preload Opening
        if err := db.Preload("Opening").
            Where("hashed_fen = ?", hashedFen).
            Find(&positionsWithSameFen).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        if len(positionsWithSameFen) == 0 {
            c.JSON(http.StatusOK, gin.H{"message": "No matching positions found"})
            return
        }

        // Collect OpeningIDs and MoveNumbers from these positions
        type openingMove struct {
            OpeningID  uint
            MoveNumber int
        }
        var openingMoves []openingMove
        for _, pos := range positionsWithSameFen {
            if pos.OpeningID != 0 {
                openingMoves = append(openingMoves, openingMove{pos.OpeningID, int(pos.MoveNumber)})
            }
        }

        if len(openingMoves) == 0 {
            c.JSON(http.StatusOK, gin.H{"message": "No openings associated with matched positions"})
            return
        }

        // Step 2: For each opening & moveNumber, find positions with same OpeningID and moveNumber + 1
        var candidatePositions []models.Position
        query := db.Preload("Opening").Where("0=1") // start with false condition to OR later

        for _, om := range openingMoves {
            query = query.Or("opening_id = ? AND move_number = ?", om.OpeningID, om.MoveNumber+1)
        }

        if err := query.Find(&candidatePositions).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, candidatePositions)
    }
	}
