package models

import "time"

// Order - Orders (orders)
type Order struct {
	OrderID         int       `gorm:"primaryKey;autoIncrement;column:order_id" json:"orderId"`
	KitchenID       string    `gorm:"column:kitchen_id;not null" json:"kitchenId"`
	OrderDate       time.Time `gorm:"column:order_date;not null" json:"orderDate"`
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
