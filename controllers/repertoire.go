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
				userIDValue,_ := c.Get("userID")
				userIDFloat := userIDValue.(float64)
        input.UserID = uint(userIDFloat)
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
