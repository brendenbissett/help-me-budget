package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/brendenbissett/help-me-budget/api/internal/database"
)

// CreateRedisSession creates a session record in Redis for admin tracking
func CreateRedisSession(ctx context.Context, userID, email, name, provider string) error {
	redisSessionKey := fmt.Sprintf("session:%s", userID)
	sessionData := map[string]interface{}{
		"user_id":  userID,
		"email":    email,
		"name":     name,
		"provider": provider,
		"login_at": time.Now().Format(time.RFC3339),
	}
	sessionJSON, _ := json.Marshal(sessionData)
	return database.RedisClient.Set(ctx, redisSessionKey, sessionJSON, 24*time.Hour).Err()
}
