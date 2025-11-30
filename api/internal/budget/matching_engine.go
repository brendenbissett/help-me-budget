package budget

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

// MatchSuggestion represents a suggested match between a transaction and budget entry
type MatchSuggestion struct {
	BudgetEntry      BudgetEntry `json:"budget_entry"`
	ConfidenceScore  float64     `json:"confidence_score"` // 0-100
	ConfidenceLevel  string      `json:"confidence_level"` // 'auto_high', 'auto_low'
	MatchReasons     []string    `json:"match_reasons"`
}

// MatchingCriteria defines how to match transactions to budget entries
type MatchingCriteria struct {
	DescriptionContains  []string `json:"description_contains,omitempty"`
	AmountTolerance      float64  `json:"amount_tolerance,omitempty"`      // e.g., 2.00 = plus/minus $2
	TolerancePercentage  float64  `json:"tolerance_percentage,omitempty"`  // e.g., 5.0 = plus/minus 5%
	MerchantName         string   `json:"merchant_name,omitempty"`
	CategoryID           *uuid.UUID `json:"category_id,omitempty"`
	MinAmount            float64  `json:"min_amount,omitempty"`
	MaxAmount            float64  `json:"max_amount,omitempty"`
}

// SuggestMatches finds potential budget entry matches for a transaction
func SuggestMatches(ctx context.Context, transaction *Transaction, userID uuid.UUID) ([]MatchSuggestion, error) {
	// Get active budget
	activeBudget, err := GetActiveBudget(ctx, userID)
	if err != nil || activeBudget == nil {
		return []MatchSuggestion{}, nil // No active budget
	}

	// Get budget entries
	entries, err := GetBudgetEntriesByBudgetID(ctx, activeBudget.ID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget entries: %w", err)
	}

	// Filter to same transaction type (income/expense)
	var candidates []BudgetEntry
	for _, entry := range entries {
		if entry.IsActive && entry.EntryType == transaction.TransactionType {
			candidates = append(candidates, entry)
		}
	}

	// Score each candidate
	var suggestions []MatchSuggestion
	for _, entry := range candidates {
		score, reasons := scoreMatch(transaction, &entry)

		if score > 0 {
			confidenceLevel := "auto_low"
			if score >= 70 {
				confidenceLevel = "auto_high"
			}

			suggestions = append(suggestions, MatchSuggestion{
				BudgetEntry:     entry,
				ConfidenceScore: score,
				ConfidenceLevel: confidenceLevel,
				MatchReasons:    reasons,
			})
		}
	}

	// Sort by confidence score (highest first)
	for i := 0; i < len(suggestions); i++ {
		for j := i + 1; j < len(suggestions); j++ {
			if suggestions[j].ConfidenceScore > suggestions[i].ConfidenceScore {
				suggestions[i], suggestions[j] = suggestions[j], suggestions[i]
			}
		}
	}

	return suggestions, nil
}

// scoreMatch calculates match confidence score and reasons
func scoreMatch(transaction *Transaction, entry *BudgetEntry) (float64, []string) {
	score := 0.0
	reasons := []string{}

	// 1. Check matching rules if defined
	if entry.MatchingRules != nil {
		ruleScore, ruleReasons := scoreByRules(transaction, entry)
		score += ruleScore
		reasons = append(reasons, ruleReasons...)
	}

	// 2. Description matching (case-insensitive substring match)
	if transaction.Description != nil && *transaction.Description != "" {
		descScore, descReason := scoreByDescription(transaction, entry)
		score += descScore
		if descReason != "" {
			reasons = append(reasons, descReason)
		}
	}

	// 3. Amount matching
	amountScore, amountReason := scoreByAmount(transaction, entry)
	score += amountScore
	if amountReason != "" {
		reasons = append(reasons, amountReason)
	}

	// 4. Category matching
	if transaction.CategoryID != nil && entry.CategoryID != nil {
		if *transaction.CategoryID == *entry.CategoryID {
			score += 20
			reasons = append(reasons, "Same category")
		}
	}

	// 5. Frequency/timing matching
	if entry.Frequency != "once_off" {
		timingScore, timingReason := scoreByTiming(transaction, entry)
		score += timingScore
		if timingReason != "" {
			reasons = append(reasons, timingReason)
		}
	}

	// Cap at 100
	if score > 100 {
		score = 100
	}

	return score, reasons
}

// scoreByRules evaluates matching rules from JSONB field
func scoreByRules(transaction *Transaction, entry *BudgetEntry) (float64, []string) {
	score := 0.0
	reasons := []string{}

	rules := entry.MatchingRules

	// Check description_contains
	if descContains, ok := rules["description_contains"].([]interface{}); ok {
		if transaction.Description != nil {
			transDesc := strings.ToLower(*transaction.Description)
			for _, pattern := range descContains {
				if patternStr, ok := pattern.(string); ok {
					if strings.Contains(transDesc, strings.ToLower(patternStr)) {
						score += 30
						reasons = append(reasons, fmt.Sprintf("Description contains '%s'", patternStr))
						break
					}
				}
			}
		}
	}

	// Check merchant_name
	if merchantName, ok := rules["merchant_name"].(string); ok {
		if transaction.Description != nil {
			if strings.Contains(strings.ToLower(*transaction.Description), strings.ToLower(merchantName)) {
				score += 25
				reasons = append(reasons, fmt.Sprintf("Merchant name: %s", merchantName))
			}
		}
	}

	// Check amount_tolerance
	if amountTol, ok := rules["amount_tolerance"].(float64); ok {
		diff := math.Abs(transaction.Amount - entry.Amount)
		if diff <= amountTol {
			score += 20
			reasons = append(reasons, fmt.Sprintf("Amount within $%.2f", amountTol))
		}
	}

	return score, reasons
}

// scoreByDescription checks if transaction description matches entry name
func scoreByDescription(transaction *Transaction, entry *BudgetEntry) (float64, string) {
	if transaction.Description == nil {
		return 0, ""
	}

	transDesc := strings.ToLower(*transaction.Description)
	entryName := strings.ToLower(entry.Name)

	// Exact match
	if transDesc == entryName {
		return 40, "Exact description match"
	}

	// Partial match (contains)
	if strings.Contains(transDesc, entryName) || strings.Contains(entryName, transDesc) {
		return 25, "Partial description match"
	}

	// Word-based matching (at least 2 common words)
	transWords := strings.Fields(transDesc)
	entryWords := strings.Fields(entryName)

	commonWords := 0
	for _, tw := range transWords {
		for _, ew := range entryWords {
			if len(tw) > 3 && tw == ew { // Only count words longer than 3 chars
				commonWords++
			}
		}
	}

	if commonWords >= 2 {
		return 15, fmt.Sprintf("%d common words", commonWords)
	}

	return 0, ""
}

// scoreByAmount checks if transaction amount matches entry amount
func scoreByAmount(transaction *Transaction, entry *BudgetEntry) (float64, string) {
	diff := math.Abs(transaction.Amount - entry.Amount)

	// Exact match
	if diff < 0.01 {
		return 30, "Exact amount match"
	}

	// Within $2
	if diff <= 2.00 {
		return 20, fmt.Sprintf("Amount within $%.2f", diff)
	}

	// Within 5%
	tolerance := entry.Amount * 0.05
	if diff <= tolerance {
		return 15, "Amount within 5%"
	}

	// Within $10
	if diff <= 10.00 {
		return 5, "Amount within $10"
	}

	return 0, ""
}

// scoreByTiming checks if transaction date aligns with budget entry frequency
func scoreByTiming(transaction *Transaction, entry *BudgetEntry) (float64, string) {
	transDate, err := time.Parse("2006-01-02", transaction.TransactionDate)
	if err != nil {
		return 0, ""
	}

	startDate, err := time.Parse("2006-01-02", entry.StartDate)
	if err != nil {
		return 0, ""
	}

	// Check if transaction is after start date
	if transDate.Before(startDate) {
		return 0, ""
	}

	// Check if before end date (if set)
	if entry.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *entry.EndDate)
		if err == nil && transDate.After(endDate) {
			return 0, ""
		}
	}

	switch entry.Frequency {
	case "monthly":
		// Check if day of month matches
		if entry.DayOfMonth != nil && transDate.Day() == *entry.DayOfMonth {
			return 15, "Matches monthly schedule"
		}
		// Allow 3 days tolerance for monthly
		if entry.DayOfMonth != nil {
			dayDiff := math.Abs(float64(transDate.Day() - *entry.DayOfMonth))
			if dayDiff <= 3 {
				return 10, "Close to monthly schedule"
			}
		}

	case "weekly":
		// Check if day of week matches
		if entry.DayOfWeek != nil && int(transDate.Weekday()) == *entry.DayOfWeek {
			return 15, "Matches weekly schedule"
		}

	case "fortnightly":
		// Calculate weeks since start
		daysSince := transDate.Sub(startDate).Hours() / 24
		weeksSince := int(daysSince / 7)

		if weeksSince%2 == 0 && entry.DayOfWeek != nil && int(transDate.Weekday()) == *entry.DayOfWeek {
			return 15, "Matches fortnightly schedule"
		}

	case "annually":
		// Check if month and day match
		if transDate.Month() == startDate.Month() && transDate.Day() == startDate.Day() {
			return 15, "Matches annual schedule"
		}
	}

	return 0, ""
}

// AutoMatchTransaction attempts to automatically match a transaction
func AutoMatchTransaction(ctx context.Context, transactionID, userID uuid.UUID) (*Transaction, error) {
	// Get transaction
	transaction, err := GetTransactionByID(ctx, transactionID, userID)
	if err != nil {
		return nil, err
	}

	// Already matched
	if transaction.MatchConfidence == "manual" || transaction.BudgetEntryID != nil {
		return transaction, nil
	}

	// Get suggestions
	suggestions, err := SuggestMatches(ctx, transaction, userID)
	if err != nil {
		return nil, err
	}

	// Auto-link if high confidence match found
	if len(suggestions) > 0 && suggestions[0].ConfidenceScore >= 70 {
		bestMatch := suggestions[0]

		// Link transaction to budget entry
		updated, err := LinkTransactionToBudgetEntry(ctx, transactionID, userID, bestMatch.BudgetEntry.ID, "auto_high")
		if err != nil {
			return nil, err
		}

		return updated, nil
	}

	return transaction, nil
}

// BulkAutoMatch attempts to auto-match multiple unmatched transactions
func BulkAutoMatch(ctx context.Context, userID uuid.UUID) (int, error) {
	// Get unmatched transactions
	unmatched, err := GetUnmatchedTransactions(ctx, userID)
	if err != nil {
		return 0, err
	}

	matchedCount := 0
	for _, transaction := range unmatched {
		matched, err := AutoMatchTransaction(ctx, transaction.ID, userID)
		if err != nil {
			continue // Skip errors, keep processing
		}

		if matched.BudgetEntryID != nil {
			matchedCount++
		}
	}

	return matchedCount, nil
}
