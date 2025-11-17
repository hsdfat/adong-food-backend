package handler

import (
	"adong-be/logger"
	"adong-be/models"
	"adong-be/store"
	"adong-be/utils"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetOrders(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("GetOrders called", "user_id", uid)

	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Log.Error("GetOrders bind query error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params = models.GetPaginationParams(params.Page, params.PageSize, params.Search, params.SortBy, params.SortDir)

	kitchenID := c.Query("kitchen_id")
	status := c.Query("status")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	dishID := c.Query("dish_id")
	ingredientID := c.Query("ingredient_id")

	var orders []models.Order
	var total int64

	dataDB := store.DB.GormClient.Model(&models.Order{})
	countDB := store.DB.GormClient.Model(&models.Order{})

	if dishID != "" || ingredientID != "" {
		countDB = countDB.Distinct("orders.order_id")
		dataDB = dataDB.Distinct("orders.order_id")
	}

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

	if err := countDB.Distinct("orders.order_id").Count(&total).Error; err != nil {
		logger.Log.Error("GetOrders count error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	allowedSort := map[string]string{
		"order_id":     "orders.order_id",
		"order_date":   "orders.order_date",
		"status":       "orders.status",
		"created_date": "orders.created_date",
	}
	dataDB = utils.ApplySort(dataDB, params.SortBy, params.SortDir, allowedSort)
	dataDB = utils.ApplyPagination(dataDB, params.Page, params.PageSize)

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

	dtos := make([]models.OrderDTO, len(orders))
	for i := range orders {
		dtos[i] = convertOrderToDTO(&orders[i], true)
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{Data: dtos, Meta: meta})
}

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

func CreateOrder(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("CreateOrder called", "user_id", uid)
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		logger.Log.Error("CreateOrder bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if identity, ok := c.Get("identity"); ok {
		if v, ok2 := identity.(string); ok2 {
			order.CreatedByUserID = v
		}
	}

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

	details := order.Details
	supplementaryFoods := order.SupplementaryFoods
	order.Details = nil
	order.SupplementaryFoods = nil

	if err := tx.Create(&order).Error; err != nil {
		logger.Log.Error("CreateOrder create header error", "error", err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := range details {
		details[i].OrderID = order.OrderID
		details[i].OrderDetailID = 0

		ingredients := details[i].Ingredients
		details[i].Ingredients = nil

		if err := tx.Create(&details[i]).Error; err != nil {
			logger.Log.Error("CreateOrder create detail error", "error", err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for j := range ingredients {
			ingredients[j].OrderDetailID = details[i].OrderDetailID
			ingredients[j].OrderIngredientID = 0

			if ingredients[j].Quantity <= 0 {
				if ingredients[j].StandardPerPortion > 0 && details[i].Portions > 0 {
					ingredients[j].Quantity = ingredients[j].StandardPerPortion * float64(details[i].Portions)
				} else {
					logger.Log.Warn("CreateOrder skipping ingredient with invalid quantity", "ingredient_id", ingredients[j].IngredientID)
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

	for i := range supplementaryFoods {
		supplementaryFoods[i].OrderID = order.OrderID
		supplementaryFoods[i].SupplementaryID = 0

		if supplementaryFoods[i].Quantity <= 0 {
			if supplementaryFoods[i].StandardPerPortion > 0 && supplementaryFoods[i].Portions > 0 {
				supplementaryFoods[i].Quantity = supplementaryFoods[i].StandardPerPortion * float64(supplementaryFoods[i].Portions)
			} else {
				logger.Log.Warn("CreateOrder skipping supplementary with invalid quantity", "ingredient_id", supplementaryFoods[i].IngredientID)
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

func UpdateOrderStatus(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("UpdateOrderStatus called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("UpdateOrderStatus bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.Order
	if err := store.DB.GormClient.First(&order, "order_id = ?", id).Error; err != nil {
		logger.Log.Error("UpdateOrderStatus not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	order.Status = req.Status
	if err := store.DB.GormClient.Save(&order).Error; err != nil {
		logger.Log.Error("UpdateOrderStatus db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully", "status": order.Status})
}

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

type IngredientTotal struct {
	IngredientID   string  `json:"ingredientId"`
	IngredientName string  `json:"ingredientName"`
	Unit           string  `json:"unit"`
	TotalQuantity  float64 `json:"totalQuantity"`
}

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

type SupplierOption struct {
	ProductID      int     `json:"productId"`
	ProductName    string  `json:"productName"`
	SupplierID     string  `json:"supplierId"`
	SupplierName   string  `json:"supplierName"`
	UnitPrice      float64 `json:"unitPrice"`
	Unit           string  `json:"unit"`
	Specification  string  `json:"specification"`
	IsFavorite     bool    `json:"isFavorite"`
	IsLowestPrice  bool    `json:"isLowestPrice"`
	TotalCost      float64 `json:"totalCost"`
}

type IngredientSuppliers struct {
	IngredientID   string           `json:"ingredientId"`
	IngredientName string           `json:"ingredientName"`
	TotalQuantity  float64          `json:"totalQuantity"`
	Unit           string           `json:"unit"`
	BestSupplier   *SupplierOption  `json:"bestSupplier"`
	AllSuppliers   []SupplierOption `json:"allSuppliers"`
}

// GetBestSuppliersForOrder returns best supplier recommendations for all ingredients
func GetBestSuppliersForOrder(c *gin.Context) {
	uid, _ := c.Get("identity")
	orderID := c.Param("id")
	logger.Log.Info("GetBestSuppliersForOrder called", "order_id", orderID, "user_id", uid)

	var order models.Order
	if err := store.DB.GormClient.First(&order, "order_id = ?", orderID).Error; err != nil {
		logger.Log.Error("GetBestSuppliersForOrder order not found", "order_id", orderID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var ingredients []IngredientTotal
	sql := `
        SELECT DISTINCT x.ingredient_id AS ingredient_id,
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
        GROUP BY x.ingredient_id, mi.ingredient_name, x.unit`

	if err := store.DB.GormClient.Raw(sql, orderID, orderID).Scan(&ingredients).Error; err != nil {
		logger.Log.Error("GetBestSuppliersForOrder ingredients query error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var favorites []models.KitchenFavoriteSupplier
	store.DB.GormClient.Where("kitchen_id = ?", order.KitchenID).Find(&favorites)
	
	favoriteMap := make(map[string]bool)
	for _, fav := range favorites {
		favoriteMap[fav.SupplierID] = true
	}

	var results []IngredientSuppliers

	for _, ing := range ingredients {
		var prices []models.SupplierPrice
		if err := store.DB.GormClient.
			Preload("Supplier").
			Where("ingredient_id = ? AND active = true", ing.IngredientID).
			Where("(effective_from IS NULL OR effective_from <= NOW())").
			Where("(effective_to IS NULL OR effective_to >= NOW())").
			Order("unit_price ASC").
			Find(&prices).Error; err != nil {
			logger.Log.Error("GetBestSuppliersForOrder prices query error", "ingredient_id", ing.IngredientID, "error", err)
			continue
		}

		if len(prices) == 0 {
			logger.Log.Warn("GetBestSuppliersForOrder no prices found", "ingredient_id", ing.IngredientID)
			results = append(results, IngredientSuppliers{
				IngredientID:   ing.IngredientID,
				IngredientName: ing.IngredientName,
				TotalQuantity:  ing.TotalQuantity,
				Unit:           ing.Unit,
				BestSupplier:   nil,
				AllSuppliers:   []SupplierOption{},
			})
			continue
		}

		lowestPrice := prices[0].UnitPrice

		var allSuppliers []SupplierOption
		var bestSupplier *SupplierOption
		var favoriteSuppliers []SupplierOption

		for _, price := range prices {
			isFavorite := favoriteMap[price.SupplierID]
			isLowestPrice := (price.UnitPrice == lowestPrice)
			totalCost := ing.TotalQuantity * price.UnitPrice

			option := SupplierOption{
				ProductID:     price.ProductID,
				ProductName:   price.ProductName,
				SupplierID:    price.SupplierID,
				UnitPrice:     price.UnitPrice,
				Unit:          price.Unit,
				Specification: price.Specification,
				IsFavorite:    isFavorite,
				IsLowestPrice: isLowestPrice,
				TotalCost:     totalCost,
			}

			if price.Supplier != nil {
				option.SupplierName = price.Supplier.SupplierName
			}

			allSuppliers = append(allSuppliers, option)

			if isFavorite {
				favoriteSuppliers = append(favoriteSuppliers, option)
			}
		}

		if len(favoriteSuppliers) > 0 {
			bestSupplier = &favoriteSuppliers[0]
		} else {
			bestSupplier = &allSuppliers[0]
		}

		results = append(results, IngredientSuppliers{
			IngredientID:   ing.IngredientID,
			IngredientName: ing.IngredientName,
			TotalQuantity:  ing.TotalQuantity,
			Unit:           ing.Unit,
			BestSupplier:   bestSupplier,
			AllSuppliers:   allSuppliers,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"orderId":     orderID,
		"kitchenId":   order.KitchenID,
		"ingredients": results,
	})
}

func SaveOrderIngredientsWithSupplier(c *gin.Context) {
	uid, _ := c.Get("identity")
	orderID := c.Param("id")
	logger.Log.Info("SaveOrderIngredientsWithSupplier called", "order_id", orderID, "user_id", uid)

	var userID string
	if identity, ok := c.Get("identity"); ok {
		if v, ok2 := identity.(string); ok2 {
			userID = v
		}
	}

	var request struct {
		Selections []struct {
			IngredientID       string  `json:"ingredientId" binding:"required"`
			SelectedSupplierID string  `json:"selectedSupplierId" binding:"required"`
			SelectedProductID  int     `json:"selectedProductId" binding:"required"`
			Quantity           float64 `json:"quantity" binding:"required,gt=0"`
			Unit               string  `json:"unit" binding:"required"`
			UnitPrice          float64 `json:"unitPrice" binding:"required,gte=0"`
			Notes              string  `json:"notes"`
		} `json:"selections" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Log.Error("SaveOrderIngredientsWithSupplier bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.Order
	if err := store.DB.GormClient.First(&order, "order_id = ?", orderID).Error; err != nil {
		logger.Log.Error("SaveOrderIngredientsWithSupplier order not found", "order_id", orderID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	for i, sel := range request.Selections {
		var ingredient models.Ingredient
		if err := store.DB.GormClient.First(&ingredient, "ingredient_id = ?", sel.IngredientID).Error; err != nil {
			logger.Log.Error("SaveOrderIngredientsWithSupplier ingredient not found", "ingredient_id", sel.IngredientID, "error", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Ingredient not found: " + sel.IngredientID})
			return
		}

		var supplier models.Supplier
		if err := store.DB.GormClient.First(&supplier, "supplier_id = ?", sel.SelectedSupplierID).Error; err != nil {
			logger.Log.Error("SaveOrderIngredientsWithSupplier supplier not found", "supplier_id", sel.SelectedSupplierID, "error", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found: " + sel.SelectedSupplierID})
			return
		}

		var product models.SupplierPrice
		if err := store.DB.GormClient.First(&product, "product_id = ? AND supplier_id = ? AND ingredient_id = ?",
			sel.SelectedProductID, sel.SelectedSupplierID, sel.IngredientID).Error; err != nil {
			logger.Log.Error("SaveOrderIngredientsWithSupplier product mismatch",
				"product_id", sel.SelectedProductID,
				"supplier_id", sel.SelectedSupplierID,
				"ingredient_id", sel.IngredientID,
				"error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found or does not match supplier/ingredient"})
			return
		}

		var presentCount int64
		presentSQL := `
			SELECT COUNT(*) AS cnt FROM (
				SELECT 1
				FROM order_details od
				JOIN order_ingredients oi ON oi.order_detail_id = od.order_detail_id
				WHERE od.order_id = ? AND oi.ingredient_id = ?
				UNION ALL
				SELECT 1
				FROM order_supplementary_foods osf
				WHERE osf.order_id = ? AND osf.ingredient_id = ?
			) x`
		if err := store.DB.GormClient.Raw(presentSQL, orderID, sel.IngredientID, orderID, sel.IngredientID).Scan(&presentCount).Error; err != nil {
			logger.Log.Error("SaveOrderIngredientsWithSupplier validate ingredient error",
				"order_id", orderID, "ingredient_id", sel.IngredientID, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if presentCount == 0 {
			logger.Log.Error("SaveOrderIngredientsWithSupplier ingredient not in order",
				"order_id", orderID, "ingredient_id", sel.IngredientID)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ingredient does not belong to the order: " + sel.IngredientID})
			return
		}

		for j := i + 1; j < len(request.Selections); j++ {
			if request.Selections[j].IngredientID == sel.IngredientID {
				logger.Log.Error("SaveOrderIngredientsWithSupplier duplicate ingredient", "ingredient_id", sel.IngredientID)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate ingredient_id in request: " + sel.IngredientID})
				return
			}
		}
	}

	tx := store.DB.GormClient.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var savedSelections []models.OrderIngredientSupplier
	for _, sel := range request.Selections {
		totalCost := sel.Quantity * sel.UnitPrice

		var existing models.OrderIngredientSupplier
		findErr := tx.Where("order_id = ? AND ingredient_id = ?", orderID, sel.IngredientID).First(&existing).Error

		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			newSelection := models.OrderIngredientSupplier{
				OrderID:            orderID,
				IngredientID:       sel.IngredientID,
				SelectedSupplierID: sel.SelectedSupplierID,
				SelectedProductID:  sel.SelectedProductID,
				Quantity:           sel.Quantity,
				Unit:               sel.Unit,
				UnitPrice:          sel.UnitPrice,
				TotalCost:          totalCost,
				SelectedByUserID:   userID,
				Notes:              sel.Notes,
			}

			if err := tx.Create(&newSelection).Error; err != nil {
				logger.Log.Error("SaveOrderIngredientsWithSupplier create error", "error", err)
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			savedSelections = append(savedSelections, newSelection)
		} else if findErr != nil {
			logger.Log.Error("SaveOrderIngredientsWithSupplier find error", "error", findErr)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": findErr.Error()})
			return
		} else {
			existing.SelectedSupplierID = sel.SelectedSupplierID
			existing.SelectedProductID = sel.SelectedProductID
			existing.Quantity = sel.Quantity
			existing.Unit = sel.Unit
			existing.UnitPrice = sel.UnitPrice
			existing.TotalCost = totalCost
			existing.SelectedByUserID = userID
			existing.Notes = sel.Notes

			if err := tx.Save(&existing).Error; err != nil {
				logger.Log.Error("SaveOrderIngredientsWithSupplier update error", "error", err)
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			savedSelections = append(savedSelections, existing)
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.Log.Error("SaveOrderIngredientsWithSupplier commit error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responseSelections []models.OrderIngredientSupplier
	if err := store.DB.GormClient.
		Preload("Ingredient").
		Preload("SelectedSupplier").
		Preload("SelectedProduct").
		Preload("SelectedBy").
		Where("order_id = ?", orderID).
		Find(&responseSelections).Error; err != nil {
		logger.Log.Error("SaveOrderIngredientsWithSupplier reload error", "error", err)
		responseSelections = savedSelections
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Supplier selections saved successfully",
		"orderId":    orderID,
		"selections": responseSelections,
		"count":      len(savedSelections),
	})
}

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