package middleware

import (
	"net/http"

	"example.com/go-xampp-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Middleware ສຳລັບກວດສອບ JWT Token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ດຶງ token ຈາກ Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ບໍ່ມີ token, ກະລຸນາເຂົ້າສູ່ລະບົບກ່ອນ"})
			c.Abort()
			return
		}

		// ກຳจັດ "Bearer " ອອກຈາກ token (ຖ້າມີ)
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// Validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// ກວດສອບວ່າໃຊ້ signing method ທີ່ຖືກຕ້ອງ
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return utils.JWTSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token ບໍ່ຖືກຕ້ອງຫຼືໝົດອາຍຸແລ້ວ"})
			c.Abort()
			return
		}

		// ເອົາຂໍ້ມູນຈາກ token ໃສ່ context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", uint(claims["user_id"].(float64)))
			c.Set("username", claims["username"].(string))
		}

		c.Next()
	}
}
