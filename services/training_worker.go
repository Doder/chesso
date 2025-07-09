package services

import (
	"log"
	"time"

	"github.com/Doder/chesso/db"
	"github.com/Doder/chesso/models"
	"github.com/robfig/cron/v3"
)

// TrainingWorker handles scheduled training reminder emails
type TrainingWorker struct {
	cron *cron.Cron
}

// NewTrainingWorker creates a new training worker
func NewTrainingWorker() *TrainingWorker {
	// Create cron with timezone support
	c := cron.New()

	return &TrainingWorker{
		cron: c,
	}
}

// Start begins the training worker with scheduled jobs
func (tw *TrainingWorker) Start() {
	// Schedule daily training reminder check at midnight
	// This will run every day at midnight
	tw.cron.AddFunc("0 0 */3 * *", tw.checkAndSendTrainingReminders)

	// For testing purposes, you can uncomment this line to run every 5 minutes
	// tw.cron.AddFunc("*/5 * * * *", tw.checkAndSendTrainingReminders)

	tw.cron.Start()
	log.Println("Training worker started - will check for training reminders daily at midnight")
}

// Stop shuts down the training worker
func (tw *TrainingWorker) Stop() {
	tw.cron.Stop()
	log.Println("Training worker stopped")
}

// checkAndSendTrainingReminders checks for users who need training reminders
func (tw *TrainingWorker) checkAndSendTrainingReminders() {
	log.Println("Checking for users who need training reminders...")

	// Get all users who have logged in within the last 30 days
	// This prevents spamming inactive users
	cutoffDate := time.Now().AddDate(0, 0, -30) // 30 days ago

	var users []models.User
	err := db.DB.Where("last_logged_in IS NOT NULL AND last_logged_in > ?", cutoffDate).Find(&users).Error
	if err != nil {
		log.Printf("Error fetching active users: %v", err)
		return
	}

	log.Printf("Found %d active users to check for training reminders", len(users))

	for _, user := range users {
		tw.processUserForTrainingReminder(user)
	}
}

// processUserForTrainingReminder processes a single user for training reminder
func (tw *TrainingWorker) processUserForTrainingReminder(user models.User) {
	// Get user's repertoires
	var repertoires []models.Repertoire
	err := db.DB.Where("user_id = ?", user.ID).Find(&repertoires).Error
	if err != nil {
		log.Printf("Error fetching repertoires for user %d: %v", user.ID, err)
		return
	}

	if len(repertoires) == 0 {
		log.Printf("User %d has no repertoires, skipping", user.ID)
		return
	}

	// Get all openings for this user
	var openings []models.Opening
	repertoireIDs := make([]uint, len(repertoires))
	for i, rep := range repertoires {
		repertoireIDs[i] = rep.ID
	}

	err = db.DB.Where("repertoire_id IN ?", repertoireIDs).Find(&openings).Error
	if err != nil {
		log.Printf("Error fetching openings for user %d: %v", user.ID, err)
		return
	}

	if len(openings) == 0 {
		log.Printf("User %d has no openings, skipping", user.ID)
		return
	}

	// Calculate positions to review for each opening
	trainingData := tw.calculateTrainingData(openings)

	if len(trainingData) == 0 {
		log.Printf("User %d has no positions to train, skipping", user.ID)
		return
	}

	// Calculate total positions
	totalPositions := 0
	for _, opening := range trainingData {
		totalPositions += opening.PositionCount
	}

	// Send training reminder email
	reminderData := TrainingReminderData{
		Username:       user.Username,
		Openings:       trainingData,
		TotalPositions: totalPositions,
	}

	err = SendTrainingReminderEmail(user.Email, reminderData)
	if err != nil {
		log.Printf("Error sending training reminder to user %d (%s): %v", user.ID, user.Email, err)
	} else {
		log.Printf("Training reminder sent successfully to user %d (%s) - %d positions", user.ID, user.Email, totalPositions)
	}
}

// calculateTrainingData calculates positions to review for each opening
func (tw *TrainingWorker) calculateTrainingData(openings []models.Opening) []OpeningTrainingData {
	var trainingData []OpeningTrainingData

	for _, opening := range openings {
		count := tw.getPositionsDueForReview(opening.ID)
		if count > 0 {
			trainingData = append(trainingData, OpeningTrainingData{
				Name:          opening.Name,
				Side:          opening.Side,
				PositionCount: count,
			})
		}
	}

	return trainingData
}

// getPositionsDueForReview returns the number of positions due for review for a specific opening
// This implements the same logic as used in the frontend Train page
func (tw *TrainingWorker) getPositionsDueForReview(openingID uint) int {
	// Get all positions for this opening
	var positions []models.Position
	err := db.DB.Where("opening_id = ?", openingID).Find(&positions).Error
	if err != nil {
		log.Printf("Error fetching positions for opening %d: %v", openingID, err)
		return 0
	}

	now := time.Now()
	dueCount := 0

	for _, position := range positions {
		if tw.isPositionDueForReview(position, now) {
			dueCount++
		}
	}

	return dueCount
}

// isPositionDueForReview determines if a position is due for review based on spaced repetition
func (tw *TrainingWorker) isPositionDueForReview(position models.Position, now time.Time) bool {
	// If never trained, it's due for review
	if position.LastCorrectGuess == nil {
		return true
	}

	// Calculate next review date based on spaced repetition
	// The interval increases with each correct guess
	var intervalDays int
	switch position.RepetitionCount {
	case 0:
		intervalDays = 1
	case 1:
		intervalDays = 3
	case 2:
		intervalDays = 7
	case 3:
		intervalDays = 14
	case 4:
		intervalDays = 30
	case 5:
		intervalDays = 60
	default:
		intervalDays = 120 // Cap at 4 months
	}

	nextReviewDate := position.LastCorrectGuess.AddDate(0, 0, intervalDays)
	return now.After(nextReviewDate)
}

// TestTrainingReminder manually triggers a training reminder check (for testing)
func (tw *TrainingWorker) TestTrainingReminder() {
	log.Println("Manual test of training reminder system...")
	tw.checkAndSendTrainingReminders()
}
