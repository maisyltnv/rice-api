package handlers

import (
	"net/http"

	"example.com/go-xampp-api/database"
	"example.com/go-xampp-api/models"
	"github.com/gin-gonic/gin"
)

// ເບິ່ງ customers ທັງໝົດ
func GetCustomers(c *gin.Context) {
	var items []models.Customer
	if err := database.DB.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// ເບິ່ງ customer ດຽວ
func GetCustomer(c *gin.Context) {
	var customer models.Customer
	if err := database.DB.Preload("Orders.OrderItems.Product").First(&customer, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}
	c.JSON(http.StatusOK, customer)
}

// ສ້າງ customer ໃໝ່
func CreateCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.BindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, customer)
}

// ແກ້ໄຂ customer
func UpdateCustomer(c *gin.Context) {
	var customer models.Customer
	if err := database.DB.First(&customer, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}
	var input models.Customer
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customer.Name = input.Name
	customer.Email = input.Email
	customer.Phone = input.Phone
	customer.Address = input.Address
	if err := database.DB.Save(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customer)
}

// ລົບ customer
func DeleteCustomer(c *gin.Context) {
	result := database.DB.Delete(&models.Customer{}, c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
