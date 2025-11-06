package handler

import (
	"adong-be/logger"
	"adong-be/models"
	"adong-be/store"
	"adong-be/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetOrders lists orders with filters: kitchen_id, status, date range, dish_id, ingredient_id
func GetOrders(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("GetOrders called", "user_id", uid)
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Log.Error("GetOrders bind query error", "error", err)
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
	dishID := c.Query("dish_id")
	ingredientID := c.Query("ingredient_id")

	// Get user role to check if user is Admin
	var userRole string
	if identity, ok := c.Get("identity"); ok {
		if userID, ok2 := identity.(string); ok2 {
			var user models.User
			if err := store.DB.GormClient.Select("role").First(&user, "user_id = ?", userID).Error; err == nil {
				userRole = user.Role
			}
		}
	}

	var total int64
	var orders []models.Order

	// Use separate queries for counting and data to avoid DISTINCT affecting selected columns
	dataDB := store.DB.GormClient.Model(&models.Order{})
	countDB := store.DB.GormClient.Model(&models.Order{})

	// Filter by created_by_user_id if user is not Admin
	if userRole != "Admin" {
		if identity, ok := c.Get("identity"); ok {
			if userID, ok2 := identity.(string); ok2 {
				dataDB = dataDB.Where("created_by_user_id = ?", userID)
				countDB = countDB.Where("created_by_user_id = ?", userID)
			}
		}
	}

	// Filters
	if params.Search != "" {
		dataDB = dataDB.Where("note ILIKE ? OR order_id ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
		countDB = countDB.Where("note ILIKE ? OR order_id ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}
	if kitchenID != "" {
		dataDB = dataDB.Where("kitchen_id = ?", kitchenID)
		countDB = countDB.Where("kitchen_id = ?", kitchenID)
	}
	if status != "" {
		dataDB = dataDB.Where("status = ?", status)
		countDB = countDB.Where("status = ?", status)
	}
	if fromDate != "" {
		if t, err := time.Parse("2006-01-02", fromDate); err == nil {
			dataDB = dataDB.Where("order_date >= ?", t)
			countDB = countDB.Where("order_date >= ?", t)
		} else {
			dataDB = dataDB.Where("order_date >= ?", fromDate)
			countDB = countDB.Where("order_date >= ?", fromDate)
		}
	}
	if toDate != "" {
		if t, err := time.Parse("2006-01-02", toDate); err == nil {
			dataDB = dataDB.Where("order_date < ?", t.Add(24*time.Hour))
			countDB = countDB.Where("order_date < ?", t.Add(24*time.Hour))
		} else {
			dataDB = dataDB.Where("order_date <= ?", toDate)
			countDB = countDB.Where("order_date <= ?", toDate)
		}
	}
	if dishID != "" {
		dataDB = dataDB.Joins("JOIN order_details od ON od.order_id = orders.order_id").Where("od.dish_id = ?", dishID)
		countDB = countDB.Joins("JOIN order_details od ON od.order_id = orders.order_id").Where("od.dish_id = ?", dishID)
	}
	if ingredientID != "" {
		dataDB = dataDB.Joins("JOIN order_details od2 ON od2.order_id = orders.order_id").
			Joins("JOIN order_ingredients oi ON oi.order_detail_id = od2.order_detail_id").
			Where("oi.ingredient_id = ?", ingredientID)
		countDB = countDB.Joins("JOIN order_details od2 ON od2.order_id = orders.order_id").
			Joins("JOIN order_ingredients oi ON oi.order_detail_id = od2.order_detail_id").
			Where("oi.ingredient_id = ?", ingredientID)
	}

	// Count distinct orders
	if err := countDB.Distinct("orders.order_id").Count(&total).Error; err != nil {
		logger.Log.Error("GetOrders count error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Sorting
	allowedSort := map[string]string{
		"order_id":     "orders.order_id",
		"order_date":   "orders.order_date",
		"status":       "orders.status",
		"created_date": "orders.created_date",
	}
	dataDB = utils.ApplySort(dataDB, params.SortBy, params.SortDir, allowedSort)

	// Pagination
	dataDB = utils.ApplyPagination(dataDB, params.Page, params.PageSize)

	// Fetch and preload relations for DTO
	if err := dataDB.Select("orders.*").
		Preload("Kitchen").
		Preload("CreatedBy").
		Preload("Details.Dish").
		Preload("Details.Ingredients.Ingredient").
		Preload("SupplementaryFoods.Ingredient").
		Find(&orders).Error; err != nil {
		logger.Log.Error("GetOrders query error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map to DTOs
	dtos := make([]models.OrderDTO, len(orders))
	for i := range orders {
		dtos[i] = convertOrderToDTO(&orders[i], true)
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{Data: dtos, Meta: meta})
}

// GetOrder returns a single order with full details
func GetOrder(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("GetOrder called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	var order models.Order
	if err := store.DB.GormClient.
		Preload("Kitchen").
		Preload("CreatedBy").
		Preload("Details.Dish").
		Preload("Details.Ingredients.Ingredient").
		Preload("SupplementaryFoods.Ingredient").
		First(&order, "order_id = ?", id).Error; err != nil {
		logger.Log.Error("GetOrder not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	dto := convertOrderToDTO(&order, true)
	c.JSON(http.StatusOK, dto)
}

// CreateOrder creates a new order with nested details/ingredients/supplementary foods
func CreateOrder(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("CreateOrder called", "user_id", uid)
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		logger.Log.Error("CreateOrder bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from authentication middleware
	if identity, ok := c.Get("identity"); ok {
		if v, ok2 := identity.(string); ok2 {
			order.CreatedByUserID = v
		}
	}

	// Auto-generate OrderID if not provided
	if order.OrderID == "" {
		order.OrderID = uuid.New().String()
		logger.Log.Info("CreateOrder auto-generated OrderID", "orderId", order.OrderID)
	}

	tx := store.DB.GormClient.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Store details and supplementary foods temporarily to avoid GORM auto-saving them
	details := order.Details
	supplementaryFoods := order.SupplementaryFoods
	order.Details = nil
	order.SupplementaryFoods = nil

	// Create order without details/supplementary foods
	if err := tx.Create(&order).Error; err != nil {
		logger.Log.Error("CreateOrder create header error", "error", err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create details and nested ingredients
	for i := range details {
		details[i].OrderID = order.OrderID
		details[i].OrderDetailID = 0 // Ensure auto-increment

		// Store ingredients temporarily to avoid GORM auto-saving them
		ingredients := details[i].Ingredients
		details[i].Ingredients = nil

		// Create order detail without ingredients
		if err := tx.Create(&details[i]).Error; err != nil {
			logger.Log.Error("CreateOrder create detail error", "error", err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Create ingredients manually after order detail is created
		for j := range ingredients {
			ingredients[j].OrderDetailID = details[i].OrderDetailID
			ingredients[j].OrderIngredientID = 0 // Ensure auto-increment

			// Calculate quantity if it's 0 or missing (similar to how summary queries work)
			if ingredients[j].Quantity <= 0 {
				if ingredients[j].StandardPerPortion > 0 && details[i].Portions > 0 {
					ingredients[j].Quantity = ingredients[j].StandardPerPortion * float64(details[i].Portions)
				} else {
					// If quantity can't be calculated and is 0, skip this ingredient
					logger.Log.Warn("CreateOrder skipping ingredient with invalid quantity",
						"ingredient_id", ingredients[j].IngredientID,
						"quantity", ingredients[j].Quantity,
						"standard_per_portion", ingredients[j].StandardPerPortion,
						"portions", details[i].Portions)
					continue
				}
			}

			if err := tx.Create(&ingredients[j]).Error; err != nil {
				logger.Log.Error("CreateOrder create ingredient error", "error", err)
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	// Create supplementary foods
	for i := range supplementaryFoods {
		supplementaryFoods[i].OrderID = order.OrderID
		supplementaryFoods[i].SupplementaryID = 0 // Ensure auto-increment

		// Calculate quantity if it's 0 or missing (similar to how summary queries work)
		if supplementaryFoods[i].Quantity <= 0 {
			if supplementaryFoods[i].StandardPerPortion > 0 && supplementaryFoods[i].Portions > 0 {
				supplementaryFoods[i].Quantity = supplementaryFoods[i].StandardPerPortion * float64(supplementaryFoods[i].Portions)
			} else {
				// If quantity can't be calculated and is 0, skip this supplementary food
				logger.Log.Warn("CreateOrder skipping supplementary food with invalid quantity",
					"ingredient_id", supplementaryFoods[i].IngredientID,
					"quantity", supplementaryFoods[i].Quantity,
					"standard_per_portion", supplementaryFoods[i].StandardPerPortion,
					"portions", supplementaryFoods[i].Portions)
				continue
			}
		}

		if err := tx.Create(&supplementaryFoods[i]).Error; err != nil {
			logger.Log.Error("CreateOrder create supplementary error", "error", err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.Log.Error("CreateOrder commit error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload with relations
	store.DB.GormClient.
		Preload("Kitchen").
		Preload("CreatedBy").
		Preload("Details.Dish").
		Preload("Details.Ingredients.Ingredient").
		Preload("SupplementaryFoods.Ingredient").
		First(&order, "order_id = ?", order.OrderID)

	dto := convertOrderToDTO(&order, true)
	c.JSON(http.StatusCreated, dto)
}

// UpdateOrderStatus updates only the status of an order (PATCH method)
func UpdateOrderStatus(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("UpdateOrderStatus called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")

	// Check if order exists
	var order models.Order
	if err := store.DB.GormClient.First(&order, "order_id = ?", id).Error; err != nil {
		logger.Log.Error("UpdateOrderStatus not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Define a struct to accept only status field
	var updateData struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Log.Error("UpdateOrderStatus bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update only the status field
	if err := store.DB.GormClient.Model(&order).Update("status", updateData.Status).Error; err != nil {
		logger.Log.Error("UpdateOrderStatus db error", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload with relations
	if err := store.DB.GormClient.
		Preload("Kitchen").
		Preload("CreatedBy").
		Preload("Details.Dish").
		Preload("Details.Ingredients.Ingredient").
		Preload("SupplementaryFoods.Ingredient").
		First(&order, "order_id = ?", id).Error; err != nil {
		logger.Log.Error("UpdateOrderStatus reload error", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dto := convertOrderToDTO(&order, true)
	c.JSON(http.StatusOK, dto)
}

// DeleteOrder deletes an order by id (cascade removes children)
func DeleteOrder(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("DeleteOrder called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.Order{}, "order_id = ?", id).Error; err != nil {
		logger.Log.Error("DeleteOrder db error", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

// IngredientTotal represents total usage per ingredient for an order
type IngredientTotal struct {
	IngredientID   string  `json:"ingredientId"`
	IngredientName string  `json:"ingredientName"`
	Unit           string  `json:"unit"`
	TotalQuantity  float64 `json:"totalQuantity"`
}

// GetOrderIngredientsSummary returns totals of ingredients for an order (details + supplementary)
func GetOrderIngredientsSummary(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("GetOrderIngredientsSummary called", "order_id", c.Param("id"), "user_id", uid)
	orderID := c.Param("id")

	var results []IngredientTotal
	sql := `
        SELECT x.ingredient_id AS ingredient_id,
               COALESCE(mi.ingredient_name, '') AS ingredient_name,
               x.unit AS unit,
               COALESCE(SUM(x.total_qty)::double precision, 0) AS total_quantity
        FROM (
            SELECT oi.ingredient_id,
                   oi.unit,
                   COALESCE(oi.quantity, oi.standard_per_portion * od.portions) AS total_qty
            FROM order_ingredients oi
            JOIN order_details od ON od.order_detail_id = oi.order_detail_id
            WHERE od.order_id = ?
            UNION ALL
            SELECT osf.ingredient_id,
                   osf.unit,
                   COALESCE(osf.quantity, osf.standard_per_portion * osf.portions) AS total_qty
            FROM order_supplementary_foods osf
            WHERE osf.order_id = ?
        ) x
        LEFT JOIN master_ingredients mi ON mi.ingredient_id = x.ingredient_id
        GROUP BY x.ingredient_id, mi.ingredient_name, x.unit
        ORDER BY mi.ingredient_name`

	if err := store.DB.GormClient.Raw(sql, orderID, orderID).Scan(&results).Error; err != nil {
		logger.Log.Error("GetOrderIngredientsSummary db error", "order_id", orderID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetOrderIngredientSummary returns total for a specific ingredient in an order
func GetOrderIngredientSummary(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("GetOrderIngredientSummary called", "order_id", c.Param("id"), "ingredient_id", c.Param("ingredientId"), "user_id", uid)
	orderID := c.Param("id")
	ingredientID := c.Param("ingredientId")

	var result IngredientTotal
	sql := `
        SELECT x.ingredient_id AS ingredient_id,
               COALESCE(mi.ingredient_name, '') AS ingredient_name,
               x.unit AS unit,
               COALESCE(SUM(x.total_qty)::double precision, 0) AS total_quantity
        FROM (
            SELECT oi.ingredient_id,
                   oi.unit,
                   COALESCE(oi.quantity, oi.standard_per_portion * od.portions) AS total_qty
            FROM order_ingredients oi
            JOIN order_details od ON od.order_detail_id = oi.order_detail_id
            WHERE od.order_id = ? AND oi.ingredient_id = ?
            UNION ALL
            SELECT osf.ingredient_id,
                   osf.unit,
                   COALESCE(osf.quantity, osf.standard_per_portion * osf.portions) AS total_qty
            FROM order_supplementary_foods osf
            WHERE osf.order_id = ? AND osf.ingredient_id = ?
        ) x
        LEFT JOIN master_ingredients mi ON mi.ingredient_id = x.ingredient_id
        GROUP BY x.ingredient_id, mi.ingredient_name, x.unit
        ORDER BY mi.ingredient_name`

	if err := store.DB.GormClient.Raw(sql, orderID, ingredientID, orderID, ingredientID).Scan(&result).Error; err != nil {
		logger.Log.Error("GetOrderIngredientSummary db error", "order_id", orderID, "ingredient_id", ingredientID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// SaveOrderIngredientsWithSupplier - Save ingredients with selected supplier for an order
func SaveOrderIngredientsWithSupplier(c *gin.Context) {
	uid, _ := c.Get("identity")
	orderID := c.Param("id")
	logger.Log.Info("SaveOrderIngredientsWithSupplier called", "order_id", orderID, "user_id", uid)

	// Define request structure
	var request struct {
		SupplierID  string `json:"supplierId" binding:"required"`
		Ingredients []struct {
			IngredientID string  `json:"ingredientId" binding:"required"`
			Quantity     float64 `json:"quantity" binding:"required,gt=0"`
			Unit         string  `json:"unit" binding:"required"`
			UnitPrice    float64 `json:"unitPrice" binding:"required,gte=0"`
		} `json:"ingredients" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Log.Error("SaveOrderIngredientsWithSupplier bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate order exists
	var order models.Order
	if err := store.DB.GormClient.First(&order, "order_id = ?", orderID).Error; err != nil {
		logger.Log.Error("SaveOrderIngredientsWithSupplier order not found", "order_id", orderID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Validate supplier exists
	var supplier models.Supplier
	if err := store.DB.GormClient.First(&supplier, "supplier_id = ?", request.SupplierID).Error; err != nil {
		logger.Log.Error("SaveOrderIngredientsWithSupplier supplier not found", "supplier_id", request.SupplierID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
		return
	}

	// Validate ingredients exist
	for _, ing := range request.Ingredients {
		var ingredient models.Ingredient
		if err := store.DB.GormClient.First(&ingredient, "ingredient_id = ?", ing.IngredientID).Error; err != nil {
			logger.Log.Error("SaveOrderIngredientsWithSupplier ingredient not found", "ingredient_id", ing.IngredientID, "error", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Ingredient not found: " + ing.IngredientID})
			return
		}
	}

	// Start transaction
	tx := store.DB.GormClient.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find or create SupplierRequest for this order + supplier
	var supplierRequest models.SupplierRequest
	err := tx.Where("order_id = ? AND supplier_id = ?", orderID, request.SupplierID).First(&supplierRequest).Error
	if err != nil {
		// Create new supplier request
		supplierRequest = models.SupplierRequest{
			OrderID:    orderID,
			SupplierID: request.SupplierID,
			Status:     "Pending",
		}
		if err := tx.Create(&supplierRequest).Error; err != nil {
			logger.Log.Error("SaveOrderIngredientsWithSupplier create request error", "error", err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		logger.Log.Info("SaveOrderIngredientsWithSupplier created new supplier request", "request_id", supplierRequest.RequestID)
	} else {
		// Update status if needed (keep existing status if not explicitly set)
		logger.Log.Info("SaveOrderIngredientsWithSupplier found existing supplier request", "request_id", supplierRequest.RequestID)
	}

	// Delete existing request details for this request
	if err := tx.Where("request_id = ?", supplierRequest.RequestID).Delete(&models.SupplierRequestDetail{}).Error; err != nil {
		logger.Log.Error("SaveOrderIngredientsWithSupplier delete existing details error", "error", err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create new request details
	for _, ing := range request.Ingredients {
		detail := models.SupplierRequestDetail{
			RequestID:    supplierRequest.RequestID,
			IngredientID: ing.IngredientID,
			Quantity:     ing.Quantity,
			Unit:         ing.Unit,
			UnitPrice:    ing.UnitPrice,
			// TotalPrice is a generated column in the database (quantity * unit_price)
		}

		if err := tx.Create(&detail).Error; err != nil {
			logger.Log.Error("SaveOrderIngredientsWithSupplier create detail error", "error", err, "ingredient_id", ing.IngredientID)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		logger.Log.Error("SaveOrderIngredientsWithSupplier commit error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload supplier request with relations
	if err := store.DB.GormClient.
		Preload("Order").
		Preload("Supplier").
		Preload("Details.Ingredient").
		First(&supplierRequest, "request_id = ?", supplierRequest.RequestID).Error; err != nil {
		logger.Log.Error("SaveOrderIngredientsWithSupplier reload error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message":      "Ingredients saved with supplier successfully",
		"requestId":    supplierRequest.RequestID,
		"orderId":      supplierRequest.OrderID,
		"supplierId":   supplierRequest.SupplierID,
		"status":       supplierRequest.Status,
		"detailsCount": len(supplierRequest.Details),
	})
}

// convertOrderToDTO maps model to DTO
func convertOrderToDTO(o *models.Order, includeChildren bool) models.OrderDTO {
	dto := models.OrderDTO{
		OrderID:         o.OrderID,
		KitchenID:       o.KitchenID,
		OrderDate:       o.OrderDate,
		Note:            o.Note,
		Status:          o.Status,
		CreatedByUserID: o.CreatedByUserID,
		CreatedDate:     o.CreatedDate,
		ModifiedDate:    o.ModifiedDate,
	}
	if o.Kitchen != nil {
		dto.KitchenName = o.Kitchen.KitchenName
	}
	if o.CreatedBy != nil {
		dto.CreatedByName = o.CreatedBy.FullName
	}
	if includeChildren {
		if len(o.Details) > 0 {
			dto.Details = make([]models.OrderDetailDTO, len(o.Details))
			for i, d := range o.Details {
				dto.Details[i] = models.OrderDetailDTO{
					OrderDetailID: d.OrderDetailID,
					DishID:        d.DishID,
					Portions:      d.Portions,
					Note:          d.Note,
				}
				if d.Dish != nil {
					dto.Details[i].DishName = d.Dish.DishName
				}
				if len(d.Ingredients) > 0 {
					dto.Details[i].Ingredients = make([]models.OrderIngredientDTO, len(d.Ingredients))
					for j, ing := range d.Ingredients {
						dto.Details[i].Ingredients[j] = models.OrderIngredientDTO{
							OrderIngredientID:  ing.OrderIngredientID,
							IngredientID:       ing.IngredientID,
							Quantity:           ing.Quantity,
							Unit:               ing.Unit,
							StandardPerPortion: ing.StandardPerPortion,
						}
						if ing.Ingredient != nil {
							dto.Details[i].Ingredients[j].IngredientName = ing.Ingredient.IngredientName
						}
					}
				}
			}
		}
		if len(o.SupplementaryFoods) > 0 {
			dto.Supplementaries = make([]models.OrderSupplementaryDTO, len(o.SupplementaryFoods))
			for i, s := range o.SupplementaryFoods {
				dto.Supplementaries[i] = models.OrderSupplementaryDTO{
					SupplementaryID:    s.SupplementaryID,
					IngredientID:       s.IngredientID,
					Quantity:           s.Quantity,
					Unit:               s.Unit,
					StandardPerPortion: s.StandardPerPortion,
					Portions:           s.Portions,
					Note:               s.Note,
				}
				if s.Ingredient != nil {
					dto.Supplementaries[i].IngredientName = s.Ingredient.IngredientName
				}
			}
		}
	}
	return dto
}
