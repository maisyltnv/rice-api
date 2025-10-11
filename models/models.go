package models

import "time"

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
