package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"popda_bulutangkis/internal/atlet"
	"popda_bulutangkis/internal/auth"
	"popda_bulutangkis/internal/cabor"
	"popda_bulutangkis/internal/kontingen"
	"popda_bulutangkis/internal/kontingenidentitas"
	"popda_bulutangkis/internal/nomor"
	"popda_bulutangkis/internal/sekolah"
	"popda_bulutangkis/internal/shared/database"
	"popda_bulutangkis/internal/shared/middleware"
	"popda_bulutangkis/internal/tahap1"
	"popda_bulutangkis/internal/tahap2"
	"popda_bulutangkis/internal/tahap3"
	"popda_bulutangkis/internal/transaksi"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db := database.Init()

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// ===== CORS =====
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://localhost:8000",
			"http://127.0.0.1:8000",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ===== STATIC FILES =====
	r.Static("/avatar", "./avatar")

	// ===== Dependency Injection =====
	authRepo := auth.NewRepository(db.DB)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	kontingenRepo := kontingen.NewRepository(db.DB)
	kontingenService := kontingen.NewService(kontingenRepo)
	kontingenHandler := kontingen.NewHandler(kontingenService)

	tahap1Repo := tahap1.NewRepository(db.DB)
	tahap1Service := tahap1.NewService(tahap1Repo)
	tahap1Handler := tahap1.NewHandler(tahap1Service)

	tahap2Repo := tahap2.NewRepository(db.DB)
	tahap2Service := tahap2.NewService(tahap2Repo)
	tahap2Handler := tahap2.NewHandler(tahap2Service)

	tahap3Repo := tahap3.NewRepository(db.DB)
	tahap3Service := tahap3.NewService(tahap3Repo)
	tahap3Handler := tahap3.NewHandler(tahap3Service)

	transaksiRepo := transaksi.NewRepository(db.DB)
	transaksiService := transaksi.NewService(transaksiRepo)
	transaksiHandler := transaksi.NewHandler(transaksiService)

	// Master Data
	caborRepo := cabor.NewRepository(db.DB)
	caborService := cabor.NewService(caborRepo)
	caborHandler := cabor.NewHandler(caborService)

	nomorRepo := nomor.NewRepository(db.DB)
	nomorService := nomor.NewService(nomorRepo)
	nomorHandler := nomor.NewHandler(nomorService)

	sekolahRepo := sekolah.NewRepository(db.DB)
	sekolahService := sekolah.NewService(sekolahRepo)
	sekolahHandler := sekolah.NewHandler(sekolahService)

	// Atlet
	atletRepo := atlet.NewRepository(db.DB)
	atletService := atlet.NewService(atletRepo)
	atletHandler := atlet.NewHandler(atletService)

	// Kontingen Identitas
	kontingenIdentitasRepo := kontingenidentitas.NewRepository(db.DB)
	kontingenIdentitasService := kontingenidentitas.NewService(kontingenIdentitasRepo)
	kontingenIdentitasHandler := kontingenidentitas.NewHandler(kontingenIdentitasService)

	// ===== ROUTES =====

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "🚀 POPDA Backend API is running!",
			"version": "1.0.0",
			"endpoints": gin.H{
				"login":               "POST /login",
				"admin_routes":        "/admin/* (requires authentication)",
				"master_data":         "/admin/master/* (CRUD for cabor, nomor, sekolah)",
				"atlet":               "/admin/atlet/* (CRUD for atlet with verification)",
				"kontingen_identitas": "/admin/kontingen-identitas/* (CRUD for kontingen identitas)",
			},
		})
	})

	// Public
	r.POST("/login", authHandler.Login)

	// Protected
	admin := r.Group("/admin")
	admin.Use(middleware.AuthRequired())
	{
		admin.GET("/identitas", kontingenHandler.GetIdentitas)
		admin.PUT("/identitas", kontingenHandler.Update)

		admin.GET("/tahap1", tahap1Handler.Get)
		admin.PUT("/tahap1", tahap1Handler.Update)
		admin.POST("/tahap1/submit", tahap1Handler.Submit)

		admin.GET("/tahap2", tahap2Handler.Get)
		admin.PUT("/tahap2", tahap2Handler.Update)
		admin.POST("/tahap2/submit", tahap2Handler.Submit)

		admin.GET("/tahap3", tahap3Handler.Get)
		admin.POST("/tahap3/submit", tahap3Handler.Submit)

		admin.POST("/trx/cabor", transaksiHandler.CreateTrxKontingenCabor)
		admin.GET("/trx/cabor", transaksiHandler.GetTrxKontingenCabor)
		admin.PUT("/trx/cabor", transaksiHandler.UpdateTrxKontingenCabor)

		admin.POST("/trx/nomor", transaksiHandler.CreateTrxKontingenNomor)
		admin.GET("/trx/nomor", transaksiHandler.GetTrxKontingenNomor)
		admin.DELETE("/trx/nomor/:nomor_id", transaksiHandler.DeleteTrxKontingenNomor)

		admin.POST("/trx/atlet", transaksiHandler.CreateTrxPendaftaranAtlet)
		admin.GET("/trx/atlet", transaksiHandler.GetTrxPendaftaranAtlet)
		admin.PUT("/trx/atlet/:atlet_id/:nomor_id", transaksiHandler.UpdateTrxPendaftaranAtlet)

		// Master Data Routes
		admin.GET("/master/cabor", caborHandler.GetAll)
		admin.GET("/master/cabor/:id", caborHandler.GetByID)
		admin.POST("/master/cabor", caborHandler.Create)
		admin.PUT("/master/cabor/:id", caborHandler.Update)
		admin.DELETE("/master/cabor/:id", caborHandler.Delete)

		admin.GET("/master/nomor", nomorHandler.GetAll)
		admin.GET("/master/nomor/:id", nomorHandler.GetByID)
		admin.GET("/master/nomor/cabor/:cabor_id", nomorHandler.GetByCaborID)
		admin.POST("/master/nomor", nomorHandler.Create)
		admin.PUT("/master/nomor/:id", nomorHandler.Update)
		admin.DELETE("/master/nomor/:id", nomorHandler.Delete)

		admin.GET("/master/sekolah", sekolahHandler.GetAll)
		admin.GET("/master/sekolah/:id", sekolahHandler.GetByID)
		admin.GET("/master/sekolah/search", sekolahHandler.Search)
		admin.POST("/master/sekolah", sekolahHandler.Create)
		admin.PUT("/master/sekolah/:id", sekolahHandler.Update)
		admin.DELETE("/master/sekolah/:id", sekolahHandler.Delete)

		// Atlet Routes
		admin.GET("/atlet", atletHandler.GetAll)
		admin.GET("/atlet/:id", atletHandler.GetByID)
		admin.GET("/atlet/kontingen/:kontingen_id", atletHandler.GetByKontingenID)
		admin.GET("/atlet/sekolah/:sekolah_id", atletHandler.GetBySekolahID)
		admin.GET("/atlet/status/:status", atletHandler.GetByStatus)
		admin.POST("/atlet", atletHandler.Create)
		admin.PUT("/atlet/:id", atletHandler.Update)
		admin.DELETE("/atlet/:id", atletHandler.Delete)
		admin.PUT("/atlet/:id/status", atletHandler.UpdateStatus)
		admin.PUT("/atlet/:id/foto", atletHandler.UpdateFoto)

		// Kontingen Identitas Routes
		admin.GET("/kontingen-identitas", kontingenIdentitasHandler.GetAll)
		admin.GET("/kontingen-identitas/:id", kontingenIdentitasHandler.GetByID)
		admin.GET("/kontingen-identitas/kontingen/:kontingen_id", kontingenIdentitasHandler.GetByKontingenID)
		admin.POST("/kontingen-identitas", kontingenIdentitasHandler.Create)
		admin.PUT("/kontingen-identitas/:id", kontingenIdentitasHandler.Update)
		admin.DELETE("/kontingen-identitas/:id", kontingenIdentitasHandler.Delete)
		admin.PUT("/kontingen-identitas/:id/kepala-foto", kontingenIdentitasHandler.UpdateKepalaFoto)
		admin.PUT("/kontingen-identitas/:id/pic-foto", kontingenIdentitasHandler.UpdatePICFoto)
	}

	log.Println("🚀 Starting server...")
	dsn := os.Getenv("DB_DSN")
	log.Println("🔗 Database DSN:", dsn)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Println("🚀 Server running on port:", port)
	log.Fatal(r.Run(":" + port))
}
