-- =========================================================
-- SCHEMA: CENTRAL KITCHEN MANAGEMENT
-- =========================================================

-- ========================
-- MASTER TABLES
-- ========================

CREATE TABLE IF NOT EXISTS public.master_suppliers (
    supplier_id VARCHAR(50) PRIMARY KEY,
    supplier_name VARCHAR(255) NOT NULL,
    zalo_link TEXT,
    address TEXT,
    phone VARCHAR(20),
    email VARCHAR(255),
    active BOOLEAN DEFAULT TRUE NOT NULL,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS public.master_kitchens (
    kitchen_id VARCHAR(50) PRIMARY KEY,
    kitchen_name VARCHAR(255) NOT NULL,
    address TEXT,
    phone VARCHAR(20),
    active BOOLEAN DEFAULT TRUE NOT NULL,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS public.master_users (
    user_id VARCHAR(50) PRIMARY KEY,
    user_name VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role VARCHAR(50),
    kitchen_id VARCHAR(50),
    email VARCHAR(255),
    phone VARCHAR(20),
    active BOOLEAN DEFAULT TRUE NOT NULL,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_users_kitchen FOREIGN KEY (kitchen_id)
        REFERENCES public.master_kitchens(kitchen_id)
        ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS public.master_dishes (
    dish_id VARCHAR(50) PRIMARY KEY,
    dish_name VARCHAR(255) NOT NULL UNIQUE,
    cooking_method VARCHAR(100),
    category VARCHAR(100),
    description TEXT,
    active BOOLEAN DEFAULT TRUE NOT NULL,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS public.master_ingredients (
    ingredient_id VARCHAR(50) PRIMARY KEY,
    ingredient_name VARCHAR(255) NOT NULL UNIQUE,
    properties VARCHAR(100),
    material_group VARCHAR(100),
    unit VARCHAR(50) NOT NULL,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- ========================
-- SUPPLIER PRICE LIST
-- ========================

CREATE TABLE IF NOT EXISTS public.supplier_price_list (
    product_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    product_name VARCHAR(255),
    ingredient_id VARCHAR(50) NOT NULL,
    classification VARCHAR(100),
    supplier_id VARCHAR(50) NOT NULL,
    manufacturer_name VARCHAR(255),
    unit VARCHAR(50),
    specification VARCHAR(100),
    unit_price NUMERIC(15,2) NOT NULL CHECK (unit_price >= 0),
    price_per_item NUMERIC(15,2),
    effective_from TIMESTAMP,
    effective_to TIMESTAMP,
    active BOOLEAN DEFAULT TRUE NOT NULL,
    new_buying_price NUMERIC(15,2),
    promotion CHAR(1),
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_price_ingredient FOREIGN KEY (ingredient_id)
        REFERENCES public.master_ingredients(ingredient_id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_price_supplier FOREIGN KEY (supplier_id)
        REFERENCES public.master_suppliers(supplier_id)
        ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_supplier_price_ingredient ON public.supplier_price_list (ingredient_id);
CREATE INDEX IF NOT EXISTS idx_supplier_price_supplier ON public.supplier_price_list (supplier_id);

-- ========================
-- DISH RECIPE STANDARDS
-- ========================

CREATE TABLE IF NOT EXISTS public.dish_recipe_standards (
    recipe_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    dish_id VARCHAR(50) NOT NULL,
    ingredient_id VARCHAR(50) NOT NULL,
    unit VARCHAR(50) NOT NULL,
    quantity_per_serving NUMERIC(10,4) NOT NULL CHECK (quantity_per_serving > 0),
    notes TEXT,
    cost NUMERIC(15,2),
    updated_by_user_id VARCHAR(50),
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_recipe_dish FOREIGN KEY (dish_id)
        REFERENCES public.master_dishes(dish_id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_recipe_ingredient FOREIGN KEY (ingredient_id)
        REFERENCES public.master_ingredients(ingredient_id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_recipe_user FOREIGN KEY (updated_by_user_id)
        REFERENCES public.master_users(user_id)
        ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_recipe_dish ON public.dish_recipe_standards (dish_id);
CREATE INDEX IF NOT EXISTS idx_recipe_ingredient ON public.dish_recipe_standards (ingredient_id);

-- ========================
-- ORDERS AND DETAILS
-- ========================

CREATE TABLE IF NOT EXISTS public.orders (
    order_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    kitchen_id VARCHAR(50) NOT NULL,
    order_date DATE NOT NULL,
    note TEXT,
    status VARCHAR(50) DEFAULT 'Pending' NOT NULL,
    created_by_user_id VARCHAR(50),
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_order_kitchen FOREIGN KEY (kitchen_id)
        REFERENCES public.master_kitchens(kitchen_id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_order_user FOREIGN KEY (created_by_user_id)
        REFERENCES public.master_users(user_id)
        ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS public.order_details (
    order_detail_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id INTEGER NOT NULL,
    dish_id VARCHAR(50) NOT NULL,
    portions INTEGER NOT NULL CHECK (portions > 0),
    note TEXT,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_detail_order FOREIGN KEY (order_id)
        REFERENCES public.orders(order_id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_detail_dish FOREIGN KEY (dish_id)
        REFERENCES public.master_dishes(dish_id)
        ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS public.order_ingredients (
    order_ingredient_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_detail_id INTEGER NOT NULL,
    ingredient_id VARCHAR(50) NOT NULL,
    quantity NUMERIC(15,4) NOT NULL CHECK (quantity > 0),
    unit VARCHAR(50) NOT NULL,
    standard_per_portion NUMERIC(10,4),
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_order_ing_detail FOREIGN KEY (order_detail_id)
        REFERENCES public.order_details(order_detail_id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_order_ing_ingredient FOREIGN KEY (ingredient_id)
        REFERENCES public.master_ingredients(ingredient_id)
        ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS public.order_supplementary_foods (
    supplementary_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id INTEGER NOT NULL,
    ingredient_id VARCHAR(50) NOT NULL,
    quantity NUMERIC(15,4) NOT NULL CHECK (quantity > 0),
    unit VARCHAR(50) NOT NULL,
    standard_per_portion NUMERIC(10,4),
    portions INTEGER,
    note TEXT,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_supp_order FOREIGN KEY (order_id)
        REFERENCES public.orders(order_id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_supp_ingredient FOREIGN KEY (ingredient_id)
        REFERENCES public.master_ingredients(ingredient_id)
        ON UPDATE CASCADE ON DELETE RESTRICT
);

-- ========================
-- SUPPLIER REQUESTS
-- ========================

CREATE TABLE IF NOT EXISTS public.supplier_requests (
    request_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id INTEGER NOT NULL,
    supplier_id VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'Pending' NOT NULL,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_req_order FOREIGN KEY (order_id)
        REFERENCES public.orders(order_id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_req_supplier FOREIGN KEY (supplier_id)
        REFERENCES public.master_suppliers(supplier_id)
        ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS public.supplier_request_details (
    request_detail_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    request_id INTEGER NOT NULL,
    ingredient_id VARCHAR(50) NOT NULL,
    quantity NUMERIC(15,4) NOT NULL CHECK (quantity > 0),
    unit VARCHAR(50) NOT NULL,
    unit_price NUMERIC(15,2) NOT NULL CHECK (unit_price >= 0),
    total_price NUMERIC(15,2) GENERATED ALWAYS AS (quantity * unit_price) STORED,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_req_detail_request FOREIGN KEY (request_id)
        REFERENCES public.supplier_requests(request_id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_req_detail_ingredient FOREIGN KEY (ingredient_id)
        REFERENCES public.master_ingredients(ingredient_id)
        ON UPDATE CASCADE ON DELETE RESTRICT
);

-- ========================
-- INDEXES (PERFORMANCE)
-- ========================

CREATE INDEX IF NOT EXISTS idx_orders_kitchen ON public.orders (kitchen_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON public.orders (status);
CREATE INDEX IF NOT EXISTS idx_order_details_order ON public.order_details (order_id);
CREATE INDEX IF NOT EXISTS idx_order_ing_detail ON public.order_ingredients (order_detail_id);
CREATE INDEX IF NOT EXISTS idx_supplementary_order ON public.order_supplementary_foods (order_id);
CREATE INDEX IF NOT EXISTS idx_supplier_requests_order ON public.supplier_requests (order_id);
CREATE INDEX IF NOT EXISTS idx_supplier_requests_supplier ON public.supplier_requests (supplier_id);
