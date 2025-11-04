package server

import (
	"adong-be/handler"
	"adong-be/store"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hsdfat/go-auth-middleware/core"
	"github.com/hsdfat/go-auth-middleware/ginauth"
)

func SetupRouter() *gin.Engine {
	// Initialize Gin router
	r := gin.Default()

	// Create user provider

	// Create enhanced token storage
	tokenStorage := core.NewInMemoryTokenStorage()

	// Create enhanced auth middleware
	authMiddleware := ginauth.NewEnhanced(ginauth.EnhancedAuthConfig{
		SecretKey:           "your-access-token-secret-key",
		RefreshSecretKey:    "your-refresh-token-secret-key", // Should be different
		AccessTokenTimeout:  24 * time.Hour,                  // Short-lived access tokens
		RefreshTokenTimeout: 7 * 24 * time.Hour,              // 7 days refresh tokens

		TokenLookup:   "header:Authorization,cookie:jwt",
		TokenHeadName: "Bearer",
		Realm:         "enhanced-auth",
		IdentityKey:   "identity",

		// Cookie configuration
		SendCookie:        true,
		CookieName:        "access_token",
		RefreshCookieName: "refresh_token",
		CookieHTTPOnly:    true,
		CookieSecure:      false, // Set to true in production with HTTPS
		CookieDomain:      "",

		// Storage and providers
		TokenStorage: tokenStorage,
		UserProvider: store.DB,

		// Authentication function
		Authenticator: ginauth.CreateEnhancedAuthenticator(store.DB),

		// Role-based authorization (example: only admin and user roles allowed)
		RoleAuthorizator: ginauth.CreateRoleAuthorizator("Admin", "user", "moderator"),

		// Security settings
		MaxConcurrentSessions: 5,         // Max 5 concurrent sessions per user
		SingleSessionMode:     false,     // Allow multiple sessions
		EnableTokenRevocation: true,      // Enable token revocation on logout
		CleanupInterval:       time.Hour, // Cleanup expired tokens every hour
	})

	// Public routes
	r.POST("/auth/login", authMiddleware.LoginHandler)
	r.POST("/auth/refresh", authMiddleware.RefreshHandler)
	authenticated := r.Group("/auth")
	authenticated.Use(authMiddleware.MiddlewareFunc())
	{
		authenticated.POST("/logout", authMiddleware.LogoutHandler)
		authenticated.POST("/logout-all", authMiddleware.LogoutAllHandler)
		authenticated.GET("/sessions", authMiddleware.GetUserSessionsHandler)
	}
	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API routes
	api := r.Group("/api")
	api.Use(authMiddleware.MiddlewareFunc())
	{
		// Master data routes
		api.GET("/ingredients", handler.GetIngredients)
		api.GET("/ingredients/:id", handler.GetIngredient)
		api.POST("/ingredients", handler.CreateIngredient)
		api.PUT("/ingredients/:id", handler.UpdateIngredient)
		api.DELETE("/ingredients/:id", handler.DeleteIngredient)

		api.GET("/kitchens", handler.GetKitchens)
		api.GET("/kitchens/:id", handler.GetKitchen)
		api.POST("/kitchens", handler.CreateKitchen)
		api.PUT("/kitchens/:id", handler.UpdateKitchen)
		api.DELETE("/kitchens/:id", handler.DeleteKitchen)

		api.GET("/users", handler.GetUsers)
		api.GET("/users/:id", handler.GetUser)
		api.POST("/users", handler.CreateUser)
		api.PUT("/users/:id", handler.UpdateUser)
		api.DELETE("/users/:id", handler.DeleteUser)

		api.GET("/dishes", handler.GetDishes)
		api.GET("/dishes/:id", handler.GetDish)
		api.POST("/dishes", handler.CreateDish)
		api.PUT("/dishes/:id", handler.UpdateDish)
		api.DELETE("/dishes/:id", handler.DeleteDish)

		api.GET("/suppliers", handler.GetSuppliers)
		api.GET("/suppliers/:id", handler.GetSupplier)
		api.POST("/suppliers", handler.CreateSupplier)
		api.PUT("/suppliers/:id", handler.UpdateSupplier)
		api.DELETE("/suppliers/:id", handler.DeleteSupplier)

		// Recipe standards
		api.GET("/recipe-standards", handler.GetRecipeStandards)
		api.GET("/recipe-standards/:id", handler.GetRecipeStandard)
		api.POST("/recipe-standards", handler.CreateRecipeStandard)
		api.PUT("/recipe-standards/:id", handler.UpdateRecipeStandard)
		api.DELETE("/recipe-standards/:id", handler.DeleteRecipeStandard)
		api.GET("/recipe-standards/dish/:dishId", handler.GetRecipeStandardsByDish)

		// Supplier price list
		api.GET("/supplier-prices", handler.GetSupplierPrices)
		api.GET("/supplier-prices/ingredient/:ingredientId", handler.GetSupplierPricesByIngredient) 
		api.GET("/supplier-prices/supplier/:supplierId", handler.GetSupplierPricesBySupplier)        
		api.GET("/supplier-prices/:id", handler.GetSupplierPrice)
		api.POST("/supplier-prices", handler.CreateSupplierPrice)
		api.PUT("/supplier-prices/:id", handler.UpdateSupplierPrice)
		api.DELETE("/supplier-prices/:id", handler.DeleteSupplierPrice)


		// Order forms
		api.GET("/orders", handler.GetOrders)
		api.GET("/orders/:id", handler.GetOrder)
		api.POST("/orders", handler.CreateOrder)
		api.PUT("/orders/:id", handler.UpdateOrder)
		api.DELETE("/orders/:id", handler.DeleteOrder)
		api.PATCH("/orders/:id/status", handler.UpdateOrderStatus)

		// Add dish with ingredients routes
		api.GET("/dishes/:id/with-ingredients", handler.GetDishWithIngredients)
		api.GET("/dishes/with-ingredients", handler.GetDishesWithIngredients)
		// // Order details
		// api.GET("/order-details", handler.GetOrderDetails)
		// api.GET("/order-details/:id", handler.GetOrderDetail)
		// api.POST("/order-details", CreateOrderDetail)
		// api.PUT("/order-details/:id", UpdateOrderDetail)
		// api.DELETE("/order-details/:id", DeleteOrderDetail)
		// api.GET("/order-details/order/:orderId", GetOrderDetailsByOrder)

		// // Ingredient requests
		// api.GET("/ingredient-requests", GetIngredientRequests)
		// api.GET("/ingredient-requests/:id", GetIngredientRequest)
		// api.POST("/ingredient-requests", CreateIngredientRequest)
		// api.PUT("/ingredient-requests/:id", UpdateIngredientRequest)
		// api.DELETE("/ingredient-requests/:id", DeleteIngredientRequest)

		// // Receiving documents
		// api.GET("/receiving-docs", GetReceivingDocs)
		// api.GET("/receiving-docs/:id", GetReceivingDoc)
		// api.POST("/receiving-docs", CreateReceivingDoc)
		// api.PUT("/receiving-docs/:id", UpdateReceivingDoc)
		// api.DELETE("/receiving-docs/:id", DeleteReceivingDoc)

		// // Receiving details
		// api.GET("/receiving-details", GetReceivingDetails)
		// api.GET("/receiving-details/:id", GetReceivingDetail)
		// api.POST("/receiving-details", CreateReceivingDetail)
		// api.PUT("/receiving-details/:id", UpdateReceivingDetail)
		// api.DELETE("/receiving-details/:id", DeleteReceivingDetail)

		// // Inventory
		// api.GET("/inventory", GetInventory)
		// api.GET("/inventory/:id", GetInventoryItem)
		// api.POST("/inventory", CreateInventoryItem)
		// api.PUT("/inventory/:id", UpdateInventoryItem)

		// // Accounts payable
		// api.GET("/payables", GetPayables)
		// api.GET("/payables/:id", GetPayable)
		// api.POST("/payables", CreatePayable)
		// api.PUT("/payables/:id", UpdatePayable)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}
