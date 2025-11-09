// models/order.go
package models

import "time"

// Order - Orders (orders)
type Order struct {
	OrderID         string    `gorm:"primaryKey;column:order_id;type:varchar(50)" json:"orderId"`
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
	OrderID       string    `gorm:"column:order_id;type:varchar(50);not null" json:"orderId"`
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
	OrderID            string    `gorm:"column:order_id;type:varchar(50);not null" json:"orderId"`
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

// OrderIngredientSupplier - Selected supplier for each ingredient in an order (order_ingredient_suppliers)
type OrderIngredientSupplier struct {
	OrderIngredientSupplierID int       `gorm:"primaryKey;autoIncrement;column:order_ingredient_supplier_id" json:"orderIngredientSupplierId"`
	OrderID                  string    `gorm:"column:order_id;type:varchar(50);not null" json:"orderId"`
	IngredientID             string    `gorm:"column:ingredient_id;not null" json:"ingredientId"`
	SelectedSupplierID       string    `gorm:"column:selected_supplier_id;not null" json:"selectedSupplierId"`
	SelectedProductID        int       `gorm:"column:selected_product_id;not null" json:"selectedProductId"`
	Quantity                 float64   `gorm:"column:quantity;type:numeric(15,4);not null" json:"quantity"`
	Unit                     string    `gorm:"column:unit;not null" json:"unit"`
	UnitPrice                float64   `gorm:"column:unit_price;type:numeric(15,2);not null" json:"unitPrice"`
	TotalCost                float64   `gorm:"column:total_cost;type:numeric(15,2);not null" json:"totalCost"`
	SelectionDate            time.Time `gorm:"column:selection_date;default:CURRENT_TIMESTAMP" json:"selectionDate"`
	SelectedByUserID         string    `gorm:"column:selected_by_user_id" json:"selectedByUserId"`
	Notes                    string    `gorm:"column:notes;type:text" json:"notes"`
	CreatedDate              time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate             time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Order            *Order        `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
	Ingredient       *Ingredient   `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
	SelectedSupplier *Supplier     `gorm:"foreignKey:SelectedSupplierID;references:SupplierID" json:"selectedSupplier,omitempty"`
	SelectedProduct  *SupplierPrice `gorm:"foreignKey:SelectedProductID;references:ProductID" json:"selectedProduct,omitempty"`
	SelectedBy       *User         `gorm:"foreignKey:SelectedByUserID;references:UserID" json:"selectedBy,omitempty"`
}

func (OrderIngredientSupplier) TableName() string {
	return "order_ingredient_suppliers"
}
