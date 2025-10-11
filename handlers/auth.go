package handlers

import (
	"net/http"

	"example.com/go-xampp-api/database"
	"example.com/go-xampp-api/models"
	"example.com/go-xampp-api/utils"
	"github.com/gin-gonic/gin"
)

// ການລົງທະບຽນ
func Register(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ກວດສອບວ່າມີ username ຫຼື email ຊ້ຳບໍ່
	var existingUser models.User
	if err := database.DB.Where("username = ? OR email = ?", input.Username, input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username ຫຼື email ຖືກໃຊ້ແລ້ວ"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ບໍ່ສາມາດ hash password ໄດ້"})
		return
	}

	// ສ້າງຜູ້ໃຊ້ໃໝ່
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ສ້າງ token
	token, err := utils.GenerateToken(user.ID, user.Username)
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
}

// ເຂົ້າສູ່ລະບົບ
func Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ຊອກຫາຜູ້ໃຊ້
	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username ຫຼື password ບໍ່ຖືກຕ້ອງ"})
		return
	}

	// ກວດສອບ password
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username ຫຼື password ບໍ່ຖືກຕ້ອງ"})
		return
	}

	// ສ້າງ token
	token, err := utils.GenerateToken(user.ID, user.Username)
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
}
