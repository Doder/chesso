package main

import (
	"github.com/gin-gonic/gin"
	"github.com/Doder/chesso/db"
	"github.com/Doder/chesso/models"
	"github.com/Doder/chesso/routes"
	"github.com/Doder/chesso/services"
)

func main() {
	router := gin.Default()

	db.Connect()
	db.DB.AutoMigrate(&models.User{}, &models.Repertoire{}, &models.Opening{}, &models.Position{}, &models.PasswordReset{})

	// Start training worker
	trainingWorker := services.NewTrainingWorker()
	trainingWorker.Start()

	routes.RegisterUserRoutes(router)
	routes.RegisterRepertoirRoutes(router, db.DB)
	router.Run(":8080")
}

