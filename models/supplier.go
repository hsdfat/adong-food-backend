package models

import "time"

// Supplier - Master data for suppliers (dm_ncc)
type Supplier struct {
	SupplierID   string    `gorm:"primaryKey;column:supplier_id" json:"supplierId"`
	SupplierName string    `gorm:"column:supplier_name;not null" json:"supplierName"`
	ZaloLink     string    `gorm:"column:zalo_link;type:text" json:"zaloLink"`
	Address      string    `gorm:"column:address;type:text" json:"address"`
	Phone        string    `gorm:"column:phone" json:"phone"`
	Email        string    `gorm:"column:email" json:"email"`
	Active       *bool     `gorm:"column:active;default:true" json:"active"`
	LegacyID     *string   `gorm:"column:legacy_id" json:"legacyId,omitempty"`
	CreatedDate  time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`
}

func (Supplier) TableName() string {
	return "master_suppliers"
}

// BestSupplierRequest - Request to find best suppliers for ingredients in an order
type BestSupplierRequest struct {
	OrderID       string   `json:"orderId" binding:"required"`
	KitchenID     string   `json:"kitchenId" binding:"required"`
	IngredientIDs []string `json:"ingredientIds" binding:"required,min=1"`
}

// SupplierInfo - Information about a selected supplier for an ingredient
type SupplierInfo struct {
	SupplierID   string  `json:"supplierId"`
	SupplierName string  `json:"supplierName"`
	Phone        string  `json:"phone"`
	Email        string  `json:"email"`
	Address      string  `json:"address"`
	UnitPrice    float64 `json:"unitPrice"`
	Unit         string  `json:"unit"`
	ProductName  string  `json:"productName"`
	ProductID    int     `json:"productId"`
}

// IngredientSupplierInfo - Supplier information for a specific ingredient
type IngredientSupplierInfo struct {
	IngredientID     string        `json:"ingredientId"`
	IngredientName   string        `json:"ingredientName"`
	IngredientType   string        `json:"ingredientType"`
	MaterialGroup    string        `json:"materialGroup"`
	SelectedSupplier *SupplierInfo `json:"selectedSupplier"`
	SelectionReason  string        `json:"selectionReason"`
}

// BestSupplierResponse - Response containing best suppliers for all ingredients
type BestSupplierResponse struct {
	OrderID   string                   `json:"orderId"`
	KitchenID string                   `json:"kitchenId"`
	Suppliers []IngredientSupplierInfo `json:"suppliers"`
}
