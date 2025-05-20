package routes

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/Doder/chesso/controllers"
)

func RegisterRepertoirRoutes(router *gin.Engine, db *gorm.DB) {
		rep := router.Group("/repertoires")
    {
        rep.POST("/", controllers.CreateRepertoire(db))
        rep.GET("/", controllers.ListRepertoires(db))
        rep.GET("/:id", controllers.GetRepertoire(db))
        rep.DELETE("/:id", controllers.DeleteRepertoire(db))
    }
		openings := router.Group("/openings")
    {
        openings.GET("/", controllers.ListOpenings(db))
        openings.GET("/:id", controllers.GetOpening(db))
        openings.DELETE("/:id", controllers.DeleteOpening(db))
    }
    positions := router.Group("/positions")
    {
        positions.POST("/", controllers.CreatePosition(db))
        positions.GET("/search", controllers.FindPositionsByFEN(db))
				positions.DELETE("/:id", controllers.DeletePosition(db))
				positions.GET("/search-candidate", controllers.SearchCandidatePositions(db)) // ?fen=xxx}
		}
}
