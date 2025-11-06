package handler

import (
	"adong-be/logger"
	"adong-be/models"
	"adong-be/store"
	"adong-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetRecipeStandards with pagination and search - Returns ResourceCollection format with DTOs
func GetRecipeStandards(c *gin.Context) {
	logger.Log.Info("GetRecipeStandards called")
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Log.Error("GetRecipeStandards bind query error", "error", err)
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
	countDB := store.DB.GormClient.Model(&models.RecipeStandard{})

	searchConfig := utils.SearchConfig{
		Fields: []string{"dish_id", "ingredient_id"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		logger.Log.Error("GetRecipeStandards count error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var recipes []models.RecipeStandard
	db := store.DB.GormClient.Model(&models.RecipeStandard{})
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"standardId":   "recipe_id",
		"dishId":       "dish_id",
		"ingredientId": "ingredient_id",
		"standardPer1": "quantity_per_serving",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	// Preload related entities to get names
	db = db.Preload("Dish").Preload("Ingredient").Preload("UpdatedBy")

	if err := db.Find(&recipes).Error; err != nil {
		logger.Log.Error("GetRecipeStandards query error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to DTOs
	dtos := models.ConvertRecipeStandardsToDTO(recipes)

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: dtos,
		Meta: meta,
	})
}

func GetRecipeStandard(c *gin.Context) {
	logger.Log.Info("GetRecipeStandard called", "id", c.Param("id"))
	id := c.Param("id")
	var recipe models.RecipeStandard

	// Preload related entities
	if err := store.DB.GormClient.
		Preload("Dish").
		Preload("Ingredient").
		Preload("UpdatedBy").
		First(&recipe, "recipe_id = ?", id).Error; err != nil {
		logger.Log.Error("GetRecipeStandard not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe standard not found"})
		return
	}

	// Convert to DTO and return
	dto := recipe.ToDTO()
	c.JSON(http.StatusOK, dto)
}

func CreateRecipeStandard(c *gin.Context) {
	logger.Log.Info("CreateRecipeStandard called")
	var recipe models.RecipeStandard
	if err := c.ShouldBindJSON(&recipe); err != nil {
		logger.Log.Error("CreateRecipeStandard bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Create(&recipe).Error; err != nil {
		logger.Log.Error("CreateRecipeStandard db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload with relationships
	store.DB.GormClient.
		Preload("Dish").
		Preload("Ingredient").
		Preload("UpdatedBy").
		First(&recipe, "recipe_id = ?", recipe.StandardID)

	// Return DTO
	dto := recipe.ToDTO()
	c.JSON(http.StatusCreated, dto)
}

func UpdateRecipeStandard(c *gin.Context) {
	logger.Log.Info("UpdateRecipeStandard called", "id", c.Param("id"))
	id := c.Param("id")
	var recipe models.RecipeStandard
	if err := store.DB.GormClient.First(&recipe, "recipe_id = ?", id).Error; err != nil {
		logger.Log.Error("UpdateRecipeStandard not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe standard not found"})
		return
	}
	if err := c.ShouldBindJSON(&recipe); err != nil {
		logger.Log.Error("UpdateRecipeStandard bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Save(&recipe).Error; err != nil {
		logger.Log.Error("UpdateRecipeStandard db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload with relationships
	store.DB.GormClient.
		Preload("Dish").
		Preload("Ingredient").
		Preload("UpdatedBy").
		First(&recipe, "recipe_id = ?", recipe.StandardID)

	// Return DTO
	dto := recipe.ToDTO()
	c.JSON(http.StatusOK, dto)
}

func DeleteRecipeStandard(c *gin.Context) {
	logger.Log.Info("DeleteRecipeStandard called", "id", c.Param("id"))
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.RecipeStandard{}, "recipe_id = ?", id).Error; err != nil {
		logger.Log.Error("DeleteRecipeStandard db error", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe standard deleted successfully"})
}

// GetRecipeStandardsByDish with pagination and search - Returns ResourceCollection format with DTOs
func GetRecipeStandardsByDish(c *gin.Context) {
	logger.Log.Info("GetRecipeStandardsByDish called", "dishId", c.Param("dishId"))
	dishId := c.Param("dishId")

	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
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
	countDB := store.DB.GormClient.Model(&models.RecipeStandard{}).Where("dish_id = ?", dishId)

	searchConfig := utils.SearchConfig{
		Fields: []string{"ingredient_id"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		logger.Log.Error("GetRecipeStandardsByDish count error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var recipes []models.RecipeStandard
	db := store.DB.GormClient.Model(&models.RecipeStandard{}).Where("dish_id = ?", dishId)
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"ingredientId": "ingredient_id",
		"standardPer1": "quantity_per_serving",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	// Preload related entities to get names
	db = db.Preload("Dish").Preload("Ingredient").Preload("UpdatedBy")

	if err := db.Find(&recipes).Error; err != nil {
		logger.Log.Error("GetRecipeStandardsByDish query error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to DTOs
	dtos := models.ConvertRecipeStandardsToDTO(recipes)

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: dtos,
		Meta: meta,
	})
}
