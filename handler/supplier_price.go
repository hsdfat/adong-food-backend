package handler

import (
	"adong-be/models"
	"adong-be/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Supplier Price handlers
func GetSupplierPrices(c *gin.Context) {
	var prices []models.SupplierPrice
	if err := store.DB.GormClient.Find(&prices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prices)
}

func GetSupplierPrice(c *gin.Context) {
	id := c.Param("id")
	var price models.SupplierPrice
	if err := store.DB.GormClient.First(&price, "sanphamid = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Supplier price not found"})
		return
	}
	c.JSON(http.StatusOK, price)
}

func CreateSupplierPrice(c *gin.Context) {
	var price models.SupplierPrice
	if err := c.ShouldBindJSON(&price); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Create(&price).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, price)
}

func UpdateSupplierPrice(c *gin.Context) {
	id := c.Param("id")
	var price models.SupplierPrice
	if err := store.DB.GormClient.First(&price, "sanphamid = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Supplier price not found"})
		return
	}
	if err := c.ShouldBindJSON(&price); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Save(&price).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, price)
}

func DeleteSupplierPrice(c *gin.Context) {
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.SupplierPrice{}, "sanphamid = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Supplier price deleted successfully"})
}
