package handler

import (
	"adong-be/models"
	"adong-be/store"
	"adong-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDishes with pagination and search
func GetDishes(c *gin.Context) {
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
	countDB := store.DB.GormClient.Model(&models.Dish{})

	searchConfig := utils.SearchConfig{
		Fields: []string{"tenmonan", "monanid", "mota"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var dishes []models.Dish
	db := store.DB.GormClient.Model(&models.Dish{})
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"monanid":   "monanid",
		"tenmonan":  "tenmonan",
		"loaimonan": "loaimonan",
		"dongia":    "dongia",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	if err := db.Find(&dishes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: dishes,
		Meta: meta,
	})
}

func GetDish(c *gin.Context) {
	id := c.Param("id")
	var dish models.Dish
	if err := store.DB.GormClient.First(&dish, "monanid = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dish not found"})
		return
	}
	c.JSON(http.StatusOK, dish)
}

func CreateDish(c *gin.Context) {
	var dish models.Dish
	if err := c.ShouldBindJSON(&dish); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Create(&dish).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dish)
}

func UpdateDish(c *gin.Context) {
	id := c.Param("id")
	var dish models.Dish
	if err := store.DB.GormClient.First(&dish, "monanid = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dish not found"})
		return
	}
	if err := c.ShouldBindJSON(&dish); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Save(&dish).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dish)
}

func DeleteDish(c *gin.Context) {
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.Dish{}, "monanid = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Dish deleted successfully"})
}
