package handler

import (
	"adong-be/models"
	"adong-be/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Kitchen handlers
func GetKitchens(c *gin.Context) {
	var kitchens []models.Kitchen
	if err := store.DB.GormClient.Find(&kitchens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kitchens)
}

func GetKitchen(c *gin.Context) {
	id := c.Param("id")
	var kitchen models.Kitchen
	if err := store.DB.GormClient.First(&kitchen, "bepid = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitchen not found"})
		return
	}
	c.JSON(http.StatusOK, kitchen)
}

func CreateKitchen(c *gin.Context) {
	var kitchen models.Kitchen
	if err := c.ShouldBindJSON(&kitchen); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Create(&kitchen).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, kitchen)
}

func UpdateKitchen(c *gin.Context) {
	id := c.Param("id")
	var kitchen models.Kitchen
	if err := store.DB.GormClient.First(&kitchen, "bepid = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitchen not found"})
		return
	}
	if err := c.ShouldBindJSON(&kitchen); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Save(&kitchen).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kitchen)
}

func DeleteKitchen(c *gin.Context) {
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.Kitchen{}, "bepid = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kitchen deleted successfully"})
}
