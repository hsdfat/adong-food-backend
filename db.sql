CREATE TABLE public.master_suppliers (
    supplier_id character varying(50) NOT NULL,
    supplier_name character varying(255) NOT NULL,
    zalo_link text,
    address text,
    phone character varying(20),
    email character varying(255),
    active boolean DEFAULT true,
    created_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    modified_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE public.master_suppliers OWNER TO adong;

--
-- Table: master_users
--
CREATE TABLE public.master_users (
    user_id character varying(50) NOT NULL,
    user_name character varying(50) NOT NULL,
    password character varying(255) NOT NULL,
    full_name character varying(255) NOT NULL,
    role character varying(50),
    kitchen_id character varying(50),
    email character varying(255),
    phone character varying(20),
    active boolean DEFAULT true,
    created_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    modified_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE public.master_users OWNER TO adong;
--
-- Table: master_ingredients
--
CREATE TABLE public.master_ingredients (
    ingredient_id character varying(50) NOT NULL,
    ingredient_name character varying(255) NOT NULL,
    properties character varying(100),
    material_group character varying(100),
    unit character varying(50),
    created_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    modified_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE public.master_ingredients OWNER TO adong;
-- Table: master_kitchens
--
CREATE TABLE public.master_kitchens (
    kitchen_id character varying(50) NOT NULL,
    kitchen_name character varying(255) NOT NULL,
    address text,
    phone character varying(20),
    active boolean DEFAULT true,
    created_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    modified_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE public.master_kitchens OWNER TO adong;

-- Table: supplementary_items
--

--
-- Table: supplier_price_list
--
CREATE TABLE public.supplier_price_list (
    product_id integer NOT NULL,
    product_name character varying(255),
    ingredient_id character varying(50),
    classification character varying(100),
    supplier_id character varying(50),
    manufacturer_name character varying(255),
    unit character varying(50),
    specification character varying(100),
    unit_price numeric(15,2),
    price_per_item numeric(15,2),
    effective_from timestamp without time zone,
    effective_to timestamp without time zone,
    active boolean DEFAULT true,
    new_buying_price numeric(15,2),
    promotion "char"
);

ALTER TABLE public.supplier_price_list OWNER TO adong;

CREATE SEQUENCE public.supplier_price_list_product_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.supplier_price_list_product_id_seq OWNER TO adong;
ALTER SEQUENCE public.supplier_price_list_product_id_seq OWNED BY public.supplier_price_list.product_id;

--
-- Table: dish_recipe_standards
--
CREATE TABLE public.dish_recipe_standards (
    recipe_id integer NOT NULL,
    dish_id character varying(50),
    ingredient_id character varying(50),
    unit character varying(50),
    quantity_per_serving numeric(10,4),
    notes text,
    cost numeric(15,2),
    updated_by_user_id character varying(50),
    created_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    modified_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE public.dish_recipe_standards OWNER TO adong;

CREATE SEQUENCE public.dish_recipe_standards_recipe_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.dish_recipe_standards_recipe_id_seq OWNER TO adong;
ALTER SEQUENCE public.dish_recipe_standards_recipe_id_seq OWNED BY public.dish_recipe_standards.recipe_id;

--
-- Constraints
--
ALTER TABLE ONLY public.supplier_price_list
    ADD CONSTRAINT supplier_price_list_pkey PRIMARY KEY (product_id);

ALTER TABLE ONLY public.master_kitchens
    ADD CONSTRAINT master_kitchens_pkey PRIMARY KEY (kitchen_id);

ALTER TABLE ONLY public.master_dishes
    ADD CONSTRAINT master_dishes_pkey PRIMARY KEY (dish_id);

ALTER TABLE ONLY public.master_suppliers
    ADD CONSTRAINT master_suppliers_pkey PRIMARY KEY (supplier_id);

ALTER TABLE ONLY public.master_users
    ADD CONSTRAINT master_users_pkey PRIMARY KEY (user_id);

ALTER TABLE ONLY public.master_ingredients
    ADD CONSTRAINT master_ingredients_pkey PRIMARY KEY (ingredient_id);

ALTER TABLE ONLY public.dish_recipe_standards
    ADD CONSTRAINT dish_recipe_standards_pkey PRIMARY KEY (recipe_id);
--
-- Indexes
--
CREATE INDEX idx_supplier_price_list_supplier ON public.supplier_price_list USING btree (supplier_id);
CREATE INDEX idx_supplier_price_list_ingredient ON public.supplier_price_list USING btree (ingredient_id);
CREATE INDEX idx_dish_standards_dish ON public.dish_recipe_standards USING btree (dish_id);
CREATE INDEX idx_dish_standards_ingredient ON public.dish_recipe_standards USING btree (ingredient_id);

CREATE TRIGGER update_order_forms_updated_at BEFORE UPDATE ON public.order_forms FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

--
-- Foreign Keys
--
ALTER TABLE ONLY public.supplier_price_list
    ADD CONSTRAINT supplier_price_list_ingredient_fkey FOREIGN KEY (ingredient_id) REFERENCES public.master_ingredients(ingredient_id);

ALTER TABLE ONLY public.supplier_price_list
    ADD CONSTRAINT supplier_price_list_supplier_fkey FOREIGN KEY (supplier_id) REFERENCES public.master_suppliers(supplier_id);

ALTER TABLE ONLY public.dish_recipe_standards
    ADD CONSTRAINT recipe_standards_dish_fkey FOREIGN KEY (dish_id) REFERENCES public.master_dishes(dish_id);

ALTER TABLE ONLY public.dish_recipe_standards
    ADD CONSTRAINT recipe_standards_updated_by_fkey FOREIGN KEY (updated_by_user_id) REFERENCES public.master_users(user_id);

ALTER TABLE ONLY public.dish_recipe_standards
    ADD CONSTRAINT recipe_standards_ingredient_fkey FOREIGN KEY (ingredient_id) REFERENCES public.master_ingredients(ingredient_id);

ALTER TABLE ONLY public.master_users
    ADD CONSTRAINT users_kitchen_fkey FOREIGN KEY (kitchen_id) REFERENCES public.master_kitchens(kitchen_id);


-- Completed on 2025-11-02 14:25:51 UTC
