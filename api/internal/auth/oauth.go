package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/brendenbissett/help-me-budget/api/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis/v3"
	"github.com/google/uuid"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
)

const (
	// session key names stored in fiber session
	gothSessionKeyPrefix = "goth_session_" // + provider name
	gothStateKeyPrefix   = "goth_state_"   // + provider name
)

// InitializeOAuthProviders sets up OAuth providers (Google and Facebook)
func InitializeOAuthProviders() error {
	googleKey := os.Getenv("GOOGLE_KEY")
	if googleKey == "" {
		log.Fatal("Error loading the GOOGLE_KEY from environment variables")
	}

	googleSecret := os.Getenv("GOOGLE_SECRET")
	if googleSecret == "" {
		log.Fatal("Error loading the GOOGLE_SECRET from environment variables.")
	}

	facebookKey := os.Getenv("FACEBOOK_KEY")
	if facebookKey == "" {
		log.Fatal("Error loading the FACEBOOK_KEY from environment variables")
	}

	facebookSecret := os.Getenv("FACEBOOK_SECRET")
	if facebookSecret == "" {
		log.Fatal("Error loading the FACEBOOK_SECRET from environment variables.")
	}

	// Use SvelteKit backend as callback URL for better security
	googleCallbackURL := os.Getenv("GOOGLE_CALLBACK_URL")
	if googleCallbackURL == "" {
		googleCallbackURL = "http://localhost:5173/api/auth/callback/google"
	}

	facebookCallbackURL := os.Getenv("FACEBOOK_CALLBACK_URL")
	if facebookCallbackURL == "" {
		facebookCallbackURL = "http://localhost:5173/api/auth/callback/facebook"
	}

	goth.UseProviders(
		google.New(googleKey, googleSecret, googleCallbackURL, "email", "profile"),
		facebook.New(facebookKey, facebookSecret, facebookCallbackURL),
	)

	return nil
}

// SetupAuthRoutes configures all authentication routes
// Note: We're now using Supabase for OAuth, so the old Goth-based OAuth routes are disabled
func SetupAuthRoutes(app *fiber.App, store *session.Store) {
	// IMPORTANT: Register specific routes BEFORE wildcard routes

	// ROUTE: get user roles -> /auth/roles (must be before :provider wildcard)
	app.Get("/auth/roles", HandleGetUserRoles)

	// ROUTE: get user roles by email -> /auth/roles/by-email (must be before :provider wildcard)
	app.Get("/auth/roles/by-email", HandleGetUserRolesByEmail)

	// ROUTE: logout -> /auth/logout/:userId
	app.Delete("/auth/logout/:userId", handleLogout)

	// ROUTE: check session -> /auth/session/:userId
	app.Get("/auth/session/:userId", handleCheckSession)

	// ROUTE: sync Supabase user to local PostgreSQL -> /auth/sync
	app.Post("/auth/sync", HandleSupabaseUserSync)

	// OLD OAUTH ROUTES (DISABLED - Using Supabase OAuth instead)
	// These routes are kept for reference but not registered
	// If you need to revert to the old OAuth flow, uncomment these:
	// app.Get("/auth/:provider", handleAuthStart(store))
	// app.Get("/auth/:provider/callback", handleAuthCallback(store))
}

// handleAuthStart initiates the OAuth flow
func handleAuthStart(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		providerName := c.Params("provider")
		provider, err := goth.GetProvider(providerName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid provider")
		}

		// create state and begin auth (provider returns a Session)
		state := uuid.New().String()
		sess, err := provider.BeginAuth(state)
		if err != nil {
			log.Println("BeginAuth error:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Authentication begin failed")
		}

		// get auth URL to redirect user to provider's consent page
		authURL, err := sess.GetAuthURL()
		if err != nil {
			log.Println("GetAuthURL error:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to get auth URL")
		}

		// store marshalled session and state in fiber session store
		fsess, ferr := store.Get(c)
		if ferr != nil {
			log.Println("session get error:", ferr)
			return c.Status(fiber.StatusInternalServerError).SendString("Session error")
		}

		// Save session string (to retrieve on callback)
		sessionKey := gothSessionKeyPrefix + providerName
		stateKey := gothStateKeyPrefix + providerName

		fsess.Set(sessionKey, sess.Marshal())
		fsess.Set(stateKey, state)
		if err := fsess.Save(); err != nil {
			log.Println("session save error:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Session save error")
		}

		// Redirect to the provider's auth URL
		return c.Redirect(authURL, fiber.StatusTemporaryRedirect)
	}
}

// handleAuthCallback completes the OAuth flow and returns user info
func handleAuthCallback(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		providerName := c.Params("provider")
		provider, err := goth.GetProvider(providerName)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid provider")
		}

		// get saved session and state from fiber session store
		fsess, ferr := store.Get(c)
		if ferr != nil {
			log.Println("session get error:", ferr)
			return c.Status(fiber.StatusInternalServerError).SendString("Session error")
		}
		sessionKey := gothSessionKeyPrefix + providerName
		stateKey := gothStateKeyPrefix + providerName

		rawSession := fsess.Get(sessionKey)
		rawState := fsess.Get(stateKey)
		if rawSession == nil {
			return c.Status(fiber.StatusBadRequest).SendString("No auth session found (start auth flow first)")
		}

		// optional: verify state param to avoid CSRF
		receivedState := c.Query("state")
		if rawState != nil {
			if s, ok := rawState.(string); ok {
				if receivedState == "" || receivedState != s {
					return c.Status(fiber.StatusBadRequest).SendString("Invalid state parameter")
				}
			}
		}

		// Recreate session object from marshalled string
		sessionStr, ok := rawSession.(string)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("Invalid stored session data")
		}

		// Unmarshal to a Session object for provider
		sess, err := provider.UnmarshalSession(sessionStr)
		if err != nil {
			log.Println("unmarshal session error:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to unmarshal session")
		}

		// Convert query string into url.Values for Authorize
		qs := string(c.Request().URI().QueryString())
		values, _ := url.ParseQuery(qs)

		// Authorize the session using callback query (code/ oauth_token etc.)
		if _, err := sess.Authorize(provider, values); err != nil {
			// some providers may have already authorized; attempt to fetch user anyway
			log.Println("session authorize error:", err)
			// continue to FetchUser; often malformed authorize causes fetch to fail too, so will be handled below
		}

		// Fetch user using the session
		user, err := provider.FetchUser(sess)
		if err != nil {
			// If the provider requires a token exchange that wasn't completed, this may fail.
			// Log and return an error to the client.
			log.Println("fetch user error:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch user: " + err.Error())
		}

		// OPTIONAL: clear stored session now that we've completed auth
		fsess.Delete(sessionKey)
		fsess.Delete(stateKey)
		if err := fsess.Save(); err != nil {
			log.Println("session save error:", err)
		}

		// Persist user to database (create or update)
		// Note: We don't store OAuth tokens - only using them for authentication
		dbUser, err := UpsertUserWithOAuth(
			c.Context(),
			user.Email,
			user.Name,
			user.AvatarURL,
			providerName,
			user.UserID,
		)
		if err != nil {
			log.Println("failed to persist user to database:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save user: " + err.Error())
		}

		// Create a session record in Redis for admin tracking
		redisSessionKey := fmt.Sprintf("session:%s", dbUser.ID.String())
		sessionData := map[string]interface{}{
			"user_id":   dbUser.ID.String(),
			"email":     dbUser.Email,
			"name":      dbUser.Name,
			"provider":  providerName,
			"login_at":  time.Now().Format(time.RFC3339),
		}
		sessionJSON, _ := json.Marshal(sessionData)
		if err := database.RedisClient.Set(c.Context(), redisSessionKey, sessionJSON, 24*time.Hour).Err(); err != nil {
			log.Println("failed to create session in Redis:", err)
			// Don't fail the auth flow if session creation fails
		}

		// Encode user data as URL parameters for safe transmission
		userDataJSON, _ := json.Marshal(map[string]interface{}{
			"provider":   providerName,
			"name":       dbUser.Name,
			"email":      dbUser.Email,
			"user_id":    dbUser.ID.String(),
			"avatar_url": dbUser.AvatarURL,
		})

		// Redirect to SvelteKit callback with user data as query parameter
		redirectURL := "http://localhost:5173/api/auth/callback/" + providerName + "?user=" + url.QueryEscape(string(userDataJSON))
		return c.Redirect(redirectURL, fiber.StatusTemporaryRedirect)
	}
}

// NewSessionStore creates and returns a new fiber session store with Redis storage
func NewSessionStore() *session.Store {
	// Use the existing Redis client from database package
	// This avoids creating duplicate Redis connections
	storage := redis.New(redis.Config{
		Host:      os.Getenv("REDIS_ADDR"),
		Port:      0, // Port is included in Host string
		Password:  os.Getenv("REDIS_PASSWORD"),
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
	})

	// Create session store with Redis backend
	return session.New(session.Config{
		Storage:        storage,
		KeyLookup:      "cookie:session_id",
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
		CookiePath:     "/",
		Expiration:     24 * time.Hour,
	})
}

// handleLogout deletes the user's Redis session
func handleLogout(c *fiber.Ctx) error {
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	// Delete the session from Redis
	sessionKey := fmt.Sprintf("session:%s", userID)
	if err := database.RedisClient.Del(c.Context(), sessionKey).Err(); err != nil {
		log.Println("failed to delete session from Redis:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete session",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

// handleCheckSession validates if a user's Redis session exists
func handleCheckSession(c *fiber.Ctx) error {
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	// Check if the session exists in Redis
	sessionKey := fmt.Sprintf("session:%s", userID)
	exists, err := database.RedisClient.Exists(c.Context(), sessionKey).Result()
	if err != nil {
		log.Println("failed to check session in Redis:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check session",
		})
	}

	if exists == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Session not found",
		})
	}

	return c.JSON(fiber.Map{
		"valid": true,
	})
}
