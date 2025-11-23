package handler

import (
	"adong-be/models"
	"adong-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InventoryStockHandler struct {
	DB *gorm.DB
}

func NewInventoryStockHandler(db *gorm.DB) *InventoryStockHandler {
	return &InventoryStockHandler{DB: db}
}

// GetAllStocks retrieves all inventory stocks with pagination and filters
func (h *InventoryStockHandler) GetAllStocks(c *gin.Context) {
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

	kitchenID := c.Query("kitchen_id")
	lowStock := c.Query("low_stock") // "true" to filter items below min stock

	var stocks []models.InventoryStock
	var total int64

	countQuery := h.DB.Model(&models.InventoryStock{})

	if kitchenID != "" {
		countQuery = countQuery.Where("kitchen_id = ?", kitchenID)
	}

	if params.Search != "" {
		countQuery = countQuery.Joins("JOIN master_ingredients ON master_ingredients.ingredient_id = inventory_stocks.ingredient_id").
			Where("master_ingredients.ingredient_name ILIKE ?", "%"+params.Search+"%")
	}

	if lowStock == "true" {
		countQuery = countQuery.Where("min_stock_level IS NOT NULL AND quantity < min_stock_level")
	}

	if err := countQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi đếm tồn kho"})
		return
	}

	query := h.DB.Model(&models.InventoryStock{})

	if kitchenID != "" {
		query = query.Where("kitchen_id = ?", kitchenID)
	}

	if params.Search != "" {
		query = query.Joins("JOIN master_ingredients ON master_ingredients.ingredient_id = inventory_stocks.ingredient_id").
			Where("master_ingredients.ingredient_name ILIKE ?", "%"+params.Search+"%")
	}

	if lowStock == "true" {
		query = query.Where("min_stock_level IS NOT NULL AND quantity < min_stock_level")
	}

	allowedSortFields := map[string]string{
		"last_updated":  "last_updated",
		"quantity":      "quantity",
		"ingredient_id": "ingredient_id",
	}
	query = utils.ApplySort(query, params.SortBy, params.SortDir, allowedSortFields)
	if params.SortBy == "" {
		query = query.Order("last_updated DESC")
	}
	query = utils.ApplyPagination(query, params.Page, params.PageSize)

	if err := query.Preload("Kitchen").
		Preload("Ingredient").
		Find(&stocks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy danh sách tồn kho"})
		return
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: stocks,
		Meta: meta,
	})
}

// GetStockByID retrieves a specific inventory stock
func (h *InventoryStockHandler) GetStockByID(c *gin.Context) {
	stockID := c.Param("id")

	var stock models.InventoryStock
	if err := h.DB.Preload("Kitchen").
		Preload("Ingredient").
		Where("stock_id = ?", stockID).
		First(&stock).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy tồn kho"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy thông tin tồn kho"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stock})
}

// GetStockByKitchenAndIngredient retrieves stock for a specific kitchen and ingredient
func (h *InventoryStockHandler) GetStockByKitchenAndIngredient(c *gin.Context) {
	kitchenID := c.Query("kitchen_id")
	ingredientID := c.Query("ingredient_id")

	if kitchenID == "" || ingredientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cần có kitchen_id và ingredient_id"})
		return
	}

	var stock models.InventoryStock
	if err := h.DB.Preload("Kitchen").
		Preload("Ingredient").
		Where("kitchen_id = ? AND ingredient_id = ?", kitchenID, ingredientID).
		First(&stock).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy tồn kho"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy thông tin tồn kho"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stock})
}

// UpdateStockLevels updates min/max stock levels
type UpdateStockLevelsRequest struct {
	MinStockLevel *float64 `json:"minStockLevel"`
	MaxStockLevel *float64 `json:"maxStockLevel"`
}

func (h *InventoryStockHandler) UpdateStockLevels(c *gin.Context) {
	stockID := c.Param("id")

	var req UpdateStockLevelsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var stock models.InventoryStock
	if err := h.DB.Where("stock_id = ?", stockID).First(&stock).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy tồn kho"})
		return
	}

	updates := map[string]interface{}{
		"min_stock_level": req.MinStockLevel,
		"max_stock_level": req.MaxStockLevel,
	}

	if err := h.DB.Model(&stock).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật mức tồn"})
		return
	}

	h.DB.Preload("Kitchen").Preload("Ingredient").First(&stock, "stock_id = ?", stockID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật mức tồn thành công",
		"data":    stock,
	})
}

// GetLowStockAlerts retrieves all items below minimum stock level
func (h *InventoryStockHandler) GetLowStockAlerts(c *gin.Context) {
	kitchenID := c.Query("kitchen_id")

	var stocks []models.InventoryStock
	query := h.DB.Model(&models.InventoryStock{}).
		Where("min_stock_level IS NOT NULL AND quantity < min_stock_level")

	if kitchenID != "" {
		query = query.Where("kitchen_id = ?", kitchenID)
	}

	if err := query.Preload("Kitchen").
		Preload("Ingredient").
		Order("(quantity / NULLIF(min_stock_level, 0))").
		Find(&stocks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy cảnh báo tồn kho"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  stocks,
		"count": len(stocks),
	})
}

// GetStockTransactions retrieves transaction history for a stock
func (h *InventoryStockHandler) GetStockTransactions(c *gin.Context) {
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

	kitchenID := c.Query("kitchen_id")
	ingredientID := c.Query("ingredient_id")
	transactionType := c.Query("transaction_type")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")

	if kitchenID == "" || ingredientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cần có kitchen_id và ingredient_id"})
		return
	}

	var transactions []models.InventoryTransaction
	var total int64

	countQuery := h.DB.Model(&models.InventoryTransaction{}).
		Where("kitchen_id = ? AND ingredient_id = ?", kitchenID, ingredientID)

	if transactionType != "" {
		countQuery = countQuery.Where("transaction_type = ?", transactionType)
	}
	if fromDate != "" {
		countQuery = countQuery.Where("transaction_date >= ?", fromDate)
	}
	if toDate != "" {
		countQuery = countQuery.Where("transaction_date <= ?", toDate)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi đếm giao dịch"})
		return
	}

	query := h.DB.Model(&models.InventoryTransaction{}).
		Where("kitchen_id = ? AND ingredient_id = ?", kitchenID, ingredientID)

	if transactionType != "" {
		query = query.Where("transaction_type = ?", transactionType)
	}
	if fromDate != "" {
		query = query.Where("transaction_date >= ?", fromDate)
	}
	if toDate != "" {
		query = query.Where("transaction_date <= ?", toDate)
	}

	allowedSortFields := map[string]string{
		"transaction_date": "transaction_date",
		"transaction_type": "transaction_type",
		"quantity":         "quantity",
	}
	query = utils.ApplySort(query, params.SortBy, params.SortDir, allowedSortFields)
	if params.SortBy == "" {
		query = query.Order("transaction_date DESC")
	}
	query = utils.ApplyPagination(query, params.Page, params.PageSize)

	if err := query.Preload("Kitchen").
		Preload("Ingredient").
		Preload("CreatedBy").
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy lịch sử giao dịch"})
		return
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: transactions,
		Meta: meta,
	})
}

// GetStockSummary retrieves summary statistics for inventory
func (h *InventoryStockHandler) GetStockSummary(c *gin.Context) {
	kitchenID := c.Query("kitchen_id")

	if kitchenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cần có kitchen_id"})
		return
	}

	type Summary struct {
		TotalItems      int64   `json:"totalItems"`
		LowStockItems   int64   `json:"lowStockItems"`
		OutOfStockItems int64   `json:"outOfStockItems"`
		TotalValue      float64 `json:"totalValue"`
	}

	var summary Summary

	// Total items
	h.DB.Model(&models.InventoryStock{}).
		Where("kitchen_id = ?", kitchenID).
		Count(&summary.TotalItems)

	// Low stock items
	h.DB.Model(&models.InventoryStock{}).
		Where("kitchen_id = ? AND min_stock_level IS NOT NULL AND quantity < min_stock_level", kitchenID).
		Count(&summary.LowStockItems)

	// Out of stock items
	h.DB.Model(&models.InventoryStock{}).
		Where("kitchen_id = ? AND quantity = 0", kitchenID).
		Count(&summary.OutOfStockItems)

	// Total value - would need to join with latest prices
	// This is a simplified version
	h.DB.Model(&models.InventoryStock{}).
		Select("COALESCE(SUM(quantity), 0)").
		Where("kitchen_id = ?", kitchenID).
		Scan(&summary.TotalValue)

	c.JSON(http.StatusOK, gin.H{"data": summary})
}

// GetStockValuation retrieves stock valuation based on latest prices
func (h *InventoryStockHandler) GetStockValuation(c *gin.Context) {
	kitchenID := c.Query("kitchen_id")

	if kitchenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cần có kitchen_id"})
		return
	}

	type ValuationItem struct {
		IngredientID   string  `json:"ingredientId"`
		IngredientName string  `json:"ingredientName"`
		Quantity       float64 `json:"quantity"`
		Unit           string  `json:"unit"`
		AveragePrice   float64 `json:"averagePrice"`
		TotalValue     float64 `json:"totalValue"`
	}

	var items []ValuationItem

	query := `
		SELECT 
			s.ingredient_id,
			i.ingredient_name,
			s.quantity,
			s.unit,
			COALESCE(AVG(p.unit_price), 0) as average_price,
			s.quantity * COALESCE(AVG(p.unit_price), 0) as total_value
		FROM inventory_stocks s
		JOIN master_ingredients i ON i.ingredient_id = s.ingredient_id
		LEFT JOIN supplier_price_list p ON p.ingredient_id = s.ingredient_id AND p.active = true
		WHERE s.kitchen_id = ?
		GROUP BY s.ingredient_id, i.ingredient_name, s.quantity, s.unit
		ORDER BY total_value DESC
	`

	if err := h.DB.Raw(query, kitchenID).Scan(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tính giá trị tồn kho"})
		return
	}

	var totalValue float64
	for _, item := range items {
		totalValue += item.TotalValue
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       items,
		"totalValue": totalValue,
		"count":      len(items),
	})
}
