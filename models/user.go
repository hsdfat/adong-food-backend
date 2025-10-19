package models

import (
	"time"
)






// User - Master data for user accounts (dm_nguoidung)
type User struct {
	UserID       string    `gorm:"primaryKey;column:userid" json:"userId"`
	Password     string    `gorm:"column:password;not null" json:"password,omitempty"`
	FullName     string    `gorm:"column:hoten;not null" json:"fullName"`
	Role         string    `gorm:"column:vaitro" json:"role"`
	KitchenID    string    `gorm:"column:bepid" json:"kitchenId"`
	Email        string    `gorm:"column:email" json:"email"`
	Phone        string    `gorm:"column:sodienthoai" json:"phone"`
	Active       *bool     `gorm:"column:active;default:true" json:"active"`
	CreatedDate  time.Time `gorm:"column:createddate;autoCreateTime" json:"createdDate"`
	ModifiedDate time.Time `gorm:"column:modifieddate;autoUpdateTime" json:"modifiedDate"`
	
	// Relationships
	Kitchen      *Kitchen  `gorm:"foreignKey:KitchenID;references:KitchenID" json:"kitchen,omitempty"`
}

func (User) TableName() string {
	return "dm_nguoidung"
}






