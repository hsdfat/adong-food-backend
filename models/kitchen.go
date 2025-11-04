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
}

func (Kitchen) TableName() string {
    return "master_kitchens"
}
