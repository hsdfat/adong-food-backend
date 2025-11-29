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

	// Relationships
	// Many-to-many: one kitchen can have many users
	Users             []User                    `gorm:"many2many:user_kitchens" json:"users,omitempty"`
	FavoriteSuppliers []KitchenFavoriteSupplier `gorm:"foreignKey:KitchenID;references:KitchenID" json:"favoriteSuppliers,omitempty"`
}

func (Kitchen) TableName() string {
	return "master_kitchens"
}

// KitchenFavoriteSupplier - Kitchen's favorite suppliers
type KitchenFavoriteSupplier struct {
	FavoriteID      int       `gorm:"primaryKey;autoIncrement;column:favorite_id" json:"favoriteId"`
	KitchenID       string    `gorm:"column:kitchen_id;not null" json:"kitchenId"`
	SupplierID      string    `gorm:"column:supplier_id;not null" json:"supplierId"`
	Notes           string    `gorm:"column:notes;type:text" json:"notes"`
	DisplayOrder    int       `gorm:"column:display_order" json:"displayOrder"`
	CreatedByUserID string    `gorm:"column:created_by_user_id" json:"createdByUserId"`
	CreatedDate     time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate    time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Kitchen    *Kitchen  `gorm:"foreignKey:KitchenID;references:KitchenID" json:"kitchen,omitempty"`
	Supplier   *Supplier `gorm:"foreignKey:SupplierID;references:SupplierID" json:"supplier,omitempty"`
	CreatedBy  *User     `gorm:"foreignKey:CreatedByUserID;references:UserID" json:"createdBy,omitempty"`
}

func (KitchenFavoriteSupplier) TableName() string {
	return "kitchen_favorite_suppliers"
}