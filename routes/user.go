package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/Doder/chesso/db"
    "github.com/Doder/chesso/models"
    "github.com/Doder/chesso/middleware"
    "github.com/Doder/chesso/controllers"
    "github.com/Doder/chesso/services"

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
		r.POST("/auth/forgot-password", controllers.ForgotPassword)
		r.POST("/auth/reset-password", controllers.ResetPassword)

    r.GET("/me", middleware.AuthMiddleware(), controllers.GetCurrentUser)
    
    // Test endpoint for training reminders (for development/testing)
    r.POST("/test-training-reminder", middleware.AuthMiddleware(), func(c *gin.Context) {
        worker := services.NewTrainingWorker()
        worker.TestTrainingReminder()
        c.JSON(http.StatusOK, gin.H{"message": "Training reminder test triggered"})
    })
}

