package models

import "time"

// Dish - Master data for dishes/menu items (dm_monan)
type Dish struct {
	DishID        string    `gorm:"primaryKey;column:monanid" json:"dishId"`
	DishName      string    `gorm:"column:tenmonan;not null" json:"dishName"`
	CookingMethod string    `gorm:"column:kieuchebien" json:"cookingMethod"`
	Group         string    `gorm:"column:nhom" json:"group"`
	Description   string    `gorm:"column:mota;type:text" json:"description"`
	Active        *bool     `gorm:"column:active;default:true" json:"active"`
	CreatedDate   time.Time `gorm:"column:createddate;autoCreateTime" json:"createdDate"`
	ModifiedDate  time.Time `gorm:"column:modifieddate;autoUpdateTime" json:"modifiedDate"`
}

func (Dish) TableName() string {
	return "dm_monan"
}
