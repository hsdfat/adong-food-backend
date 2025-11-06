package models

import "time"

// SupplierPrice - Supplier price list (supplier_price_list)
type SupplierPrice struct {
	ProductID     int        `gorm:"primaryKey;autoIncrement;column:product_id" json:"productId"`
	ProductName   string     `gorm:"column:product_name" json:"productName"`
	IngredientID  string     `gorm:"column:ingredient_id" json:"ingredientId"`
	Category      string     `gorm:"column:classification" json:"category"`
	SupplierID    string     `gorm:"column:supplier_id" json:"supplierId"`
	Manufacturer  string     `gorm:"column:manufacturer_name" json:"manufacturer"`
	Unit          string     `gorm:"column:unit" json:"unit"`
	Specification string     `gorm:"column:specification" json:"specification"`
	UnitPrice     float64    `gorm:"column:unit_price;type:decimal(15,2)" json:"unitPrice"`
	PricePer1     float64    `gorm:"column:price_per_item;type:decimal(15,2)" json:"pricePer1"`
	EffectiveFrom *time.Time `gorm:"column:effective_from" json:"effectiveFrom"`
	EffectiveTo   *time.Time `gorm:"column:effective_to" json:"effectiveTo"`
	Active        *bool      `gorm:"column:active;default:true" json:"active"`
	NewPrice      float64    `gorm:"column:new_buying_price;type:decimal(15,2)" json:"newPrice"`
	Promotion     string     `gorm:"column:promotion;type:char(1)" json:"promotion"`
	CreatedDate   time.Time  `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate  time.Time  `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Ingredient *Ingredient `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
	Supplier   *Supplier   `gorm:"foreignKey:SupplierID;references:SupplierID" json:"supplier,omitempty"`
}

func (SupplierPrice) TableName() string {
	return "supplier_price_list"
}
