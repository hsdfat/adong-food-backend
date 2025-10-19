package models

import "time"

// SupplierPrice - Supplier price list (banggia_ncc)
type SupplierPrice struct {
	ProductID     int        `gorm:"primaryKey;autoIncrement;column:sanphamid" json:"productId"`
	ProductName   string     `gorm:"column:tensanpham" json:"productName"`
	IngredientID  string     `gorm:"column:nguyenlieuid" json:"ingredientId"`
	Category      string     `gorm:"column:phanloai" json:"category"`
	SupplierID    string     `gorm:"column:nhacungcapid" json:"supplierId"`
	Manufacturer  string     `gorm:"column:tencososx" json:"manufacturer"`
	Unit          string     `gorm:"column:donvitinh" json:"unit"`
	Specification string     `gorm:"column:quycach" json:"specification"`
	UnitPrice     float64    `gorm:"column:dongia;type:decimal(15,2)" json:"unitPrice"`
	PricePer1     float64    `gorm:"column:dongia1sp;type:decimal(15,2)" json:"pricePer1"`
	EffectiveFrom *time.Time `gorm:"column:hieuluctu" json:"effectiveFrom"`
	EffectiveTo   *time.Time `gorm:"column:hieulucden" json:"effectiveTo"`
	Active        *bool      `gorm:"column:active;default:true" json:"active"`
	NewPrice      float64    `gorm:"column:giakimoi;type:decimal(15,2)" json:"newPrice"`
	Promotion     string     `gorm:"column:khuyenmai;type:char(1)" json:"promotion"`

	// Relationships
	Ingredient *Ingredient `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
	Supplier   *Supplier   `gorm:"foreignKey:SupplierID;references:SupplierID" json:"supplier,omitempty"`
}

func (SupplierPrice) TableName() string {
	return "banggia_ncc"
}
