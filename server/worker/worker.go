package worker

import (
	"context"
	"picourl-backend/db"
	"picourl-backend/logger"
	"time"
)

const workerInterval = 24 * time.Hour

func Setup() {
	ticker := time.NewTicker(workerInterval)
	defer ticker.Stop()

	CleanupOldClicks()

	for range ticker.C {
		CleanupOldClicks()
	}
}

func CleanupOldClicks() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := db.Queries.DeleteOldClicks(ctx)
	if err != nil {
		logger.Log.Error("Failed to clean up old clicks", "error", err)
		return
	}

	logger.Log.Info("Successfully cleaned up old clicks")
}
