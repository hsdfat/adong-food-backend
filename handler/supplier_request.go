package handler

import (
	"adong-be/logger"
	"adong-be/models"
	"adong-be/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSupplierRequestsByOrder returns all supplier requests for a given order id
func GetSupplierRequestsByOrder(c *gin.Context) {
    orderID := c.Param("id")
    logger.Log.Info("GetSupplierRequestsByOrder called", "order_id", orderID)

    var requests []models.SupplierRequest
    if err := store.DB.GormClient.
        Preload("Order").
        Preload("Supplier").
        Preload("Details.Ingredient").
        Where("order_id = ?", orderID).
        Find(&requests).Error; err != nil {
        logger.Log.Error("GetSupplierRequestsByOrder db error", "order_id", orderID, "error", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, requests)
}

// GetSupplierRequestByOrderAndIngredient returns request details for a given order and ingredient
func GetSupplierRequestByOrderAndIngredient(c *gin.Context) {
    orderID := c.Param("id")
    ingredientID := c.Param("ingredientId")
    logger.Log.Info("GetSupplierRequestByOrderAndIngredient called", "order_id", orderID, "ingredient_id", ingredientID)

    // Join details with requests to filter by order
    var details []models.SupplierRequestDetail
    db := store.DB.GormClient.Model(&models.SupplierRequestDetail{})
    db = db.Joins("JOIN supplier_requests sr ON sr.request_id = supplier_request_details.request_id")
    db = db.Where("sr.order_id = ? AND supplier_request_details.ingredient_id = ?", orderID, ingredientID)
    db = db.Preload("Ingredient").
        Preload("Request").
        Preload("Request.Supplier")

    if err := db.Find(&details).Error; err != nil {
        logger.Log.Error("GetSupplierRequestByOrderAndIngredient db error", "order_id", orderID, "ingredient_id", ingredientID, "error", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, details)
}


