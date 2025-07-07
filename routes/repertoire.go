package routes

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Doder/chesso/controllers"
		"github.com/Doder/chesso/middleware"
)

func RegisterRepertoirRoutes(router *gin.Engine, db *gorm.DB) {
		rep := router.Group("/repertoires")
		rep.Use(middleware.AuthMiddleware())
    {
        rep.POST("/", controllers.CreateRepertoire(db))
        rep.GET("/", controllers.ListRepertoires(db))
        rep.GET("/:id", controllers.GetRepertoire(db))
				rep.PATCH("/:id", controllers.UpdateRepertoire(db))
        rep.DELETE("/:id", controllers.DeleteRepertoire(db))
    }
		openings := router.Group("/openings")
		openings.Use(middleware.AuthMiddleware())
    {
			  openings.POST("/", controllers.CreateOpening(db))
        openings.GET("/", controllers.ListOpenings(db))
        openings.GET("/:id", controllers.GetOpening(db))
				openings.PATCH("/:id", controllers.UpdateOpening(db))
        openings.DELETE("/:id", controllers.DeleteOpening(db))
    }
    positions := router.Group("/positions")
		positions.Use(middleware.AuthMiddleware())
    {
				positions.GET("/search-candidate", controllers.SearchCandidatePositions(db)) 
				positions.GET("/by-openings", controllers.GetPositionsByOpeningIds(db))
				positions.GET("/counts-by-openings", controllers.GetPositionCountsByOpeningIds(db))
				positions.PATCH("/:id", controllers.UpdatePositionMeta(db))
				positions.POST("/:id/correct", controllers.UpdatePositionCorrectGuess(db))
				positions.POST("/:id/incorrect", controllers.ResetPositionProgress(db))
				positions.DELETE("/:id", controllers.DeletePosition(db))
		}
}
