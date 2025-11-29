package models

import "time"

// IngredientType - Lookup table for ingredient categories
type IngredientType struct {
	IngredientTypeID   string    `gorm:"primaryKey;column:ingredient_type_id" json:"ingredientTypeId"`
	IngredientTypeName string    `gorm:"column:ingredient_type_name;not null;unique" json:"ingredientTypeName"`
	Description        string    `gorm:"column:description" json:"description"`
	Active             bool      `gorm:"column:active;not null;default:true" json:"active"`
	CreatedDate        time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate       time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`
}

func (IngredientType) TableName() string {
	return "ingredient_types"
}

// Ingredient - Master data for raw materials and ingredients (dm_nvl)
type Ingredient struct {
	IngredientID     string          `gorm:"primaryKey;column:ingredient_id" json:"ingredientId"`
	IngredientName   string          `gorm:"column:ingredient_name;not null;unique" json:"ingredientName"`
	IngredientTypeID *string         `gorm:"column:ingredient_type_id" json:"ingredientTypeId"`
	Property         string          `gorm:"column:properties" json:"property"`
	MaterialGroup    string          `gorm:"column:material_group" json:"materialGroup"`
	Unit             string          `gorm:"column:unit;not null" json:"unit"`
	LegacyID         *string         `gorm:"column:legacy_id" json:"legacyId,omitempty"`
	CreatedDate      time.Time       `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate     time.Time       `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`
	IngredientType   *IngredientType `gorm:"foreignKey:IngredientTypeID;references:IngredientTypeID" json:"ingredientType,omitempty"`
}

func (Ingredient) TableName() string {
	return "master_ingredients"
}
