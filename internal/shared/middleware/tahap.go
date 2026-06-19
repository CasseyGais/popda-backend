package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TahapOpen memblokir operasi tulis ke endpoint tahap yang sedang ditutup.
// Superadmin dikecualikan — superadmin tetap bisa akses meskipun tahap tutup,
// agar bisa lihat data dan mengatur kembali pengaturan.
//
// Contoh pemakaian di route:
//
//	admin.PUT("/tahap1", middleware.TahapOpen(db, 1), tahap1Handler.Update)
//	admin.POST("/tahap1/submit", middleware.TahapOpen(db, 1), tahap1Handler.Submit)
func TahapOpen(db *gorm.DB, tahap int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Superadmin bypass — bisa akses meskipun tahap tutup
		role, _ := c.Get("role")
		if strings.ToLower(role.(string)) == "superadmin" {
			c.Next()
			return
		}

		// Cek is_open di tabel pengaturan_tahap
		var isOpen bool
		err := db.Table("pengaturan_tahap").
			Select("is_open").
			Where("tahap = ?", tahap).
			Scan(&isOpen).Error

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal mengecek status tahap",
			})
			return
		}

		if !isOpen {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success":    false,
				"message":    "Tahap " + itoa(tahap) + " belum dibuka. Silakan hubungi panitia.",
				"error_code": "TAHAP_CLOSED",
			})
			return
		}

		c.Next()
	}
}

// itoa konversi int ke string tanpa import strconv yang redundan di package ini
func itoa(i int) string {
	switch i {
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	default:
		return "?"
	}
}
