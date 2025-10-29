// File: handler/menu_card.go
package handler

import (
	"adong-be/models"
	"adong-be/store"
	"adong-be/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetMenuCards - Lấy danh sách phiếu thực đơn
func GetMenuCards(c *gin.Context) {
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

	var total int64
	countDB := store.DB.GormClient.Model(&models.MenuCard{})

	searchConfig := utils.SearchConfig{
		Fields: []string{"tenphieu", "phieuthucdonid", "ghichu"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var menuCards []models.MenuCard
	db := store.DB.GormClient.Model(&models.MenuCard{})
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"phieuthucdonid": "phieuthucdonid",
		"tenphieu":       "tenphieu",
		"ngaytao":        "ngaytao",
		"trangthai":      "trangthai",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	db = db.Preload("Kitchen").Preload("CreatedBy")

	if err := db.Find(&menuCards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to DTOs
	dtos := make([]models.MenuCardDTO, len(menuCards))
	for i, mc := range menuCards {
		dtos[i] = convertMenuCardToDTO(&mc)
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: dtos,
		Meta: meta,
	})
}

// GetMenuCard - Lấy chi tiết một phiếu thực đơn
func GetMenuCard(c *gin.Context) {
	id := c.Param("id")
	var menuCard models.MenuCard

	if err := store.DB.GormClient.
		Preload("Kitchen").
		Preload("CreatedBy").
		Preload("Details.Dish").
		Preload("Details.Ingredients.Ingredient").
		First(&menuCard, "phieuthucdonid = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu card not found"})
		return
	}

	dto := convertMenuCardToDTO(&menuCard)
	c.JSON(http.StatusOK, dto)
}

// CreateMenuCard - Tạo phiếu thực đơn mới
func CreateMenuCard(c *gin.Context) {
	var req models.MenuCardCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Lấy user ID từ context (giả sử đã có middleware auth)
	userID, exists := c.Get("userID")
	if !exists {
		userID = "USR001" // Default cho test
	}

	// Tạo MenuCard
	menuCard := models.MenuCard{
		MenuCardID:   generateMenuCardID(),
		MenuCardName: req.MenuCardName,
		KitchenID:    req.KitchenID,
		CreatedByID:  userID.(string),
		Status:       "DRAFT",
		Note:         req.Note,
	}

	if req.CreatedDate != nil && *req.CreatedDate != "" {
		t, err := time.Parse(time.RFC3339, *req.CreatedDate)
		if err == nil {
			menuCard.CreatedDate = &t
		}
	}

	// Bắt đầu transaction
	tx := store.DB.GormClient.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Tạo menu card
	if err := tx.Create(&menuCard).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Tạo chi tiết món ăn
	for _, detailReq := range req.Details {
		detail := models.MenuCardDetail{
			DetailID:   generateDetailID(),
			MenuCardID: menuCard.MenuCardID,
			DishID:     detailReq.DishID,
			Servings:   detailReq.Servings,
			Note:       detailReq.Note,
		}

		// Lấy tên món ăn
		var dish models.Dish
		if err := tx.First(&dish, "monanid = ?", detailReq.DishID).Error; err == nil {
			detail.DishName = dish.DishName
		}

		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Nếu có custom ingredients, thêm vào
		if len(detailReq.Ingredients) > 0 {
			for _, ingReq := range detailReq.Ingredients {
				ingredient := models.MenuCardDetailIngredient{
					ID:           generateIngredientID(),
					DetailID:     detail.DetailID,
					IngredientID: ingReq.IngredientID,
					Standard:     ingReq.Standard,
					Unit:         ingReq.Unit,
					Note:         ingReq.Note,
				}

				if err := tx.Create(&ingredient).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		} else {
			// Nếu không có custom, copy từ định mức chuẩn
			var recipeStandards []models.RecipeStandard
			if err := tx.Where("monanid = ?", detailReq.DishID).Find(&recipeStandards).Error; err == nil {
				for _, rs := range recipeStandards {
					ingredient := models.MenuCardDetailIngredient{
						ID:           generateIngredientID(),
						DetailID:     detail.DetailID,
						IngredientID: rs.IngredientID,
						Standard:     rs.StandardPer1 * float64(detailReq.Servings),
						Unit:         rs.Unit,
					}

					if err := tx.Create(&ingredient).Error; err != nil {
						tx.Rollback()
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
				}
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Load lại với đầy đủ thông tin
	store.DB.GormClient.
		Preload("Kitchen").
		Preload("CreatedBy").
		Preload("Details.Dish").
		Preload("Details.Ingredients.Ingredient").
		First(&menuCard, "phieuthucdonid = ?", menuCard.MenuCardID)

	dto := convertMenuCardToDTO(&menuCard)
	c.JSON(http.StatusCreated, dto)
}

// UpdateMenuCard - Cập nhật phiếu thực đơn
func UpdateMenuCard(c *gin.Context) {
	id := c.Param("id")
	var req models.MenuCardCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var menuCard models.MenuCard
	if err := store.DB.GormClient.First(&menuCard, "phieuthucdonid = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu card not found"})
		return
	}

	// Chỉ cho phép update nếu status là DRAFT
	if menuCard.Status != "DRAFT" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can only update draft menu cards"})
		return
	}

	tx := store.DB.GormClient.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update menu card info
	menuCard.MenuCardName = req.MenuCardName
	menuCard.KitchenID = req.KitchenID
	menuCard.Note = req.Note

	if req.CreatedDate != nil && *req.CreatedDate != "" {
		t, err := time.Parse(time.RFC3339, *req.CreatedDate)
		if err == nil {
			menuCard.CreatedDate = &t
		}
	}

	if err := tx.Save(&menuCard).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Xóa các details cũ
	if err := tx.Where("phieuthucdonid = ?", id).Delete(&models.MenuCardDetail{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Tạo lại details mới (giống CreateMenuCard)
	for _, detailReq := range req.Details {
		detail := models.MenuCardDetail{
			DetailID:   generateDetailID(),
			MenuCardID: menuCard.MenuCardID,
			DishID:     detailReq.DishID,
			Servings:   detailReq.Servings,
			Note:       detailReq.Note,
		}

		var dish models.Dish
		if err := tx.First(&dish, "monanid = ?", detailReq.DishID).Error; err == nil {
			detail.DishName = dish.DishName
		}

		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(detailReq.Ingredients) > 0 {
			for _, ingReq := range detailReq.Ingredients {
				ingredient := models.MenuCardDetailIngredient{
					ID:           generateIngredientID(),
					DetailID:     detail.DetailID,
					IngredientID: ingReq.IngredientID,
					Standard:     ingReq.Standard,
					Unit:         ingReq.Unit,
					Note:         ingReq.Note,
				}

				if err := tx.Create(&ingredient).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	store.DB.GormClient.
		Preload("Kitchen").
		Preload("CreatedBy").
		Preload("Details.Dish").
		Preload("Details.Ingredients.Ingredient").
		First(&menuCard, "phieuthucdonid = ?", menuCard.MenuCardID)

	dto := convertMenuCardToDTO(&menuCard)
	c.JSON(http.StatusOK, dto)
}

// DeleteMenuCard - Xóa phiếu thực đơn
func DeleteMenuCard(c *gin.Context) {
	id := c.Param("id")
	
	var menuCard models.MenuCard
	if err := store.DB.GormClient.First(&menuCard, "phieuthucdonid = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu card not found"})
		return
	}

	// Chỉ cho phép xóa nếu status là DRAFT
	if menuCard.Status != "DRAFT" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can only delete draft menu cards"})
		return
	}

	if err := store.DB.GormClient.Delete(&menuCard).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu card deleted successfully"})
}

// ApproveMenuCard - Duyệt phiếu thực đơn
func ApproveMenuCard(c *gin.Context) {
	id := c.Param("id")
	
	var menuCard models.MenuCard
	if err := store.DB.GormClient.First(&menuCard, "phieuthucdonid = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu card not found"})
		return
	}

	if menuCard.Status != "DRAFT" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Menu card is not in draft status"})
		return
	}

	menuCard.Status = "APPROVED"
	if err := store.DB.GormClient.Save(&menuCard).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu card approved successfully"})
}

// Helper functions
func generateMenuCardID() string {
	return fmt.Sprintf("MC%s", uuid.New().String()[:8])
}

func generateDetailID() string {
	return fmt.Sprintf("MCD%s", uuid.New().String()[:8])
}

func generateIngredientID() string {
	return fmt.Sprintf("MCI%s", uuid.New().String()[:8])
}

func convertMenuCardToDTO(mc *models.MenuCard) models.MenuCardDTO {
	dto := models.MenuCardDTO{
		MenuCardID:   mc.MenuCardID,
		MenuCardName: mc.MenuCardName,
		KitchenID:    mc.KitchenID,
		CreatedByID:  mc.CreatedByID,
		Status:       mc.Status,
		Note:         mc.Note,
		CreatedAt:    mc.CreatedAt.Format(time.RFC3339),
		ModifiedDate: mc.ModifiedDate.Format(time.RFC3339),
	}

	if mc.CreatedDate != nil {
		t := mc.CreatedDate.Format(time.RFC3339)
		dto.CreatedDate = &t
	}

	if mc.Kitchen != nil {
		dto.KitchenName = mc.Kitchen.KitchenName
	}

	if mc.CreatedBy != nil {
		dto.CreatedByName = mc.CreatedBy.FullName
	}

	if mc.Details != nil {
		dto.Details = make([]models.MenuCardDetailDTO, len(mc.Details))
		for i, detail := range mc.Details {
			dto.Details[i] = models.MenuCardDetailDTO{
				DetailID: detail.DetailID,
				DishID:   detail.DishID,
				DishName: detail.DishName,
				Servings: detail.Servings,
				Note:     detail.Note,
			}

			if detail.Ingredients != nil {
				dto.Details[i].Ingredients = make([]models.MenuCardDetailIngredientDTO, len(detail.Ingredients))
				for j, ing := range detail.Ingredients {
					dto.Details[i].Ingredients[j] = models.MenuCardDetailIngredientDTO{
						ID:           ing.ID,
						IngredientID: ing.IngredientID,
						Standard:     ing.Standard,
						Unit:         ing.Unit,
						Note:         ing.Note,
					}

					if ing.Ingredient != nil {
						dto.Details[i].Ingredients[j].IngredientName = ing.Ingredient.IngredientName
					}
				}
			}
		}
	}

	return dto
}