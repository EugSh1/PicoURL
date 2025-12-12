package utils

import (
	"context"
	"encoding/json"
	"picourl-backend/logger"
	"picourl-backend/redis"
	"time"
)

func CacheResult(key string, data any, cacheTTL time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var payload any

	switch v := data.(type) {
	case string:
		payload = v
	default:
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			logger.Log.Warn("JSON marshal error", "key", key, "error", err)
			return
		}

		payload = jsonBytes
	}

	err := redis.RedisClient.Set(ctx, key, payload, cacheTTL).Err()
	if err != nil {
		logger.Log.Warn("Failed to add data to the redis cache", "key", key, "error", err)
	}
}
