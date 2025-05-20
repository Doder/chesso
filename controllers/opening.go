package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Doder/chesso/models"
)

// List all openings
func ListOpenings(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var openings []models.Opening
        if err := db.Find(&openings).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, openings)
    }
}

// Get single opening by id
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

// Delete opening by id
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

