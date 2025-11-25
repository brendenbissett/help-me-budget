package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSession "github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
)

const (
	// session key names stored in fiber session
	gothSessionKeyPrefix = "goth_session_" // + provider name
	gothStateKeyPrefix   = "goth_state_"   // + provider name
)

func main() {

	API_VERSION := "0.0.1"

	// Load .env file, oly if in non-production environment
	//if os.Getenv("APP_ENV") != "production" {
	//	err := godotenv.Load()
	//	if err != nil {
	//		log.Fatal("Error loading .env file")
	//	}
	//}

	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Load environment variables
	GOOGLE_KEY := os.Getenv("GOOGLE_KEY")
	if GOOGLE_KEY == "" {
		log.Fatal("Error loading the GOOGLE_KEY from environment variables")
	}

	GOOGLE_SECRET := os.Getenv("GOOGLE_SECRET")
	if GOOGLE_SECRET == "" {
		log.Fatal("Error loading the GOOGLE_SECRET from environment variables.")
	}

	FACEBOOK_KEY := os.Getenv("FACEBOOK_KEY")
	if FACEBOOK_KEY == "" {
		log.Fatal("Error loading the FACEBOOK_KEY from environment variables")
	}

	FACEBOOK_SECRET := os.Getenv("FACEBOOK_SECRET")
	if FACEBOOK_SECRET == "" {
		log.Fatal("Error loading the FACEBOOK_SECRET from environment variables.")
	}

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:3000/auth/google/callback", "email", "profile"),
		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), "http://localhost:3000/auth/facebook/callback"),
	)

	// TODO: Load auth

	app := fiber.New()

	// Set up middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	store := fiberSession.New(fiberSession.Config{
		CookieHTTPOnly: true,
		Expiration:     24 * time.Hour,
	})
	app.Get("/", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Welcome to Help-Me-Budget API (V %s)", API_VERSION)
		return c.SendString(msg)
	})

	// ROUTE: start OAuth flow -> /auth/:provider
	app.Get("/auth/:provider", func(c *fiber.Ctx) error {
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
	})

	// ROUTE: callback -> /auth/:provider/callback
	app.Get("/auth/:provider/callback", func(c *fiber.Ctx) error {
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

		// For demo: return a simple JSON of the fetched user
		return c.JSON(fiber.Map{
			"provider":      providerName,
			"name":          user.Name,
			"email":         user.Email,
			"user_id":       user.UserID,
			"avatar_url":    user.AvatarURL,
			"access_token":  user.AccessToken,
			"refresh_token": user.RefreshToken,
			// do NOT include RefreshToken or AccessToken raw in production responses
		})
	})

	log.Fatal(app.Listen(":3000"))
}
