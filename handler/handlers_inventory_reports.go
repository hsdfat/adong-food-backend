package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InventoryReportsHandler struct {
	DB *gorm.DB
}

func NewInventoryReportsHandler(db *gorm.DB) *InventoryReportsHandler {
	return &InventoryReportsHandler{DB: db}
}

// StockMovementReport represents stock movement for a period
type StockMovementReport struct {
	IngredientID   string  `json:"ingredientId"`
	IngredientName string  `json:"ingredientName"`
	Unit           string  `json:"unit"`
	OpeningStock   float64 `json:"openingStock"`
	StockIn        float64 `json:"stockIn"`
	StockOut       float64 `json:"stockOut"`
	Adjustment     float64 `json:"adjustment"`
	ClosingStock   float64 `json:"closingStock"`
}

// GetStockMovementReport retrieves stock movement report for a period
func (h *InventoryReportsHandler) GetStockMovementReport(c *gin.Context) {
	kitchenID := c.Query("kitchen_id")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")

	if kitchenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cần có kitchen_id"})
		return
	}
	if fromDate == "" || toDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cần có from_date và to_date"})
		return
	}

	fromDateTime, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Định dạng from_date không hợp lệ"})
		return
	}

	toDateTime, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Định dạng to_date không hợp lệ"})
		return
	}

	var movements []StockMovementReport

	query := `
		WITH opening_stocks AS (
			SELECT
				it.ingredient_id,
				i.ingredient_name,
				COALESCE(SUM(CASE
					WHEN it.transaction_date < ? THEN it.quantity
					ELSE 0
				END), 0) as opening_stock,
				MAX(it.unit) as unit
			FROM inventory_transactions it
			JOIN master_ingredients i ON i.ingredient_id = it.ingredient_id
			WHERE it.kitchen_id = ?
			GROUP BY it.ingredient_id, i.ingredient_name
		),
		period_movements AS (
			SELECT
				it.ingredient_id,
				SUM(CASE WHEN it.transaction_type IN ('IMPORT', 'TRANSFER_IN') THEN it.quantity ELSE 0 END) as stock_in,
				SUM(CASE WHEN it.transaction_type IN ('EXPORT') THEN ABS(it.quantity) ELSE 0 END) as stock_out,
				SUM(CASE WHEN it.transaction_type LIKE 'ADJUSTMENT%' THEN it.quantity ELSE 0 END) as adjustment
			FROM inventory_transactions it
			WHERE it.kitchen_id = ?
				AND it.transaction_date >= ?
				AND it.transaction_date <= ?
			GROUP BY it.ingredient_id
		)
		SELECT
			os.ingredient_id,
			os.ingredient_name,
			os.unit,
			os.opening_stock,
			COALESCE(pm.stock_in, 0) as stock_in,
			COALESCE(pm.stock_out, 0) as stock_out,
			COALESCE(pm.adjustment, 0) as adjustment,
			os.opening_stock + COALESCE(pm.stock_in, 0) - COALESCE(pm.stock_out, 0) + COALESCE(pm.adjustment, 0) as closing_stock
		FROM opening_stocks os
		LEFT JOIN period_movements pm ON pm.ingredient_id = os.ingredient_id
		WHERE os.opening_stock != 0 OR pm.stock_in IS NOT NULL OR pm.stock_out IS NOT NULL OR pm.adjustment IS NOT NULL
		ORDER BY os.ingredient_name
	`

	if err := h.DB.Raw(query, fromDateTime, kitchenID, kitchenID, fromDateTime, toDateTime).Scan(&movements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy báo cáo xuất nhập tồn"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      movements,
		"from_date": fromDate,
		"to_date":   toDate,
		"count":     len(movements),
	})
}

// ExpiryAlert represents ingredients nearing expiry
type ExpiryAlert struct {
	ImportDetailID int     `json:"importDetailId"`
	ImportID       string  `json:"importId"`
	IngredientID   string  `json:"ingredientId"`
	IngredientName string  `json:"ingredientName"`
	Quantity       float64 `json:"quantity"`
	Unit           string  `json:"unit"`
	ExpiryDate     string  `json:"expiryDate"`
	DaysToExpiry   int     `json:"daysToExpiry"`
	BatchNumber    *string `json:"batchNumber,omitempty"`
}

// GetExpiryAlerts retrieves ingredients nearing expiry
func (h *InventoryReportsHandler) GetExpiryAlerts(c *gin.Context) {
	kitchenID := c.Query("kitchen_id")
	daysAhead := c.DefaultQuery("days_ahead", "30") // Default 30 days

	if kitchenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cần có kitchen_id"})
		return
	}

	var alerts []ExpiryAlert

	query := `
		SELECT
			iid.import_detail_id,
			iid.import_id,
			iid.ingredient_id,
			i.ingredient_name,
			iid.quantity,
			iid.unit,
			iid.expiry_date::text as expiry_date,
			DATE_PART('day', iid.expiry_date - CURRENT_DATE)::int as days_to_expiry,
			iid.batch_number
		FROM inventory_import_details iid
		JOIN inventory_imports ii ON ii.import_id = iid.import_id
		JOIN master_ingredients i ON i.ingredient_id = iid.ingredient_id
		WHERE ii.kitchen_id = ?
			AND ii.status = 'approved'
			AND iid.expiry_date IS NOT NULL
			AND iid.expiry_date <= CURRENT_DATE + INTERVAL '` + daysAhead + ` days'
			AND iid.expiry_date >= CURRENT_DATE
		ORDER BY iid.expiry_date ASC, i.ingredient_name
	`

	if err := h.DB.Raw(query, kitchenID).Scan(&alerts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy cảnh báo hết hạn"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       alerts,
		"days_ahead": daysAhead,
		"count":      len(alerts),
	})
}

// StockValueTrend represents stock value over time
type StockValueTrend struct {
	Date       string  `json:"date"`
	TotalValue float64 `json:"totalValue"`
	TotalItems int     `json:"totalItems"`
}

// GetStockValueTrend retrieves stock value trend over time
func (h *InventoryReportsHandler) GetStockValueTrend(c *gin.Context) {
	kitchenID := c.Query("kitchen_id")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	interval := c.DefaultQuery("interval", "day") // day, week, month

	if kitchenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cần có kitchen_id"})
		return
	}
	if fromDate == "" || toDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cần có from_date và to_date"})
		return
	}

	var dateFormat string
	switch interval {
	case "week":
		dateFormat = "YYYY-IW"
	case "month":
		dateFormat = "YYYY-MM"
	default:
		dateFormat = "YYYY-MM-DD"
	}

	var trends []StockValueTrend

	query := `
		WITH date_series AS (
			SELECT generate_series(
				?::date,
				?::date,
				'1 ` + interval + `'::interval
			)::date as date
		),
		daily_values AS (
			SELECT
				ds.date,
				COALESCE(SUM(
					(SELECT quantity_after
					 FROM inventory_transactions it2
					 WHERE it2.ingredient_id = it.ingredient_id
					 	AND it2.kitchen_id = ?
					 	AND it2.transaction_date <= ds.date
					 ORDER BY it2.transaction_date DESC, it2.transaction_id DESC
					 LIMIT 1) *
					COALESCE((SELECT unit_price
							  FROM supplier_price_list spl
							  WHERE spl.ingredient_id = it.ingredient_id
							  	AND spl.active = true
							  LIMIT 1), 0)
				), 0) as total_value,
				COUNT(DISTINCT it.ingredient_id) as total_items
			FROM date_series ds
			LEFT JOIN inventory_transactions it ON it.kitchen_id = ?
			GROUP BY ds.date
		)
		SELECT
			TO_CHAR(date, ?) as date,
			SUM(total_value) as total_value,
			MAX(total_items) as total_items
		FROM daily_values
		GROUP BY TO_CHAR(date, ?)
		ORDER BY date
	`

	if err := h.DB.Raw(query, fromDate, toDate, kitchenID, kitchenID, dateFormat, dateFormat).Scan(&trends).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy xu hướng giá trị tồn kho"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      trends,
		"from_date": fromDate,
		"to_date":   toDate,
		"interval":  interval,
		"count":     len(trends),
	})
}

// TransactionSummary represents transaction summary by type
type TransactionSummary struct {
	TransactionType  string  `json:"transactionType"`
	TotalQuantity    float64 `json:"totalQuantity"`
	TransactionCount int     `json:"transactionCount"`
}

// GetTransactionSummary retrieves transaction summary by type
func (h *InventoryReportsHandler) GetTransactionSummary(c *gin.Context) {
	kitchenID := c.Query("kitchen_id")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")

	if kitchenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cần có kitchen_id"})
		return
	}

	var summary []TransactionSummary

	query := `
		SELECT
			transaction_type,
			SUM(ABS(quantity)) as total_quantity,
			COUNT(*) as transaction_count
		FROM inventory_transactions
		WHERE kitchen_id = ?
	`

	params := []interface{}{kitchenID}

	if fromDate != "" {
		query += " AND transaction_date >= ?"
		params = append(params, fromDate)
	}
	if toDate != "" {
		query += " AND transaction_date <= ?"
		params = append(params, toDate)
	}

	query += " GROUP BY transaction_type ORDER BY transaction_type"

	if err := h.DB.Raw(query, params...).Scan(&summary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy tổng hợp giao dịch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  summary,
		"count": len(summary),
	})
}

// TopConsumedIngredient represents most consumed ingredients
type TopConsumedIngredient struct {
	IngredientID   string  `json:"ingredientId"`
	IngredientName string  `json:"ingredientName"`
	TotalConsumed  float64 `json:"totalConsumed"`
	Unit           string  `json:"unit"`
	ExportCount    int     `json:"exportCount"`
}

// GetTopConsumedIngredients retrieves most consumed ingredients
func (h *InventoryReportsHandler) GetTopConsumedIngredients(c *gin.Context) {
	kitchenID := c.Query("kitchen_id")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	limit := c.DefaultQuery("limit", "10")

	if kitchenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cần có kitchen_id"})
		return
	}

	var topIngredients []TopConsumedIngredient

	query := `
		SELECT
			it.ingredient_id,
			i.ingredient_name,
			SUM(ABS(it.quantity)) as total_consumed,
			it.unit,
			COUNT(*) as export_count
		FROM inventory_transactions it
		JOIN master_ingredients i ON i.ingredient_id = it.ingredient_id
		WHERE it.kitchen_id = ?
			AND it.transaction_type = 'EXPORT'
	`

	params := []interface{}{kitchenID}

	if fromDate != "" {
		query += " AND it.transaction_date >= ?"
		params = append(params, fromDate)
	}
	if toDate != "" {
		query += " AND it.transaction_date <= ?"
		params = append(params, toDate)
	}

	query += `
		GROUP BY it.ingredient_id, i.ingredient_name, it.unit
		ORDER BY total_consumed DESC
		LIMIT ` + limit

	if err := h.DB.Raw(query, params...).Scan(&topIngredients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi lấy nguyên liệu tiêu thụ nhiều nhất"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  topIngredients,
		"count": len(topIngredients),
		"limit": limit,
	})
}
