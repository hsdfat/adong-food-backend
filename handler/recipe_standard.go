package handler

import (
	"adong-be/models"
	"adong-be/store"
	"adong-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetRecipeStandards with pagination and search - Returns ResourceCollection format
func GetRecipeStandards(c *gin.Context) {
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
	countDB := store.DB.GormClient.Model(&models.RecipeStandard{})

	searchConfig := utils.SearchConfig{
		Fields: []string{"monanid", "nguyenlieuid"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var recipes []models.RecipeStandard
	db := store.DB.GormClient.Model(&models.RecipeStandard{})
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"dinhmucid":    "dinhmucid",
		"monanid":      "monanid",
		"nguyenlieuid": "nguyenlieuid",
		"dinhmuc":      "dinhmuc",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	if err := db.Find(&recipes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: recipes,
		Meta: meta,
	})
}

func GetRecipeStandard(c *gin.Context) {
	id := c.Param("id")
	var recipe models.RecipeStandard
	if err := store.DB.GormClient.First(&recipe, "dinhmucid = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe standard not found"})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

func CreateRecipeStandard(c *gin.Context) {
	var recipe models.RecipeStandard
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Create(&recipe).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, recipe)
}

func UpdateRecipeStandard(c *gin.Context) {
	id := c.Param("id")
	var recipe models.RecipeStandard
	if err := store.DB.GormClient.First(&recipe, "dinhmucid = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe standard not found"})
		return
	}
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Save(&recipe).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

func DeleteRecipeStandard(c *gin.Context) {
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.RecipeStandard{}, "dinhmucid = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe standard deleted successfully"})
}

// GetRecipeStandardsByDish with pagination and search - Returns ResourceCollection format
func GetRecipeStandardsByDish(c *gin.Context) {
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
	countDB := store.DB.GormClient.Model(&models.RecipeStandard{}).Where("monanid = ?", dishId)

	searchConfig := utils.SearchConfig{
		Fields: []string{"nguyenlieuid"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var recipes []models.RecipeStandard
	db := store.DB.GormClient.Model(&models.RecipeStandard{}).Where("monanid = ?", dishId)
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"nguyenlieuid": "nguyenlieuid",
		"dinhmuc":      "dinhmuc",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	if err := db.Find(&recipes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: recipes,
		Meta: meta,
	})
}
