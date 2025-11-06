package models

import "time"

// RecipeStandard - Bill of materials for dishes (dish_recipe_standards)
type RecipeStandard struct {
	StandardID   int       `gorm:"primaryKey;autoIncrement;column:recipe_id" json:"standardId"`
	DishID       string    `gorm:"column:dish_id" json:"dishId"`
	IngredientID string    `gorm:"column:ingredient_id" json:"ingredientId"`
	Unit         string    `gorm:"column:unit" json:"unit"`
	StandardPer1 float64   `gorm:"column:quantity_per_serving;type:decimal(10,4)" json:"standardPer1"`
	Note         string    `gorm:"column:notes;type:text" json:"note"`
	Amount       float64   `gorm:"column:cost;type:decimal(15,2)" json:"amount"`
	UpdatedByID  string    `gorm:"column:updated_by_user_id" json:"updatedById"`
	CreatedDate  time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Dish       *Dish       `gorm:"foreignKey:DishID;references:DishID" json:"dish,omitempty"`
	Ingredient *Ingredient `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
	UpdatedBy  *User       `gorm:"foreignKey:UpdatedByID;references:UserID" json:"updatedBy,omitempty"`
}

func (RecipeStandard) TableName() string {
	return "dish_recipe_standards"
}
