package handler

import (
	"adong-be/logger"
	"adong-be/models"
	"adong-be/store"
	"adong-be/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetSupplierPrices(c *gin.Context) {
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get date range parameters
	effectiveFrom := c.Query("effective_from")
	effectiveTo := c.Query("effective_to")
	logger.Log.Debug("receive query", "Effective From:", effectiveFrom, "Effective To:", effectiveTo)

	params = models.GetPaginationParams(
		params.Page,
		params.PageSize,
		params.Search,
		params.SortBy,
		params.SortDir,
	)

	var total int64
	countDB := store.DB.GormClient.Model(&models.SupplierPrice{})

	searchConfig := utils.SearchConfig{
		Fields: []string{"product_name", "ingredient_id", "supplier_id", "classification"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	// Apply date range filters for counting
	countDB = applyDateRangeFilter(countDB, effectiveFrom, effectiveTo)

	fmt.Println(params.Search)
	if err := countDB.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var prices []models.SupplierPrice
	db := store.DB.GormClient.Model(&models.SupplierPrice{})
	db = utils.ApplySearch(db, params.Search, searchConfig)

	// Apply date range filters for data query
	db = applyDateRangeFilter(db, effectiveFrom, effectiveTo)

	allowedSortFields := map[string]string{
		"product_id":     "product_id",
		"product_name":   "product_name",
		"ingredient_id":  "ingredient_id",
		"supplier_id":    "supplier_id",
		"unit_price":     "unit_price",
		"effective_from": "effective_from",
		"effective_to":   "effective_to",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	// Preload related entities to get names
	db = db.Preload("Ingredient").Preload("Supplier")

	if err := db.Find(&prices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to DTOs
	dtos := models.ConvertSupplierPricesToDTO(prices)

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: dtos,
		Meta: meta,
	})
}

// Helper function to apply date range filters
func applyDateRangeFilter(db *gorm.DB, effectiveFrom, effectiveTo string) *gorm.DB {
	// // Parse and validate effectiveFrom date
	// if effectiveFrom != "" {
	// 	// Parse the date string (format: YYYY-MM-DD)
	// 	fromDate, err := time.Parse("2006-01-02", effectiveFrom)
	// 	if err == nil {
	// 		// Filter records where hieuluctu >= effectiveFrom OR hieulucden >= effectiveFrom
	// 		// This ensures we get prices that are effective during or after the from date
	// 		db = db.Where("hieuluctu >= ?", fromDate)
	// 	}
	// }

	// // Parse and validate effectiveTo date
	// if effectiveTo != "" {
	// 	// Parse the date string (format: YYYY-MM-DD)
	// 	toDate, err := time.Parse("2006-01-02", effectiveTo)
	// 	if err == nil {
	// 		// Add 1 day to include the entire end date
	// 		toDateEnd := toDate.Add(24 * time.Hour)
	// 		// Filter records where hieuluctu <= effectiveTo
	// 		// This ensures we get prices that start before or on the to date
	// 		db = db.Where("hieulucden < ?", toDateEnd)
	// 	}
	// }

	if effectiveFrom == "" {
		effectiveFrom = "0001-01-01"
	}
	if effectiveTo == "" {
		effectiveTo = "9999-12-31"
	}

	// If both dates are provided, find prices that overlap with the date range
	if effectiveFrom != "" && effectiveTo != "" {
		fromDate, errFrom := time.Parse("2006-01-02", effectiveFrom)
		toDate, errTo := time.Parse("2006-01-02", effectiveTo)

		if errFrom == nil && errTo == nil {
			// toDateEnd := toDate.Add(24 * time.Hour)
			// Records where:
			// - Start date is within range, OR
			// - End date is within range, OR
			// - The price period encompasses the entire search range
			db = db.Where(
				"hieuluctu >= ? AND hieulucden <= ?",
				fromDate, toDate,
			)
		}
	}

	return db
}

func GetSupplierPrice(c *gin.Context) {
	id := c.Param("id")
	var price models.SupplierPrice

	// Preload related entities to get names
	if err := store.DB.GormClient.
		Preload("Ingredient").
		Preload("Supplier").
		First(&price, "product_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Supplier price not found"})
		return
	}

	// Convert to DTO and return
	dto := price.ToDTO()
	c.JSON(http.StatusOK, dto)
}

// GetSupplierPricesByIngredient - Get all supplier prices for a specific ingredient
func GetSupplierPricesByIngredient(c *gin.Context) {
	ingredientId := c.Param("ingredientId")

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
	countDB := store.DB.GormClient.Model(&models.SupplierPrice{}).Where("ingredient_id = ?", ingredientId)

	searchConfig := utils.SearchConfig{
		Fields: []string{"tensanpham", "nhacungcapid", "phanloai"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var prices []models.SupplierPrice
	db := store.DB.GormClient.Model(&models.SupplierPrice{}).Where("ingredient_id = ?", ingredientId)
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"product_id":     "product_id",
		"product_name":   "product_name",
		"supplier_id":    "supplier_id",
		"unit_price":     "unit_price",
		"effective_from": "effective_from",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	// Preload related entities to get names
	db = db.Preload("Ingredient").Preload("Supplier")

	if err := db.Find(&prices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to DTOs
	dtos := models.ConvertSupplierPricesToDTO(prices)

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: dtos,
		Meta: meta,
	})
}

// GetSupplierPricesBySupplier - Get all supplier prices for a specific supplier
func GetSupplierPricesBySupplier(c *gin.Context) {
	supplierId := c.Param("supplierId")

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
	countDB := store.DB.GormClient.Model(&models.SupplierPrice{}).Where("supplier_id = ?", supplierId)

	searchConfig := utils.SearchConfig{
		Fields: []string{"tensanpham", "nguyenlieuid", "phanloai"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var prices []models.SupplierPrice
	db := store.DB.GormClient.Model(&models.SupplierPrice{}).Where("supplier_id = ?", supplierId)
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"product_id":     "product_id",
		"product_name":   "product_name",
		"ingredient_id":  "ingredient_id",
		"unit_price":     "unit_price",
		"effective_from": "effective_from",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	// Preload related entities to get names
	db = db.Preload("Ingredient").Preload("Supplier")

	if err := db.Find(&prices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to DTOs
	dtos := models.ConvertSupplierPricesToDTO(prices)

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: dtos,
		Meta: meta,
	})
}

func CreateSupplierPrice(c *gin.Context) {
	var price models.SupplierPrice
	if err := c.ShouldBindJSON(&price); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Create(&price).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload with relationships to get names
	store.DB.GormClient.
		Preload("Ingredient").
		Preload("Supplier").
		First(&price, "product_id = ?", price.ProductID)

	// Return DTO with names
	dto := price.ToDTO()
	c.JSON(http.StatusCreated, dto)
}

func UpdateSupplierPrice(c *gin.Context) {
	id := c.Param("id")
	var price models.SupplierPrice
	if err := store.DB.GormClient.First(&price, "product_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Supplier price not found"})
		return
	}
	if err := c.ShouldBindJSON(&price); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Save(&price).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload with relationships to get names
	store.DB.GormClient.
		Preload("Ingredient").
		Preload("Supplier").
		First(&price, "product_id = ?", price.ProductID)

	// Return DTO with names
	dto := price.ToDTO()
	c.JSON(http.StatusOK, dto)
}

func DeleteSupplierPrice(c *gin.Context) {
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.SupplierPrice{}, "product_id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Supplier price deleted successfully"})
}
