# Project Structure Documentation

## 📁 Directory Structure

```
go-api/
├── main.go                 # Main application (modular version)
├── main_old.go             # Original monolithic version
├── go.mod                  # Go module definition
├── go.sum                  # Dependencies checksum
├── README.md               # Main documentation
├── STRUCTURE.md            # This file
│
├── models/                 # Data Models Package
│   └── models.go          # All struct definitions and input types
│
├── handlers/               # API Handlers Package
│   ├── auth.go            # Authentication handlers (register, login)
│   ├── category.go        # Category CRUD operations
│   ├── product.go         # Product CRUD operations
│   ├── customer.go        # Customer CRUD operations
│   └── order.go           # Order CRUD operations
│
├── middleware/             # Middleware Package
│   └── auth.go            # JWT authentication middleware
│
├── database/               # Database Package
│   └── database.go        # Database connection and migration
│
└── utils/                  # Utilities Package
    └── auth.go            # Authentication utilities (hash, token)
```

## 🔧 Package Responsibilities

### `models/` Package
- **Purpose**: ກຳນົດ data structures ແລະ input validation
- **Files**: `models.go`
- **Contains**:
  - Database models (Category, Product, User, Customer, Order, OrderItem)
  - Input structs (RegisterInput, LoginInput, CreateOrderInput, etc.)
  - JSON tags ແລະ GORM tags

### `handlers/` Package
- **Purpose**: ຈັດການ HTTP requests ແລະ responses
- **Files**: ແຍກຕາມ domain (auth, category, product, customer, order)
- **Contains**:
  - HTTP handler functions
  - Business logic
  - Database operations
  - Response formatting

### `middleware/` Package
- **Purpose**: ປະຕິບັດ middleware functions
- **Files**: `auth.go`
- **Contains**:
  - JWT authentication middleware
  - Request validation
  - Authorization checks

### `database/` Package
- **Purpose**: ຈັດການ database connection ແລະ configuration
- **Files**: `database.go`
- **Contains**:
  - Database connection setup
  - Auto-migration
  - Database configuration
  - Global DB instance

### `utils/` Package
- **Purpose**: ຟັງຊັນຊ່ວຍເຫຼືອ
- **Files**: `auth.go`
- **Contains**:
  - Password hashing
  - JWT token generation
  - Authentication utilities

### `main.go`
- **Purpose**: Application entry point ແລະ routing
- **Contains**:
  - Database initialization
  - Route definitions
  - Middleware registration
  - Server startup

## 🔄 Data Flow

```
Request → main.go → middleware → handlers → database → models
                ↓
Response ← JSON ← handlers ← database ← models
```

### Example Flow for POST /orders:
1. **main.go**: Route to `handlers.CreateOrder`
2. **middleware**: Check JWT token
3. **handlers/order.go**: Validate input, business logic
4. **database**: Execute database operations
5. **models**: Data structure validation
6. **Response**: JSON response to client

## 🚀 Benefits of This Structure

### 1. **Separation of Concerns**
- ແຍກການຈັດການຂໍ້ມູນຈາກ business logic
- ແຍກ authentication ຈາກ API handlers
- ແຍກ database logic ຈາກ HTTP handling

### 2. **Maintainability**
- ງ່າຍຕໍ່ການແກ້ໄຂ ແລະ debug
- ແຍກຟາຍລ໌ຕາມຫນ້າທີ່
- ງ່າຍຕໍ່ການທົດສອບ

### 3. **Scalability**
- ເພີ່ມ handlers ໃໝ່ໄດ້ງ່າຍ
- ເພີ່ມ middleware ໃໝ່ໄດ້ງ່າຍ
- ແຍກການພັດທະນາໃນທີມງານ

### 4. **Reusability**
- Models ໃຊ້ໄດ້ໃນ handlers ຫຼາຍໆຕົວ
- Middleware ໃຊ້ໄດ້ກັບ routes ຫຼາຍໆຕົວ
- Utils ໃຊ້ໄດ້ທົ່ວໂປຣເຈັກ

## 🔧 Development Workflow

### Adding New Feature:
1. **Models**: ເພີ່ມ struct ໃນ `models/models.go`
2. **Database**: Update migration ໃນ `database/database.go`
3. **Handlers**: ເພີ່ມ handler functions
4. **Routes**: ເພີ່ມ routes ໃນ `main.go`
5. **Test**: ທົດສອບ API endpoints

### Modifying Existing Feature:
1. ຫາຟາຍລ໌ທີ່ກ່ຽວຂ້ອງໃນ `handlers/`
2. ແກ້ໄຂ business logic
3. Update models ຖ້າຈຳເປັນ
4. ທົດສອບການແກ້ໄຂ

## 📋 File Naming Conventions

- **Package files**: `package_name.go`
- **Handlers**: `domain.go` (auth.go, category.go, etc.)
- **Main files**: `main.go`, `main_old.go`
- **Documentation**: `README.md`, `STRUCTURE.md`

## 🎯 Best Practices Applied

1. **Single Responsibility**: ແຕ່ລະ package ມີຫນ້າທີ່ຊັດເຈນ
2. **Dependency Injection**: ໃຊ້ global DB instance
3. **Error Handling**: ຈັດການ errors ຢ່າງສອດຄ່ອງ
4. **Input Validation**: ໃຊ້ Gin binding tags
5. **Security**: JWT middleware ສຳລັບ protected routes

---

**This structure makes the codebase professional, maintainable, and scalable! 🚀**
