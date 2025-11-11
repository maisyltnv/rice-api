package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"example.com/go-xampp-api/database"
	"example.com/go-xampp-api/models"
	"example.com/go-xampp-api/utils"
	"github.com/gin-gonic/gin"
)

// ເບິ່ງ orders ທັງໝົດ
func GetOrders(c *gin.Context) {
	var items []models.Order
	query := database.DB.Preload("Customer").Preload("OrderItems.Product")

	// Determine access: default to customer-restricted unless role is explicitly admin
	if v, ok := c.Get("role"); ok {
		if role, ok2 := v.(string); ok2 && role == "admin" {
			// admin: no filtering
		} else {
			// customer or unknown role: restrict
			if cid, ok3 := c.Get("customer_id"); ok3 {
				query = query.Where("customer_id = ?", cid)
			} else if uname, ok4 := c.Get("username"); ok4 {
				// Backward compatibility: try to resolve by email from token
				if email, ok5 := uname.(string); ok5 && email != "" {
					var cust models.Customer
					if err := database.DB.Where("email = ?", email).First(&cust).Error; err == nil {
						query = query.Where("customer_id = ?", cust.ID)
					}
				}
			}
		}
	} else {
		// No role in token: attempt email fallback
		if uname, ok := c.Get("username"); ok {
			if email, ok2 := uname.(string); ok2 && email != "" {
				var cust models.Customer
				if err := database.DB.Where("email = ?", email).First(&cust).Error; err == nil {
					query = query.Where("customer_id = ?", cust.ID)
				}
			}
		}
	}

	if err := query.Find(&items).Error; err != nil {
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

	// Enforce customer-only access unless explicit admin role
	if v, ok := c.Get("role"); ok {
		if role, ok2 := v.(string); ok2 && role == "admin" {
			// admin allowed
		} else {
			if cid, ok3 := c.Get("customer_id"); ok3 {
				if order.CustomerID != cid {
					c.JSON(http.StatusForbidden, gin.H{"error": "ບໍ່ມີສິດເຂົ້າເຖິງ order ນີ້"})
					return
				}
			} else if uname, ok4 := c.Get("username"); ok4 {
				if email, ok5 := uname.(string); ok5 && email != "" {
					var cust models.Customer
					if err := database.DB.Where("email = ?", email).First(&cust).Error; err == nil {
						if order.CustomerID != cust.ID {
							c.JSON(http.StatusForbidden, gin.H{"error": "ບໍ່ມີສິດເຂົ້າເຖິງ order ນີ້"})
							return
						}
					}
				}
			}
		}
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

	var customer models.Customer
	var customerID uint

	// ກວດສອບວ່າ customer ເຂົ້າສູ່ລະບົບບໍ່ (ມີ token)
	role, exists := c.Get("role")
	if exists && role == "customer" {
		// ຖ້າເປັນ customer ທີ່ເຂົ້າສູ່ລະບົບ, ໃຊ້ customer_id ຈາກ token
		if cid, ok := c.Get("customer_id"); ok {
			customerID = cid.(uint)
			if err := database.DB.First(&customer, customerID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "customer not found"})
				return
			}
		}
	} else {
		// Guest checkout ຫຼື customer ໃໝ່ - ຊອກຫາຫຼືສ້າງ customer ຈາກ email
		if input.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email is required for guest checkout"})
			return
		}

		// ຊອກຫາ customer ດ້ວຍ email
		err := database.DB.Where("email = ?", input.Email).First(&customer).Error
		if err != nil {
			// ຖ້າບໍ່ພົບ, ສ້າງ customer ໃໝ່ (guest checkout)
			customerName := input.CustomerName
			if customerName == "" {
				customerName = fmt.Sprintf("%s %s", input.FirstName, input.LastName)
				customerName = strings.TrimSpace(customerName)
			}
			if customerName == "" {
				customerName = input.Email // Fallback to email if no name provided
			}

			// Generate a temporary password for guest customers (they can set a real one later if they register)
			// Using a placeholder that won't work for login unless they set a real password
			tempPassword := fmt.Sprintf("guest_%d", time.Now().UnixNano())
			hashedPassword, err := utils.HashPassword(tempPassword)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password: " + err.Error()})
				return
			}

			customer = models.Customer{
				Name:     customerName,
				Email:    input.Email,
				Phone:    "", // Can be added later
				Address:  "", // Can be added later
				Password: hashedPassword,
			}

			// ສ້າງ customer ໃໝ່
			if err := database.DB.Create(&customer).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create customer: " + err.Error()})
				return
			}
		}
		customerID = customer.ID
	}

	// ສ້າງ shipping address string ຈາກ structured object
	shippingAddrParts := []string{}
	if input.ShippingAddress.Street != "" {
		shippingAddrParts = append(shippingAddrParts, input.ShippingAddress.Street)
	}
	if input.ShippingAddress.City != "" {
		shippingAddrParts = append(shippingAddrParts, input.ShippingAddress.City)
	}
	if input.ShippingAddress.State != "" {
		shippingAddrParts = append(shippingAddrParts, input.ShippingAddress.State)
	}
	if input.ShippingAddress.ZipCode != "" {
		shippingAddrParts = append(shippingAddrParts, input.ShippingAddress.ZipCode)
	}
	if input.ShippingAddress.Country != "" {
		shippingAddrParts = append(shippingAddrParts, input.ShippingAddress.Country)
	}
	shippingAddressStr := strings.Join(shippingAddrParts, ", ")

	// ສ້າງ order ໃໝ່
	order := models.Order{
		CustomerID:      customerID,
		Status:          "pending",
		ShippingAddress: shippingAddressStr,
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
			Image:     product.Image,
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
