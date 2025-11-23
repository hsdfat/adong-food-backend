package handler

import (
	"adong-be/models"
	"adong-be/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InventoryExportHandler struct {
	DB *gorm.DB
}

func NewInventoryExportHandler(db *gorm.DB) *InventoryExportHandler {
	return &InventoryExportHandler{DB: db}
}

// CreateExportRequest represents the request body for creating an export
type CreateExportRequest struct {
	KitchenID            string                      `json:"kitchenId" binding:"required"`
	ExportDate           string                      `json:"exportDate" binding:"required"`
	ExportType           string                      `json:"exportType" binding:"required"`
	DestinationKitchenID *string                     `json:"destinationKitchenId"`
	OrderID              *string                     `json:"orderId"`
	Status               string                      `json:"status"`
	Notes                *string                     `json:"notes"`
	ExportDetails        []CreateExportDetailRequest `json:"exportDetails" binding:"required,min=1"`
}

type CreateExportDetailRequest struct {
	IngredientID string   `json:"ingredientId" binding:"required"`
	Quantity     float64  `json:"quantity" binding:"required,gt=0"`
	Unit         string   `json:"unit" binding:"required"`
	UnitCost     *float64 `json:"unitCost"`
	BatchNumber  *string  `json:"batchNumber"`
	Notes        *string  `json:"notes"`
}

// GetAllExports retrieves all inventory exports with pagination and filters
func (h *InventoryExportHandler) GetAllExports(c *gin.Context) {
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
	exportType := c.Query("export_type")
	status := c.Query("status")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")

	var exports []models.InventoryExport
	var total int64

	countQuery := h.DB.Model(&models.InventoryExport{})

	if kitchenID != "" {
		countQuery = countQuery.Where("kitchen_id = ?", kitchenID)
	}
	if exportType != "" {
		countQuery = countQuery.Where("export_type = ?", exportType)
	}
	if status != "" {
		countQuery = countQuery.Where("status = ?", status)
	}
	if fromDate != "" {
		countQuery = countQuery.Where("export_date >= ?", fromDate)
	}
	if toDate != "" {
		countQuery = countQuery.Where("export_date <= ?", toDate)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi đếm phiếu xuất"})
		return
	}

	query := h.DB.Model(&models.InventoryExport{})

	if kitchenID != "" {
		query = query.Where("kitchen_id = ?", kitchenID)
	}
	if exportType != "" {
		query = query.Where("export_type = ?", exportType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if fromDate != "" {
		query = query.Where("export_date >= ?", fromDate)
	}
	if toDate != "" {
		query = query.Where("export_date <= ?", toDate)
	}

	allowedSortFields := map[string]string{
		"export_date":  "export_date",
		"created_date": "created_date",
		"status":       "status",
		"export_type":  "export_type",
		"total_amount": "total_amount",
	}
	query = utils.ApplySort(query, params.SortBy, params.SortDir, allowedSortFields)
	if params.SortBy == "" {
		query = query.Order("export_date DESC, created_date DESC")
	}
	query = utils.ApplyPagination(query, params.Page, params.PageSize)

	if err := query.Preload("Kitchen").
		Preload("DestinationKitchen").
		Preload("IssuedBy").
		Preload("ApprovedBy").
		Preload("CreatedBy").
		Find(&exports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy danh sách phiếu xuất"})
		return
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: exports,
		Meta: meta,
	})
}

// GetExportByID retrieves a specific export with details
func (h *InventoryExportHandler) GetExportByID(c *gin.Context) {
	exportID := c.Param("id")

	var exportRecord models.InventoryExport
	if err := h.DB.Preload("Kitchen").
		Preload("DestinationKitchen").
		Preload("Order").
		Preload("IssuedBy").
		Preload("ApprovedBy").
		Preload("CreatedBy").
		Preload("ExportDetails.Ingredient").
		Where("export_id = ?", exportID).
		First(&exportRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy phiếu xuất"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy thông tin phiếu xuất"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": exportRecord})
}

// CreateExport creates a new inventory export with details
func (h *InventoryExportHandler) CreateExport(c *gin.Context) {
	var req CreateExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userID string
	if identity, ok := c.Get("identity"); ok {
		if v, ok2 := identity.(string); ok2 {
			userID = v
		}
	}

	exportDate, err := time.Parse("2006-01-02", req.ExportDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Định dạng ngày không hợp lệ"})
		return
	}

	// Validate export type
	validTypes := map[string]bool{
		"production": true, // Xuất cho sản xuất
		"transfer":   true, // Chuyển kho
		"disposal":   true, // Hủy bỏ
		"return":     true, // Trả hàng
		"sample":     true, // Xuất mẫu
	}
	if !validTypes[req.ExportType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Loại xuất kho không hợp lệ"})
		return
	}

	// Generate export ID
	exportID := generateExportID(exportDate, req.ExportType)

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	status := req.Status
	if status == "" {
		status = "draft"
	}

	// Create export header
	exportRecord := models.InventoryExport{
		ExportID:             exportID,
		KitchenID:            req.KitchenID,
		ExportDate:           exportDate,
		ExportType:           req.ExportType,
		DestinationKitchenID: req.DestinationKitchenID,
		OrderID:              req.OrderID,
		Status:               status,
		Notes:                req.Notes,
		CreatedByUserID:      &userID,
	}

	if err := tx.Create(&exportRecord).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo phiếu xuất"})
		return
	}

	// Create export details and calculate total
	var totalAmount float64
	for _, detail := range req.ExportDetails {
		var totalCost *float64
		if detail.UnitCost != nil {
			cost := *detail.UnitCost * detail.Quantity
			totalCost = &cost
			totalAmount += cost
		}

		exportDetail := models.InventoryExportDetail{
			ExportID:     exportID,
			IngredientID: detail.IngredientID,
			Quantity:     detail.Quantity,
			Unit:         detail.Unit,
			UnitCost:     detail.UnitCost,
			TotalCost:    totalCost,
			BatchNumber:  detail.BatchNumber,
			Notes:        detail.Notes,
		}

		if err := tx.Create(&exportDetail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo chi tiết phiếu xuất"})
			return
		}
	}

	// Update total amount
	if err := tx.Model(&exportRecord).Update("total_amount", totalAmount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật tổng tiền"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lưu phiếu xuất"})
		return
	}

	h.DB.Preload("Kitchen").
		Preload("DestinationKitchen").
		Preload("ExportDetails.Ingredient").
		First(&exportRecord, "export_id = ?", exportID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo phiếu xuất thành công",
		"data":    exportRecord,
	})
}

// UpdateExport updates an existing export
func (h *InventoryExportHandler) UpdateExport(c *gin.Context) {
	exportID := c.Param("id")

	var req CreateExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingExport models.InventoryExport
	if err := h.DB.Where("export_id = ?", exportID).First(&existingExport).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy phiếu xuất"})
		return
	}

	if existingExport.Status == "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể sửa phiếu xuất đã duyệt"})
		return
	}

	exportDate, err := time.Parse("2006-01-02", req.ExportDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Định dạng ngày không hợp lệ"})
		return
	}

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete existing details
	if err := tx.Where("export_id = ?", exportID).Delete(&models.InventoryExportDetail{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi xóa chi tiết cũ"})
		return
	}

	// Update header
	updates := map[string]interface{}{
		"kitchen_id":             req.KitchenID,
		"export_date":            exportDate,
		"export_type":            req.ExportType,
		"destination_kitchen_id": req.DestinationKitchenID,
		"order_id":               req.OrderID,
		"notes":                  req.Notes,
	}

	if err := tx.Model(&existingExport).Updates(updates).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật phiếu xuất"})
		return
	}

	// Create new details
	var totalAmount float64
	for _, detail := range req.ExportDetails {
		var totalCost *float64
		if detail.UnitCost != nil {
			cost := *detail.UnitCost * detail.Quantity
			totalCost = &cost
			totalAmount += cost
		}

		exportDetail := models.InventoryExportDetail{
			ExportID:     exportID,
			IngredientID: detail.IngredientID,
			Quantity:     detail.Quantity,
			Unit:         detail.Unit,
			UnitCost:     detail.UnitCost,
			TotalCost:    totalCost,
			BatchNumber:  detail.BatchNumber,
			Notes:        detail.Notes,
		}

		if err := tx.Create(&exportDetail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo chi tiết phiếu xuất"})
			return
		}
	}

	// Update total amount
	if err := tx.Model(&existingExport).Update("total_amount", totalAmount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật tổng tiền"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lưu phiếu xuất"})
		return
	}

	h.DB.Preload("Kitchen").
		Preload("DestinationKitchen").
		Preload("ExportDetails.Ingredient").
		First(&existingExport, "export_id = ?", exportID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật phiếu xuất thành công",
		"data":    existingExport,
	})
}

// ApproveExport approves an export and updates inventory stocks
func (h *InventoryExportHandler) ApproveExport(c *gin.Context) {
	exportID := c.Param("id")
	var userID string
	if identity, ok := c.Get("identity"); ok {
		if v, ok2 := identity.(string); ok2 {
			userID = v
		}
	}

	var exportRecord models.InventoryExport
	if err := h.DB.Preload("ExportDetails").
		Where("export_id = ?", exportID).
		First(&exportRecord).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy phiếu xuất"})
		return
	}

	if exportRecord.Status == "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phiếu xuất đã được duyệt"})
		return
	}

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check stock availability
	for _, detail := range exportRecord.ExportDetails {
		var stock models.InventoryStock
		if err := tx.Where("kitchen_id = ? AND ingredient_id = ?",
			exportRecord.KitchenID, detail.IngredientID).
			First(&stock).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error":         "Nguyên liệu không tồn tại trong kho",
				"ingredient_id": detail.IngredientID,
			})
			return
		}

		if stock.Quantity < detail.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error":         "Số lượng tồn kho không đủ",
				"ingredient_id": detail.IngredientID,
				"available":     stock.Quantity,
				"required":      detail.Quantity,
			})
			return
		}
	}

	// Update export status
	now := time.Now()
	updates := map[string]interface{}{
		"status":              "approved",
		"approved_by_user_id": userID,
		"approved_date":       now,
	}

	if err := tx.Model(&exportRecord).Updates(updates).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi duyệt phiếu xuất"})
		return
	}

	// Update inventory stocks
	for _, detail := range exportRecord.ExportDetails {
		var stock models.InventoryStock
		tx.Where("kitchen_id = ? AND ingredient_id = ?",
			exportRecord.KitchenID, detail.IngredientID).
			First(&stock)

		quantityBefore := stock.Quantity

		// Decrease stock
		if err := tx.Model(&stock).Updates(map[string]interface{}{
			"quantity":     gorm.Expr("quantity - ?", detail.Quantity),
			"last_updated": now,
		}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật tồn kho"})
			return
		}

		// Log transaction
		transaction := models.InventoryTransaction{
			KitchenID:       exportRecord.KitchenID,
			IngredientID:    detail.IngredientID,
			TransactionType: "EXPORT",
			TransactionDate: now,
			Quantity:        -detail.Quantity, // Negative for exports
			Unit:            detail.Unit,
			QuantityBefore:  quantityBefore,
			QuantityAfter:   quantityBefore - detail.Quantity,
			ReferenceType:   strPtr(exportRecord.ExportType),
			ReferenceID:     &exportID,
			CreatedByUserID: &userID,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi ghi log giao dịch"})
			return
		}

		// If transfer to another kitchen, create import record there
		if exportRecord.ExportType == "transfer" && exportRecord.DestinationKitchenID != nil {
			var destStock models.InventoryStock
			destResult := tx.Where("kitchen_id = ? AND ingredient_id = ?",
				*exportRecord.DestinationKitchenID, detail.IngredientID).
				First(&destStock)

			destQuantityBefore := 0.0
			if destResult.Error == nil {
				destQuantityBefore = destStock.Quantity
				// Update destination stock
				if err := tx.Model(&destStock).Updates(map[string]interface{}{
					"quantity":     gorm.Expr("quantity + ?", detail.Quantity),
					"unit":         detail.Unit,
					"last_updated": now,
				}).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật kho đích"})
					return
				}
			} else {
				// Create new stock entry at destination
				destStock = models.InventoryStock{
					KitchenID:    *exportRecord.DestinationKitchenID,
					IngredientID: detail.IngredientID,
					Quantity:     detail.Quantity,
					Unit:         detail.Unit,
					LastUpdated:  now,
				}
				if err := tx.Create(&destStock).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo tồn kho đích"})
					return
				}
			}

			// Log destination transaction
			destTransaction := models.InventoryTransaction{
				KitchenID:       *exportRecord.DestinationKitchenID,
				IngredientID:    detail.IngredientID,
				TransactionType: "TRANSFER_IN",
				TransactionDate: now,
				Quantity:        detail.Quantity,
				Unit:            detail.Unit,
				QuantityBefore:  destQuantityBefore,
				QuantityAfter:   destQuantityBefore + detail.Quantity,
				ReferenceType:   strPtr("EXPORT"),
				ReferenceID:     &exportID,
				CreatedByUserID: &userID,
			}
			if err := tx.Create(&destTransaction).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi ghi log chuyển kho"})
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi hoàn tất duyệt phiếu"})
		return
	}

	h.DB.Preload("Kitchen").
		Preload("DestinationKitchen").
		Preload("ApprovedBy").
		Preload("ExportDetails.Ingredient").
		First(&exportRecord, "export_id = ?", exportID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Duyệt phiếu xuất thành công",
		"data":    exportRecord,
	})
}

// DeleteExport deletes a draft export
func (h *InventoryExportHandler) DeleteExport(c *gin.Context) {
	exportID := c.Param("id")

	var exportRecord models.InventoryExport
	if err := h.DB.Where("export_id = ?", exportID).First(&exportRecord).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy phiếu xuất"})
		return
	}

	if exportRecord.Status == "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể xóa phiếu xuất đã duyệt"})
		return
	}

	if err := h.DB.Delete(&exportRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi xóa phiếu xuất"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xóa phiếu xuất thành công"})
}

// Helper function
func generateExportID(exportDate time.Time, exportType string) string {
	prefix := "EX"
	if exportType == "transfer" {
		prefix = "TR"
	} else if exportType == "disposal" {
		prefix = "DS"
	}
	return prefix + exportDate.Format("20060102") + "-" + strconv.FormatInt(time.Now().UnixNano()%100000, 10)
}
