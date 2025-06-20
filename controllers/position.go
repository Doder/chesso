package controllers

import (
    "net/http"
		"errors"
		"strconv"

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

type PositionMetaInput struct {
	Eval         string `json:"eval"`
	Comment      string `json:"comment"`
	Order        int    `json:"order"`
}

func SearchCandidatePositions(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
			  rep_id := c.Query("repertoire_id")
				op_id := c.Query("opening_id")
        to_fen := c.Query("to_fen")
				from_fen := c.Query("from_fen")
        if to_fen == "" || rep_id == "" || to_fen == from_fen || op_id == ""{
            c.JSON(http.StatusBadRequest, gin.H{"error": "to_fen, repertoire_id, opening_id query params are required"})
            return
        }
				opening_id, err := strconv.ParseUint(op_id, 10, 0)
				if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		Preload("NextPositions.Opening").
    First(&position2WithSameFen).Error

				position := models.Position{
						FEN:        to_fen,
						OpeningID:  uint(opening_id),
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
				for _,p := range candidatePositions {
					p.OpeningName = p.Opening.Name
				}
				c.JSON(http.StatusOK, candidatePositions)
    }
	}

func UpdatePositionMeta(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input PositionMetaInput
		if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
		}
		var position models.Position
		if err := db.First(&position, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := db.Model(&position).Updates(input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, position)
		return
	}
}

func DeletePosition(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		visited := map[uint]bool{}
		pos_id, err := strconv.ParseUint(id, 10, 0)
		if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
		}
		if err := recursiveDelete(db, uint(pos_id), visited); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, nil)
	}
}

func recursiveDelete(db *gorm.DB, id uint, visited map[uint]bool) error {
	if visited[id] {
		return nil // avoid cycles
	}
	visited[id] = true

	// Step 1: Find children (positions this one points to)
	var children []uint
	db.
		Table("position_prevposition").
		Where("prev_position_id = ?", id).
		Pluck("position_id", &children)

	for _, childID := range children {
		// Step 2: Check if any other parent exists for this child
		var count int64
		db.
			Table("position_prevposition").
			Where("position_id = ? AND prev_position_id != ?", childID, id).
			Count(&count)

		if count == 0 {
			// No other parents — delete child recursively
			if err := recursiveDelete(db, childID, visited); err != nil {
				return err
			}
		}
	}

	// Step 3: Delete position_relations where this position is involved
	err := db.
		Table("position_prevposition").
		Where("position_id = ? OR prev_position_id = ?", id, id).
		Delete(nil).Error
	if err != nil {
		return err
	}

	// Step 4: Delete the position itself
	db.Delete(&models.Position{}, id)
	return nil
}

