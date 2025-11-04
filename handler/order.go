package handler

import (
	"adong-be/models"
	"adong-be/store"
	"adong-be/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetOrders lists orders with filters: kitchen_id, status, date range, dish_id, ingredient_id
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

	kitchenID := c.Query("kitchen_id")
	status := c.Query("status")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	dishID := c.Query("dish_id")
	ingredientID := c.Query("ingredient_id")

	var total int64
	var orders []models.Order

	db := store.DB.GormClient.Model(&models.Order{})

	// Filters
	if params.Search != "" {
		db = db.Where("note ILIKE ? OR CAST(order_id AS TEXT) ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}
	if kitchenID != "" {
		db = db.Where("kitchen_id = ?", kitchenID)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if fromDate != "" {
		db = db.Where("order_date >= ?", fromDate)
	}
	if toDate != "" {
		if t, err := time.Parse("2006-01-02", toDate); err == nil {
			db = db.Where("order_date < ?", t.Add(24*time.Hour))
		} else {
			db = db.Where("order_date <= ?", toDate)
		}
	}
	if dishID != "" {
		db = db.Joins("JOIN order_details od ON od.order_id = orders.order_id").Where("od.dish_id = ?", dishID)
	}
	if ingredientID != "" {
		db = db.Joins("JOIN order_details od2 ON od2.order_id = orders.order_id").
			Joins("JOIN order_ingredients oi ON oi.order_detail_id = od2.order_detail_id").
			Where("oi.ingredient_id = ?", ingredientID)
	}

	// Count distinct orders
	if err := db.Distinct("orders.order_id").Count(&total).Error; err != nil {
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
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSort)

	// Pagination
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	// Fetch and preload minimal relations
	if err := db.Preload("Kitchen").Preload("CreatedBy").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map to DTOs
	dtos := make([]models.OrderDTO, len(orders))
	for i := range orders {
		dtos[i] = convertOrderToDTO(&orders[i], false)
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{Data: dtos, Meta: meta})
}

// GetOrder returns a single order with full details
func GetOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := store.DB.GormClient.
		Preload("Kitchen").
		Preload("CreatedBy").
		Preload("Details.Dish").
		Preload("Details.Ingredients.Ingredient").
		Preload("SupplementaryFoods.Ingredient").
		First(&order, "order_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	dto := convertOrderToDTO(&order, true)
	c.JSON(http.StatusOK, dto)
}

// CreateOrder creates a new order with nested details/ingredients/supplementary foods
func CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if identity, ok := c.Get("userID"); ok {
		if v, ok2 := identity.(string); ok2 {
			order.CreatedByUserID = v
		}
	}

	tx := store.DB.GormClient.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create details and nested ingredients
	for i := range order.Details {
		order.Details[i].OrderID = order.OrderID
		if err := tx.Create(&order.Details[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for j := range order.Details[i].Ingredients {
			order.Details[i].Ingredients[j].OrderDetailID = order.Details[i].OrderDetailID
			if err := tx.Create(&order.Details[i].Ingredients[j]).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	// Create supplementary foods
	for i := range order.SupplementaryFoods {
		order.SupplementaryFoods[i].OrderID = order.OrderID
		if err := tx.Create(&order.SupplementaryFoods[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
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

// DeleteOrder deletes an order by id (cascade removes children)
func DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.Order{}, "order_id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
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
