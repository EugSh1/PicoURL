package utils

import (
	"context"
	"encoding/json"
	"log"
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
			log.Println("JSON marshal error for cache:", key, "error:", err)
			return
		}

		payload = jsonBytes
	}

	err := redis.RedisClient.Set(ctx, key, payload, cacheTTL).Err()
	if err != nil {
		log.Println("Failed to add data to the redis cache:", key, "error:", err)
	}
}
