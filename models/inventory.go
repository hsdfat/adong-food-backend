package models

import "time"

// InventoryStock represents the current stock level of an ingredient in a kitchen
type InventoryStock struct {
	StockID       int       `gorm:"column:stock_id;primaryKey;autoIncrement" json:"stockId"`
	KitchenID     string    `gorm:"column:kitchen_id;not null" json:"kitchenId"`
	IngredientID  string    `gorm:"column:ingredient_id;not null" json:"ingredientId"`
	Quantity      float64   `gorm:"column:quantity;not null;default:0" json:"quantity"`
	Unit          string    `gorm:"column:unit;not null" json:"unit"`
	MinStockLevel *float64  `gorm:"column:min_stock_level" json:"minStockLevel,omitempty"`
	MaxStockLevel *float64  `gorm:"column:max_stock_level" json:"maxStockLevel,omitempty"`
	LastUpdated   time.Time `gorm:"column:last_updated;autoUpdateTime" json:"lastUpdated"`
	CreatedDate   time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate  time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Kitchen    *Kitchen    `gorm:"foreignKey:KitchenID;references:KitchenID" json:"kitchen,omitempty"`
	Ingredient *Ingredient `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
}

func (InventoryStock) TableName() string {
	return "inventory_stocks"
}

// InventoryImport represents a goods receipt/import transaction
type InventoryImport struct {
	ImportID         string     `gorm:"column:import_id;primaryKey" json:"importId"`
	KitchenID        string     `gorm:"column:kitchen_id;not null" json:"kitchenId"`
	ImportDate       time.Time  `gorm:"column:import_date;type:date;not null" json:"importDate"`
	OrderID          *string    `gorm:"column:order_id" json:"orderId,omitempty"`
	SupplierID       *string    `gorm:"column:supplier_id" json:"supplierId,omitempty"`
	TotalAmount      float64    `gorm:"column:total_amount;default:0" json:"totalAmount"`
	Status           string     `gorm:"column:status;not null;default:draft" json:"status"`
	Notes            *string    `gorm:"column:notes" json:"notes,omitempty"`
	ReceivedByUserID *string    `gorm:"column:received_by_user_id" json:"receivedByUserId,omitempty"`
	ApprovedByUserID *string    `gorm:"column:approved_by_user_id" json:"approvedByUserId,omitempty"`
	ApprovedDate     *time.Time `gorm:"column:approved_date" json:"approvedDate,omitempty"`
	CreatedByUserID  *string    `gorm:"column:created_by_user_id" json:"createdByUserId,omitempty"`
	CreatedDate      time.Time  `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate     time.Time  `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Kitchen      *Kitchen                   `gorm:"foreignKey:KitchenID;references:KitchenID" json:"kitchen,omitempty"`
	Order        *Order                     `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
	Supplier     *Supplier                  `gorm:"foreignKey:SupplierID;references:SupplierID" json:"supplier,omitempty"`
	ReceivedBy   *User                      `gorm:"foreignKey:ReceivedByUserID;references:UserID" json:"receivedBy,omitempty"`
	ApprovedBy   *User                      `gorm:"foreignKey:ApprovedByUserID;references:UserID" json:"approvedBy,omitempty"`
	CreatedBy    *User                      `gorm:"foreignKey:CreatedByUserID;references:UserID" json:"createdBy,omitempty"`
	ImportDetails []InventoryImportDetail    `gorm:"foreignKey:ImportID;references:ImportID" json:"importDetails,omitempty"`
}

func (InventoryImport) TableName() string {
	return "inventory_imports"
}

// InventoryImportDetail represents details of imported ingredients
type InventoryImportDetail struct {
	ImportDetailID int        `gorm:"column:import_detail_id;primaryKey;autoIncrement" json:"importDetailId"`
	ImportID       string     `gorm:"column:import_id;not null" json:"importId"`
	IngredientID   string     `gorm:"column:ingredient_id;not null" json:"ingredientId"`
	SupplierID     *string    `gorm:"column:supplier_id" json:"supplierId,omitempty"`
	Quantity       float64    `gorm:"column:quantity;not null" json:"quantity"`
	Unit           string     `gorm:"column:unit;not null" json:"unit"`
	UnitPrice      float64    `gorm:"column:unit_price;not null" json:"unitPrice"`
	TotalPrice     float64    `gorm:"column:total_price;not null" json:"totalPrice"`
	ExpiryDate     *time.Time `gorm:"column:expiry_date;type:date" json:"expiryDate,omitempty"`
	BatchNumber    *string    `gorm:"column:batch_number" json:"batchNumber,omitempty"`
	Notes          *string    `gorm:"column:notes" json:"notes,omitempty"`
	CreatedDate    time.Time  `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate   time.Time  `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Import     *InventoryImport `gorm:"foreignKey:ImportID;references:ImportID" json:"import,omitempty"`
	Ingredient *Ingredient      `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
	Supplier   *Supplier        `gorm:"foreignKey:SupplierID;references:SupplierID" json:"supplier,omitempty"`
}

func (InventoryImportDetail) TableName() string {
	return "inventory_import_details"
}

// InventoryExport represents a goods issue/export transaction
type InventoryExport struct {
	ExportID            string     `gorm:"column:export_id;primaryKey" json:"exportId"`
	KitchenID           string     `gorm:"column:kitchen_id;not null" json:"kitchenId"`
	ExportDate          time.Time  `gorm:"column:export_date;type:date;not null" json:"exportDate"`
	ExportType          string     `gorm:"column:export_type;not null" json:"exportType"`
	DestinationKitchenID *string    `gorm:"column:destination_kitchen_id" json:"destinationKitchenId,omitempty"`
	OrderID             *string    `gorm:"column:order_id" json:"orderId,omitempty"`
	TotalAmount         float64    `gorm:"column:total_amount;default:0" json:"totalAmount"`
	Status              string     `gorm:"column:status;not null;default:draft" json:"status"`
	Notes               *string    `gorm:"column:notes" json:"notes,omitempty"`
	IssuedByUserID      *string    `gorm:"column:issued_by_user_id" json:"issuedByUserId,omitempty"`
	ApprovedByUserID    *string    `gorm:"column:approved_by_user_id" json:"approvedByUserId,omitempty"`
	ApprovedDate        *time.Time `gorm:"column:approved_date" json:"approvedDate,omitempty"`
	CreatedByUserID     *string    `gorm:"column:created_by_user_id" json:"createdByUserId,omitempty"`
	CreatedDate         time.Time  `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate        time.Time  `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Kitchen            *Kitchen                  `gorm:"foreignKey:KitchenID;references:KitchenID" json:"kitchen,omitempty"`
	DestinationKitchen *Kitchen                  `gorm:"foreignKey:DestinationKitchenID;references:KitchenID" json:"destinationKitchen,omitempty"`
	Order              *Order                    `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
	IssuedBy           *User                     `gorm:"foreignKey:IssuedByUserID;references:UserID" json:"issuedBy,omitempty"`
	ApprovedBy         *User                     `gorm:"foreignKey:ApprovedByUserID;references:UserID" json:"approvedBy,omitempty"`
	CreatedBy          *User                     `gorm:"foreignKey:CreatedByUserID;references:UserID" json:"createdBy,omitempty"`
	ExportDetails      []InventoryExportDetail   `gorm:"foreignKey:ExportID;references:ExportID" json:"exportDetails,omitempty"`
}

func (InventoryExport) TableName() string {
	return "inventory_exports"
}

// InventoryExportDetail represents details of exported ingredients
type InventoryExportDetail struct {
	ExportDetailID int       `gorm:"column:export_detail_id;primaryKey;autoIncrement" json:"exportDetailId"`
	ExportID       string    `gorm:"column:export_id;not null" json:"exportId"`
	IngredientID   string    `gorm:"column:ingredient_id;not null" json:"ingredientId"`
	Quantity       float64   `gorm:"column:quantity;not null" json:"quantity"`
	Unit           string    `gorm:"column:unit;not null" json:"unit"`
	UnitCost       *float64  `gorm:"column:unit_cost" json:"unitCost,omitempty"`
	TotalCost      *float64  `gorm:"column:total_cost" json:"totalCost,omitempty"`
	BatchNumber    *string   `gorm:"column:batch_number" json:"batchNumber,omitempty"`
	Notes          *string   `gorm:"column:notes" json:"notes,omitempty"`
	CreatedDate    time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate   time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Export     *InventoryExport `gorm:"foreignKey:ExportID;references:ExportID" json:"export,omitempty"`
	Ingredient *Ingredient      `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
}

func (InventoryExportDetail) TableName() string {
	return "inventory_export_details"
}

// InventoryTransaction represents a log entry for all inventory movements
type InventoryTransaction struct {
	TransactionID   int       `gorm:"column:transaction_id;primaryKey;autoIncrement" json:"transactionId"`
	KitchenID       string    `gorm:"column:kitchen_id;not null" json:"kitchenId"`
	IngredientID    string    `gorm:"column:ingredient_id;not null" json:"ingredientId"`
	TransactionType string    `gorm:"column:transaction_type;not null" json:"transactionType"`
	TransactionDate time.Time `gorm:"column:transaction_date;autoCreateTime" json:"transactionDate"`
	Quantity        float64   `gorm:"column:quantity;not null" json:"quantity"`
	Unit            string    `gorm:"column:unit;not null" json:"unit"`
	QuantityBefore  float64   `gorm:"column:quantity_before;not null" json:"quantityBefore"`
	QuantityAfter   float64   `gorm:"column:quantity_after;not null" json:"quantityAfter"`
	ReferenceType   *string   `gorm:"column:reference_type" json:"referenceType,omitempty"`
	ReferenceID     *string   `gorm:"column:reference_id" json:"referenceId,omitempty"`
	Notes           *string   `gorm:"column:notes" json:"notes,omitempty"`
	CreatedByUserID *string   `gorm:"column:created_by_user_id" json:"createdByUserId,omitempty"`
	CreatedDate     time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`

	// Relationships
	Kitchen    *Kitchen    `gorm:"foreignKey:KitchenID;references:KitchenID" json:"kitchen,omitempty"`
	Ingredient *Ingredient `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
	CreatedBy  *User       `gorm:"foreignKey:CreatedByUserID;references:UserID" json:"createdBy,omitempty"`
}

func (InventoryTransaction) TableName() string {
	return "inventory_transactions"
}

// InventoryAdjustment represents inventory count adjustments
type InventoryAdjustment struct {
	AdjustmentID     string                        `gorm:"column:adjustment_id;primaryKey" json:"adjustmentId"`
	KitchenID        string                        `gorm:"column:kitchen_id;not null" json:"kitchenId"`
	AdjustmentDate   time.Time                     `gorm:"column:adjustment_date;type:date;not null" json:"adjustmentDate"`
	AdjustmentType   string                        `gorm:"column:adjustment_type;not null" json:"adjustmentType"`
	Reason           *string                       `gorm:"column:reason" json:"reason,omitempty"`
	Status           string                        `gorm:"column:status;not null;default:draft" json:"status"`
	TotalValue       *float64                      `gorm:"column:total_value" json:"totalValue,omitempty"`
	ApprovedByUserID *string                       `gorm:"column:approved_by_user_id" json:"approvedByUserId,omitempty"`
	ApprovedDate     *time.Time                    `gorm:"column:approved_date" json:"approvedDate,omitempty"`
	CreatedByUserID  *string                       `gorm:"column:created_by_user_id" json:"createdByUserId,omitempty"`
	CreatedDate      time.Time                     `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate     time.Time                     `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Kitchen           *Kitchen                      `gorm:"foreignKey:KitchenID;references:KitchenID" json:"kitchen,omitempty"`
	ApprovedBy        *User                         `gorm:"foreignKey:ApprovedByUserID;references:UserID" json:"approvedBy,omitempty"`
	CreatedBy         *User                         `gorm:"foreignKey:CreatedByUserID;references:UserID" json:"createdBy,omitempty"`
	AdjustmentDetails []InventoryAdjustmentDetail   `gorm:"foreignKey:AdjustmentID;references:AdjustmentID" json:"adjustmentDetails,omitempty"`
}

func (InventoryAdjustment) TableName() string {
	return "inventory_adjustments"
}

// InventoryAdjustmentDetail represents details of inventory adjustments
type InventoryAdjustmentDetail struct {
	AdjustmentDetailID int       `gorm:"column:adjustment_detail_id;primaryKey;autoIncrement" json:"adjustmentDetailId"`
	AdjustmentID       string    `gorm:"column:adjustment_id;not null" json:"adjustmentId"`
	IngredientID       string    `gorm:"column:ingredient_id;not null" json:"ingredientId"`
	QuantityBefore     float64   `gorm:"column:quantity_before;not null" json:"quantityBefore"`
	QuantityAfter      float64   `gorm:"column:quantity_after;not null" json:"quantityAfter"`
	QuantityDifference float64   `gorm:"column:quantity_difference;not null" json:"quantityDifference"`
	Unit               string    `gorm:"column:unit;not null" json:"unit"`
	UnitCost           *float64  `gorm:"column:unit_cost" json:"unitCost,omitempty"`
	TotalValue         *float64  `gorm:"column:total_value" json:"totalValue,omitempty"`
	Reason             *string   `gorm:"column:reason" json:"reason,omitempty"`
	CreatedDate        time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate       time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Adjustment *InventoryAdjustment `gorm:"foreignKey:AdjustmentID;references:AdjustmentID" json:"adjustment,omitempty"`
	Ingredient *Ingredient          `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
}

func (InventoryAdjustmentDetail) TableName() string {
	return "inventory_adjustment_details"
}
