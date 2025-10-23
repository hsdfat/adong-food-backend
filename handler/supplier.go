package handler

import (
	"adong-be/models"
	"adong-be/store"
	"adong-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSuppliers with pagination and search - Returns ResourceCollection format
func GetSuppliers(c *gin.Context) {
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
	countDB := store.DB.GormClient.Model(&models.Supplier{})

	searchConfig := utils.SearchConfig{
		Fields: []string{"tenncc", "nhacungcapid", "diachi", "dienthoai"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var items []models.Supplier
	db := store.DB.GormClient.Model(&models.Supplier{})
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"nhacungcapid": "nhacungcapid",
		"tenncc":       "tenncc",
		"diachi":       "diachi",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	if err := db.Find(&items).Error; err != nil {
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
	id := c.Param("id")
	var item models.Supplier
	if err := store.DB.GormClient.First(&item, "nhacungcapid = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func CreateSupplier(c *gin.Context) {
	var item models.Supplier
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func UpdateSupplier(c *gin.Context) {
	id := c.Param("id")
	var item models.Supplier
	if err := store.DB.GormClient.First(&item, "nhacungcapid = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
		return
	}
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func DeleteSupplier(c *gin.Context) {
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.Supplier{}, "nhacungcapid = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Supplier deleted successfully"})
}
