# Inventory Management API Documentation

## Overview

The Inventory Management API provides endpoints for managing inventory stocks, imports, and exports across multiple kitchens. All endpoints require authentication via Bearer token.

**Base URL:** `/api/inventory`

**Authentication:** All endpoints require a valid JWT token in the `Authorization` header:
```
Authorization: Bearer <your-access-token>
```

---

## Table of Contents

1. [Stock Management](#stock-management)
2. [Import Management](#import-management)
3. [Export Management](#export-management)
4. [Data Models](#data-models)
5. [Error Handling](#error-handling)

---

## Stock Management

### Get All Stocks

Retrieve paginated list of inventory stocks with optional filters.

**Endpoint:** `GET /api/inventory/stocks`

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `per_page` (optional, default: 50) - Items per page
- `search` (optional) - Search by ingredient name (case-insensitive)
- `sort_by` (optional) - Sort field (e.g., "last_updated", "quantity", "ingredient_id")
- `sort_dir` (optional) - Sort direction ("asc" or "desc")
- `kitchen_id` (optional) - Filter by kitchen ID
- `low_stock` (optional) - Set to "true" to filter items below minimum stock level

**Example Request:**
```http
GET /api/inventory/stocks?kitchen_id=K001&page=1&per_page=50&low_stock=true&sort_by=last_updated&sort_dir=desc
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "stockId": 1,
      "kitchenId": "K001",
      "ingredientId": "NL001",
      "quantity": 100.5,
      "unit": "kg",
      "minStockLevel": 50.0,
      "maxStockLevel": 200.0,
      "lastUpdated": "2024-01-15T10:30:00Z",
      "createdDate": "2024-01-01T00:00:00Z",
      "modifiedDate": "2024-01-15T10:30:00Z",
      "kitchen": {
        "kitchenId": "K001",
        "kitchenName": "Main Kitchen"
      },
      "ingredient": {
        "ingredientId": "NL001",
        "ingredientName": "Flour"
      }
    }
  ],
  "meta": {
    "current_page": 1,
    "last_page": 3,
    "from": 1,
    "to": 50,
    "per_page": 50,
    "total": 150
  }
}
```

---

### Get Stock by ID

Retrieve a specific stock record by its ID.

**Endpoint:** `GET /api/inventory/stocks/:id`

**Path Parameters:**
- `id` (required) - Stock ID

**Example Request:**
```http
GET /api/inventory/stocks/1
```

**Response (200 OK):**
```json
{
  "data": {
    "stockId": 1,
    "kitchenId": "K001",
    "ingredientId": "NL001",
    "quantity": 100.5,
    "unit": "kg",
    "minStockLevel": 50.0,
    "maxStockLevel": 200.0,
    "lastUpdated": "2024-01-15T10:30:00Z",
    "kitchen": { ... },
    "ingredient": { ... }
  }
}
```

**Error Response (404 Not Found):**
```json
{
  "error": "Không tìm thấy tồn kho"
}
```

---

### Get Stock by Kitchen and Ingredient

Retrieve stock for a specific kitchen and ingredient combination.

**Endpoint:** `GET /api/inventory/stocks/query`

**Query Parameters:**
- `kitchen_id` (required) - Kitchen ID
- `ingredient_id` (required) - Ingredient ID

**Example Request:**
```http
GET /api/inventory/stocks/query?kitchen_id=K001&ingredient_id=NL001
```

**Response (200 OK):**
```json
{
  "data": {
    "stockId": 1,
    "kitchenId": "K001",
    "ingredientId": "NL001",
    "quantity": 100.5,
    "unit": "kg",
    ...
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "Cần có kitchen_id và ingredient_id"
}
```

---

### Update Stock Levels

Update minimum and/or maximum stock levels for a stock item.

**Endpoint:** `PUT /api/inventory/stocks/:id/levels`

**Path Parameters:**
- `id` (required) - Stock ID

**Request Body:**
```json
{
  "minStockLevel": 50.0,
  "maxStockLevel": 200.0
}
```

**Response (200 OK):**
```json
{
  "message": "Cập nhật mức tồn thành công",
  "data": {
    "stockId": 1,
    "minStockLevel": 50.0,
    "maxStockLevel": 200.0,
    ...
  }
}
```

---

### Get Low Stock Alerts

Retrieve all items that are below their minimum stock level.

**Endpoint:** `GET /api/inventory/stocks/alerts/low`

**Query Parameters:**
- `kitchen_id` (optional) - Filter by kitchen ID

**Example Request:**
```http
GET /api/inventory/stocks/alerts/low?kitchen_id=K001
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "stockId": 1,
      "kitchenId": "K001",
      "ingredientId": "NL001",
      "quantity": 30.0,
      "minStockLevel": 50.0,
      ...
    }
  ],
  "count": 5
}
```

---

### Get Stock Transactions

Retrieve transaction history for a specific stock item.

**Endpoint:** `GET /api/inventory/stocks/transactions`

**Query Parameters:**
- `kitchen_id` (required) - Kitchen ID
- `ingredient_id` (required) - Ingredient ID
- `page` (optional, default: 1) - Page number
- `per_page` (optional, default: 50) - Items per page
- `search` (optional) - Search term
- `sort_by` (optional) - Sort field (e.g., "transaction_date", "transaction_type", "quantity")
- `sort_dir` (optional) - Sort direction ("asc" or "desc")
- `transaction_type` (optional) - Filter by transaction type (IMPORT, EXPORT, TRANSFER_IN, etc.)
- `from_date` (optional) - Filter from date (YYYY-MM-DD)
- `to_date` (optional) - Filter to date (YYYY-MM-DD)

**Example Request:**
```http
GET /api/inventory/stocks/transactions?kitchen_id=K001&ingredient_id=NL001&transaction_type=IMPORT&page=1&per_page=50&sort_by=transaction_date&sort_dir=desc
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "transactionId": 1,
      "kitchenId": "K001",
      "ingredientId": "NL001",
      "transactionType": "IMPORT",
      "transactionDate": "2024-01-15T10:30:00Z",
      "quantity": 50.0,
      "unit": "kg",
      "quantityBefore": 50.0,
      "quantityAfter": 100.0,
      "referenceType": "IMPORT",
      "referenceId": "IM20240115-12345",
      "createdBy": { ... }
    }
  ],
  "meta": {
    "current_page": 1,
    "last_page": 2,
    "from": 1,
    "to": 50,
    "per_page": 50,
    "total": 100
  }
}
```

---

### Get Stock Summary

Retrieve summary statistics for a kitchen's inventory.

**Endpoint:** `GET /api/inventory/stocks/summary`

**Query Parameters:**
- `kitchen_id` (required) - Kitchen ID

**Example Request:**
```http
GET /api/inventory/stocks/summary?kitchen_id=K001
```

**Response (200 OK):**
```json
{
  "data": {
    "totalItems": 150,
    "lowStockItems": 5,
    "outOfStockItems": 2,
    "totalValue": 50000.0
  }
}
```

---

### Get Stock Valuation

Retrieve stock valuation based on latest supplier prices.

**Endpoint:** `GET /api/inventory/stocks/valuation`

**Query Parameters:**
- `kitchen_id` (required) - Kitchen ID

**Example Request:**
```http
GET /api/inventory/stocks/valuation?kitchen_id=K001
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "ingredientId": "NL001",
      "ingredientName": "Flour",
      "quantity": 100.5,
      "unit": "kg",
      "averagePrice": 50.0,
      "totalValue": 5025.0
    }
  ],
  "totalValue": 50000.0,
  "count": 150
}
```

---

## Import Management

### Get All Imports

Retrieve paginated list of inventory imports with optional filters.

**Endpoint:** `GET /api/inventory/imports`

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `per_page` (optional, default: 20) - Items per page
- `search` (optional) - Search term
- `sort_by` (optional) - Sort field (e.g., "import_date", "created_date", "status", "total_amount")
- `sort_dir` (optional) - Sort direction ("asc" or "desc")
- `kitchen_id` (optional) - Filter by kitchen ID
- `status` (optional) - Filter by status (draft, approved)
- `from_date` (optional) - Filter from date (YYYY-MM-DD)
- `to_date` (optional) - Filter to date (YYYY-MM-DD)

**Example Request:**
```http
GET /api/inventory/imports?kitchen_id=K001&status=draft&page=1&per_page=20&sort_by=import_date&sort_dir=desc
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "importId": "IM20240115-12345",
      "kitchenId": "K001",
      "importDate": "2024-01-15T00:00:00Z",
      "orderId": "ORD001",
      "supplierId": "SUP001",
      "totalAmount": 5000.0,
      "status": "draft",
      "notes": "Initial import",
      "createdBy": { ... },
      "kitchen": { ... },
      "supplier": { ... }
    }
  ],
  "meta": {
    "current_page": 1,
    "last_page": 3,
    "from": 1,
    "to": 20,
    "per_page": 20,
    "total": 50
  }
}
```

---

### Get Import by ID

Retrieve a specific import record with all details.

**Endpoint:** `GET /api/inventory/imports/:id`

**Path Parameters:**
- `id` (required) - Import ID (e.g., "IM20240115-12345")

**Example Request:**
```http
GET /api/inventory/imports/IM20240115-12345
```

**Response (200 OK):**
```json
{
  "data": {
    "importId": "IM20240115-12345",
    "kitchenId": "K001",
    "importDate": "2024-01-15T00:00:00Z",
    "orderId": "ORD001",
    "supplierId": "SUP001",
    "totalAmount": 5000.0,
    "status": "draft",
    "notes": "Initial import",
    "createdDate": "2024-01-15T10:00:00Z",
    "kitchen": { ... },
    "supplier": { ... },
    "order": { ... },
    "importDetails": [
      {
        "importDetailId": 1,
        "importId": "IM20240115-12345",
        "ingredientId": "NL001",
        "quantity": 100.0,
        "unit": "kg",
        "unitPrice": 50.0,
        "totalPrice": 5000.0,
        "expiryDate": "2025-01-15T00:00:00Z",
        "batchNumber": "BATCH001",
        "notes": "High quality flour",
        "ingredient": {
          "ingredientId": "NL001",
          "ingredientName": "Flour"
        }
      }
    ]
  }
}
```

---

### Create Import

Create a new inventory import record.

**Endpoint:** `POST /api/inventory/imports`

**Request Body:**
```json
{
  "kitchenId": "K001",
  "importDate": "2024-01-15",
  "orderId": "ORD001",
  "supplierId": "SUP001",
  "status": "draft",
  "notes": "Initial import",
  "importDetails": [
    {
      "ingredientId": "NL001",
      "quantity": 100.0,
      "unit": "kg",
      "unitPrice": 50.0,
      "expiryDate": "2025-01-15",
      "batchNumber": "BATCH001",
      "notes": "High quality flour"
    }
  ]
}
```

**Field Requirements:**
- `kitchenId` (required) - Kitchen ID
- `importDate` (required) - Date in YYYY-MM-DD format
- `orderId` (optional) - Related order ID
- `supplierId` (optional) - Supplier ID
- `status` (optional) - Status, defaults to "draft"
- `notes` (optional) - Additional notes
- `importDetails` (required, min: 1) - Array of import detail items
  - `ingredientId` (required) - Ingredient ID
  - `quantity` (required, > 0) - Quantity
  - `unit` (required) - Unit of measurement
  - `unitPrice` (required, > 0) - Unit price
  - `expiryDate` (optional) - Expiry date in YYYY-MM-DD format
  - `batchNumber` (optional) - Batch number
  - `notes` (optional) - Detail notes

**Response (201 Created):**
```json
{
  "message": "Tạo phiếu nhập thành công",
  "data": {
    "importId": "IM20240115-12345",
    "kitchenId": "K001",
    "totalAmount": 5000.0,
    "status": "draft",
    "importDetails": [ ... ]
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "Định dạng ngày không hợp lệ"
}
```

---

### Update Import

Update an existing import record. Only draft imports can be updated.

**Endpoint:** `PUT /api/inventory/imports/:id`

**Path Parameters:**
- `id` (required) - Import ID

**Request Body:** Same as Create Import

**Response (200 OK):**
```json
{
  "message": "Cập nhật phiếu nhập thành công",
  "data": { ... }
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "Không thể sửa phiếu nhập đã duyệt"
}
```

---

### Approve Import

Approve an import and update inventory stocks. This action:
- Changes status to "approved"
- Updates inventory stock quantities
- Creates transaction logs
- Cannot be undone

**Endpoint:** `POST /api/inventory/imports/:id/approve`

**Path Parameters:**
- `id` (required) - Import ID

**Example Request:**
```http
POST /api/inventory/imports/IM20240115-12345/approve
```

**Response (200 OK):**
```json
{
  "message": "Duyệt phiếu nhập thành công",
  "data": {
    "importId": "IM20240115-12345",
    "status": "approved",
    "approvedBy": { ... },
    "approvedDate": "2024-01-15T10:30:00Z",
    ...
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "Phiếu nhập đã được duyệt"
}
```

---

### Delete Import

Delete a draft import. Approved imports cannot be deleted.

**Endpoint:** `DELETE /api/inventory/imports/:id`

**Path Parameters:**
- `id` (required) - Import ID

**Response (200 OK):**
```json
{
  "message": "Xóa phiếu nhập thành công"
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "Không thể xóa phiếu nhập đã duyệt"
}
```

---

## Export Management

### Get All Exports

Retrieve paginated list of inventory exports with optional filters.

**Endpoint:** `GET /api/inventory/exports`

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `per_page` (optional, default: 20) - Items per page
- `search` (optional) - Search term
- `sort_by` (optional) - Sort field (e.g., "export_date", "created_date", "status", "export_type", "total_amount")
- `sort_dir` (optional) - Sort direction ("asc" or "desc")
- `kitchen_id` (optional) - Filter by kitchen ID
- `export_type` (optional) - Filter by export type (production, transfer, disposal, return, sample)
- `status` (optional) - Filter by status (draft, approved)
- `from_date` (optional) - Filter from date (YYYY-MM-DD)
- `to_date` (optional) - Filter to date (YYYY-MM-DD)

**Example Request:**
```http
GET /api/inventory/exports?kitchen_id=K001&export_type=production&status=draft&page=1&per_page=20&sort_by=export_date&sort_dir=desc
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "exportId": "EX20240115-12345",
      "kitchenId": "K001",
      "exportDate": "2024-01-15T00:00:00Z",
      "exportType": "production",
      "destinationKitchenId": null,
      "orderId": "ORD001",
      "totalAmount": 3000.0,
      "status": "draft",
      "notes": "For production",
      "createdBy": { ... },
      "kitchen": { ... }
    }
  ],
  "meta": {
    "current_page": 1,
    "last_page": 2,
    "from": 1,
    "to": 20,
    "per_page": 20,
    "total": 30
  }
}
```

---

### Get Export by ID

Retrieve a specific export record with all details.

**Endpoint:** `GET /api/inventory/exports/:id`

**Path Parameters:**
- `id` (required) - Export ID (e.g., "EX20240115-12345")

**Example Request:**
```http
GET /api/inventory/exports/EX20240115-12345
```

**Response (200 OK):**
```json
{
  "data": {
    "exportId": "EX20240115-12345",
    "kitchenId": "K001",
    "exportDate": "2024-01-15T00:00:00Z",
    "exportType": "production",
    "destinationKitchenId": null,
    "orderId": "ORD001",
    "totalAmount": 3000.0,
    "status": "draft",
    "exportDetails": [
      {
        "exportDetailId": 1,
        "exportId": "EX20240115-12345",
        "ingredientId": "NL001",
        "quantity": 50.0,
        "unit": "kg",
        "unitCost": 50.0,
        "totalCost": 2500.0,
        "batchNumber": "BATCH001",
        "notes": "For production",
        "ingredient": { ... }
      }
    ]
  }
}
```

---

### Create Export

Create a new inventory export record.

**Endpoint:** `POST /api/inventory/exports`

**Request Body:**
```json
{
  "kitchenId": "K001",
  "exportDate": "2024-01-15",
  "exportType": "production",
  "destinationKitchenId": null,
  "orderId": "ORD001",
  "status": "draft",
  "notes": "For production",
  "exportDetails": [
    {
      "ingredientId": "NL001",
      "quantity": 50.0,
      "unit": "kg",
      "unitCost": 50.0,
      "batchNumber": "BATCH001",
      "notes": "For production"
    }
  ]
}
```

**Field Requirements:**
- `kitchenId` (required) - Kitchen ID
- `exportDate` (required) - Date in YYYY-MM-DD format
- `exportType` (required) - One of: "production", "transfer", "disposal", "return", "sample"
- `destinationKitchenId` (optional) - Required for "transfer" type
- `orderId` (optional) - Related order ID
- `status` (optional) - Status, defaults to "draft"
- `notes` (optional) - Additional notes
- `exportDetails` (required, min: 1) - Array of export detail items
  - `ingredientId` (required) - Ingredient ID
  - `quantity` (required, > 0) - Quantity
  - `unit` (required) - Unit of measurement
  - `unitCost` (optional) - Unit cost
  - `batchNumber` (optional) - Batch number
  - `notes` (optional) - Detail notes

**Export Types:**
- `production` - For production use
- `transfer` - Transfer to another kitchen (requires `destinationKitchenId`)
- `disposal` - Disposal/waste
- `return` - Return to supplier
- `sample` - Sample export

**Response (201 Created):**
```json
{
  "message": "Tạo phiếu xuất thành công",
  "data": {
    "exportId": "EX20240115-12345",
    "kitchenId": "K001",
    "exportType": "production",
    "totalAmount": 3000.0,
    "status": "draft",
    "exportDetails": [ ... ]
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "Loại xuất kho không hợp lệ"
}
```

---

### Update Export

Update an existing export record. Only draft exports can be updated.

**Endpoint:** `PUT /api/inventory/exports/:id`

**Path Parameters:**
- `id` (required) - Export ID

**Request Body:** Same as Create Export

**Response (200 OK):**
```json
{
  "message": "Cập nhật phiếu xuất thành công",
  "data": { ... }
}
```

---

### Approve Export

Approve an export and update inventory stocks. This action:
- Validates stock availability
- Changes status to "approved"
- Decreases inventory stock quantities
- Creates transaction logs
- For "transfer" type, increases destination kitchen stock
- Cannot be undone

**Endpoint:** `POST /api/inventory/exports/:id/approve`

**Path Parameters:**
- `id` (required) - Export ID

**Example Request:**
```http
POST /api/inventory/exports/EX20240115-12345/approve
```

**Response (200 OK):**
```json
{
  "message": "Duyệt phiếu xuất thành công",
  "data": {
    "exportId": "EX20240115-12345",
    "status": "approved",
    "approvedBy": { ... },
    "approvedDate": "2024-01-15T10:30:00Z",
    ...
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "Số lượng tồn kho không đủ",
  "ingredient_id": "NL001",
  "available": 30.0,
  "required": 50.0
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "Nguyên liệu không tồn tại trong kho",
  "ingredient_id": "NL001"
}
```

---

### Delete Export

Delete a draft export. Approved exports cannot be deleted.

**Endpoint:** `DELETE /api/inventory/exports/:id`

**Path Parameters:**
- `id` (required) - Export ID

**Response (200 OK):**
```json
{
  "message": "Xóa phiếu xuất thành công"
}
```

---

## Data Models

### InventoryStock

```typescript
interface InventoryStock {
  stockId: number;
  kitchenId: string;
  ingredientId: string;
  quantity: number;
  unit: string;
  minStockLevel?: number;
  maxStockLevel?: number;
  lastUpdated: string; // ISO 8601 datetime
  createdDate: string; // ISO 8601 datetime
  modifiedDate: string; // ISO 8601 datetime
  kitchen?: Kitchen;
  ingredient?: Ingredient;
}
```

### InventoryImport

```typescript
interface InventoryImport {
  importId: string; // Format: IMYYYYMMDD-XXXXX
  kitchenId: string;
  importDate: string; // ISO 8601 date
  orderId?: string;
  supplierId?: string;
  totalAmount: number;
  status: "draft" | "approved";
  notes?: string;
  receivedByUserId?: string;
  approvedByUserId?: string;
  approvedDate?: string; // ISO 8601 datetime
  createdByUserId?: string;
  createdDate: string; // ISO 8601 datetime
  modifiedDate: string; // ISO 8601 datetime
  kitchen?: Kitchen;
  supplier?: Supplier;
  order?: Order;
  receivedBy?: User;
  approvedBy?: User;
  createdBy?: User;
  importDetails?: InventoryImportDetail[];
}
```

### InventoryImportDetail

```typescript
interface InventoryImportDetail {
  importDetailId: number;
  importId: string;
  ingredientId: string;
  quantity: number;
  unit: string;
  unitPrice: number;
  totalPrice: number;
  expiryDate?: string; // ISO 8601 date
  batchNumber?: string;
  notes?: string;
  createdDate: string; // ISO 8601 datetime
  modifiedDate: string; // ISO 8601 datetime
  ingredient?: Ingredient;
}
```

### InventoryExport

```typescript
interface InventoryExport {
  exportId: string; // Format: EX/TR/DS + YYYYMMDD-XXXXX
  kitchenId: string;
  exportDate: string; // ISO 8601 date
  exportType: "production" | "transfer" | "disposal" | "return" | "sample";
  destinationKitchenId?: string; // Required for "transfer" type
  orderId?: string;
  totalAmount: number;
  status: "draft" | "approved";
  notes?: string;
  issuedByUserId?: string;
  approvedByUserId?: string;
  approvedDate?: string; // ISO 8601 datetime
  createdByUserId?: string;
  createdDate: string; // ISO 8601 datetime
  modifiedDate: string; // ISO 8601 datetime
  kitchen?: Kitchen;
  destinationKitchen?: Kitchen;
  order?: Order;
  issuedBy?: User;
  approvedBy?: User;
  createdBy?: User;
  exportDetails?: InventoryExportDetail[];
}
```

### InventoryExportDetail

```typescript
interface InventoryExportDetail {
  exportDetailId: number;
  exportId: string;
  ingredientId: string;
  quantity: number;
  unit: string;
  unitCost?: number;
  totalCost?: number;
  batchNumber?: string;
  notes?: string;
  createdDate: string; // ISO 8601 datetime
  modifiedDate: string; // ISO 8601 datetime
  ingredient?: Ingredient;
}
```

### InventoryTransaction

```typescript
interface InventoryTransaction {
  transactionId: number;
  kitchenId: string;
  ingredientId: string;
  transactionType: "IMPORT" | "EXPORT" | "TRANSFER_IN" | "TRANSFER_OUT";
  transactionDate: string; // ISO 8601 datetime
  quantity: number; // Negative for exports
  unit: string;
  quantityBefore: number;
  quantityAfter: number;
  referenceType?: string;
  referenceId?: string;
  notes?: string;
  createdByUserId?: string;
  createdDate: string; // ISO 8601 datetime
  kitchen?: Kitchen;
  ingredient?: Ingredient;
  createdBy?: User;
}
```

---

## Error Handling

### HTTP Status Codes

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data or business rule violation
- `401 Unauthorized` - Missing or invalid authentication token
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

### Error Response Format

All error responses follow this format:

```json
{
  "error": "Error message in Vietnamese"
}
```

### Common Error Messages

| Error Message | Description | HTTP Status |
|--------------|-------------|-------------|
| `Không tìm thấy tồn kho` | Stock not found | 404 |
| `Cần có kitchen_id và ingredient_id` | Missing required parameters | 400 |
| `Không thể sửa phiếu nhập đã duyệt` | Cannot edit approved import | 400 |
| `Phiếu nhập đã được duyệt` | Import already approved | 400 |
| `Số lượng tồn kho không đủ` | Insufficient stock quantity | 400 |
| `Nguyên liệu không tồn tại trong kho` | Ingredient not in stock | 400 |
| `Loại xuất kho không hợp lệ` | Invalid export type | 400 |
| `Định dạng ngày không hợp lệ` | Invalid date format | 400 |

### Validation Errors

When validation fails, the error response includes the validation error message:

```json
{
  "error": "Key: 'CreateImportRequest.ImportDetails' Error:Field validation for 'ImportDetails' failed on the 'required' tag"
}
```

---

## Important Notes

### Pagination

All paginated endpoints use a consistent pagination format:

**Query Parameters:**
- `page` - Page number (default: 1)
- `per_page` - Items per page (default varies by endpoint)
- `search` - Search term for filtering
- `sort_by` - Field to sort by (varies by endpoint)
- `sort_dir` - Sort direction: "asc" or "desc" (default: "asc")

**Response Format:**
All paginated responses follow the `ResourceCollection` format:
```json
{
  "data": [...],
  "meta": {
    "current_page": 1,
    "last_page": 10,
    "from": 1,
    "to": 20,
    "per_page": 20,
    "total": 200
  }
}
```

### Date Formats

- All date fields in request bodies should be in `YYYY-MM-DD` format
- All datetime fields in responses are in ISO 8601 format (e.g., `2024-01-15T10:30:00Z`)

### ID Generation

- **Import IDs**: Format `IMYYYYMMDD-XXXXX` (e.g., `IM20240115-12345`)
- **Export IDs**: 
  - Production: `EXYYYYMMDD-XXXXX`
  - Transfer: `TRYYYYMMDD-XXXXX`
  - Disposal: `DSYYYYMMDD-XXXXX`
  - Other types: `EXYYYYMMDD-XXXXX`

### Status Workflow

**Import/Export Status:**
- `draft` - Can be edited or deleted
- `approved` - Cannot be edited or deleted, stock has been updated

### Stock Updates

- Stock quantities are only updated when an import/export is **approved**
- Draft imports/exports do not affect stock levels
- Approving an import increases stock
- Approving an export decreases stock
- Transfer exports automatically update both source and destination kitchen stocks

### Transaction Logging

All approved imports/exports create transaction logs that can be queried via the stock transactions endpoint.

---

## Example Integration Flow

### Creating and Approving an Import

1. **Create Import (Draft)**
   ```http
   POST /api/inventory/imports
   ```
   Returns: `importId` with status "draft"

2. **Review Import**
   ```http
   GET /api/inventory/imports/{importId}
   ```

3. **Approve Import** (updates stock)
   ```http
   POST /api/inventory/imports/{importId}/approve
   ```
   Stock quantities are now updated, transaction logs created

### Creating and Approving an Export

1. **Check Stock Availability**
   ```http
   GET /api/inventory/stocks/query?kitchen_id=K001&ingredient_id=NL001
   ```

2. **Create Export (Draft)**
   ```http
   POST /api/inventory/exports
   ```
   Returns: `exportId` with status "draft"

3. **Approve Export** (validates and updates stock)
   ```http
   POST /api/inventory/exports/{exportId}/approve
   ```
   Stock quantities are decreased, transaction logs created

---

## Support

For questions or issues, please contact the backend development team.

