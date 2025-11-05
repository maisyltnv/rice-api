package handlers

import (
	"net/http"

	"example.com/go-xampp-api/database"
	"example.com/go-xampp-api/models"
	"github.com/gin-gonic/gin"
)

// ເບິ່ງ products ທັງໝົດ
func GetProducts(c *gin.Context) {
	var items []models.Product
	// Preload Category ເພື່ອສະແດງຂໍ້ມູນ category ພ້ອມກັບ product
	if err := database.DB.Preload("Category").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// ເບິ່ງ product ດຽວ
func GetProduct(c *gin.Context) {
	var p models.Product
	// Preload Category ເພື່ອສະແດງຂໍ້ມູນ category
	if err := database.DB.Preload("Category").First(&p, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

// ສ້າງ product ໃໝ່
func CreateProduct(c *gin.Context) {
	var p models.Product
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

// ແກ້ໄຂ product
func UpdateProduct(c *gin.Context) {
	var p models.Product
	if err := database.DB.First(&p, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	// Use a dedicated input with pointers to detect which fields are present
	type updateProductInput struct {
		Name       *string `json:"name"`
		Price      *int    `json:"price"`
		CategoryID *uint   `json:"category_id"`
	}

	var in updateProductInput
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if in.Name != nil {
		updates["name"] = *in.Name
	}
	if in.Price != nil {
		updates["price"] = *in.Price
	}
	// Explicitly include category_id when present so it never gets ignored
	if in.CategoryID != nil {
		updates["category_id"] = *in.CategoryID
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&p).Select("name", "price", "category_id").Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Reload with Category for response
	if err := database.DB.Preload("Category").First(&p, p.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

// ລົບ product
func DeleteProduct(c *gin.Context) {
	result := database.DB.Delete(&models.Product{}, c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
