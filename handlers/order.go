package handlers

import (
	"net/http"

	"example.com/go-xampp-api/database"
	"example.com/go-xampp-api/models"
	"github.com/gin-gonic/gin"
)

// ເບິ່ງ orders ທັງໝົດ
func GetOrders(c *gin.Context) {
	var items []models.Order
	if err := database.DB.Preload("Customer").Preload("OrderItems.Product").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// ເບິ່ງ order ດຽວ
func GetOrder(c *gin.Context) {
	var order models.Order
	if err := database.DB.Preload("Customer").Preload("OrderItems.Product").First(&order, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

// ສ້າງ order ໃໝ່
func CreateOrder(c *gin.Context) {
	var input models.CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ກວດສອບວ່າ customer ມີຢູ່ບໍ່
	var customer models.Customer
	if err := database.DB.First(&customer, input.CustomerID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer not found"})
		return
	}

	// ສ້າງ order ໃໝ່
	order := models.Order{
		CustomerID: input.CustomerID,
		Status:     "pending",
	}

	// ເລີ່ມ transaction
	tx := database.DB.Begin()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalAmount := 0
	// ສ້າງ order items
	for _, item := range input.Items {
		// ກວດສອບວ່າ product ມີຢູ່ບໍ່
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
			return
		}

		orderItem := models.OrderItem{
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
	if err := database.DB.Preload("Customer").Preload("OrderItems.Product").First(&order, order.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// ອັບເດດ order status
func UpdateOrderStatus(c *gin.Context) {
	var order models.Order
	if err := database.DB.First(&order, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	var input models.UpdateOrderStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order.Status = input.Status
	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ໂຫຼດ order ພ້ອມກັບ relationships
	if err := database.DB.Preload("Customer").Preload("OrderItems.Product").First(&order, order.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// ລົບ order
func DeleteOrder(c *gin.Context) {
	// ເລີ່ມ transaction ເພື່ອລົບ order items ກ່ອນ
	tx := database.DB.Begin()

	// ລົບ order items ກ່ອນ
	if err := tx.Where("order_id = ?", c.Param("id")).Delete(&models.OrderItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ລົບ order
	result := tx.Delete(&models.Order{}, c.Param("id"))
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
}
