package models

import "time"

// Kitchen - Master data for kitchen/location information (dm_bep)
type Kitchen struct {
	KitchenID    string    `gorm:"primaryKey;column:bepid" json:"kitchenId"`
	KitchenName  string    `gorm:"column:tenbep;not null" json:"kitchenName"`
	Address      string    `gorm:"column:diachi;type:text" json:"address"`
	Phone        string    `gorm:"column:sodienthoai" json:"phone"`
	Active       *bool     `gorm:"column:active;default:true" json:"active"`
	CreatedDate  time.Time `gorm:"column:createddate;autoCreateTime" json:"createdDate"`
	ModifiedDate time.Time `gorm:"column:modifieddate;autoUpdateTime" json:"modifiedDate"`
}

func (Kitchen) TableName() string {
	return "dm_bep"
}
