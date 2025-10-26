package models

import "time"

// RecipeStandardDTO - Data Transfer Object for Recipe Standard with related names
type RecipeStandardDTO struct {
	StandardID     int       `json:"standardId"`
	DishID         string    `json:"dishId"`
	DishName       string    `json:"dishName"`           // Added: Dish name
	IngredientID   string    `json:"ingredientId"`
	IngredientName string    `json:"ingredientName"`     // Added: Ingredient name
	Unit           string    `json:"unit"`
	StandardPer1   float64   `json:"standardPer1"`
	Note           string    `json:"note"`
	Amount         float64   `json:"amount"`
	UpdatedByID    string    `json:"updatedById"`
	UpdatedByName  string    `json:"updatedByName"`      // Added: User name (optional)
	CreatedDate    time.Time `json:"createdDate"`
	ModifiedDate   time.Time `json:"modifiedDate"`
}

// ToDTO converts RecipeStandard model to DTO
func (r *RecipeStandard) ToDTO() RecipeStandardDTO {
	dto := RecipeStandardDTO{
		StandardID:   r.StandardID,
		DishID:       r.DishID,
		IngredientID: r.IngredientID,
		Unit:         r.Unit,
		StandardPer1: r.StandardPer1,
		Note:         r.Note,
		Amount:       r.Amount,
		UpdatedByID:  r.UpdatedByID,
		CreatedDate:  r.CreatedDate,
		ModifiedDate: r.ModifiedDate,
	}

	// Populate names from relationships if available
	if r.Dish != nil {
		dto.DishName = r.Dish.DishName
	}
	if r.Ingredient != nil {
		dto.IngredientName = r.Ingredient.IngredientName
	}
	if r.UpdatedBy != nil {
		dto.UpdatedByName = r.UpdatedBy.FullName
	}

	return dto
}

// ConvertToDTO converts a slice of RecipeStandard to a slice of RecipeStandardDTO
func ConvertRecipeStandardsToDTO(recipes []RecipeStandard) []RecipeStandardDTO {
	dtos := make([]RecipeStandardDTO, len(recipes))
	for i, recipe := range recipes {
		dtos[i] = recipe.ToDTO()
	}
	return dtos
}