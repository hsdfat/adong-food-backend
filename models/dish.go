package models

import "time"

// Dish - Master data for dishes/menu items (dm_monan)
type Dish struct {
    DishID        string    `gorm:"primaryKey;column:dish_id" json:"dishId"`
    DishName      string    `gorm:"column:dish_name;not null" json:"dishName"`
    CookingMethod string    `gorm:"column:cooking_method" json:"cookingMethod"`
    Group         string    `gorm:"column:category" json:"group"`
    Description   string    `gorm:"column:description;type:text" json:"description"`
	Active        *bool     `gorm:"column:active;default:true" json:"active"`
    CreatedDate   time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
    ModifiedDate  time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`
}

func (Dish) TableName() string {
    return "master_dishes"
}
