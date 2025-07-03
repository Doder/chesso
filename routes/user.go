package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/Doder/chesso/db"
    "github.com/Doder/chesso/models"
    "github.com/Doder/chesso/middleware"
    "github.com/Doder/chesso/controllers"

    "net/http"
)

func RegisterUserRoutes(r *gin.Engine) {
    r.GET("/users", middleware.AuthMiddleware(), func(c *gin.Context) {
        var users []models.User
        db.DB.Find(&users)
        c.JSON(http.StatusOK, users)
    })

		r.POST("/register", controllers.RegisterUser)
		r.POST("/login", controllers.LoginUser)

    r.GET("/me", middleware.AuthMiddleware(), controllers.GetCurrentUser)
}

