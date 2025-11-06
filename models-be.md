# ProjectDump Analysis

**Generated on:** 2025-11-06 04:38:03
**Project Path:** models/

## Project Summary

- **Primary Language:** Go
- **Total Files:** 13
- **Processed Files:** 13
- **Project Size:** 22.53 KB

## Detected Technologies

### Go (100.0% confidence)
*Go programming language*

**Related files:**
- common.go
- dish.go
- ingredient.go
- kitchen.go
- order.go
- ... and 8 more files

### CSS (100.0% confidence)
*Cascading Style Sheets*

**Related files:**
- dish.go
- ingredient.go
- kitchen.go
- order.go
- order_dto.go
- ... and 7 more files

### Python (75.0% confidence)
*Python programming language*

**Related files:**
- dish.go
- ingredient.go
- kitchen.go
- order.go
- order_dto.go
- ... and 7 more files

### Java (65.0% confidence)
*Java programming language*

**Related files:**
- common.go
- dish.go
- ingredient.go
- kitchen.go
- order.go
- ... and 8 more files

### TypeScript (60.0% confidence)
*TypeScript - JavaScript with static typing*

**Related files:**
- dish.go
- ingredient.go
- kitchen.go
- order.go
- order_dto.go
- ... and 7 more files

### JavaScript (55.0% confidence)
*JavaScript runtime and ecosystem*

**Related files:**
- dish.go
- ingredient.go
- kitchen.go
- order.go
- order_dto.go
- ... and 6 more files

### Docker (20.0% confidence)
*Docker containerization platform*

**Related files:**
- pagination.go
- recipe_standard_dto.go
- supplier_price.go
- supplier_price_dto.go

### Ruby (10.0% confidence)
*Ruby programming language*

**Related files:**
- order.go
- pagination.go

## Directory Structure

```
├── common.go
├── dish.go
├── ingredient.go
├── kitchen.go
├── order.go
├── order_dto.go
├── pagination.go
├── recipe_standard.go
├── recipe_standard_dto.go
├── supplier.go
├── supplier_price.go
├── supplier_price_dto.go
└── user.go
```

## Source Code

#### common.go
*Language: Go | Size: 17 bytes*

```go
package models


```

#### dish.go
*Language: Go | Size: 816 bytes*

```go
package models

import "time"

// Dish - Master data for dishes/menu items (dm_monan)
type Dish struct {
    DishID        string    `gorm:"primaryKey;column:dish_id" json:"dishId"`
    DishName      string    `gorm:"column:dish_name;not null" json:"dishName"`
    CookingMethod string    `gorm:"column:cooking_method" json:"cookingMethod"`
    Group         string    `gorm:"column:category" json:"group"`
    Description   string    `gorm:"column:description;type:text" json:"description"`
	Active        *bool     `gorm:"column:active;default:true" json:"active"`
    CreatedDate   time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
    ModifiedDate  time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`
}

func (Dish) TableName() string {
    return "master_dishes"
}
```

#### ingredient.go
*Language: Go | Size: 786 bytes*

```go
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
```

#### kitchen.go
*Language: Go | Size: 750 bytes*

```go
package models

import "time"

// Kitchen - Master data for kitchen/location information (dm_bep)
type Kitchen struct {
    KitchenID    string    `gorm:"primaryKey;column:kitchen_id" json:"kitchenId"`
    KitchenName  string    `gorm:"column:kitchen_name;not null" json:"kitchenName"`
    Address      string    `gorm:"column:address;type:text" json:"address"`
    Phone        string    `gorm:"column:phone" json:"phone"`
	Active       *bool     `gorm:"column:active;default:true" json:"active"`
    CreatedDate  time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
    ModifiedDate time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`
}

func (Kitchen) TableName() string {
    return "master_kitchens"
}
```

#### order.go
*Language: Go | Size: 7218 bytes*

```go
// models/order.go
package models

import "time"

// Order - Orders (orders)
type Order struct {
	OrderID         int       `gorm:"primaryKey;autoIncrement;column:order_id" json:"orderId"`
	KitchenID       string    `gorm:"column:kitchen_id;not null" json:"kitchenId"`
	OrderDate       string    `gorm:"column:order_date;not null" json:"orderDate"`
	Note            string    `gorm:"column:note;type:text" json:"note"`
	Status          string    `gorm:"column:status;default:Pending;not null" json:"status"`
	CreatedByUserID string    `gorm:"column:created_by_user_id" json:"createdByUserId"`
	CreatedDate     time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate    time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Kitchen            *Kitchen                 `gorm:"foreignKey:KitchenID;references:KitchenID" json:"kitchen,omitempty"`
	CreatedBy          *User                    `gorm:"foreignKey:CreatedByUserID;references:UserID" json:"createdBy,omitempty"`
	Details            []OrderDetail            `gorm:"foreignKey:OrderID;references:OrderID" json:"details,omitempty"`
	SupplementaryFoods []OrderSupplementaryFood `gorm:"foreignKey:OrderID;references:OrderID" json:"supplementaryFoods,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

// OrderDetail - Order details (order_details)
type OrderDetail struct {
	OrderDetailID int       `gorm:"primaryKey;autoIncrement;column:order_detail_id" json:"orderDetailId"`
	OrderID       int       `gorm:"column:order_id;not null" json:"orderId"`
	DishID        string    `gorm:"column:dish_id;not null" json:"dishId"`
	Portions      int       `gorm:"column:portions;not null" json:"portions"`
	Note          string    `gorm:"column:note;type:text" json:"note"`
	CreatedDate   time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate  time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Order       *Order            `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
	Dish        *Dish             `gorm:"foreignKey:DishID;references:DishID" json:"dish,omitempty"`
	Ingredients []OrderIngredient `gorm:"foreignKey:OrderDetailID;references:OrderDetailID" json:"ingredients,omitempty"`
}

func (OrderDetail) TableName() string {
	return "order_details"
}

// OrderIngredient - Ingredients calculated for an order detail (order_ingredients)
type OrderIngredient struct {
	OrderIngredientID  int       `gorm:"primaryKey;autoIncrement;column:order_ingredient_id" json:"orderIngredientId"`
	OrderDetailID      int       `gorm:"column:order_detail_id;not null" json:"orderDetailId"`
	IngredientID       string    `gorm:"column:ingredient_id;not null" json:"ingredientId"`
	Quantity           float64   `gorm:"column:quantity;type:numeric(15,4);not null" json:"quantity"`
	Unit               string    `gorm:"column:unit;not null" json:"unit"`
	StandardPerPortion float64   `gorm:"column:standard_per_portion;type:numeric(10,4)" json:"standardPerPortion"`
	CreatedDate        time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate       time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	OrderDetail *OrderDetail `gorm:"foreignKey:OrderDetailID;references:OrderDetailID" json:"orderDetail,omitempty"`
	Ingredient  *Ingredient  `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
}

func (OrderIngredient) TableName() string {
	return "order_ingredients"
}

// OrderSupplementaryFood - Extra items for an order (order_supplementary_foods)
type OrderSupplementaryFood struct {
	SupplementaryID    int       `gorm:"primaryKey;autoIncrement;column:supplementary_id" json:"supplementaryId"`
	OrderID            int       `gorm:"column:order_id;not null" json:"orderId"`
	IngredientID       string    `gorm:"column:ingredient_id;not null" json:"ingredientId"`
	Quantity           float64   `gorm:"column:quantity;type:numeric(15,4);not null" json:"quantity"`
	Unit               string    `gorm:"column:unit;not null" json:"unit"`
	StandardPerPortion float64   `gorm:"column:standard_per_portion;type:numeric(10,4)" json:"standardPerPortion"`
	Portions           int       `gorm:"column:portions" json:"portions"`
	Note               string    `gorm:"column:note;type:text" json:"note"`
	CreatedDate        time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate       time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Order      *Order      `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
	Ingredient *Ingredient `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
}

func (OrderSupplementaryFood) TableName() string {
	return "order_supplementary_foods"
}

// SupplierRequest - Requests sent to suppliers (supplier_requests)
type SupplierRequest struct {
	RequestID    int       `gorm:"primaryKey;autoIncrement;column:request_id" json:"requestId"`
	OrderID      int       `gorm:"column:order_id;not null" json:"orderId"`
	SupplierID   string    `gorm:"column:supplier_id;not null" json:"supplierId"`
	Status       string    `gorm:"column:status;default:Pending;not null" json:"status"`
	CreatedDate  time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Order    *Order                  `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
	Supplier *Supplier               `gorm:"foreignKey:SupplierID;references:SupplierID" json:"supplier,omitempty"`
	Details  []SupplierRequestDetail `gorm:"foreignKey:RequestID;references:RequestID" json:"details,omitempty"`
}

func (SupplierRequest) TableName() string {
	return "supplier_requests"
}

// SupplierRequestDetail - Line items for a supplier request (supplier_request_details)
type SupplierRequestDetail struct {
	RequestDetailID int       `gorm:"primaryKey;autoIncrement;column:request_detail_id" json:"requestDetailId"`
	RequestID       int       `gorm:"column:request_id;not null" json:"requestId"`
	IngredientID    string    `gorm:"column:ingredient_id;not null" json:"ingredientId"`
	Quantity        float64   `gorm:"column:quantity;type:numeric(15,4);not null" json:"quantity"`
	Unit            string    `gorm:"column:unit;not null" json:"unit"`
	UnitPrice       float64   `gorm:"column:unit_price;type:numeric(15,2);not null" json:"unitPrice"`
	TotalPrice      float64   `gorm:"column:total_price;type:numeric(15,2)" json:"totalPrice"`
	CreatedDate     time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate    time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Request    *SupplierRequest `gorm:"foreignKey:RequestID;references:RequestID" json:"request,omitempty"`
	Ingredient *Ingredient      `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
}

func (SupplierRequestDetail) TableName() string {
	return "supplier_request_details"
}
```

#### order_dto.go
*Language: Go | Size: 2131 bytes*

```go
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
```

#### pagination.go
*Language: Go | Size: 2231 bytes*

```go
package models

// PaginationParams contains pagination parameters from query string
type PaginationParams struct {
	Page     int    `form:"page" binding:"omitempty,min=0"`
	PageSize int    `form:"per_page" binding:"omitempty,min=0,max=100"` // Changed to per_page to match common conventions
	Search   string `form:"search"`
	SortBy   string `form:"sort_by"`
	SortDir  string `form:"sort_dir" binding:"omitempty,oneof=asc desc"`
}

// PaginationMeta contains pagination metadata matching frontend ResourceCollection interface
type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	From        int `json:"from"`
	To          int `json:"to"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
}

// ResourceCollection is the response wrapper matching frontend interface
type ResourceCollection struct {
	Data interface{}     `json:"data"`
	Meta *PaginationMeta `json:"meta"`
}

// GetPaginationParams extracts and validates pagination parameters with defaults
func GetPaginationParams(page, pageSize int, search, sortBy, sortDir string) PaginationParams {
	// Set defaults
	// if page < 1 {
	// 	page = 1
	// }
	// if pageSize < 1 {
	// 	pageSize = 10
	// }
	// if pageSize > 100 {
	// 	pageSize = 100
	// }
	if sortDir == "" {
		sortDir = "asc"
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
		SortBy:   sortBy,
		SortDir:  sortDir,
	}
}

// CalculatePaginationMeta calculates pagination metadata
func CalculatePaginationMeta(page, perPage int, total int64) *PaginationMeta {
	if perPage < 1 || page < 1 {
		return &PaginationMeta{
			CurrentPage: page,
			LastPage:    1,
			From:        0,
			To:          int(total),
			PerPage:     perPage,
			Total:       int(total),
		}
	}
	totalInt := int(total)
	lastPage := (totalInt + perPage - 1) / perPage
	if lastPage < 1 {
		lastPage = 1
	}

	// Calculate from and to
	from := 0
	to := 0
	if totalInt > 0 {
		from = (page-1)*perPage + 1
		to = from + perPage - 1
		if to > totalInt {
			to = totalInt
		}
	}

	return &PaginationMeta{
		CurrentPage: page,
		LastPage:    lastPage,
		From:        from,
		To:          to,
		PerPage:     perPage,
		Total:       totalInt,
	}
}
```

#### recipe_standard.go
*Language: Go | Size: 1319 bytes*

```go
package models

import "time"

// RecipeStandard - Bill of materials for dishes (dish_recipe_standards)
type RecipeStandard struct {
	StandardID   int       `gorm:"primaryKey;autoIncrement;column:recipe_id" json:"standardId"`
	DishID       string    `gorm:"column:dish_id" json:"dishId"`
	IngredientID string    `gorm:"column:ingredient_id" json:"ingredientId"`
	Unit         string    `gorm:"column:unit" json:"unit"`
	StandardPer1 float64   `gorm:"column:quantity_per_serving;type:decimal(10,4)" json:"standardPer1"`
	Note         string    `gorm:"column:notes;type:text" json:"note"`
	Amount       float64   `gorm:"column:cost;type:decimal(15,2)" json:"amount"`
	UpdatedByID  string    `gorm:"column:updated_by_user_id" json:"updatedById"`
	CreatedDate  time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Dish       *Dish       `gorm:"foreignKey:DishID;references:DishID" json:"dish,omitempty"`
	Ingredient *Ingredient `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
	UpdatedBy  *User       `gorm:"foreignKey:UpdatedByID;references:UserID" json:"updatedBy,omitempty"`
}

func (RecipeStandard) TableName() string {
	return "dish_recipe_standards"
}
```

#### recipe_standard_dto.go
*Language: Go | Size: 1819 bytes*

```go
package models

import "time"

// RecipeStandardDTO - Data Transfer Object for Recipe Standard with related names
type RecipeStandardDTO struct {
	StandardID     int       `json:"standardId"`
	DishID         string    `json:"dishId"`
	DishName       string    `json:"dishName"`           // Added: Dish name
	IngredientID   string    `json:"ingredientId"`
	IngredientName string    `json:"ingredientName"`     // Added: Ingredient name
	Unit           string    `json:"unit"`
	StandardPer1   float64   `json:"standardPer1"`
	Note           string    `json:"note"`
	Amount         float64   `json:"amount"`
	UpdatedByID    string    `json:"updatedById"`
	UpdatedByName  string    `json:"updatedByName"`      // Added: User name (optional)
	CreatedDate    time.Time `json:"createdDate"`
	ModifiedDate   time.Time `json:"modifiedDate"`
}

// ToDTO converts RecipeStandard model to DTO
func (r *RecipeStandard) ToDTO() RecipeStandardDTO {
	dto := RecipeStandardDTO{
		StandardID:   r.StandardID,
		DishID:       r.DishID,
		IngredientID: r.IngredientID,
		Unit:         r.Unit,
		StandardPer1: r.StandardPer1,
		Note:         r.Note,
		Amount:       r.Amount,
		UpdatedByID:  r.UpdatedByID,
		CreatedDate:  r.CreatedDate,
		ModifiedDate: r.ModifiedDate,
	}

	// Populate names from relationships if available
	if r.Dish != nil {
		dto.DishName = r.Dish.DishName
	}
	if r.Ingredient != nil {
		dto.IngredientName = r.Ingredient.IngredientName
	}
	if r.UpdatedBy != nil {
		dto.UpdatedByName = r.UpdatedBy.FullName
	}

	return dto
}

// ConvertToDTO converts a slice of RecipeStandard to a slice of RecipeStandardDTO
func ConvertRecipeStandardsToDTO(recipes []RecipeStandard) []RecipeStandardDTO {
	dtos := make([]RecipeStandardDTO, len(recipes))
	for i, recipe := range recipes {
		dtos[i] = recipe.ToDTO()
	}
	return dtos
}
```

#### supplier.go
*Language: Go | Size: 877 bytes*

```go
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
    CreatedDate  time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
    ModifiedDate time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`
}

func (Supplier) TableName() string {
    return "master_suppliers"
}
```

#### supplier_price.go
*Language: Go | Size: 1800 bytes*

```go
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
```

#### supplier_price_dto.go
*Language: Go | Size: 2175 bytes*

```go
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
```

#### user.go
*Language: Go | Size: 1136 bytes*

```go
package models

import (
	"time"
)






// User - Master data for user accounts (dm_nguoidung)
type User struct {
    UserID       string    `gorm:"primaryKey;column:user_id" json:"userId"`
    UserName     string    `gorm:"column:user_name;not null;unique" json:"userName"`
	Password     string    `gorm:"column:password;not null" json:"password,omitempty"`
    FullName     string    `gorm:"column:full_name;not null" json:"fullName"`
    Role         string    `gorm:"column:role" json:"role"`
    KitchenID    string    `gorm:"column:kitchen_id" json:"kitchenId"`
	Email        string    `gorm:"column:email" json:"email"`
    Phone        string    `gorm:"column:phone" json:"phone"`
	Active       *bool     `gorm:"column:active;default:true" json:"active"`
    CreatedDate  time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
    ModifiedDate time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`
	
	// Relationships
	Kitchen      *Kitchen  `gorm:"foreignKey:KitchenID;references:KitchenID" json:"kitchen,omitempty"`
}

func (User) TableName() string {
    return "master_users"
}






```

