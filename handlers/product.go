package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
	// Support both JSON and multipart/form-data
	contentType := c.Request.Header.Get("Content-Type")
	var p models.Product
	if strings.Contains(strings.ToLower(contentType), "multipart/form-data") {
		name := c.PostForm("name")
		priceStr := c.PostForm("price")
		if name == "" || priceStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name and price are required"})
			return
		}
		price, err := strconv.Atoi(priceStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "price must be a number"})
			return
		}
		var categoryID *uint
		if cid := c.PostForm("category_id"); cid != "" {
			if v, err := strconv.Atoi(cid); err == nil {
				vv := uint(v)
				categoryID = &vv
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "category_id must be a number"})
				return
			}
		}

		// Handle image upload
		file, _ := c.FormFile("image")
		var imagePath *string
		if file != nil {
			savedPath, err := saveUploadedFile(c, file)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			imagePath = &savedPath
		}

		p = models.Product{Name: name, Price: price, CategoryID: categoryID, Image: imagePath}
	} else {
		if err := c.BindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
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

	if strings.Contains(strings.ToLower(c.Request.Header.Get("Content-Type")), "multipart/form-data") {
		updates := map[string]interface{}{}
		if v := c.PostForm("name"); v != "" {
			updates["name"] = v
		}
		if v := c.PostForm("price"); v != "" {
			if iv, err := strconv.Atoi(v); err == nil {
				updates["price"] = iv
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "price must be a number"})
				return
			}
		}
		if v := c.PostForm("category_id"); v != "" {
			if iv, err := strconv.Atoi(v); err == nil {
				updates["category_id"] = uint(iv)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "category_id must be a number"})
				return
			}
		}
		// image file
		file, _ := c.FormFile("image")
		if file != nil {
			savedPath, err := saveUploadedFile(c, file)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			updates["image"] = savedPath
		}
		if len(updates) > 0 {
			if err := database.DB.Model(&p).Select("name", "price", "category_id", "image").Updates(updates).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	} else {
		// JSON update (backward compatible)
		type updateProductInput struct {
			Name       *string `json:"name"`
			Price      *int    `json:"price"`
			CategoryID *uint   `json:"category_id"`
			Image      *string `json:"image"`
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
		if in.CategoryID != nil {
			updates["category_id"] = *in.CategoryID
		}
		if in.Image != nil {
			updates["image"] = *in.Image
		}
		if len(updates) > 0 {
			if err := database.DB.Model(&p).Select("name", "price", "category_id", "image").Updates(updates).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	// Reload with Category for response
	if err := database.DB.Preload("Category").First(&p, p.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

// saveUploadedFile stores the file under ./uploads and returns the relative path
func saveUploadedFile(c *gin.Context, file *multipart.FileHeader) (string, error) {
	// Ensure uploads directory exists
	if err := os.MkdirAll("uploads", 0755); err != nil {
		return "", err
	}
	// Keep original extension, generate unique name
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join("uploads", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}
	// Expose via /uploads in main.go
	return "/uploads/" + filename, nil
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
