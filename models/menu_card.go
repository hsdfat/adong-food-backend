// File: models/menu_card.go
package models

import "time"

// MenuCard - Phiếu thực đơn (phieuthucdon)
type MenuCard struct {
	MenuCardID   string     `gorm:"primaryKey;column:phieuthucdonid" json:"menuCardId"`
	MenuCardName string     `gorm:"column:tenphieu;not null" json:"menuCardName"`
	CreatedDate  *time.Time `gorm:"column:ngaytao" json:"createdDate"`
	KitchenID    string     `gorm:"column:bepid" json:"kitchenId"`
	CreatedByID  string     `gorm:"column:nguoitaoid" json:"createdById"`
	Status       string     `gorm:"column:trangthai;default:DRAFT" json:"status"`
	Note         string     `gorm:"column:ghichu;type:text" json:"note"`
	CreatedAt    time.Time  `gorm:"column:createddate;autoCreateTime" json:"createdAt"`
	ModifiedDate time.Time  `gorm:"column:modifieddate;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Kitchen   *Kitchen         `gorm:"foreignKey:KitchenID;references:KitchenID" json:"kitchen,omitempty"`
	CreatedBy *User            `gorm:"foreignKey:CreatedByID;references:UserID" json:"createdBy,omitempty"`
	Details   []MenuCardDetail `gorm:"foreignKey:MenuCardID;references:MenuCardID" json:"details,omitempty"`
}

func (MenuCard) TableName() string {
	return "phieuthucdon"
}

// MenuCardDetail - Chi tiết món ăn trong phiếu thực đơn (chitietphieuthucdon)
type MenuCardDetail struct {
	DetailID     string                     `gorm:"primaryKey;column:chitietid" json:"detailId"`
	MenuCardID   string                     `gorm:"column:phieuthucdonid;not null" json:"menuCardId"`
	DishID       string                     `gorm:"column:monanid;not null" json:"dishId"`
	DishName     string                     `gorm:"column:tenmonan" json:"dishName"`
	Servings     int                        `gorm:"column:sosuat;default:1" json:"servings"`
	Note         string                     `gorm:"column:ghichu;type:text" json:"note"`
	CreatedDate  time.Time                  `gorm:"column:createddate;autoCreateTime" json:"createdDate"`

	// Relationships
	MenuCard    *MenuCard                   `gorm:"foreignKey:MenuCardID;references:MenuCardID" json:"menuCard,omitempty"`
	Dish        *Dish                       `gorm:"foreignKey:DishID;references:DishID" json:"dish,omitempty"`
	Ingredients []MenuCardDetailIngredient  `gorm:"foreignKey:DetailID;references:DetailID" json:"ingredients,omitempty"`
}

func (MenuCardDetail) TableName() string {
	return "chitietphieuthucdon"
}

// MenuCardDetailIngredient - Nguyên liệu tùy chỉnh (nguyenlieu_phieuthucdon)
type MenuCardDetailIngredient struct {
	ID           string    `gorm:"primaryKey;column:id" json:"id"`
	DetailID     string    `gorm:"column:chitietid;not null" json:"detailId"`
	IngredientID string    `gorm:"column:nguyenlieuid;not null" json:"ingredientId"`
	Standard     float64   `gorm:"column:dinhmuc;type:decimal(10,4);not null" json:"standard"`
	Unit         string    `gorm:"column:donvitinh" json:"unit"`
	Note         string    `gorm:"column:ghichu;type:text" json:"note"`
	CreatedDate  time.Time `gorm:"column:createddate;autoCreateTime" json:"createdDate"`

	// Relationships
	Detail     *MenuCardDetail `gorm:"foreignKey:DetailID;references:DetailID" json:"detail,omitempty"`
	Ingredient *Ingredient     `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
}

func (MenuCardDetailIngredient) TableName() string {
	return "nguyenlieu_phieuthucdon"
}