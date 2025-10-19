package handler

import (
	"adong-be/models"
	"adong-be/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ingredient handlers
func GetIngredients(c *gin.Context) {
	var ingredients []models.Ingredient
	if err := store.DB.GormClient.Find(&ingredients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ingredients)
}

func GetIngredient(c *gin.Context) {
	id := c.Param("id")
	var ingredient models.Ingredient
	if err := store.DB.GormClient.First(&ingredient, "nguyenlieuid = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ingredient not found"})
		return
	}
	c.JSON(http.StatusOK, ingredient)
}

func CreateIngredient(c *gin.Context) {
	var ingredient models.Ingredient
	if err := c.ShouldBindJSON(&ingredient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Create(&ingredient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ingredient)
}

func UpdateIngredient(c *gin.Context) {
	id := c.Param("id")
	var ingredient models.Ingredient
	if err := store.DB.GormClient.First(&ingredient, "nguyenlieuid = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ingredient not found"})
		return
	}
	if err := c.ShouldBindJSON(&ingredient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Save(&ingredient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ingredient)
}

func DeleteIngredient(c *gin.Context) {
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.Ingredient{}, "nguyenlieuid = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ingredient deleted successfully"})
}
