package budget

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetBudgetsHandler returns all budgets for the authenticated user
func GetBudgetsHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	budgets, err := GetBudgetsByUserID(c.Context(), userID)
	if err != nil {
		log.Printf("Error fetching budgets for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch budgets",
		})
	}

	if budgets == nil {
		budgets = []Budget{}
	}

	return c.JSON(fiber.Map{
		"budgets": budgets,
	})
}

// GetBudgetHandler returns a specific budget by ID
func GetBudgetHandler(c *fiber.Ctx) error {
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

	budget, err := GetBudgetByID(c.Context(), budgetID, userID)
	if err != nil {
		if err.Error() == "budget not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Budget not found",
			})
		}
		log.Printf("Error fetching budget %s: %v", budgetID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch budget",
		})
	}

	return c.JSON(budget)
}

// GetBudgetWithEntriesHandler returns a budget with all its entries
func GetBudgetWithEntriesHandler(c *fiber.Ctx) error {
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

	budgetWithEntries, err := GetBudgetWithEntries(c.Context(), budgetID, userID)
	if err != nil {
		if err.Error() == "budget not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Budget not found",
			})
		}
		log.Printf("Error fetching budget with entries %s: %v", budgetID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch budget",
		})
	}

	return c.JSON(budgetWithEntries)
}

// CreateBudgetHandler creates a new budget
func CreateBudgetHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req CreateBudgetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Budget name is required",
		})
	}

	budget, err := CreateBudget(c.Context(), userID, req)
	if err != nil {
		log.Printf("Error creating budget for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create budget",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(budget)
}

// UpdateBudgetHandler updates an existing budget
func UpdateBudgetHandler(c *fiber.Ctx) error {
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

	var req UpdateBudgetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	budget, err := UpdateBudget(c.Context(), budgetID, userID, req)
	if err != nil {
		if err.Error() == "budget not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Budget not found",
			})
		}
		log.Printf("Error updating budget %s: %v", budgetID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update budget",
		})
	}

	return c.JSON(budget)
}

// DeleteBudgetHandler deletes a budget (soft delete)
func DeleteBudgetHandler(c *fiber.Ctx) error {
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

	err = DeleteBudget(c.Context(), budgetID, userID)
	if err != nil {
		if err.Error() == "budget not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Budget not found",
			})
		}
		log.Printf("Error deleting budget %s: %v", budgetID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete budget",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Budget deleted successfully",
	})
}

// GetBudgetEntriesHandler returns all entries for a budget
func GetBudgetEntriesHandler(c *fiber.Ctx) error {
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

	entries, err := GetBudgetEntriesByBudgetID(c.Context(), budgetID, userID)
	if err != nil {
		log.Printf("Error fetching budget entries for budget %s: %v", budgetID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch budget entries",
		})
	}

	if entries == nil {
		entries = []BudgetEntry{}
	}

	return c.JSON(fiber.Map{
		"entries": entries,
	})
}

// CreateBudgetEntryHandler creates a new budget entry
func CreateBudgetEntryHandler(c *fiber.Ctx) error {
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

	var req CreateBudgetEntryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Entry name is required",
		})
	}
	if req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Amount must be greater than zero",
		})
	}
	if req.EntryType != "income" && req.EntryType != "expense" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Entry type must be 'income' or 'expense'",
		})
	}

	validFrequencies := map[string]bool{
		"once_off":    true,
		"daily":       true,
		"weekly":      true,
		"fortnightly": true,
		"monthly":     true,
		"annually":    true,
	}
	if !validFrequencies[req.Frequency] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid frequency",
		})
	}

	entry, err := CreateBudgetEntry(c.Context(), budgetID, userID, req)
	if err != nil {
		log.Printf("Error creating budget entry for budget %s: %v", budgetID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create budget entry",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(entry)
}

// UpdateBudgetEntryHandler updates an existing budget entry
func UpdateBudgetEntryHandler(c *fiber.Ctx) error {
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

	var req UpdateBudgetEntryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	entry, err := UpdateBudgetEntry(c.Context(), entryID, budgetID, userID, req)
	if err != nil {
		if err.Error() == "budget entry not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Budget entry not found",
			})
		}
		log.Printf("Error updating budget entry %s: %v", entryID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update budget entry",
		})
	}

	return c.JSON(entry)
}

// DeleteBudgetEntryHandler deletes a budget entry (soft delete)
func DeleteBudgetEntryHandler(c *fiber.Ctx) error {
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

	err = DeleteBudgetEntry(c.Context(), entryID, budgetID, userID)
	if err != nil {
		if err.Error() == "budget entry not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Budget entry not found",
			})
		}
		log.Printf("Error deleting budget entry %s: %v", entryID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete budget entry",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Budget entry deleted successfully",
	})
}

// GetBudgetSummaryHandler returns summary statistics for a budget
func GetBudgetSummaryHandler(c *fiber.Ctx) error {
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

	summary, err := CalculateBudgetSummary(c.Context(), budgetID, userID)
	if err != nil {
		log.Printf("Error calculating budget summary for budget %s: %v", budgetID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate budget summary",
		})
	}

	// Add health status
	healthStatus := GetBudgetHealthStatus(summary)

	return c.JSON(fiber.Map{
		"summary": summary,
		"health":  healthStatus,
	})
}

// ProjectCashFlowHandler projects future cash flow based on budget
func ProjectCashFlowHandler(c *fiber.Ctx) error {
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

	// Get query parameters
	startingBalance, _ := strconv.ParseFloat(c.Query("starting_balance", "0"), 64)
	days, _ := strconv.Atoi(c.Query("days", "90"))

	// Cap days at 365
	if days > 365 {
		days = 365
	}
	if days < 1 {
		days = 30
	}

	projection, err := ProjectCashFlow(c.Context(), budgetID, userID, startingBalance, days)
	if err != nil {
		log.Printf("Error projecting cash flow for budget %s: %v", budgetID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to project cash flow",
		})
	}

	return c.JSON(projection)
}
