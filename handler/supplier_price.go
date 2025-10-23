package handler

import (
	"adong-be/models"
	"adong-be/store"
	"adong-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSupplierPrices with pagination and search - Returns ResourceCollection format
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

	if err := db.Find(&prices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: prices,
		Meta: meta,
	})
}

func GetSupplierPrice(c *gin.Context) {
	id := c.Param("id")
	var price models.SupplierPrice
	if err := store.DB.GormClient.First(&price, "sanphamid = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Supplier price not found"})
		return
	}
	c.JSON(http.StatusOK, price)
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
	c.JSON(http.StatusCreated, price)
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
	c.JSON(http.StatusOK, price)
}

func DeleteSupplierPrice(c *gin.Context) {
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.SupplierPrice{}, "sanphamid = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Supplier price deleted successfully"})
}
