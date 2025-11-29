package models

import "time"






// User - Master data for user accounts (dm_nguoidung)
type User struct {
	UserID       string    `gorm:"primaryKey;column:user_id" json:"userId"`
	UserName     string    `gorm:"column:user_name;not null;unique" json:"userName"`
	Password     string    `gorm:"column:password;not null" json:"password,omitempty"`
	FullName     string    `gorm:"column:full_name;not null" json:"fullName"`
	Role         string    `gorm:"column:role" json:"role"`
	Email        string    `gorm:"column:email" json:"email"`
	Phone        string    `gorm:"column:phone" json:"phone"`
	Active       *bool     `gorm:"column:active;default:true" json:"active"`
	CreatedDate  time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	// Many-to-many: one user can work in many kitchens, one kitchen can have many users
	Kitchens []Kitchen `gorm:"many2many:user_kitchens" json:"kitchens,omitempty"`
}

func (User) TableName() string {
	return "master_users"
}

