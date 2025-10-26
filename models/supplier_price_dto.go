package models

import "time"

// SupplierPriceDTO - Data Transfer Object for Supplier Price with related names
type SupplierPriceDTO struct {
	ProductID        int        `json:"productId"`
	ProductName      string     `json:"productName"`
	IngredientID     string     `json:"ingredientId"`
	IngredientName   string     `json:"ingredientName"`   // Ingredient name from relationship
	Category         string     `json:"category"`
	SupplierID       string     `json:"supplierId"`
	SupplierName     string     `json:"supplierName"`     // Supplier name from relationship
	Manufacturer     string     `json:"manufacturer"`
	Unit             string     `json:"unit"`
	Specification    string     `json:"specification"`
	UnitPrice        float64    `json:"unitPrice"`
	PricePer1        float64    `json:"pricePer1"`
	EffectiveFrom    *time.Time `json:"effectiveFrom"`
	EffectiveTo      *time.Time `json:"effectiveTo"`
	Active           *bool      `json:"active"`
	NewPrice         float64    `json:"newPrice"`
	Promotion        string     `json:"promotion"`
}

// ToDTO converts SupplierPrice model to DTO
func (sp *SupplierPrice) ToDTO() SupplierPriceDTO {
	dto := SupplierPriceDTO{
		ProductID:     sp.ProductID,
		ProductName:   sp.ProductName,
		IngredientID:  sp.IngredientID,
		Category:      sp.Category,
		SupplierID:    sp.SupplierID,
		Manufacturer:  sp.Manufacturer,
		Unit:          sp.Unit,
		Specification: sp.Specification,
		UnitPrice:     sp.UnitPrice,
		PricePer1:     sp.PricePer1,
		EffectiveFrom: sp.EffectiveFrom,
		EffectiveTo:   sp.EffectiveTo,
		Active:        sp.Active,
		NewPrice:      sp.NewPrice,
		Promotion:     sp.Promotion,
	}

	// Populate names from relationships if available
	if sp.Ingredient != nil {
		dto.IngredientName = sp.Ingredient.IngredientName
	}
	if sp.Supplier != nil {
		dto.SupplierName = sp.Supplier.SupplierName
	}

	return dto
}

// ConvertSupplierPricesToDTO converts a slice of SupplierPrice to a slice of SupplierPriceDTO
func ConvertSupplierPricesToDTO(prices []SupplierPrice) []SupplierPriceDTO {
	dtos := make([]SupplierPriceDTO, len(prices))
	for i, price := range prices {
		dtos[i] = price.ToDTO()
	}
	return dtos
}