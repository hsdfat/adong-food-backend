package models

import "time"

// Order - Order forms (phieulendon)
type Order struct {
	OrderID      string     `gorm:"primaryKey;column:phieulendonid" json:"orderId"`
	OrderDate    *time.Time `gorm:"column:ngaylendon" json:"orderDate"`
	KitchenID    string     `gorm:"column:bepid" json:"kitchenId"`
	Status       string     `gorm:"column:trangthai" json:"status"`
	CreatedByID  string     `gorm:"column:nguoitaoid" json:"createdById"`
	Note         string     `gorm:"column:ghichu;type:text" json:"note"`
	CreatedDate  time.Time  `gorm:"column:createddate;autoCreateTime" json:"createdDate"`
	ModifiedDate time.Time  `gorm:"column:modifieddate;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Kitchen      *Kitchen      `gorm:"foreignKey:KitchenID;references:KitchenID" json:"kitchen,omitempty"`
	CreatedBy    *User         `gorm:"foreignKey:CreatedByID;references:UserID" json:"createdBy,omitempty"`
	OrderDetails []OrderDetail `gorm:"foreignKey:OrderID;references:OrderID" json:"orderDetails,omitempty"`
}

func (Order) TableName() string {
	return "phieulendon"
}

// OrderDetail - Order line items (chitietlendon)
type OrderDetail struct {
	OrderDetailID  string     `gorm:"primaryKey;column:chitietlendonid" json:"orderDetailId"`
	OrderID        string     `gorm:"column:phieulendonid" json:"orderId"`
	DishID         string     `gorm:"column:monanid" json:"dishId"`
	DishName       string     `gorm:"column:tenmonan" json:"dishName"`
	Servings       int        `gorm:"column:sosuat" json:"servings"`
	IngredientList string     `gorm:"column:listnguyenlieu;type:text" json:"ingredientList"`
	LastModify     *time.Time `gorm:"column:lastmodify" json:"lastModify"`
	CreatedDate    time.Time  `gorm:"column:createddate;autoCreateTime" json:"createdDate"`

	// Relationships
	Order *Order `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
	Dish  *Dish  `gorm:"foreignKey:DishID;references:DishID" json:"dish,omitempty"`
}

func (OrderDetail) TableName() string {
	return "chitietlendon"
}
