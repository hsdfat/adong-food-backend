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

type InventoryAdjustmentHandler struct {
	DB *gorm.DB
}

func NewInventoryAdjustmentHandler(db *gorm.DB) *InventoryAdjustmentHandler {
	return &InventoryAdjustmentHandler{DB: db}
}

// CreateAdjustmentRequest represents the request body for creating an adjustment
type CreateAdjustmentRequest struct {
	KitchenID        string                            `json:"kitchenId" binding:"required"`
	AdjustmentDate   string                            `json:"adjustmentDate" binding:"required"`
	AdjustmentType   string                            `json:"adjustmentType" binding:"required"`
	Reason           *string                           `json:"reason"`
	Status           string                            `json:"status"`
	AdjustmentDetails []CreateAdjustmentDetailRequest `json:"adjustmentDetails" binding:"required,min=1"`
}

type CreateAdjustmentDetailRequest struct {
	IngredientID       string   `json:"ingredientId" binding:"required"`
	QuantityBefore     float64  `json:"quantityBefore" binding:"required"`
	QuantityAfter      float64  `json:"quantityAfter" binding:"required"`
	Unit               string   `json:"unit" binding:"required"`
	UnitCost           *float64 `json:"unitCost"`
	Reason             *string  `json:"reason"`
}

// GetAllAdjustments retrieves all inventory adjustments with pagination and filters
func (h *InventoryAdjustmentHandler) GetAllAdjustments(c *gin.Context) {
	// Kitchen-based authorization
	scope, err := utils.GetUserKitchenScope(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
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
	adjustmentType := c.Query("adjustment_type")
	status := c.Query("status")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")

	var adjustments []models.InventoryAdjustment
	var total int64

	countQuery := h.DB.Model(&models.InventoryAdjustment{})

	// Apply kitchen auth for count query
	if scope.IsAdmin {
		if kitchenID != "" {
			countQuery = countQuery.Where("kitchen_id = ?", kitchenID)
		}
	} else {
		if len(scope.KitchenIDs) == 0 {
			meta := models.CalculatePaginationMeta(params.Page, params.PageSize, 0)
			c.JSON(http.StatusOK, models.ResourceCollection{
				Data: []models.InventoryAdjustment{},
				Meta: meta,
			})
			return
		}
		if kitchenID != "" {
			allowed := false
			for _, kid := range scope.KitchenIDs {
				if kid == kitchenID {
					allowed = true
					break
				}
			}
			if !allowed {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access to this kitchen is not allowed"})
				return
			}
			countQuery = countQuery.Where("kitchen_id = ?", kitchenID)
		} else {
			countQuery = countQuery.Where("kitchen_id IN ?", scope.KitchenIDs)
		}
	}
	if adjustmentType != "" {
		countQuery = countQuery.Where("adjustment_type = ?", adjustmentType)
	}
	if status != "" {
		countQuery = countQuery.Where("status = ?", status)
	}
	if fromDate != "" {
		countQuery = countQuery.Where("adjustment_date >= ?", fromDate)
	}
	if toDate != "" {
		countQuery = countQuery.Where("adjustment_date <= ?", toDate)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi đếm phiếu kiểm kê"})
		return
	}

	query := h.DB.Model(&models.InventoryAdjustment{})

	// Apply same kitchen restriction to data query
	if scope.IsAdmin {
		if kitchenID != "" {
			query = query.Where("kitchen_id = ?", kitchenID)
		}
	} else {
		if kitchenID != "" {
			query = query.Where("kitchen_id = ?", kitchenID)
		} else {
			query = query.Where("kitchen_id IN ?", scope.KitchenIDs)
		}
	}
	if adjustmentType != "" {
		query = query.Where("adjustment_type = ?", adjustmentType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if fromDate != "" {
		query = query.Where("adjustment_date >= ?", fromDate)
	}
	if toDate != "" {
		query = query.Where("adjustment_date <= ?", toDate)
	}

	allowedSortFields := map[string]string{
		"adjustment_date": "adjustment_date",
		"created_date":    "created_date",
		"status":          "status",
		"adjustment_type": "adjustment_type",
	}
	query = utils.ApplySort(query, params.SortBy, params.SortDir, allowedSortFields)
	if params.SortBy == "" {
		query = query.Order("adjustment_date DESC, created_date DESC")
	}
	query = utils.ApplyPagination(query, params.Page, params.PageSize)

	if err := query.Preload("Kitchen").
		Preload("ApprovedBy").
		Preload("CreatedBy").
		Find(&adjustments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy danh sách phiếu kiểm kê"})
		return
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: adjustments,
		Meta: meta,
	})
}

// GetAdjustmentByID retrieves a specific adjustment with details
func (h *InventoryAdjustmentHandler) GetAdjustmentByID(c *gin.Context) {
	adjustmentID := c.Param("id")

	var adjustment models.InventoryAdjustment
	if err := h.DB.Preload("Kitchen").
		Preload("ApprovedBy").
		Preload("CreatedBy").
		Preload("AdjustmentDetails.Ingredient").
		Where("adjustment_id = ?", adjustmentID).
		First(&adjustment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy phiếu kiểm kê"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy thông tin phiếu kiểm kê"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": adjustment})
}

// CreateAdjustment creates a new inventory adjustment with details
func (h *InventoryAdjustmentHandler) CreateAdjustment(c *gin.Context) {
	var req CreateAdjustmentRequest
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

	adjustmentDate, err := time.Parse("2006-01-02", req.AdjustmentDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Định dạng ngày không hợp lệ"})
		return
	}

	// Validate adjustment type
	validTypes := map[string]bool{
		"count":     true, // Kiểm kê định kỳ
		"damage":    true, // Hư hỏng
		"loss":      true, // Mất mát
		"found":     true, // Tìm thấy thừa
		"expired":   true, // Hết hạn
		"other":     true, // Khác
	}
	if !validTypes[req.AdjustmentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Loại điều chỉnh không hợp lệ"})
		return
	}

	// Generate adjustment ID
	adjustmentID := generateAdjustmentID(adjustmentDate)

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

	// Create adjustment header
	adjustment := models.InventoryAdjustment{
		AdjustmentID:    adjustmentID,
		KitchenID:       req.KitchenID,
		AdjustmentDate:  adjustmentDate,
		AdjustmentType:  req.AdjustmentType,
		Reason:          req.Reason,
		Status:          status,
		CreatedByUserID: &userID,
	}

	if err := tx.Create(&adjustment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo phiếu kiểm kê"})
		return
	}

	// Create adjustment details and calculate total value
	var totalValue float64
	for _, detail := range req.AdjustmentDetails {
		quantityDifference := detail.QuantityAfter - detail.QuantityBefore

		var totalItemValue *float64
		if detail.UnitCost != nil {
			value := quantityDifference * (*detail.UnitCost)
			totalItemValue = &value
			totalValue += value
		}

		adjustmentDetail := models.InventoryAdjustmentDetail{
			AdjustmentID:       adjustmentID,
			IngredientID:       detail.IngredientID,
			QuantityBefore:     detail.QuantityBefore,
			QuantityAfter:      detail.QuantityAfter,
			QuantityDifference: quantityDifference,
			Unit:               detail.Unit,
			UnitCost:           detail.UnitCost,
			TotalValue:         totalItemValue,
			Reason:             detail.Reason,
		}

		if err := tx.Create(&adjustmentDetail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo chi tiết phiếu kiểm kê"})
			return
		}
	}

	// Update total value
	if err := tx.Model(&adjustment).Update("total_value", totalValue).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật tổng giá trị"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lưu phiếu kiểm kê"})
		return
	}

	// Reload with relationships
	h.DB.Preload("Kitchen").
		Preload("AdjustmentDetails.Ingredient").
		First(&adjustment, "adjustment_id = ?", adjustmentID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo phiếu kiểm kê thành công",
		"data":    adjustment,
	})
}

// UpdateAdjustment updates an existing adjustment
func (h *InventoryAdjustmentHandler) UpdateAdjustment(c *gin.Context) {
	adjustmentID := c.Param("id")

	var req CreateAdjustmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingAdjustment models.InventoryAdjustment
	if err := h.DB.Where("adjustment_id = ?", adjustmentID).First(&existingAdjustment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy phiếu kiểm kê"})
		return
	}

	if existingAdjustment.Status == "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể sửa phiếu kiểm kê đã duyệt"})
		return
	}

	adjustmentDate, err := time.Parse("2006-01-02", req.AdjustmentDate)
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
	if err := tx.Where("adjustment_id = ?", adjustmentID).Delete(&models.InventoryAdjustmentDetail{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi xóa chi tiết cũ"})
		return
	}

	// Update header
	updates := map[string]interface{}{
		"kitchen_id":      req.KitchenID,
		"adjustment_date": adjustmentDate,
		"adjustment_type": req.AdjustmentType,
		"reason":          req.Reason,
	}

	if err := tx.Model(&existingAdjustment).Updates(updates).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật phiếu kiểm kê"})
		return
	}

	// Create new details
	var totalValue float64
	for _, detail := range req.AdjustmentDetails {
		quantityDifference := detail.QuantityAfter - detail.QuantityBefore

		var totalItemValue *float64
		if detail.UnitCost != nil {
			value := quantityDifference * (*detail.UnitCost)
			totalItemValue = &value
			totalValue += value
		}

		adjustmentDetail := models.InventoryAdjustmentDetail{
			AdjustmentID:       adjustmentID,
			IngredientID:       detail.IngredientID,
			QuantityBefore:     detail.QuantityBefore,
			QuantityAfter:      detail.QuantityAfter,
			QuantityDifference: quantityDifference,
			Unit:               detail.Unit,
			UnitCost:           detail.UnitCost,
			TotalValue:         totalItemValue,
			Reason:             detail.Reason,
		}

		if err := tx.Create(&adjustmentDetail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo chi tiết phiếu kiểm kê"})
			return
		}
	}

	// Update total value
	if err := tx.Model(&existingAdjustment).Update("total_value", totalValue).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật tổng giá trị"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lưu phiếu kiểm kê"})
		return
	}

	h.DB.Preload("Kitchen").
		Preload("AdjustmentDetails.Ingredient").
		First(&existingAdjustment, "adjustment_id = ?", adjustmentID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật phiếu kiểm kê thành công",
		"data":    existingAdjustment,
	})
}

// ApproveAdjustment approves an adjustment and updates inventory stocks
func (h *InventoryAdjustmentHandler) ApproveAdjustment(c *gin.Context) {
	adjustmentID := c.Param("id")
	var userID string
	if identity, ok := c.Get("identity"); ok {
		if v, ok2 := identity.(string); ok2 {
			userID = v
		}
	}

	var adjustment models.InventoryAdjustment
	if err := h.DB.Preload("AdjustmentDetails").
		Where("adjustment_id = ?", adjustmentID).
		First(&adjustment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy phiếu kiểm kê"})
		return
	}

	if adjustment.Status == "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phiếu kiểm kê đã được duyệt"})
		return
	}

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update adjustment status
	now := time.Now()
	updates := map[string]interface{}{
		"status":              "approved",
		"approved_by_user_id": userID,
		"approved_date":       now,
	}

	if err := tx.Model(&adjustment).Updates(updates).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi duyệt phiếu kiểm kê"})
		return
	}

	// Update inventory stocks
	for _, detail := range adjustment.AdjustmentDetails {
		var stock models.InventoryStock
		result := tx.Where("kitchen_id = ? AND ingredient_id = ?",
			adjustment.KitchenID, detail.IngredientID).
			First(&stock)

		if result.Error == nil {
			// Update existing stock to the adjusted quantity
			if err := tx.Model(&stock).Updates(map[string]interface{}{
				"quantity":     detail.QuantityAfter,
				"unit":         detail.Unit,
				"last_updated": now,
			}).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật tồn kho"})
				return
			}
		} else {
			// Create new stock entry
			stock = models.InventoryStock{
				KitchenID:    adjustment.KitchenID,
				IngredientID: detail.IngredientID,
				Quantity:     detail.QuantityAfter,
				Unit:         detail.Unit,
				LastUpdated:  now,
			}
			if err := tx.Create(&stock).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo tồn kho"})
				return
			}
		}

		// Log transaction
		transactionType := "ADJUSTMENT"
		if detail.QuantityDifference > 0 {
			transactionType = "ADJUSTMENT_IN"
		} else if detail.QuantityDifference < 0 {
			transactionType = "ADJUSTMENT_OUT"
		}

		transaction := models.InventoryTransaction{
			KitchenID:       adjustment.KitchenID,
			IngredientID:    detail.IngredientID,
			TransactionType: transactionType,
			TransactionDate: now,
			Quantity:        detail.QuantityDifference,
			Unit:            detail.Unit,
			QuantityBefore:  detail.QuantityBefore,
			QuantityAfter:   detail.QuantityAfter,
			ReferenceType:   strPtr("ADJUSTMENT"),
			ReferenceID:     &adjustmentID,
			CreatedByUserID: &userID,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi ghi log giao dịch"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi hoàn tất duyệt phiếu"})
		return
	}

	h.DB.Preload("Kitchen").
		Preload("ApprovedBy").
		Preload("AdjustmentDetails.Ingredient").
		First(&adjustment, "adjustment_id = ?", adjustmentID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Duyệt phiếu kiểm kê thành công",
		"data":    adjustment,
	})
}

// DeleteAdjustment deletes a draft adjustment
func (h *InventoryAdjustmentHandler) DeleteAdjustment(c *gin.Context) {
	adjustmentID := c.Param("id")

	var adjustment models.InventoryAdjustment
	if err := h.DB.Where("adjustment_id = ?", adjustmentID).First(&adjustment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy phiếu kiểm kê"})
		return
	}

	if adjustment.Status == "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể xóa phiếu kiểm kê đã duyệt"})
		return
	}

	if err := h.DB.Delete(&adjustment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi xóa phiếu kiểm kê"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xóa phiếu kiểm kê thành công"})
}

// Helper function
func generateAdjustmentID(adjustmentDate time.Time) string {
	return "ADJ" + adjustmentDate.Format("20060102") + "-" + strconv.FormatInt(time.Now().UnixNano()%100000, 10)
}
