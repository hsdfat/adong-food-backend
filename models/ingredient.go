package models

import "time"

// Ingredient - Master data for raw materials and ingredients (dm_nvl)
type Ingredient struct {
	IngredientID   string    `gorm:"primaryKey;column:nguyenlieuid" json:"ingredientId"`
	IngredientName string    `gorm:"column:tennguyenlieu;not null" json:"ingredientName"`
	Property       string    `gorm:"column:tinhchat" json:"property"`
	MaterialGroup  string    `gorm:"column:nhomvthh" json:"materialGroup"`
	Unit           string    `gorm:"column:donvitinh" json:"unit"`
	CreatedDate    time.Time `gorm:"column:createddate;autoCreateTime" json:"createdDate"`
	ModifiedDate   time.Time `gorm:"column:modifieddate;autoUpdateTime" json:"modifiedDate"`
}

func (Ingredient) TableName() string {
	return "dm_nvl"
}
