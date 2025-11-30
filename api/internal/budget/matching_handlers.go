package budget

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetSuggestedMatchesHandler returns potential matches for a transaction
func GetSuggestedMatchesHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	transactionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction ID",
		})
	}

	// Get transaction
	transaction, err := GetTransactionByID(c.Context(), transactionID, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	// Get match suggestions
	suggestions, err := SuggestMatches(c.Context(), transaction, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate match suggestions",
		})
	}

	return c.JSON(fiber.Map{
		"transaction": transaction,
		"suggestions": suggestions,
	})
}

// AutoMatchTransactionHandler attempts to auto-match a single transaction
func AutoMatchTransactionHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	transactionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction ID",
		})
	}

	// Attempt auto-match
	transaction, err := AutoMatchTransaction(c.Context(), transactionID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to auto-match transaction",
		})
	}

	matched := transaction.BudgetEntryID != nil

	return c.JSON(fiber.Map{
		"transaction": transaction,
		"matched":     matched,
	})
}

// BulkAutoMatchHandler attempts to auto-match all unmatched transactions
func BulkAutoMatchHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Run bulk auto-match
	matchedCount, err := BulkAutoMatch(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to bulk auto-match transactions",
		})
	}

	return c.JSON(fiber.Map{
		"matched_count": matchedCount,
		"message":       "Auto-match completed",
	})
}

// UpdateBudgetEntryMatchingRulesHandler updates matching rules for a budget entry
func UpdateBudgetEntryMatchingRulesHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	budgetID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid budget ID",
		})
	}

	entryID, err := uuid.Parse(c.Params("entryId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid entry ID",
		})
	}

	var req struct {
		MatchingRules map[string]interface{} `json:"matching_rules"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update budget entry with matching rules
	updateReq := UpdateBudgetEntryRequest{
		MatchingRules: req.MatchingRules,
	}

	entry, err := UpdateBudgetEntry(c.Context(), budgetID, entryID, userID, updateReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update matching rules",
		})
	}

	return c.JSON(entry)
}

// TeachMatchHandler links a transaction to a budget entry AND creates matching rules
func TeachMatchHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	transactionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction ID",
		})
	}

	var req struct {
		BudgetEntryID   uuid.UUID `json:"budget_entry_id"`
		CreateRules     bool      `json:"create_rules"`
		AmountTolerance float64   `json:"amount_tolerance"` // Optional
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get transaction
	transaction, err := GetTransactionByID(c.Context(), transactionID, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	// Link transaction
	updated, err := LinkTransactionToBudgetEntry(c.Context(), transactionID, userID, req.BudgetEntryID, "manual")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to link transaction",
		})
	}

	// If create_rules is true, extract matching rules from transaction
	if req.CreateRules && transaction.Description != nil {
		// Get budget entry to find budget ID
		entry, err := GetBudgetEntryByID(c.Context(), req.BudgetEntryID, userID)
		if err == nil {
			// Create matching rules based on transaction
			rules := make(map[string]interface{})

			// Add description pattern
			desc := *transaction.Description
			rules["description_contains"] = []string{desc}

			// Add amount tolerance if specified
			if req.AmountTolerance > 0 {
				rules["amount_tolerance"] = req.AmountTolerance
			} else {
				rules["amount_tolerance"] = 2.0 // Default $2 tolerance
			}

			// Update budget entry with rules
			updateReq := UpdateBudgetEntryRequest{
				MatchingRules: rules,
			}

			_, _ = UpdateBudgetEntry(c.Context(), entry.BudgetID, req.BudgetEntryID, userID, updateReq)
		}
	}

	return c.JSON(fiber.Map{
		"transaction":   updated,
		"rules_created": req.CreateRules,
	})
}

// GetBudgetEntryByID helper function (assuming it exists or needs to be created)
func GetBudgetEntryByID(ctx context.Context, entryID, userID uuid.UUID) (*BudgetEntry, error) {
	// This would need to be implemented in budget_repository.go if it doesn't exist
	// For now, we'll use a workaround by getting all entries and filtering
	// In production, you'd want a dedicated function

	// Get all budgets
	budgets, err := GetBudgetsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Search through budget entries
	for _, budget := range budgets {
		entries, err := GetBudgetEntriesByBudgetID(ctx, budget.ID, userID)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.ID == entryID {
				return &entry, nil
			}
		}
	}

	return nil, fmt.Errorf("budget entry not found")
}
