package models

import "time"

// OrderDTO - Aggregated response for an order
type OrderDTO struct {
	OrderID         string                  `json:"orderId"`
	KitchenID       string                  `json:"kitchenId"`
	KitchenName     string                  `json:"kitchenName"`
	OrderDate       string                  `json:"orderDate"`
	Note            string                  `json:"note"`
	Status          string                  `json:"status"`
	CreatedByUserID string                  `json:"createdByUserId"`
	CreatedByName   string                  `json:"createdByName"`
	CreatedDate     time.Time               `json:"createdDate"`
	ModifiedDate    time.Time               `json:"modifiedDate"`
	Details         []OrderDetailDTO        `json:"details"`
	Supplementaries []OrderSupplementaryDTO `json:"supplementaries"`
}

// OrderDetailDTO - Detail lines with dish name and ingredients
type OrderDetailDTO struct {
	OrderDetailID int                  `json:"orderDetailId"`
	DishID        string               `json:"dishId"`
	DishName      string               `json:"dishName"`
	Portions      int                  `json:"portions"`
	Note          string               `json:"note"`
	Ingredients   []OrderIngredientDTO `json:"ingredients"`
}

// OrderIngredientDTO - Ingredient usage per detail
type OrderIngredientDTO struct {
	OrderIngredientID  int     `json:"orderIngredientId"`
	IngredientID       string  `json:"ingredientId"`
	IngredientName     string  `json:"ingredientName"`
	Quantity           float64 `json:"quantity"`
	Unit               string  `json:"unit"`
	StandardPerPortion float64 `json:"standardPerPortion"`
}

// OrderSupplementaryDTO - Supplementary items for an order
type OrderSupplementaryDTO struct {
	SupplementaryID    int     `json:"supplementaryId"`
	IngredientID       string  `json:"ingredientId"`
	IngredientName     string  `json:"ingredientName"`
	Quantity           float64 `json:"quantity"`
	Unit               string  `json:"unit"`
	StandardPerPortion float64 `json:"standardPerPortion"`
	Portions           int     `json:"portions"`
	Note               string  `json:"note"`
}

// OrderListItem represents a simplified order item for selection lists
type OrderListItem struct {
	OrderID    string `json:"orderId"`
	OrderDate  string `json:"orderDate"`
	KitchenID  string `json:"kitchenId"`
	Status     string `json:"status"`
	Note       string `json:"note"`
}

// OrderIngredientWithSupplier represents ingredient details with supplier for inventory operations
type OrderIngredientWithSupplier struct {
	OrderID            string  `json:"orderId"`
	IngredientID       string  `json:"ingredientId"`
	IngredientName     string  `json:"ingredientName"`
	Quantity           float64 `json:"quantity"`
	Unit               string  `json:"unit"`
	SupplierID         string  `json:"supplierId,omitempty"`
	SupplierName       string  `json:"supplierName,omitempty"`
	UnitPrice          float64 `json:"unitPrice,omitempty"`
	TotalCost          float64 `json:"totalCost,omitempty"`
}

// GetOrderSuppliersResponse represents the response for getting order suppliers
type GetOrderSuppliersResponse struct {
	OrderID     string                            `json:"orderId"`
	OrderDate   string                            `json:"orderDate"`
	Status      string                            `json:"status"`
	Suppliers   []OrderIngredientWithSupplier     `json:"suppliers"`
}

// SupplierWithOrderFlag represents a supplier with flag indicating if it's used in an order
type SupplierWithOrderFlag struct {
	SupplierID      string `json:"supplierId"`
	SupplierName    string `json:"supplierName"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Address         string `json:"address"`
	Active          *bool  `json:"active"`
	IsUsedInOrder   bool   `json:"isUsedInOrder"`   // Flag indicating if supplier is used in the order
	IngredientCount int    `json:"ingredientCount"` // Number of ingredients from this supplier in the order
}

// GetSuppliersForOrderResponse represents all suppliers with order usage highlighted
type GetSuppliersForOrderResponse struct {
	OrderID   string                    `json:"orderId"`
	Suppliers []SupplierWithOrderFlag   `json:"suppliers"`
}
