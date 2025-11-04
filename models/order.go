// models/order.go
package models

import (
	"time"
)

type OrderForm struct {
	OrderFormID        string                `json:"orderFormId" gorm:"primaryKey;column:phieu_len_don_id"`
	KitchenID          string                `json:"kitchenId" gorm:"column:bep_id"`
	OrderDate          string            `json:"orderDate" gorm:"column:ngay_len;type:date"`
	Note               string                `json:"note" gorm:"column:ghi_chu"`
	Status             string                `json:"status" gorm:"column:trang_thai;default:'Pending'"`
	CreatedBy          string                `json:"createdBy" gorm:"column:created_by"`
	CreatedAt          time.Time             `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          time.Time             `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
	Details            []OrderFormDetail     `json:"details" gorm:"foreignKey:OrderFormID;references:OrderFormID"`
	SupplementaryFoods []SupplementaryFood   `json:"supplementaryFoods" gorm:"foreignKey:OrderFormID;references:OrderFormID"`
}

type OrderFormDetail struct {
	ID             int                      `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	OrderFormID    string                   `json:"orderFormId" gorm:"column:phieu_len_don_id"`
	DishID         string                   `json:"dishId" gorm:"column:monan_id"`
	DishName       string                   `json:"dishName" gorm:"column:ten_mon_an"`
	Portions       int                      `json:"portions" gorm:"column:so_suat"`
	Ingredients    []IngredientDetail       `json:"ingredients" gorm:"foreignKey:DetailID;references:ID"`
}

type IngredientDetail struct {
	ID               int     `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	DetailID         int     `json:"detailId" gorm:"column:chi_tiet_id"`
	IngredientID     string  `json:"ingredientId" gorm:"column:nguyen_lieu_id"`
	IngredientName   string  `json:"ingredientName" gorm:"column:ten_nguyen_lieu"`
	Quantity         float64 `json:"quantity" gorm:"column:so_luong;type:decimal(10,2)"`
	Unit             string  `json:"unit" gorm:"column:don_vi_tinh"`
	StandardPerPortion float64 `json:"standardPerPortion" gorm:"column:dinh_muc;type:decimal(10,2);default:0"`
}

type SupplementaryFood struct {
	ID               int     `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	OrderFormID      string  `json:"orderFormId" gorm:"column:phieu_len_don_id"`
	IngredientID     string  `json:"ingredientId" gorm:"column:nguyen_lieu_id"`
	IngredientName   string  `json:"ingredientName" gorm:"column:ten_nguyen_lieu"`
	Quantity         float64 `json:"quantity" gorm:"column:so_luong;type:decimal(10,2)"`
	Unit             string  `json:"unit" gorm:"column:don_vi_tinh"`
	StandardPerPortion float64 `json:"standardPerPortion" gorm:"column:dinh_muc;type:decimal(10,2);default:0"`
	Portions         float64 `json:"portions" gorm:"column:so_suat;type:decimal(10,2);default:1"`
	Note             string  `json:"note" gorm:"column:ghi_chu"`
}

type TotalIngredientSummary struct {
	IngredientID     string  `json:"ingredientId"`
	IngredientName   string  `json:"ingredientName"`
	Quantity         float64 `json:"quantity"`
	Unit             string  `json:"unit"`
}

type DishWithIngredients struct {
	DishID      string                `json:"dishId"`
	DishName    string                `json:"dishName"`
	Ingredients []IngredientInRecipe  `json:"ingredients"`
}

type IngredientInRecipe struct {
	IngredientID       string  `json:"ingredientId"`
	IngredientName     string  `json:"ingredientName"`
	Unit               string  `json:"unit"`
	StandardPerPortion float64 `json:"standardPerPortion"`
}

func (OrderForm) TableName() string {
	return "phieu_len_don"
}

func (OrderFormDetail) TableName() string {
	return "chi_tiet_phieu_len_don"
}

func (IngredientDetail) TableName() string {
	return "nguyen_lieu_chi_tiet"
}

func (SupplementaryFood) TableName() string {
	return "thuc_pham_bo_sung"
}