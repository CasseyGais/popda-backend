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
	"popda_bulutangkis/internal/masterofficial"
	"popda_bulutangkis/internal/masterpelatih"
	"popda_bulutangkis/internal/modules"
	"popda_bulutangkis/internal/nomor"
	"popda_bulutangkis/internal/permissions"
	"popda_bulutangkis/internal/rolepermissions"
	"popda_bulutangkis/internal/roles"
	"popda_bulutangkis/internal/sekolah"
	"popda_bulutangkis/internal/shared/database"
	"popda_bulutangkis/internal/shared/middleware"
	"popda_bulutangkis/internal/tahap1"
	"popda_bulutangkis/internal/tahap2"
	"popda_bulutangkis/internal/tahap3"
	"popda_bulutangkis/internal/territories"
	"popda_bulutangkis/internal/transaksi"
	"popda_bulutangkis/internal/users"
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
	r.Static("/uploads", "./uploads")

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

	// Role Permissions
	rolePermissionsRepo := rolepermissions.NewRepository(db.DB)
	rolePermissionsService := rolepermissions.NewService(rolePermissionsRepo)
	rolePermissionsHandler := rolepermissions.NewHandler(rolePermissionsService)

	// Territories
	territoriesRepo := territories.NewRepository(db.DB)
	territoriesService := territories.NewService(territoriesRepo)
	territoriesHandler := territories.NewHandler(territoriesService)

	// Users
	usersRepo := users.NewRepository(db.DB)
	usersService := users.NewService(usersRepo)
	usersHandler := users.NewHandler(usersService)

	// Roles
	rolesRepo := roles.NewRepository(db.DB)
	rolesService := roles.NewService(rolesRepo)
	rolesHandler := roles.NewHandler(rolesService)

	// Permissions
	permissionsRepo := permissions.NewRepository(db.DB)
	permissionsService := permissions.NewService(permissionsRepo)
	permissionsHandler := permissions.NewHandler(permissionsService)

	// Master Pelatih
	masterpelatihRepo := masterpelatih.NewRepository(db.DB)
	masterpelatihService := masterpelatih.NewService(masterpelatihRepo)
	masterpelatihHandler := masterpelatih.NewHandler(masterpelatihService)

	// Master Official
	masterofficialRepo := masterofficial.NewRepository(db.DB)
	masterofficialService := masterofficial.NewService(masterofficialRepo)
	masterofficialHandler := masterofficial.NewHandler(masterofficialService)

	// Modules
	modulesRepo := modules.NewRepository(db.DB)
	modulesService := modules.NewService(modulesRepo)
	modulesHandler := modules.NewHandler(modulesService)

	// ===== ROUTES =====

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "🚀 POPDA Backend API is running!",
			"version": "1.0.0",
			"endpoints": gin.H{
				"login":               "POST /login",
				"logout":              "POST /logout (requires authentication)",
				"admin_routes":        "/admin/* (requires authentication)",
				"master_data":         "/admin/master/* (CRUD for cabor, nomor, sekolah)",
				"atlet":               "/admin/atlet/* (CRUD for atlet with verification)",
				"kontingen_identitas": "/admin/kontingen-identitas/* (CRUD for kontingen identitas)",
			},
		})
	})

	// Public
	r.POST("/login", authHandler.Login)
	r.GET("/territories", territoriesHandler.GetAll)

	// Protected
	r.POST("/logout", middleware.AuthRequired(), authHandler.Logout)

	admin := r.Group("/admin")
	admin.Use(middleware.AuthRequired())
	{
		admin.GET("/identitas", kontingenHandler.GetIdentitas)
		admin.PUT("/identitas", kontingenHandler.Update)

		// Kontingen CRUD routes
		admin.GET("/kontingen/:id", kontingenHandler.GetKontingen)
		admin.GET("/kontingen/territory/:territory_id", kontingenHandler.GetKontingenByTerritory)
		admin.POST("/kontingen", kontingenHandler.CreateKontingen)
		admin.PUT("/kontingen/:id", kontingenHandler.UpdateKontingen)

		admin.GET("/tahap1", tahap1Handler.Get)
		admin.PUT("/tahap1", tahap1Handler.Update)
		admin.DELETE("/tahap1/:cabor_id", tahap1Handler.DeleteCabor)
		admin.POST("/tahap1/submit", tahap1Handler.Submit)
		admin.GET("/tahap1/export/pdf", tahap1Handler.ExportPDF)
		admin.GET("/tahap1/export/excel", tahap1Handler.ExportExcel)

		admin.GET("/tahap2", tahap2Handler.Get)
		admin.POST("/tahap2/nomor/:nomor_id", tahap2Handler.DaftarNomor)
		admin.DELETE("/tahap2/nomor/:nomor_id", tahap2Handler.BatalNomor)
		admin.POST("/tahap2/submit", tahap2Handler.Submit)
		admin.GET("/tahap2/export/pdf", tahap2Handler.ExportPDF)
		admin.GET("/tahap2/export/excel", tahap2Handler.ExportExcel)

		admin.GET("/tahap3", tahap3Handler.Get)
		admin.POST("/tahap3/submit", tahap3Handler.Submit)
		admin.GET("/tahap3/export/pdf", tahap3Handler.ExportPDF)
		admin.GET("/tahap3/export/excel", tahap3Handler.ExportExcel)
		admin.GET("/tahap3/export/atlet/pdf", tahap3Handler.ExportAtletPDF)
		admin.GET("/tahap3/export/atlet/excel", tahap3Handler.ExportAtletExcel)
		admin.GET("/tahap3/export/pelatih/pdf", tahap3Handler.ExportPelatihPDF)
		admin.GET("/tahap3/export/pelatih/excel", tahap3Handler.ExportPelatihExcel)
		admin.GET("/tahap3/export/official/pdf", tahap3Handler.ExportOfficialPDF)
		admin.GET("/tahap3/export/official/excel", tahap3Handler.ExportOfficialExcel)

		// Tahap 3 — Atlet
		admin.GET("/tahap3/atlet", tahap3Handler.GetAtlets)
		admin.GET("/tahap3/atlet/:id", tahap3Handler.GetAtletByID)
		admin.POST("/tahap3/atlet", tahap3Handler.CreateAtlet)
		admin.PUT("/tahap3/atlet/:id", tahap3Handler.UpdateAtlet)
		admin.DELETE("/tahap3/atlet/:id", tahap3Handler.DeleteAtlet)
		admin.PUT("/tahap3/atlet/:id/foto", tahap3Handler.UploadAtletFoto)
		admin.PUT("/tahap3/atlet/:id/file/:kolom", tahap3Handler.UploadAtletFile)

		// Tahap 3 — Pelatih
		admin.GET("/tahap3/pelatih", tahap3Handler.GetPelatihs)
		admin.GET("/tahap3/pelatih/:id", tahap3Handler.GetPelatihByID)
		admin.POST("/tahap3/pelatih", tahap3Handler.CreatePelatih)
		admin.PUT("/tahap3/pelatih/:id", tahap3Handler.UpdatePelatih)
		admin.DELETE("/tahap3/pelatih/:id", tahap3Handler.DeletePelatih)
		admin.PUT("/tahap3/pelatih/:id/file/:kolom", tahap3Handler.UploadPelatihFile)

		// Tahap 3 — Official
		admin.GET("/tahap3/official", tahap3Handler.GetOfficials)
		admin.GET("/tahap3/official/:id", tahap3Handler.GetOfficialByID)
		admin.POST("/tahap3/official", tahap3Handler.CreateOfficial)
		admin.PUT("/tahap3/official/:id", tahap3Handler.UpdateOfficial)
		admin.DELETE("/tahap3/official/:id", tahap3Handler.DeleteOfficial)
		admin.PUT("/tahap3/official/:id/file/:kolom", tahap3Handler.UploadOfficialFile)

		// Tahap 3 — Transaksi Pendaftaran
		admin.POST("/tahap3/trx/atlet", tahap3Handler.CreateTrxAtlet)
		admin.DELETE("/tahap3/trx/atlet/:id", tahap3Handler.DeleteTrxAtlet)
		admin.POST("/tahap3/trx/pelatih", tahap3Handler.CreateTrxPelatih)
		admin.DELETE("/tahap3/trx/pelatih/:id", tahap3Handler.DeleteTrxPelatih)
		admin.POST("/tahap3/trx/official", tahap3Handler.CreateTrxOfficial)
		admin.DELETE("/tahap3/trx/official/:id", tahap3Handler.DeleteTrxOfficial)

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

		// Role Permissions Routes
		admin.GET("/role-permissions", rolePermissionsHandler.GetAll)
		admin.GET("/role-permissions/role/:id", rolePermissionsHandler.GetByRoleID)
		admin.GET("/role-permissions/permission/:id", rolePermissionsHandler.GetByPermissionID)
		admin.POST("/role-permissions", rolePermissionsHandler.AssignPermissionToRole)
		admin.DELETE("/role-permissions/role/:id/permission/:permissionId", rolePermissionsHandler.RemovePermissionFromRole)
		admin.DELETE("/role-permissions/role/:id", rolePermissionsHandler.DeleteByRoleID)
		admin.DELETE("/role-permissions/permission/:id", rolePermissionsHandler.DeleteByPermissionID)

		// Territories Routes
		admin.GET("/territories", territoriesHandler.GetAll)
		admin.GET("/territories/:id", territoriesHandler.GetByID)
		admin.GET("/territories/type/:type", territoriesHandler.GetByType)
		admin.GET("/territories/provinces", territoriesHandler.GetProvinces)
		admin.GET("/territories/kabupatens", territoriesHandler.GetKabupatens)
		admin.GET("/territories/kotas", territoriesHandler.GetKotas)
		admin.GET("/territories/user/:user_id", territoriesHandler.GetByUserID)
		admin.POST("/territories", territoriesHandler.Create)
		admin.PUT("/territories/:id", territoriesHandler.Update)
		admin.DELETE("/territories/:id", territoriesHandler.Delete)
		admin.POST("/territories/user/:user_id/:territory_id", territoriesHandler.AssignToUser)
		admin.DELETE("/territories/user/:user_id/:territory_id", territoriesHandler.RemoveFromUser)

		// Users Routes
		admin.GET("/users", usersHandler.GetAll)
		admin.GET("/users/:id", usersHandler.GetByID)
		admin.GET("/users/email/:email", usersHandler.GetByEmail)
		admin.POST("/users", usersHandler.Create)
		admin.PUT("/users/:id", usersHandler.Update)
		admin.DELETE("/users/:id", usersHandler.Delete)
		admin.PUT("/users/:id/avatar", usersHandler.UpdateAvatar)
		admin.PUT("/users/:id/password", usersHandler.UpdatePassword)
		admin.PUT("/users/:id/status", usersHandler.UpdateStatus)
		admin.GET("/users/:id/roles", usersHandler.GetRoles)
		admin.POST("/users/:id/roles/:role_id", usersHandler.AssignRole)
		admin.DELETE("/users/:id/roles/:role_id", usersHandler.RemoveRole)
		admin.GET("/users/:id/territories", usersHandler.GetTerritories)
		admin.POST("/users/:id/territories/:territory_id", usersHandler.AssignTerritory)
		admin.DELETE("/users/:id/territories/:territory_id", usersHandler.RemoveTerritory)

		// Roles Routes
		admin.GET("/roles", rolesHandler.GetAll)
		admin.GET("/roles/:id", rolesHandler.GetByID)
		admin.GET("/roles/user/:user_id", rolesHandler.GetByUserID)
		admin.POST("/roles", rolesHandler.Create)
		admin.PUT("/roles/:id", rolesHandler.Update)
		admin.DELETE("/roles/:id", rolesHandler.Delete)
		admin.POST("/roles/:id/permissions/:permission_id", rolesHandler.AssignPermission)
		admin.DELETE("/roles/:id/permissions/:permission_id", rolesHandler.RemovePermission)
		admin.GET("/roles/:id/permissions", rolesHandler.GetPermissions)

		// Permissions Routes
		admin.GET("/permissions", permissionsHandler.GetAll)
		admin.GET("/permissions/:id", permissionsHandler.GetByID)
		admin.GET("/permissions/role/:role_id", permissionsHandler.GetByRoleID)
		admin.GET("/permissions/module/:module_id", permissionsHandler.GetByModuleID)
		admin.POST("/permissions", permissionsHandler.Create)
		admin.PUT("/permissions/:id", permissionsHandler.Update)
		admin.DELETE("/permissions/:id", permissionsHandler.Delete)

		// Modules Routes
		admin.GET("/modules", modulesHandler.GetAll)
		admin.GET("/modules/:id", modulesHandler.GetByID)
		admin.POST("/modules", modulesHandler.Create)
		admin.PUT("/modules/:id", modulesHandler.Update)
		admin.DELETE("/modules/:id", modulesHandler.Delete)

		// Master Pelatih Routes
		admin.GET("/master/pelatih", masterpelatihHandler.GetAll)
		admin.GET("/master/pelatih/:id", masterpelatihHandler.GetByID)
		admin.GET("/master/pelatih/kontingen/:kontingen_id", masterpelatihHandler.GetByKontingenID)
		admin.POST("/master/pelatih", masterpelatihHandler.Create)
		admin.PUT("/master/pelatih/:id", masterpelatihHandler.Update)
		admin.DELETE("/master/pelatih/:id", masterpelatihHandler.Delete)
		admin.PUT("/master/pelatih/:id/foto", masterpelatihHandler.UpdateFoto)
		admin.PUT("/master/pelatih/:id/sertifikat", masterpelatihHandler.UpdateSertifikat)

		// Master Official Routes
		admin.GET("/master/official", masterofficialHandler.GetAll)
		admin.GET("/master/official/:id", masterofficialHandler.GetByID)
		admin.GET("/master/official/kontingen/:kontingen_id", masterofficialHandler.GetByKontingenID)
		admin.POST("/master/official", masterofficialHandler.Create)
		admin.PUT("/master/official/:id", masterofficialHandler.Update)
		admin.DELETE("/master/official/:id", masterofficialHandler.Delete)
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
