package budget

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetCategoriesHandler returns all categories for the authenticated user
func GetCategoriesHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Check for optional type filter
	categoryType := c.Query("type")
	var categories []Category
	var err error

	if categoryType != "" {
		// Validate category type
		if categoryType != "income" && categoryType != "expense" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid category type. Must be 'income' or 'expense'",
			})
		}
		categories, err = GetCategoriesByType(c.Context(), userID, categoryType)
	} else {
		categories, err = GetCategoriesByUserID(c.Context(), userID)
	}

	if err != nil {
		log.Printf("Error fetching categories for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch categories",
		})
	}

	// Return empty array instead of null if no categories
	if categories == nil {
		categories = []Category{}
	}

	return c.JSON(fiber.Map{
		"categories": categories,
	})
}

// GetCategoryHandler returns a specific category by ID
func GetCategoryHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	categoryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID",
		})
	}

	category, err := GetCategoryByID(c.Context(), categoryID, userID)
	if err != nil {
		if err.Error() == "category not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Category not found",
			})
		}
		log.Printf("Error fetching category %s: %v", categoryID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch category",
		})
	}

	return c.JSON(category)
}

// CreateCategoryHandler creates a new category
func CreateCategoryHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category name is required",
		})
	}
	if req.CategoryType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category type is required",
		})
	}

	// Validate category type
	if req.CategoryType != "income" && req.CategoryType != "expense" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category type. Must be 'income' or 'expense'",
		})
	}

	// Validate parent category if provided
	if req.ParentCategoryID != nil {
		parentCategory, err := GetCategoryByID(c.Context(), *req.ParentCategoryID, userID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Parent category not found",
			})
		}
		// Ensure parent category is the same type
		if parentCategory.CategoryType != req.CategoryType {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Parent category must be the same type (income or expense)",
			})
		}
	}

	category, err := CreateCategory(c.Context(), userID, req)
	if err != nil {
		log.Printf("Error creating category for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create category",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}

// UpdateCategoryHandler updates an existing category
func UpdateCategoryHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	categoryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID",
		})
	}

	var req UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate category type if provided
	if req.CategoryType != nil {
		if *req.CategoryType != "income" && *req.CategoryType != "expense" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid category type. Must be 'income' or 'expense'",
			})
		}
	}

	// Validate parent category if provided
	if req.ParentCategoryID != nil {
		parentCategory, err := GetCategoryByID(c.Context(), *req.ParentCategoryID, userID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Parent category not found",
			})
		}
		// Get the current category to check type compatibility
		currentCategory, err := GetCategoryByID(c.Context(), categoryID, userID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Category not found",
			})
		}
		// Ensure parent category is the same type
		if parentCategory.CategoryType != currentCategory.CategoryType {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Parent category must be the same type (income or expense)",
			})
		}
	}

	category, err := UpdateCategory(c.Context(), categoryID, userID, req)
	if err != nil {
		if err.Error() == "category not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Category not found",
			})
		}
		log.Printf("Error updating category %s: %v", categoryID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update category",
		})
	}

	return c.JSON(category)
}

// DeleteCategoryHandler deletes a category (soft delete)
func DeleteCategoryHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	categoryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID",
		})
	}

	err = DeleteCategory(c.Context(), categoryID, userID)
	if err != nil {
		if err.Error() == "category not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Category not found",
			})
		}
		log.Printf("Error deleting category %s: %v", categoryID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete category",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Category deleted successfully",
	})
}

// SeedDefaultCategoriesHandler seeds default categories for the authenticated user
func SeedDefaultCategoriesHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Check if user already has categories
	existingCategories, err := GetCategoriesByUserID(c.Context(), userID)
	if err != nil {
		log.Printf("Error checking existing categories for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check existing categories",
		})
	}

	if len(existingCategories) > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   "User already has categories",
			"message": "Default categories can only be seeded for users with no existing categories",
		})
	}

	err = SeedDefaultCategories(c.Context(), userID)
	if err != nil {
		log.Printf("Error seeding default categories for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to seed default categories",
		})
	}

	// Return the newly created categories
	categories, err := GetCategoriesByUserID(c.Context(), userID)
	if err != nil {
		log.Printf("Error fetching seeded categories for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Categories seeded but failed to fetch them",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Default categories seeded successfully",
		"categories": categories,
	})
}
