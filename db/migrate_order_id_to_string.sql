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

