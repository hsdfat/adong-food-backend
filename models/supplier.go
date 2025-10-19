package models

import "time"

// Supplier - Master data for suppliers (dm_ncc)
type Supplier struct {
	SupplierID   string    `gorm:"primaryKey;column:nhacungcapid" json:"supplierId"`
	SupplierName string    `gorm:"column:tennhacungcap;not null" json:"supplierName"`
	ZaloLink     string    `gorm:"column:linkzalo;type:text" json:"zaloLink"`
	Address      string    `gorm:"column:diachi;type:text" json:"address"`
	Phone        string    `gorm:"column:sodienthoai" json:"phone"`
	Email        string    `gorm:"column:email" json:"email"`
	Active       *bool     `gorm:"column:active;default:true" json:"active"`
	CreatedDate  time.Time `gorm:"column:createddate;autoCreateTime" json:"createdDate"`
	ModifiedDate time.Time `gorm:"column:modifieddate;autoUpdateTime" json:"modifiedDate"`
}

func (Supplier) TableName() string {
	return "dm_ncc"
}
