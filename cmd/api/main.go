package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"popda_bulutangkis/internal/auth"
	"popda_bulutangkis/internal/kontingen"
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

	// ===== ROUTES =====

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
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Println("🚀 Server running on port:", port)
	log.Fatal(r.Run(":" + port))
}
