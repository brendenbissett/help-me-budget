package budget

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// getUserIDFromContext extracts the user ID from the request context
func getUserIDFromContext(c *fiber.Ctx) uuid.UUID {
	userIDStr := c.Get("X-User-ID")
	if userIDStr == "" {
		return uuid.Nil
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil
	}

	return userID
}

// GetAccountsHandler returns all accounts for the authenticated user
func GetAccountsHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	accounts, err := GetAccountsByUserID(c.Context(), userID)
	if err != nil {
		log.Printf("Error fetching accounts for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch accounts",
		})
	}

	// Return empty array instead of null if no accounts
	if accounts == nil {
		accounts = []Account{}
	}

	return c.JSON(fiber.Map{
		"accounts": accounts,
	})
}

// GetAccountHandler returns a specific account by ID
func GetAccountHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	accountID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account ID",
		})
	}

	account, err := GetAccountByID(c.Context(), accountID, userID)
	if err != nil {
		if err.Error() == "account not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Account not found",
			})
		}
		log.Printf("Error fetching account %s: %v", accountID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch account",
		})
	}

	return c.JSON(account)
}

// CreateAccountHandler creates a new account
func CreateAccountHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req CreateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Account name is required",
		})
	}
	if req.AccountType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Account type is required",
		})
	}
	if req.Currency == "" {
		req.Currency = "USD" // Default to USD
	}

	// Validate account type
	validTypes := map[string]bool{
		"checking":    true,
		"savings":     true,
		"credit_card": true,
		"cash":        true,
		"investment":  true,
	}
	if !validTypes[req.AccountType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account type. Must be one of: checking, savings, credit_card, cash, investment",
		})
	}

	account, err := CreateAccount(c.Context(), userID, req)
	if err != nil {
		log.Printf("Error creating account for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create account",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(account)
}

// UpdateAccountHandler updates an existing account
func UpdateAccountHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	accountID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account ID",
		})
	}

	var req UpdateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate account type if provided
	if req.AccountType != nil {
		validTypes := map[string]bool{
			"checking":    true,
			"savings":     true,
			"credit_card": true,
			"cash":        true,
			"investment":  true,
		}
		if !validTypes[*req.AccountType] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid account type. Must be one of: checking, savings, credit_card, cash, investment",
			})
		}
	}

	account, err := UpdateAccount(c.Context(), accountID, userID, req)
	if err != nil {
		if err.Error() == "account not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Account not found",
			})
		}
		log.Printf("Error updating account %s: %v", accountID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update account",
		})
	}

	return c.JSON(account)
}

// DeleteAccountHandler deletes an account (soft delete)
func DeleteAccountHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	accountID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account ID",
		})
	}

	err = DeleteAccount(c.Context(), accountID, userID)
	if err != nil {
		if err.Error() == "account not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Account not found",
			})
		}
		log.Printf("Error deleting account %s: %v", accountID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete account",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Account deleted successfully",
	})
}

// GetTotalBalanceHandler returns the total balance across all active accounts
func GetTotalBalanceHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	totalBalance, err := GetTotalBalance(c.Context(), userID)
	if err != nil {
		log.Printf("Error calculating total balance for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate total balance",
		})
	}

	return c.JSON(fiber.Map{
		"total_balance": totalBalance,
		"currency":      "USD", // TODO: Support multiple currencies
	})
}
