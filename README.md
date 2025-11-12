# Go E-commerce API

ໂປຣເຈັກ API ສຳລັບ E-commerce ທີ່ສ້າງດ້ວຍ Go, Gin Framework, ແລະ MySQL

## 🏗️ Project Structure

```
go-api/
├── main_new.go           # Main application entry point
├── go.mod               # Go module file
├── go.sum               # Go dependencies checksum
├── models/              # Data models
│   └── models.go        # All struct definitions
├── handlers/            # API handlers
│   ├── auth.go          # Authentication handlers
│   ├── category.go      # Category CRUD handlers
│   ├── product.go       # Product CRUD handlers
│   ├── customer.go      # Customer CRUD handlers
│   └── order.go         # Order CRUD handlers
├── middleware/          # Middleware functions
│   └── auth.go          # JWT authentication middleware
├── database/            # Database connection
│   └── database.go      # Database initialization
├── utils/               # Utility functions
│   └── auth.go          # Authentication utilities
└── README.md            # This file
```

## 🚀 Features

### Models
- **Category** - ປະເພດສິນຄ້າ
- **Product** - ສິນຄ້າ
- **User** - ຜູ້ໃຊ້ງານລະບົບ
- **Customer** - ລູກຄ້າ
- **Order** - ຄຳສັ່ງຊື້
- **OrderItem** - ລາຍການໃນຄຳສັ່ງຊື້

### API Endpoints

#### 🔐 Authentication
- `POST /register` - ລົງທະບຽນຜູ້ໃຊ້ໃໝ່
- `POST /login` - ເຂົ້າສູ່ລະບົບ

#### 📦 Categories (Public GET, Protected POST/PUT/DELETE)
- `GET /categories` - ເບິ່ງທັງໝົດ
- `GET /categories/:id` - ເບິ່ງດຽວ
- `POST /categories` - ສ້າງໃໝ່ 🔒
- `PUT /categories/:id` - ແກ້ໄຂ 🔒
- `DELETE /categories/:id` - ລົບ 🔒

#### 🛍️ Products (Public GET, Protected POST/PUT/DELETE)
- `GET /products` - ເບິ່ງທັງໝົດ
- `GET /products/:id` - ເບິ່ງດຽວ
- `POST /products` - ສ້າງໃໝ່ 🔒
- `PUT /products/:id` - ແກ້ໄຂ 🔒
- `DELETE /products/:id` - ລົບ 🔒

#### 👥 Customers (Public GET, Protected POST/PUT/DELETE)
- `GET /customers` - ເບິ່ງທັງໝົດ
- `GET /customers/:id` - ເບິ່ງດຽວ (ພ້ອມ orders)
- `POST /customers` - ສ້າງໃໝ່ 🔒
- `PUT /customers/:id` - ແກ້ໄຂ 🔒
- `DELETE /customers/:id` - ລົບ 🔒

#### 📋 Orders (Public GET, Protected POST/PUT/DELETE)
- `GET /orders` - ເບິ່ງທັງໝົດ (ພ້ອມ customer ແລະ items)
- `GET /orders/:id` - ເບິ່ງດຽວ (ພ້ອມ customer ແລະ items)
- `POST /orders` - ສ້າງໃໝ່ 🔒
- `PUT /orders/:id/status` - ອັບເດດ status 🔒
- `DELETE /orders/:id` - ລົບ 🔒

#### 🔧 Utility
- `GET /health` - Health check

## 🔑 Authentication

### JWT Token Usage
```bash
# Login to get token
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password123"}'

# Use token in protected endpoints
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{"name": "Product Name", "price": 1000}'
```

### Order Status Values
- `pending` - ລໍຖ້າການອະນຸມັດ
- `processing` - ກຳລັງດຳເນີນການ
- `shipped` - ສົ່ງແລ້ວ
- `delivered` - ສົ່ງເຖິງແລ້ວ
- `cancelled` - ຍົກເລີກ

## 🗄️ Database Setup

### Requirements
- MySQL/MariaDB
- XAMPP (recommended for development)

### Environment Variables (Optional)
```bash
DB_USER=root
DB_PASS=
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=go_api_db
JWT_SECRET=your-secret-key
```

### Database Schema
API ຈະສ້າງ tables ອັດຕະໂນມັດເມື່ອເລີ່ມ server:
- `categories`
- `products`
- `users`
- `customers`
- `orders`
- `order_items`

## 🛠️ Installation & Usage

### Prerequisites
- Go 1.24.3+
- MySQL/MariaDB
- Git

### Setup
```bash
# Clone repository
git clone <repository-url>
cd go-api

# Install dependencies
go mod tidy

# Start MySQL server (XAMPP)
# Create database: go_api_db

# Run application
go run main_new.go
```

### Development vs Production

#### Development (main_new.go)
```bash
go run main_new.go
```

#### Production (Original main.go)
```bash
go run main.go
```

## 📊 Example API Usage

### 1. Create Customer
```bash
curl -X POST http://localhost:8080/customers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "ສົມຊາຍ ວິລາວັນ",
    "email": "somchai@test.com",
    "phone": "020-1234567",
    "address": "ວຽງຈັນ, ລາວ"
  }'
```

### 2. Create Order
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "customer_id": 1,
    "items": [
      {"product_id": 1, "quantity": 2},
      {"product_id": 2, "quantity": 1}
    ]
  }'
```

### 3. Update Order Status
```bash
curl -X PUT http://localhost:8080/orders/1/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"status": "processing"}'
```

## 🔒 Security Features

- JWT Token Authentication
- Password Hashing (bcrypt)
- Protected Admin Endpoints
- Public Read Access
- Input Validation

## 📝 Notes

- GET endpoints ເປີດໃຫ້ທຸກຄົນເຂົ້າເຖິງໄດ້
- POST/PUT/DELETE endpoints ຕ້ອງການ JWT token
- Token ໝົດອາຍຸໃນ 7 ວັນ
- ການສ້າງ Order ໃຊ້ database transactions
- Price ໃນ OrderItem ເກັບລາຄາໃນຕອນທີ່ສັ່ງຊື້

## 🚀 Production Deployment

1. Set environment variables
2. Use `gin.SetMode(gin.ReleaseMode)`
3. Set proper JWT secret
4. Configure database connection
5. Use reverse proxy (nginx)
6. Enable HTTPS

---

**Created with ❤️ using Go, Gin, GORM, and MySQL**