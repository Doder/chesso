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

