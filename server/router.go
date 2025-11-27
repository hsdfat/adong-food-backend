package server

import (
	"adong-be/handler"
	"adong-be/logger"
	"adong-be/store"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hsdfat/go-auth-middleware/ginauth"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Create user provider

	// Create enhanced token storage
	// tokenStorage := core.NewInMemoryTokenStorage()

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
		TokenStorage: store.DB,
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
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Request logging middleware with user identity
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		userIDAfter, _ := c.Get("identity")
		if len(c.Errors) > 0 {
			logger.Log.Error("handler returned error",
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"status", status,
				"errors", c.Errors.String(),
				"latency", latency.String(),
				"user_id", userIDAfter,
			)
		} else if status >= 400 {
			logger.Log.Error("request completed with error status",
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"status", status,
				"latency", latency.String(),
				"user_id", userIDAfter,
			)
		} else {
			logger.Log.Info("request completed",
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"status", status,
				"latency", latency.String(),
				"user_id", userIDAfter,
			)
		}
	})

	// API routes
	api := r.Group("/api")
	api.Use(authMiddleware.MiddlewareFunc())
	{
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

		api.GET("/recipe-standards", handler.GetRecipeStandards)
		api.GET("/recipe-standards/:id", handler.GetRecipeStandard)
		api.POST("/recipe-standards", handler.CreateRecipeStandard)
		api.PUT("/recipe-standards/:id", handler.UpdateRecipeStandard)
		api.DELETE("/recipe-standards/:id", handler.DeleteRecipeStandard)
		api.GET("/recipe-standards/dish/:dishId", handler.GetRecipeStandardsByDish)

		api.GET("/supplier-prices", handler.GetSupplierPrices)
		api.GET("/supplier-prices/ingredient/:ingredientId", handler.GetSupplierPricesByIngredient)
		api.GET("/supplier-prices/supplier/:supplierId", handler.GetSupplierPricesBySupplier)
		api.GET("/supplier-prices/:id", handler.GetSupplierPrice)
		api.POST("/supplier-prices", handler.CreateSupplierPrice)
		api.PUT("/supplier-prices/:id", handler.UpdateSupplierPrice)
		api.DELETE("/supplier-prices/:id", handler.DeleteSupplierPrice)

		api.GET("/orders", handler.GetOrders)
		api.GET("/orders/:id", handler.GetOrder)
		api.GET("/orders/:id/ingredients/summary", handler.GetOrderIngredientsSummary)
		api.GET("/orders/:id/ingredients/:ingredientId/summary", handler.GetOrderIngredientSummary)
		api.GET("/orders/:id/selected-suppliers", handler.GetOrderSelectedSuppliers)
		api.POST("/orders", handler.CreateOrder)
		api.POST("/orders/:id/supplier-requests", handler.SaveOrderIngredientsWithSupplier)
		api.PATCH("/orders/:id/status", handler.UpdateOrderStatus)
		api.DELETE("/orders/:id", handler.DeleteOrder)

		// Best supplier selection - returns data to frontend only
		api.GET("/orders/:id/best-suppliers", handler.GetBestSuppliersForOrder)
		api.POST("/orders/best-suppliers", handler.GetBestSuppliersForIngredients)

		// Initialize inventory handlers
		stockHandler := handler.NewInventoryStockHandler(store.DB.GormClient)
		importHandler := handler.NewInventoryImportHandler(store.DB.GormClient)
		exportHandler := handler.NewInventoryExportHandler(store.DB.GormClient)
		adjustmentHandler := handler.NewInventoryAdjustmentHandler(store.DB.GormClient)
		requestHandler := handler.NewIngredientRequestHandler(store.DB.GormClient)
		reportsHandler := handler.NewInventoryReportsHandler(store.DB.GormClient)

		// Inventory routes group
		inventory := api.Group("/inventory")
		{
			// Stock management
			stock := inventory.Group("/stocks")
			{
				stock.GET("", stockHandler.GetAllStocks)                         // GET /api/inventory/stocks?kitchen_id=K001&page=1&limit=50
				stock.GET("/:id", stockHandler.GetStockByID)                     // GET /api/inventory/stocks/1
				stock.GET("/query", stockHandler.GetStockByKitchenAndIngredient) // GET /api/inventory/stocks/query?kitchen_id=K001&ingredient_id=NL001
				stock.PUT("/:id/levels", stockHandler.UpdateStockLevels)         // PUT /api/inventory/stocks/1/levels
				stock.GET("/alerts/low", stockHandler.GetLowStockAlerts)         // GET /api/inventory/stocks/alerts/low?kitchen_id=K001
				stock.GET("/transactions", stockHandler.GetStockTransactions)    // GET /api/inventory/stocks/transactions?kitchen_id=K001&ingredient_id=NL001
				stock.GET("/summary", stockHandler.GetStockSummary)              // GET /api/inventory/stocks/summary?kitchen_id=K001
				stock.GET("/valuation", stockHandler.GetStockValuation)          // GET /api/inventory/stocks/valuation?kitchen_id=K001
			}

			// Import management
			imports := inventory.Group("/imports")
			{
				imports.GET("", importHandler.GetAllImports)                                // GET /api/inventory/imports?kitchen_id=K001&status=draft
				imports.GET("/:id", importHandler.GetImportByID)                            // GET /api/inventory/imports/IM20240520-12345
				imports.POST("", importHandler.CreateImport)                                // POST /api/inventory/imports
				imports.POST("/from-request/:requestId", importHandler.CreateImportFromRequest) // POST /api/inventory/imports/from-request/RQ20240520-12345
				imports.PUT("/:id", importHandler.UpdateImport)                             // PUT /api/inventory/imports/IM20240520-12345
				imports.POST("/:id/approve", importHandler.ApproveImport)                   // POST /api/inventory/imports/IM20240520-12345/approve
				imports.DELETE("/:id", importHandler.DeleteImport)                          // DELETE /api/inventory/imports/IM20240520-12345
			}

			// Export management
			exports := inventory.Group("/exports")
			{
				exports.GET("", exportHandler.GetAllExports)              // GET /api/inventory/exports?kitchen_id=K001&export_type=production
				exports.GET("/:id", exportHandler.GetExportByID)          // GET /api/inventory/exports/EX20240520-12345
				exports.POST("", exportHandler.CreateExport)              // POST /api/inventory/exports
				exports.PUT("/:id", exportHandler.UpdateExport)           // PUT /api/inventory/exports/EX20240520-12345
				exports.POST("/:id/approve", exportHandler.ApproveExport) // POST /api/inventory/exports/EX20240520-12345/approve
				exports.DELETE("/:id", exportHandler.DeleteExport)        // DELETE /api/inventory/exports/EX20240520-12345
			}

			// Adjustment management
			adjustments := inventory.Group("/adjustments")
			{
				adjustments.GET("", adjustmentHandler.GetAllAdjustments)                  // GET /api/inventory/adjustments?kitchen_id=K001&adjustment_type=count
				adjustments.GET("/:id", adjustmentHandler.GetAdjustmentByID)              // GET /api/inventory/adjustments/ADJ20240520-12345
				adjustments.POST("", adjustmentHandler.CreateAdjustment)                  // POST /api/inventory/adjustments
				adjustments.PUT("/:id", adjustmentHandler.UpdateAdjustment)               // PUT /api/inventory/adjustments/ADJ20240520-12345
				adjustments.POST("/:id/approve", adjustmentHandler.ApproveAdjustment)     // POST /api/inventory/adjustments/ADJ20240520-12345/approve
				adjustments.DELETE("/:id", adjustmentHandler.DeleteAdjustment)            // DELETE /api/inventory/adjustments/ADJ20240520-12345
			}

			// Ingredient Request management
			requests := inventory.Group("/requests")
			{
				requests.GET("", requestHandler.GetAllRequests)                           // GET /api/inventory/requests?kitchen_id=K001&status=pending
				requests.GET("/:id", requestHandler.GetRequestByID)                       // GET /api/inventory/requests/RQ20240520-12345
				requests.POST("", requestHandler.CreateRequest)                           // POST /api/inventory/requests
				requests.POST("/from-order/:orderId", requestHandler.CreateRequestFromOrder) // POST /api/inventory/requests/from-order/OR001
				requests.PUT("/:id", requestHandler.UpdateRequest)                        // PUT /api/inventory/requests/RQ20240520-12345
				requests.POST("/:id/approve", requestHandler.ApproveRequest)              // POST /api/inventory/requests/RQ20240520-12345/approve
				requests.DELETE("/:id", requestHandler.DeleteRequest)                     // DELETE /api/inventory/requests/RQ20240520-12345
			}

			// Inventory Reports
			reports := inventory.Group("/reports")
			{
				reports.GET("/stock-movement", reportsHandler.GetStockMovementReport)         // GET /api/inventory/reports/stock-movement?kitchen_id=K001&from_date=2024-01-01&to_date=2024-01-31
				reports.GET("/expiry-alerts", reportsHandler.GetExpiryAlerts)                 // GET /api/inventory/reports/expiry-alerts?kitchen_id=K001&days_ahead=30
				reports.GET("/stock-value-trend", reportsHandler.GetStockValueTrend)          // GET /api/inventory/reports/stock-value-trend?kitchen_id=K001&from_date=2024-01-01&to_date=2024-01-31
				reports.GET("/transaction-summary", reportsHandler.GetTransactionSummary)     // GET /api/inventory/reports/transaction-summary?kitchen_id=K001&from_date=2024-01-01&to_date=2024-01-31
				reports.GET("/top-consumed", reportsHandler.GetTopConsumedIngredients)        // GET /api/inventory/reports/top-consumed?kitchen_id=K001&from_date=2024-01-01&to_date=2024-01-31&limit=10
			}
		}
	}
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}
