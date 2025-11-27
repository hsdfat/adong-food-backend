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

type InventoryImportHandler struct {
	DB *gorm.DB
}

func NewInventoryImportHandler(db *gorm.DB) *InventoryImportHandler {
	return &InventoryImportHandler{DB: db}
}

// CreateImportRequest represents the request body for creating an import
type CreateImportRequest struct {
	KitchenID     string                      `json:"kitchenId" binding:"required"`
	ImportDate    string                      `json:"importDate" binding:"required"`
	OrderID       *string                     `json:"orderId"`
	SupplierID    *string                     `json:"supplierId"`
	Status        string                      `json:"status"`
	Notes         *string                     `json:"notes"`
	ImportDetails []CreateImportDetailRequest `json:"importDetails" binding:"required,min=1"`
}

type CreateImportDetailRequest struct {
	IngredientID string  `json:"ingredientId" binding:"required"`
	Quantity     float64 `json:"quantity" binding:"required,gt=0"`
	Unit         string  `json:"unit" binding:"required"`
	UnitPrice    float64 `json:"unitPrice" binding:"required,gt=0"`
	ExpiryDate   *string `json:"expiryDate"`
	BatchNumber  *string `json:"batchNumber"`
	Notes        *string `json:"notes"`
}

// GetAllImports retrieves all inventory imports with pagination and filters
func (h *InventoryImportHandler) GetAllImports(c *gin.Context) {
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
	status := c.Query("status")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")

	var imports []models.InventoryImport
	var total int64

	countQuery := h.DB.Model(&models.InventoryImport{})

	if kitchenID != "" {
		countQuery = countQuery.Where("kitchen_id = ?", kitchenID)
	}
	if status != "" {
		countQuery = countQuery.Where("status = ?", status)
	}
	if fromDate != "" {
		countQuery = countQuery.Where("import_date >= ?", fromDate)
	}
	if toDate != "" {
		countQuery = countQuery.Where("import_date <= ?", toDate)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi đếm phiếu nhập"})
		return
	}

	query := h.DB.Model(&models.InventoryImport{})

	if kitchenID != "" {
		query = query.Where("kitchen_id = ?", kitchenID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if fromDate != "" {
		query = query.Where("import_date >= ?", fromDate)
	}
	if toDate != "" {
		query = query.Where("import_date <= ?", toDate)
	}

	allowedSortFields := map[string]string{
		"import_date":  "import_date",
		"created_date": "created_date",
		"status":       "status",
		"total_amount": "total_amount",
	}
	query = utils.ApplySort(query, params.SortBy, params.SortDir, allowedSortFields)
	if params.SortBy == "" {
		query = query.Order("import_date DESC, created_date DESC")
	}
	query = utils.ApplyPagination(query, params.Page, params.PageSize)

	if err := query.Preload("Kitchen").
		Preload("Supplier").
		Preload("ReceivedBy").
		Preload("ApprovedBy").
		Preload("CreatedBy").
		Find(&imports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy danh sách phiếu nhập"})
		return
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: imports,
		Meta: meta,
	})
}

// GetImportByID retrieves a specific import with details
func (h *InventoryImportHandler) GetImportByID(c *gin.Context) {
	importID := c.Param("id")

	var importRecord models.InventoryImport
	if err := h.DB.Preload("Kitchen").
		Preload("Supplier").
		Preload("Order").
		Preload("ReceivedBy").
		Preload("ApprovedBy").
		Preload("CreatedBy").
		Preload("ImportDetails.Ingredient").
		Where("import_id = ?", importID).
		First(&importRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy phiếu nhập"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy thông tin phiếu nhập"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": importRecord})
}

// CreateImport creates a new inventory import with details
func (h *InventoryImportHandler) CreateImport(c *gin.Context) {
	var req CreateImportRequest
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

	importDate, err := time.Parse("2006-01-02", req.ImportDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Định dạng ngày không hợp lệ"})
		return
	}

	// Generate import ID
	importID := generateImportID(importDate)

	// Start transaction
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

	// Create import header
	importRecord := models.InventoryImport{
		ImportID:        importID,
		KitchenID:       req.KitchenID,
		ImportDate:      importDate,
		OrderID:         req.OrderID,
		SupplierID:      req.SupplierID,
		Status:          status,
		Notes:           req.Notes,
		CreatedByUserID: &userID,
	}

	if err := tx.Create(&importRecord).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo phiếu nhập"})
		return
	}

	// Create import details and calculate total
	var totalAmount float64
	for _, detail := range req.ImportDetails {
		totalPrice := detail.Quantity * detail.UnitPrice
		totalAmount += totalPrice

		var expiryDate *time.Time
		if detail.ExpiryDate != nil {
			expDate, err := time.Parse("2006-01-02", *detail.ExpiryDate)
			if err == nil {
				expiryDate = &expDate
			}
		}

		importDetail := models.InventoryImportDetail{
			ImportID:     importID,
			IngredientID: detail.IngredientID,
			Quantity:     detail.Quantity,
			Unit:         detail.Unit,
			UnitPrice:    detail.UnitPrice,
			TotalPrice:   totalPrice,
			ExpiryDate:   expiryDate,
			BatchNumber:  detail.BatchNumber,
			Notes:        detail.Notes,
		}

		if err := tx.Create(&importDetail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo chi tiết phiếu nhập"})
			return
		}
	}

	// Update total amount
	if err := tx.Model(&importRecord).Update("total_amount", totalAmount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật tổng tiền"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lưu phiếu nhập"})
		return
	}

	// Reload with relationships
	h.DB.Preload("Kitchen").
		Preload("Supplier").
		Preload("ImportDetails.Ingredient").
		First(&importRecord, "import_id = ?", importID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo phiếu nhập thành công",
		"data":    importRecord,
	})
}

// UpdateImport updates an existing import
func (h *InventoryImportHandler) UpdateImport(c *gin.Context) {
	importID := c.Param("id")

	var req CreateImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingImport models.InventoryImport
	if err := h.DB.Where("import_id = ?", importID).First(&existingImport).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy phiếu nhập"})
		return
	}

	if existingImport.Status == "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể sửa phiếu nhập đã duyệt"})
		return
	}

	importDate, err := time.Parse("2006-01-02", req.ImportDate)
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
	if err := tx.Where("import_id = ?", importID).Delete(&models.InventoryImportDetail{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi xóa chi tiết cũ"})
		return
	}

	// Update header
	updates := map[string]interface{}{
		"kitchen_id":  req.KitchenID,
		"import_date": importDate,
		"order_id":    req.OrderID,
		"supplier_id": req.SupplierID,
		"notes":       req.Notes,
	}

	if err := tx.Model(&existingImport).Updates(updates).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật phiếu nhập"})
		return
	}

	// Create new details
	var totalAmount float64
	for _, detail := range req.ImportDetails {
		totalPrice := detail.Quantity * detail.UnitPrice
		totalAmount += totalPrice

		var expiryDate *time.Time
		if detail.ExpiryDate != nil {
			expDate, err := time.Parse("2006-01-02", *detail.ExpiryDate)
			if err == nil {
				expiryDate = &expDate
			}
		}

		importDetail := models.InventoryImportDetail{
			ImportID:     importID,
			IngredientID: detail.IngredientID,
			Quantity:     detail.Quantity,
			Unit:         detail.Unit,
			UnitPrice:    detail.UnitPrice,
			TotalPrice:   totalPrice,
			ExpiryDate:   expiryDate,
			BatchNumber:  detail.BatchNumber,
			Notes:        detail.Notes,
		}

		if err := tx.Create(&importDetail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo chi tiết phiếu nhập"})
			return
		}
	}

	// Update total amount
	if err := tx.Model(&existingImport).Update("total_amount", totalAmount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật tổng tiền"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lưu phiếu nhập"})
		return
	}

	h.DB.Preload("Kitchen").
		Preload("Supplier").
		Preload("ImportDetails.Ingredient").
		First(&existingImport, "import_id = ?", importID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật phiếu nhập thành công",
		"data":    existingImport,
	})
}

// ApproveImport approves an import and updates inventory stocks
func (h *InventoryImportHandler) ApproveImport(c *gin.Context) {
	importID := c.Param("id")
	var userID string
	if identity, ok := c.Get("identity"); ok {
		if v, ok2 := identity.(string); ok2 {
			userID = v
		}
	}

	var importRecord models.InventoryImport
	if err := h.DB.Preload("ImportDetails").
		Where("import_id = ?", importID).
		First(&importRecord).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy phiếu nhập"})
		return
	}

	if importRecord.Status == "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phiếu nhập đã được duyệt"})
		return
	}

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update import status
	now := time.Now()
	updates := map[string]interface{}{
		"status":              "approved",
		"approved_by_user_id": userID,
		"approved_date":       now,
	}

	if err := tx.Model(&importRecord).Updates(updates).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi duyệt phiếu nhập"})
		return
	}

	// Update inventory stocks
	for _, detail := range importRecord.ImportDetails {
		var stock models.InventoryStock
		result := tx.Where("kitchen_id = ? AND ingredient_id = ?",
			importRecord.KitchenID, detail.IngredientID).
			First(&stock)

		quantityBefore := 0.0
		if result.Error == nil {
			quantityBefore = stock.Quantity
			// Update existing stock
			if err := tx.Model(&stock).Updates(map[string]interface{}{
				"quantity":     gorm.Expr("quantity + ?", detail.Quantity),
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
				KitchenID:    importRecord.KitchenID,
				IngredientID: detail.IngredientID,
				Quantity:     detail.Quantity,
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
		transaction := models.InventoryTransaction{
			KitchenID:       importRecord.KitchenID,
			IngredientID:    detail.IngredientID,
			TransactionType: "IMPORT",
			TransactionDate: now,
			Quantity:        detail.Quantity,
			Unit:            detail.Unit,
			QuantityBefore:  quantityBefore,
			QuantityAfter:   quantityBefore + detail.Quantity,
			ReferenceType:   strPtr("IMPORT"),
			ReferenceID:     &importID,
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
		Preload("Supplier").
		Preload("ApprovedBy").
		Preload("ImportDetails.Ingredient").
		First(&importRecord, "import_id = ?", importID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Duyệt phiếu nhập thành công",
		"data":    importRecord,
	})
}

// DeleteImport deletes a draft import
func (h *InventoryImportHandler) DeleteImport(c *gin.Context) {
	importID := c.Param("id")

	var importRecord models.InventoryImport
	if err := h.DB.Where("import_id = ?", importID).First(&importRecord).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy phiếu nhập"})
		return
	}

	if importRecord.Status == "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể xóa phiếu nhập đã duyệt"})
		return
	}

	if err := h.DB.Delete(&importRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi xóa phiếu nhập"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xóa phiếu nhập thành công"})
}

// CreateImportFromRequest creates an import record from an ingredient request
func (h *InventoryImportHandler) CreateImportFromRequest(c *gin.Context) {
	requestID := c.Param("requestId")

	var userID string
	if identity, ok := c.Get("identity"); ok {
		if v, ok2 := identity.(string); ok2 {
			userID = v
		}
	}

	// Get ingredient request with details
	var request models.IngredientRequest
	if err := h.DB.Preload("RequestDetails").
		Where("request_id = ?", requestID).
		First(&request).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy phiếu yêu cầu"})
		return
	}

	if request.Status != "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phiếu yêu cầu chưa được duyệt"})
		return
	}

	// Get the most common supplier from request details (or first one)
	var mainSupplierID *string
	if len(request.RequestDetails) > 0 {
		supplierCount := make(map[string]int)
		for _, detail := range request.RequestDetails {
			if detail.SupplierID != nil {
				supplierCount[*detail.SupplierID]++
			}
		}

		maxCount := 0
		for supplierID, count := range supplierCount {
			if count > maxCount {
				maxCount = count
				sid := supplierID
				mainSupplierID = &sid
			}
		}
	}

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Generate import ID
	importDate := time.Now()
	importID := generateImportID(importDate)

	// Create import header
	importRecord := models.InventoryImport{
		ImportID:        importID,
		KitchenID:       request.KitchenID,
		ImportDate:      importDate,
		OrderID:         &request.OrderID,
		SupplierID:      mainSupplierID,
		Status:          "draft",
		CreatedByUserID: &userID,
	}

	if err := tx.Create(&importRecord).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo phiếu nhập"})
		return
	}

	// Create import details from request details
	var totalAmount float64
	for _, reqDetail := range request.RequestDetails {
		var unitPrice float64
		if reqDetail.UnitPrice != nil {
			unitPrice = *reqDetail.UnitPrice
		}

		totalPrice := reqDetail.Quantity * unitPrice
		totalAmount += totalPrice

		importDetail := models.InventoryImportDetail{
			ImportID:     importID,
			IngredientID: reqDetail.IngredientID,
			Quantity:     reqDetail.Quantity,
			Unit:         reqDetail.Unit,
			UnitPrice:    unitPrice,
			TotalPrice:   totalPrice,
		}

		if err := tx.Create(&importDetail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo chi tiết phiếu nhập"})
			return
		}
	}

	// Update total amount
	if err := tx.Model(&importRecord).Update("total_amount", totalAmount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật tổng tiền"})
		return
	}

	// Update request status to received
	if err := tx.Model(&request).Update("status", "received").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật trạng thái yêu cầu"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lưu phiếu nhập"})
		return
	}

	// Reload with relationships
	h.DB.Preload("Kitchen").
		Preload("Supplier").
		Preload("ImportDetails.Ingredient").
		First(&importRecord, "import_id = ?", importID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo phiếu nhập từ yêu cầu thành công",
		"data":    importRecord,
	})
}

// Helper functions
func generateImportID(importDate time.Time) string {
	return "IM" + importDate.Format("20060102") + "-" + strconv.FormatInt(time.Now().UnixNano()%100000, 10)
}

func strPtr(s string) *string {
	return &s
}
