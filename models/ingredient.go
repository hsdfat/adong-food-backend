package models

import "time"

// Ingredient - Master data for raw materials and ingredients (dm_nvl)
type Ingredient struct {
    IngredientID   string    `gorm:"primaryKey;column:ingredient_id" json:"ingredientId"`
    IngredientName string    `gorm:"column:ingredient_name;not null" json:"ingredientName"`
    Property       string    `gorm:"column:properties" json:"property"`
    MaterialGroup  string    `gorm:"column:material_group" json:"materialGroup"`
    Unit           string    `gorm:"column:unit" json:"unit"`
    CreatedDate    time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
    ModifiedDate   time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`
}

func (Ingredient) TableName() string {
    return "master_ingredients"
}
