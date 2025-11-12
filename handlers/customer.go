package handlers

import (
	"net/http"

	"example.com/go-xampp-api/database"
	"example.com/go-xampp-api/models"
	"example.com/go-xampp-api/utils"
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

// ສ້າງ customer ໃໝ່ (Admin only - requires password)
func CreateCustomer(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Phone    string `json:"phone"`
		Address  string `json:"address"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ກວດສອບວ່າມີ email ຊ້ຳບໍ່
	var existingCustomer models.Customer
	if err := database.DB.Where("email = ?", input.Email).First(&existingCustomer).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email ຖືກໃຊ້ແລ້ວ"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ບໍ່ສາມາດ hash password ໄດ້"})
		return
	}

	customer := models.Customer{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		Phone:    input.Phone,
		Address:  input.Address,
	}

	if err := database.DB.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Don't return password in response
	customer.Password = ""
	c.JSON(http.StatusCreated, customer)
}

// ແກ້ໄຂ customer
func UpdateCustomer(c *gin.Context) {
	var customer models.Customer
	if err := database.DB.First(&customer, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}

	var input struct {
		Name     *string `json:"name"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
		Phone    *string `json:"phone"`
		Address  *string `json:"address"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if input.Name != nil {
		customer.Name = *input.Name
	}
	if input.Email != nil {
		// Check if email is already taken by another customer
		var existingCustomer models.Customer
		if err := database.DB.Where("email = ? AND id != ?", *input.Email, customer.ID).First(&existingCustomer).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "email ຖືກໃຊ້ແລ້ວ"})
			return
		}
		customer.Email = *input.Email
	}
	if input.Password != nil {
		// Hash new password
		hashedPassword, err := utils.HashPassword(*input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ບໍ່ສາມາດ hash password ໄດ້"})
			return
		}
		customer.Password = hashedPassword
	}
	if input.Phone != nil {
		customer.Phone = *input.Phone
	}
	if input.Address != nil {
		customer.Address = *input.Address
	}

	if err := database.DB.Save(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Don't return password in response
	customer.Password = ""
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

// ການລົງທະບຽນ Customer
func CustomerRegister(c *gin.Context) {
	var input models.CustomerRegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ກວດສອບວ່າມີ email ຊ້ຳບໍ່
	var existingCustomer models.Customer
	if err := database.DB.Where("email = ?", input.Email).First(&existingCustomer).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email ຖືກໃຊ້ແລ້ວ"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ບໍ່ສາມາດ hash password ໄດ້"})
		return
	}

	// ສ້າງ customer ໃໝ່
	customer := models.Customer{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		Phone:    input.Phone,
		Address:  input.Address,
	}

	if err := database.DB.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ສ້າງ token
	token, err := utils.GenerateToken(customer.ID, customer.Email, "customer")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ບໍ່ສາມາດສ້າງ token ໄດ້"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "ລົງທະບຽນສຳເລັດ",
		"customer": gin.H{
			"id":    customer.ID,
			"name":  customer.Name,
			"email": customer.Email,
			"phone": customer.Phone,
		},
		"token": token,
	})
}

// ເຂົ້າສູ່ລະບົບ Customer
func CustomerLogin(c *gin.Context) {
	var input models.CustomerLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ຊອກຫາ customer
	var customer models.Customer
	if err := database.DB.Where("email = ?", input.Email).First(&customer).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email ຫຼື password ບໍ່ຖືກຕ້ອງ"})
		return
	}

	// ກວດສອບ password
	if !utils.CheckPasswordHash(input.Password, customer.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email ຫຼື password ບໍ່ຖືກຕ້ອງ"})
		return
	}

	// ສ້າງ token
	token, err := utils.GenerateToken(customer.ID, customer.Email, "customer")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ບໍ່ສາມາດສ້າງ token ໄດ້"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ເຂົ້າສູ່ລະບົບສຳເລັດ",
		"customer": gin.H{
			"id":    customer.ID,
			"name":  customer.Name,
			"email": customer.Email,
			"phone": customer.Phone,
		},
		"token": token,
	})
}
