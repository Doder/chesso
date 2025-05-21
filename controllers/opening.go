package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Doder/chesso/models"
)

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

        c.JSON(http.StatusCreated, input)
    }
}

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

