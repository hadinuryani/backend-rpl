package routes

import (
	"ic-plus-backend/handlers"
	"ic-plus-backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine,
	authH *handlers.AuthHandler,
	pasienH *handlers.PasienHandler,
	bidanH *handlers.BidanHandler,
	klinikH *handlers.KlinikHandler,
	antrianH *handlers.AntrianHandler,
	rmH *handlers.RekamMedisHandler,
	resepH *handlers.ResepHandler,
	jadwalH *handlers.JadwalHandler,
	notifH *handlers.NotifikasiHandler,
	invH *handlers.InventoriHandler,
	dashH *handlers.DashboardHandler,
	monitorH *handlers.MonitorHandler,
) {
	api := r.Group("/api/v1")

	// === AUTH (public) ===
	auth := api.Group("/auth")
	auth.POST("/register", authH.Register)
	auth.POST("/login", authH.Login)
	auth.POST("/logout", middleware.AuthMiddleware(), authH.Logout)
	auth.GET("/me", middleware.AuthMiddleware(), authH.GetMe)

	// === PUBLIC ===
	api.GET("/klinik/status", klinikH.GetStatus)

	// === PASIEN routes ===
	pasien := api.Group("/pasien")
	pasien.Use(middleware.AuthMiddleware(), middleware.RequireRole("pasien"))
	pasien.GET("/profile", pasienH.GetProfile)
	pasien.PUT("/profile", pasienH.UpdateProfile)
	pasien.POST("/antrian", antrianH.CreateAntrian)
	pasien.GET("/antrian", antrianH.GetMyAntrian)
	pasien.GET("/antrian/:id", antrianH.GetAntrianDetail)
	pasien.GET("/rekam-medis", rmH.GetMyRekamMedis)
	pasien.GET("/rekam-medis/:id", rmH.GetByID)
	pasien.GET("/notifikasi", notifH.GetMyNotifikasi)
	pasien.PUT("/notifikasi/:id/read", notifH.MarkAsRead)
	pasien.GET("/jadwal-kontrol", jadwalH.GetMyJadwal)

	// === BIDAN routes ===
	bidan := api.Group("/bidan")
	bidan.Use(middleware.AuthMiddleware(), middleware.RequireRole("bidan"))

	// Klinik
	bidan.PUT("/klinik/status", klinikH.SetStatus)

	// Dashboard
	bidan.GET("/dashboard", dashH.GetStats)

	// Antrian
	bidan.GET("/antrian", antrianH.GetTodayAntrian)
	bidan.GET("/antrian/:id", antrianH.GetAntrianDetail)
	bidan.PUT("/antrian/:id/batal", antrianH.CancelAntrian)

	// Rekam Medis
	bidan.POST("/rekam-medis", rmH.Create)
	bidan.GET("/rekam-medis/:id", rmH.GetByID)
	bidan.PUT("/rekam-medis/:id", rmH.Update)

	// Resep
	bidan.GET("/resep/:rekam_medis_id", resepH.GetByRekamMedisID)

	// Jadwal Kontrol
	bidan.POST("/jadwal-kontrol", jadwalH.Create)
	bidan.GET("/jadwal-kontrol", jadwalH.GetAll)
	bidan.PUT("/jadwal-kontrol/:id", jadwalH.Update)
	bidan.DELETE("/jadwal-kontrol/:id", jadwalH.Delete)

	// Monitor Kunjungan
	bidan.GET("/monitor-kunjungan", monitorH.GetVisitHistory)

	// Data Pasien
	bidan.GET("/pasien", bidanH.GetAllPasien)
	bidan.GET("/pasien/:id", bidanH.GetPasienByID)
	bidan.POST("/pasien", bidanH.CreatePasien)
	bidan.PUT("/pasien/:id", bidanH.UpdatePasien)

	// Inventori Obat
	bidan.GET("/obat", invH.GetAllObat)
	bidan.POST("/obat", invH.CreateObat)
	bidan.PUT("/obat/:id", invH.UpdateObat)
	bidan.DELETE("/obat/:id", invH.DeleteObat)
	bidan.GET("/inventori", invH.GetInventori)
	bidan.POST("/inventori/stok-masuk", invH.StokMasuk)
	bidan.POST("/inventori/stok-keluar", invH.StokKeluar)
	bidan.GET("/inventori/riwayat", invH.GetRiwayat)
}
