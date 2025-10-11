package database

import (
	"log"
	"os"

	"example.com/go-xampp-api/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Connect DB
	db, err := gorm.Open(mysql.Open(dsn()), &gorm.Config{})
	if err != nil {
		log.Fatal("cannot connect database: ", err)
	}

	DB = db

	// Auto-migrate
	if err := DB.AutoMigrate(&models.Category{}, &models.Product{}, &models.User{}, &models.Customer{}, &models.Order{}, &models.OrderItem{}); err != nil {
		log.Fatal(err)
	}
}

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
