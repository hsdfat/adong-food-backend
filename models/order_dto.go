package models

import "time"

// OrderDTO - Aggregated response for an order
type OrderDTO struct {
	OrderID         int                     `json:"orderId"`
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
