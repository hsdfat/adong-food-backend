package models

import "time"

// IngredientRequest represents a purchase request for ingredients from an order
type IngredientRequest struct {
	RequestID        string    `gorm:"column:request_id;primaryKey" json:"requestId"`
	OrderID          string    `gorm:"column:order_id;not null" json:"orderId"`
	KitchenID        string    `gorm:"column:kitchen_id;not null" json:"kitchenId"`
	RequestDate      time.Time `gorm:"column:request_date;type:date;not null" json:"requestDate"`
	RequiredDate     time.Time `gorm:"column:required_date;type:date" json:"requiredDate"`
	Status           string    `gorm:"column:status;not null;default:pending" json:"status"`
	TotalAmount      float64   `gorm:"column:total_amount;default:0" json:"totalAmount"`
	Notes            *string   `gorm:"column:notes" json:"notes,omitempty"`
	CreatedByUserID  *string   `gorm:"column:created_by_user_id" json:"createdByUserId,omitempty"`
	ApprovedByUserID *string   `gorm:"column:approved_by_user_id" json:"approvedByUserId,omitempty"`
	ApprovedDate     *time.Time `gorm:"column:approved_date" json:"approvedDate,omitempty"`
	CreatedDate      time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate     time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Order        *Order                      `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
	Kitchen      *Kitchen                    `gorm:"foreignKey:KitchenID;references:KitchenID" json:"kitchen,omitempty"`
	CreatedBy    *User                       `gorm:"foreignKey:CreatedByUserID;references:UserID" json:"createdBy,omitempty"`
	ApprovedBy   *User                       `gorm:"foreignKey:ApprovedByUserID;references:UserID" json:"approvedBy,omitempty"`
	RequestDetails []IngredientRequestDetail `gorm:"foreignKey:RequestID;references:RequestID" json:"requestDetails,omitempty"`
}

func (IngredientRequest) TableName() string {
	return "ingredient_requests"
}

// IngredientRequestDetail represents details of requested ingredients with selected supplier
type IngredientRequestDetail struct {
	RequestDetailID int       `gorm:"column:request_detail_id;primaryKey;autoIncrement" json:"requestDetailId"`
	RequestID       string    `gorm:"column:request_id;not null" json:"requestId"`
	IngredientID    string    `gorm:"column:ingredient_id;not null" json:"ingredientId"`
	Quantity        float64   `gorm:"column:quantity;not null" json:"quantity"`
	Unit            string    `gorm:"column:unit;not null" json:"unit"`
	SupplierID      *string   `gorm:"column:supplier_id" json:"supplierId,omitempty"`
	UnitPrice       *float64  `gorm:"column:unit_price" json:"unitPrice,omitempty"`
	TotalPrice      *float64  `gorm:"column:total_price" json:"totalPrice,omitempty"`
	Notes           *string   `gorm:"column:notes" json:"notes,omitempty"`
	CreatedDate     time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate    time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Request    *IngredientRequest `gorm:"foreignKey:RequestID;references:RequestID" json:"request,omitempty"`
	Ingredient *Ingredient        `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
	Supplier   *Supplier          `gorm:"foreignKey:SupplierID;references:SupplierID" json:"supplier,omitempty"`
}

func (IngredientRequestDetail) TableName() string {
	return "ingredient_request_details"
}
