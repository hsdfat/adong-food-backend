package handler

import (
	"adong-be/models"
	"adong-be/logger"
	"adong-be/store"
	"adong-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUsers with pagination and search - Returns ResourceCollection format
func GetUsers(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("GetUsers called", "user_id", uid)
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Log.Error("GetUsers bind query error", "error", err)
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
	countDB := store.DB.GormClient.Model(&models.User{})

	searchConfig := utils.SearchConfig{
		Fields: []string{"user_id", "user_name", "full_name", "email", "phone"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		logger.Log.Error("GetUsers count error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var items []models.User
	db := store.DB.GormClient.Model(&models.User{})
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"user_id":   "user_id",
		"user_name": "user_name",
		"full_name": "full_name",
		"email":     "email",
		"role":      "role",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	if err := db.Find(&items).Error; err != nil {
		logger.Log.Error("GetUsers query error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: items,
		Meta: meta,
	})
}

func GetUser(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("GetUser called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	var item models.User
	if err := store.DB.GormClient.First(&item, "user_id = ?", id).Error; err != nil {
		logger.Log.Error("GetUser not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func CreateUser(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("CreateUser called", "user_id", uid)
	var item models.User
	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Log.Error("CreateUser bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Create(&item).Error; err != nil {
		logger.Log.Error("CreateUser db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func UpdateUser(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("UpdateUser called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	var item models.User
	if err := store.DB.GormClient.First(&item, "user_id = ?", id).Error; err != nil {
		logger.Log.Error("UpdateUser not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Log.Error("UpdateUser bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Save(&item).Error; err != nil {
		logger.Log.Error("UpdateUser db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func DeleteUser(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("DeleteUser called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.User{}, "user_id = ?", id).Error; err != nil {
		logger.Log.Error("DeleteUser db error", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
