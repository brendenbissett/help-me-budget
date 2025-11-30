package budget

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetTransactionsHandler returns all transactions for a user with optional filters
func GetTransactionsHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Parse optional query parameters
	var accountID *uuid.UUID
	if accountIDStr := c.Query("account_id"); accountIDStr != "" {
		parsedID, err := uuid.Parse(accountIDStr)
		if err == nil {
			accountID = &parsedID
		}
	}

	var categoryID *uuid.UUID
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		parsedID, err := uuid.Parse(categoryIDStr)
		if err == nil {
			categoryID = &parsedID
		}
	}

	var startDate *string
	if start := c.Query("start_date"); start != "" {
		startDate = &start
	}

	var endDate *string
	if end := c.Query("end_date"); end != "" {
		endDate = &end
	}

	transactions, err := GetTransactionsByUserID(c.Context(), userID, accountID, categoryID, startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve transactions",
		})
	}

	return c.JSON(fiber.Map{
		"transactions": transactions,
	})
}

// GetTransactionHandler returns a specific transaction by ID
func GetTransactionHandler(c *fiber.Ctx) error {
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

	transaction, err := GetTransactionByID(c.Context(), transactionID, userID)
	if err != nil {
		if err.Error() == "transaction not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Transaction not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve transaction",
		})
	}

	return c.JSON(transaction)
}

// CreateTransactionHandler creates a new transaction
func CreateTransactionHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	transaction, err := CreateTransaction(c.Context(), userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create transaction",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(transaction)
}

// UpdateTransactionHandler updates an existing transaction
func UpdateTransactionHandler(c *fiber.Ctx) error {
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

	var req UpdateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	transaction, err := UpdateTransaction(c.Context(), transactionID, userID, req)
	if err != nil {
		if err.Error() == "transaction not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Transaction not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update transaction",
		})
	}

	return c.JSON(transaction)
}

// DeleteTransactionHandler deletes a transaction
func DeleteTransactionHandler(c *fiber.Ctx) error {
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

	err = DeleteTransaction(c.Context(), transactionID, userID)
	if err != nil {
		if err.Error() == "transaction not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Transaction not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete transaction",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// GetUnmatchedTransactionsHandler returns transactions that haven't been matched to budget entries
func GetUnmatchedTransactionsHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	transactions, err := GetUnmatchedTransactions(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve unmatched transactions",
		})
	}

	return c.JSON(fiber.Map{
		"transactions": transactions,
	})
}

// CategorizeTransactionHandler assigns a category to a transaction
func CategorizeTransactionHandler(c *fiber.Ctx) error {
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
		CategoryID uuid.UUID `json:"category_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	transaction, err := CategorizeTransaction(c.Context(), transactionID, userID, req.CategoryID)
	if err != nil {
		if err.Error() == "transaction not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Transaction not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to categorize transaction",
		})
	}

	return c.JSON(transaction)
}

// LinkTransactionHandler links a transaction to a budget entry
func LinkTransactionHandler(c *fiber.Ctx) error {
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
		MatchConfidence string    `json:"match_confidence"` // 'manual', 'auto_high', 'auto_low'
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Default to manual if not specified
	if req.MatchConfidence == "" {
		req.MatchConfidence = "manual"
	}

	transaction, err := LinkTransactionToBudgetEntry(c.Context(), transactionID, userID, req.BudgetEntryID, req.MatchConfidence)
	if err != nil {
		if err.Error() == "transaction not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Transaction not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to link transaction",
		})
	}

	return c.JSON(transaction)
}
