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
