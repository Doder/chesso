package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Doder/chesso/models"
)

type RepertoireWithOpenings struct {
	gorm.Model
	Name string `json:"name"`
	Openings []models.Opening `gorm:"foreignKey:RepertoireID" json:"openings"`
}

type RepertoireInput struct {
	Name string `json:"name"`
}

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
				userID,_ := c.Get("userID")
        if err := db.Where("user_id = ?", userID).Order("created_at ASC").Find(&repertoires).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, repertoires)
    }
}

func GetRepertoire(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var rep RepertoireWithOpenings
        id := c.Param("id")
				
        if err := db.Model(&models.Repertoire{}).Where("id=?", id).Preload("Openings", func(db *gorm.DB) *gorm.DB {
					return db.Order("openings.created_at ASC")
				}).First(&rep).Error; err != nil {
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

func UpdateRepertoire(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input RepertoireInput
		if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
		}
		var repertoire models.Repertoire
		if err := db.First(&repertoire, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := db.Model(&repertoire).Updates(input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, repertoire)
	}
}
