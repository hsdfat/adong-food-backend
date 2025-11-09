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