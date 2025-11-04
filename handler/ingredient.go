package handler

import (
	"adong-be/logger"
	"adong-be/models"
	"adong-be/store"
	"adong-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetIngredients with pagination and search - Returns ResourceCollection format
func GetIngredients(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("GetIngredients called", "user_id", uid)
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Log.Error("GetIngredients bind query error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	params = models.GetPaginationParams(
		params.Page,
		params.PageSize,
		params.Search,
		params.SortBy,
		params.SortDir,
	)

	var total int64
	countDB := store.DB.GormClient.Model(&models.Ingredient{})

	searchConfig := utils.SearchConfig{
		Fields: []string{"ingredient_name", "ingredient_id"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		logger.Log.Error("GetIngredients count error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var items []models.Ingredient
	db := store.DB.GormClient.Model(&models.Ingredient{})
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"ingredient_id":   "ingredient_id",
		"ingredient_name": "ingredient_name",
		"unit":            "unit",
		"created_date":    "created_date",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	if err := db.Find(&items).Error; err != nil {
		logger.Log.Error("GetIngredients query error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: items,
		Meta: meta,
	})
}

func GetIngredient(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("GetIngredient called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	var item models.Ingredient
	if err := store.DB.GormClient.First(&item, "ingredient_id = ?", id).Error; err != nil {
		logger.Log.Error("GetIngredient not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Ingredient not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func CreateIngredient(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("CreateIngredient called", "user_id", uid)
	var item models.Ingredient
	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Log.Error("CreateIngredient bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Create(&item).Error; err != nil {
		logger.Log.Error("CreateIngredient db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func UpdateIngredient(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("UpdateIngredient called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	var item models.Ingredient
	if err := store.DB.GormClient.First(&item, "ingredient_id = ?", id).Error; err != nil {
		logger.Log.Error("UpdateIngredient not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Ingredient not found"})
		return
	}
	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Log.Error("UpdateIngredient bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Save(&item).Error; err != nil {
		logger.Log.Error("UpdateIngredient db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func DeleteIngredient(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("DeleteIngredient called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.Ingredient{}, "ingredient_id = ?", id).Error; err != nil {
		logger.Log.Error("DeleteIngredient db error", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ingredient deleted successfully"})
}
