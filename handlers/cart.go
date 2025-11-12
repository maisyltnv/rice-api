package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"example.com/go-xampp-api/database"
	"example.com/go-xampp-api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetCart returns the authenticated customer's cart.
func GetCart(c *gin.Context) {
	customerID, ok := getCustomerID(c)
	if !ok {
		return
	}

	cart, err := loadOrCreateCart(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hydrateCart(&cart)
	c.JSON(http.StatusOK, cart)
}

// AddCartItem adds or increments a product in the customer's cart.
func AddCartItem(c *gin.Context) {
	customerID, ok := getCustomerID(c)
	if !ok {
		return
	}

	var input models.AddCartItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	cart, err := findOrCreateCart(tx, customerID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := tx.First(&product, input.ProductID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var item models.CartItem
	err = tx.Where("cart_id = ? AND product_id = ?", cart.ID, input.ProductID).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			item = models.CartItem{
				CartID:       cart.ID,
				ProductID:    product.ID,
				ProductName:  product.Name,
				ProductImage: product.Image,
				UnitPrice:    product.Price,
				Quantity:     input.Quantity,
			}
			if err := tx.Create(&item).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		item.Quantity += input.Quantity
		item.ProductName = product.Name
		item.ProductImage = product.Image
		item.UnitPrice = product.Price
		if err := tx.Save(&item).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if _, err := recalcCartTotals(tx, cart.ID); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cart, err = loadCartByID(cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hydrateCart(&cart)
	c.JSON(http.StatusOK, cart)
}

// UpdateCartItem updates the quantity of a cart item.
func UpdateCartItem(c *gin.Context) {
	customerID, ok := getCustomerID(c)
	if !ok {
		return
	}

	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cart item id"})
		return
	}

	var input models.UpdateCartItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := database.DB.Begin()

	cart, err := findOrCreateCart(tx, customerID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var item models.CartItem
	if err := tx.Where("id = ? AND cart_id = ?", itemID, cart.ID).First(&item).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "cart item not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	item.Quantity = input.Quantity
	if err := tx.Save(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if _, err := recalcCartTotals(tx, cart.ID); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cart, err = loadCartByID(cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hydrateCart(&cart)
	c.JSON(http.StatusOK, cart)
}

// DeleteCartItem removes a specific item from the cart.
func DeleteCartItem(c *gin.Context) {
	customerID, ok := getCustomerID(c)
	if !ok {
		return
	}

	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cart item id"})
		return
	}

	tx := database.DB.Begin()

	cart, err := findOrCreateCart(tx, customerID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := tx.Where("id = ? AND cart_id = ?", itemID, cart.ID).Delete(&models.CartItem{})
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "cart item not found"})
		return
	}

	if _, err := recalcCartTotals(tx, cart.ID); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cart, err = loadCartByID(cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hydrateCart(&cart)
	c.JSON(http.StatusOK, cart)
}

// ClearCart removes all items in the customer's cart.
func ClearCart(c *gin.Context) {
	customerID, ok := getCustomerID(c)
	if !ok {
		return
	}

	tx := database.DB.Begin()

	cart, err := findOrCreateCart(tx, customerID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if _, err := recalcCartTotals(tx, cart.ID); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cart, err = loadCartByID(cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hydrateCart(&cart)
	c.JSON(http.StatusOK, cart)
}

func getCustomerID(c *gin.Context) (uint, bool) {
	role, ok := c.Get("role")
	if !ok || role != "customer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "customer access required"})
		return 0, false
	}

	if cid, ok := c.Get("customer_id"); ok {
		switch v := cid.(type) {
		case uint:
			return v, true
		case int:
			return uint(v), true
		case int64:
			return uint(v), true
		case float64:
			return uint(v), true
		}
	}

	if uid, ok := c.Get("user_id"); ok {
		switch v := uid.(type) {
		case uint:
			return v, true
		case int:
			return uint(v), true
		case int64:
			return uint(v), true
		case float64:
			return uint(v), true
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "customer identification failed"})
	return 0, false
}

func findOrCreateCart(tx *gorm.DB, customerID uint) (models.Cart, error) {
	var cart models.Cart
	err := tx.Where("customer_id = ?", customerID).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart = models.Cart{
				CustomerID: customerID,
			}
			if err := tx.Create(&cart).Error; err != nil {
				return models.Cart{}, err
			}
		} else {
			return models.Cart{}, err
		}
	}
	return cart, nil
}

func recalcCartTotals(tx *gorm.DB, cartID uint) (int, error) {
	var result struct {
		Total int
	}
	if err := tx.Model(&models.CartItem{}).
		Select("COALESCE(SUM(unit_price * quantity), 0) AS total").
		Where("cart_id = ?", cartID).
		Scan(&result).Error; err != nil {
		return 0, err
	}

	if err := tx.Model(&models.Cart{}).
		Where("id = ?", cartID).
		Update("total_amount", result.Total).Error; err != nil {
		return 0, err
	}

	return result.Total, nil
}

func loadOrCreateCart(customerID uint) (models.Cart, error) {
	var cart models.Cart
	err := database.DB.Preload("Items.Product").
		Where("customer_id = ?", customerID).
		First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart = models.Cart{
				CustomerID: customerID,
			}
			if err := database.DB.Create(&cart).Error; err != nil {
				return models.Cart{}, err
			}
			if err := database.DB.Preload("Items.Product").
				First(&cart, cart.ID).Error; err != nil {
				return models.Cart{}, err
			}
		} else {
			return models.Cart{}, err
		}
	}
	return cart, nil
}

func loadCartByID(cartID uint) (models.Cart, error) {
	var cart models.Cart
	if err := database.DB.Preload("Items.Product").First(&cart, cartID).Error; err != nil {
		return models.Cart{}, err
	}
	return cart, nil
}

func hydrateCart(cart *models.Cart) {
	total := 0
	for i := range cart.Items {
		subtotal := cart.Items[i].UnitPrice * cart.Items[i].Quantity
		cart.Items[i].Subtotal = subtotal
		total += subtotal
	}
	cart.TotalAmount = total
}
