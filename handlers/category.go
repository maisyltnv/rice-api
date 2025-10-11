package handlers

import (
	"net/http"

	"example.com/go-xampp-api/database"
	"example.com/go-xampp-api/models"
	"github.com/gin-gonic/gin"
)

// ເບິ່ງ categories ທັງໝົດ
func GetCategories(c *gin.Context) {
	var items []models.Category
	if err := database.DB.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// ເບິ່ງ category ດຽວ
func GetCategory(c *gin.Context) {
	var cat models.Category
	if err := database.DB.First(&cat, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}
	c.JSON(http.StatusOK, cat)
}

// ສ້າງ category ໃໝ່
func CreateCategory(c *gin.Context) {
	var cat models.Category
	if err := c.BindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Create(&cat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

// ແກ້ໄຂ category
func UpdateCategory(c *gin.Context) {
	var cat models.Category
	if err := database.DB.First(&cat, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}
	var input models.Category
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat.Name = input.Name
	cat.Description = input.Description
	if err := database.DB.Save(&cat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cat)
}

// ລົບ category
func DeleteCategory(c *gin.Context) {
	result := database.DB.Delete(&models.Category{}, c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
