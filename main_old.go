package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Category struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Product struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name"`
	Price      int       `json:"price"`
	CategoryID *uint     `json:"category_id"`                                     // ໃຊ້ pointer ເພື່ອໃຫ້ສາມາດເປັນ null ໄດ້
	Category   *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"` // Eager loading
}

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"` // "-" ບໍ່ສົ່ງ password ອອກໄປໃນ JSON
	CreatedAt time.Time `json:"created_at"`
}

type Customer struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	Orders    []Order   `json:"orders,omitempty" gorm:"foreignKey:CustomerID"`
}

type Order struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	CustomerID  uint        `json:"customer_id" gorm:"not null"`
	Customer    *Customer   `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Status      string      `json:"status" gorm:"default:'pending'"` // pending, processing, shipped, delivered, cancelled
	TotalAmount int         `json:"total_amount"`
	OrderItems  []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:OrderID"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID        uint     `json:"id" gorm:"primaryKey"`
	OrderID   uint     `json:"order_id" gorm:"not null"`
	Order     *Order   `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	ProductID uint     `json:"product_id" gorm:"not null"`
	Product   *Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Quantity  int      `json:"quantity" gorm:"not null"`
	Price     int      `json:"price" gorm:"not null"` // ລາຄາໃນຕອນທີ່ສັ່ງຊື້
}

// Struct ສຳລັບຮັບຂໍ້ມູນການລົງທະບຽນ
type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Struct ສຳລັບຮັບຂໍ້ມູນການເຂົ້າສູ່ລະບົບ
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Struct ສຳລັບຮັບຂໍ້ມູນການສ້າງ Order
type CreateOrderInput struct {
	CustomerID uint                   `json:"customer_id" binding:"required"`
	Items      []CreateOrderItemInput `json:"items" binding:"required,min=1"`
}

type CreateOrderItemInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

// Struct ສຳລັບອັບເດດ Order status
type UpdateOrderStatusInput struct {
	Status string `json:"status" binding:"required,oneof=pending processing shipped delivered cancelled"`
}

// JWT secret key (ໃນໂປຣເຈັກຈິງຄວນເກັບໃນ environment variable)
var jwtSecret = []byte(getenv("JWT_SECRET", "my-secret-key-change-in-production"))

func dsn() string {
	// ปรับค่าตาม XAMPP ของคุณ
	user := getenv("DB_USER", "root")
	pass := getenv("DB_PASS", "")          // ค่าเริ่มต้น XAMPP ส่วนมากว่าง
	host := getenv("DB_HOST", "127.0.0.1") // หลีกเลี่ยง "localhost" เพื่อไม่ให้ใช้ socket
	port := getenv("DB_PORT", "3306")
	name := getenv("DB_NAME", "go_api_db")

	// parseTime=1 ให้ scan เป็น time.Time ได้, charset utf8mb4 รองรับ emoji/ພາສາລາວ
	return user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + name + "?parseTime=true&charset=utf8mb4&loc=Local"
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// Hash password ດ້ວຍ bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// ກວດສອບວ່າ password ຖືກຕ້ອງບໍ່
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ສ້າງ JWT token
func generateToken(userID uint, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // Token ໝົດອາຍຸໃນ 7 ວັນ
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Middleware ສຳລັບກວດສອບ JWT Token
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ດຶງ token ຈາກ Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ບໍ່ມີ token, ກະລຸນາເຂົ້າສູ່ລະບົບກ່ອນ"})
			c.Abort()
			return
		}

		// ກຳจັດ "Bearer " ອອກຈາກ token (ຖ້າມີ)
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// Validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// ກວດສອບວ່າໃຊ້ signing method ທີ່ຖືກຕ້ອງ
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token ບໍ່ຖືກຕ້ອງຫຼືໝົດອາຍຸແລ້ວ"})
			c.Abort()
			return
		}

		// ເອົາຂໍ້ມູນຈາກ token ໃສ່ context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", uint(claims["user_id"].(float64)))
			c.Set("username", claims["username"].(string))
		}

		c.Next()
	}
}

func main() {
	// Connect DB
	db, err := gorm.Open(mysql.Open(dsn()), &gorm.Config{})
	if err != nil {
		log.Fatal("cannot connect database: ", err)
	}
	// Auto-migrate
	if err := db.AutoMigrate(&Category{}, &Product{}, &User{}, &Customer{}, &Order{}, &OrderItem{}); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	// health
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	// AUTH - ການລົງທະບຽນ
	r.POST("/register", func(c *gin.Context) {
		var input RegisterInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ກວດສອບວ່າມີ username ຫຼື email ຊ້ຳບໍ່
		var existingUser User
		if err := db.Where("username = ? OR email = ?", input.Username, input.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "username ຫຼື email ຖືກໃຊ້ແລ້ວ"})
			return
		}

		// Hash password
		hashedPassword, err := hashPassword(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ບໍ່ສາມາດ hash password ໄດ້"})
			return
		}

		// ສ້າງຜູ້ໃຊ້ໃໝ່
		user := User{
			Username: input.Username,
			Email:    input.Email,
			Password: hashedPassword,
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// ສ້າງ token
		token, err := generateToken(user.ID, user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ບໍ່ສາມາດສ້າງ token ໄດ້"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "ລົງທະບຽນສຳເລັດ",
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
			},
			"token": token,
		})
	})

	// AUTH - ເຂົ້າສູ່ລະບົບ
	r.POST("/login", func(c *gin.Context) {
		var input LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ຊອກຫາຜູ້ໃຊ້
		var user User
		if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "username ຫຼື password ບໍ່ຖືກຕ້ອງ"})
			return
		}

		// ກວດສອບ password
		if !checkPasswordHash(input.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "username ຫຼື password ບໍ່ຖືກຕ້ອງ"})
			return
		}

		// ສ້າງ token
		token, err := generateToken(user.ID, user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ບໍ່ສາມາດສ້າງ token ໄດ້"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "ເຂົ້າສູ່ລະບົບສຳເລັດ",
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
			},
			"token": token,
		})
	})

	// CATEGORY CRUD
	// ເບິ່ງ categories ທັງໝົດ
	r.GET("/categories", func(c *gin.Context) {
		var items []Category
		if err := db.Find(&items).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, items)
	})

	// ເບິ່ງ category ດຽວ
	r.GET("/categories/:id", func(c *gin.Context) {
		var cat Category
		if err := db.First(&cat, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		c.JSON(http.StatusOK, cat)
	})

	// ສ້າງ category ໃໝ່ (ຕ້ອງມີ token)
	r.POST("/categories", authMiddleware(), func(c *gin.Context) {
		var cat Category
		if err := c.BindJSON(&cat); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&cat).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, cat)
	})

	// ແກ້ໄຂ category (ຕ້ອງມີ token)
	r.PUT("/categories/:id", authMiddleware(), func(c *gin.Context) {
		var cat Category
		if err := db.First(&cat, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		var input Category
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cat.Name = input.Name
		cat.Description = input.Description
		if err := db.Save(&cat).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, cat)
	})

	// ລົບ category (ຕ້ອງມີ token)
	r.DELETE("/categories/:id", authMiddleware(), func(c *gin.Context) {
		result := db.Delete(&Category{}, c.Param("id"))
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		c.Status(http.StatusNoContent)
	})

	// PRODUCT CRUD
	r.GET("/products", func(c *gin.Context) {
		var items []Product
		// Preload Category ເພື່ອສະແດງຂໍ້ມູນ category ພ້ອມກັບ product
		if err := db.Preload("Category").Find(&items).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, items)
	})

	r.GET("/products/:id", func(c *gin.Context) {
		var p Product
		// Preload Category ເພື່ອສະແດງຂໍ້ມູນ category
		if err := db.Preload("Category").First(&p, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, p)
	})

	// ສ້າງ product ໃໝ່ (ຕ້ອງມີ token)
	r.POST("/products", authMiddleware(), func(c *gin.Context) {
		var p Product
		if err := c.BindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&p).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, p)
	})

	// ແກ້ໄຂ product (ຕ້ອງມີ token)
	r.PUT("/products/:id", authMiddleware(), func(c *gin.Context) {
		var p Product
		if err := db.First(&p, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		var in Product
		if err := c.BindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		p.Name, p.Price = in.Name, in.Price
		if err := db.Save(&p).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, p)
	})

	// ລົບ product (ຕ້ອງມີ token)
	r.DELETE("/products/:id", authMiddleware(), func(c *gin.Context) {
		result := db.Delete(&Product{}, c.Param("id"))
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.Status(http.StatusNoContent)
	})

	// CUSTOMER CRUD
	// ເບິ່ງ customers ທັງໝົດ
	r.GET("/customers", func(c *gin.Context) {
		var items []Customer
		if err := db.Find(&items).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, items)
	})

	// ເບິ່ງ customer ດຽວ
	r.GET("/customers/:id", func(c *gin.Context) {
		var customer Customer
		if err := db.Preload("Orders.OrderItems.Product").First(&customer, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
			return
		}
		c.JSON(http.StatusOK, customer)
	})

	// ສ້າງ customer ໃໝ່ (ຕ້ອງມີ token)
	r.POST("/customers", authMiddleware(), func(c *gin.Context) {
		var customer Customer
		if err := c.BindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&customer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, customer)
	})

	// ແກ້ໄຂ customer (ຕ້ອງມີ token)
	r.PUT("/customers/:id", authMiddleware(), func(c *gin.Context) {
		var customer Customer
		if err := db.First(&customer, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
			return
		}
		var input Customer
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		customer.Name = input.Name
		customer.Email = input.Email
		customer.Phone = input.Phone
		customer.Address = input.Address
		if err := db.Save(&customer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, customer)
	})

	// ລົບ customer (ຕ້ອງມີ token)
	r.DELETE("/customers/:id", authMiddleware(), func(c *gin.Context) {
		result := db.Delete(&Customer{}, c.Param("id"))
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
			return
		}
		c.Status(http.StatusNoContent)
	})

	// ORDER CRUD
	// ເບິ່ງ orders ທັງໝົດ
	r.GET("/orders", func(c *gin.Context) {
		var items []Order
		if err := db.Preload("Customer").Preload("OrderItems.Product").Find(&items).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, items)
	})

	// ເບິ່ງ order ດຽວ
	r.GET("/orders/:id", func(c *gin.Context) {
		var order Order
		if err := db.Preload("Customer").Preload("OrderItems.Product").First(&order, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		c.JSON(http.StatusOK, order)
	})

	// ສ້າງ order ໃໝ່ (ຕ້ອງມີ token)
	r.POST("/orders", authMiddleware(), func(c *gin.Context) {
		var input CreateOrderInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ກວດສອບວ່າ customer ມີຢູ່ບໍ່
		var customer Customer
		if err := db.First(&customer, input.CustomerID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "customer not found"})
			return
		}

		// ສ້າງ order ໃໝ່
		order := Order{
			CustomerID: input.CustomerID,
			Status:     "pending",
		}

		// ເລີ່ມ transaction
		tx := db.Begin()
		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		totalAmount := 0
		// ສ້າງ order items
		for _, item := range input.Items {
			// ກວດສອບວ່າ product ມີຢູ່ບໍ່
			var product Product
			if err := tx.First(&product, item.ProductID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
				return
			}

			orderItem := OrderItem{
				OrderID:   order.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     product.Price,
			}

			if err := tx.Create(&orderItem).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			totalAmount += product.Price * item.Quantity
		}

		// ອັບເດດ total amount
		order.TotalAmount = totalAmount
		if err := tx.Save(&order).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		tx.Commit()

		// ໂຫຼດ order ພ້ອມກັບ relationships
		if err := db.Preload("Customer").Preload("OrderItems.Product").First(&order, order.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, order)
	})

	// ອັບເດດ order status (ຕ້ອງມີ token)
	r.PUT("/orders/:id/status", authMiddleware(), func(c *gin.Context) {
		var order Order
		if err := db.First(&order, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}

		var input UpdateOrderStatusInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		order.Status = input.Status
		if err := db.Save(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// ໂຫຼດ order ພ້ອມກັບ relationships
		if err := db.Preload("Customer").Preload("OrderItems.Product").First(&order, order.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, order)
	})

	// ລົບ order (ຕ້ອງມີ token)
	r.DELETE("/orders/:id", authMiddleware(), func(c *gin.Context) {
		// ເລີ່ມ transaction ເພື່ອລົບ order items ກ່ອນ
		tx := db.Begin()

		// ລົບ order items ກ່ອນ
		if err := tx.Where("order_id = ?", c.Param("id")).Delete(&OrderItem{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// ລົບ order
		result := tx.Delete(&Order{}, c.Param("id"))
		if result.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		if result.RowsAffected == 0 {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}

		tx.Commit()
		c.Status(http.StatusNoContent)
	})

	r.Run(":8080")
}
