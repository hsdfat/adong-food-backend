package handler

import (
	"adong-be/logger"
	"adong-be/models"
	"adong-be/store"
	"adong-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDishes with pagination and search
func GetDishes(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("GetDishes called", "user_id", uid)
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Log.Error("GetDishes bind query error", "error", err)
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
		Fields: []string{"dish_name", "dish_id", "description"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		logger.Log.Error("GetDishes count error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var dishes []models.Dish
	db := store.DB.GormClient.Model(&models.Dish{})
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"dish_id":        "dish_id",
		"dish_name":      "dish_name",
		"cooking_method": "cooking_method",
		"category":       "category",
		"created_date":   "created_date",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	if err := db.Find(&dishes).Error; err != nil {
		logger.Log.Error("GetDishes query error", "error", err)
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
	uid, _ := c.Get("identity")
	logger.Log.Info("GetDish called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	var dish models.Dish
	if err := store.DB.GormClient.First(&dish, "dish_id = ?", id).Error; err != nil {
		logger.Log.Error("GetDish not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Dish not found"})
		return
	}
	c.JSON(http.StatusOK, dish)
}

func CreateDish(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("CreateDish called", "user_id", uid)
	var dish models.Dish
	if err := c.ShouldBindJSON(&dish); err != nil {
		logger.Log.Error("CreateDish bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Create(&dish).Error; err != nil {
		logger.Log.Error("CreateDish db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dish)
}

func UpdateDish(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("UpdateDish called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	var dish models.Dish
	if err := store.DB.GormClient.First(&dish, "dish_id = ?", id).Error; err != nil {
		logger.Log.Error("UpdateDish not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Dish not found"})
		return
	}
	if err := c.ShouldBindJSON(&dish); err != nil {
		logger.Log.Error("UpdateDish bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Save(&dish).Error; err != nil {
		logger.Log.Error("UpdateDish db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dish)
}

func DeleteDish(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("DeleteDish called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.Dish{}, "dish_id = ?", id).Error; err != nil {
		logger.Log.Error("DeleteDish db error", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Dish deleted successfully"})
}
