package handler

import (
	"adong-be/models"
	"adong-be/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recipe Standard handlers
func GetRecipeStandards(c *gin.Context) {
	var standards []models.RecipeStandard
	if err := store.DB.GormClient.Find(&standards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, standards)
}

func GetRecipeStandard(c *gin.Context) {
	id := c.Param("id")
	var standard models.RecipeStandard
	if err := store.DB.GormClient.First(&standard, "dinhmucid = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe standard not found"})
		return
	}
	c.JSON(http.StatusOK, standard)
}

func GetRecipeStandardsByDish(c *gin.Context) {
	dishId := c.Param("dishId")
	var standards []models.RecipeStandard
	if err := store.DB.GormClient.Where("monanid = ?", dishId).Find(&standards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, standards)
}

func CreateRecipeStandard(c *gin.Context) {
	var standard models.RecipeStandard
	if err := c.ShouldBindJSON(&standard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Create(&standard).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, standard)
}

func UpdateRecipeStandard(c *gin.Context) {
	id := c.Param("id")
	var standard models.RecipeStandard
	if err := store.DB.GormClient.First(&standard, "dinhmucid = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe standard not found"})
		return
	}
	if err := c.ShouldBindJSON(&standard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Save(&standard).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, standard)
}

func DeleteRecipeStandard(c *gin.Context) {
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.RecipeStandard{}, "dinhmucid = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe standard deleted successfully"})
}
