package middleware

import (
	"net/http"
	"strings"

	"popda_bulutangkis/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthRequired memvalidasi JWT token dan menyimpan claims ke context.
// Semua protected route wajib pakai middleware ini.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Header Authorization diperlukan",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Format header authorization tidak valid",
			})
			return
		}

		tokenStr := parts[1]

		claims, err := jwt.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token tidak valid atau kadaluarsa",
			})
			return
		}

		c.Set("user", claims)
		c.Set("user_id", claims.UserID)
		c.Set("kontingen_id", claims.KontingenID)
		c.Set("role", claims.Role)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// PermissionRequired mengecek apakah user punya permission tertentu.
// Superadmin selalu lolos tanpa perlu assign permission satu-satu.
// User lain dicek via tabel role_permissions → permissions.
//
// Contoh pemakaian di route:
//
//	admin.GET("/master/cabor", middleware.PermissionRequired(db, "cabor.read"), caborHandler.GetAll)
func PermissionRequired(db *gorm.DB, permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil role dari context (sudah di-set oleh AuthRequired)
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Role tidak ditemukan di token",
			})
			return
		}

		// Superadmin bypass semua permission check
		if role.(string) == "superadmin" {
			c.Next()
			return
		}

		// Ambil user_id dari context
		userID, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "User ID tidak ditemukan di token",
			})
			return
		}

		// Cek apakah user punya permission via role_permissions
		// users → user_roles → role_permissions → permissions
		var count int64
		err := db.Table("permissions p").
			Joins("INNER JOIN role_permissions rp ON rp.permission_id = p.id").
			Joins("INNER JOIN user_roles ur ON ur.role_id = rp.role_id").
			Where("ur.user_id = ? AND p.name = ?", userID, permissionName).
			Count(&count).Error

		if err != nil || count == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak — tidak punya permission: " + permissionName,
			})
			return
		}

		c.Next()
	}
}

// SuperadminOnly hanya mengizinkan user dengan role superadmin.
// Dipakai untuk endpoint yang benar-benar hanya boleh diakses superadmin
// seperti manajemen user, roles, permissions, modules.
//
// Contoh pemakaian:
//
//	admin.GET("/users", middleware.SuperadminOnly(), usersHandler.GetAll)
func SuperadminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(string) != "superadmin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak — hanya superadmin yang diizinkan",
			})
			return
		}
		c.Next()
	}
}
