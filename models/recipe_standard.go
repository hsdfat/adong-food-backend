package models

import "time"

// RecipeStandard - Bill of materials for dishes (dinhmuc_monan)
type RecipeStandard struct {
	StandardID   int       `gorm:"primaryKey;autoIncrement;column:dinhmucid" json:"standardId"`
	DishID       string    `gorm:"column:monanid" json:"dishId"`
	IngredientID string    `gorm:"column:nguyenlieuid" json:"ingredientId"`
	Unit         string    `gorm:"column:donvitinh" json:"unit"`
	StandardPer1 float64   `gorm:"column:dinhmuc1suat;type:decimal(10,4)" json:"standardPer1"`
	Note         string    `gorm:"column:ghichu;type:text" json:"note"`
	Amount       float64   `gorm:"column:thanhtien;type:decimal(15,2)" json:"amount"`
	UpdatedByID  string    `gorm:"column:nguoicapnhatid" json:"updatedById"`
	CreatedDate  time.Time `gorm:"column:createddate;autoCreateTime" json:"createdDate"`
	ModifiedDate time.Time `gorm:"column:modifieddate;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Dish       *Dish       `gorm:"foreignKey:DishID;references:DishID" json:"dish,omitempty"`
	Ingredient *Ingredient `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
	UpdatedBy  *User       `gorm:"foreignKey:UpdatedByID;references:UserID" json:"updatedBy,omitempty"`
}

func (RecipeStandard) TableName() string {
	return "dinhmuc_monan"
}
