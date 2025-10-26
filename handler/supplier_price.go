package handler

import (
	"adong-be/models"
	"adong-be/store"
	"adong-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSupplierPrices with pagination and search - Returns ResourceCollection format with DTOs
func GetSupplierPrices(c *gin.Context) {
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
	countDB := store.DB.GormClient.Model(&models.SupplierPrice{})

	searchConfig := utils.SearchConfig{
		Fields: []string{"tensanpham", "nguyenlieuid", "nhacungcapid", "phanloai"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var prices []models.SupplierPrice
	db := store.DB.GormClient.Model(&models.SupplierPrice{})
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"sanphamid":    "sanphamid",
		"tensanpham":   "tensanpham",
		"nguyenlieuid": "nguyenlieuid",
		"nhacungcapid": "nhacungcapid",
		"dongia":       "dongia",
		"hieuluctu":    "hieuluctu",
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

func GetSupplierPrice(c *gin.Context) {
	id := c.Param("id")
	var price models.SupplierPrice

	// Preload related entities to get names
	if err := store.DB.GormClient.
		Preload("Ingredient").
		Preload("Supplier").
		First(&price, "sanphamid = ?", id).Error; err != nil {
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
	countDB := store.DB.GormClient.Model(&models.SupplierPrice{}).Where("nguyenlieuid = ?", ingredientId)

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
	db := store.DB.GormClient.Model(&models.SupplierPrice{}).Where("nguyenlieuid = ?", ingredientId)
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"sanphamid":    "sanphamid",
		"tensanpham":   "tensanpham",
		"nhacungcapid": "nhacungcapid",
		"dongia":       "dongia",
		"hieuluctu":    "hieuluctu",
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
	countDB := store.DB.GormClient.Model(&models.SupplierPrice{}).Where("nhacungcapid = ?", supplierId)

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
	db := store.DB.GormClient.Model(&models.SupplierPrice{}).Where("nhacungcapid = ?", supplierId)
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"sanphamid":    "sanphamid",
		"tensanpham":   "tensanpham",
		"nguyenlieuid": "nguyenlieuid",
		"dongia":       "dongia",
		"hieuluctu":    "hieuluctu",
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
		First(&price, "sanphamid = ?", price.ProductID)

	// Return DTO with names
	dto := price.ToDTO()
	c.JSON(http.StatusCreated, dto)
}

func UpdateSupplierPrice(c *gin.Context) {
	id := c.Param("id")
	var price models.SupplierPrice
	if err := store.DB.GormClient.First(&price, "sanphamid = ?", id).Error; err != nil {
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
		First(&price, "sanphamid = ?", price.ProductID)

	// Return DTO with names
	dto := price.ToDTO()
	c.JSON(http.StatusOK, dto)
}

func DeleteSupplierPrice(c *gin.Context) {
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.SupplierPrice{}, "sanphamid = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Supplier price deleted successfully"})
}