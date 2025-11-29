package utils

import (
	"adong-be/models"
	"adong-be/store"
	"errors"

	"github.com/gin-gonic/gin"
)

// UserKitchenScope holds authorization info for the current user
type UserKitchenScope struct {
	User        *models.User
	IsAdmin     bool
	KitchenIDs  []string
}

// GetUserKitchenScope loads the current user and the list of kitchens they can access.
// Admin users have IsAdmin=true and KitchenIDs left empty (meaning "all kitchens").
func GetUserKitchenScope(c *gin.Context) (*UserKitchenScope, error) {
	identity, ok := c.Get("identity")
	if !ok {
		return nil, errors.New("missing identity in context")
	}

	userID, ok := identity.(string)
	if !ok || userID == "" {
		return nil, errors.New("invalid identity in context")
	}

	var user models.User
	if err := store.DB.GormClient.First(&user, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}

	scope := &UserKitchenScope{
		User: &user,
	}

	if user.Role == "Admin" {
		scope.IsAdmin = true
		return scope, nil
	}

	// Load kitchen IDs from the join table user_kitchens
	type userKitchenRow struct {
		KitchenID string `gorm:"column:kitchen_id"`
	}
	var rows []userKitchenRow
	if err := store.DB.GormClient.
		Table("user_kitchens").
		Where("user_id = ?", userID).
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	for _, r := range rows {
		if r.KitchenID != "" {
			scope.KitchenIDs = append(scope.KitchenIDs, r.KitchenID)
		}
	}

	return scope, nil
}


