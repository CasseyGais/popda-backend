// internal/middleware/auth.go
package middleware

import (
	"net/http"
	"strings"

	"popda_bulutangkis/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Header Authorization diperlukan"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Format header authorization tidak valid"})
			return
		}

		tokenStr := parts[1]

		claims, err := jwt.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau kadaluarsa"})
			return
		}

		// TODO: Ambil kontingenID dari user_territories berdasarkan userID
		// Untuk sementara, gunakan userID sebagai kontingenID
		// claims.KontingenID = claims.UserID

		// Simpan claims ke context
		c.Set("user", claims)

		c.Next()
	}
}
