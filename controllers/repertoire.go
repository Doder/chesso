package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Doder/chesso/models"
)

func CreateRepertoire(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input models.Repertoire
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if err := db.Create(&input).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusCreated, input)
    }
}

func ListRepertoires(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var repertoires []models.Repertoire
        if err := db.Find(&repertoires).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, repertoires)
    }
}

func GetRepertoire(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var rep models.Repertoire
        id := c.Param("id")

        if err := db.First(&rep, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Repertoire not found"})
            return
        }

        c.JSON(http.StatusOK, rep)
    }
}

func DeleteRepertoire(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")

        if err := db.Delete(&models.Repertoire{}, id).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.Status(http.StatusNoContent)
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
