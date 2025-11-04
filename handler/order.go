// handler/order.go
package handler

import (
	"adong-be/logger"
	"adong-be/models"
	"adong-be/store"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetOrders - Get all orders with pagination and filters
func GetOrders(c *gin.Context) {
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
	page, pageSize := params.Page, params.PageSize
	search := c.DefaultQuery("search", "")
	kitchenID := c.DefaultQuery("kitchenId", "")
	status := c.DefaultQuery("status", "")
	fromDate := c.DefaultQuery("fromDate", "")
	toDate := c.DefaultQuery("toDate", "")
	// sortBy := c.DefaultQuery("sortBy", "ngay_len")
	// sortDir := c.DefaultQuery("sortDir", "desc")

	var orders []models.OrderForm
	var total int64

	db := store.DB.GormClient.Model(&models.OrderForm{})

	// Apply filters
	if search != "" {
		db = db.Where("phieu_len_don_id LIKE ? OR ghi_chu LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if kitchenID != "" {
		db = db.Where("bep_id = ?", kitchenID)
	}
	if status != "" {
		db = db.Where("trang_thai = ?", status)
	}
	if fromDate != "" {
		db = db.Where("ngay_len >= ?", fromDate)
	}
	if toDate != "" {
		db = db.Where("ngay_len <= ?", toDate)
	}

	// Count total
	db.Count(&total)

	// Sort
	// orderBy := sortBy + " " + sortDir
	// db = db.Order(orderBy)

	// Paginate
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get kitchen names
	for i := range orders {
		var kitchen models.Kitchen
		if err := store.DB.GormClient.First(&kitchen, "bepid = ?", orders[i].KitchenID).Error; err == nil {
			orders[i].KitchenID = kitchen.KitchenName
		}
	}

	totalPages := (int(total) + pageSize - 1) / pageSize
	c.JSON(http.StatusOK, gin.H{
		"data": orders,
		"pagination": gin.H{
			"page":       page,
			"pageSize":   pageSize,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}

// GetOrder - Get single order by ID with full details
func GetOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.OrderForm

	if err := store.DB.GormClient.
		Preload("Details.Ingredients").
		Preload("SupplementaryFoods").
		First(&order, "phieu_len_don_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Get kitchen name
	var kitchen models.Kitchen
	if err := store.DB.GormClient.First(&kitchen, "bep_id = ?", order.KitchenID).Error; err == nil {
		// Kitchen name can be added to response if needed
	}

	// Calculate total ingredients
	totalSummary := calculateTotalIngredients(order)

	c.JSON(http.StatusOK, gin.H{
		"orderFormId":         order.OrderFormID,
		"kitchenId":           order.KitchenID,
		"orderDate":           order.OrderDate,
		"note":                order.Note,
		"status":              order.Status,
		"details":             order.Details,
		"supplementaryFoods":  order.SupplementaryFoods,
		"totalIngredients":    totalSummary,
		"createdBy":           order.CreatedBy,
		"createdAt":           order.CreatedAt,
		"updatedAt":           order.UpdatedAt,
	})
}

// CreateOrder - Create new order
func CreateOrder(c *gin.Context) {
	var order models.OrderForm
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set created by from auth context
	if userID, exists := c.Get("identity"); exists {
		order.CreatedBy = userID.(string)
		logger.Log.Info("Create order", "identity", userID)
	}

	// Validate kitchen exists
	var kitchen models.Kitchen
	if err := store.DB.GormClient.First(&kitchen, "bepid = ?", order.KitchenID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kitchen not found"})
		return
	}

	// Start transaction
	tx := store.DB.GormClient.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create order
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create order details
	for i := range order.Details {
		order.Details[i].OrderFormID = order.OrderFormID
		order.Details[i].ID = 0 // Ensure ID is zero for new insert
		if err := tx.Create(&order.Details[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Create ingredients for this detail
		for j := range order.Details[i].Ingredients {
			order.Details[i].Ingredients[j].DetailID = order.Details[i].ID
			order.Details[i].Ingredients[j].ID = 0
			if err := tx.Create(&order.Details[i].Ingredients[j]).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	// Create supplementary foods
	for i := range order.SupplementaryFoods {
		order.SupplementaryFoods[i].OrderFormID = order.OrderFormID
		order.SupplementaryFoods[i].ID = 0
		if err := tx.Create(&order.SupplementaryFoods[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order form created successfully",
		"data": gin.H{
			"orderFormId": order.OrderFormID,
			"kitchenId":   order.KitchenID,
			"orderDate":   order.OrderDate,
			"status":      order.Status,
			"createdAt":   order.CreatedAt,
		},
	})
}

// UpdateOrder - Update existing order
func UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.OrderForm

	// Check if order exists
	if err := store.DB.GormClient.First(&order, "phieu_len_don_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Bind new data
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order.OrderFormID = id // Ensure ID doesn't change

	// Start transaction
	tx := store.DB.GormClient.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update order header
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete existing details and ingredients
	tx.Where("phieu_len_don_id = ?", id).Delete(&models.SupplementaryFood{})
	tx.Exec("DELETE FROM nguyen_lieu_chi_tiet WHERE chi_tiet_id IN (SELECT id FROM chi_tiet_phieu_len_don WHERE phieu_len_don_id = ?)", id)
	tx.Where("phieu_len_don_id = ?", id).Delete(&models.OrderFormDetail{})

	// Recreate details
	for i := range order.Details {
		order.Details[i].OrderFormID = order.OrderFormID
		order.Details[i].ID = 0 // Reset ID for new insert
		if err := tx.Create(&order.Details[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for j := range order.Details[i].Ingredients {
			order.Details[i].Ingredients[j].DetailID = order.Details[i].ID
			order.Details[i].Ingredients[j].ID = 0
			if err := tx.Create(&order.Details[i].Ingredients[j]).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	// Recreate supplementary foods
	for i := range order.SupplementaryFoods {
		order.SupplementaryFoods[i].OrderFormID = order.OrderFormID
		order.SupplementaryFoods[i].ID = 0
		if err := tx.Create(&order.SupplementaryFoods[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Order form updated successfully",
		"data":    order,
	})
}

// DeleteOrder - Delete order
func DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.OrderForm

	if err := store.DB.GormClient.First(&order, "phieu_len_don_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Delete with cascade (foreign keys with ON DELETE CASCADE will handle details)
	if err := store.DB.GormClient.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order form deleted successfully"})
}

// UpdateOrderStatus - Update order status
func UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var request struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.OrderForm
	if err := store.DB.GormClient.First(&order, "phieu_len_don_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	order.Status = request.Status
	order.UpdatedAt = time.Now()

	if err := store.DB.GormClient.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status updated successfully",
		"data": gin.H{
			"orderFormId": order.OrderFormID,
			"status":      order.Status,
			"updatedAt":   order.UpdatedAt,
		},
	})
}

// GetDishWithIngredients - Get single dish with ingredients
func GetDishWithIngredients(c *gin.Context) {
	id := c.Param("id")
	var dish models.Dish

	if err := store.DB.GormClient.First(&dish, "monanid = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dish not found"})
		return
	}

	// Get recipe standards for this dish
	var recipeStandards []models.RecipeStandard
	if err := store.DB.GormClient.Where("dish_id = ?", id).Find(&recipeStandards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Build ingredients list
	var ingredients []models.IngredientInRecipe
	for _, rs := range recipeStandards {
		var ingredient models.Ingredient
		if err := store.DB.GormClient.First(&ingredient, "nguyenlieu_id = ?", rs.IngredientID).Error; err == nil {
			ingredients = append(ingredients, models.IngredientInRecipe{
				IngredientID:       ingredient.IngredientID,
				IngredientName:     ingredient.IngredientName,
				Unit:               ingredient.Unit,
				StandardPerPortion: rs.StandardPer1,
			})
		}
	}

	result := models.DishWithIngredients{
		DishID:      dish.DishID,
		DishName:    dish.DishName,
		Ingredients: ingredients,
	}

	c.JSON(http.StatusOK, result)
}

// GetDishesWithIngredients - Get all dishes with ingredients
func GetDishesWithIngredients(c *gin.Context) {
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
	page, pageSize := params.Page, params.PageSize
	search := c.DefaultQuery("search", "")

	var dishes []models.Dish
	db := store.DB.GormClient.Model(&models.Dish{})

	if search != "" {
		db = db.Where("ten_mon_an LIKE ? OR monanid LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&dishes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var results []models.DishWithIngredients
	for _, dish := range dishes {
		// Get recipe standards for this dish
		var recipeStandards []models.RecipeStandard
		store.DB.GormClient.Where("dish_id = ?", dish.DishID).Find(&recipeStandards)

		// Build ingredients list
		var ingredients []models.IngredientInRecipe
		for _, rs := range recipeStandards {
			var ingredient models.Ingredient
			if err := store.DB.GormClient.First(&ingredient, "nguyenlieu_id = ?", rs.IngredientID).Error; err == nil {
				ingredients = append(ingredients, models.IngredientInRecipe{
					IngredientID:       ingredient.IngredientID,
					IngredientName:     ingredient.IngredientName,
					Unit:               ingredient.Unit,
					StandardPerPortion: rs.StandardPer1,
				})
			}
		}

		results = append(results, models.DishWithIngredients{
			DishID:      dish.DishID,
			DishName:    dish.DishName,
			Ingredients: ingredients,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": results,
		"pagination": gin.H{
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// Helper function to calculate total ingredients
func calculateTotalIngredients(order models.OrderForm) []models.TotalIngredientSummary {
	totals := make(map[string]*models.TotalIngredientSummary)

	// From dishes
	for _, detail := range order.Details {
		for _, ing := range detail.Ingredients {
			if existing, ok := totals[ing.IngredientID]; ok {
				existing.Quantity += ing.Quantity
			} else {
				totals[ing.IngredientID] = &models.TotalIngredientSummary{
					IngredientID:   ing.IngredientID,
					IngredientName: ing.IngredientName,
					Quantity:       ing.Quantity,
					Unit:           ing.Unit,
				}
			}
		}
	}

	// From supplementary foods
	for _, item := range order.SupplementaryFoods {
		if existing, ok := totals[item.IngredientID]; ok {
			existing.Quantity += item.Quantity
		} else {
			totals[item.IngredientID] = &models.TotalIngredientSummary{
				IngredientID:   item.IngredientID,
				IngredientName: item.IngredientName,
				Quantity:       item.Quantity,
				Unit:           item.Unit,
			}
		}
	}

	// Convert map to slice
	result := make([]models.TotalIngredientSummary, 0, len(totals))
	for _, v := range totals {
		result = append(result, *v)
	}

	return result
}