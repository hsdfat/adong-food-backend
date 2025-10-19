package models

import "time"

// IngredientByOrder - Ingredients by order (nguyenlieutheodon)
type IngredientByOrder struct {
	ID            string    `gorm:"primaryKey;column:id" json:"id"`
	OrderDetailID string    `gorm:"column:chitietlendonid" json:"orderDetailId"`
	IngredientID  string    `gorm:"column:nguyenlieuid" json:"ingredientId"`
	Standard      float64   `gorm:"column:dinhmuc;type:decimal(10,4)" json:"standard"`
	CreatedDate   time.Time `gorm:"column:createddate;autoCreateTime" json:"createdDate"`

	// Relationships
	OrderDetail *OrderDetail `gorm:"foreignKey:OrderDetailID;references:OrderDetailID" json:"orderDetail,omitempty"`
	Ingredient  *Ingredient  `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
}

func (IngredientByOrder) TableName() string {
	return "nguyenlieutheodon"
}

// IngredientRequest - Ingredient purchase requests (yeucaunguyenlieu)
type IngredientRequest struct {
	RequestID    string    `gorm:"primaryKey;column:yeucaunlid" json:"requestId"`
	OrderID      string    `gorm:"column:phieulendonid" json:"orderId"`
	IngredientID string    `gorm:"column:nguyenlieuid" json:"ingredientId"`
	RequestQty   float64   `gorm:"column:soluongyeucau;type:decimal(10,4)" json:"requestQty"`
	Unit         string    `gorm:"column:donvitinh" json:"unit"`
	ProductID    *int      `gorm:"column:sanphamid" json:"productId"`
	SupplierID   string    `gorm:"column:nhacungcapid" json:"supplierId"`
	UnitPrice    float64   `gorm:"column:dongia;type:decimal(15,2)" json:"unitPrice"`
	Amount       float64   `gorm:"column:thanhtien;type:decimal(15,2)" json:"amount"`
	Status       string    `gorm:"column:trangthai" json:"status"`
	CreatedDate  time.Time `gorm:"column:createddate;autoCreateTime" json:"createdDate"`

	// Relationships
	Order      *Order         `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
	Ingredient *Ingredient    `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
	Supplier   *Supplier      `gorm:"foreignKey:SupplierID;references:SupplierID" json:"supplier,omitempty"`
	Product    *SupplierPrice `gorm:"foreignKey:ProductID;references:ProductID" json:"product,omitempty"`
}

func (IngredientRequest) TableName() string {
	return "yeucaunguyenlieu"
}

// ReceivingDoc - Goods receiving documents (phieunhanhang)
type ReceivingDoc struct {
	ReceivingDocID string     `gorm:"primaryKey;column:phieunhanhangid" json:"receivingDocId"`
	OrderID        string     `gorm:"column:phieulendonid" json:"orderId"`
	SupplierID     string     `gorm:"column:nhacungcapid" json:"supplierId"`
	ReceivingDate  *time.Time `gorm:"column:ngaynhanhang" json:"receivingDate"`
	KitchenID      string     `gorm:"column:bepid" json:"kitchenId"`
	ReceivedByID   string     `gorm:"column:nguoinhanid" json:"receivedById"`
	Status         string     `gorm:"column:trangthai" json:"status"`
	Note           string     `gorm:"column:ghichu;type:text" json:"note"`
	CreatedDate    time.Time  `gorm:"column:createddate;autoCreateTime" json:"createdDate"`

	// Relationships
	Order            *Order            `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
	Supplier         *Supplier         `gorm:"foreignKey:SupplierID;references:SupplierID" json:"supplier,omitempty"`
	Kitchen          *Kitchen          `gorm:"foreignKey:KitchenID;references:KitchenID" json:"kitchen,omitempty"`
	ReceivedBy       *User             `gorm:"foreignKey:ReceivedByID;references:UserID" json:"receivedBy,omitempty"`
	ReceivingDetails []ReceivingDetail `gorm:"foreignKey:ReceivingDocID;references:ReceivingDocID" json:"receivingDetails,omitempty"`
}

func (ReceivingDoc) TableName() string {
	return "phieunhanhang"
}

// ReceivingDetail - Receiving line items (chitietnhanhang)
type ReceivingDetail struct {
	ReceivingDetailID int       `gorm:"primaryKey;autoIncrement;column:chitietnhanhangid" json:"receivingDetailId"`
	ReceivingDocID    string    `gorm:"column:phieunhanhangid" json:"receivingDocId"`
	IngredientID      string    `gorm:"column:nguyenlieuid" json:"ingredientId"`
	SupplierID        string    `gorm:"column:nhacungcapid" json:"supplierId"`
	RequestQty        float64   `gorm:"column:soluongyeucau;type:decimal(10,4)" json:"requestQty"`
	ActualQty         float64   `gorm:"column:soluongthucnhan;type:decimal(10,4)" json:"actualQty"`
	Unit              string    `gorm:"column:donvitinh" json:"unit"`
	UnitPrice         float64   `gorm:"column:dongia;type:decimal(15,2)" json:"unitPrice"`
	Amount            float64   `gorm:"column:thanhtien;type:decimal(15,2)" json:"amount"`
	ReceivingStatus   string    `gorm:"column:trangthainhan" json:"receivingStatus"`
	Note              string    `gorm:"column:ghichu;type:text" json:"note"`
	Difference        float64   `gorm:"column:chenhlech;type:decimal(10,4)" json:"difference"`
	CreatedDate       time.Time `gorm:"column:createddate;autoCreateTime" json:"createdDate"`

	// Relationships
	ReceivingDoc *ReceivingDoc `gorm:"foreignKey:ReceivingDocID;references:ReceivingDocID" json:"receivingDoc,omitempty"`
	Ingredient   *Ingredient   `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
	Supplier     *Supplier     `gorm:"foreignKey:SupplierID;references:SupplierID" json:"supplier,omitempty"`
}

func (ReceivingDetail) TableName() string {
	return "chitietnhanhang"
}

// OutboundStock - Stock outbound (xuatkho)
type OutboundStock struct {
	OutboundID   int        `gorm:"primaryKey;autoIncrement;column:xuatkhoid" json:"outboundId"`
	OutboundDate *time.Time `gorm:"column:ngayxuat" json:"outboundDate"`
	KitchenID    string     `gorm:"column:bepid" json:"kitchenId"`
	DishID       string     `gorm:"column:monanid" json:"dishId"`
	IngredientID string     `gorm:"column:nguyenlieuid" json:"ingredientId"`
	OutboundQty  float64    `gorm:"column:soluongxuat;type:decimal(10,4)" json:"outboundQty"`
	Unit         string     `gorm:"column:donvitinh" json:"unit"`
	UnitPrice    float64    `gorm:"column:dongia;type:decimal(15,2)" json:"unitPrice"`
	SupplierID   string     `gorm:"column:nhacungcapid" json:"supplierId"`
	CreatedDate  time.Time  `gorm:"column:createddate;autoCreateTime" json:"createdDate"`

	// Relationships
	Kitchen    *Kitchen    `gorm:"foreignKey:KitchenID;references:KitchenID" json:"kitchen,omitempty"`
	Dish       *Dish       `gorm:"foreignKey:DishID;references:DishID" json:"dish,omitempty"`
	Ingredient *Ingredient `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
	Supplier   *Supplier   `gorm:"foreignKey:SupplierID;references:SupplierID" json:"supplier,omitempty"`
}

func (OutboundStock) TableName() string {
	return "xuatkho"
}

// Inventory - Inventory balance by location (tonkho)
type Inventory struct {
	InventoryID   int        `gorm:"primaryKey;autoIncrement;column:tonkhoid" json:"inventoryId"`
	InventoryDate *time.Time `gorm:"column:ngayton" json:"inventoryDate"`
	KitchenID     string     `gorm:"column:bepid" json:"kitchenId"`
	IngredientID  string     `gorm:"column:nguyenlieuid" json:"ingredientId"`
	OpeningQty    float64    `gorm:"column:tondauky;type:decimal(10,4)" json:"openingQty"`
	InboundQty    float64    `gorm:"column:nhapkho;type:decimal(10,4)" json:"inboundQty"`
	OutboundQty   float64    `gorm:"column:xuatkho;type:decimal(10,4)" json:"outboundQty"`
	ClosingQty    float64    `gorm:"column:toncuoiky;type:decimal(10,4)" json:"closingQty"`
	Unit          string     `gorm:"column:donvitinh" json:"unit"`
	CreatedDate   time.Time  `gorm:"column:createddate;autoCreateTime" json:"createdDate"`

	// Relationships
	Kitchen    *Kitchen    `gorm:"foreignKey:KitchenID;references:KitchenID" json:"kitchen,omitempty"`
	Ingredient *Ingredient `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
}

func (Inventory) TableName() string {
	return "tonkho"
}

// Payable - Accounts payable to suppliers (congno)
type Payable struct {
	PayableID       int        `gorm:"primaryKey;autoIncrement;column:congnoid" json:"payableId"`
	Date            *time.Time `gorm:"column:ngay" json:"date"`
	SupplierID      string     `gorm:"column:nhacungcapid" json:"supplierId"`
	ReceivingDocID  string     `gorm:"column:phieunhanhangid" json:"receivingDocId"`
	IngredientID    string     `gorm:"column:nguyenlieuid" json:"ingredientId"`
	Quantity        float64    `gorm:"column:soluong;type:decimal(10,4)" json:"quantity"`
	UnitPrice       float64    `gorm:"column:dongia;type:decimal(15,2)" json:"unitPrice"`
	Amount          float64    `gorm:"column:thanhtien;type:decimal(15,2)" json:"amount"`
	PaymentStatus   string     `gorm:"column:trangthaithanhtoan" json:"paymentStatus"`
	PaidAmount      float64    `gorm:"column:sotiendatra;type:decimal(15,2)" json:"paidAmount"`
	RemainingAmount float64    `gorm:"column:conlai;type:decimal(15,2)" json:"remainingAmount"`
	CreatedDate     time.Time  `gorm:"column:createddate;autoCreateTime" json:"createdDate"`

	// Relationships
	Supplier     *Supplier     `gorm:"foreignKey:SupplierID;references:SupplierID" json:"supplier,omitempty"`
	ReceivingDoc *ReceivingDoc `gorm:"foreignKey:ReceivingDocID;references:ReceivingDocID" json:"receivingDoc,omitempty"`
	Ingredient   *Ingredient   `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
}

func (Payable) TableName() string {
	return "congno"
}

// StandardAdjustment - Recipe standard adjustments (dinhmuc_dieuchinh)
type StandardAdjustment struct {
	AdjustmentID     int        `gorm:"primaryKey;autoIncrement;column:dieuchinhid" json:"adjustmentId"`
	OrderID          string     `gorm:"column:phieulendonid" json:"orderId"`
	KitchenID        string     `gorm:"column:bepid" json:"kitchenId"`
	DishID           string     `gorm:"column:monanid" json:"dishId"`
	IngredientID     string     `gorm:"column:nguyenlieuid" json:"ingredientId"`
	AdjustedStandard float64    `gorm:"column:dinhmucdieuchinh;type:decimal(10,4)" json:"adjustedStandard"`
	AdjustmentReason string     `gorm:"column:lydodieuchinh;type:text" json:"adjustmentReason"`
	AdjustmentDate   *time.Time `gorm:"column:ngaydieuchinh" json:"adjustmentDate"`
	AdjustedByID     string     `gorm:"column:nguoidieuchinhid" json:"adjustedById"`
	Status           string     `gorm:"column:trangthai" json:"status"`
	ApprovedByID     string     `gorm:"column:nguoiduyetid" json:"approvedById"`
	ApprovedDate     *time.Time `gorm:"column:ngayduyet" json:"approvedDate"`
	CreatedDate      time.Time  `gorm:"column:createddate;autoCreateTime" json:"createdDate"`

	// Relationships
	Order      *Order      `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
	Kitchen    *Kitchen    `gorm:"foreignKey:KitchenID;references:KitchenID" json:"kitchen,omitempty"`
	Dish       *Dish       `gorm:"foreignKey:DishID;references:DishID" json:"dish,omitempty"`
	Ingredient *Ingredient `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
	AdjustedBy *User       `gorm:"foreignKey:AdjustedByID;references:UserID" json:"adjustedBy,omitempty"`
	ApprovedBy *User       `gorm:"foreignKey:ApprovedByID;references:UserID" json:"approvedBy,omitempty"`
}

func (StandardAdjustment) TableName() string {
	return "dinhmuc_dieuchinh"
}

// StandardIngredient - Links between recipe standards and ingredients (dinhmuc_nguyenlieu)
type StandardIngredient struct {
	StandardIngredientID string    `gorm:"primaryKey;column:dinhmucnlid" json:"standardIngredientId"`
	StandardID           *int      `gorm:"column:dinhmucid" json:"standardId"`
	OrderDetailID        string    `gorm:"column:chitietlendonid" json:"orderDetailId"`
	CreatedDate          time.Time `gorm:"column:createddate;autoCreateTime" json:"createdDate"`
	ModifiedDate         time.Time `gorm:"column:modifieddate;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	RecipeStandard *RecipeStandard `gorm:"foreignKey:StandardID;references:StandardID" json:"recipeStandard,omitempty"`
	OrderDetail    *OrderDetail    `gorm:"foreignKey:OrderDetailID;references:OrderDetailID" json:"orderDetail,omitempty"`
}

func (StandardIngredient) TableName() string {
	return "dinhmuc_nguyenlieu"
}
