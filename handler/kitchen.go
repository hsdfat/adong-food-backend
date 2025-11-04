package handler

import (
	"adong-be/models"
	"adong-be/logger"
	"adong-be/store"
	"adong-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetKitchens with pagination and search - Returns ResourceCollection format
func GetKitchens(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("GetKitchens called", "user_id", uid)
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Log.Error("GetKitchens bind query error", "error", err)
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
	countDB := store.DB.GormClient.Model(&models.Kitchen{})

	searchConfig := utils.SearchConfig{
		Fields: []string{"kitchen_name", "kitchen_id", "address"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		logger.Log.Error("GetKitchens count error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var items []models.Kitchen
	db := store.DB.GormClient.Model(&models.Kitchen{})
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"kitchen_id":   "kitchen_id",
		"kitchen_name": "kitchen_name",
		"address":      "address",
		"created_date": "created_date",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	if err := db.Find(&items).Error; err != nil {
		logger.Log.Error("GetKitchens query error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: items,
		Meta: meta,
	})
}

func GetKitchen(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("GetKitchen called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	var item models.Kitchen
	if err := store.DB.GormClient.First(&item, "kitchen_id = ?", id).Error; err != nil {
		logger.Log.Error("GetKitchen not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitchen not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func CreateKitchen(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("CreateKitchen called", "user_id", uid)
	var item models.Kitchen
	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Log.Error("CreateKitchen bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Create(&item).Error; err != nil {
		logger.Log.Error("CreateKitchen db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func UpdateKitchen(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("UpdateKitchen called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	var item models.Kitchen
	if err := store.DB.GormClient.First(&item, "kitchen_id = ?", id).Error; err != nil {
		logger.Log.Error("UpdateKitchen not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitchen not found"})
		return
	}
	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Log.Error("UpdateKitchen bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Save(&item).Error; err != nil {
		logger.Log.Error("UpdateKitchen db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func DeleteKitchen(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("DeleteKitchen called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.Kitchen{}, "kitchen_id = ?", id).Error; err != nil {
		logger.Log.Error("DeleteKitchen db error", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kitchen deleted successfully"})
}
