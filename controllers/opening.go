package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Doder/chesso/models"
    "github.com/Doder/chesso/utils"
)

type OpeningInput struct {
	Name string `json:"name"`
	Side string `json:"side"`
}

func CreateOpening(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input models.Opening
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        if err := db.Create(&input).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
			  //create also first position
				fpFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
				fpFENHashed := utils.NormalizeHashFEN(fpFEN)
				firstPosition := models.Position{
					FEN: fpFEN,
					HashedFEN: fpFENHashed,
					OpeningID: input.ID,
				}

        if err := db.Create(&firstPosition).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusCreated, input)
    }
}

func ListOpenings(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("userID")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
            return
        }

        var openings []models.Opening
        if err := db.
            Joins("JOIN repertoires ON repertoires.id = openings.repertoire_id").
            Where("repertoires.user_id = ?", userID).
            Order("openings.created_at ASC").
            Find(&openings).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, openings)
    }
}

func UpdateOpening(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input OpeningInput
		if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
		}
		var opening models.Opening
		if err := db.First(&opening, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := db.Model(&opening).Updates(input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, opening)
	}
}

func GetOpening(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        var opening models.Opening

        if err := db.First(&opening, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Opening not found"})
            return
        }
        c.JSON(http.StatusOK, opening)
    }
}

func DeleteOpening(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        if err := db.Delete(&models.Opening{}, id).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.Status(http.StatusNoContent)
    }
}

