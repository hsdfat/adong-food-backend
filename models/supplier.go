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
