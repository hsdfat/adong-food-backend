-- ============================================================================
-- SAMPLE DATA FOR ADONG FOOD MANAGEMENT SYSTEM
-- ============================================================================
-- This file contains comprehensive sample data for testing and demonstration
-- Includes: 25 ingredients, 15 dishes, 7 suppliers, 3 kitchens, users, recipes, and orders
-- ============================================================================

BEGIN;

-- ============================================================================
-- INGREDIENT TYPES
-- ============================================================================
INSERT INTO public.ingredient_types (ingredient_type_id, ingredient_type_name, description, active) VALUES
('MEAT', 'Thịt', 'Các loại thịt động vật', true),
('SEAFOOD', 'Hải sản', 'Các loại hải sản tươi sống', true),
('VEGETABLE', 'Rau củ', 'Rau và củ quả tươi', true),
('SPICE', 'Gia vị', 'Các loại gia vị và hương liệu', true),
('GRAIN', 'Ngũ cốc', 'Gạo, bột và ngũ cốc', true),
('DAIRY', 'Sữa và trứng', 'Các sản phẩm từ sữa và trứng', true),
('OIL', 'Dầu mỡ', 'Các loại dầu ăn và mỡ', true);

-- ============================================================================
-- MASTER INGREDIENTS (25 ingredients)
-- ============================================================================
INSERT INTO public.master_ingredients (ingredient_id, ingredient_name, ingredient_type_id, properties, material_group, unit) VALUES
-- Thịt (5 items)
('ING001', 'Thịt heo ba chỉ', 'MEAT', 'Tươi', 'Thịt heo', 'kg'),
('ING002', 'Thịt bò bắp', 'MEAT', 'Tươi', 'Thịt bò', 'kg'),
('ING003', 'Thịt gà ta', 'MEAT', 'Tươi', 'Thịt gia cầm', 'kg'),
('ING004', 'Sườn heo', 'MEAT', 'Tươi', 'Thịt heo', 'kg'),
('ING005', 'Thịt vịt', 'MEAT', 'Tươi', 'Thịt gia cầm', 'kg'),

-- Hải sản (5 items)
('ING006', 'Tôm sú', 'SEAFOOD', 'Tươi sống', 'Giáp xác', 'kg'),
('ING007', 'Cá lóc', 'SEAFOOD', 'Tươi sống', 'Cá nước ngọt', 'kg'),
('ING008', 'Mực ống', 'SEAFOOD', 'Tươi', 'Động vật thân mềm', 'kg'),
('ING009', 'Nghêu', 'SEAFOOD', 'Tươi sống', 'Động vật có vỏ', 'kg'),
('ING010', 'Cá basa', 'SEAFOOD', 'Tươi', 'Cá nước ngọt', 'kg'),

-- Rau củ (6 items)
('ING011', 'Cà chua', 'VEGETABLE', 'Tươi', 'Củ quả', 'kg'),
('ING012', 'Hành tây', 'VEGETABLE', 'Tươi', 'Củ', 'kg'),
('ING013', 'Khoai tây', 'VEGETABLE', 'Tươi', 'Củ', 'kg'),
('ING014', 'Rau muống', 'VEGETABLE', 'Tươi', 'Rau xanh', 'kg'),
('ING015', 'Cải bắp', 'VEGETABLE', 'Tươi', 'Rau xanh', 'kg'),
('ING016', 'Ớt', 'VEGETABLE', 'Tươi', 'Gia vị tươi', 'kg'),

-- Gia vị (4 items)
('ING017', 'Nước mắm', 'SPICE', 'Đóng chai', 'Nước chấm', 'lít'),
('ING018', 'Dầu ăn', 'OIL', 'Đóng chai', 'Dầu thực vật', 'lít'),
('ING019', 'Hạt tiêu', 'SPICE', 'Khô', 'Gia vị khô', 'kg'),
('ING020', 'Muối', 'SPICE', 'Khô', 'Gia vị cơ bản', 'kg'),

-- Ngũ cốc (3 items)
('ING021', 'Gạo tẻ', 'GRAIN', 'Khô', 'Gạo', 'kg'),
('ING022', 'Bún tươi', 'GRAIN', 'Tươi', 'Bún phở', 'kg'),
('ING023', 'Bánh phở', 'GRAIN', 'Tươi', 'Bún phở', 'kg'),

-- Sữa và trứng (2 items)
('ING024', 'Trứng gà', 'DAIRY', 'Tươi', 'Trứng', 'quả'),
('ING025', 'Sữa tươi', 'DAIRY', 'Tươi', 'Sữa', 'lít');

-- ============================================================================
-- MASTER KITCHENS (3 kitchens)
-- ============================================================================
INSERT INTO public.master_kitchens (kitchen_id, kitchen_name, address, phone, active) VALUES
('KIT001', 'Bếp Trung Tâm Á Đông', '123 Đường Lê Lợi, Quận 1, TP.HCM', '0283456789', true),
('KIT002', 'Bếp Chi Nhánh Quận 7', '456 Đường Nguyễn Văn Linh, Quận 7, TP.HCM', '0287654321', true),
('KIT003', 'Bếp Chi Nhánh Thủ Đức', '789 Đường Võ Văn Ngân, Thủ Đức, TP.HCM', '0289876543', true);

-- ============================================================================
-- MASTER USERS (5 users)
-- ============================================================================
INSERT INTO public.master_users (user_id, user_name, password, full_name, role, kitchen_id, email, phone, active) VALUES
('USR001', 'admin', '$2a$10$rF8Z9K1qXJx5yGqL3pQy0.KJhL7vXmZ9F3uY5tKnM8wQxPzN4bC2K', 'Nguyễn Văn An', 'Admin', NULL, 'admin@adongfood.vn', '0901234567', true),
('USR002', 'chef_k001', '$2a$10$rF8Z9K1qXJx5yGqL3pQy0.KJhL7vXmZ9F3uY5tKnM8wQxPzN4bC2K', 'Trần Thị Bình', 'Chef', 'KIT001', 'binh@adongfood.vn', '0901234568', true),
('USR003', 'chef_k002', '$2a$10$rF8Z9K1qXJx5yGqL3pQy0.KJhL7vXmZ9F3uY5tKnM8wQxPzN4bC2K', 'Lê Văn Cường', 'Chef', 'KIT002', 'cuong@adongfood.vn', '0901234569', true),
('USR004', 'manager', '$2a$10$rF8Z9K1qXJx5yGqL3pQy0.KJhL7vXmZ9F3uY5tKnM8wQxPzN4bC2K', 'Phạm Thị Dung', 'Manager', 'KIT001', 'dung@adongfood.vn', '0901234570', true),
('USR005', 'staff_k003', '$2a$10$rF8Z9K1qXJx5yGqL3pQy0.KJhL7vXmZ9F3uY5tKnM8wQxPzN4bC2K', 'Hoàng Văn Em', 'Staff', 'KIT003', 'em@adongfood.vn', '0901234571', true);

-- Note: All passwords are hashed version of 'password123'

-- ============================================================================
-- MASTER SUPPLIERS (7 suppliers)
-- ============================================================================
INSERT INTO public.master_suppliers (supplier_id, supplier_name, zalo_link, address, phone, email, active) VALUES
('SUP001', 'Công ty Thực phẩm Sạch Việt', 'https://zalo.me/sachviet', '45 Đường Bến Vân Đồn, Quận 4, TP.HCM', '0283567890', 'sachviet@gmail.com', true),
('SUP002', 'Nhà cung cấp Hải sản Tươi Sống', 'https://zalo.me/haisantuoisong', '78 Đường Đinh Tiên Hoàng, Quận Bình Thạnh, TP.HCM', '0287890123', 'haisantuoi@gmail.com', true),
('SUP003', 'Cửa hàng Rau Củ Đà Lạt', 'https://zalo.me/raucudalat', '123 Đường Lý Thường Kiệt, Quận 10, TP.HCM', '0289012345', 'raudalat@gmail.com', true),
('SUP004', 'Công ty Gia vị Việt Nam', 'https://zalo.me/giavivietnam', '567 Đường Nguyễn Tri Phương, Quận 5, TP.HCM', '0281234567', 'giavi@gmail.com', true),
('SUP005', 'Nhà phân phối Thịt Sạch An Toàn', 'https://zalo.me/thitsach', '234 Đường Võ Thị Sáu, Quận 3, TP.HCM', '0283456123', 'thitsach@gmail.com', true),
('SUP006', 'Cửa hàng Gạo Đồng Tháp', 'https://zalo.me/gaodongthap', '890 Đường Lê Hồng Phong, Quận 10, TP.HCM', '0287891234', 'gaodongthap@gmail.com', true),
('SUP007', 'Công ty Dầu Ăn Cao Cấp', 'https://zalo.me/dauan', '345 Đường Trần Hưng Đạo, Quận 1, TP.HCM', '0289012567', 'dauan@gmail.com', true);

-- ============================================================================
-- MASTER DISHES (15 dishes)
-- ============================================================================
INSERT INTO public.master_dishes (dish_id, dish_name, cooking_method, category, description, active) VALUES
-- Món thịt (5 dishes)
('DISH001', 'Thịt kho tàu', 'Kho', 'Món mặn', 'Thịt heo ba chỉ kho với trứng, nước dừa và nước mắm', true),
('DISH002', 'Bò lúc lắc', 'Xào', 'Món mặn', 'Thịt bò bắp cắt khối xào với hành tây và sốt tiêu đen', true),
('DISH003', 'Gà kho gừng', 'Kho', 'Món mặn', 'Thịt gà kho với gừng và nước mắm đường', true),
('DISH004', 'Sườn xào chua ngọt', 'Xào', 'Món mặn', 'Sườn heo xào với dứa và cà chua sốt chua ngọt', true),
('DISH005', 'Vịt nấu chao', 'Nấu', 'Món mặn', 'Thịt vịt nấu với chao và gừng', true),

-- Món hải sản (4 dishes)
('DISH006', 'Tôm rim', 'Rim', 'Món mặn', 'Tôm sú rim với nước mắm và tiêu', true),
('DISH007', 'Cá lóc kho tộ', 'Kho', 'Món mặn', 'Cá lóc kho tộ với nước dừa và ớt', true),
('DISH008', 'Mực xào chua ngọt', 'Xào', 'Món mặn', 'Mực ống xào với dứa và ớt chuông', true),
('DISH009', 'Nghêu hấp xả', 'Hấp', 'Món mặn', 'Nghêu hấp với sả và ớt', true),

-- Món canh và rau (3 dishes)
('DISH010', 'Canh chua cá', 'Nấu', 'Canh', 'Canh chua với cá basa, cà chua và rau thơm', true),
('DISH011', 'Rau muống xào tỏi', 'Xào', 'Rau', 'Rau muống xào với tỏi và nước mắm', true),
('DISH012', 'Cải bắp xào', 'Xào', 'Rau', 'Cải bắp xào với tỏi và nước mắm', true),

-- Món bún phở (3 dishes)
('DISH013', 'Bún bò Huế', 'Nấu', 'Bún phở', 'Bún bò với nước dùng cay và thịt bò', true),
('DISH014', 'Phở bò', 'Nấu', 'Bún phở', 'Phở với nước dùng xương bò và thịt bò', true),
('DISH015', 'Bún chả', 'Nướng', 'Bún phở', 'Bún với chả nướng và nước mắm chua ngọt', true);

-- ============================================================================
-- DISH RECIPE STANDARDS (Recipes for all dishes)
-- ============================================================================

-- DISH001: Thịt kho tàu
INSERT INTO public.dish_recipe_standards (dish_id, ingredient_id, unit, quantity_per_serving, notes, cost, updated_by_user_id) VALUES
('DISH001', 'ING001', 'kg', 0.2000, 'Thịt ba chỉ cắt miếng vừa', 40000.00, 'USR002'),
('DISH001', 'ING024', 'quả', 2.0000, 'Trứng luộc chín', 8000.00, 'USR002'),
('DISH001', 'ING017', 'lít', 0.0500, 'Nước mắm ngon', 2500.00, 'USR002'),
('DISH001', 'ING020', 'kg', 0.0050, 'Muối hạt', 50.00, 'USR002'),
('DISH001', 'ING012', 'kg', 0.0500, 'Hành tây thái múi cau', 1500.00, 'USR002');

-- DISH002: Bò lúc lắc
INSERT INTO public.dish_recipe_standards (dish_id, ingredient_id, unit, quantity_per_serving, notes, cost, updated_by_user_id) VALUES
('DISH002', 'ING002', 'kg', 0.1800, 'Thịt bò bắp cắt khối', 72000.00, 'USR002'),
('DISH002', 'ING012', 'kg', 0.0800, 'Hành tây thái miếng', 2400.00, 'USR002'),
('DISH002', 'ING011', 'kg', 0.0500, 'Cà chua bi', 1500.00, 'USR002'),
('DISH002', 'ING019', 'kg', 0.0020, 'Hạt tiêu đen', 200.00, 'USR002'),
('DISH002', 'ING018', 'lít', 0.0200, 'Dầu ăn', 600.00, 'USR002');

-- DISH003: Gà kho gừng
INSERT INTO public.dish_recipe_standards (dish_id, ingredient_id, unit, quantity_per_serving, notes, cost, updated_by_user_id) VALUES
('DISH003', 'ING003', 'kg', 0.2500, 'Thịt gà ta cắt miếng', 50000.00, 'USR002'),
('DISH003', 'ING017', 'lít', 0.0300, 'Nước mắm', 1500.00, 'USR002'),
('DISH003', 'ING020', 'kg', 0.0030, 'Muối', 30.00, 'USR002'),
('DISH003', 'ING018', 'lít', 0.0150, 'Dầu ăn', 450.00, 'USR002');

-- DISH004: Sườn xào chua ngọt
INSERT INTO public.dish_recipe_standards (dish_id, ingredient_id, unit, quantity_per_serving, notes, cost, updated_by_user_id) VALUES
('DISH004', 'ING004', 'kg', 0.2500, 'Sườn heo cắt miếng', 50000.00, 'USR002'),
('DISH004', 'ING011', 'kg', 0.1000, 'Cà chua thái múi cau', 3000.00, 'USR002'),
('DISH004', 'ING012', 'kg', 0.0600, 'Hành tây thái miếng', 1800.00, 'USR002'),
('DISH004', 'ING016', 'kg', 0.0100, 'Ớt chuông', 300.00, 'USR002'),
('DISH004', 'ING018', 'lít', 0.0200, 'Dầu ăn', 600.00, 'USR002');

-- DISH005: Vịt nấu chao
INSERT INTO public.dish_recipe_standards (dish_id, ingredient_id, unit, quantity_per_serving, notes, cost, updated_by_user_id) VALUES
('DISH005', 'ING005', 'kg', 0.3000, 'Thịt vịt cắt miếng', 60000.00, 'USR003'),
('DISH005', 'ING017', 'lít', 0.0300, 'Nước mắm', 1500.00, 'USR003'),
('DISH005', 'ING018', 'lít', 0.0200, 'Dầu ăn', 600.00, 'USR003');

-- DISH006: Tôm rim
INSERT INTO public.dish_recipe_standards (dish_id, ingredient_id, unit, quantity_per_serving, notes, cost, updated_by_user_id) VALUES
('DISH006', 'ING006', 'kg', 0.2000, 'Tôm sú tươi', 100000.00, 'USR002'),
('DISH006', 'ING017', 'lít', 0.0400, 'Nước mắm', 2000.00, 'USR002'),
('DISH006', 'ING019', 'kg', 0.0030, 'Hạt tiêu', 300.00, 'USR002'),
('DISH006', 'ING018', 'lít', 0.0150, 'Dầu ăn', 450.00, 'USR002');

-- DISH007: Cá lóc kho tộ
INSERT INTO public.dish_recipe_standards (dish_id, ingredient_id, unit, quantity_per_serving, notes, cost, updated_by_user_id) VALUES
('DISH007', 'ING007', 'kg', 0.3000, 'Cá lóc cắt khúc', 90000.00, 'USR002'),
('DISH007', 'ING017', 'lít', 0.0500, 'Nước mắm', 2500.00, 'USR002'),
('DISH007', 'ING016', 'kg', 0.0200, 'Ớt cắt khúc', 600.00, 'USR002'),
('DISH007', 'ING018', 'lít', 0.0200, 'Dầu ăn', 600.00, 'USR002');

-- DISH008: Mực xào chua ngọt
INSERT INTO public.dish_recipe_standards (dish_id, ingredient_id, unit, quantity_per_serving, notes, cost, updated_by_user_id) VALUES
('DISH008', 'ING008', 'kg', 0.2500, 'Mực ống tươi', 75000.00, 'USR002'),
('DISH008', 'ING012', 'kg', 0.0700, 'Hành tây', 2100.00, 'USR002'),
('DISH008', 'ING016', 'kg', 0.0150, 'Ớt chuông', 450.00, 'USR002'),
('DISH008', 'ING018', 'lít', 0.0200, 'Dầu ăn', 600.00, 'USR002');

-- DISH009: Nghêu hấp xả
INSERT INTO public.dish_recipe_standards (dish_id, ingredient_id, unit, quantity_per_serving, notes, cost, updated_by_user_id) VALUES
('DISH009', 'ING009', 'kg', 0.3000, 'Nghêu tươi sống', 60000.00, 'USR002'),
('DISH009', 'ING016', 'kg', 0.0150, 'Ớt cắt lát', 450.00, 'USR002'),
('DISH009', 'ING017', 'lít', 0.0200, 'Nước mắm', 1000.00, 'USR002');

-- DISH010: Canh chua cá
INSERT INTO public.dish_recipe_standards (dish_id, ingredient_id, unit, quantity_per_serving, notes, cost, updated_by_user_id) VALUES
('DISH010', 'ING010', 'kg', 0.2000, 'Cá basa cắt khúc', 40000.00, 'USR003'),
('DISH010', 'ING011', 'kg', 0.1000, 'Cà chua thái múi', 3000.00, 'USR003'),
('DISH010', 'ING017', 'lít', 0.0300, 'Nước mắm', 1500.00, 'USR003');

-- DISH011: Rau muống xào tỏi
INSERT INTO public.dish_recipe_standards (dish_id, ingredient_id, unit, quantity_per_serving, notes, cost, updated_by_user_id) VALUES
('DISH011', 'ING014', 'kg', 0.2000, 'Rau muống tươi', 4000.00, 'USR002'),
('DISH011', 'ING017', 'lít', 0.0150, 'Nước mắm', 750.00, 'USR002'),
('DISH011', 'ING018', 'lít', 0.0150, 'Dầu ăn', 450.00, 'USR002');

-- DISH012: Cải bắp xào
INSERT INTO public.dish_recipe_standards (dish_id, ingredient_id, unit, quantity_per_serving, notes, cost, updated_by_user_id) VALUES
('DISH012', 'ING015', 'kg', 0.2000, 'Cải bắp thái sợi', 6000.00, 'USR002'),
('DISH012', 'ING017', 'lít', 0.0150, 'Nước mắm', 750.00, 'USR002'),
('DISH012', 'ING018', 'lít', 0.0150, 'Dầu ăn', 450.00, 'USR002');

-- DISH013: Bún bò Huế
INSERT INTO public.dish_recipe_standards (dish_id, ingredient_id, unit, quantity_per_serving, notes, cost, updated_by_user_id) VALUES
('DISH013', 'ING022', 'kg', 0.2500, 'Bún tươi', 10000.00, 'USR003'),
('DISH013', 'ING002', 'kg', 0.1500, 'Thịt bò bắp', 60000.00, 'USR003'),
('DISH013', 'ING017', 'lít', 0.0300, 'Nước mắm', 1500.00, 'USR003'),
('DISH013', 'ING016', 'kg', 0.0200, 'Ớt', 600.00, 'USR003');

-- DISH014: Phở bò
INSERT INTO public.dish_recipe_standards (dish_id, ingredient_id, unit, quantity_per_serving, notes, cost, updated_by_user_id) VALUES
('DISH014', 'ING023', 'kg', 0.2500, 'Bánh phở tươi', 10000.00, 'USR003'),
('DISH014', 'ING002', 'kg', 0.1500, 'Thịt bò bắp', 60000.00, 'USR003'),
('DISH014', 'ING012', 'kg', 0.0300, 'Hành tây', 900.00, 'USR003'),
('DISH014', 'ING017', 'lít', 0.0300, 'Nước mắm', 1500.00, 'USR003');

-- DISH015: Bún chả
INSERT INTO public.dish_recipe_standards (dish_id, ingredient_id, unit, quantity_per_serving, notes, cost, updated_by_user_id) VALUES
('DISH015', 'ING022', 'kg', 0.2500, 'Bún tươi', 10000.00, 'USR003'),
('DISH015', 'ING001', 'kg', 0.1500, 'Thịt heo ba chỉ', 30000.00, 'USR003'),
('DISH015', 'ING017', 'lít', 0.0400, 'Nước mắm', 2000.00, 'USR003'),
('DISH015', 'ING011', 'kg', 0.0500, 'Cà chua', 1500.00, 'USR003');

-- ============================================================================
-- SUPPLIER PRICE LIST (Multiple suppliers for each ingredient)
-- ============================================================================
-- Note: product_id is auto-increment, so we don't specify it

-- Thịt heo ba chỉ (ING001) - 3 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP001', 'ING001', 'kg', 150000, true),
('SUP005', 'ING001', 'kg', 145000, true),
('SUP001', 'ING001', 'kg', 148000, true);

-- Thịt bò bắp (ING002) - 3 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP001', 'ING002', 'kg', 350000, true),
('SUP005', 'ING002', 'kg', 345000, true),
('SUP001', 'ING002', 'kg', 355000, true);

-- Thịt gà ta (ING003) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP001', 'ING003', 'kg', 180000, true),
('SUP005', 'ING003', 'kg', 175000, true);

-- Sườn heo (ING004) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP005', 'ING004', 'kg', 165000, true),
('SUP001', 'ING004', 'kg', 170000, true);

-- Thịt vịt (ING005) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP005', 'ING005', 'kg', 190000, true),
('SUP001', 'ING005', 'kg', 195000, true);

-- Tôm sú (ING006) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP002', 'ING006', 'kg', 450000, true),
('SUP002', 'ING006', 'kg', 460000, true);

-- Cá lóc (ING007) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP002', 'ING007', 'kg', 280000, true),
('SUP002', 'ING007', 'kg', 285000, true);

-- Mực ống (ING008) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP002', 'ING008', 'kg', 250000, true),
('SUP002', 'ING008', 'kg', 255000, true);

-- Nghêu (ING009) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP002', 'ING009', 'kg', 180000, true),
('SUP002', 'ING009', 'kg', 175000, true);

-- Cá basa (ING010) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP002', 'ING010', 'kg', 120000, true),
('SUP002', 'ING010', 'kg', 115000, true);

-- Cà chua (ING011) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP003', 'ING011', 'kg', 25000, true),
('SUP003', 'ING011', 'kg', 23000, true);

-- Hành tây (ING012) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP003', 'ING012', 'kg', 28000, true),
('SUP003', 'ING012', 'kg', 26000, true);

-- Khoai tây (ING013) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP003', 'ING013', 'kg', 22000, true),
('SUP003', 'ING013', 'kg', 20000, true);

-- Rau muống (ING014) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP003', 'ING014', 'kg', 18000, true),
('SUP003', 'ING014', 'kg', 17000, true);

-- Cải bắp (ING015) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP003', 'ING015', 'kg', 28000, true),
('SUP003', 'ING015', 'kg', 26000, true);

-- Ớt (ING016) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP003', 'ING016', 'kg', 35000, true),
('SUP004', 'ING016', 'kg', 33000, true);

-- Nước mắm (ING017) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP004', 'ING017', 'lít', 45000, true),
('SUP004', 'ING017', 'lít', 43000, true);

-- Dầu ăn (ING018) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP007', 'ING018', 'lít', 35000, true),
('SUP004', 'ING018', 'lít', 33000, true);

-- Hạt tiêu (ING019) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP004', 'ING019', 'kg', 180000, true),
('SUP004', 'ING019', 'kg', 175000, true);

-- Muối (ING020) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP004', 'ING020', 'kg', 8000, true),
('SUP004', 'ING020', 'kg', 7500, true);

-- Gạo tẻ (ING021) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP006', 'ING021', 'kg', 22000, true),
('SUP006', 'ING021', 'kg', 21000, true);

-- Bún tươi (ING022) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP006', 'ING022', 'kg', 18000, true),
('SUP006', 'ING022', 'kg', 17000, true);

-- Bánh phở (ING023) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP006', 'ING023', 'kg', 18000, true),
('SUP006', 'ING023', 'kg', 17500, true);

-- Trứng gà (ING024) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP001', 'ING024', 'quả', 3500, true),
('SUP005', 'ING024', 'quả', 3300, true);

-- Sữa tươi (ING025) - 2 suppliers
INSERT INTO public.supplier_price_list (supplier_id, ingredient_id, unit, unit_price, active) VALUES
('SUP001', 'ING025', 'lít', 35000, true),
('SUP005', 'ING025', 'lít', 33000, true);

-- ============================================================================
-- KITCHEN FAVORITE SUPPLIERS
-- ============================================================================
-- Note: kitchen_favorite_suppliers table does NOT have ingredient_id column
INSERT INTO public.kitchen_favorite_suppliers (kitchen_id, supplier_id, notes, created_by_user_id) VALUES
-- KIT001 favorites
('KIT001', 'SUP001', 'Thịt tươi chất lượng tốt', 'USR002'),
('KIT001', 'SUP002', 'Tôm luôn tươi sống', 'USR002'),
('KIT001', 'SUP003', 'Rau sạch Đà Lạt', 'USR002'),
('KIT001', 'SUP004', 'Nước mắm ngon', 'USR002'),

-- KIT002 favorites
('KIT002', 'SUP005', 'Giá tốt, giao hàng đúng giờ', 'USR003'),
('KIT002', 'SUP002', 'Cá tươi mỗi ngày', 'USR003'),
('KIT002', 'SUP006', 'Gạo ngon', 'USR003'),

-- KIT003 favorites
('KIT003', 'SUP001', 'Gà ta chất lượng', 'USR005'),
('KIT003', 'SUP003', 'Rau củ tươi', 'USR005'),
('KIT003', 'SUP007', 'Dầu ăn chất lượng cao', 'USR005');

-- ============================================================================
-- SAMPLE ORDERS
-- ============================================================================

-- Order 1: KIT001 - Pending
INSERT INTO public.orders (order_id, kitchen_id, order_date, note, status, created_by_user_id) VALUES
('ORD001', 'KIT001', '2025-11-15', 'Đơn hàng cho thực đơn tuần tới', 'Pending', 'USR002');

-- Order 1 Details - 3 dishes
INSERT INTO public.order_details (order_id, dish_id, portions, note) VALUES
('ORD001', 'DISH001', 50, 'Thịt kho cho 50 suất'),
('ORD001', 'DISH006', 30, 'Tôm rim cho 30 suất'),
('ORD001', 'DISH011', 60, 'Rau muống cho 60 suất');

-- Order 1 Ingredients (auto-calculated from recipe standards)
INSERT INTO public.order_ingredients (order_detail_id, ingredient_id, quantity, unit, standard_per_portion) VALUES
-- DISH001 (50 portions) - order_detail_id will be 1
(1, 'ING001', 10.0000, 'kg', 0.2000),
(1, 'ING024', 100.0000, 'quả', 2.0000),
(1, 'ING017', 2.5000, 'lít', 0.0500),
(1, 'ING020', 0.2500, 'kg', 0.0050),
(1, 'ING012', 2.5000, 'kg', 0.0500),

-- DISH006 (30 portions) - order_detail_id will be 2
(2, 'ING006', 6.0000, 'kg', 0.2000),
(2, 'ING017', 1.2000, 'lít', 0.0400),
(2, 'ING019', 0.0900, 'kg', 0.0030),
(2, 'ING018', 0.4500, 'lít', 0.0150),

-- DISH011 (60 portions) - order_detail_id will be 3
(3, 'ING014', 12.0000, 'kg', 0.2000),
(3, 'ING017', 0.9000, 'lít', 0.0150),
(3, 'ING018', 0.9000, 'lít', 0.0150);

-- Order 1 Supplementary Foods (additional ingredients not from recipes)
INSERT INTO public.order_supplementary_foods (order_id, ingredient_id, quantity, unit, standard_per_portion, portions) VALUES
('ORD001', 'ING021', 30.0000, 'kg', 0.3000, 100),  -- Gạo tẻ cho 100 suất
('ORD001', 'ING016', 2.0000, 'kg', NULL, NULL);     -- Ớt phụ trội

-- Order 2: KIT002 - Approved
INSERT INTO public.orders (order_id, kitchen_id, order_date, note, status, created_by_user_id) VALUES
('ORD002', 'KIT002', '2025-11-10', 'Đơn hàng đã duyệt', 'Approved', 'USR003');

-- Order 2 Details - 2 dishes
INSERT INTO public.order_details (order_id, dish_id, portions, note) VALUES
('ORD002', 'DISH002', 40, 'Bò lúc lắc cho 40 suất'),
('ORD002', 'DISH013', 50, 'Bún bò Huế cho 50 suất');

-- Order 2 Ingredients - order_detail_id will be 4, 5
INSERT INTO public.order_ingredients (order_detail_id, ingredient_id, quantity, unit, standard_per_portion) VALUES
-- DISH002 (40 portions)
(4, 'ING002', 7.2000, 'kg', 0.1800),
(4, 'ING012', 3.2000, 'kg', 0.0800),
(4, 'ING011', 2.0000, 'kg', 0.0500),
(4, 'ING019', 0.0800, 'kg', 0.0020),
(4, 'ING018', 0.8000, 'lít', 0.0200),

-- DISH013 (50 portions)
(5, 'ING022', 12.5000, 'kg', 0.2500),
(5, 'ING002', 7.5000, 'kg', 0.1500),
(5, 'ING017', 1.5000, 'lít', 0.0300),
(5, 'ING016', 1.0000, 'kg', 0.0200);

-- Order 2 Supplier Selections (linking to supplier price list)
-- Note: selected_product_id must reference actual product_id from supplier_price_list
-- Since product_id is auto-increment, we use subquery to find matching product
INSERT INTO public.order_ingredient_suppliers (order_id, ingredient_id, selected_supplier_id, selected_product_id, quantity, unit, unit_price, total_cost, selected_by_user_id, notes) VALUES
('ORD002', 'ING002', 'SUP005', (SELECT product_id FROM supplier_price_list WHERE supplier_id = 'SUP005' AND ingredient_id = 'ING002' AND unit_price = 345000 LIMIT 1), 14.7000, 'kg', 345000, 5071500, 'USR003', 'Nhà cung cấp ưa thích'),
('ORD002', 'ING012', 'SUP003', (SELECT product_id FROM supplier_price_list WHERE supplier_id = 'SUP003' AND ingredient_id = 'ING012' AND unit_price = 26000 LIMIT 1), 3.2000, 'kg', 26000, 83200, 'USR003', NULL),
('ORD002', 'ING022', 'SUP006', (SELECT product_id FROM supplier_price_list WHERE supplier_id = 'SUP006' AND ingredient_id = 'ING022' AND unit_price = 17000 LIMIT 1), 12.5000, 'kg', 17000, 212500, 'USR003', NULL);

-- Order 3: KIT003 - Completed
INSERT INTO public.orders (order_id, kitchen_id, order_date, note, status, created_by_user_id) VALUES
('ORD003', 'KIT003', '2025-11-05', 'Đơn hàng đã hoàn thành', 'Completed', 'USR005');

-- Order 3 Details
INSERT INTO public.order_details (order_id, dish_id, portions, note) VALUES
('ORD003', 'DISH003', 35, 'Gà kho gừng'),
('ORD003', 'DISH012', 45, 'Cải bắp xào');

-- Order 3 Ingredients - order_detail_id will be 6, 7
INSERT INTO public.order_ingredients (order_detail_id, ingredient_id, quantity, unit, standard_per_portion) VALUES
-- DISH003 (35 portions)
(6, 'ING003', 8.7500, 'kg', 0.2500),
(6, 'ING017', 1.0500, 'lít', 0.0300),
(6, 'ING020', 0.1050, 'kg', 0.0030),
(6, 'ING018', 0.5250, 'lít', 0.0150),

-- DISH012 (45 portions)
(7, 'ING015', 9.0000, 'kg', 0.2000),
(7, 'ING017', 0.6750, 'lít', 0.0150),
(7, 'ING018', 0.6750, 'lít', 0.0150);

COMMIT;

-- ============================================================================
-- VERIFICATION QUERIES (Run these to verify the data)
-- ============================================================================

-- Count all records
-- SELECT 'ingredient_types' as table_name, COUNT(*) as count FROM public.ingredient_types
-- UNION ALL
-- SELECT 'master_ingredients', COUNT(*) FROM public.master_ingredients
-- UNION ALL
-- SELECT 'master_dishes', COUNT(*) FROM public.master_dishes
-- UNION ALL
-- SELECT 'master_kitchens', COUNT(*) FROM public.master_kitchens
-- UNION ALL
-- SELECT 'master_suppliers', COUNT(*) FROM public.master_suppliers
-- UNION ALL
-- SELECT 'master_users', COUNT(*) FROM public.master_users
-- UNION ALL
-- SELECT 'dish_recipe_standards', COUNT(*) FROM public.dish_recipe_standards
-- UNION ALL
-- SELECT 'supplier_price_list', COUNT(*) FROM public.supplier_price_list
-- UNION ALL
-- SELECT 'kitchen_favorite_suppliers', COUNT(*) FROM public.kitchen_favorite_suppliers
-- UNION ALL
-- SELECT 'orders', COUNT(*) FROM public.orders
-- UNION ALL
-- SELECT 'order_details', COUNT(*) FROM public.order_details
-- UNION ALL
-- SELECT 'order_ingredients', COUNT(*) FROM public.order_ingredients
-- UNION ALL
-- SELECT 'order_supplementary_foods', COUNT(*) FROM public.order_supplementary_foods
-- UNION ALL
-- SELECT 'order_ingredient_suppliers', COUNT(*) FROM public.order_ingredient_suppliers;

-- View sample dishes with their recipes
-- SELECT 
--     md.dish_name,
--     mi.ingredient_name,
--     drs.quantity_per_serving,
--     drs.unit,
--     drs.cost
-- FROM dish_recipe_standards drs
-- JOIN master_dishes md ON drs.dish_id = md.dish_id
-- JOIN master_ingredients mi ON drs.ingredient_id = mi.ingredient_id
-- ORDER BY md.dish_name, mi.ingredient_name;

-- View orders with details
-- SELECT 
--     o.order_id,
--     mk.kitchen_name,
--     o.order_date,
--     o.status,
--     md.dish_name,
--     od.portions
-- FROM orders o
-- JOIN master_kitchens mk ON o.kitchen_id = mk.kitchen_id
-- JOIN order_details od ON o.order_id = od.order_id
-- JOIN master_dishes md ON od.dish_id = md.dish_id
-- ORDER BY o.order_date DESC, o.order_id, md.dish_name;