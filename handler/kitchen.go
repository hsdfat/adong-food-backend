package handler

import (
	"adong-be/logger"
	"adong-be/models"
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

// ============================================================================
// Kitchen Favorite Suppliers Handlers
// ============================================================================

// GetKitchenFavoriteSuppliers returns all favorite suppliers for a kitchen
func GetKitchenFavoriteSuppliers(c *gin.Context) {
	uid, _ := c.Get("identity")
	kitchenID := c.Param("id")
	logger.Log.Info("GetKitchenFavoriteSuppliers called", "kitchen_id", kitchenID, "user_id", uid)

	// Validate kitchen exists
	var kitchen models.Kitchen
	if err := store.DB.GormClient.First(&kitchen, "kitchen_id = ?", kitchenID).Error; err != nil {
		logger.Log.Error("GetKitchenFavoriteSuppliers kitchen not found", "kitchen_id", kitchenID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitchen not found"})
		return
	}

	var favorites []models.KitchenFavoriteSupplier
	query := store.DB.GormClient.
		Where("kitchen_id = ?", kitchenID).
		Preload("Supplier").
		Preload("CreatedBy")

	// Order by display_order if set, otherwise by created_date
	query = query.Order("COALESCE(display_order, 999999), created_date ASC")

	if err := query.Find(&favorites).Error; err != nil {
		logger.Log.Error("GetKitchenFavoriteSuppliers db error", "kitchen_id", kitchenID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Count total favorites for meta info
	var total int64
	if err := store.DB.GormClient.Model(&models.KitchenFavoriteSupplier{}).Where("kitchen_id = ?", kitchenID).Count(&total).Error; err != nil {
		logger.Log.Error("GetKitchenFavoriteSuppliers count error", "kitchen_id", kitchenID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create pagination meta (using page 1, all items as per_page)
	meta := models.CalculatePaginationMeta(1, len(favorites), total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: favorites,
		Meta: meta,
	})
}

// GetKitchenFavoriteSupplier returns a single favorite supplier by ID
func GetKitchenFavoriteSupplier(c *gin.Context) {
	uid, _ := c.Get("identity")
	kitchenID := c.Param("id")
	favoriteID := c.Param("favoriteId")
	logger.Log.Info("GetKitchenFavoriteSupplier called", "kitchen_id", kitchenID, "favorite_id", favoriteID, "user_id", uid)

	var favorite models.KitchenFavoriteSupplier
	if err := store.DB.GormClient.
		Where("favorite_id = ? AND kitchen_id = ?", favoriteID, kitchenID).
		Preload("Kitchen").
		Preload("Supplier").
		Preload("CreatedBy").
		First(&favorite).Error; err != nil {
		logger.Log.Error("GetKitchenFavoriteSupplier not found", "kitchen_id", kitchenID, "favorite_id", favoriteID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Favorite supplier not found"})
		return
	}

	c.JSON(http.StatusOK, favorite)
}

// CreateKitchenFavoriteSupplier adds a supplier to a kitchen's favorites
func CreateKitchenFavoriteSupplier(c *gin.Context) {
	uid, _ := c.Get("identity")
	kitchenID := c.Param("id")
	logger.Log.Info("CreateKitchenFavoriteSupplier called", "kitchen_id", kitchenID, "user_id", uid)

	// Get user ID from authentication middleware
	var userID string
	if identity, ok := c.Get("identity"); ok {
		if v, ok2 := identity.(string); ok2 {
			userID = v
		}
	}

	var favorite models.KitchenFavoriteSupplier
	if err := c.ShouldBindJSON(&favorite); err != nil {
		logger.Log.Error("CreateKitchenFavoriteSupplier bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set kitchen_id from URL parameter
	favorite.KitchenID = kitchenID
	favorite.CreatedByUserID = userID

	// Validate kitchen exists
	var kitchen models.Kitchen
	if err := store.DB.GormClient.First(&kitchen, "kitchen_id = ?", kitchenID).Error; err != nil {
		logger.Log.Error("CreateKitchenFavoriteSupplier kitchen not found", "kitchen_id", kitchenID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitchen not found"})
		return
	}

	// Validate supplier exists
	var supplier models.Supplier
	if err := store.DB.GormClient.First(&supplier, "supplier_id = ?", favorite.SupplierID).Error; err != nil {
		logger.Log.Error("CreateKitchenFavoriteSupplier supplier not found", "supplier_id", favorite.SupplierID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
		return
	}

	// Check if favorite already exists (unique constraint: kitchen_id + supplier_id)
	var existing models.KitchenFavoriteSupplier
	if err := store.DB.GormClient.Where("kitchen_id = ? AND supplier_id = ?", kitchenID, favorite.SupplierID).First(&existing).Error; err == nil {
		logger.Log.Error("CreateKitchenFavoriteSupplier duplicate favorite", "kitchen_id", kitchenID, "supplier_id", favorite.SupplierID)
		c.JSON(http.StatusConflict, gin.H{"error": "This supplier is already in the kitchen's favorites"})
		return
	}

	// Create favorite
	if err := store.DB.GormClient.Create(&favorite).Error; err != nil {
		logger.Log.Error("CreateKitchenFavoriteSupplier db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload with relations
	store.DB.GormClient.
		Preload("Kitchen").
		Preload("Supplier").
		Preload("CreatedBy").
		First(&favorite, "favorite_id = ?", favorite.FavoriteID)

	c.JSON(http.StatusCreated, favorite)
}
