# ProjectDump Analysis

**Generated on:** 2025-11-09 15:53:26
**Project Path:** .

## Project Summary

- **Primary Language:** Python
- **Total Files:** 37
- **Processed Files:** 37
- **Project Size:** 159.23 KB

## Detected Technologies

### Python (100.0% confidence)
*Python programming language*

**Related files:**
- Dockerfile
- auth/ginauth.go
- cmd/main.go
- db/db.sql
- db/migrate_order_id_to_string.sql
- ... and 25 more files

### Go (100.0% confidence)
*Go programming language*

**Related files:**
- auth/ginauth.go
- cmd/main.go
- db/db.sql
- db/migrate_order_id_to_string.sql
- go.mod
- ... and 27 more files

### JavaScript (100.0% confidence)
*JavaScript runtime and ecosystem*

**Related files:**
- auth/ginauth.go
- cmd/main.go
- handler/dish.go
- handler/dish_test.go
- handler/ingredient.go
- ... and 21 more files

### TypeScript (100.0% confidence)
*TypeScript - JavaScript with static typing*

**Related files:**
- auth/ginauth.go
- db/db.sql
- db/migrate_order_id_to_string.sql
- go.mod
- go.sum
- ... and 16 more files

### Docker (100.0% confidence)
*Docker containerization platform*

**Related files:**
- Dockerfile
- db/db.sql
- db/migrate_order_id_to_string.sql
- handler/kitchen.go
- handler/order.go
- ... and 5 more files

### CSS (100.0% confidence)
*Cascading Style Sheets*

**Related files:**
- Dockerfile
- auth/ginauth.go
- cmd/main.go
- db/db.sql
- db/migrate_order_id_to_string.sql
- ... and 27 more files

### Java (100.0% confidence)
*Java programming language*

**Related files:**
- auth/ginauth.go
- cmd/main.go
- handler/dish.go
- handler/dish_test.go
- handler/ingredient.go
- ... and 23 more files

### Ruby (55.0% confidence)
*Ruby programming language*

**Related files:**
- .gitignore
- README.md
- db/db.sql
- db/samples.sql
- go.mod
- ... and 6 more files

### Rust (30.0% confidence)
*Rust systems programming language*

**Related files:**
- Dockerfile
- db/migrate_order_id_to_string.sql
- go.mod
- go.sum
- handler/order.go
- ... and 1 more files

### C (10.0% confidence)
*C programming language*

**Related files:**
- cmd/main.go
- utils/search.go

## Directory Structure

```
├── .env
├── .gitignore
├── Dockerfile
├── README.md
├── auth
│   └── ginauth.go
├── cmd
│   └── main.go
├── db
│   ├── db.sql
│   ├── migrate_order_id_to_string.sql
│   └── samples.sql
├── go.mod
├── go.sum
├── handler
│   ├── dish.go
│   ├── dish_test.go
│   ├── ingredient.go
│   ├── kitchen.go
│   ├── order.go
│   ├── recipe_standard.go
│   ├── supplier.go
│   ├── supplier_price.go
│   └── user.go
├── logger
│   └── logger.go
├── models
│   ├── common.go
│   ├── dish.go
│   ├── ingredient.go
│   ├── kitchen.go
│   ├── order.go
│   ├── order_dto.go
│   ├── pagination.go
│   ├── recipe_standard.go
│   ├── recipe_standard_dto.go
│   ├── supplier.go
│   ├── supplier_price.go
│   ├── supplier_price_dto.go
│   └── user.go
├── server
│   └── router.go
├── store
│   └── gorm.go
└── utils
    └── search.go
```

## Source Code

### auth/

#### auth/ginauth.go
*Language: Go | Size: 399 bytes*

```go
package auth

import (
	"time"

	"github.com/hsdfat/go-auth-middleware/core"
)

// Enhanced UserProvider interface
type UserProvider interface {
	GetUserByUsername(username string) (*core.User, error)
	GetUserByID(userID int) (*core.User, error)
	GetUserByEmail(email string) (*core.User, error)
	UpdateUserLastLogin(userID int, lastLogin time.Time) error
	IsUserActive(userID int) (bool, error)
}

```

### cmd/

#### cmd/main.go
*Language: Go | Size: 908 bytes*

```go
package main

import (
	"adong-be/server"
	"adong-be/store"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	var err error
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=adong password=adong123 dbname=adongfood port=5432 sslmode=disable"
	}

	store.DB.GormClient, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")
	s := server.SetupRouter() 
	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "18080"
	}

	log.Printf("Server starting on port %s", port)
	if err := s.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
```

### db/

#### db/db.sql
*Language: SQL | Size: 21008 bytes*

```sql
-- ============================================================================
-- SIMPLIFIED DATABASE SCHEMA - WITHOUT SUPPLIER_REQUESTS
-- ============================================================================
-- This version assumes:
-- 1. All supplier prices are pre-loaded in supplier_price_list
-- 2. You select suppliers directly from the price list
-- 3. No need for quote request workflow
-- ============================================================================

BEGIN;

-- ============================================================================
-- MASTER DATA TABLES (Unchanged)
-- ============================================================================

-- Ingredient Types Lookup Table
CREATE TABLE IF NOT EXISTS public.ingredient_types
(
    ingredient_type_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    ingredient_type_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default",
    active boolean NOT NULL DEFAULT true,
    created_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT ingredient_types_pkey PRIMARY KEY (ingredient_type_id),
    CONSTRAINT ingredient_types_name_key UNIQUE (ingredient_type_name)
);

CREATE TABLE IF NOT EXISTS public.master_dishes
(
    dish_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    dish_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    cooking_method character varying(100) COLLATE pg_catalog."default",
    category character varying(100) COLLATE pg_catalog."default",
    description text COLLATE pg_catalog."default",
    active boolean NOT NULL DEFAULT true,
    created_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT master_dishes_pkey PRIMARY KEY (dish_id),
    CONSTRAINT master_dishes_dish_name_key UNIQUE (dish_name)
);

CREATE TABLE IF NOT EXISTS public.master_ingredients
(
    ingredient_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    ingredient_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    ingredient_type_id character varying(50) COLLATE pg_catalog."default",
    properties character varying(100) COLLATE pg_catalog."default",
    material_group character varying(100) COLLATE pg_catalog."default",
    unit character varying(50) COLLATE pg_catalog."default" NOT NULL,
    created_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT master_ingredients_pkey PRIMARY KEY (ingredient_id),
    CONSTRAINT master_ingredients_ingredient_name_key UNIQUE (ingredient_name)
);

CREATE TABLE IF NOT EXISTS public.master_kitchens
(
    kitchen_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    kitchen_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    address text COLLATE pg_catalog."default",
    phone character varying(20) COLLATE pg_catalog."default",
    active boolean NOT NULL DEFAULT true,
    created_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT master_kitchens_pkey PRIMARY KEY (kitchen_id)
);

CREATE TABLE IF NOT EXISTS public.master_suppliers
(
    supplier_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    supplier_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    zalo_link text COLLATE pg_catalog."default",
    address text COLLATE pg_catalog."default",
    phone character varying(20) COLLATE pg_catalog."default",
    email character varying(255) COLLATE pg_catalog."default",
    active boolean NOT NULL DEFAULT true,
    created_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT master_suppliers_pkey PRIMARY KEY (supplier_id)
);

CREATE TABLE IF NOT EXISTS public.master_users
(
    user_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    user_name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    password character varying(255) COLLATE pg_catalog."default" NOT NULL,
    full_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    role character varying(50) COLLATE pg_catalog."default",
    kitchen_id character varying(50) COLLATE pg_catalog."default",
    email character varying(255) COLLATE pg_catalog."default",
    phone character varying(20) COLLATE pg_catalog."default",
    active boolean NOT NULL DEFAULT true,
    created_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT master_users_pkey PRIMARY KEY (user_id),
    CONSTRAINT master_users_user_name_key UNIQUE (user_name)
);

-- ============================================================================
-- RECIPE STANDARDS (Unchanged)
-- ============================================================================

CREATE TABLE IF NOT EXISTS public.dish_recipe_standards
(
    recipe_id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    dish_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    ingredient_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    unit character varying(50) COLLATE pg_catalog."default" NOT NULL,
    quantity_per_serving numeric(10, 4) NOT NULL,
    notes text COLLATE pg_catalog."default",
    cost numeric(15, 2),
    updated_by_user_id character varying(50) COLLATE pg_catalog."default",
    created_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT dish_recipe_standards_pkey PRIMARY KEY (recipe_id)
);

-- ============================================================================
-- ORDER TABLES (Unchanged)
-- ============================================================================

CREATE TABLE IF NOT EXISTS public.orders
(
    order_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    kitchen_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    order_date date NOT NULL,
    note text COLLATE pg_catalog."default",
    status character varying(50) COLLATE pg_catalog."default" NOT NULL DEFAULT 'Pending'::character varying,
    created_by_user_id character varying(50) COLLATE pg_catalog."default",
    created_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT orders_pkey PRIMARY KEY (order_id)
);

CREATE TABLE IF NOT EXISTS public.order_details
(
    order_detail_id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    order_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    dish_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    portions integer NOT NULL,
    note text COLLATE pg_catalog."default",
    created_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT order_details_pkey PRIMARY KEY (order_detail_id)
);

CREATE TABLE IF NOT EXISTS public.order_ingredients
(
    order_ingredient_id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    order_detail_id integer NOT NULL,
    ingredient_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    quantity numeric(15, 4) NOT NULL,
    unit character varying(50) COLLATE pg_catalog."default" NOT NULL,
    standard_per_portion numeric(10, 4),
    created_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT order_ingredients_pkey PRIMARY KEY (order_ingredient_id)
);

CREATE TABLE IF NOT EXISTS public.order_supplementary_foods
(
    supplementary_id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    order_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    ingredient_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    quantity numeric(15, 4) NOT NULL,
    unit character varying(50) COLLATE pg_catalog."default" NOT NULL,
    standard_per_portion numeric(10, 4),
    portions integer,
    note text COLLATE pg_catalog."default",
    created_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT order_supplementary_foods_pkey PRIMARY KEY (supplementary_id)
);

-- ============================================================================
-- SUPPLIER PRICE LIST (Unchanged)
-- ============================================================================

CREATE TABLE IF NOT EXISTS public.supplier_price_list
(
    product_id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    product_name character varying(255) COLLATE pg_catalog."default",
    ingredient_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    classification character varying(100) COLLATE pg_catalog."default",
    supplier_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    manufacturer_name character varying(255) COLLATE pg_catalog."default",
    unit character varying(50) COLLATE pg_catalog."default",
    specification character varying(100) COLLATE pg_catalog."default",
    unit_price numeric(15, 2) NOT NULL,
    price_per_item numeric(15, 2),
    effective_from timestamp without time zone,
    effective_to timestamp without time zone,
    active boolean NOT NULL DEFAULT true,
    new_buying_price numeric(15, 2),
    promotion character(1) COLLATE pg_catalog."default",
    created_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT supplier_price_list_pkey PRIMARY KEY (product_id)
);

-- ============================================================================
-- SIMPLIFIED: DIRECT SUPPLIER SELECTION
-- ============================================================================
-- No supplier_requests table needed!
-- Just record which supplier was selected for each ingredient in the order
-- ============================================================================

CREATE TABLE IF NOT EXISTS public.order_ingredient_suppliers
(
    order_ingredient_supplier_id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
    order_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    ingredient_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    selected_supplier_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    selected_product_id integer NOT NULL,  -- Links directly to supplier_price_list
    quantity numeric(15, 4) NOT NULL,
    unit character varying(50) COLLATE pg_catalog."default" NOT NULL,
    unit_price numeric(15, 2) NOT NULL,  -- Price at time of selection
    total_cost numeric(15, 2) NOT NULL,
    selection_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    selected_by_user_id character varying(50) COLLATE pg_catalog."default",
    notes text COLLATE pg_catalog."default",
    created_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT order_ingredient_suppliers_pkey PRIMARY KEY (order_ingredient_supplier_id),
    CONSTRAINT uq_order_ingredient_supplier UNIQUE (order_id, ingredient_id)
);

-- Create the kitchen favorite suppliers table
CREATE TABLE IF NOT EXISTS public.kitchen_favorite_suppliers
(
    favorite_id integer NOT NULL GENERATED ALWAYS AS IDENTITY 
        ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    kitchen_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    supplier_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    notes text COLLATE pg_catalog."default",  -- Optional notes about why this supplier is favorited
    display_order integer,  -- Optional: For custom ordering of favorites
    created_by_user_id character varying(50) COLLATE pg_catalog."default",
    created_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_date timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT kitchen_favorite_suppliers_pkey PRIMARY KEY (favorite_id),
    -- Prevent duplicate entries: one kitchen cannot favorite the same supplier twice
    CONSTRAINT uq_kitchen_supplier_favorite UNIQUE (kitchen_id, supplier_id)
);


-- ============================================================================
-- FOREIGN KEY CONSTRAINTS
-- ============================================================================

-- Ingredient Types
ALTER TABLE IF EXISTS public.master_ingredients
    ADD CONSTRAINT fk_ingredient_type FOREIGN KEY (ingredient_type_id)
    REFERENCES public.ingredient_types (ingredient_type_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE SET NULL;

-- Dish Recipe Standards
ALTER TABLE IF EXISTS public.dish_recipe_standards
    ADD CONSTRAINT fk_recipe_dish FOREIGN KEY (dish_id)
    REFERENCES public.master_dishes (dish_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE CASCADE;

ALTER TABLE IF EXISTS public.dish_recipe_standards
    ADD CONSTRAINT fk_recipe_ingredient FOREIGN KEY (ingredient_id)
    REFERENCES public.master_ingredients (ingredient_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE CASCADE;

ALTER TABLE IF EXISTS public.dish_recipe_standards
    ADD CONSTRAINT fk_recipe_user FOREIGN KEY (updated_by_user_id)
    REFERENCES public.master_users (user_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE SET NULL;

-- Master Users
ALTER TABLE IF EXISTS public.master_users
    ADD CONSTRAINT fk_users_kitchen FOREIGN KEY (kitchen_id)
    REFERENCES public.master_kitchens (kitchen_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE SET NULL;

-- Orders
ALTER TABLE IF EXISTS public.orders
    ADD CONSTRAINT fk_order_kitchen FOREIGN KEY (kitchen_id)
    REFERENCES public.master_kitchens (kitchen_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE RESTRICT;

ALTER TABLE IF EXISTS public.orders
    ADD CONSTRAINT fk_order_user FOREIGN KEY (created_by_user_id)
    REFERENCES public.master_users (user_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE SET NULL;

-- Order Details
ALTER TABLE IF EXISTS public.order_details
    ADD CONSTRAINT fk_detail_order FOREIGN KEY (order_id)
    REFERENCES public.orders (order_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE CASCADE;

ALTER TABLE IF EXISTS public.order_details
    ADD CONSTRAINT fk_detail_dish FOREIGN KEY (dish_id)
    REFERENCES public.master_dishes (dish_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE RESTRICT;

-- Order Ingredients
ALTER TABLE IF EXISTS public.order_ingredients
    ADD CONSTRAINT fk_order_ing_detail FOREIGN KEY (order_detail_id)
    REFERENCES public.order_details (order_detail_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE CASCADE;

ALTER TABLE IF EXISTS public.order_ingredients
    ADD CONSTRAINT fk_order_ing_ingredient FOREIGN KEY (ingredient_id)
    REFERENCES public.master_ingredients (ingredient_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE RESTRICT;

-- Order Supplementary Foods
ALTER TABLE IF EXISTS public.order_supplementary_foods
    ADD CONSTRAINT fk_supp_order FOREIGN KEY (order_id)
    REFERENCES public.orders (order_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE CASCADE;

ALTER TABLE IF EXISTS public.order_supplementary_foods
    ADD CONSTRAINT fk_supp_ingredient FOREIGN KEY (ingredient_id)
    REFERENCES public.master_ingredients (ingredient_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE RESTRICT;

-- Supplier Price List
ALTER TABLE IF EXISTS public.supplier_price_list
    ADD CONSTRAINT fk_price_ingredient FOREIGN KEY (ingredient_id)
    REFERENCES public.master_ingredients (ingredient_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE RESTRICT;

ALTER TABLE IF EXISTS public.supplier_price_list
    ADD CONSTRAINT fk_price_supplier FOREIGN KEY (supplier_id)
    REFERENCES public.master_suppliers (supplier_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE RESTRICT;

-- Order Ingredient Suppliers (SIMPLIFIED - Direct link to price list)
ALTER TABLE IF EXISTS public.order_ingredient_suppliers
    ADD CONSTRAINT fk_ois_order FOREIGN KEY (order_id)
    REFERENCES public.orders (order_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE CASCADE;

ALTER TABLE IF EXISTS public.order_ingredient_suppliers
    ADD CONSTRAINT fk_ois_ingredient FOREIGN KEY (ingredient_id)
    REFERENCES public.master_ingredients (ingredient_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE RESTRICT;

ALTER TABLE IF EXISTS public.order_ingredient_suppliers
    ADD CONSTRAINT fk_ois_supplier FOREIGN KEY (selected_supplier_id)
    REFERENCES public.master_suppliers (supplier_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE RESTRICT;

ALTER TABLE IF EXISTS public.order_ingredient_suppliers
    ADD CONSTRAINT fk_ois_product FOREIGN KEY (selected_product_id)
    REFERENCES public.supplier_price_list (product_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE RESTRICT;

ALTER TABLE IF EXISTS public.order_ingredient_suppliers
    ADD CONSTRAINT fk_ois_user FOREIGN KEY (selected_by_user_id)
    REFERENCES public.master_users (user_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE SET NULL;

ALTER TABLE IF EXISTS public.kitchen_favorite_suppliers
    ADD CONSTRAINT fk_favorite_kitchen FOREIGN KEY (kitchen_id)
    REFERENCES public.master_kitchens (kitchen_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE CASCADE;  -- If kitchen is deleted, remove all its favorites

ALTER TABLE IF EXISTS public.kitchen_favorite_suppliers
    ADD CONSTRAINT fk_favorite_supplier FOREIGN KEY (supplier_id)
    REFERENCES public.master_suppliers (supplier_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE CASCADE;  -- If supplier is deleted, remove it from all favorites

ALTER TABLE IF EXISTS public.kitchen_favorite_suppliers
    ADD CONSTRAINT fk_favorite_user FOREIGN KEY (created_by_user_id)
    REFERENCES public.master_users (user_id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE SET NULL;  -- If user is deleted, keep the favorite but nullify the user

-- ============================================================================
-- INDEXES FOR PERFORMANCE
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_recipe_dish ON public.dish_recipe_standards(dish_id);
CREATE INDEX IF NOT EXISTS idx_recipe_ingredient ON public.dish_recipe_standards(ingredient_id);

CREATE INDEX IF NOT EXISTS idx_orders_kitchen ON public.orders(kitchen_id);
CREATE INDEX IF NOT EXISTS idx_orders_date ON public.orders(order_date);
CREATE INDEX IF NOT EXISTS idx_orders_status ON public.orders(status);

CREATE INDEX IF NOT EXISTS idx_order_details_order ON public.order_details(order_id);
CREATE INDEX IF NOT EXISTS idx_order_details_dish ON public.order_details(dish_id);

CREATE INDEX IF NOT EXISTS idx_order_ing_detail ON public.order_ingredients(order_detail_id);
CREATE INDEX IF NOT EXISTS idx_order_ing_ingredient ON public.order_ingredients(ingredient_id);

CREATE INDEX IF NOT EXISTS idx_supplementary_order ON public.order_supplementary_foods(order_id);
CREATE INDEX IF NOT EXISTS idx_supplementary_ingredient ON public.order_supplementary_foods(ingredient_id);

CREATE INDEX IF NOT EXISTS idx_supplier_price_ingredient ON public.supplier_price_list(ingredient_id);
CREATE INDEX IF NOT EXISTS idx_supplier_price_supplier ON public.supplier_price_list(supplier_id);
CREATE INDEX IF NOT EXISTS idx_supplier_price_active ON public.supplier_price_list(active);

CREATE INDEX IF NOT EXISTS idx_ois_order ON public.order_ingredient_suppliers(order_id);
CREATE INDEX IF NOT EXISTS idx_ois_ingredient ON public.order_ingredient_suppliers(ingredient_id);
CREATE INDEX IF NOT EXISTS idx_ois_supplier ON public.order_ingredient_suppliers(selected_supplier_id);
CREATE INDEX IF NOT EXISTS idx_ois_product ON public.order_ingredient_suppliers(selected_product_id);

CREATE INDEX IF NOT EXISTS idx_favorite_kitchen ON public.kitchen_favorite_suppliers(kitchen_id);    
CREATE INDEX IF NOT EXISTS idx_favorite_supplier ON public.kitchen_favorite_suppliers(supplier_id);

END;
```

#### db/migrate_order_id_to_string.sql
*Language: SQL | Size: 5270 bytes*

```sql
-- =========================================================
-- MIGRATION: Change order_id from INTEGER to VARCHAR(50)
-- =========================================================
-- This script converts the order_id column from INTEGER to VARCHAR(50)
-- and migrates all existing records and foreign key relationships.

BEGIN;

-- Step 1: Drop all foreign key constraints that reference orders.order_id
-- This is necessary before we can change the column type

ALTER TABLE IF EXISTS public.order_details 
    DROP CONSTRAINT IF EXISTS fk_detail_order;

ALTER TABLE IF EXISTS public.order_supplementary_foods 
    DROP CONSTRAINT IF EXISTS fk_supp_order;

ALTER TABLE IF EXISTS public.supplier_requests 
    DROP CONSTRAINT IF EXISTS fk_req_order;

-- Step 2: Change column types to VARCHAR(50) in all foreign key tables first
-- PostgreSQL will automatically convert INTEGER to TEXT/VARCHAR during type change
ALTER TABLE public.order_details 
    ALTER COLUMN order_id TYPE VARCHAR(50) USING CAST(order_id AS VARCHAR(50));

ALTER TABLE public.order_supplementary_foods 
    ALTER COLUMN order_id TYPE VARCHAR(50) USING CAST(order_id AS VARCHAR(50));

ALTER TABLE public.supplier_requests 
    ALTER COLUMN order_id TYPE VARCHAR(50) USING CAST(order_id AS VARCHAR(50));

-- Step 3: Change the primary key column in orders table
-- Since we can't directly change an IDENTITY column, we'll:
-- 1. Create a new VARCHAR column
-- 2. Copy data (convert integer IDs to strings)
-- 3. Drop old column and constraints
-- 4. Rename new column

-- Add new column for order_id as VARCHAR
ALTER TABLE public.orders 
    ADD COLUMN order_id_new VARCHAR(50);

-- Convert existing integer IDs to strings
UPDATE public.orders 
SET order_id_new = CAST(order_id AS VARCHAR(50));

-- Make sure the new column is NOT NULL
ALTER TABLE public.orders 
    ALTER COLUMN order_id_new SET NOT NULL;

-- Step 4: Update all foreign key references to use the new string column
-- At this point, FK columns are already VARCHAR, so we can match by string value
UPDATE public.order_details od
SET order_id = o.order_id_new
FROM public.orders o
WHERE od.order_id = CAST(o.order_id AS VARCHAR(50));

UPDATE public.order_supplementary_foods osf
SET order_id = o.order_id_new
FROM public.orders o
WHERE osf.order_id = CAST(o.order_id AS VARCHAR(50));

UPDATE public.supplier_requests sr
SET order_id = o.order_id_new
FROM public.orders o
WHERE sr.order_id = CAST(o.order_id AS VARCHAR(50));

-- Step 5: Drop the old order_id column and rename the new one
-- First drop the primary key constraint
ALTER TABLE public.orders 
    DROP CONSTRAINT IF EXISTS orders_pkey;

-- Drop the old order_id column
ALTER TABLE public.orders 
    DROP COLUMN order_id;

-- Rename the new column to order_id
ALTER TABLE public.orders 
    RENAME COLUMN order_id_new TO order_id;

-- Add primary key constraint on the new order_id column
ALTER TABLE public.orders 
    ADD PRIMARY KEY (order_id);

-- Step 7: Recreate all foreign key constraints
ALTER TABLE public.order_details 
    ADD CONSTRAINT fk_detail_order 
    FOREIGN KEY (order_id) 
    REFERENCES public.orders(order_id) 
    ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE public.order_supplementary_foods 
    ADD CONSTRAINT fk_supp_order 
    FOREIGN KEY (order_id) 
    REFERENCES public.orders(order_id) 
    ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE public.supplier_requests 
    ADD CONSTRAINT fk_req_order 
    FOREIGN KEY (order_id) 
    REFERENCES public.orders(order_id) 
    ON UPDATE CASCADE ON DELETE CASCADE;

-- Step 8: Recreate indexes (they should still exist, but ensure they're correct)
CREATE INDEX IF NOT EXISTS idx_order_details_order ON public.order_details (order_id);
CREATE INDEX IF NOT EXISTS idx_supplementary_order ON public.order_supplementary_foods (order_id);
CREATE INDEX IF NOT EXISTS idx_supplier_requests_order ON public.supplier_requests (order_id);

COMMIT;

-- =========================================================
-- VERIFICATION QUERIES (run these after migration to verify)
-- =========================================================
-- 
-- SELECT 'orders' as table_name, COUNT(*) as total_rows, 
--        COUNT(DISTINCT order_id) as unique_ids,
--        MIN(order_id) as min_id, MAX(order_id) as max_id
-- FROM public.orders
-- UNION ALL
-- SELECT 'order_details', COUNT(*), COUNT(DISTINCT order_id), MIN(order_id), MAX(order_id)
-- FROM public.order_details
-- UNION ALL
-- SELECT 'order_supplementary_foods', COUNT(*), COUNT(DISTINCT order_id), MIN(order_id), MAX(order_id)
-- FROM public.order_supplementary_foods
-- UNION ALL
-- SELECT 'supplier_requests', COUNT(*), COUNT(DISTINCT order_id), MIN(order_id), MAX(order_id)
-- FROM public.supplier_requests;
--
-- -- Check foreign key integrity
-- SELECT COUNT(*) as orphaned_details
-- FROM public.order_details od
-- LEFT JOIN public.orders o ON od.order_id = o.order_id
-- WHERE o.order_id IS NULL;
--
-- SELECT COUNT(*) as orphaned_supplementary
-- FROM public.order_supplementary_foods osf
-- LEFT JOIN public.orders o ON osf.order_id = o.order_id
-- WHERE o.order_id IS NULL;
--
-- SELECT COUNT(*) as orphaned_requests
-- FROM public.supplier_requests sr
-- LEFT JOIN public.orders o ON sr.order_id = o.order_id
-- WHERE o.order_id IS NULL;

```

#### db/samples.sql
*Language: SQL | Size: 5961 bytes*

```sql
-- =========================================================
-- SAMPLE DATA FOR CENTRAL KITCHEN MANAGEMENT
-- =========================================================

-- ========================
-- MASTER TABLES
-- ========================

INSERT INTO public.master_suppliers (supplier_id, supplier_name, zalo_link, address, phone, email)
VALUES
('NCC001', 'Công ty TNHH Thực Phẩm Sạch An Tâm', 'https://zalo.me/antamfood', '123 Lý Thường Kiệt, Q.10, TP.HCM', '0909123456', 'contact@antam.vn'),
('NCC002', 'Nhà Cung Cấp Rau Củ Quả Việt Xanh', 'https://zalo.me/vietxanh', '25 Nguyễn Văn Cừ, Q.5, TP.HCM', '0912345678', 'info@vietxanh.vn'),
('NCC003', 'Công ty TNHH Hải Sản Biển Đông', 'https://zalo.me/biendongseafood', '88 Trần Hưng Đạo, Q.1, TP.HCM', '0987654321', 'sales@biendong.vn');

INSERT INTO public.master_kitchens (kitchen_id, kitchen_name, address, phone)
VALUES
('BEP001', 'Bếp Trung Tâm Quận 1', '12 Nguyễn Thị Minh Khai, Q.1, TP.HCM', '02839123456'),
('BEP002', 'Bếp Trung Tâm Quận 7', '45 Nguyễn Văn Linh, Q.7, TP.HCM', '02837771234');

INSERT INTO public.master_users (user_id, user_name, password, full_name, role, kitchen_id, email, phone)
VALUES
('USR001', 'admin', 'hashed_password', 'Nguyễn Văn Quản Lý', 'Admin', 'BEP001', 'admin@beptrungtam.vn', '0909000001'),
('NV001', 'NV001', '1234', 'Nguyễn Văn Quản Lý', 'Admin', 'BEP001', 'admin@beptrungtam.vn', '0909000001'),
('USR002', 'beptruong1', 'hashed_password', 'Trần Thị Bếp Trưởng', 'Bếp trưởng', 'BEP001', 'beptruong1@beptrungtam.vn', '0909000002'),
('USR003', 'nhanvien1', 'hashed_password', 'Lê Văn Nhân Viên', 'Nhân viên', 'BEP002', 'nhanvien1@beptrungtam.vn', '0909000003');

INSERT INTO public.master_dishes (dish_id, dish_name, cooking_method, category, description)
VALUES
('PHO001', 'Phở Bò Tái', 'Luộc/Nấu nước dùng', 'Món nước', 'Phở bò với nước dùng trong và thịt bò tái mỏng.'),
('COM001', 'Cơm Gà Xối Mỡ', 'Chiên', 'Món chính', 'Cơm trắng với gà chiên giòn rưới mỡ hành.'),
('GOI001', 'Gỏi Cuốn Tôm Thịt', 'Cuốn', 'Khai vị', 'Gỏi cuốn với tôm, thịt, bún và rau sống.');

INSERT INTO public.master_ingredients (ingredient_id, ingredient_name, properties, material_group, unit)
VALUES
('ING001', 'Thịt bò thăn', 'Tươi', 'Thịt', 'kg'),
('ING002', 'Bánh phở', 'Khô', 'Tinh bột', 'kg'),
('ING003', 'Hành lá', 'Tươi', 'Rau củ', 'kg'),
('ING004', 'Gà ta', 'Tươi', 'Thịt', 'kg'),
('ING005', 'Tôm sú', 'Tươi', 'Hải sản', 'kg'),
('ING006', 'Bún tươi', 'Tươi', 'Tinh bột', 'kg'),
('ING007', 'Bánh tráng', 'Khô', 'Tinh bột', 'bó'),
('ING008', 'Rau thơm', 'Tươi', 'Rau củ', 'kg');

-- ========================
-- SUPPLIER PRICE LIST
-- ========================

INSERT INTO public.supplier_price_list (ingredient_id, classification, supplier_id, manufacturer_name, unit, specification, unit_price, effective_from)
VALUES
('ING001', 'Hàng loại 1', 'NCC001', 'Trang Trại Bò Việt', 'kg', 'Bò thăn tươi loại A', 250000, NOW()),
('ING002', 'Hàng tiêu chuẩn', 'NCC002', 'Công ty Bột Gạo An Bình', 'kg', 'Bánh phở khô 1kg', 40000, NOW()),
('ING003', 'Hàng tươi', 'NCC002', 'Trang Trại Rau Việt', 'kg', 'Hành lá loại 1', 30000, NOW()),
('ING004', 'Hàng tươi sống', 'NCC001', 'Trang Trại Gà Long Thành', 'kg', 'Gà ta làm sẵn', 120000, NOW()),
('ING005', 'Hải sản đông lạnh', 'NCC003', 'Biển Đông Seafood', 'kg', 'Tôm sú size 30-40', 280000, NOW());

-- ========================
-- DISH RECIPE STANDARDS
-- ========================

INSERT INTO public.dish_recipe_standards (dish_id, ingredient_id, unit, quantity_per_serving, notes, updated_by_user_id, cost)
VALUES
('PHO001', 'ING001', 'kg', 0.12, 'Thịt bò tái', 'USR002', 30000),
('PHO001', 'ING002', 'kg', 0.08, 'Bánh phở mềm', 'USR002', 3200),
('PHO001', 'ING003', 'kg', 0.02, 'Hành lá cắt nhỏ', 'USR002', 600),
('COM001', 'ING004', 'kg', 0.25, 'Gà chiên giòn', 'USR002', 30000),
('GOI001', 'ING005', 'kg', 0.05, 'Tôm luộc bóc vỏ', 'USR002', 14000),
('GOI001', 'ING007', 'bó', 0.2, 'Bánh tráng mỏng', 'USR002', 2000),
('GOI001', 'ING008', 'kg', 0.05, 'Rau thơm tươi', 'USR002', 2500);

-- ========================
-- ORDERS AND DETAILS
-- ========================

INSERT INTO public.orders (kitchen_id, order_date, note, status, created_by_user_id)
VALUES
('BEP001', '2025-11-02', 'Chuẩn bị cho bữa trưa thứ Hai', 'Pending', 'USR002'),
('BEP002', '2025-11-02', 'Đơn hàng cho sự kiện công ty', 'Pending', 'USR003');

INSERT INTO public.order_details (order_id, dish_id, portions, note)
VALUES
(1, 'PHO001', 30, 'Phở bò cho nhân viên văn phòng'),
(1, 'GOI001', 20, 'Gỏi cuốn khai vị'),
(2, 'COM001', 50, 'Cơm gà phục vụ sự kiện');

INSERT INTO public.order_ingredients (order_detail_id, ingredient_id, quantity, unit, standard_per_portion)
VALUES
(1, 'ING001', 3.6, 'kg', 0.12),
(1, 'ING002', 2.4, 'kg', 0.08),
(1, 'ING003', 0.6, 'kg', 0.02),
(2, 'ING005', 1.0, 'kg', 0.05),
(2, 'ING007', 4.0, 'bó', 0.2),
(2, 'ING008', 1.0, 'kg', 0.05),
(3, 'ING004', 12.5, 'kg', 0.25);

INSERT INTO public.order_supplementary_foods (order_id, ingredient_id, quantity, unit, note)
VALUES
(1, 'ING003', 0.2, 'kg', 'Thêm hành lá dự phòng'),
(2, 'ING008', 0.5, 'kg', 'Thêm rau thơm cho món trang trí');

-- ========================
-- SUPPLIER REQUESTS
-- ========================

INSERT INTO public.supplier_requests (order_id, supplier_id, status)
VALUES
(1, 'NCC001', 'Pending'),
(2, 'NCC002', 'Pending');

INSERT INTO public.supplier_request_details (request_id, ingredient_id, quantity, unit, unit_price)
VALUES
(1, 'ING001', 3.6, 'kg', 250000),
(1, 'ING004', 12.5, 'kg', 120000),
(2, 'ING002', 2.4, 'kg', 40000),
(2, 'ING003', 0.6, 'kg', 30000),
(2, 'ING008', 1.0, 'kg', 25000);
```

### handler/

#### handler/dish.go
*Language: Go | Size: 4221 bytes*

```go
package handler

import (
	"adong-be/logger"
	"adong-be/models"
	"adong-be/store"
	"adong-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDishes with pagination and search
func GetDishes(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("GetDishes called", "user_id", uid)
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Log.Error("GetDishes bind query error", "error", err)
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
	countDB := store.DB.GormClient.Model(&models.Dish{})

	searchConfig := utils.SearchConfig{
		Fields: []string{"dish_name", "dish_id", "description"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		logger.Log.Error("GetDishes count error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var dishes []models.Dish
	db := store.DB.GormClient.Model(&models.Dish{})
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"dish_id":        "dish_id",
		"dish_name":      "dish_name",
		"cooking_method": "cooking_method",
		"category":       "category",
		"created_date":   "created_date",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	if err := db.Find(&dishes).Error; err != nil {
		logger.Log.Error("GetDishes query error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: dishes,
		Meta: meta,
	})
}

func GetDish(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("GetDish called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	var dish models.Dish
	if err := store.DB.GormClient.First(&dish, "dish_id = ?", id).Error; err != nil {
		logger.Log.Error("GetDish not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Dish not found"})
		return
	}
	c.JSON(http.StatusOK, dish)
}

func CreateDish(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("CreateDish called", "user_id", uid)
	var dish models.Dish
	if err := c.ShouldBindJSON(&dish); err != nil {
		logger.Log.Error("CreateDish bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Create(&dish).Error; err != nil {
		logger.Log.Error("CreateDish db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dish)
}

func UpdateDish(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("UpdateDish called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	var dish models.Dish
	if err := store.DB.GormClient.First(&dish, "dish_id = ?", id).Error; err != nil {
		logger.Log.Error("UpdateDish not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Dish not found"})
		return
	}
	if err := c.ShouldBindJSON(&dish); err != nil {
		logger.Log.Error("UpdateDish bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Save(&dish).Error; err != nil {
		logger.Log.Error("UpdateDish db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dish)
}

func DeleteDish(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("DeleteDish called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.Dish{}, "dish_id = ?", id).Error; err != nil {
		logger.Log.Error("DeleteDish db error", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Dish deleted successfully"})
}
```

#### handler/dish_test.go
*Language: Go | Size: 1643 bytes*

```go
package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetDishes_WithPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	router := gin.Default()
	router.GET("/dishes", GetDishes)

	// Test case 1: Default pagination
	req, _ := http.NewRequest("GET", "/dishes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "data")
	assert.Contains(t, response, "meta")

	meta := response["meta"].(map[string]interface{})
	assert.Equal(t, float64(1), meta["current_page"])
	assert.Equal(t, float64(10), meta["page_size"])

	// Test case 2: Custom pagination
	req, _ = http.NewRequest("GET", "/dishes?page=2&page_size=5", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	json.Unmarshal(w.Body.Bytes(), &response)
	meta = response["meta"].(map[string]interface{})
	assert.Equal(t, float64(2), meta["current_page"])
	assert.Equal(t, float64(5), meta["page_size"])
}

func TestGetDishes_WithSearch(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/dishes", GetDishes)

	req, _ := http.NewRequest("GET", "/dishes?search=gà", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "data")
	assert.Contains(t, response, "meta")
}
```

#### handler/ingredient.go
*Language: Go | Size: 4422 bytes*

```go
package handler

import (
	"adong-be/logger"
	"adong-be/models"
	"adong-be/store"
	"adong-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetIngredients with pagination and search - Returns ResourceCollection format
func GetIngredients(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("GetIngredients called", "user_id", uid)
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Log.Error("GetIngredients bind query error", "error", err)
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
	countDB := store.DB.GormClient.Model(&models.Ingredient{})

	searchConfig := utils.SearchConfig{
		Fields: []string{"ingredient_name", "ingredient_id"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		logger.Log.Error("GetIngredients count error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var items []models.Ingredient
	db := store.DB.GormClient.Model(&models.Ingredient{})
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"ingredient_id":   "ingredient_id",
		"ingredient_name": "ingredient_name",
		"unit":            "unit",
		"created_date":    "created_date",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	if err := db.Find(&items).Error; err != nil {
		logger.Log.Error("GetIngredients query error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: items,
		Meta: meta,
	})
}

func GetIngredient(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("GetIngredient called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	var item models.Ingredient
	if err := store.DB.GormClient.First(&item, "ingredient_id = ?", id).Error; err != nil {
		logger.Log.Error("GetIngredient not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Ingredient not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func CreateIngredient(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("CreateIngredient called", "user_id", uid)
	var item models.Ingredient
	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Log.Error("CreateIngredient bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Create(&item).Error; err != nil {
		logger.Log.Error("CreateIngredient db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func UpdateIngredient(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("UpdateIngredient called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	var item models.Ingredient
	if err := store.DB.GormClient.First(&item, "ingredient_id = ?", id).Error; err != nil {
		logger.Log.Error("UpdateIngredient not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Ingredient not found"})
		return
	}
	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Log.Error("UpdateIngredient bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Save(&item).Error; err != nil {
		logger.Log.Error("UpdateIngredient db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func DeleteIngredient(c *gin.Context) {
	uid, _ := c.Get("identity")
	logger.Log.Info("DeleteIngredient called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.Ingredient{}, "ingredient_id = ?", id).Error; err != nil {
		logger.Log.Error("DeleteIngredient db error", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ingredient deleted successfully"})
}
```

#### handler/kitchen.go
*Language: Go | Size: 11556 bytes*

```go
package handler

import (
	"adong-be/models"
	"adong-be/logger"
	"adong-be/store"
	"adong-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetKitchens with pagination and search - Returns ResourceCollection format
func GetKitchens(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("GetKitchens called", "user_id", uid)
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Log.Error("GetKitchens bind query error", "error", err)
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
	countDB := store.DB.GormClient.Model(&models.Kitchen{})

	searchConfig := utils.SearchConfig{
		Fields: []string{"kitchen_name", "kitchen_id", "address"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		logger.Log.Error("GetKitchens count error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var items []models.Kitchen
	db := store.DB.GormClient.Model(&models.Kitchen{})
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"kitchen_id":   "kitchen_id",
		"kitchen_name": "kitchen_name",
		"address":      "address",
		"created_date": "created_date",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	if err := db.Find(&items).Error; err != nil {
		logger.Log.Error("GetKitchens query error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: items,
		Meta: meta,
	})
}

func GetKitchen(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("GetKitchen called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	var item models.Kitchen
	if err := store.DB.GormClient.First(&item, "kitchen_id = ?", id).Error; err != nil {
		logger.Log.Error("GetKitchen not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitchen not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func CreateKitchen(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("CreateKitchen called", "user_id", uid)
	var item models.Kitchen
	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Log.Error("CreateKitchen bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Create(&item).Error; err != nil {
		logger.Log.Error("CreateKitchen db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func UpdateKitchen(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("UpdateKitchen called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	var item models.Kitchen
	if err := store.DB.GormClient.First(&item, "kitchen_id = ?", id).Error; err != nil {
		logger.Log.Error("UpdateKitchen not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitchen not found"})
		return
	}
	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Log.Error("UpdateKitchen bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Save(&item).Error; err != nil {
		logger.Log.Error("UpdateKitchen db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func DeleteKitchen(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("DeleteKitchen called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.Kitchen{}, "kitchen_id = ?", id).Error; err != nil {
		logger.Log.Error("DeleteKitchen db error", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kitchen deleted successfully"})
}

// ============================================================================
// Kitchen Favorite Suppliers Handlers
// ============================================================================

// GetKitchenFavoriteSuppliers returns all favorite suppliers for a kitchen
func GetKitchenFavoriteSuppliers(c *gin.Context) {
	uid, _ := c.Get("identity")
	kitchenID := c.Param("id")
	logger.Log.Info("GetKitchenFavoriteSuppliers called", "kitchen_id", kitchenID, "user_id", uid)

	// Validate kitchen exists
	var kitchen models.Kitchen
	if err := store.DB.GormClient.First(&kitchen, "kitchen_id = ?", kitchenID).Error; err != nil {
		logger.Log.Error("GetKitchenFavoriteSuppliers kitchen not found", "kitchen_id", kitchenID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitchen not found"})
		return
	}

	var favorites []models.KitchenFavoriteSupplier
	query := store.DB.GormClient.
		Where("kitchen_id = ?", kitchenID).
		Preload("Supplier").
		Preload("CreatedBy")

	// Order by display_order if set, otherwise by created_date
	query = query.Order("COALESCE(display_order, 999999), created_date ASC")

	if err := query.Find(&favorites).Error; err != nil {
		logger.Log.Error("GetKitchenFavoriteSuppliers db error", "kitchen_id", kitchenID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, favorites)
}

// GetKitchenFavoriteSupplier returns a single favorite supplier by ID
func GetKitchenFavoriteSupplier(c *gin.Context) {
	uid, _ := c.Get("identity")
	kitchenID := c.Param("id")
	favoriteID := c.Param("favoriteId")
	logger.Log.Info("GetKitchenFavoriteSupplier called", "kitchen_id", kitchenID, "favorite_id", favoriteID, "user_id", uid)

	var favorite models.KitchenFavoriteSupplier
	if err := store.DB.GormClient.
		Where("favorite_id = ? AND kitchen_id = ?", favoriteID, kitchenID).
		Preload("Kitchen").
		Preload("Supplier").
		Preload("CreatedBy").
		First(&favorite).Error; err != nil {
		logger.Log.Error("GetKitchenFavoriteSupplier not found", "kitchen_id", kitchenID, "favorite_id", favoriteID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Favorite supplier not found"})
		return
	}

	c.JSON(http.StatusOK, favorite)
}

// CreateKitchenFavoriteSupplier adds a supplier to a kitchen's favorites
func CreateKitchenFavoriteSupplier(c *gin.Context) {
	uid, _ := c.Get("identity")
	kitchenID := c.Param("id")
	logger.Log.Info("CreateKitchenFavoriteSupplier called", "kitchen_id", kitchenID, "user_id", uid)

	// Get user ID from authentication middleware
	var userID string
	if identity, ok := c.Get("identity"); ok {
		if v, ok2 := identity.(string); ok2 {
			userID = v
		}
	}

	var favorite models.KitchenFavoriteSupplier
	if err := c.ShouldBindJSON(&favorite); err != nil {
		logger.Log.Error("CreateKitchenFavoriteSupplier bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set kitchen_id from URL parameter
	favorite.KitchenID = kitchenID
	favorite.CreatedByUserID = userID

	// Validate kitchen exists
	var kitchen models.Kitchen
	if err := store.DB.GormClient.First(&kitchen, "kitchen_id = ?", kitchenID).Error; err != nil {
		logger.Log.Error("CreateKitchenFavoriteSupplier kitchen not found", "kitchen_id", kitchenID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitchen not found"})
		return
	}

	// Validate supplier exists
	var supplier models.Supplier
	if err := store.DB.GormClient.First(&supplier, "supplier_id = ?", favorite.SupplierID).Error; err != nil {
		logger.Log.Error("CreateKitchenFavoriteSupplier supplier not found", "supplier_id", favorite.SupplierID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
		return
	}

	// Check if favorite already exists (unique constraint: kitchen_id + supplier_id)
	var existing models.KitchenFavoriteSupplier
	if err := store.DB.GormClient.Where("kitchen_id = ? AND supplier_id = ?", kitchenID, favorite.SupplierID).First(&existing).Error; err == nil {
		logger.Log.Error("CreateKitchenFavoriteSupplier duplicate favorite", "kitchen_id", kitchenID, "supplier_id", favorite.SupplierID)
		c.JSON(http.StatusConflict, gin.H{"error": "This supplier is already in the kitchen's favorites"})
		return
	}

	// Create favorite
	if err := store.DB.GormClient.Create(&favorite).Error; err != nil {
		logger.Log.Error("CreateKitchenFavoriteSupplier db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload with relations
	store.DB.GormClient.
		Preload("Kitchen").
		Preload("Supplier").
		Preload("CreatedBy").
		First(&favorite, "favorite_id = ?", favorite.FavoriteID)

	c.JSON(http.StatusCreated, favorite)
}

// UpdateKitchenFavoriteSupplier updates a favorite supplier entry
func UpdateKitchenFavoriteSupplier(c *gin.Context) {
	uid, _ := c.Get("identity")
	kitchenID := c.Param("id")
	favoriteID := c.Param("favoriteId")
	logger.Log.Info("UpdateKitchenFavoriteSupplier called", "kitchen_id", kitchenID, "favorite_id", favoriteID, "user_id", uid)

	var favorite models.KitchenFavoriteSupplier
	if err := store.DB.GormClient.Where("favorite_id = ? AND kitchen_id = ?", favoriteID, kitchenID).First(&favorite).Error; err != nil {
		logger.Log.Error("UpdateKitchenFavoriteSupplier not found", "kitchen_id", kitchenID, "favorite_id", favoriteID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Favorite supplier not found"})
		return
	}

	// Define update struct (don't allow changing kitchen_id or supplier_id)
	var updateData struct {
		Notes        string `json:"notes"`
		DisplayOrder *int   `json:"displayOrder"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Log.Error("UpdateKitchenFavoriteSupplier bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update allowed fields
	favorite.Notes = updateData.Notes
	favorite.DisplayOrder = updateData.DisplayOrder

	if err := store.DB.GormClient.Save(&favorite).Error; err != nil {
		logger.Log.Error("UpdateKitchenFavoriteSupplier db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload with relations
	store.DB.GormClient.
		Preload("Kitchen").
		Preload("Supplier").
		Preload("CreatedBy").
		First(&favorite, "favorite_id = ?", favoriteID)

	c.JSON(http.StatusOK, favorite)
}

// DeleteKitchenFavoriteSupplier removes a supplier from a kitchen's favorites
func DeleteKitchenFavoriteSupplier(c *gin.Context) {
	uid, _ := c.Get("identity")
	kitchenID := c.Param("id")
	favoriteID := c.Param("favoriteId")
	logger.Log.Info("DeleteKitchenFavoriteSupplier called", "kitchen_id", kitchenID, "favorite_id", favoriteID, "user_id", uid)

	if err := store.DB.GormClient.Where("favorite_id = ? AND kitchen_id = ?", favoriteID, kitchenID).Delete(&models.KitchenFavoriteSupplier{}).Error; err != nil {
		logger.Log.Error("DeleteKitchenFavoriteSupplier db error", "kitchen_id", kitchenID, "favorite_id", favoriteID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite supplier removed successfully"})
}
```

#### handler/order.go
*Language: Go | Size: 27259 bytes*

```go
package handler

import (
	"adong-be/logger"
	"adong-be/models"
	"adong-be/store"
	"adong-be/utils"
	"net/http"
	"time"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

// SaveOrderIngredientsWithSupplier - Save selected suppliers for ingredients in an order
// This saves (order_id, ingredient_id) pairs with selected supplier/product information
func SaveOrderIngredientsWithSupplier(c *gin.Context) {
	uid, _ := c.Get("identity")
	orderID := c.Param("id")
	logger.Log.Info("SaveOrderIngredientsWithSupplier called", "order_id", orderID, "user_id", uid)

	// Get user ID from authentication middleware
	var userID string
	if identity, ok := c.Get("identity"); ok {
		if v, ok2 := identity.(string); ok2 {
			userID = v
		}
	}

	// Define request structure - list of ingredient selections
	var request struct {
		Selections []struct {
			IngredientID       string  `json:"ingredientId" binding:"required"`
			SelectedSupplierID string  `json:"selectedSupplierId" binding:"required"`
			SelectedProductID int     `json:"selectedProductId" binding:"required"`
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

	// Validate order exists
	var order models.Order
	if err := store.DB.GormClient.First(&order, "order_id = ?", orderID).Error; err != nil {
		logger.Log.Error("SaveOrderIngredientsWithSupplier order not found", "order_id", orderID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Validate all selections before processing
	for i, sel := range request.Selections {
		// Validate ingredient exists
		var ingredient models.Ingredient
		if err := store.DB.GormClient.First(&ingredient, "ingredient_id = ?", sel.IngredientID).Error; err != nil {
			logger.Log.Error("SaveOrderIngredientsWithSupplier ingredient not found", "ingredient_id", sel.IngredientID, "error", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Ingredient not found: " + sel.IngredientID})
			return
		}

		// Validate supplier exists
		var supplier models.Supplier
		if err := store.DB.GormClient.First(&supplier, "supplier_id = ?", sel.SelectedSupplierID).Error; err != nil {
			logger.Log.Error("SaveOrderIngredientsWithSupplier supplier not found", "supplier_id", sel.SelectedSupplierID, "error", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found: " + sel.SelectedSupplierID})
			return
		}

		// Validate product exists and belongs to supplier and ingredient
		var product models.SupplierPrice
		if err := store.DB.GormClient.First(&product, "product_id = ? AND supplier_id = ? AND ingredient_id = ?", 
			sel.SelectedProductID, sel.SelectedSupplierID, sel.IngredientID).Error; err != nil {
			logger.Log.Error("SaveOrderIngredientsWithSupplier product not found or mismatch", 
				"product_id", sel.SelectedProductID, 
				"supplier_id", sel.SelectedSupplierID, 
				"ingredient_id", sel.IngredientID, 
				"error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found or does not match supplier/ingredient"})
			return
		}

		// Validate ingredient belongs to the specified order (either in order_ingredients or order_supplementary_foods)
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
			logger.Log.Error("SaveOrderIngredientsWithSupplier validate ingredient in order error", 
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

		// Check for duplicate ingredient_id in request
		for j := i + 1; j < len(request.Selections); j++ {
			if request.Selections[j].IngredientID == sel.IngredientID {
				logger.Log.Error("SaveOrderIngredientsWithSupplier duplicate ingredient in request", 
					"ingredient_id", sel.IngredientID)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate ingredient_id in request: " + sel.IngredientID})
				return
			}
		}
	}

	// Start transaction
	tx := store.DB.GormClient.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Process each selection - save or update (order_id, ingredient_id) pair
	var savedSelections []models.OrderIngredientSupplier
	for _, sel := range request.Selections {
		// Calculate total cost
		totalCost := sel.Quantity * sel.UnitPrice

		// Check if selection already exists for this (order_id, ingredient_id) pair
		var existing models.OrderIngredientSupplier
		findErr := tx.Where("order_id = ? AND ingredient_id = ?", orderID, sel.IngredientID).First(&existing).Error

		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			// Create new selection
			newSelection := models.OrderIngredientSupplier{
				OrderID:           orderID,
				IngredientID:      sel.IngredientID,
				SelectedSupplierID: sel.SelectedSupplierID,
				SelectedProductID: sel.SelectedProductID,
				Quantity:          sel.Quantity,
				Unit:              sel.Unit,
				UnitPrice:         sel.UnitPrice,
				TotalCost:         totalCost,
				SelectedByUserID:  userID,
				Notes:             sel.Notes,
			}

			if err := tx.Create(&newSelection).Error; err != nil {
				logger.Log.Error("SaveOrderIngredientsWithSupplier create selection error", 
					"error", err, "ingredient_id", sel.IngredientID)
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			savedSelections = append(savedSelections, newSelection)
			logger.Log.Info("SaveOrderIngredientsWithSupplier created new selection", 
				"order_id", orderID, "ingredient_id", sel.IngredientID)
		} else if findErr != nil {
			logger.Log.Error("SaveOrderIngredientsWithSupplier find existing selection error", 
				"error", findErr, "ingredient_id", sel.IngredientID)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": findErr.Error()})
			return
		} else {
			// Update existing selection
			existing.SelectedSupplierID = sel.SelectedSupplierID
			existing.SelectedProductID = sel.SelectedProductID
			existing.Quantity = sel.Quantity
			existing.Unit = sel.Unit
			existing.UnitPrice = sel.UnitPrice
			existing.TotalCost = totalCost
			existing.SelectedByUserID = userID
			existing.Notes = sel.Notes

			if err := tx.Save(&existing).Error; err != nil {
				logger.Log.Error("SaveOrderIngredientsWithSupplier update selection error", 
					"error", err, "ingredient_id", sel.IngredientID)
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			savedSelections = append(savedSelections, existing)
			logger.Log.Info("SaveOrderIngredientsWithSupplier updated existing selection", 
				"order_id", orderID, "ingredient_id", sel.IngredientID)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		logger.Log.Error("SaveOrderIngredientsWithSupplier commit error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload selections with relations for response
	var responseSelections []models.OrderIngredientSupplier
	if err := store.DB.GormClient.
		Preload("Ingredient").
		Preload("SelectedSupplier").
		Preload("SelectedProduct").
		Preload("SelectedBy").
		Where("order_id = ?", orderID).
		Find(&responseSelections).Error; err != nil {
		logger.Log.Error("SaveOrderIngredientsWithSupplier reload error", "error", err)
		// Don't fail the request, just return what we saved
		responseSelections = savedSelections
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message":   "Supplier selections saved successfully",
		"orderId":   orderID,
		"selections": responseSelections,
		"count":     len(savedSelections),
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
```

#### handler/recipe_standard.go
*Language: Go | Size: 7118 bytes*

```go
package handler

import (
	"adong-be/logger"
	"adong-be/models"
	"adong-be/store"
	"adong-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetRecipeStandards with pagination and search - Returns ResourceCollection format with DTOs
func GetRecipeStandards(c *gin.Context) {
	logger.Log.Info("GetRecipeStandards called")
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Log.Error("GetRecipeStandards bind query error", "error", err)
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
	countDB := store.DB.GormClient.Model(&models.RecipeStandard{})

	searchConfig := utils.SearchConfig{
		Fields: []string{"dish_id", "ingredient_id"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		logger.Log.Error("GetRecipeStandards count error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var recipes []models.RecipeStandard
	db := store.DB.GormClient.Model(&models.RecipeStandard{})
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"standardId":   "recipe_id",
		"dishId":       "dish_id",
		"ingredientId": "ingredient_id",
		"standardPer1": "quantity_per_serving",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	// Preload related entities to get names
	db = db.Preload("Dish").Preload("Ingredient").Preload("UpdatedBy")

	if err := db.Find(&recipes).Error; err != nil {
		logger.Log.Error("GetRecipeStandards query error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to DTOs
	dtos := models.ConvertRecipeStandardsToDTO(recipes)

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: dtos,
		Meta: meta,
	})
}

func GetRecipeStandard(c *gin.Context) {
	logger.Log.Info("GetRecipeStandard called", "id", c.Param("id"))
	id := c.Param("id")
	var recipe models.RecipeStandard

	// Preload related entities
	if err := store.DB.GormClient.
		Preload("Dish").
		Preload("Ingredient").
		Preload("UpdatedBy").
		First(&recipe, "recipe_id = ?", id).Error; err != nil {
		logger.Log.Error("GetRecipeStandard not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe standard not found"})
		return
	}

	// Convert to DTO and return
	dto := recipe.ToDTO()
	c.JSON(http.StatusOK, dto)
}

func CreateRecipeStandard(c *gin.Context) {
	logger.Log.Info("CreateRecipeStandard called")
	var recipe models.RecipeStandard
	if err := c.ShouldBindJSON(&recipe); err != nil {
		logger.Log.Error("CreateRecipeStandard bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Create(&recipe).Error; err != nil {
		logger.Log.Error("CreateRecipeStandard db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload with relationships
	store.DB.GormClient.
		Preload("Dish").
		Preload("Ingredient").
		Preload("UpdatedBy").
		First(&recipe, "recipe_id = ?", recipe.StandardID)

	// Return DTO
	dto := recipe.ToDTO()
	c.JSON(http.StatusCreated, dto)
}

func UpdateRecipeStandard(c *gin.Context) {
	logger.Log.Info("UpdateRecipeStandard called", "id", c.Param("id"))
	id := c.Param("id")
	var recipe models.RecipeStandard
	if err := store.DB.GormClient.First(&recipe, "recipe_id = ?", id).Error; err != nil {
		logger.Log.Error("UpdateRecipeStandard not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe standard not found"})
		return
	}
	if err := c.ShouldBindJSON(&recipe); err != nil {
		logger.Log.Error("UpdateRecipeStandard bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Save(&recipe).Error; err != nil {
		logger.Log.Error("UpdateRecipeStandard db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload with relationships
	store.DB.GormClient.
		Preload("Dish").
		Preload("Ingredient").
		Preload("UpdatedBy").
		First(&recipe, "recipe_id = ?", recipe.StandardID)

	// Return DTO
	dto := recipe.ToDTO()
	c.JSON(http.StatusOK, dto)
}

func DeleteRecipeStandard(c *gin.Context) {
	logger.Log.Info("DeleteRecipeStandard called", "id", c.Param("id"))
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.RecipeStandard{}, "recipe_id = ?", id).Error; err != nil {
		logger.Log.Error("DeleteRecipeStandard db error", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe standard deleted successfully"})
}

// GetRecipeStandardsByDish with pagination and search - Returns ResourceCollection format with DTOs
func GetRecipeStandardsByDish(c *gin.Context) {
	logger.Log.Info("GetRecipeStandardsByDish called", "dishId", c.Param("dishId"))
	dishId := c.Param("dishId")

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
	countDB := store.DB.GormClient.Model(&models.RecipeStandard{}).Where("dish_id = ?", dishId)

	searchConfig := utils.SearchConfig{
		Fields: []string{"ingredient_id"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		logger.Log.Error("GetRecipeStandardsByDish count error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var recipes []models.RecipeStandard
	db := store.DB.GormClient.Model(&models.RecipeStandard{}).Where("dish_id = ?", dishId)
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"ingredientId": "ingredient_id",
		"standardPer1": "quantity_per_serving",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	// Preload related entities to get names
	db = db.Preload("Dish").Preload("Ingredient").Preload("UpdatedBy")

	if err := db.Find(&recipes).Error; err != nil {
		logger.Log.Error("GetRecipeStandardsByDish query error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to DTOs
	dtos := models.ConvertRecipeStandardsToDTO(recipes)

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: dtos,
		Meta: meta,
	})
}
```

#### handler/supplier.go
*Language: Go | Size: 4392 bytes*

```go
package handler

import (
	"adong-be/models"
	"adong-be/logger"
	"adong-be/store"
	"adong-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSuppliers with pagination and search - Returns ResourceCollection format
func GetSuppliers(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("GetSuppliers called", "user_id", uid)
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Log.Error("GetSuppliers bind query error", "error", err)
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
	countDB := store.DB.GormClient.Model(&models.Supplier{})

	searchConfig := utils.SearchConfig{
		Fields: []string{"supplier_name", "supplier_id", "address", "phone"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		logger.Log.Error("GetSuppliers count error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var items []models.Supplier
	db := store.DB.GormClient.Model(&models.Supplier{})
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"supplier_id":   "supplier_id",
		"supplier_name": "supplier_name",
		"address":       "address",
		"created_date":  "created_date",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	if err := db.Find(&items).Error; err != nil {
		logger.Log.Error("GetSuppliers query error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: items,
		Meta: meta,
	})
}

func GetSupplier(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("GetSupplier called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	var item models.Supplier
	if err := store.DB.GormClient.First(&item, "supplier_id = ?", id).Error; err != nil {
		logger.Log.Error("GetSupplier not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func CreateSupplier(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("CreateSupplier called", "user_id", uid)
	var item models.Supplier
	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Log.Error("CreateSupplier bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Create(&item).Error; err != nil {
		logger.Log.Error("CreateSupplier db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func UpdateSupplier(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("UpdateSupplier called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	var item models.Supplier
	if err := store.DB.GormClient.First(&item, "supplier_id = ?", id).Error; err != nil {
		logger.Log.Error("UpdateSupplier not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
		return
	}
	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Log.Error("UpdateSupplier bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Save(&item).Error; err != nil {
		logger.Log.Error("UpdateSupplier db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func DeleteSupplier(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("DeleteSupplier called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.Supplier{}, "supplier_id = ?", id).Error; err != nil {
		logger.Log.Error("DeleteSupplier db error", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Supplier deleted successfully"})
}
```

#### handler/supplier_price.go
*Language: Go | Size: 11747 bytes*

```go
package handler

import (
	"adong-be/logger"
	"adong-be/models"
	"adong-be/store"
	"adong-be/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetSupplierPrices(c *gin.Context) {
	logger.Log.Info("GetSupplierPrices called")
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Log.Error("GetSupplierPrices bind query error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get date range parameters
	effectiveFrom := c.Query("effective_from")
	effectiveTo := c.Query("effective_to")
	logger.Log.Debug("receive query", "Effective From:", effectiveFrom, "Effective To:", effectiveTo)

	params = models.GetPaginationParams(
		params.Page,
		params.PageSize,
		params.Search,
		params.SortBy,
		params.SortDir,
	)

	var total int64
	countDB := store.DB.GormClient.Model(&models.SupplierPrice{})

	searchConfig := utils.SearchConfig{
		Fields: []string{"product_name", "ingredient_id", "supplier_id",
		 "classification", "specification", "manufacturer_name"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	// Apply date range filters for counting
	countDB = applyDateRangeFilter(countDB, effectiveFrom, effectiveTo)

	if err := countDB.Count(&total).Error; err != nil {
		logger.Log.Error("GetSupplierPrices count error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var prices []models.SupplierPrice
	db := store.DB.GormClient.Model(&models.SupplierPrice{})
	db = utils.ApplySearch(db, params.Search, searchConfig)

	// Apply date range filters for data query
	db = applyDateRangeFilter(db, effectiveFrom, effectiveTo)

	allowedSortFields := map[string]string{
		"product_id":     "product_id",
		"product_name":   "product_name",
		"ingredient_id":  "ingredient_id",
		"supplier_id":    "supplier_id",
		"unit_price":     "unit_price",
		"effective_from": "effective_from",
		"effective_to":   "effective_to",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	// Preload related entities to get names
	db = db.Preload("Ingredient").Preload("Supplier")

	if err := db.Find(&prices).Error; err != nil {
		logger.Log.Error("GetSupplierPrices query error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to DTOs
	dtos := models.ConvertSupplierPricesToDTO(prices)

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: dtos,
		Meta: meta,
	})
}

// Helper function to apply date range filters
func applyDateRangeFilter(db *gorm.DB, effectiveFrom, effectiveTo string) *gorm.DB {
	// // Parse and validate effectiveFrom date
	// if effectiveFrom != "" {
	// 	// Parse the date string (format: YYYY-MM-DD)
	// 	fromDate, err := time.Parse("2006-01-02", effectiveFrom)
	// 	if err == nil {
	// 		// Filter records where hieuluctu >= effectiveFrom OR hieulucden >= effectiveFrom
	// 		// This ensures we get prices that are effective during or after the from date
	// 		db = db.Where("hieuluctu >= ?", fromDate)
	// 	}
	// }

	// // Parse and validate effectiveTo date
	// if effectiveTo != "" {
	// 	// Parse the date string (format: YYYY-MM-DD)
	// 	toDate, err := time.Parse("2006-01-02", effectiveTo)
	// 	if err == nil {
	// 		// Add 1 day to include the entire end date
	// 		toDateEnd := toDate.Add(24 * time.Hour)
	// 		// Filter records where hieuluctu <= effectiveTo
	// 		// This ensures we get prices that start before or on the to date
	// 		db = db.Where("hieulucden < ?", toDateEnd)
	// 	}
	// }

	if effectiveFrom == "" {
		effectiveFrom = "0001-01-01"
	}
	if effectiveTo == "" {
		effectiveTo = "9999-12-31"
	}

	// If both dates are provided, find prices that overlap with the date range
	if effectiveFrom != "" && effectiveTo != "" {
		fromDate, errFrom := time.Parse("2006-01-02", effectiveFrom)
		toDate, errTo := time.Parse("2006-01-02", effectiveTo)

		if errFrom == nil && errTo == nil {
			// toDateEnd := toDate.Add(24 * time.Hour)
			// Records where:
			// - Start date is within range, OR
			// - End date is within range, OR
			// - The price period encompasses the entire search range
			db = db.Where(
				"(effective_from IS NULL OR effective_to IS NULL) OR (effective_from >= ? AND effective_to <= ?)",
				fromDate, toDate,
			)
		}
	}

	return db
}

func GetSupplierPrice(c *gin.Context) {
	logger.Log.Info("GetSupplierPrice called", "id", c.Param("id"))
	id := c.Param("id")
	var price models.SupplierPrice

	// Preload related entities to get names
	if err := store.DB.GormClient.
		Preload("Ingredient").
		Preload("Supplier").
		First(&price, "product_id = ?", id).Error; err != nil {
		logger.Log.Error("GetSupplierPrice not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Supplier price not found"})
		return
	}

	// Convert to DTO and return
	dto := price.ToDTO()
	c.JSON(http.StatusOK, dto)
}

// GetSupplierPricesByIngredient - Get all supplier prices for a specific ingredient
func GetSupplierPricesByIngredient(c *gin.Context) {
	logger.Log.Info("GetSupplierPricesByIngredient called", "ingredientId", c.Param("ingredientId"))
	ingredientId := c.Param("ingredientId")

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
	countDB := store.DB.GormClient.Model(&models.SupplierPrice{}).Where("ingredient_id = ?", ingredientId)

	searchConfig := utils.SearchConfig{
		Fields: []string{"tensanpham", "nhacungcapid", "phanloai"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		logger.Log.Error("GetSupplierPricesByIngredient count error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var prices []models.SupplierPrice
	db := store.DB.GormClient.Model(&models.SupplierPrice{}).Where("ingredient_id = ?", ingredientId)
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"product_id":     "product_id",
		"product_name":   "product_name",
		"supplier_id":    "supplier_id",
		"unit_price":     "unit_price",
		"effective_from": "effective_from",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	// Preload related entities to get names
	db = db.Preload("Ingredient").Preload("Supplier")

	if err := db.Find(&prices).Error; err != nil {
		logger.Log.Error("GetSupplierPricesByIngredient query error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to DTOs
	dtos := models.ConvertSupplierPricesToDTO(prices)

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: dtos,
		Meta: meta,
	})
}

// GetSupplierPricesBySupplier - Get all supplier prices for a specific supplier
func GetSupplierPricesBySupplier(c *gin.Context) {
	logger.Log.Info("GetSupplierPricesBySupplier called", "supplierId", c.Param("supplierId"))
	supplierId := c.Param("supplierId")

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
	countDB := store.DB.GormClient.Model(&models.SupplierPrice{}).Where("supplier_id = ?", supplierId)

	searchConfig := utils.SearchConfig{
		Fields: []string{"tensanpham", "nguyenlieuid", "phanloai"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		logger.Log.Error("GetSupplierPricesBySupplier count error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var prices []models.SupplierPrice
	db := store.DB.GormClient.Model(&models.SupplierPrice{}).Where("supplier_id = ?", supplierId)
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"product_id":     "product_id",
		"product_name":   "product_name",
		"ingredient_id":  "ingredient_id",
		"unit_price":     "unit_price",
		"effective_from": "effective_from",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	// Preload related entities to get names
	db = db.Preload("Ingredient").Preload("Supplier")

	if err := db.Find(&prices).Error; err != nil {
		logger.Log.Error("GetSupplierPricesBySupplier query error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to DTOs
	dtos := models.ConvertSupplierPricesToDTO(prices)

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: dtos,
		Meta: meta,
	})
}

func CreateSupplierPrice(c *gin.Context) {
	logger.Log.Info("CreateSupplierPrice called")
	var price models.SupplierPrice
	if err := c.ShouldBindJSON(&price); err != nil {
		logger.Log.Error("CreateSupplierPrice bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Create(&price).Error; err != nil {
		logger.Log.Error("CreateSupplierPrice db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload with relationships to get names
	store.DB.GormClient.
		Preload("Ingredient").
		Preload("Supplier").
		First(&price, "product_id = ?", price.ProductID)

	// Return DTO with names
	dto := price.ToDTO()
	c.JSON(http.StatusCreated, dto)
}

func UpdateSupplierPrice(c *gin.Context) {
	logger.Log.Info("UpdateSupplierPrice called", "id", c.Param("id"))
	id := c.Param("id")
	var price models.SupplierPrice
	if err := store.DB.GormClient.First(&price, "product_id = ?", id).Error; err != nil {
		logger.Log.Error("UpdateSupplierPrice not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Supplier price not found"})
		return
	}
	if err := c.ShouldBindJSON(&price); err != nil {
		logger.Log.Error("UpdateSupplierPrice bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Save(&price).Error; err != nil {
		logger.Log.Error("UpdateSupplierPrice db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload with relationships to get names
	store.DB.GormClient.
		Preload("Ingredient").
		Preload("Supplier").
		First(&price, "product_id = ?", price.ProductID)

	// Return DTO with names
	dto := price.ToDTO()
	c.JSON(http.StatusOK, dto)
}

func DeleteSupplierPrice(c *gin.Context) {
	logger.Log.Info("DeleteSupplierPrice called", "id", c.Param("id"))
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.SupplierPrice{}, "product_id = ?", id).Error; err != nil {
		logger.Log.Error("DeleteSupplierPrice db error", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Supplier price deleted successfully"})
}
```

#### handler/user.go
*Language: Go | Size: 4253 bytes*

```go
package handler

import (
	"adong-be/models"
	"adong-be/logger"
	"adong-be/store"
	"adong-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUsers with pagination and search - Returns ResourceCollection format
func GetUsers(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("GetUsers called", "user_id", uid)
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Log.Error("GetUsers bind query error", "error", err)
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
	countDB := store.DB.GormClient.Model(&models.User{})

	searchConfig := utils.SearchConfig{
		Fields: []string{"user_id", "user_name", "full_name", "email", "phone"},
		Fuzzy:  true,
	}
	countDB = utils.ApplySearch(countDB, params.Search, searchConfig)

	if err := countDB.Count(&total).Error; err != nil {
		logger.Log.Error("GetUsers count error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var items []models.User
	db := store.DB.GormClient.Model(&models.User{})
	db = utils.ApplySearch(db, params.Search, searchConfig)

	allowedSortFields := map[string]string{
		"user_id":   "user_id",
		"user_name": "user_name",
		"full_name": "full_name",
		"email":     "email",
		"role":      "role",
	}
	db = utils.ApplySort(db, params.SortBy, params.SortDir, allowedSortFields)
	db = utils.ApplyPagination(db, params.Page, params.PageSize)

	if err := db.Find(&items).Error; err != nil {
		logger.Log.Error("GetUsers query error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	meta := models.CalculatePaginationMeta(params.Page, params.PageSize, total)
	c.JSON(http.StatusOK, models.ResourceCollection{
		Data: items,
		Meta: meta,
	})
}

func GetUser(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("GetUser called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	var item models.User
	if err := store.DB.GormClient.First(&item, "user_id = ?", id).Error; err != nil {
		logger.Log.Error("GetUser not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func CreateUser(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("CreateUser called", "user_id", uid)
	var item models.User
	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Log.Error("CreateUser bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Create(&item).Error; err != nil {
		logger.Log.Error("CreateUser db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func UpdateUser(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("UpdateUser called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	var item models.User
	if err := store.DB.GormClient.First(&item, "user_id = ?", id).Error; err != nil {
		logger.Log.Error("UpdateUser not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Log.Error("UpdateUser bind error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := store.DB.GormClient.Save(&item).Error; err != nil {
		logger.Log.Error("UpdateUser db error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func DeleteUser(c *gin.Context) {
    uid, _ := c.Get("identity")
    logger.Log.Info("DeleteUser called", "id", c.Param("id"), "user_id", uid)
	id := c.Param("id")
	if err := store.DB.GormClient.Delete(&models.User{}, "user_id = ?", id).Error; err != nil {
		logger.Log.Error("DeleteUser db error", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
```

### logger/

#### logger/logger.go
*Language: Go | Size: 1074 bytes*

```go
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerI interface {

	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})

}


type Logger struct {
	*zap.SugaredLogger
}

func NewLogger() *Logger {
	// set caller skip to 2
	cfg := zap.NewProductionConfig()
	
	cfg.Level.SetLevel(zapcore.DebugLevel) // Set the desired level (e.g., InfoLevel)
	logger, _ := cfg.Build()
	logger = logger.WithOptions(zap.AddCallerSkip(1))
	
	sugar := logger.Sugar()
	
	return &Logger{
		SugaredLogger: sugar,
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.SugaredLogger.With(args...).Info(msg)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.SugaredLogger.With(args...).Warn(msg)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.SugaredLogger.With(args...).Error(msg)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.SugaredLogger.With(args...).Debug(msg)
}

var (
	Log  LoggerI = NewLogger()
)
```

### models/

#### models/common.go
*Language: Go | Size: 17 bytes*

```go
package models


```

#### models/dish.go
*Language: Go | Size: 816 bytes*

```go
package models

import "time"

// Dish - Master data for dishes/menu items (dm_monan)
type Dish struct {
    DishID        string    `gorm:"primaryKey;column:dish_id" json:"dishId"`
    DishName      string    `gorm:"column:dish_name;not null" json:"dishName"`
    CookingMethod string    `gorm:"column:cooking_method" json:"cookingMethod"`
    Group         string    `gorm:"column:category" json:"group"`
    Description   string    `gorm:"column:description;type:text" json:"description"`
	Active        *bool     `gorm:"column:active;default:true" json:"active"`
    CreatedDate   time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
    ModifiedDate  time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`
}

func (Dish) TableName() string {
    return "master_dishes"
}
```

#### models/ingredient.go
*Language: Go | Size: 786 bytes*

```go
package models

import "time"

// Ingredient - Master data for raw materials and ingredients (dm_nvl)
type Ingredient struct {
    IngredientID   string    `gorm:"primaryKey;column:ingredient_id" json:"ingredientId"`
    IngredientName string    `gorm:"column:ingredient_name;not null" json:"ingredientName"`
    Property       string    `gorm:"column:properties" json:"property"`
    MaterialGroup  string    `gorm:"column:material_group" json:"materialGroup"`
    Unit           string    `gorm:"column:unit" json:"unit"`
    CreatedDate    time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
    ModifiedDate   time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`
}

func (Ingredient) TableName() string {
    return "master_ingredients"
}
```

#### models/kitchen.go
*Language: Go | Size: 2034 bytes*

```go
package models

import "time"

// Kitchen - Master data for kitchen/location information (dm_bep)
type Kitchen struct {
    KitchenID    string    `gorm:"primaryKey;column:kitchen_id" json:"kitchenId"`
    KitchenName  string    `gorm:"column:kitchen_name;not null" json:"kitchenName"`
    Address      string    `gorm:"column:address;type:text" json:"address"`
    Phone        string    `gorm:"column:phone" json:"phone"`
	Active       *bool     `gorm:"column:active;default:true" json:"active"`
    CreatedDate  time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
    ModifiedDate time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`
}

func (Kitchen) TableName() string {
    return "master_kitchens"
}

// KitchenFavoriteSupplier - Favorite suppliers for each kitchen (kitchen_favorite_suppliers)
type KitchenFavoriteSupplier struct {
	FavoriteID      int       `gorm:"primaryKey;autoIncrement;column:favorite_id" json:"favoriteId"`
	KitchenID       string    `gorm:"column:kitchen_id;type:varchar(50);not null" json:"kitchenId"`
	SupplierID      string    `gorm:"column:supplier_id;type:varchar(50);not null" json:"supplierId"`
	Notes           string    `gorm:"column:notes;type:text" json:"notes"`
	DisplayOrder    *int      `gorm:"column:display_order" json:"displayOrder"`
	CreatedByUserID string    `gorm:"column:created_by_user_id;type:varchar(50)" json:"createdByUserId"`
	CreatedDate     time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate    time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Kitchen       *Kitchen  `gorm:"foreignKey:KitchenID;references:KitchenID" json:"kitchen,omitempty"`
	Supplier      *Supplier `gorm:"foreignKey:SupplierID;references:SupplierID" json:"supplier,omitempty"`
	CreatedBy     *User     `gorm:"foreignKey:CreatedByUserID;references:UserID" json:"createdBy,omitempty"`
}

func (KitchenFavoriteSupplier) TableName() string {
	return "kitchen_favorite_suppliers"
}
```

#### models/order.go
*Language: Go | Size: 7247 bytes*

```go
// models/order.go
package models

import "time"

// Order - Orders (orders)
type Order struct {
	OrderID         string    `gorm:"primaryKey;column:order_id;type:varchar(50)" json:"orderId"`
	KitchenID       string    `gorm:"column:kitchen_id;not null" json:"kitchenId"`
	OrderDate       string    `gorm:"column:order_date;not null" json:"orderDate"`
	Note            string    `gorm:"column:note;type:text" json:"note"`
	Status          string    `gorm:"column:status;default:Pending;not null" json:"status"`
	CreatedByUserID string    `gorm:"column:created_by_user_id" json:"createdByUserId"`
	CreatedDate     time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate    time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Kitchen            *Kitchen                 `gorm:"foreignKey:KitchenID;references:KitchenID" json:"kitchen,omitempty"`
	CreatedBy          *User                    `gorm:"foreignKey:CreatedByUserID;references:UserID" json:"createdBy,omitempty"`
	Details            []OrderDetail            `gorm:"foreignKey:OrderID;references:OrderID" json:"details,omitempty"`
	SupplementaryFoods []OrderSupplementaryFood `gorm:"foreignKey:OrderID;references:OrderID" json:"supplementaryFoods,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

// OrderDetail - Order details (order_details)
type OrderDetail struct {
	OrderDetailID int       `gorm:"primaryKey;autoIncrement;column:order_detail_id" json:"orderDetailId"`
	OrderID       string    `gorm:"column:order_id;type:varchar(50);not null" json:"orderId"`
	DishID        string    `gorm:"column:dish_id;not null" json:"dishId"`
	Portions      int       `gorm:"column:portions;not null" json:"portions"`
	Note          string    `gorm:"column:note;type:text" json:"note"`
	CreatedDate   time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate  time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Order       *Order            `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
	Dish        *Dish             `gorm:"foreignKey:DishID;references:DishID" json:"dish,omitempty"`
	Ingredients []OrderIngredient `gorm:"foreignKey:OrderDetailID;references:OrderDetailID" json:"ingredients,omitempty"`
}

func (OrderDetail) TableName() string {
	return "order_details"
}

// OrderIngredient - Ingredients calculated for an order detail (order_ingredients)
type OrderIngredient struct {
	OrderIngredientID  int       `gorm:"primaryKey;autoIncrement;column:order_ingredient_id" json:"orderIngredientId"`
	OrderDetailID      int       `gorm:"column:order_detail_id;not null" json:"orderDetailId"`
	IngredientID       string    `gorm:"column:ingredient_id;not null" json:"ingredientId"`
	Quantity           float64   `gorm:"column:quantity;type:numeric(15,4);not null" json:"quantity"`
	Unit               string    `gorm:"column:unit;not null" json:"unit"`
	StandardPerPortion float64   `gorm:"column:standard_per_portion;type:numeric(10,4)" json:"standardPerPortion"`
	CreatedDate        time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate       time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	OrderDetail *OrderDetail `gorm:"foreignKey:OrderDetailID;references:OrderDetailID" json:"orderDetail,omitempty"`
	Ingredient  *Ingredient  `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
}

func (OrderIngredient) TableName() string {
	return "order_ingredients"
}

// OrderSupplementaryFood - Extra items for an order (order_supplementary_foods)
type OrderSupplementaryFood struct {
	SupplementaryID    int       `gorm:"primaryKey;autoIncrement;column:supplementary_id" json:"supplementaryId"`
	OrderID            string    `gorm:"column:order_id;type:varchar(50);not null" json:"orderId"`
	IngredientID       string    `gorm:"column:ingredient_id;not null" json:"ingredientId"`
	Quantity           float64   `gorm:"column:quantity;type:numeric(15,4);not null" json:"quantity"`
	Unit               string    `gorm:"column:unit;not null" json:"unit"`
	StandardPerPortion float64   `gorm:"column:standard_per_portion;type:numeric(10,4)" json:"standardPerPortion"`
	Portions           int       `gorm:"column:portions" json:"portions"`
	Note               string    `gorm:"column:note;type:text" json:"note"`
	CreatedDate        time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate       time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Order      *Order      `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
	Ingredient *Ingredient `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
}

func (OrderSupplementaryFood) TableName() string {
	return "order_supplementary_foods"
}

// OrderIngredientSupplier - Selected supplier for each ingredient in an order (order_ingredient_suppliers)
type OrderIngredientSupplier struct {
	OrderIngredientSupplierID int       `gorm:"primaryKey;autoIncrement;column:order_ingredient_supplier_id" json:"orderIngredientSupplierId"`
	OrderID                  string    `gorm:"column:order_id;type:varchar(50);not null" json:"orderId"`
	IngredientID             string    `gorm:"column:ingredient_id;not null" json:"ingredientId"`
	SelectedSupplierID       string    `gorm:"column:selected_supplier_id;not null" json:"selectedSupplierId"`
	SelectedProductID        int       `gorm:"column:selected_product_id;not null" json:"selectedProductId"`
	Quantity                 float64   `gorm:"column:quantity;type:numeric(15,4);not null" json:"quantity"`
	Unit                     string    `gorm:"column:unit;not null" json:"unit"`
	UnitPrice                float64   `gorm:"column:unit_price;type:numeric(15,2);not null" json:"unitPrice"`
	TotalCost                float64   `gorm:"column:total_cost;type:numeric(15,2);not null" json:"totalCost"`
	SelectionDate            time.Time `gorm:"column:selection_date;default:CURRENT_TIMESTAMP" json:"selectionDate"`
	SelectedByUserID         string    `gorm:"column:selected_by_user_id" json:"selectedByUserId"`
	Notes                    string    `gorm:"column:notes;type:text" json:"notes"`
	CreatedDate              time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate             time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Order            *Order        `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
	Ingredient       *Ingredient   `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
	SelectedSupplier *Supplier     `gorm:"foreignKey:SelectedSupplierID;references:SupplierID" json:"selectedSupplier,omitempty"`
	SelectedProduct  *SupplierPrice `gorm:"foreignKey:SelectedProductID;references:ProductID" json:"selectedProduct,omitempty"`
	SelectedBy       *User         `gorm:"foreignKey:SelectedByUserID;references:UserID" json:"selectedBy,omitempty"`
}

func (OrderIngredientSupplier) TableName() string {
	return "order_ingredient_suppliers"
}
```

#### models/order_dto.go
*Language: Go | Size: 2131 bytes*

```go
package models

import "time"

// OrderDTO - Aggregated response for an order
type OrderDTO struct {
	OrderID         string                  `json:"orderId"`
	KitchenID       string                  `json:"kitchenId"`
	KitchenName     string                  `json:"kitchenName"`
	OrderDate       string                  `json:"orderDate"`
	Note            string                  `json:"note"`
	Status          string                  `json:"status"`
	CreatedByUserID string                  `json:"createdByUserId"`
	CreatedByName   string                  `json:"createdByName"`
	CreatedDate     time.Time               `json:"createdDate"`
	ModifiedDate    time.Time               `json:"modifiedDate"`
	Details         []OrderDetailDTO        `json:"details"`
	Supplementaries []OrderSupplementaryDTO `json:"supplementaries"`
}

// OrderDetailDTO - Detail lines with dish name and ingredients
type OrderDetailDTO struct {
	OrderDetailID int                  `json:"orderDetailId"`
	DishID        string               `json:"dishId"`
	DishName      string               `json:"dishName"`
	Portions      int                  `json:"portions"`
	Note          string               `json:"note"`
	Ingredients   []OrderIngredientDTO `json:"ingredients"`
}

// OrderIngredientDTO - Ingredient usage per detail
type OrderIngredientDTO struct {
	OrderIngredientID  int     `json:"orderIngredientId"`
	IngredientID       string  `json:"ingredientId"`
	IngredientName     string  `json:"ingredientName"`
	Quantity           float64 `json:"quantity"`
	Unit               string  `json:"unit"`
	StandardPerPortion float64 `json:"standardPerPortion"`
}

// OrderSupplementaryDTO - Supplementary items for an order
type OrderSupplementaryDTO struct {
	SupplementaryID    int     `json:"supplementaryId"`
	IngredientID       string  `json:"ingredientId"`
	IngredientName     string  `json:"ingredientName"`
	Quantity           float64 `json:"quantity"`
	Unit               string  `json:"unit"`
	StandardPerPortion float64 `json:"standardPerPortion"`
	Portions           int     `json:"portions"`
	Note               string  `json:"note"`
}
```

#### models/pagination.go
*Language: Go | Size: 2231 bytes*

```go
package models

// PaginationParams contains pagination parameters from query string
type PaginationParams struct {
	Page     int    `form:"page" binding:"omitempty,min=0"`
	PageSize int    `form:"per_page" binding:"omitempty,min=0,max=100"` // Changed to per_page to match common conventions
	Search   string `form:"search"`
	SortBy   string `form:"sort_by"`
	SortDir  string `form:"sort_dir" binding:"omitempty,oneof=asc desc"`
}

// PaginationMeta contains pagination metadata matching frontend ResourceCollection interface
type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	From        int `json:"from"`
	To          int `json:"to"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
}

// ResourceCollection is the response wrapper matching frontend interface
type ResourceCollection struct {
	Data interface{}     `json:"data"`
	Meta *PaginationMeta `json:"meta"`
}

// GetPaginationParams extracts and validates pagination parameters with defaults
func GetPaginationParams(page, pageSize int, search, sortBy, sortDir string) PaginationParams {
	// Set defaults
	// if page < 1 {
	// 	page = 1
	// }
	// if pageSize < 1 {
	// 	pageSize = 10
	// }
	// if pageSize > 100 {
	// 	pageSize = 100
	// }
	if sortDir == "" {
		sortDir = "asc"
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
		SortBy:   sortBy,
		SortDir:  sortDir,
	}
}

// CalculatePaginationMeta calculates pagination metadata
func CalculatePaginationMeta(page, perPage int, total int64) *PaginationMeta {
	if perPage < 1 || page < 1 {
		return &PaginationMeta{
			CurrentPage: page,
			LastPage:    1,
			From:        0,
			To:          int(total),
			PerPage:     perPage,
			Total:       int(total),
		}
	}
	totalInt := int(total)
	lastPage := (totalInt + perPage - 1) / perPage
	if lastPage < 1 {
		lastPage = 1
	}

	// Calculate from and to
	from := 0
	to := 0
	if totalInt > 0 {
		from = (page-1)*perPage + 1
		to = from + perPage - 1
		if to > totalInt {
			to = totalInt
		}
	}

	return &PaginationMeta{
		CurrentPage: page,
		LastPage:    lastPage,
		From:        from,
		To:          to,
		PerPage:     perPage,
		Total:       totalInt,
	}
}
```

#### models/recipe_standard.go
*Language: Go | Size: 1319 bytes*

```go
package models

import "time"

// RecipeStandard - Bill of materials for dishes (dish_recipe_standards)
type RecipeStandard struct {
	StandardID   int       `gorm:"primaryKey;autoIncrement;column:recipe_id" json:"standardId"`
	DishID       string    `gorm:"column:dish_id" json:"dishId"`
	IngredientID string    `gorm:"column:ingredient_id" json:"ingredientId"`
	Unit         string    `gorm:"column:unit" json:"unit"`
	StandardPer1 float64   `gorm:"column:quantity_per_serving;type:decimal(10,4)" json:"standardPer1"`
	Note         string    `gorm:"column:notes;type:text" json:"note"`
	Amount       float64   `gorm:"column:cost;type:decimal(15,2)" json:"amount"`
	UpdatedByID  string    `gorm:"column:updated_by_user_id" json:"updatedById"`
	CreatedDate  time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Dish       *Dish       `gorm:"foreignKey:DishID;references:DishID" json:"dish,omitempty"`
	Ingredient *Ingredient `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
	UpdatedBy  *User       `gorm:"foreignKey:UpdatedByID;references:UserID" json:"updatedBy,omitempty"`
}

func (RecipeStandard) TableName() string {
	return "dish_recipe_standards"
}
```

#### models/recipe_standard_dto.go
*Language: Go | Size: 1819 bytes*

```go
package models

import "time"

// RecipeStandardDTO - Data Transfer Object for Recipe Standard with related names
type RecipeStandardDTO struct {
	StandardID     int       `json:"standardId"`
	DishID         string    `json:"dishId"`
	DishName       string    `json:"dishName"`           // Added: Dish name
	IngredientID   string    `json:"ingredientId"`
	IngredientName string    `json:"ingredientName"`     // Added: Ingredient name
	Unit           string    `json:"unit"`
	StandardPer1   float64   `json:"standardPer1"`
	Note           string    `json:"note"`
	Amount         float64   `json:"amount"`
	UpdatedByID    string    `json:"updatedById"`
	UpdatedByName  string    `json:"updatedByName"`      // Added: User name (optional)
	CreatedDate    time.Time `json:"createdDate"`
	ModifiedDate   time.Time `json:"modifiedDate"`
}

// ToDTO converts RecipeStandard model to DTO
func (r *RecipeStandard) ToDTO() RecipeStandardDTO {
	dto := RecipeStandardDTO{
		StandardID:   r.StandardID,
		DishID:       r.DishID,
		IngredientID: r.IngredientID,
		Unit:         r.Unit,
		StandardPer1: r.StandardPer1,
		Note:         r.Note,
		Amount:       r.Amount,
		UpdatedByID:  r.UpdatedByID,
		CreatedDate:  r.CreatedDate,
		ModifiedDate: r.ModifiedDate,
	}

	// Populate names from relationships if available
	if r.Dish != nil {
		dto.DishName = r.Dish.DishName
	}
	if r.Ingredient != nil {
		dto.IngredientName = r.Ingredient.IngredientName
	}
	if r.UpdatedBy != nil {
		dto.UpdatedByName = r.UpdatedBy.FullName
	}

	return dto
}

// ConvertToDTO converts a slice of RecipeStandard to a slice of RecipeStandardDTO
func ConvertRecipeStandardsToDTO(recipes []RecipeStandard) []RecipeStandardDTO {
	dtos := make([]RecipeStandardDTO, len(recipes))
	for i, recipe := range recipes {
		dtos[i] = recipe.ToDTO()
	}
	return dtos
}
```

#### models/supplier.go
*Language: Go | Size: 877 bytes*

```go
package models

import "time"

// Supplier - Master data for suppliers (dm_ncc)
type Supplier struct {
    SupplierID   string    `gorm:"primaryKey;column:supplier_id" json:"supplierId"`
    SupplierName string    `gorm:"column:supplier_name;not null" json:"supplierName"`
    ZaloLink     string    `gorm:"column:zalo_link;type:text" json:"zaloLink"`
    Address      string    `gorm:"column:address;type:text" json:"address"`
    Phone        string    `gorm:"column:phone" json:"phone"`
	Email        string    `gorm:"column:email" json:"email"`
	Active       *bool     `gorm:"column:active;default:true" json:"active"`
    CreatedDate  time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
    ModifiedDate time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`
}

func (Supplier) TableName() string {
    return "master_suppliers"
}
```

#### models/supplier_price.go
*Language: Go | Size: 1800 bytes*

```go
package models

import "time"

// SupplierPrice - Supplier price list (supplier_price_list)
type SupplierPrice struct {
	ProductID     int        `gorm:"primaryKey;autoIncrement;column:product_id" json:"productId"`
	ProductName   string     `gorm:"column:product_name" json:"productName"`
	IngredientID  string     `gorm:"column:ingredient_id" json:"ingredientId"`
	Category      string     `gorm:"column:classification" json:"category"`
	SupplierID    string     `gorm:"column:supplier_id" json:"supplierId"`
	Manufacturer  string     `gorm:"column:manufacturer_name" json:"manufacturer"`
	Unit          string     `gorm:"column:unit" json:"unit"`
	Specification string     `gorm:"column:specification" json:"specification"`
	UnitPrice     float64    `gorm:"column:unit_price;type:decimal(15,2)" json:"unitPrice"`
	PricePer1     float64    `gorm:"column:price_per_item;type:decimal(15,2)" json:"pricePer1"`
	EffectiveFrom *time.Time `gorm:"column:effective_from" json:"effectiveFrom"`
	EffectiveTo   *time.Time `gorm:"column:effective_to" json:"effectiveTo"`
	Active        *bool      `gorm:"column:active;default:true" json:"active"`
	NewPrice      float64    `gorm:"column:new_buying_price;type:decimal(15,2)" json:"newPrice"`
	Promotion     string     `gorm:"column:promotion;type:char(1)" json:"promotion"`
	CreatedDate   time.Time  `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	ModifiedDate  time.Time  `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`

	// Relationships
	Ingredient *Ingredient `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient,omitempty"`
	Supplier   *Supplier   `gorm:"foreignKey:SupplierID;references:SupplierID" json:"supplier,omitempty"`
}

func (SupplierPrice) TableName() string {
	return "supplier_price_list"
}
```

#### models/supplier_price_dto.go
*Language: Go | Size: 2175 bytes*

```go
package models

import "time"

// SupplierPriceDTO - Data Transfer Object for Supplier Price with related names
type SupplierPriceDTO struct {
	ProductID        int        `json:"productId"`
	ProductName      string     `json:"productName"`
	IngredientID     string     `json:"ingredientId"`
	IngredientName   string     `json:"ingredientName"`   // Ingredient name from relationship
	Category         string     `json:"category"`
	SupplierID       string     `json:"supplierId"`
	SupplierName     string     `json:"supplierName"`     // Supplier name from relationship
	Manufacturer     string     `json:"manufacturer"`
	Unit             string     `json:"unit"`
	Specification    string     `json:"specification"`
	UnitPrice        float64    `json:"unitPrice"`
	PricePer1        float64    `json:"pricePer1"`
	EffectiveFrom    *time.Time `json:"effectiveFrom"`
	EffectiveTo      *time.Time `json:"effectiveTo"`
	Active           *bool      `json:"active"`
	NewPrice         float64    `json:"newPrice"`
	Promotion        string     `json:"promotion"`
}

// ToDTO converts SupplierPrice model to DTO
func (sp *SupplierPrice) ToDTO() SupplierPriceDTO {
	dto := SupplierPriceDTO{
		ProductID:     sp.ProductID,
		ProductName:   sp.ProductName,
		IngredientID:  sp.IngredientID,
		Category:      sp.Category,
		SupplierID:    sp.SupplierID,
		Manufacturer:  sp.Manufacturer,
		Unit:          sp.Unit,
		Specification: sp.Specification,
		UnitPrice:     sp.UnitPrice,
		PricePer1:     sp.PricePer1,
		EffectiveFrom: sp.EffectiveFrom,
		EffectiveTo:   sp.EffectiveTo,
		Active:        sp.Active,
		NewPrice:      sp.NewPrice,
		Promotion:     sp.Promotion,
	}

	// Populate names from relationships if available
	if sp.Ingredient != nil {
		dto.IngredientName = sp.Ingredient.IngredientName
	}
	if sp.Supplier != nil {
		dto.SupplierName = sp.Supplier.SupplierName
	}

	return dto
}

// ConvertSupplierPricesToDTO converts a slice of SupplierPrice to a slice of SupplierPriceDTO
func ConvertSupplierPricesToDTO(prices []SupplierPrice) []SupplierPriceDTO {
	dtos := make([]SupplierPriceDTO, len(prices))
	for i, price := range prices {
		dtos[i] = price.ToDTO()
	}
	return dtos
}
```

#### models/user.go
*Language: Go | Size: 1136 bytes*

```go
package models

import (
	"time"
)






// User - Master data for user accounts (dm_nguoidung)
type User struct {
    UserID       string    `gorm:"primaryKey;column:user_id" json:"userId"`
    UserName     string    `gorm:"column:user_name;not null;unique" json:"userName"`
	Password     string    `gorm:"column:password;not null" json:"password,omitempty"`
    FullName     string    `gorm:"column:full_name;not null" json:"fullName"`
    Role         string    `gorm:"column:role" json:"role"`
    KitchenID    string    `gorm:"column:kitchen_id" json:"kitchenId"`
	Email        string    `gorm:"column:email" json:"email"`
    Phone        string    `gorm:"column:phone" json:"phone"`
	Active       *bool     `gorm:"column:active;default:true" json:"active"`
    CreatedDate  time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
    ModifiedDate time.Time `gorm:"column:modified_date;autoUpdateTime" json:"modifiedDate"`
	
	// Relationships
	Kitchen      *Kitchen  `gorm:"foreignKey:KitchenID;references:KitchenID" json:"kitchen,omitempty"`
}

func (User) TableName() string {
    return "master_users"
}






```

#### .env
*Language: Text | Size: 119 bytes*

```text
DATABASE_URL="host=14.225.198.206 user=adong password=adong123 dbname=adongfoodv2 port=5432 sslmode=disable"
PORT=18080
```

#### .gitignore
*Language: Text | Size: 15 bytes*

```text
backend.md
.env
```

#### Dockerfile
*Language: Text | Size: 1121 bytes*

```text
# ============================================
# Stage 1: Build Stage
# ============================================
FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN --mount=type=bind,target=/app \
    go build -o /tmp/main cmd/main.go

# ============================================
# Stage 2: Runtime Stage
# ============================================
FROM alpine:latest

# Set timezone
ENV TZ=Asia/Ho_Chi_Minh

# Create app user for security (non-root)
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder --chown=appuser:appuser /tmp/main .

# Copy .env file if exists (optional - can use env vars instead)
COPY --from=builder --chown=appuser:appuser /app/.env* ./

# Set ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose application port
EXPOSE 18080

# Health check
HEALTHCHECK --interval=30s \
    --timeout=10s \
    --start-period=5s \
    --retries=3 \
    CMD curl -f http://localhost:18080/health || exit 1

# Run the application
CMD ["./main"]
```

#### README.md
*Language: Markdown | Size: 21 bytes*

```markdown
# adong-food-backend
```

### server/

#### server/router.go
*Language: Go | Size: 8746 bytes*

```go
package server

import (
	"adong-be/handler"
	"adong-be/logger"
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
		// Master data routes
		api.GET("/ingredients", handler.GetIngredients)
		api.GET("/ingredients/:id", handler.GetIngredient)
		api.POST("/ingredients", handler.CreateIngredient)
		api.PUT("/ingredients/:id", handler.UpdateIngredient)
		api.DELETE("/ingredients/:id", handler.DeleteIngredient)

		api.GET("/kitchens", handler.GetKitchens)
		api.POST("/kitchens", handler.CreateKitchen)
		
		// Kitchen favorite suppliers (must be before /kitchens/:id to avoid route conflict)
		api.GET("/kitchens/:id/favorite-suppliers", handler.GetKitchenFavoriteSuppliers)
		api.GET("/kitchens/:id/favorite-suppliers/:favoriteId", handler.GetKitchenFavoriteSupplier)
		api.POST("/kitchens/:id/favorite-suppliers", handler.CreateKitchenFavoriteSupplier)
		api.PUT("/kitchens/:id/favorite-suppliers/:favoriteId", handler.UpdateKitchenFavoriteSupplier)
		api.DELETE("/kitchens/:id/favorite-suppliers/:favoriteId", handler.DeleteKitchenFavoriteSupplier)
		
		api.GET("/kitchens/:id", handler.GetKitchen)
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
		api.GET("/orders/:id/ingredients/summary", handler.GetOrderIngredientsSummary)
		api.GET("/orders/:id/ingredients/:ingredientId/summary", handler.GetOrderIngredientSummary)
		api.POST("/orders", handler.CreateOrder)
		api.POST("/orders/:id/supplier-requests", handler.SaveOrderIngredientsWithSupplier)
		api.PATCH("/orders/:id/status", handler.UpdateOrderStatus)
		api.DELETE("/orders/:id", handler.DeleteOrder)

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
```

### store/

#### store/gorm.go
*Language: Go | Size: 1630 bytes*

```go
package store

import (
	"adong-be/models"
	"time"

	"github.com/hsdfat/go-auth-middleware/core"
	"gorm.io/gorm"
)

type Store struct {
	GormClient *gorm.DB
}

var DB *Store = &Store{}

// Enhanced UserProvider interface
type UserProvider interface {
	GetUserByUsername(username string) (*core.User, error)
	GetUserByID(userID string) (*core.User, error)
	GetUserByEmail(email string) (*core.User, error)
	UpdateUserLastLogin(userID string, lastLogin time.Time) error
	IsUserActive(userID string) (bool, error)
}

func (s *Store) GetUserByUsername(username string) (*core.User, error) {
	var dbUser models.User
	if err := s.GormClient.First(&dbUser, "user_name = ?", username).Error; err != nil {
		return nil, err
	}
	return convertToCoreUser(dbUser), nil
}

func (s *Store) GetUserByID(userID string) (*core.User, error) {
	var dbUser models.User
	if err := s.GormClient.First(&dbUser, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return convertToCoreUser(dbUser), nil
}

func (s *Store) GetUserByEmail(email string) (*core.User, error) {
	var user models.User
	if err := s.GormClient.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return convertToCoreUser(user), nil
}

func (s *Store) UpdateUserLastLogin(userID string, lastLogin time.Time) error {
	return nil
}

func (s *Store) IsUserActive(userID string) (bool, error) {
	return true, nil
}

func convertToCoreUser(dbUser models.User) *core.User {
	return &core.User{
		ID: dbUser.UserID,
		Username: dbUser.UserID,
		Email:    dbUser.Email,
		Password: dbUser.Password,
		Role:     dbUser.Role,
		IsActive: true,
	}
}
```

### utils/

#### utils/search.go
*Language: Go | Size: 1635 bytes*

```go
package utils

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// SearchConfig defines which fields to search and their weight
type SearchConfig struct {
	Fields []string
	Fuzzy  bool // Use ILIKE (case-insensitive) vs exact match
}

// ApplySearch applies search conditions to a GORM query
func ApplySearch(db *gorm.DB, search string, config SearchConfig) *gorm.DB {
	if search == "" || len(config.Fields) == 0 {
		return db
	}

	// Trim and prepare search term
	search = strings.TrimSpace(search)

	// Build search conditions
	var conditions []string
	var args []interface{}

	for _, field := range config.Fields {
		if config.Fuzzy {
			conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", field))
			args = append(args, "%"+search+"%")
		} else {
			conditions = append(conditions, fmt.Sprintf("%s = ?", field))
			args = append(args, search)
		}
	}

	// Combine with OR
	query := strings.Join(conditions, " OR ")
	return db.Where(query, args...)
}

// ApplySort applies sorting to a GORM query
func ApplySort(db *gorm.DB, sortBy, sortDir string, allowedFields map[string]string) *gorm.DB {
	if sortBy == "" {
		return db
	}

	// Validate sort field
	dbField, ok := allowedFields[sortBy]
	if !ok {
		return db
	}

	// Validate sort direction
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "asc"
	}

	return db.Order(fmt.Sprintf("%s %s", dbField, strings.ToUpper(sortDir)))
}

// ApplyPagination applies pagination to a GORM query
func ApplyPagination(db *gorm.DB, page, pageSize int) *gorm.DB {
	if page < 1 || pageSize < 1 {
		return db
	}
	offset := (page - 1) * pageSize
	return db.Offset(offset).Limit(pageSize)
}
```

