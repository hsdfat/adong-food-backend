// File: models/menu_card_dto.go
package models

// MenuCardCreateRequest - Request để tạo phiếu thực đơn
type MenuCardCreateRequest struct {
	MenuCardName string                      `json:"menuCardName" binding:"required"`
	CreatedDate  *string                     `json:"createdDate"`
	KitchenID    string                      `json:"kitchenId" binding:"required"`
	Note         string                      `json:"note"`
	Details      []MenuCardDetailRequest     `json:"details" binding:"required,min=1"`
}

// MenuCardDetailRequest - Request cho món ăn trong phiếu
type MenuCardDetailRequest struct {
	DishID      string                              `json:"dishId" binding:"required"`
	Servings    int                                 `json:"servings" binding:"required,min=1"`
	Note        string                              `json:"note"`
	Ingredients []MenuCardDetailIngredientRequest   `json:"ingredients"`
}

// MenuCardDetailIngredientRequest - Request cho nguyên liệu tùy chỉnh
type MenuCardDetailIngredientRequest struct {
	IngredientID string  `json:"ingredientId" binding:"required"`
	Standard     float64 `json:"standard" binding:"required,gt=0"`
	Unit         string  `json:"unit"`
	Note         string  `json:"note"`
}

// MenuCardDTO - Response với đầy đủ thông tin
type MenuCardDTO struct {
	MenuCardID   string              `json:"menuCardId"`
	MenuCardName string              `json:"menuCardName"`
	CreatedDate  *string             `json:"createdDate"`
	KitchenID    string              `json:"kitchenId"`
	KitchenName  string              `json:"kitchenName"`
	CreatedByID  string              `json:"createdById"`
	CreatedByName string             `json:"createdByName"`
	Status       string              `json:"status"`
	Note         string              `json:"note"`
	CreatedAt    string              `json:"createdAt"`
	ModifiedDate string              `json:"modifiedDate"`
	Details      []MenuCardDetailDTO `json:"details"`
}

// MenuCardDetailDTO - Response cho món ăn
type MenuCardDetailDTO struct {
	DetailID    string                          `json:"detailId"`
	DishID      string                          `json:"dishId"`
	DishName    string                          `json:"dishName"`
	Servings    int                             `json:"servings"`
	Note        string                          `json:"note"`
	Ingredients []MenuCardDetailIngredientDTO   `json:"ingredients"`
}

// MenuCardDetailIngredientDTO - Response cho nguyên liệu
type MenuCardDetailIngredientDTO struct {
	ID             string  `json:"id"`
	IngredientID   string  `json:"ingredientId"`
	IngredientName string  `json:"ingredientName"`
	Standard       float64 `json:"standard"`
	Unit           string  `json:"unit"`
	Note           string  `json:"note"`
}