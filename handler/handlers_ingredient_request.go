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

type IngredientRequestHandler struct {
	DB *gorm.DB
}

func NewIngredientRequestHandler(db *gorm.DB) *IngredientRequestHandler {
	return &IngredientRequestHandler{DB: db}
}

// CreateRequestInput represents the request body for creating an ingredient request
type CreateRequestInput struct {
	OrderID        string                     `json:"orderId" binding:"required"`
	KitchenID      string                     `json:"kitchenId" binding:"required"`
	RequestDate    string                     `json:"requestDate" binding:"required"`
	RequiredDate   string                     `json:"requiredDate" binding:"required"`
	Status         string                     `json:"status"`
	Notes          *string                    `json:"notes"`
	RequestDetails []CreateRequestDetailInput `json:"requestDetails" binding:"required,min=1"`
}

type CreateRequestDetailInput struct {
	IngredientID string   `json:"ingredientId" binding:"required"`
	Quantity     float64  `json:"quantity" binding:"required,gt=0"`
	Unit         string   `json:"unit" binding:"required"`
	SupplierID   *string  `json:"supplierId"`
	UnitPrice    *float64 `json:"unitPrice"`
	Notes        *string  `json:"notes"`
}

// GetAllRequests retrieves all ingredient requests with pagination and filters
func (h *IngredientRequestHandler) GetAllRequests(c *gin.Context) {
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
	orderID := c.Query("order_id")
	status := c.Query("status")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")

	var requests []models.IngredientRequest
	var total int64

	countQuery := h.DB.Model(&models.IngredientRequest{})

	if kitchenID != "" {
		countQuery = countQuery.Where("kitchen_id = ?", kitchenID)
	}
	if orderID != "" {
		countQuery = countQuery.Where("order_id = ?", orderID)
	}
	if status != "" {
		countQuery = countQuery.Where("status = ?", status)
	}
	if fromDate != "" {
		countQuery = countQuery.Where("request_date >= ?", fromDate)
	}
	if toDate != "" {
		countQuery = countQuery.Where("request_date <= ?", toDate)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi đếm phiếu yêu cầu"})
		return
	}

	query := h.DB.Model(&models.IngredientRequest{})

	if kitchenID != "" {
		query = query.Where("kitchen_id = ?", kitchenID)
	}
	if orderID != "" {
		query = query.Where("order_id = ?", orderID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if fromDate != "" {
		query = query.Where("request_date >= ?", fromDate)
	}
	if toDate != "" {
		query = query.Where("request_date <= ?", toDate)
	}

	allowedSortFields := map[string]string{
		"request_date":  "request_date",
		"required_date": "required_date",
		"created_date":  "created_date",
		"status":        "status",
		"total_amount":  "total_amount",
	}
	query = utils.ApplySort(query, params.SortBy, params.SortDir, allowedSortFields)
	if params.SortBy == "" {
		query = query.Order("request_date DESC, created_date DESC")
	}
	query = utils.ApplyPagination(query, params.Page, params.PageSize)

	if err := query.Preload("Kitchen").
		Preload("Order").
		Preload("CreatedBy").
		Preload("ApprovedBy").
		Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy danh sách phiếu yêu cầu"})
		return
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: requests,
		Meta: meta,
	})
}

// GetRequestByID retrieves a specific ingredient request with details
func (h *IngredientRequestHandler) GetRequestByID(c *gin.Context) {
	requestID := c.Param("id")

	var request models.IngredientRequest
	if err := h.DB.Preload("Kitchen").
		Preload("Order").
		Preload("CreatedBy").
		Preload("ApprovedBy").
		Preload("RequestDetails.Ingredient").
		Preload("RequestDetails.Supplier").
		Where("request_id = ?", requestID).
		First(&request).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy phiếu yêu cầu"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy thông tin phiếu yêu cầu"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": request})
}

// CreateRequest creates a new ingredient request
func (h *IngredientRequestHandler) CreateRequest(c *gin.Context) {
	var req CreateRequestInput
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

	requestDate, err := time.Parse("2006-01-02", req.RequestDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Định dạng ngày yêu cầu không hợp lệ"})
		return
	}

	requiredDate, err := time.Parse("2006-01-02", req.RequiredDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Định dạng ngày cần không hợp lệ"})
		return
	}

	// Generate request ID
	requestID := generateRequestID(requestDate)

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	status := req.Status
	if status == "" {
		status = "pending"
	}

	// Create request header
	request := models.IngredientRequest{
		RequestID:       requestID,
		OrderID:         req.OrderID,
		KitchenID:       req.KitchenID,
		RequestDate:     requestDate,
		RequiredDate:    requiredDate,
		Status:          status,
		Notes:           req.Notes,
		CreatedByUserID: &userID,
	}

	if err := tx.Create(&request).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo phiếu yêu cầu"})
		return
	}

	// Create request details and calculate total
	var totalAmount float64
	for _, detail := range req.RequestDetails {
		var totalPrice *float64
		if detail.UnitPrice != nil {
			price := *detail.UnitPrice * detail.Quantity
			totalPrice = &price
			totalAmount += price
		}

		requestDetail := models.IngredientRequestDetail{
			RequestID:    requestID,
			IngredientID: detail.IngredientID,
			Quantity:     detail.Quantity,
			Unit:         detail.Unit,
			SupplierID:   detail.SupplierID,
			UnitPrice:    detail.UnitPrice,
			TotalPrice:   totalPrice,
			Notes:        detail.Notes,
		}

		if err := tx.Create(&requestDetail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo chi tiết phiếu yêu cầu"})
			return
		}
	}

	// Update total amount
	if err := tx.Model(&request).Update("total_amount", totalAmount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật tổng tiền"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lưu phiếu yêu cầu"})
		return
	}

	// Reload with relationships
	h.DB.Preload("Kitchen").
		Preload("Order").
		Preload("RequestDetails.Ingredient").
		Preload("RequestDetails.Supplier").
		First(&request, "request_id = ?", requestID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo phiếu yêu cầu thành công",
		"data":    request,
	})
}

// CreateRequestFromOrder creates an ingredient request from an order
func (h *IngredientRequestHandler) CreateRequestFromOrder(c *gin.Context) {
	orderID := c.Param("orderId")

	var userID string
	if identity, ok := c.Get("identity"); ok {
		if v, ok2 := identity.(string); ok2 {
			userID = v
		}
	}

	// Get order with selected suppliers
	var order models.Order
	if err := h.DB.Preload("Kitchen").
		Where("order_id = ?", orderID).
		First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy đơn hàng"})
		return
	}

	// Get order ingredients with selected suppliers
	type OrderIngredientWithSupplier struct {
		IngredientID string
		Quantity     float64
		Unit         string
		SupplierID   *string
		UnitPrice    *float64
	}

	var ingredients []OrderIngredientWithSupplier
	query := `
		SELECT
			oi.ingredient_id,
			oi.quantity,
			oi.unit,
			ois.supplier_id,
			ois.unit_price
		FROM order_ingredients oi
		LEFT JOIN order_ingredient_suppliers ois ON ois.order_id = oi.order_id
			AND ois.ingredient_id = oi.ingredient_id
			AND ois.is_selected = true
		WHERE oi.order_id = ?
	`

	if err := h.DB.Raw(query, orderID).Scan(&ingredients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy nguyên liệu đơn hàng"})
		return
	}

	if len(ingredients) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Đơn hàng không có nguyên liệu"})
		return
	}

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Generate request ID
	requestDate := time.Now()
	requestID := generateRequestID(requestDate)

	// Create request header
	// models.Order does not contain RequiredDate, so default to the generated requestDate.
	requiredDate := requestDate
	request := models.IngredientRequest{
		RequestID:       requestID,
		OrderID:         orderID,
		KitchenID:       order.KitchenID,
		RequestDate:     requestDate,
		RequiredDate:    requiredDate,
		Status:          "pending",
		CreatedByUserID: &userID,
	}

	if err := tx.Create(&request).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo phiếu yêu cầu"})
		return
	}

	// Create request details
	var totalAmount float64
	for _, ing := range ingredients {
		var totalPrice *float64
		if ing.UnitPrice != nil {
			price := *ing.UnitPrice * ing.Quantity
			totalPrice = &price
			totalAmount += price
		}

		requestDetail := models.IngredientRequestDetail{
			RequestID:    requestID,
			IngredientID: ing.IngredientID,
			Quantity:     ing.Quantity,
			Unit:         ing.Unit,
			SupplierID:   ing.SupplierID,
			UnitPrice:    ing.UnitPrice,
			TotalPrice:   totalPrice,
		}

		if err := tx.Create(&requestDetail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo chi tiết phiếu yêu cầu"})
			return
		}
	}

	// Update total amount
	if err := tx.Model(&request).Update("total_amount", totalAmount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật tổng tiền"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lưu phiếu yêu cầu"})
		return
	}

	// Reload with relationships
	h.DB.Preload("Kitchen").
		Preload("Order").
		Preload("RequestDetails.Ingredient").
		Preload("RequestDetails.Supplier").
		First(&request, "request_id = ?", requestID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo phiếu yêu cầu từ đơn hàng thành công",
		"data":    request,
	})
}

// UpdateRequest updates an existing ingredient request
func (h *IngredientRequestHandler) UpdateRequest(c *gin.Context) {
	requestID := c.Param("id")

	var req CreateRequestInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingRequest models.IngredientRequest
	if err := h.DB.Where("request_id = ?", requestID).First(&existingRequest).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy phiếu yêu cầu"})
		return
	}

	if existingRequest.Status == "approved" || existingRequest.Status == "received" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể sửa phiếu yêu cầu đã duyệt hoặc đã nhận hàng"})
		return
	}

	requestDate, err := time.Parse("2006-01-02", req.RequestDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Định dạng ngày yêu cầu không hợp lệ"})
		return
	}

	requiredDate, err := time.Parse("2006-01-02", req.RequiredDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Định dạng ngày cần không hợp lệ"})
		return
	}

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete existing details
	if err := tx.Where("request_id = ?", requestID).Delete(&models.IngredientRequestDetail{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi xóa chi tiết cũ"})
		return
	}

	// Update header
	updates := map[string]interface{}{
		"order_id":      req.OrderID,
		"kitchen_id":    req.KitchenID,
		"request_date":  requestDate,
		"required_date": requiredDate,
		"notes":         req.Notes,
	}

	if err := tx.Model(&existingRequest).Updates(updates).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật phiếu yêu cầu"})
		return
	}

	// Create new details
	var totalAmount float64
	for _, detail := range req.RequestDetails {
		var totalPrice *float64
		if detail.UnitPrice != nil {
			price := *detail.UnitPrice * detail.Quantity
			totalPrice = &price
			totalAmount += price
		}

		requestDetail := models.IngredientRequestDetail{
			RequestID:    requestID,
			IngredientID: detail.IngredientID,
			Quantity:     detail.Quantity,
			Unit:         detail.Unit,
			SupplierID:   detail.SupplierID,
			UnitPrice:    detail.UnitPrice,
			TotalPrice:   totalPrice,
			Notes:        detail.Notes,
		}

		if err := tx.Create(&requestDetail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi tạo chi tiết phiếu yêu cầu"})
			return
		}
	}

	// Update total amount
	if err := tx.Model(&existingRequest).Update("total_amount", totalAmount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi cập nhật tổng tiền"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lưu phiếu yêu cầu"})
		return
	}

	h.DB.Preload("Kitchen").
		Preload("Order").
		Preload("RequestDetails.Ingredient").
		Preload("RequestDetails.Supplier").
		First(&existingRequest, "request_id = ?", requestID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Cập nhật phiếu yêu cầu thành công",
		"data":    existingRequest,
	})
}

// ApproveRequest approves an ingredient request
func (h *IngredientRequestHandler) ApproveRequest(c *gin.Context) {
	requestID := c.Param("id")
	var userID string
	if identity, ok := c.Get("identity"); ok {
		if v, ok2 := identity.(string); ok2 {
			userID = v
		}
	}

	var request models.IngredientRequest
	if err := h.DB.Where("request_id = ?", requestID).First(&request).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy phiếu yêu cầu"})
		return
	}

	if request.Status == "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phiếu yêu cầu đã được duyệt"})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":              "approved",
		"approved_by_user_id": userID,
		"approved_date":       now,
	}

	if err := h.DB.Model(&request).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi duyệt phiếu yêu cầu"})
		return
	}

	h.DB.Preload("Kitchen").
		Preload("Order").
		Preload("ApprovedBy").
		Preload("RequestDetails.Ingredient").
		Preload("RequestDetails.Supplier").
		First(&request, "request_id = ?", requestID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Duyệt phiếu yêu cầu thành công",
		"data":    request,
	})
}

// DeleteRequest deletes a pending request
func (h *IngredientRequestHandler) DeleteRequest(c *gin.Context) {
	requestID := c.Param("id")

	var request models.IngredientRequest
	if err := h.DB.Where("request_id = ?", requestID).First(&request).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy phiếu yêu cầu"})
		return
	}

	if request.Status == "approved" || request.Status == "received" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể xóa phiếu yêu cầu đã duyệt hoặc đã nhận hàng"})
		return
	}

	if err := h.DB.Delete(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi xóa phiếu yêu cầu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xóa phiếu yêu cầu thành công"})
}

// Helper function
func generateRequestID(requestDate time.Time) string {
	return "RQ" + requestDate.Format("20060102") + "-" + strconv.FormatInt(time.Now().UnixNano()%100000, 10)
}
