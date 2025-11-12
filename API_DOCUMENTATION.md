# API Documentation

Base URL: `http://localhost:8081`

## Authentication

All protected routes require a JWT token in the `Authorization` header:
```
Authorization: Bearer <your-jwt-token>
```

---

## 1. Health Check

### GET /health
Check if the API is running.

**Response:**
```json
{
  "ok": true
}
```

---

## 2. Authentication Endpoints

### POST /register
Register a new user.

**Request Body:**
```json
{
  "username": "string (required)",
  "email": "string (required, valid email)",
  "password": "string (required, min 6 characters)"
}
```

**Response:** `201 Created`
```json
{
  "message": "ລົງທະບຽນສຳເລັດ",
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com"
  },
  "token": "jwt-token-string"
}
```

### POST /login
Login and get JWT token.

**Request Body:**
```json
{
  "username": "string (required)",
  "password": "string (required)"
}
```

**Response:** `200 OK`
```json
{
  "message": "ເຂົ້າສູ່ລະບົບສຳເລັດ",
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com"
  },
  "token": "jwt-token-string"
}
```

---

## 3. Category Endpoints

### GET /categories
Get all categories.

**No Authentication Required**

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "name": "Electronics",
    "description": "Electronic products",
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

### GET /categories/:id
Get a single category by ID.

**No Authentication Required**

**Response:** `200 OK`
```json
{
  "id": 1,
  "name": "Electronics",
  "description": "Electronic products",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### POST /categories
Create a new category.

**Authentication Required**

**Request Body:**
```json
{
  "name": "string (required, unique)",
  "description": "string (optional)"
}
```

**Response:** `201 Created`
```json
{
  "id": 1,
  "name": "Electronics",
  "description": "Electronic products",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### PUT /categories/:id
Update a category.

**Authentication Required**

**Request Body:**
```json
{
  "name": "string (optional)",
  "description": "string (optional)"
}
```

**Response:** `200 OK`
```json
{
  "id": 1,
  "name": "Updated Name",
  "description": "Updated description",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### DELETE /categories/:id
Delete a category.

**Authentication Required**

**Response:** `204 No Content`

---

## 4. Product Endpoints

### GET /products
Get all products with their categories.

**No Authentication Required**

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "name": "Product Name",
    "price": 1000,
    "image": "/uploads/image.jpg",
    "category_id": 1,
    "category": {
      "id": 1,
      "name": "Electronics",
      "description": "Electronic products",
      "created_at": "2024-01-01T00:00:00Z"
    }
  }
]
```

### GET /products/:id
Get a single product by ID with category.

**No Authentication Required**

**Response:** `200 OK`
```json
{
  "id": 1,
  "name": "Product Name",
  "price": 1000,
  "image": "/uploads/image.jpg",
  "category_id": 1,
  "category": {
    "id": 1,
    "name": "Electronics",
    "description": "Electronic products",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### POST /products
Create a new product.

**Authentication Required**

**Supports both JSON and multipart/form-data**

**JSON Request:**
```json
{
  "name": "string (required)",
  "price": 1000,
  "category_id": 1,
  "image": "string (optional, URL)"
}
```

**Multipart/Form-Data Request:**
- `name`: string (required)
- `price`: string (required, will be converted to int)
- `category_id`: string (optional, will be converted to uint)
- `image`: file (optional, image file)

**Response:** `201 Created`
```json
{
  "id": 1,
  "name": "Product Name",
  "price": 1000,
  "image": "/uploads/image.jpg",
  "category_id": 1
}
```

### PUT /products/:id
Update a product.

**Authentication Required**

**Supports both JSON and multipart/form-data**

**JSON Request:**
```json
{
  "name": "string (optional)",
  "price": 1000,
  "category_id": 1,
  "image": "string (optional, URL)"
}
```

**Multipart/Form-Data Request:**
- `name`: string (optional)
- `price`: string (optional)
- `category_id`: string (optional)
- `image`: file (optional, image file)

**Response:** `200 OK`
```json
{
  "id": 1,
  "name": "Updated Product Name",
  "price": 1500,
  "image": "/uploads/new-image.jpg",
  "category_id": 1,
  "category": {
    "id": 1,
    "name": "Electronics",
    "description": "Electronic products",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### DELETE /products/:id
Delete a product.

**Authentication Required**

**Response:** `204 No Content`

---

## 5. Customer Endpoints

### GET /customers
Get all customers.

**No Authentication Required**

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "1234567890",
    "address": "123 Main St",
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

### GET /customers/:id
Get a single customer by ID.

**No Authentication Required**

**Response:** `200 OK`
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "1234567890",
  "address": "123 Main St",
  "created_at": "2024-01-01T00:00:00Z",
  "orders": []
}
```

### POST /customers
Create a new customer.

**Authentication Required**

**Request Body:**
```json
{
  "name": "string (required)",
  "email": "string (required, unique)",
  "phone": "string (optional)",
  "address": "string (optional)"
}
```

**Response:** `201 Created`
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "1234567890",
  "address": "123 Main St",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### PUT /customers/:id
Update a customer.

**Authentication Required**

**Request Body:**
```json
{
  "name": "string (optional)",
  "email": "string (optional)",
  "phone": "string (optional)",
  "address": "string (optional)"
}
```

**Response:** `200 OK`
```json
{
  "id": 1,
  "name": "Updated Name",
  "email": "updated@example.com",
  "phone": "9876543210",
  "address": "456 New St",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### DELETE /customers/:id
Delete a customer.

**Authentication Required**

**Response:** `204 No Content`

---

## 6. Order Endpoints

### GET /orders
Get all orders with customer and order items.

**No Authentication Required**

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "customer_id": 1,
    "customer": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com"
    },
    "status": "pending",
    "total_amount": 5000,
    "order_items": [
      {
        "id": 1,
        "order_id": 1,
        "product_id": 1,
        "product": {
          "id": 1,
          "name": "Product Name",
          "price": 1000
        },
        "quantity": 5,
        "price": 1000
      }
    ],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

### GET /orders/:id
Get a single order by ID with customer and order items.

**No Authentication Required**

**Response:** `200 OK`
```json
{
  "id": 1,
  "customer_id": 1,
  "customer": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  },
  "status": "pending",
  "total_amount": 5000,
  "order_items": [
    {
      "id": 1,
      "order_id": 1,
      "product_id": 1,
      "product": {
        "id": 1,
        "name": "Product Name",
        "price": 1000
      },
      "quantity": 5,
      "price": 1000
    }
  ],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### POST /orders
Create a new order.

**Authentication Required**

**Request Body:**
```json
{
  "customer_id": 1,
  "items": [
    {
      "product_id": 1,
      "quantity": 5
    },
    {
      "product_id": 2,
      "quantity": 3
    }
  ]
}
```

**Response:** `201 Created`
```json
{
  "id": 1,
  "customer_id": 1,
  "customer": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  },
  "status": "pending",
  "total_amount": 5000,
  "order_items": [
    {
      "id": 1,
      "order_id": 1,
      "product_id": 1,
      "product": {
        "id": 1,
        "name": "Product Name",
        "price": 1000
      },
      "quantity": 5,
      "price": 1000
    }
  ],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### PUT /orders/:id/status
Update order status.

**Authentication Required**

**Request Body:**
```json
{
  "status": "pending|processing|shipped|delivered|cancelled"
}
```

**Response:** `200 OK`
```json
{
  "id": 1,
  "customer_id": 1,
  "status": "processing",
  "total_amount": 5000,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T01:00:00Z"
}
```

### DELETE /orders/:id
Delete an order.

**Authentication Required**

**Response:** `204 No Content`

---

## 7. Static Files

### GET /uploads/:filename
Access uploaded images.

**No Authentication Required**

---

## Error Responses

### 400 Bad Request
```json
{
  "error": "error message"
}
```

### 401 Unauthorized
```json
{
  "error": "Unauthorized",
  "message": "error details"
}
```

### 404 Not Found
```json
{
  "error": "not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "error message"
}
```

---

## Notes

- All timestamps are in ISO 8601 format
- JWT tokens expire after 7 days
- Product image uploads are stored in the `./uploads` directory
- All protected routes require valid JWT token in `Authorization: Bearer <token>` header
- CORS is enabled for all origins

