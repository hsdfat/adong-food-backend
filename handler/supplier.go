package handler

import (
	"adong-be/logger"
	"adong-be/models"
	"adong-be/store"
	"adong-be/utils"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetSuppliers with pagination and search - Returns ResourceCollection format
func GetSuppliers(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("GetSuppliers called", "user_id", uid)
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Log.Error("GetSuppliers bind query error", "error", err)
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
	countDB := store.DB.GormClient.Model(&models.Supplier{})

	searchConfig := utils.SearchConfig{
		Fields: []string{"supplier_name", "supplier_id", "address", "phone"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		logger.Log.Error("GetSuppliers count error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var items []models.Supplier
	db := store.DB.GormClient.Model(&models.Supplier{})
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"supplier_id":   "supplier_id",
		"supplier_name": "supplier_name",
		"address":       "address",
		"created_date":  "created_date",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	if err := db.Find(&items).Error; err != nil {
		logger.Log.Error("GetSuppliers query error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: items,
		Meta: meta,
	})
}

func GetSupplier(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("GetSupplier called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	var item models.Supplier
	if err := store.DB.GormClient.First(&item, "supplier_id = ?", id).Error; err != nil {
		logger.Log.Error("GetSupplier not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func CreateSupplier(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("CreateSupplier called", "user_id", uid)
	var item models.Supplier
	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Log.Error("CreateSupplier bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Create(&item).Error; err != nil {
		logger.Log.Error("CreateSupplier db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func UpdateSupplier(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("UpdateSupplier called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	var item models.Supplier
	if err := store.DB.GormClient.First(&item, "supplier_id = ?", id).Error; err != nil {
		logger.Log.Error("UpdateSupplier not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
		return
	}
	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Log.Error("UpdateSupplier bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Save(&item).Error; err != nil {
		logger.Log.Error("UpdateSupplier db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func DeleteSupplier(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("DeleteSupplier called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.Supplier{}, "supplier_id = ?", id).Error; err != nil {
		logger.Log.Error("DeleteSupplier db error", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Supplier deleted successfully"})
}

// FindBestSuppliers - Find best suppliers for ingredients based on kitchen preferences and pricing
func FindBestSuppliers(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("FindBestSuppliers called", "user_id", uid)

	var req models.BestSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("FindBestSuppliers bind JSON error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate order exists and belongs to the specified kitchen
	var order models.Order
	if err := store.DB.GormClient.Where("order_id = ? AND kitchen_id = ?", req.OrderID, req.KitchenID).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found or does not belong to specified kitchen"})
			return
		}
		logger.Log.Error("FindBestSuppliers order validation error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get ingredients with their types and material groups
	var ingredients []models.Ingredient
	if err := store.DB.GormClient.Preload("IngredientType").Where("ingredient_id IN ?", req.IngredientIDs).Find(&ingredients).Error; err != nil {
		logger.Log.Error("FindBestSuppliers ingredients query error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(ingredients) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No ingredients found"})
		return
	}

	// Get kitchen favorite suppliers
	var favoriteSuppliers []models.KitchenFavoriteSupplier
	if err := store.DB.GormClient.Where("kitchen_id = ?", req.KitchenID).Find(&favoriteSuppliers).Error; err != nil {
		logger.Log.Error("FindBestSuppliers favorite suppliers query error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create map of favorite supplier IDs for quick lookup
	favoriteSupplierMap := make(map[string]bool)
	for _, fav := range favoriteSuppliers {
		favoriteSupplierMap[fav.SupplierID] = true
	}

	// Process each ingredient to find best supplier
	var result []models.IngredientSupplierInfo
	for _, ingredient := range ingredients {
		supplierInfo := findBestSupplierForIngredient(ingredient, favoriteSupplierMap)
		result = append(result, supplierInfo)
	}

	response := models.BestSupplierResponse{
		OrderID:   req.OrderID,
		KitchenID: req.KitchenID,
		Suppliers: result,
	}

	c.JSON(http.StatusOK, response)
}

// findBestSupplierForIngredient - Helper function to find best supplier for a single ingredient
func findBestSupplierForIngredient(ingredient models.Ingredient, favoriteSupplierMap map[string]bool) models.IngredientSupplierInfo {
	ingredientTypeName := ""
	if ingredient.IngredientType != nil {
		ingredientTypeName = ingredient.IngredientType.IngredientTypeName
	}

	// Determine selection strategy based on ingredient type
	useFavoriteStrategy := shouldUseFavoriteStrategy(ingredientTypeName, ingredient.MaterialGroup)

	var supplierPrices []models.SupplierPrice
	query := store.DB.GormClient.Preload("Supplier").Where("ingredient_id = ? AND active = ?", ingredient.IngredientID, true)

	if useFavoriteStrategy {
		// For favorite strategy, prioritize kitchen's favorite suppliers
		var favoriteIDs []string
		for id := range favoriteSupplierMap {
			favoriteIDs = append(favoriteIDs, id)
		}
		if len(favoriteIDs) > 0 {
			query = query.Where("supplier_id IN ?", favoriteIDs)
		}
	}

	query.Find(&supplierPrices)

	if len(supplierPrices) == 0 {
		// If no suppliers found with favorite strategy, try all suppliers
		if useFavoriteStrategy {
			store.DB.GormClient.Preload("Supplier").Where("ingredient_id = ? AND active = ?", ingredient.IngredientID, true).Find(&supplierPrices)
		}
	}

	var selectedSupplier *models.SupplierInfo
	var selectionReason string

	if len(supplierPrices) > 0 {
		if useFavoriteStrategy {
			// For favorite strategy, sort by display order (if available) then by price
			sort.Slice(supplierPrices, func(i, j int) bool {
				// Prioritize favorite suppliers
				isFavI := favoriteSupplierMap[supplierPrices[i].SupplierID]
				isFavJ := favoriteSupplierMap[supplierPrices[j].SupplierID]
				if isFavI != isFavJ {
					return isFavI
				}
				// If both are favorites or both are not, sort by price
				return supplierPrices[i].UnitPrice < supplierPrices[j].UnitPrice
			})
			selectionReason = "Kitchen favorite supplier (lowest price among favorites)"
		} else {
			// For price strategy, sort by lowest price
			sort.Slice(supplierPrices, func(i, j int) bool {
				return supplierPrices[i].UnitPrice < supplierPrices[j].UnitPrice
			})
			selectionReason = "Lowest price supplier"
		}

		bestPrice := supplierPrices[0]
		selectedSupplier = &models.SupplierInfo{
			SupplierID:   bestPrice.SupplierID,
			SupplierName: bestPrice.Supplier.SupplierName,
			Phone:        bestPrice.Supplier.Phone,
			Email:        bestPrice.Supplier.Email,
			Address:      bestPrice.Supplier.Address,
			UnitPrice:    bestPrice.UnitPrice,
			Unit:         bestPrice.Unit,
			ProductName:  bestPrice.ProductName,
			ProductID:    bestPrice.ProductID,
		}
	}

	return models.IngredientSupplierInfo{
		IngredientID:     ingredient.IngredientID,
		IngredientName:   ingredient.IngredientName,
		IngredientType:   ingredientTypeName,
		MaterialGroup:    ingredient.MaterialGroup,
		SelectedSupplier: selectedSupplier,
		SelectionReason:  selectionReason,
	}
}

// shouldUseFavoriteStrategy - Determine if ingredient should use favorite supplier strategy
func shouldUseFavoriteStrategy(ingredientType, materialGroup string) bool {
	// Use favorite strategy for vegetables, meat, beans, eggs
	favoriteTypes := map[string]bool{
		"VEGETABLE": true,
		"MEAT":      true,
		"DAIRY":     true, // includes eggs
		"GRAIN":     true, // includes beans
	}

	// Check ingredient type first
	if favoriteTypes[ingredientType] {
		return true
	}

	// Check material group for specific cases
	favoriteGroups := map[string]bool{
		"Thịt heo":     true,
		"Thịt bò":      true,
		"Thịt gia cầm": true,
		"Trứng":        true,
		"Gạo":          true,
		"Bún phở":      true,
		"Củ quả":       true,
		"Rau xanh":     true,
		"Củ":           true,
	}

	return favoriteGroups[materialGroup]
}
