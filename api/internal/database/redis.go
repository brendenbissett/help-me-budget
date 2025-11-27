package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// InitRedis initializes the Redis client connection
func InitRedis() error {
	redisURL := os.Getenv("REDIS_URL")

	var opts *redis.Options
	var err error

	if redisURL != "" {
		// Parse Redis URL if provided
		opts, err = redis.ParseURL(redisURL)
		if err != nil {
			return fmt.Errorf("unable to parse Redis URL: %w", err)
		}
	} else {
		// Build Redis options from individual env vars
		opts = &redis.Options{
			Addr:         getEnvOrDefault("REDIS_ADDR", "localhost:6379"),
			Password:     os.Getenv("REDIS_PASSWORD"), // Empty string if not set
			DB:           0,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolSize:     10,
			MinIdleConns: 5,
		}
	}

	client := redis.NewClient(opts)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("unable to connect to Redis: %w", err)
	}

	RedisClient = client
	log.Println("Redis connection established successfully")
	return nil
}

// CloseRedis closes the Redis client connection
func CloseRedis() {
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		} else {
			log.Println("Redis connection closed")
		}
	}
}
